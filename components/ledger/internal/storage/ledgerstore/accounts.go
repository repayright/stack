package ledgerstore

import (
	"context"
	"fmt"
	"math/big"
	"strings"

	ledger "github.com/formancehq/ledger/internal"
	storageerrors "github.com/formancehq/ledger/internal/storage"
	"github.com/formancehq/ledger/internal/storage/paginate"
	"github.com/formancehq/stack/libs/go-libs/api"
	"github.com/formancehq/stack/libs/go-libs/metadata"
	"github.com/formancehq/stack/libs/go-libs/pointer"
	"github.com/uptrace/bun"
)

func (store *Store) accountQueryBuilder(q PITFilter) func(query *bun.SelectQuery) *bun.SelectQuery {
	return func(query *bun.SelectQuery) *bun.SelectQuery {
		query = query.
			DistinctOn("accounts.address").
			Column("accounts.address").
			ColumnExpr("coalesce(metadata, '{}'::jsonb) as metadata").
			Table("accounts").
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

func (store *Store) GetAccountsWithVolumes(ctx context.Context, q GetAccountsQuery) (*api.Cursor[ledger.ExpandedAccount], error) {
	return paginateWithOffset[GetAccountsOptions, ledger.ExpandedAccount](store, ctx,
		paginate.OffsetPaginatedQuery[GetAccountsOptions](q),
		func(query *bun.SelectQuery) *bun.SelectQuery {
			query = store.
				accountQueryBuilder(q.Options.PITFilter)(query).
				Apply(filterMetadata(q.Options.Metadata))

			if len(q.Options.Addresses) > 0 {
				parts := make([]string, 0)
				for _, address := range q.Options.Addresses {
					parts = append(parts, filterAccountAddress(address, "accounts.address"))
				}
				query.Where("(" + strings.Join(parts, ") or (") + ")")
			}

			if q.Options.Balances != nil {
				for asset, balances := range q.Options.Balances {
					for operator, value := range balances {
						query = query.Join(fmt.Sprintf(`join lateral (
							select balance_from_volumes(post_commit_volumes) as balance
							from moves
							where moves.account_address = accounts.address and moves.asset = ?
							order by seq desc
							limit 1
						) m on m.balance %s ?`, operator), asset, (*paginate.BigInt)(value))
					}
				}
			}

			return query
		},
	)
}

type GetAccountQuery struct {
	PITFilter
	Addr string
}

func (q GetAccountQuery) WithPIT(pit ledger.Time) GetAccountQuery {
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

func (store *Store) GetAccount(ctx context.Context, address string) (*ledger.Account, error) {
	account, err := fetch[*ledger.Account](store, ctx, func(query *bun.SelectQuery) *bun.SelectQuery {
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
			return pointer.For(ledger.NewAccount(address)), nil
		}
		return nil, err
	}
	return account, nil
}

func (store *Store) GetAccountWithVolumes(ctx context.Context, q GetAccountQuery) (*ledger.ExpandedAccount, error) {
	account, err := fetch[*ledger.ExpandedAccount](store, ctx, func(query *bun.SelectQuery) *bun.SelectQuery {
		query = store.accountQueryBuilder(q.PITFilter)(query).
			Where("accounts.address = ?", q.Addr).
			Limit(1)

		return query
	})
	if err != nil {
		if storageerrors.IsNotFoundError(err) {
			return pointer.For(ledger.NewExpandedAccount(q.Addr)), nil
		}
		return nil, err
	}
	return account, nil
}

func (store *Store) CountAccounts(ctx context.Context, q GetAccountsQuery) (uint64, error) {
	return count(store, ctx, func(query *bun.SelectQuery) *bun.SelectQuery {
		query = store.
			accountQueryBuilder(q.Options.PITFilter)(query).
			Apply(filterMetadata(q.Options.Metadata))

		for _, address := range q.Options.Addresses {
			query.WhereOr(filterAccountAddress(address, "accounts.address"))
		}

		return query
	})
}

type GetAccountsQuery paginate.OffsetPaginatedQuery[GetAccountsOptions]

type PITFilter struct {
	PIT                    *ledger.Time `json:"pit"`
	ExpandVolumes          bool         `json:"volumes"`
	ExpandEffectiveVolumes bool         `json:"effectiveVolumes"`
}

type GetAccountsOptions struct {
	PITFilter
	AfterAddress string   `json:"after"`
	Addresses    []string `json:"address"`
	Balances     map[string]map[string]*big.Int
	Metadata     metadata.Metadata
}

func NewGetAccountsQuery() GetAccountsQuery {
	return GetAccountsQuery{
		PageSize: paginate.QueryDefaultPageSize,
		Order:    paginate.OrderAsc,
		Options: GetAccountsOptions{
			Metadata: metadata.Metadata{},
			Balances: map[string]map[string]*big.Int{},
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

func (a GetAccountsQuery) WithAddress(addresses ...string) GetAccountsQuery {
	a.Options.Addresses = addresses

	return a
}

func (a GetAccountsQuery) WithMetadataFilter(metadata metadata.Metadata) GetAccountsQuery {
	a.Options.Metadata = metadata

	return a
}

func (a GetAccountsQuery) WithPIT(date ledger.Time) GetAccountsQuery {
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

func (a GetAccountsQuery) WithBalances(balances map[string]map[string]*big.Int) GetAccountsQuery {
	a.Options.Balances = balances

	return a
}
