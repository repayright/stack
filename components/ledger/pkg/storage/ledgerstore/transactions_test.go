package ledgerstore_test

import (
	"context"
	"math/big"
	"testing"
	"time"

	"github.com/formancehq/ledger/pkg/core"
	"github.com/formancehq/ledger/pkg/storage/ledgerstore"
	"github.com/formancehq/stack/libs/go-libs/api"
	"github.com/formancehq/stack/libs/go-libs/metadata"
	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/require"
)

func bigIntComparer(v1 *big.Int, v2 *big.Int) bool {
	return v1.String() == v2.String()
}

func RequireEqual(t *testing.T, expected, actual any) {
	t.Helper()
	if diff := cmp.Diff(expected, actual, cmp.Comparer(bigIntComparer)); diff != "" {
		require.Failf(t, "Content not matching", diff)
	}
}

func ExpandTransactions(txs ...*core.Transaction) []core.ExpandedTransaction {
	ret := make([]core.ExpandedTransaction, len(txs))
	accumulatedVolumes := core.AccountsAssetsVolumes{}
	for ind, tx := range txs {
		ret[ind].Transaction = *tx
		for _, posting := range tx.Postings {
			ret[ind].PreCommitVolumes.AddInput(posting.Destination, posting.Asset, accumulatedVolumes.GetVolumes(posting.Destination, posting.Asset).Input)
			ret[ind].PreCommitVolumes.AddOutput(posting.Source, posting.Asset, accumulatedVolumes.GetVolumes(posting.Source, posting.Asset).Output)
		}
		for _, posting := range tx.Postings {
			accumulatedVolumes.AddOutput(posting.Source, posting.Asset, posting.Amount)
			accumulatedVolumes.AddInput(posting.Destination, posting.Asset, posting.Amount)
		}
		for _, posting := range tx.Postings {
			ret[ind].PostCommitVolumes.AddInput(posting.Destination, posting.Asset, accumulatedVolumes.GetVolumes(posting.Destination, posting.Asset).Input)
			ret[ind].PostCommitVolumes.AddOutput(posting.Source, posting.Asset, accumulatedVolumes.GetVolumes(posting.Source, posting.Asset).Output)
		}
	}
	return ret
}

func Reverse[T any](values ...T) []T {
	for i := 0; i < len(values)/2; i++ {
		values[i], values[len(values)-i-1] = values[len(values)-i-1], values[i]
	}
	return values
}

