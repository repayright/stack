package ledgerstore

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/formancehq/ledger/pkg/core"
	"github.com/formancehq/stack/libs/go-libs/api"
)

type BalancesQueryFilters struct {
	AfterAddress  string `json:"afterAddress"`
	AddressRegexp string `json:"addressRegexp"`
}

type BalancesQuery OffsetPaginatedQuery[BalancesQueryFilters]

func NewBalancesQuery() BalancesQuery {
	return BalancesQuery{
		PageSize: QueryDefaultPageSize,
		Order:    OrderAsc,
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

type balancesByAssets core.BalancesByAssets

func (b *balancesByAssets) Scan(value interface{}) error {
	var i sql.NullString
	if err := i.Scan(value); err != nil {
		return err
	}

	*b = balancesByAssets{}
	if err := json.Unmarshal([]byte(i.String), b); err != nil {
		return err
	}

	return nil
}

func (s *Store) GetBalancesAggregated(ctx context.Context, q BalancesQuery) (core.BalancesByAssets, error) {

	/**
	with
	    potentially_staled_moves as (
	        select distinct on (m.account_address, m.asset) m.*
	        from moves m
	        order by account_address, asset, m.seq desc
	    ),
	    moves as (
	        select move.*
	        from potentially_staled_moves, ensure_move_computed(potentially_staled_moves) move
	    )
	select v.asset, sum(v.post_commit_aggregated_input) as inputs, sum(v.post_commit_aggregated_output) as outputs
	from moves v
	group by v.asset
	*/

	potentiallyStaledMoves := s.schema.
		NewSelect(MovesTableName).
		ColumnExpr("distinct on (moves.account_address, moves.asset) moves.*").
		Order("account_address", "asset", "moves.seq desc")

	if q.Filters.AddressRegexp != "" {
		src := strings.Split(q.Filters.AddressRegexp, ":")
		potentiallyStaledMoves.Where(fmt.Sprintf("jsonb_array_length(account_address_array) = %d", len(src)))

		for i, segment := range src {
			if segment == "" {
				continue
			}
			potentiallyStaledMoves.Where(fmt.Sprintf("account_address_array @@ ('$[%d] == \"' || ?::text || '\"')::jsonpath", i), segment)
		}
	}

	moves := s.schema.IDB.
		NewSelect().
		ColumnExpr("move.*").
		TableExpr("potentially_staled_moves").
		TableExpr("ensure_move_computed(potentially_staled_moves) move")

	type Temp struct {
		Aggregated core.VolumesByAssets `bun:"aggregated,type:jsonb"`
	}
	temp := Temp{}

	if err := s.schema.IDB.
		NewSelect().
		With("potentially_staled_moves", potentiallyStaledMoves).
		With("moves", moves).
		TableExpr("moves").
		ColumnExpr("volumes_to_jsonb((moves.asset, sum(moves.post_commit_aggregated_input), sum(moves.post_commit_aggregated_output))) as aggregated").
		Group("moves.asset").
		Scan(ctx, &temp); err != nil {
		return nil, err
	}

	return temp.Aggregated.Balances(), nil
}

func (s *Store) GetBalances(ctx context.Context, q BalancesQuery) (*api.Cursor[core.BalancesByAssetsByAccounts], error) {

	type Temp struct {
		Aggregated core.AccountsAssetsVolumes `bun:"aggregated,type:jsonb"`
	}

	query := s.schema.NewSelect(MovesTableName).
		ColumnExpr("distinct on (moves.account_address) jsonb_build_object(moves.account_address, aggregate_objects(volumes_to_jsonb)) as aggregated").
		TableExpr(`get_account_volumes_for_asset(moves.account_address, moves.asset) volumes`).
		TableExpr("volumes_to_jsonb(volumes)").
		Group("moves.account_address", "moves.asset").
		Order("moves.account_address", "moves.asset")

	if q.Filters.AddressRegexp != "" {
		// todo(gfyrag): factorize segments handling
		src := strings.Split(q.Filters.AddressRegexp, ":")
		query.Where(fmt.Sprintf("jsonb_array_length(account_address_array) = %d", len(src)))

		for i, segment := range src {
			if len(segment) == 0 {
				continue
			}
			query.Where(fmt.Sprintf(`account_address_array @@ ('$[%d] == "' || ?::text || '"')::jsonpath`, i), segment)
		}
	}

	if q.Filters.AfterAddress != "" {
		query.Where("account_address > ?", q.Filters.AfterAddress)
	}

	cursor, err := UsingOffset[BalancesQueryFilters, Temp](ctx,
		query, OffsetPaginatedQuery[BalancesQueryFilters](q))
	if err != nil {
		return nil, err
	}

	return api.MapCursor(cursor, func(from Temp) core.BalancesByAssetsByAccounts {
		return from.Aggregated.Balances()
	}), nil
}
