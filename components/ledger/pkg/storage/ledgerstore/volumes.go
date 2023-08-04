package ledgerstore

import (
	"context"
	"fmt"

	"github.com/formancehq/ledger/pkg/core"
	storageerrors "github.com/formancehq/ledger/pkg/storage"
)

func (s *Store) GetAssetsVolumes(ctx context.Context, accountAddress string) (core.VolumesByAssets, error) {

	type Temp struct {
		Aggregated core.VolumesByAssets `bun:"aggregated"`
	}

	temp := Temp{}
	err := s.schema.IDB.NewSelect().
		ColumnExpr("aggregate_objects(aggregated_volumes) as aggregated").
		// TODO: Check SQL injections
		TableExpr(fmt.Sprintf(`aggregate_ledger_volumes(_accounts := '{"%s"}') volumes`, accountAddress)).
		TableExpr("volumes_to_jsonb(volumes) as aggregated_volumes").
		Scan(ctx, &temp)
	if err != nil {
		return nil, storageerrors.PostgresError(err)
	}
	return temp.Aggregated, err
}