func TestGetTransaction(t *testing.T) {
	t.Parallel()
	store := newLedgerStore(t)
	now := core.Now()

	tx1 := core.ExpandedTransaction{
		Transaction: core.Transaction{
			TransactionData: core.TransactionData{
				Postings: []core.Posting{
					{
						Source:      "world",
						Destination: "central_bank",
						Amount:      big.NewInt(100),
						Asset:       "USD",
					},
				},
				Reference: "tx1",
				Date:      now.Add(-3 * time.Hour),
			},
		},
		PostCommitVolumes: core.AccountsAssetsVolumes{
			"world": {
				"USD": {
					Input:  big.NewInt(0),
					Output: big.NewInt(100),
				},
			},
			"central_bank": {
				"USD": {
					Input:  big.NewInt(100),
					Output: big.NewInt(0),
				},
			},
		},
		PreCommitVolumes: core.AccountsAssetsVolumes{
			"world": {
				"USD": {
					Input:  big.NewInt(0),
					Output: big.NewInt(0),
				},
			},
			"central_bank": {
				"USD": {
					Input:  big.NewInt(0),
					Output: big.NewInt(0),
				},
			},
		},
	}
	tx2 := core.ExpandedTransaction{
		Transaction: core.Transaction{
			ID: 1,
			TransactionData: core.TransactionData{
				Postings: []core.Posting{
					{
						Source:      "world",
						Destination: "central_bank",
						Amount:      big.NewInt(100),
						Asset:       "USD",
					},
				},
				Reference: "tx2",
				Date:      now.Add(-2 * time.Hour),
			},
		},
		PostCommitVolumes: core.AccountsAssetsVolumes{
			"world": {
				"USD": {
					Input:  big.NewInt(0),
					Output: big.NewInt(200),
				},
			},
			"central_bank": {
				"USD": {
					Input:  big.NewInt(200),
					Output: big.NewInt(0),
				},
			},
		},
		PreCommitVolumes: core.AccountsAssetsVolumes{
			"world": {
				"USD": {
					Input:  big.NewInt(0),
					Output: big.NewInt(100),
				},
			},
			"central_bank": {
				"USD": {
					Input:  big.NewInt(100),
					Output: big.NewInt(0),
				},
			},
		},
	}

	require.NoError(t, insertTransactions(context.Background(), store, tx1.Transaction, tx2.Transaction))

	tx, err := store.GetTransaction(context.Background(), tx1.ID)
	require.NoError(t, err)
	require.Equal(t, tx1.Postings, tx.Postings)
	require.Equal(t, tx1.Reference, tx.Reference)
	require.Equal(t, tx1.Date, tx.Date)
	RequireEqual(t, core.AccountsAssetsVolumes{
		"world": {
			"USD": {
				Input:  big.NewInt(0),
				Output: big.NewInt(100),
			},
		},
		"central_bank": {
			"USD": {
				Input:  big.NewInt(100),
				Output: big.NewInt(0),
			},
		},
	}, tx.PostCommitVolumes)
	RequireEqual(t, core.AccountsAssetsVolumes{
		"world": {
			"USD": {
				Input:  big.NewInt(0),
				Output: big.NewInt(0),
			},
		},
		"central_bank": {
			"USD": {
				Input:  big.NewInt(0),
				Output: big.NewInt(0),
			},
		},
	}, tx.PreCommitVolumes)

	tx, err = store.GetTransaction(context.Background(), tx2.ID)
	require.Equal(t, tx2.Postings, tx.Postings)
	require.Equal(t, tx2.Reference, tx.Reference)
	require.Equal(t, tx2.Date, tx.Date)
	RequireEqual(t, core.AccountsAssetsVolumes{
		"world": {
			"USD": {
				Input:  big.NewInt(0),
				Output: big.NewInt(200),
			},
		},
		"central_bank": {
			"USD": {
				Input:  big.NewInt(200),
				Output: big.NewInt(0),
			},
		},
	}, tx.PostCommitVolumes)
	RequireEqual(t, core.AccountsAssetsVolumes{
		"world": {
			"USD": {
				Input:  big.NewInt(0),
				Output: big.NewInt(100),
			},
		},
		"central_bank": {
			"USD": {
				Input:  big.NewInt(100),
				Output: big.NewInt(0),
			},
		},
	}, tx.PreCommitVolumes)
}

