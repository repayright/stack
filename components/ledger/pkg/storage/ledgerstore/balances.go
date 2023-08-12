package ledgerstore

import (
	"context"
	"math/big"

	"github.com/formancehq/ledger/pkg/core"
	"github.com/formancehq/ledger/pkg/storage/paginate"
	"github.com/uptrace/bun"
)

func (s *Store) GetAggregatedBalances(ctx context.Context, q BalancesQuery) (core.BalancesByAssets, error) {

	type Temp struct {
		Aggregated core.VolumesByAssets `bun:"aggregated,type:jsonb"`
	}
	return fetchAndMap[*Temp, core.BalancesByAssets](s, ctx,
		func(temp *Temp) core.BalancesByAssets {
			return temp.Aggregated.Balances()
		},
		func(query *bun.SelectQuery) *bun.SelectQuery {
			moves := s.db.
				NewSelect().
				Table(MovesTableName).
				ColumnExpr("distinct on (moves.account_address, moves.asset) moves.*").
				Order("account_address", "asset", "moves.seq desc").
				Apply(filterAccountAddress(q.Options.AddressRegexp, "account_address")).
				Apply(filterPIT(q.Options.PIT, "insertion_date")) // todo(gfyrag): expose capability to use effective_date

			return query.
				With("moves", moves).
				TableExpr("moves").
				ColumnExpr("volumes_to_jsonb((moves.asset, (sum((moves.post_commit_volumes).inputs), sum((moves.post_commit_volumes).outputs))::volumes)) as aggregated").
				Group("moves.asset")
		})
}

func (s *Store) GetBalance(ctx context.Context, address, asset string) (*big.Int, error) {
	type Temp struct {
		Balance *big.Int `bun:"balance,type:numeric"`
	}
	return fetchAndMap[*Temp, *big.Int](s, ctx, func(temp *Temp) *big.Int {
		return temp.Balance
	}, func(query *bun.SelectQuery) *bun.SelectQuery {
		return query.TableExpr("get_account_balance(?, ?) as balance", address, asset)
	})
}

type BalancesQueryOptions struct {
	AfterAddress  string `json:"afterAddress"`
	AddressRegexp string `json:"addressRegexp"`

	PIT core.Time `json:"pit"`
}

type BalancesQuery paginate.OffsetPaginatedQuery[BalancesQueryOptions]

func NewBalancesQuery() BalancesQuery {
	return BalancesQuery{
		PageSize: paginate.QueryDefaultPageSize,
		Order:    paginate.OrderAsc,
		Options:  BalancesQueryOptions{},
	}
}

func (q BalancesQuery) GetPageSize() uint64 {
	return q.PageSize
}

func (q BalancesQuery) WithAfterAddress(after string) BalancesQuery {
	q.Options.AfterAddress = after

	return q
}

func (q BalancesQuery) WithAddressFilter(address string) BalancesQuery {
	q.Options.AddressRegexp = address

	return q
}

func (q BalancesQuery) WithPageSize(pageSize uint64) BalancesQuery {
	q.PageSize = pageSize
	return q
}

func (q BalancesQuery) WithPIT(pit core.Time) BalancesQuery {
	q.Options.PIT = pit
	return q
}
