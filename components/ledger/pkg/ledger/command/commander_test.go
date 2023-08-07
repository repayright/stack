package command

import (
	"context"
	"math/big"
	"testing"

	"github.com/formancehq/ledger/pkg/core"
	storageerrors "github.com/formancehq/ledger/pkg/storage"
	"github.com/formancehq/stack/libs/go-libs/logging"
	"github.com/formancehq/stack/libs/go-libs/metadata"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
)

var (
	now = core.Now()
)

type testCase struct {
	name          string
	setup         func(t *testing.T, r Store)
	script        string
	reference     string
	expectedError error
	expectedTx    *core.Transaction
	expectedLogs  []*core.Log
	parameters    Parameters
}

var testCases = []testCase{
	{
		name: "nominal",
		script: `
			send [GEM 100] (
				source = @world
				destination = @mint
			)`,
		expectedTx: core.NewTransaction().WithPostings(
			core.NewPosting("world", "mint", "GEM", big.NewInt(100)),
		),
		expectedLogs: []*core.Log{
			core.NewTransactionLog(
				core.NewTransaction().WithPostings(
					core.NewPosting("world", "mint", "GEM", big.NewInt(100))),
				map[string]metadata.Metadata{},
			),
		},
	},
	{
		name:          "no script",
		script:        ``,
		expectedError: ErrNoScript,
	},
	{
		name:          "invalid script",
		script:        `XXX`,
		expectedError: ErrCompilationFailed,
	},
	{
		name: "set reference conflict",
		setup: func(t *testing.T, store Store) {
			tx := core.NewTransaction().
				WithPostings(core.NewPosting("world", "mint", "GEM", big.NewInt(100))).
				WithReference("tx_ref")
			log := core.NewTransactionLog(tx, nil)
			err := store.InsertLogs(context.Background(), log.ChainLog(nil))
			require.NoError(t, err)
		},
		script: `
			send [GEM 100] (
				source = @world
				destination = @mint
			)`,
		reference:     "tx_ref",
		expectedError: ErrConflictError,
	},
	{
		name: "set reference",
		script: `
			send [GEM 100] (
				source = @world
				destination = @mint
			)`,
		reference: "tx_ref",
		expectedTx: core.NewTransaction().
			WithPostings(
				core.NewPosting("world", "mint", "GEM", big.NewInt(100)),
			).
			WithReference("tx_ref"),
		expectedLogs: []*core.Log{
			core.NewTransactionLog(
				core.NewTransaction().
					WithPostings(
						core.NewPosting("world", "mint", "GEM", big.NewInt(100)),
					).
					WithReference("tx_ref"),
				map[string]metadata.Metadata{},
			),
		},
	},
	{
		name: "using idempotency",
		script: `
			send [GEM 100] (
				source = @world
				destination = @mint
			)`,
		reference: "tx_ref",
		expectedTx: core.NewTransaction().
			WithPostings(
				core.NewPosting("world", "mint", "GEM", big.NewInt(100)),
			),
		expectedLogs: []*core.Log{
			core.NewTransactionLog(
				core.NewTransaction().
					WithPostings(
						core.NewPosting("world", "mint", "GEM", big.NewInt(100)),
					),
				map[string]metadata.Metadata{},
			).WithIdempotencyKey("testing"),
		},
		setup: func(t *testing.T, r Store) {
			log := core.NewTransactionLog(
				core.NewTransaction().
					WithPostings(
						core.NewPosting("world", "mint", "GEM", big.NewInt(100)),
					).
					WithTimestamp(now),
				map[string]metadata.Metadata{},
			).WithIdempotencyKey("testing")
			err := r.InsertLogs(context.Background(), log.ChainLog(nil))
			require.NoError(t, err)
		},
		parameters: Parameters{
			IdempotencyKey: "testing",
		},
	},
}

func TestCreateTransaction(t *testing.T) {
	t.Parallel()

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {

			store := storageerrors.NewInMemoryStore()
			ctx := logging.TestingContext()

			commander := New(store, NoOpLocker, NewCompiler(1024), NewReferencer())
			go commander.Run(ctx)
			defer commander.Close()

			if tc.setup != nil {
				tc.setup(t, store)
			}
			ret, err := commander.CreateTransaction(ctx, tc.parameters, core.RunScript{
				Script: core.Script{
					Plain: tc.script,
				},
				Timestamp: now,
				Reference: tc.reference,
			})

			if tc.expectedError != nil {
				require.True(t, errors.Is(err, tc.expectedError))
			} else {
				require.NoError(t, err)
				require.NotNil(t, ret)
				tc.expectedTx.Date = now
				require.Equal(t, tc.expectedTx, ret)

				for ind := range tc.expectedLogs {
					expectedLog := tc.expectedLogs[ind]
					switch v := expectedLog.Data.(type) {
					case core.NewTransactionLogPayload:
						v.Transaction.Date = now
						expectedLog.Data = v
					}
					expectedLog.Date = now
				}
			}
		})
	}
}

func TestRevert(t *testing.T) {
	txID := uint64(0)
	store := storageerrors.NewInMemoryStore()
	ctx := logging.TestingContext()

	log := core.NewTransactionLog(
		core.NewTransaction().WithPostings(
			core.NewPosting("world", "bank", "USD", big.NewInt(100)),
		),
		map[string]metadata.Metadata{},
	).ChainLog(nil)
	err := store.InsertLogs(context.Background(), log)
	require.NoError(t, err)

	commander := New(store, NoOpLocker, NewCompiler(1024), NewReferencer())
	go commander.Run(ctx)
	defer commander.Close()

	_, err = commander.RevertTransaction(ctx, Parameters{}, txID)
	require.NoError(t, err)
}

func TestRevertWithAlreadyReverted(t *testing.T) {

	store := storageerrors.NewInMemoryStore()
	ctx := logging.TestingContext()

	err := store.InsertLogs(context.Background(),
		core.NewTransactionLog(
			core.NewTransaction().WithPostings(core.NewPosting("world", "bank", "USD", big.NewInt(100))),
			map[string]metadata.Metadata{},
		).ChainLog(nil),
		core.NewRevertedTransactionLog(core.Now(), 0, core.NewTransaction()).ChainLog(nil),
	)
	require.NoError(t, err)

	commander := New(store, NoOpLocker, NewCompiler(1024), NewReferencer())
	go commander.Run(ctx)
	defer commander.Close()

	_, err = commander.RevertTransaction(context.Background(), Parameters{}, 0)
	require.True(t, errors.Is(err, ErrAlreadyReverted))
}

func TestRevertWithRevertOccurring(t *testing.T) {

	store := storageerrors.NewInMemoryStore()
	ctx := logging.TestingContext()

	log := core.NewTransactionLog(
		core.NewTransaction().WithPostings(
			core.NewPosting("world", "bank", "USD", big.NewInt(100)),
		),
		map[string]metadata.Metadata{},
	)
	err := store.InsertLogs(ctx, log.ChainLog(nil))
	require.NoError(t, err)

	referencer := NewReferencer()
	commander := New(store, NoOpLocker, NewCompiler(1024), referencer)
	go commander.Run(ctx)
	defer commander.Close()

	referencer.take(referenceReverts, uint64(0))

	_, err = commander.RevertTransaction(ctx, Parameters{}, 0)
	require.True(t, errors.Is(err, ErrRevertOccurring))
}