func TestInsertTransactions(t *testing.T) {
	t.Parallel()
	store := newLedgerStore(t)
	now := core.Now()

	t.Run("success inserting transaction", func(t *testing.T) {
		tx1 := core.ExpandedTransaction{
			Transaction: core.Transaction{
				ID: 0,
				TransactionData: core.TransactionData{
					Postings: core.Postings{
						{
							Source:      "world",
							Destination: "alice",
							Amount:      big.NewInt(100),
							Asset:       "USD",
						},
					},
					Date:     now.Add(-3 * time.Hour),
					Metadata: metadata.Metadata{},
				},
			},
			PreCommitVolumes: map[string]core.VolumesByAssets{
				"world": map[string]*core.Volumes{
					"USD": core.NewEmptyVolumes(),
				},
				"alice": map[string]*core.Volumes{
					"USD": core.NewEmptyVolumes(),
				},
			},
			PostCommitVolumes: map[string]core.VolumesByAssets{
				"world": map[string]*core.Volumes{
					"USD": core.NewEmptyVolumes().WithOutputInt64(100),
				},
				"alice": map[string]*core.Volumes{
					"USD": core.NewEmptyVolumes().WithInputInt64(100),
				},
			},
		}

		err := insertTransactions(context.Background(), store, tx1.Transaction)
		require.NoError(t, err, "inserting transaction should not fail")

		tx, err := store.GetTransaction(context.Background(), 0)
		RequireEqual(t, tx1, *tx)
	})

	t.Run("success inserting multiple transactions", func(t *testing.T) {
		tx2 := core.ExpandedTransaction{
			Transaction: core.Transaction{
				ID: 1,
				TransactionData: core.TransactionData{
					Postings: core.Postings{
						{
							Source:      "world",
							Destination: "polo",
							Amount:      big.NewInt(200),
							Asset:       "USD",
						},
					},
					Date:     now.Add(-2 * time.Hour),
					Metadata: metadata.Metadata{},
				},
			},
			PreCommitVolumes: map[string]core.VolumesByAssets{
				"world": map[string]*core.Volumes{
					"USD": core.NewEmptyVolumes().WithOutputInt64(100),
				},
				"polo": map[string]*core.Volumes{
					"USD": core.NewEmptyVolumes(),
				},
			},
			PostCommitVolumes: map[string]core.VolumesByAssets{
				"world": map[string]*core.Volumes{
					"USD": core.NewEmptyVolumes().WithOutputInt64(300),
				},
				"polo": map[string]*core.Volumes{
					"USD": core.NewEmptyVolumes().WithInputInt64(200),
				},
			},
		}

		tx3 := core.ExpandedTransaction{
			Transaction: core.Transaction{
				ID: 2,
				TransactionData: core.TransactionData{
					Postings: core.Postings{
						{
							Source:      "world",
							Destination: "gfyrag",
							Amount:      big.NewInt(150),
							Asset:       "USD",
						},
					},
					Date:     now.Add(-1 * time.Hour),
					Metadata: metadata.Metadata{},
				},
			},
			PreCommitVolumes: map[string]core.VolumesByAssets{
				"world": map[string]*core.Volumes{
					"USD": core.NewEmptyVolumes().WithOutputInt64(300),
				},
				"gfyrag": map[string]*core.Volumes{
					"USD": core.NewEmptyVolumes(),
				},
			},
			PostCommitVolumes: map[string]core.VolumesByAssets{
				"world": map[string]*core.Volumes{
					"USD": core.NewEmptyVolumes().WithOutputInt64(450),
				},
				"gfyrag": map[string]*core.Volumes{
					"USD": core.NewEmptyVolumes().WithInputInt64(150),
				},
			},
		}

		err := insertTransactions(context.Background(), store, tx2.Transaction, tx3.Transaction)
		require.NoError(t, err, "inserting multiple transactions should not fail")

		tx, err := store.GetTransaction(context.Background(), 1)
		require.NoError(t, err, "getting transaction should not fail")
		RequireEqual(t, tx2, *tx)

		tx, err = store.GetTransaction(context.Background(), 2)
		require.NoError(t, err, "getting transaction should not fail")
		RequireEqual(t, tx3, *tx)
	})
}

