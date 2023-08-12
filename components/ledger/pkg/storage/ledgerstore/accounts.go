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

func (s *Store) accountQueryBuilder(q AccountsQuery) func(query *bun.SelectQuery) *bun.SelectQuery {
	return func(query *bun.SelectQuery) *bun.SelectQuery {
		query = query.
			DistinctOn("accounts.address").
			Column("accounts.address").
			ColumnExpr("coalesce(metadata, '{}'::jsonb) as metadata").
			Table("accounts").
			Join("left join accounts_metadata on accounts_metadata.address = accounts.address").
			Apply(filterMetadata(q.Options.Metadata)).
			Apply(filterAccountAddress(q.Options.Address, "accounts.address")).
			Apply(filterPIT(q.Options.PIT, "coalesce(accounts_metadata.date, accounts.insertion_date)")).
			Order("accounts.address", "revision desc")

		if q.Options.ExpandVolumes {
			query = query.
				ColumnExpr("volumes.*").
				TableExpr("get_account_aggregated_volumes(accounts.address) volumes")
		}

		return query
	}
}

func (s *Store) GetAccounts(ctx context.Context, q AccountsQuery) (*api.Cursor[core.Account], error) {
	return paginateWithOffset[AccountsQueryOptions, core.Account](s, ctx,
		paginate.OffsetPaginatedQuery[AccountsQueryOptions](q),
		s.accountQueryBuilder(q),
	)
}

type GetAccountQuery struct {
	Addr string
	PIT  core.Time
}

func (q GetAccountQuery) WithPIT(pit core.Time) GetAccountQuery {
	q.PIT = pit

	return q
}

func NewGetAccountQuery(addr string) GetAccountQuery {
	return GetAccountQuery{
		Addr: addr,
	}
}

func (s *Store) GetAccount(ctx context.Context, address string) (*core.Account, error) {
	account, err := fetch[*core.Account](s, ctx, func(query *bun.SelectQuery) *bun.SelectQuery {
		return query.
			ColumnExpr("accounts.address").
			ColumnExpr("coalesce(metadata, '{}'::jsonb) as metadata").
			Join("left join accounts_metadata on accounts_metadata.address = accounts.address").
			Where("accounts.address = ?", address).
			Order("revision desc").
			Limit(1)
	})
	if err != nil {
		if storageerrors.IsNotFoundError(err) {
			return pointer.For(core.NewAccount(address)), nil
		}
		return nil, err
	}
	return account, nil
}

// TODO: Find another name for this method
func (s *Store) GetAccountWithQuery(ctx context.Context, q GetAccountQuery) (*core.Account, error) {
	account, err := fetch[*core.Account](s, ctx, func(query *bun.SelectQuery) *bun.SelectQuery {
		return query.
			ColumnExpr("accounts.address").
			ColumnExpr("coalesce(metadata, '{}'::jsonb) as metadata").
			Where("accounts.address = ?", q.Addr).
			Join("left join accounts_metadata on accounts_metadata.address = accounts.address").
			Order("revision desc").
			Apply(filterPIT(q.PIT, "coalesce(accounts_metadata.date, accounts.insertion_date)")).
			Limit(1)
	})
	if err != nil {
		if storageerrors.IsNotFoundError(err) {
			return pointer.For(core.NewAccount(q.Addr)), nil
		}
		return nil, err
	}
	return account, nil
}

// TODO: Add PIT
func (s *Store) GetAccountWithVolumes(ctx context.Context, account string, volumes, effectiveVolumes bool) (*core.AccountWithVolumes, error) {
	return fetch[*core.AccountWithVolumes](s, ctx, func(query *bun.SelectQuery) *bun.SelectQuery {
		query = query.
			Column("accounts.address").
			ColumnExpr("coalesce(accounts_metadata.metadata, '{}'::jsonb) as metadata").
			Join("left join accounts_metadata on accounts_metadata.address = accounts.address").
			Where("accounts.address = ?", account).
			Order("revision desc").
			Limit(1)
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

type AccountsQuery paginate.OffsetPaginatedQuery[AccountsQueryOptions]

type AccountsQueryOptions struct {
	AfterAddress string            `json:"after"`
	Address      string            `json:"address"`
	Metadata     metadata.Metadata `json:"metadata"`

	PIT                    core.Time `json:"pit"`
	ExpandVolumes          bool      `json:"volumes"`
	ExpandEffectiveVolumes bool      `json:"effectiveVolumes"`
}

func NewAccountsQuery() AccountsQuery {
	return AccountsQuery{
		PageSize: paginate.QueryDefaultPageSize,
		Order:    paginate.OrderAsc,
		Options: AccountsQueryOptions{
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
	a.Options.AfterAddress = after

	return a
}

func (a AccountsQuery) WithAddressFilter(address string) AccountsQuery {
	a.Options.Address = address

	return a
}

func (a AccountsQuery) WithMetadataFilter(metadata metadata.Metadata) AccountsQuery {
	a.Options.Metadata = metadata

	return a
}

func (a AccountsQuery) WithPIT(date core.Time) AccountsQuery {
	a.Options.PIT = date

	return a
}
