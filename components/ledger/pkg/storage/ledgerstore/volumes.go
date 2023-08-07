package ledgerstore

import (
	"context"
	"fmt"

	"github.com/formancehq/ledger/pkg/core"
	"github.com/formancehq/ledger/pkg/storage"
	"github.com/uptrace/bun"
)

func (s *Store) GetAssetsVolumes(ctx context.Context, accountAddress string) (core.VolumesByAssets, error) {
	type Temp struct {
		Aggregated core.VolumesByAssets `bun:"aggregated"`
	}

	return fetchAndMap[*Temp, core.VolumesByAssets](s, ctx,
		func(temp *Temp) core.VolumesByAssets {
			return temp.Aggregated
		},
		func(query *bun.SelectQuery) *bun.SelectQuery {
			return query.
				ColumnExpr("aggregate_objects(aggregated_volumes) as aggregated").
				// TODO: Check SQL injections
				TableExpr(fmt.Sprintf(`aggregate_ledger_volumes(_accounts := '{"%s"}') v`, accountAddress)).
				TableExpr("volumes_to_jsonb(v) as aggregated_volumes")
		})
}

func (s *Store) UpdateVolumes(ctx context.Context) error {
	_, err := s.db.ExecContext(ctx, fmt.Sprintf(`call "%s".update_moves_volumes()`, s.name))
	return storage.PostgresError(err)
}

func (s *Store) UpdateEffectiveVolumes(ctx context.Context) error {
	_, err := s.db.ExecContext(ctx, fmt.Sprintf(`call "%s".update_moves_effective_volumes()`, s.name))
	return storage.PostgresError(err)
}