func TestCountTransactions(t *testing.T) {
	t.Parallel()
	store := newLedgerStore(t)
	now := core.Now()

	tx1 := core.ExpandedTransaction{
		Transaction: core.Transaction{
			ID: 0,
			TransactionData: core.TransactionData{
				Postings: core.Postings{
					{
						Source:      "world",
						Destination: "alice",
						Amount:      big.NewInt(100),
						Asset:       "USD",
					},
				},
				Date:     now.Add(-3 * time.Hour),
				Metadata: metadata.Metadata{},
			},
		},
		PreCommitVolumes: map[string]core.VolumesByAssets{
			"world": map[string]*core.Volumes{
				"USD": core.NewEmptyVolumes(),
			},
			"alice": map[string]*core.Volumes{
				"USD": core.NewEmptyVolumes(),
			},
		},
		PostCommitVolumes: map[string]core.VolumesByAssets{
			"world": map[string]*core.Volumes{
				"USD": core.NewEmptyVolumes().WithOutputInt64(100),
			},
			"alice": map[string]*core.Volumes{
				"USD": core.NewEmptyVolumes().WithInputInt64(100),
			},
		},
	}
	tx2 := core.ExpandedTransaction{
		Transaction: core.Transaction{
			ID: 1,
			TransactionData: core.TransactionData{
				Postings: core.Postings{
					{
						Source:      "world",
						Destination: "polo",
						Amount:      big.NewInt(200),
						Asset:       "USD",
					},
				},
				Date:     now.Add(-2 * time.Hour),
				Metadata: metadata.Metadata{},
			},
		},
		PreCommitVolumes: map[string]core.VolumesByAssets{
			"world": map[string]*core.Volumes{
				"USD": core.NewEmptyVolumes().WithOutputInt64(100),
			},
			"polo": map[string]*core.Volumes{
				"USD": core.NewEmptyVolumes(),
			},
		},
		PostCommitVolumes: map[string]core.VolumesByAssets{
			"world": map[string]*core.Volumes{
				"USD": core.NewEmptyVolumes().WithOutputInt64(300),
			},
			"polo": map[string]*core.Volumes{
				"USD": core.NewEmptyVolumes().WithInputInt64(200),
			},
		},
	}

	tx3 := core.ExpandedTransaction{
		Transaction: core.Transaction{
			ID: 2,
			TransactionData: core.TransactionData{
				Postings: core.Postings{
					{
						Source:      "world",
						Destination: "gfyrag",
						Amount:      big.NewInt(150),
						Asset:       "USD",
					},
				},
				Date:     now.Add(-1 * time.Hour),
				Metadata: metadata.Metadata{},
			},
		},
		PreCommitVolumes: map[string]core.VolumesByAssets{
			"world": map[string]*core.Volumes{
				"USD": core.NewEmptyVolumes().WithOutputInt64(300),
			},
			"gfyrag": map[string]*core.Volumes{
				"USD": core.NewEmptyVolumes(),
			},
		},
		PostCommitVolumes: map[string]core.VolumesByAssets{
			"world": map[string]*core.Volumes{
				"USD": core.NewEmptyVolumes().WithOutputInt64(450),
			},
			"gfyrag": map[string]*core.Volumes{
				"USD": core.NewEmptyVolumes().WithInputInt64(150),
			},
		},
	}

	err := insertTransactions(context.Background(), store, tx1.Transaction, tx2.Transaction, tx3.Transaction)
	require.NoError(t, err, "inserting transaction should not fail")

	count, err := store.CountTransactions(context.Background(), ledgerstore.TransactionsQuery{})
	require.NoError(t, err, "counting transactions should not fail")
	require.Equal(t, uint64(3), count, "count should be equal")
}

