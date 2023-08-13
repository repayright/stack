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

func (s *Store) accountQueryBuilder(q AccountFilter) func(query *bun.SelectQuery) *bun.SelectQuery {
	return func(query *bun.SelectQuery) *bun.SelectQuery {
		query = query.
			DistinctOn("accounts.address").
			Column("accounts.address").
			ColumnExpr("coalesce(metadata, '{}'::jsonb) as metadata").
			Table("accounts").
			Apply(filterMetadata(q.Metadata)).
			Apply(filterPIT(q.PIT, "insertion_date")).
			Order("accounts.address", "revision desc")

		if q.PIT == nil {
			query = query.Join("left join accounts_metadata on accounts_metadata.address = accounts.address")
		} else {
			query = query.Join("left join accounts_metadata on accounts_metadata.address = accounts.address and accounts_metadata.date < ?", q.PIT)
		}

		if q.ExpandVolumes {
			query = query.
				ColumnExpr("volumes.*").
				Join("join get_account_aggregated_volumes(accounts.address, ?) volumes on true", q.PIT)
		}

		if q.ExpandEffectiveVolumes {
			query = query.
				ColumnExpr("effective_volumes.*").
				Join("join get_account_aggregated_effective_volumes(accounts.address, ?) effective_volumes on true", q.PIT)
		}

		return query
	}
}

func (s *Store) GetAccountsWithVolumes(ctx context.Context, q GetAccountsQuery) (*api.Cursor[core.ExpandedAccount], error) {
	return paginateWithOffset[GetAccountsOptions, core.ExpandedAccount](s, ctx,
		paginate.OffsetPaginatedQuery[GetAccountsOptions](q),
		func(query *bun.SelectQuery) *bun.SelectQuery {
			return s.
				accountQueryBuilder(q.Options.AccountFilter)(query).
				Apply(filterAccountAddress(q.Options.Address, "accounts.address"))
		},
	)
}

type GetAccountQuery struct {
	AccountFilter
	Addr string
}

func (q GetAccountQuery) WithPIT(pit core.Time) GetAccountQuery {
	q.PIT = &pit

	return q
}

func (q GetAccountQuery) WithExpandVolumes() GetAccountQuery {
	q.ExpandVolumes = true

	return q
}

func (q GetAccountQuery) WithExpandEffectiveVolumes() GetAccountQuery {
	q.ExpandEffectiveVolumes = true

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

func (s *Store) GetAccountWithVolumes(ctx context.Context, q GetAccountQuery) (*core.ExpandedAccount, error) {
	account, err := fetch[*core.ExpandedAccount](s, ctx, func(query *bun.SelectQuery) *bun.SelectQuery {
		query = s.accountQueryBuilder(q.AccountFilter)(query).
			Where("accounts.address = ?", q.Addr).
			Limit(1)

		return query
	})
	if err != nil {
		if storageerrors.IsNotFoundError(err) {
			return pointer.For(core.NewExpandedAccount(q.Addr)), nil
		}
		return nil, err
	}
	return account, nil
}

func (s *Store) CountAccounts(ctx context.Context, q GetAccountsQuery) (uint64, error) {
	return count(s, ctx, func(query *bun.SelectQuery) *bun.SelectQuery {
		return s.
			accountQueryBuilder(q.Options.AccountFilter)(query).
			Apply(filterAccountAddress(q.Options.Address, "accounts.address"))
	})
}

type GetAccountsQuery paginate.OffsetPaginatedQuery[GetAccountsOptions]

type AccountFilter struct {
	PIT                    *core.Time        `json:"pit"`
	ExpandVolumes          bool              `json:"volumes"`
	ExpandEffectiveVolumes bool              `json:"effectiveVolumes"`
	Metadata               metadata.Metadata `json:"metadata"`
}

type GetAccountsOptions struct {
	AccountFilter
	AfterAddress string `json:"after"`
	Address      string `json:"address"`
}

func NewGetAccountsQuery() GetAccountsQuery {
	return GetAccountsQuery{
		PageSize: paginate.QueryDefaultPageSize,
		Order:    paginate.OrderAsc,
		Options: GetAccountsOptions{
			AccountFilter: AccountFilter{
				Metadata: metadata.Metadata{},
			},
		},
	}
}

func (a GetAccountsQuery) WithPageSize(pageSize uint64) GetAccountsQuery {
	if pageSize != 0 {
		a.PageSize = pageSize
	}

	return a
}

func (a GetAccountsQuery) WithAfterAddress(after string) GetAccountsQuery {
	a.Options.AfterAddress = after

	return a
}

func (a GetAccountsQuery) WithAddress(address string) GetAccountsQuery {
	a.Options.Address = address

	return a
}

func (a GetAccountsQuery) WithMetadataFilter(metadata metadata.Metadata) GetAccountsQuery {
	a.Options.Metadata = metadata

	return a
}

func (a GetAccountsQuery) WithPIT(date core.Time) GetAccountsQuery {
	a.Options.PIT = &date

	return a
}

func (a GetAccountsQuery) WithExpandVolumes() GetAccountsQuery {
	a.Options.ExpandVolumes = true

	return a
}

func (a GetAccountsQuery) WithExpandEffectiveVolumes() GetAccountsQuery {
	a.Options.ExpandEffectiveVolumes = true

	return a
}
