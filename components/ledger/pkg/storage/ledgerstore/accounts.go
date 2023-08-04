package ledgerstore

import (
	"context"
	"database/sql"
	"fmt"
	"regexp"
	"strings"

	"github.com/formancehq/ledger/pkg/core"
	storageerrors "github.com/formancehq/ledger/pkg/storage"
	"github.com/formancehq/stack/libs/go-libs/api"
	"github.com/formancehq/stack/libs/go-libs/metadata"
	"github.com/uptrace/bun"
)

const (
	accountsTableName = "accounts"
)

type AccountsQuery OffsetPaginatedQuery[AccountsQueryFilters]

type AccountsQueryFilters struct {
	AfterAddress string            `json:"after"`
	Address      string            `json:"address"`
	Metadata     metadata.Metadata `json:"metadata"`
}

func NewAccountsQuery() AccountsQuery {
	return AccountsQuery{
		PageSize: QueryDefaultPageSize,
		Order:    OrderAsc,
		Filters: AccountsQueryFilters{
			Metadata: metadata.Metadata{},
		},
	}
}

func (a AccountsQuery) WithPageSize(pageSize uint64) AccountsQuery {
	if pageSize != 0 {
		a.PageSize = pageSize
	}

	return a
}

func (a AccountsQuery) WithAfterAddress(after string) AccountsQuery {
	a.Filters.AfterAddress = after

	return a
}

func (a AccountsQuery) WithAddressFilter(address string) AccountsQuery {
	a.Filters.Address = address

	return a
}

func (a AccountsQuery) WithMetadataFilter(metadata metadata.Metadata) AccountsQuery {
	a.Filters.Metadata = metadata

	return a
}

// This regexp is used to validate the account name
// If the account name is not valid, it means that the user putted a regex in
// the address filter, and we have to change the postgres operator used.
var accountNameRegex = regexp.MustCompile(`^[a-zA-Z_0-9]+$`)

type Account struct {
	bun.BaseModel `bun:"accounts,alias:accounts"`

	Address     string            `bun:"address,type:varchar,unique,notnull"`
	Metadata    map[string]string `bun:"metadata,type:jsonb,default:'{}'"`
	AddressJson []string          `bun:"address_array,type:jsonb"`
}

func (s *Store) buildAccountsQuery(p AccountsQuery) *bun.SelectQuery {
	query := s.schema.NewSelect(accountsTableName).
		Model((*Account)(nil))

	if p.Filters.Address != "" {
		src := strings.Split(p.Filters.Address, ":")
		query.Where(fmt.Sprintf("jsonb_array_length(address_array) = %d", len(src)))

		for i, segment := range src {
			if len(segment) == 0 {
				continue
			}
			query.Where(fmt.Sprintf("address_array @@ ('$[%d] == \"' || ?::text || '\"')::jsonpath", i), segment)
		}
	}

	for key, value := range p.Filters.Metadata {
		query.Where(
			fmt.Sprintf(`"%s".%s(metadata, ?, '%s')`, s.schema.Name(),
				SQLCustomFuncMetaCompare, strings.ReplaceAll(key, ".", "', '"),
			), value)
	}

	return s.schema.IDB.NewSelect().
		With("cte1", query).
		DistinctOn("cte1.address").
		ColumnExpr("cte1.address").
		ColumnExpr("cte1.metadata").
		Table("cte1")
}

func (s *Store) GetAccounts(ctx context.Context, q AccountsQuery) (*api.Cursor[core.Account], error) {
	return UsingOffset[AccountsQueryFilters, core.Account](ctx,
		s.buildAccountsQuery(q), OffsetPaginatedQuery[AccountsQueryFilters](q))
}

func (s *Store) GetAccount(ctx context.Context, addr string) (*core.Account, error) {
	account := &core.Account{}
	if err := s.schema.NewSelect(accountsTableName).
		ColumnExpr("address").
		ColumnExpr("metadata").
		Where("address = ?", addr).
		Model(account).
		Scan(ctx, account); err != nil {
		if err == sql.ErrNoRows {
			return &core.Account{
				Address:  addr,
				Metadata: metadata.Metadata{},
			}, nil
		}
		return nil, err
	}

	return account, nil
}

func (s *Store) GetAccountWithVolumes(ctx context.Context, account string) (*core.AccountWithVolumes, error) {

	accountWithVolumes := &core.AccountWithVolumes{}
	err := s.schema.NewSelect(accountsTableName).
		Column("address", "metadata").
		ColumnExpr("get_account_aggregated_volumes(accounts.address) as volumes").
		Where("address = ?", account).
		Scan(ctx, accountWithVolumes)
	if err != nil {
		return nil, storageerrors.PostgresError(err)
	}
	return accountWithVolumes, nil
}

func (s *Store) CountAccounts(ctx context.Context, q AccountsQuery) (uint64, error) {
	sb := s.buildAccountsQuery(q)
	count, err := sb.Count(ctx)
	return uint64(count), storageerrors.PostgresError(err)
}