func TestUpdateTransactionsMetadata(t *testing.T) {
	t.Parallel()
	store := newLedgerStore(t)
	now := core.Now()

	tx1 := core.ExpandedTransaction{
		Transaction: core.Transaction{
			ID: 0,
			TransactionData: core.TransactionData{
				Postings: core.Postings{
					{
						Source:      "world",
						Destination: "alice",
						Amount:      big.NewInt(100),
						Asset:       "USD",
					},
				},
				Date:     now.Add(-3 * time.Hour),
				Metadata: metadata.Metadata{},
			},
		},
		PreCommitVolumes: map[string]core.VolumesByAssets{
			"world": map[string]*core.Volumes{
				"USD": core.NewEmptyVolumes(),
			},
			"alice": map[string]*core.Volumes{
				"USD": core.NewEmptyVolumes(),
			},
		},
		PostCommitVolumes: map[string]core.VolumesByAssets{
			"world": map[string]*core.Volumes{
				"USD": core.NewEmptyVolumes().WithOutputInt64(100),
			},
			"alice": map[string]*core.Volumes{
				"USD": core.NewEmptyVolumes().WithInputInt64(100),
			},
		},
	}
	tx2 := core.ExpandedTransaction{
		Transaction: core.Transaction{
			ID: 1,
			TransactionData: core.TransactionData{
				Postings: core.Postings{
					{
						Source:      "world",
						Destination: "polo",
						Amount:      big.NewInt(200),
						Asset:       "USD",
					},
				},
				Date:     now.Add(-2 * time.Hour),
				Metadata: metadata.Metadata{},
			},
		},
		PreCommitVolumes: map[string]core.VolumesByAssets{
			"world": map[string]*core.Volumes{
				"USD": core.NewEmptyVolumes().WithOutputInt64(100),
			},
			"polo": map[string]*core.Volumes{
				"USD": core.NewEmptyVolumes(),
			},
		},
		PostCommitVolumes: map[string]core.VolumesByAssets{
			"world": map[string]*core.Volumes{
				"USD": core.NewEmptyVolumes().WithOutputInt64(300),
			},
			"polo": map[string]*core.Volumes{
				"USD": core.NewEmptyVolumes().WithInputInt64(200),
			},
		},
	}

	err := insertTransactions(context.Background(), store, tx1.Transaction, tx2.Transaction)
	require.NoError(t, err, "inserting transaction should not fail")

	err = store.InsertLogs(context.Background(),
		core.NewSetMetadataOnTransactionLog(core.Now(), 0, metadata.Metadata{"foo1": "bar2"}).ChainLog(nil),
		core.NewSetMetadataOnTransactionLog(core.Now(), 1, metadata.Metadata{"foo2": "bar2"}).ChainLog(nil),
	)
	require.NoError(t, err, "updating multiple transaction metadata should not fail")

	tx, err := store.GetTransaction(context.Background(), 0)
	require.NoError(t, err, "getting transaction should not fail")
	require.Equal(t, tx.Metadata, metadata.Metadata{"foo1": "bar2"}, "metadata should be equal")

	tx, err = store.GetTransaction(context.Background(), 1)
	require.NoError(t, err, "getting transaction should not fail")
	require.Equal(t, tx.Metadata, metadata.Metadata{"foo2": "bar2"}, "metadata should be equal")
}

func TestInsertTransactionInPast(t *testing.T) {
	t.Parallel()
	store := newLedgerStore(t)
	now := core.Now()

	tx1 := core.NewTransaction().WithPostings(
		core.NewPosting("world", "bank", "USD/2", big.NewInt(100)),
	).WithTimestamp(now)

	tx2 := core.NewTransaction().WithPostings(
		core.NewPosting("bank", "user1", "USD/2", big.NewInt(50)),
	).WithTimestamp(now.Add(time.Hour)).WithID(1)

	// Insert in past must modify pre/post commit volumes of tx2
	tx3 := core.NewTransaction().WithPostings(
		core.NewPosting("bank", "user2", "USD/2", big.NewInt(50)),
	).WithTimestamp(now.Add(30 * time.Minute)).WithID(2)

	require.NoError(t, insertTransactions(context.Background(), store, *tx1, *tx2))
	require.NoError(t, insertTransactions(context.Background(), store, *tx3))

	tx2FromDatabase, err := store.GetTransaction(context.Background(), tx2.ID)
	require.NoError(t, err)

	RequireEqual(t, core.AccountsAssetsVolumes{
		"bank": {
			"USD/2": core.NewVolumesInt64(100, 50),
		},
		"user1": {
			"USD/2": core.NewVolumesInt64(0, 0),
		},
	}, tx2FromDatabase.PreCommitVolumes)
	RequireEqual(t, core.AccountsAssetsVolumes{
		"bank": {
			"USD/2": core.NewVolumesInt64(100, 100),
		},
		"user1": {
			"USD/2": core.NewVolumesInt64(50, 0),
		},
	}, tx2FromDatabase.PostCommitVolumes)
}

