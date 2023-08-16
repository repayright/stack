package ledgerstore

import (
	"context"
	"math/big"

	ledger "github.com/formancehq/ledger/internal"
	"github.com/formancehq/ledger/internal/storage/paginate"
	"github.com/uptrace/bun"
)

func (store *Store) GetAggregatedBalances(ctx context.Context, q GetAggregatedBalancesQuery) (ledger.BalancesByAssets, error) {

	type Temp struct {
		Aggregated ledger.VolumesByAssets `bun:"aggregated,type:jsonb"`
	}
	return fetchAndMap[*Temp, ledger.BalancesByAssets](store, ctx,
		func(temp *Temp) ledger.BalancesByAssets {
			return temp.Aggregated.Balances()
		},
		func(query *bun.SelectQuery) *bun.SelectQuery {
			moves := store.db.
				NewSelect().
				Table(MovesTableName).
				ColumnExpr("distinct on (moves.account_address, moves.asset) moves.*").
				Order("account_address", "asset", "moves.seq desc").
				Apply(filterAccountAddressBuilder(q.Options.AddressRegexp, "account_address")).
				Apply(filterPIT(q.Options.PIT, "insertion_date")) // todo(gfyrag): expose capability to use effective_date

			return query.
				With("moves", moves).
				TableExpr("moves").
				ColumnExpr("volumes_to_jsonb((moves.asset, (sum((moves.post_commit_volumes).inputs), sum((moves.post_commit_volumes).outputs))::volumes)) as aggregated").
				Group("moves.asset")
		})
}

func (store *Store) GetBalance(ctx context.Context, address, asset string) (*big.Int, error) {
	type Temp struct {
		Balance *big.Int `bun:"balance,type:numeric"`
	}
	return fetchAndMap[*Temp, *big.Int](store, ctx, func(temp *Temp) *big.Int {
		return temp.Balance
	}, func(query *bun.SelectQuery) *bun.SelectQuery {
		return query.TableExpr("get_account_balance(?, ?) as balance", address, asset)
	})
}

type BalancesQueryOptions struct {
	AfterAddress  string `json:"afterAddress"`
	AddressRegexp string `json:"addressRegexp"`

	PIT *ledger.Time `json:"pit"`
}

type GetAggregatedBalancesQuery paginate.OffsetPaginatedQuery[BalancesQueryOptions]

func NewGetAggregatedBalancesQuery() GetAggregatedBalancesQuery {
	return GetAggregatedBalancesQuery{
		PageSize: paginate.QueryDefaultPageSize,
		Order:    paginate.OrderAsc,
		Options:  BalancesQueryOptions{},
	}
}

func (q GetAggregatedBalancesQuery) GetPageSize() uint64 {
	return q.PageSize
}

func (q GetAggregatedBalancesQuery) WithAfterAddress(after string) GetAggregatedBalancesQuery {
	q.Options.AfterAddress = after

	return q
}

func (q GetAggregatedBalancesQuery) WithAddressFilter(address string) GetAggregatedBalancesQuery {
	q.Options.AddressRegexp = address

	return q
}

func (q GetAggregatedBalancesQuery) WithPageSize(pageSize uint64) GetAggregatedBalancesQuery {
	q.PageSize = pageSize
	return q
}

func (q GetAggregatedBalancesQuery) WithPIT(pit ledger.Time) GetAggregatedBalancesQuery {
	q.Options.PIT = &pit
	return q
}
