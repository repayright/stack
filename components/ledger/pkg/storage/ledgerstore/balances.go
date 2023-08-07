package ledgerstore

import (
	"context"
	"math/big"

	"github.com/formancehq/ledger/pkg/core"
	"github.com/formancehq/ledger/pkg/storage/paginate"
	"github.com/formancehq/stack/libs/go-libs/api"
	"github.com/uptrace/bun"
)

type BalancesQueryFilters struct {
	AfterAddress  string `json:"afterAddress"`
	AddressRegexp string `json:"addressRegexp"`
}

type BalancesQuery paginate.OffsetPaginatedQuery[BalancesQueryFilters]

func NewBalancesQuery() BalancesQuery {
	return BalancesQuery{
		PageSize: paginate.QueryDefaultPageSize,
		Order:    paginate.OrderAsc,
		Filters:  BalancesQueryFilters{},
	}
}

func (q BalancesQuery) GetPageSize() uint64 {
	return q.PageSize
}

func (b BalancesQuery) WithAfterAddress(after string) BalancesQuery {
	b.Filters.AfterAddress = after

	return b
}

func (b BalancesQuery) WithAddressFilter(address string) BalancesQuery {
	b.Filters.AddressRegexp = address

	return b
}

func (b BalancesQuery) WithPageSize(pageSize uint64) BalancesQuery {
	b.PageSize = pageSize
	return b
}

func (s *Store) GetAggregatedBalances(ctx context.Context, q BalancesQuery) (core.BalancesByAssets, error) {

	type Temp struct {
		Aggregated core.VolumesByAssets `bun:"aggregated,type:jsonb"`
	}
	return fetchAndMap[*Temp, core.BalancesByAssets](s, ctx,
		func(temp *Temp) core.BalancesByAssets {
			return temp.Aggregated.Balances()
		},
		func(query *bun.SelectQuery) *bun.SelectQuery {
			potentiallyStaledMoves := s.schema.
				NewSelect(MovesTableName).
				ColumnExpr("distinct on (moves.account_address, moves.asset) moves.*").
				Order("account_address", "asset", "moves.seq desc").
				Apply(filterAccountAddress(q.Filters.AddressRegexp, "account_address"))

			moves := s.schema.IDB.
				NewSelect().
				ColumnExpr("move.*").
				TableExpr("potentially_staled_moves").
				TableExpr("ensure_move_computed(potentially_staled_moves) move")

			return query.
				With("potentially_staled_moves", potentiallyStaledMoves).
				With("moves", moves).
				TableExpr("moves").
				ColumnExpr("volumes_to_jsonb((moves.asset, (sum((moves.post_commit_effective_volumes).inputs), sum((moves.post_commit_effective_volumes).outputs))::volumes)) as aggregated").
				Group("moves.asset")
		})
}

func (s *Store) GetBalances(ctx context.Context, q BalancesQuery) (*api.Cursor[core.BalancesByAssetsByAccounts], error) {

	type Temp struct {
		Aggregated core.AccountsAssetsVolumes `bun:"aggregated,type:jsonb"`
	}

	ret, err := paginateWithOffset[BalancesQueryFilters, *Temp](s, ctx,
		paginate.OffsetPaginatedQuery[BalancesQueryFilters](q),
		func(query *bun.SelectQuery) *bun.SelectQuery {
			query = query.
				ColumnExpr("distinct on (moves.account_address) jsonb_build_object(moves.account_address, aggregate_objects(volumes_to_jsonb)) as aggregated").
				Table("moves").
				TableExpr(`get_account_volumes_for_asset(moves.account_address, moves.asset) v`).
				TableExpr("volumes_to_jsonb(v)").
				Group("moves.account_address", "moves.asset").
				Order("moves.account_address", "moves.asset").
				Apply(filterAccountAddress(q.Filters.AddressRegexp, "account_address"))

			if q.Filters.AfterAddress != "" {
				query.Where("account_address > ?", q.Filters.AfterAddress)
			}

			return query
		})
	if err != nil {
		return nil, err
	}
	return api.MapCursor(ret, func(from *Temp) core.BalancesByAssetsByAccounts {
		return from.Aggregated.Balances()
	}), nil
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