func TestInsertTransactionInPastInOneBatch(t *testing.T) {
	t.Parallel()
	store := newLedgerStore(t)
	now := core.Now()

	tx1 := core.NewTransaction().WithPostings(
		core.NewPosting("world", "bank", "USD/2", big.NewInt(100)),
	).WithTimestamp(now)

	tx2 := core.NewTransaction().WithPostings(
		core.NewPosting("bank", "user1", "USD/2", big.NewInt(50)),
	).WithTimestamp(now.Add(time.Hour)).WithID(1)

	// Insert in past must modify pre/post commit volumes of tx2
	tx3 := core.NewTransaction().WithPostings(
		core.NewPosting("bank", "user2", "USD/2", big.NewInt(50)),
	).WithTimestamp(now.Add(30 * time.Minute)).WithID(2)

	require.NoError(t, insertTransactions(context.Background(), store, *tx1, *tx2, *tx3))

	tx2FromDatabase, err := store.GetTransaction(context.Background(), tx2.ID)
	require.NoError(t, err)

	RequireEqual(t, core.AccountsAssetsVolumes{
		"bank": {
			"USD/2": core.NewVolumesInt64(100, 50),
		},
		"user1": {
			"USD/2": core.NewVolumesInt64(0, 0),
		},
	}, tx2FromDatabase.PreCommitVolumes)
	RequireEqual(t, core.AccountsAssetsVolumes{
		"bank": {
			"USD/2": core.NewVolumesInt64(100, 100),
		},
		"user1": {
			"USD/2": core.NewVolumesInt64(50, 0),
		},
	}, tx2FromDatabase.PostCommitVolumes)
}

func TestInsertTwoTransactionAtSameDateInSameBatch(t *testing.T) {
	t.Parallel()
	store := newLedgerStore(t)
	now := core.Now()

	tx1 := core.NewTransaction().WithPostings(
		core.NewPosting("world", "bank", "USD/2", big.NewInt(100)),
	).WithTimestamp(now.Add(-time.Hour))

	tx2 := core.NewTransaction().WithPostings(
		core.NewPosting("bank", "user1", "USD/2", big.NewInt(10)),
	).WithTimestamp(now).WithID(1)

	tx3 := core.NewTransaction().WithPostings(
		core.NewPosting("bank", "user2", "USD/2", big.NewInt(10)),
	).WithTimestamp(now).WithID(2)

	require.NoError(t, insertTransactions(context.Background(), store, *tx1, *tx2, *tx3))

	tx2FromDatabase, err := store.GetTransaction(context.Background(), tx2.ID)
	require.NoError(t, err)

	RequireEqual(t, core.AccountsAssetsVolumes{
		"bank": {
			"USD/2": core.NewVolumesInt64(100, 10),
		},
		"user1": {
			"USD/2": core.NewVolumesInt64(10, 0),
		},
	}, tx2FromDatabase.PostCommitVolumes)
	RequireEqual(t, core.AccountsAssetsVolumes{
		"bank": {
			"USD/2": core.NewVolumesInt64(100, 0),
		},
		"user1": {
			"USD/2": core.NewVolumesInt64(0, 0),
		},
	}, tx2FromDatabase.PreCommitVolumes)

	tx3FromDatabase, err := store.GetTransaction(context.Background(), tx3.ID)
	require.NoError(t, err)

	RequireEqual(t, core.AccountsAssetsVolumes{
		"bank": {
			"USD/2": core.NewVolumesInt64(100, 10),
		},
		"user2": {
			"USD/2": core.NewVolumesInt64(0, 0),
		},
	}, tx3FromDatabase.PreCommitVolumes)
	RequireEqual(t, core.AccountsAssetsVolumes{
		"bank": {
			"USD/2": core.NewVolumesInt64(100, 20),
		},
		"user2": {
			"USD/2": core.NewVolumesInt64(10, 0),
		},
	}, tx3FromDatabase.PostCommitVolumes)
}

