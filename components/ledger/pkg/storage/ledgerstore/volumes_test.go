package ledgerstore_test

import (
	"context"
	"math/big"
	"testing"
	"time"

	"github.com/formancehq/ledger/pkg/core"
	"github.com/formancehq/ledger/pkg/storage/ledgerstore"
	"github.com/formancehq/stack/libs/go-libs/metadata"
	"github.com/stretchr/testify/require"
	"github.com/uptrace/bun"
)

func TestGetAssetsVolumes(t *testing.T) {
	t.Parallel()
	store := newLedgerStore(t)
	now := core.Now()

	tx1 := core.NewTransaction().
		WithID(0).
		WithPostings(
			core.NewPosting("world", "alice", "USD", big.NewInt(100)),
		).
		WithTimestamp(now.Add(-3 * time.Hour))
	tx2 := core.NewTransaction().
		WithID(1).
		WithPostings(
			core.NewPosting("world", "bob", "USD", big.NewInt(100)),
		).
		WithTimestamp(now.Add(-2 * time.Hour))
	tx3 := core.NewTransaction().
		WithID(2).
		WithPostings(
			core.NewPosting("world", "users:marley", "USD", big.NewInt(100)),
		).
		WithTimestamp(now.Add(-time.Hour))

	require.NoError(t, insertTransactions(context.Background(), store, *tx1, *tx2, *tx3))

	assetVolumesForWorld, err := store.GetAssetsVolumes(context.Background(), "world")
	require.NoError(t, err, "get asset volumes should not fail")
	require.Equal(t, core.VolumesByAssets{
		"USD": core.NewEmptyVolumes().WithOutputInt64(300),
	}, assetVolumesForWorld, "asset volumes should be equal")

	assetVolumesForBob, err := store.GetAssetsVolumes(context.Background(), "bob")
	require.NoError(t, err, "get asset volumes should not fail")
	require.Equal(t, core.VolumesByAssets{
		"USD": core.NewEmptyVolumes().WithInputInt64(100),
	}, assetVolumesForBob, "asset volumes should be equal")
}

func TestUpdateVolumes(t *testing.T) {
	t.Parallel()
	store := newLedgerStore(t)

	modified, err := store.Migrate(context.Background())
	require.NoError(t, err)
	require.False(t, modified)

	require.NoError(t, store.InsertLogs(context.Background(),
		core.ChainLogs(
			core.NewTransactionLog(
				core.NewTransaction().WithPostings(core.NewPosting("world", "bank", "USD", big.NewInt(100))),
				map[string]metadata.Metadata{},
			),
			core.NewTransactionLog(
				core.NewTransaction().
					WithID(1).
					WithPostings(core.NewPosting("bank", "user:1", "USD", big.NewInt(50))),
				map[string]metadata.Metadata{},
			),
		)...,
	))

	require.NoError(t, store.UpdateVolumes(context.Background()))

	type Moves struct {
		bun.BaseModel `bun:"moves"`

		Inputs  *ledgerstore.BigInt `bun:"inputs,type:numeric"`
		Outputs *ledgerstore.BigInt `bun:"outputs,type:numeric"`
	}

	moves := make([]Moves, 0)
	err = store.GetDatabase().
		NewSelect().
		Model(&moves).
		ColumnExpr("(moves.post_commit_volumes).outputs as outputs").
		ColumnExpr("(moves.post_commit_volumes).inputs as inputs").
		Scan(context.Background())
	require.NoError(t, err)

	require.Len(t, moves, 4)
	for _, move := range moves {
		require.NotNil(t, move.Inputs)
		require.NotNil(t, move.Outputs)
	}
}

func TestUpdateEffectiveVolumes(t *testing.T) {
	t.Parallel()
	store := newLedgerStore(t)

	modified, err := store.Migrate(context.Background())
	require.NoError(t, err)
	require.False(t, modified)

	require.NoError(t, store.InsertLogs(context.Background(),
		core.ChainLogs(
			core.NewTransactionLog(
				core.NewTransaction().WithPostings(core.NewPosting("world", "bank", "USD", big.NewInt(100))),
				map[string]metadata.Metadata{},
			),
			core.NewTransactionLog(
				core.NewTransaction().
					WithID(1).
					WithPostings(core.NewPosting("bank", "user:1", "USD", big.NewInt(50))),
				map[string]metadata.Metadata{},
			),
		)...,
	))

	require.NoError(t, store.UpdateEffectiveVolumes(context.Background()))

	type Moves struct {
		bun.BaseModel `bun:"moves"`

		Inputs  *ledgerstore.BigInt `bun:"inputs,type:numeric"`
		Outputs *ledgerstore.BigInt `bun:"outputs,type:numeric"`
		Stale   bool                `bun:"stale"`
	}

	moves := make([]Moves, 0)
	err = store.GetDatabase().
		NewSelect().
		Model(&moves).
		ColumnExpr("(moves.post_commit_effective_volumes).outputs as outputs").
		ColumnExpr("(moves.post_commit_effective_volumes).inputs as inputs").
		Column("stale").
		Scan(context.Background())
	require.NoError(t, err)

	require.Len(t, moves, 4)
	for _, move := range moves {
		require.NotNil(t, move.Inputs)
		require.NotNil(t, move.Outputs)
		require.False(t, move.Stale)
	}
}
