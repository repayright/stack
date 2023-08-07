package ledgerstore

import (
	"context"

	"github.com/formancehq/ledger/pkg/core"
	storageerrors "github.com/formancehq/ledger/pkg/storage"
	"github.com/formancehq/ledger/pkg/storage/paginate"
	"github.com/formancehq/stack/libs/go-libs/api"
	"github.com/formancehq/stack/libs/go-libs/metadata"
	"github.com/formancehq/stack/libs/go-libs/pointer"
	"github.com/uptrace/bun"
)

func (s *Store) accountQueryBuilder(p AccountsQuery) func(query *bun.SelectQuery) *bun.SelectQuery {
	return func(query *bun.SelectQuery) *bun.SelectQuery {
		selectAccounts := s.schema.IDB.NewSelect().
			Table("accounts").
			Apply(filterMetadata(p.Filters.Metadata)).
			Apply(filterAccountAddress(p.Filters.Address, "address"))

		return query.
			With("cte1", selectAccounts).
			DistinctOn("cte1.address").
			ColumnExpr("cte1.address").
			ColumnExpr("cte1.metadata").
			Table("cte1")
	}
}

func (s *Store) GetAccounts(ctx context.Context, q AccountsQuery) (*api.Cursor[core.Account], error) {
	return paginateWithOffset[AccountsQueryFilters, core.Account](s, ctx,
		paginate.OffsetPaginatedQuery[AccountsQueryFilters](q),
		s.accountQueryBuilder(q),
	)
}

func (s *Store) GetAccount(ctx context.Context, addr string) (*core.Account, error) {
	account, err := fetch[*core.Account](s, ctx, func(query *bun.SelectQuery) *bun.SelectQuery {
		return query.
			ColumnExpr("address").
			ColumnExpr("metadata").
			Where("address = ?", addr)
	})
	if err != nil {
		if storageerrors.IsNotFoundError(err) {
			return pointer.For(core.NewAccount(addr)), nil
		}
		return nil, err
	}
	return account, nil
}

func (s *Store) GetAccountWithVolumes(ctx context.Context, account string, volumes, effectiveVolumes bool) (*core.AccountWithVolumes, error) {
	return fetch[*core.AccountWithVolumes](s, ctx, func(query *bun.SelectQuery) *bun.SelectQuery {
		query = query.
			Column("address", "metadata").
			Where("address = ?", account)
		if volumes {
			query = query.ColumnExpr("get_account_aggregated_volumes(accounts.address) as volumes")
		}
		if effectiveVolumes {
			query = query.ColumnExpr("get_account_aggregated_effective_volumes(accounts.address) as effective_volumes")
		}
		return query
	})
}

func (s *Store) CountAccounts(ctx context.Context, q AccountsQuery) (uint64, error) {
	return count(s, ctx, s.accountQueryBuilder(q))
}

type AccountsQuery paginate.OffsetPaginatedQuery[AccountsQueryFilters]

type AccountsQueryFilters struct {
	AfterAddress string            `json:"after"`
	Address      string            `json:"address"`
	Metadata     metadata.Metadata `json:"metadata"`
}

func NewAccountsQuery() AccountsQuery {
	return AccountsQuery{
		PageSize: paginate.QueryDefaultPageSize,
		Order:    paginate.OrderAsc,
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