func TestInsertTwoTransactionAtSameDateInTwoBatch(t *testing.T) {
	t.Parallel()
	store := newLedgerStore(t)
	now := core.Now()

	tx1 := core.NewTransaction().WithPostings(
		core.NewPosting("world", "bank", "USD/2", big.NewInt(100)),
	).WithTimestamp(now.Add(-time.Hour))

	tx2 := core.NewTransaction().WithPostings(
		core.NewPosting("bank", "user1", "USD/2", big.NewInt(10)),
	).WithTimestamp(now).WithID(1)

	require.NoError(t, insertTransactions(context.Background(), store, *tx1, *tx2))

	tx3 := core.NewTransaction().WithPostings(
		core.NewPosting("bank", "user2", "USD/2", big.NewInt(10)),
	).WithTimestamp(now).WithID(2)

	require.NoError(t, insertTransactions(context.Background(), store, *tx3))

	tx3FromDatabase, err := store.GetTransaction(context.Background(), tx3.ID)
	require.NoError(t, err)

	RequireEqual(t, core.AccountsAssetsVolumes{
		"bank": {
			"USD/2": core.NewVolumesInt64(100, 10),
		},
		"user2": {
			"USD/2": core.NewVolumesInt64(0, 0),
		},
	}, tx3FromDatabase.PreCommitVolumes)
	RequireEqual(t, core.AccountsAssetsVolumes{
		"bank": {
			"USD/2": core.NewVolumesInt64(100, 20),
		},
		"user2": {
			"USD/2": core.NewVolumesInt64(10, 0),
		},
	}, tx3FromDatabase.PostCommitVolumes)
}

func TestListTransactions(t *testing.T) {
	t.Parallel()
	store := newLedgerStore(t)
	now := core.Now()

	tx1 := core.NewTransaction().
		WithID(0).
		WithPostings(
			core.NewPosting("world", "alice", "USD", big.NewInt(100)),
		).
		WithMetadata(metadata.Metadata{"category": "1"}).
		WithTimestamp(now.Add(-3 * time.Hour))
	tx2 := core.NewTransaction().
		WithID(1).
		WithPostings(
			core.NewPosting("world", "bob", "USD", big.NewInt(100)),
		).
		WithMetadata(metadata.Metadata{"category": "2"}).
		WithTimestamp(now.Add(-2 * time.Hour))
	tx3 := core.NewTransaction().
		WithID(2).
		WithPostings(
			core.NewPosting("world", "users:marley", "USD", big.NewInt(100)),
		).
		WithMetadata(metadata.Metadata{"category": "3"}).
		WithTimestamp(now.Add(-time.Hour))

	require.NoError(t, insertTransactions(context.Background(), store, *tx1, *tx2, *tx3))

	type testCase struct {
		name     string
		query    ledgerstore.TransactionsQuery
		expected *api.Cursor[core.ExpandedTransaction]
	}
	testCases := []testCase{
		{
			name:  "nominal",
			query: ledgerstore.NewTransactionsQuery(),
			expected: &api.Cursor[core.ExpandedTransaction]{
				PageSize: 15,
				HasMore:  false,
				Data:     Reverse(ExpandTransactions(tx1, tx2, tx3)...),
			},
		},
		{
			name: "address filter",
			query: ledgerstore.NewTransactionsQuery().
				WithAccountFilter("bob"),
			expected: &api.Cursor[core.ExpandedTransaction]{
				PageSize: 15,
				HasMore:  false,
				Data:     ExpandTransactions(tx1, tx2)[1:],
			},
		},
		{
			name: "address filter using segment",
			query: ledgerstore.NewTransactionsQuery().
				WithAccountFilter("users:"),
			expected: &api.Cursor[core.ExpandedTransaction]{
				PageSize: 15,
				HasMore:  false,
				Data:     ExpandTransactions(tx1, tx2, tx3)[2:],
			},
		},
		{
			name: "filter using metadata",
			query: ledgerstore.NewTransactionsQuery().
				WithMetadataFilter(metadata.Metadata{
					"category": "2",
				}),
			expected: &api.Cursor[core.ExpandedTransaction]{
				PageSize: 15,
				HasMore:  false,
				Data:     ExpandTransactions(tx1, tx2, tx3)[1:2],
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			tc.query.Filters.ExpandVolumes = true
			cursor, err := store.GetTransactions(context.Background(), tc.query)
			require.NoError(t, err)
			RequireEqual(t, *tc.expected, *cursor)

			count, err := store.CountTransactions(context.Background(), tc.query)
			require.NoError(t, err)
			require.EqualValues(t, len(tc.expected.Data), count)
		})
	}
}
