package ledgerstore_test

import (
	"context"
	"math/big"
	"testing"
	"time"

	ledger "github.com/formancehq/ledger/internal"
	"github.com/formancehq/ledger/internal/storage/ledgerstore"
	internaltesting "github.com/formancehq/ledger/internal/testing"
	"github.com/formancehq/stack/libs/go-libs/api"
	"github.com/formancehq/stack/libs/go-libs/metadata"
	"github.com/stretchr/testify/require"
)

func ExpandTransactions(txs ...*ledger.Transaction) []ledger.ExpandedTransaction {
	ret := make([]ledger.ExpandedTransaction, len(txs))
	accumulatedVolumes := ledger.AccountsAssetsVolumes{}
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
	now := ledger.Now()

	tx1 := ledger.ExpandedTransaction{
		Transaction: ledger.Transaction{
			ID: big.NewInt(0),
			TransactionData: ledger.TransactionData{
				Postings: []ledger.Posting{
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
		PostCommitVolumes: ledger.AccountsAssetsVolumes{
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
		PreCommitVolumes: ledger.AccountsAssetsVolumes{
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
	tx2 := ledger.ExpandedTransaction{
		Transaction: ledger.Transaction{
			ID: big.NewInt(1),
			TransactionData: ledger.TransactionData{
				Postings: []ledger.Posting{
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
		PostCommitVolumes: ledger.AccountsAssetsVolumes{
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
		PreCommitVolumes: ledger.AccountsAssetsVolumes{
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

	tx, err := store.GetTransactionWithVolumes(context.Background(), ledgerstore.NewGetTransactionQuery(tx1.ID).
		WithExpandVolumes().
		WithExpandEffectiveVolumes())
	require.NoError(t, err)
	require.Equal(t, tx1.Postings, tx.Postings)
	require.Equal(t, tx1.Reference, tx.Reference)
	require.Equal(t, tx1.Date, tx.Date)
	internaltesting.RequireEqual(t, ledger.AccountsAssetsVolumes{
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
	internaltesting.RequireEqual(t, ledger.AccountsAssetsVolumes{
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

	tx, err = store.GetTransactionWithVolumes(context.Background(), ledgerstore.NewGetTransactionQuery(tx2.ID).
		WithExpandVolumes().
		WithExpandEffectiveVolumes())
	require.Equal(t, tx2.Postings, tx.Postings)
	require.Equal(t, tx2.Reference, tx.Reference)
	require.Equal(t, tx2.Date, tx.Date)
	internaltesting.RequireEqual(t, ledger.AccountsAssetsVolumes{
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
	internaltesting.RequireEqual(t, ledger.AccountsAssetsVolumes{
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
	now := ledger.Now()

	t.Run("success inserting transaction", func(t *testing.T) {
		tx1 := ledger.ExpandedTransaction{
			Transaction: ledger.Transaction{
				ID: big.NewInt(0),
				TransactionData: ledger.TransactionData{
					Postings: ledger.Postings{
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
			PreCommitVolumes: map[string]ledger.VolumesByAssets{
				"world": map[string]*ledger.Volumes{
					"USD": ledger.NewEmptyVolumes(),
				},
				"alice": map[string]*ledger.Volumes{
					"USD": ledger.NewEmptyVolumes(),
				},
			},
			PostCommitVolumes: map[string]ledger.VolumesByAssets{
				"world": map[string]*ledger.Volumes{
					"USD": ledger.NewEmptyVolumes().WithOutputInt64(100),
				},
				"alice": map[string]*ledger.Volumes{
					"USD": ledger.NewEmptyVolumes().WithInputInt64(100),
				},
			},
		}

		err := insertTransactions(context.Background(), store, tx1.Transaction)
		require.NoError(t, err, "inserting transaction should not fail")

		tx, err := store.GetTransactionWithVolumes(context.Background(), ledgerstore.NewGetTransactionQuery(big.NewInt(0)).
			WithExpandVolumes())
		internaltesting.RequireEqual(t, tx1, *tx)
	})

	t.Run("success inserting multiple transactions", func(t *testing.T) {
		tx2 := ledger.ExpandedTransaction{
			Transaction: ledger.Transaction{
				ID: big.NewInt(1),
				TransactionData: ledger.TransactionData{
					Postings: ledger.Postings{
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
			PreCommitVolumes: map[string]ledger.VolumesByAssets{
				"world": map[string]*ledger.Volumes{
					"USD": ledger.NewEmptyVolumes().WithOutputInt64(100),
				},
				"polo": map[string]*ledger.Volumes{
					"USD": ledger.NewEmptyVolumes(),
				},
			},
			PostCommitVolumes: map[string]ledger.VolumesByAssets{
				"world": map[string]*ledger.Volumes{
					"USD": ledger.NewEmptyVolumes().WithOutputInt64(300),
				},
				"polo": map[string]*ledger.Volumes{
					"USD": ledger.NewEmptyVolumes().WithInputInt64(200),
				},
			},
		}

		tx3 := ledger.ExpandedTransaction{
			Transaction: ledger.Transaction{
				ID: big.NewInt(2),
				TransactionData: ledger.TransactionData{
					Postings: ledger.Postings{
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
			PreCommitVolumes: map[string]ledger.VolumesByAssets{
				"world": map[string]*ledger.Volumes{
					"USD": ledger.NewEmptyVolumes().WithOutputInt64(300),
				},
				"gfyrag": map[string]*ledger.Volumes{
					"USD": ledger.NewEmptyVolumes(),
				},
			},
			PostCommitVolumes: map[string]ledger.VolumesByAssets{
				"world": map[string]*ledger.Volumes{
					"USD": ledger.NewEmptyVolumes().WithOutputInt64(450),
				},
				"gfyrag": map[string]*ledger.Volumes{
					"USD": ledger.NewEmptyVolumes().WithInputInt64(150),
				},
			},
		}

		err := insertTransactions(context.Background(), store, tx2.Transaction, tx3.Transaction)
		require.NoError(t, err, "inserting multiple transactions should not fail")

		tx, err := store.GetTransactionWithVolumes(context.Background(), ledgerstore.NewGetTransactionQuery(big.NewInt(1)).WithExpandVolumes())
		require.NoError(t, err, "getting transaction should not fail")
		internaltesting.RequireEqual(t, tx2, *tx)

		tx, err = store.GetTransactionWithVolumes(context.Background(), ledgerstore.NewGetTransactionQuery(big.NewInt(2)).WithExpandVolumes())
		require.NoError(t, err, "getting transaction should not fail")
		internaltesting.RequireEqual(t, tx3, *tx)
	})
}

func TestCountTransactions(t *testing.T) {
	t.Parallel()
	store := newLedgerStore(t)
	now := ledger.Now()

	tx1 := ledger.ExpandedTransaction{
		Transaction: ledger.Transaction{
			ID: big.NewInt(0),
			TransactionData: ledger.TransactionData{
				Postings: ledger.Postings{
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
		PreCommitVolumes: map[string]ledger.VolumesByAssets{
			"world": map[string]*ledger.Volumes{
				"USD": ledger.NewEmptyVolumes(),
			},
			"alice": map[string]*ledger.Volumes{
				"USD": ledger.NewEmptyVolumes(),
			},
		},
		PostCommitVolumes: map[string]ledger.VolumesByAssets{
			"world": map[string]*ledger.Volumes{
				"USD": ledger.NewEmptyVolumes().WithOutputInt64(100),
			},
			"alice": map[string]*ledger.Volumes{
				"USD": ledger.NewEmptyVolumes().WithInputInt64(100),
			},
		},
	}
	tx2 := ledger.ExpandedTransaction{
		Transaction: ledger.Transaction{
			ID: big.NewInt(1),
			TransactionData: ledger.TransactionData{
				Postings: ledger.Postings{
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
		PreCommitVolumes: map[string]ledger.VolumesByAssets{
			"world": map[string]*ledger.Volumes{
				"USD": ledger.NewEmptyVolumes().WithOutputInt64(100),
			},
			"polo": map[string]*ledger.Volumes{
				"USD": ledger.NewEmptyVolumes(),
			},
		},
		PostCommitVolumes: map[string]ledger.VolumesByAssets{
			"world": map[string]*ledger.Volumes{
				"USD": ledger.NewEmptyVolumes().WithOutputInt64(300),
			},
			"polo": map[string]*ledger.Volumes{
				"USD": ledger.NewEmptyVolumes().WithInputInt64(200),
			},
		},
	}

	tx3 := ledger.ExpandedTransaction{
		Transaction: ledger.Transaction{
			ID: big.NewInt(2),
			TransactionData: ledger.TransactionData{
				Postings: ledger.Postings{
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
		PreCommitVolumes: map[string]ledger.VolumesByAssets{
			"world": map[string]*ledger.Volumes{
				"USD": ledger.NewEmptyVolumes().WithOutputInt64(300),
			},
			"gfyrag": map[string]*ledger.Volumes{
				"USD": ledger.NewEmptyVolumes(),
			},
		},
		PostCommitVolumes: map[string]ledger.VolumesByAssets{
			"world": map[string]*ledger.Volumes{
				"USD": ledger.NewEmptyVolumes().WithOutputInt64(450),
			},
			"gfyrag": map[string]*ledger.Volumes{
				"USD": ledger.NewEmptyVolumes().WithInputInt64(150),
			},
		},
	}

	err := insertTransactions(context.Background(), store, tx1.Transaction, tx2.Transaction, tx3.Transaction)
	require.NoError(t, err, "inserting transaction should not fail")

	count, err := store.CountTransactions(context.Background(), ledgerstore.GetTransactionsQuery{})
	require.NoError(t, err, "counting transactions should not fail")
	require.Equal(t, uint64(3), count, "count should be equal")
}

func TestUpdateTransactionsMetadata(t *testing.T) {
	t.Parallel()
	store := newLedgerStore(t)
	now := ledger.Now()

	tx1 := ledger.ExpandedTransaction{
		Transaction: ledger.Transaction{
			ID: big.NewInt(0),
			TransactionData: ledger.TransactionData{
				Postings: ledger.Postings{
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
		PreCommitVolumes: map[string]ledger.VolumesByAssets{
			"world": map[string]*ledger.Volumes{
				"USD": ledger.NewEmptyVolumes(),
			},
			"alice": map[string]*ledger.Volumes{
				"USD": ledger.NewEmptyVolumes(),
			},
		},
		PostCommitVolumes: map[string]ledger.VolumesByAssets{
			"world": map[string]*ledger.Volumes{
				"USD": ledger.NewEmptyVolumes().WithOutputInt64(100),
			},
			"alice": map[string]*ledger.Volumes{
				"USD": ledger.NewEmptyVolumes().WithInputInt64(100),
			},
		},
	}
	tx2 := ledger.ExpandedTransaction{
		Transaction: ledger.Transaction{
			ID: big.NewInt(1),
			TransactionData: ledger.TransactionData{
				Postings: ledger.Postings{
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
		PreCommitVolumes: map[string]ledger.VolumesByAssets{
			"world": map[string]*ledger.Volumes{
				"USD": ledger.NewEmptyVolumes().WithOutputInt64(100),
			},
			"polo": map[string]*ledger.Volumes{
				"USD": ledger.NewEmptyVolumes(),
			},
		},
		PostCommitVolumes: map[string]ledger.VolumesByAssets{
			"world": map[string]*ledger.Volumes{
				"USD": ledger.NewEmptyVolumes().WithOutputInt64(300),
			},
			"polo": map[string]*ledger.Volumes{
				"USD": ledger.NewEmptyVolumes().WithInputInt64(200),
			},
		},
	}

	err := insertTransactions(context.Background(), store, tx1.Transaction, tx2.Transaction)
	require.NoError(t, err, "inserting transaction should not fail")

	err = store.InsertLogs(context.Background(),
		ledger.NewSetMetadataOnTransactionLog(ledger.Now(), 0, metadata.Metadata{"foo1": "bar2"}).ChainLog(nil),
		ledger.NewSetMetadataOnTransactionLog(ledger.Now(), 1, metadata.Metadata{"foo2": "bar2"}).ChainLog(nil),
	)
	require.NoError(t, err, "updating multiple transaction metadata should not fail")

	tx, err := store.GetTransactionWithVolumes(context.Background(), ledgerstore.NewGetTransactionQuery(big.NewInt(0)).WithExpandVolumes().WithExpandEffectiveVolumes())
	require.NoError(t, err, "getting transaction should not fail")
	require.Equal(t, tx.Metadata, metadata.Metadata{"foo1": "bar2"}, "metadata should be equal")

	tx, err = store.GetTransactionWithVolumes(context.Background(), ledgerstore.NewGetTransactionQuery(big.NewInt(1)).WithExpandVolumes().WithExpandEffectiveVolumes())
	require.NoError(t, err, "getting transaction should not fail")
	require.Equal(t, tx.Metadata, metadata.Metadata{"foo2": "bar2"}, "metadata should be equal")
}

func TestInsertTransactionInPast(t *testing.T) {
	t.Parallel()
	store := newLedgerStore(t)
	now := ledger.Now()

	tx1 := ledger.NewTransaction().WithPostings(
		ledger.NewPosting("world", "bank", "USD/2", big.NewInt(100)),
	).WithDate(now)

	tx2 := ledger.NewTransaction().WithPostings(
		ledger.NewPosting("bank", "user1", "USD/2", big.NewInt(50)),
	).WithDate(now.Add(time.Hour)).WithIDUint64(1)

	// Insert in past must modify pre/post commit volumes of tx2
	tx3 := ledger.NewTransaction().WithPostings(
		ledger.NewPosting("bank", "user2", "USD/2", big.NewInt(50)),
	).WithDate(now.Add(30 * time.Minute)).WithIDUint64(2)

	require.NoError(t, insertTransactions(context.Background(), store, *tx1, *tx2))
	require.NoError(t, insertTransactions(context.Background(), store, *tx3))

	tx2FromDatabase, err := store.GetTransactionWithVolumes(context.Background(), ledgerstore.NewGetTransactionQuery(tx2.ID).WithExpandVolumes().WithExpandEffectiveVolumes())
	require.NoError(t, err)

	internaltesting.RequireEqual(t, ledger.AccountsAssetsVolumes{
		"bank": {
			"USD/2": ledger.NewVolumesInt64(100, 50),
		},
		"user1": {
			"USD/2": ledger.NewVolumesInt64(0, 0),
		},
	}, tx2FromDatabase.PreCommitEffectiveVolumes)
	internaltesting.RequireEqual(t, ledger.AccountsAssetsVolumes{
		"bank": {
			"USD/2": ledger.NewVolumesInt64(100, 100),
		},
		"user1": {
			"USD/2": ledger.NewVolumesInt64(50, 0),
		},
	}, tx2FromDatabase.PostCommitEffectiveVolumes)
}

func TestInsertTransactionInPastInOneBatch(t *testing.T) {
	t.Parallel()
	store := newLedgerStore(t)
	now := ledger.Now()

	tx1 := ledger.NewTransaction().WithPostings(
		ledger.NewPosting("world", "bank", "USD/2", big.NewInt(100)),
	).WithDate(now)

	tx2 := ledger.NewTransaction().WithPostings(
		ledger.NewPosting("bank", "user1", "USD/2", big.NewInt(50)),
	).WithDate(now.Add(time.Hour)).WithIDUint64(1)

	// Insert in past must modify pre/post commit volumes of tx2
	tx3 := ledger.NewTransaction().WithPostings(
		ledger.NewPosting("bank", "user2", "USD/2", big.NewInt(50)),
	).WithDate(now.Add(30 * time.Minute)).WithIDUint64(2)

	require.NoError(t, insertTransactions(context.Background(), store, *tx1, *tx2, *tx3))

	tx2FromDatabase, err := store.GetTransactionWithVolumes(context.Background(), ledgerstore.NewGetTransactionQuery(tx2.ID).WithExpandVolumes().WithExpandEffectiveVolumes())
	require.NoError(t, err)

	internaltesting.RequireEqual(t, ledger.AccountsAssetsVolumes{
		"bank": {
			"USD/2": ledger.NewVolumesInt64(100, 50),
		},
		"user1": {
			"USD/2": ledger.NewVolumesInt64(0, 0),
		},
	}, tx2FromDatabase.PreCommitEffectiveVolumes)
	internaltesting.RequireEqual(t, ledger.AccountsAssetsVolumes{
		"bank": {
			"USD/2": ledger.NewVolumesInt64(100, 100),
		},
		"user1": {
			"USD/2": ledger.NewVolumesInt64(50, 0),
		},
	}, tx2FromDatabase.PostCommitEffectiveVolumes)
}

func TestInsertTwoTransactionAtSameDateInSameBatch(t *testing.T) {
	t.Parallel()
	store := newLedgerStore(t)
	now := ledger.Now()

	tx1 := ledger.NewTransaction().WithPostings(
		ledger.NewPosting("world", "bank", "USD/2", big.NewInt(100)),
	).WithDate(now.Add(-time.Hour))

	tx2 := ledger.NewTransaction().WithPostings(
		ledger.NewPosting("bank", "user1", "USD/2", big.NewInt(10)),
	).WithDate(now).WithIDUint64(1)

	tx3 := ledger.NewTransaction().WithPostings(
		ledger.NewPosting("bank", "user2", "USD/2", big.NewInt(10)),
	).WithDate(now).WithIDUint64(2)

	require.NoError(t, insertTransactions(context.Background(), store, *tx1, *tx2, *tx3))

	tx2FromDatabase, err := store.GetTransactionWithVolumes(context.Background(), ledgerstore.NewGetTransactionQuery(tx2.ID).WithExpandVolumes().WithExpandEffectiveVolumes())
	require.NoError(t, err)

	internaltesting.RequireEqual(t, ledger.AccountsAssetsVolumes{
		"bank": {
			"USD/2": ledger.NewVolumesInt64(100, 10),
		},
		"user1": {
			"USD/2": ledger.NewVolumesInt64(10, 0),
		},
	}, tx2FromDatabase.PostCommitVolumes)
	internaltesting.RequireEqual(t, ledger.AccountsAssetsVolumes{
		"bank": {
			"USD/2": ledger.NewVolumesInt64(100, 0),
		},
		"user1": {
			"USD/2": ledger.NewVolumesInt64(0, 0),
		},
	}, tx2FromDatabase.PreCommitVolumes)

	tx3FromDatabase, err := store.GetTransactionWithVolumes(context.Background(), ledgerstore.NewGetTransactionQuery(tx3.ID).WithExpandVolumes().WithExpandEffectiveVolumes())
	require.NoError(t, err)

	internaltesting.RequireEqual(t, ledger.AccountsAssetsVolumes{
		"bank": {
			"USD/2": ledger.NewVolumesInt64(100, 10),
		},
		"user2": {
			"USD/2": ledger.NewVolumesInt64(0, 0),
		},
	}, tx3FromDatabase.PreCommitVolumes)
	internaltesting.RequireEqual(t, ledger.AccountsAssetsVolumes{
		"bank": {
			"USD/2": ledger.NewVolumesInt64(100, 20),
		},
		"user2": {
			"USD/2": ledger.NewVolumesInt64(10, 0),
		},
	}, tx3FromDatabase.PostCommitVolumes)
}

func TestInsertTwoTransactionAtSameDateInTwoBatch(t *testing.T) {
	t.Parallel()
	store := newLedgerStore(t)
	now := ledger.Now()

	tx1 := ledger.NewTransaction().WithPostings(
		ledger.NewPosting("world", "bank", "USD/2", big.NewInt(100)),
	).WithDate(now.Add(-time.Hour))

	tx2 := ledger.NewTransaction().WithPostings(
		ledger.NewPosting("bank", "user1", "USD/2", big.NewInt(10)),
	).WithDate(now).WithIDUint64(1)

	require.NoError(t, insertTransactions(context.Background(), store, *tx1, *tx2))

	tx3 := ledger.NewTransaction().WithPostings(
		ledger.NewPosting("bank", "user2", "USD/2", big.NewInt(10)),
	).WithDate(now).WithIDUint64(2)

	require.NoError(t, insertTransactions(context.Background(), store, *tx3))

	tx3FromDatabase, err := store.GetTransactionWithVolumes(context.Background(), ledgerstore.NewGetTransactionQuery(tx3.ID).WithExpandVolumes().WithExpandEffectiveVolumes())
	require.NoError(t, err)

	internaltesting.RequireEqual(t, ledger.AccountsAssetsVolumes{
		"bank": {
			"USD/2": ledger.NewVolumesInt64(100, 10),
		},
		"user2": {
			"USD/2": ledger.NewVolumesInt64(0, 0),
		},
	}, tx3FromDatabase.PreCommitVolumes)
	internaltesting.RequireEqual(t, ledger.AccountsAssetsVolumes{
		"bank": {
			"USD/2": ledger.NewVolumesInt64(100, 20),
		},
		"user2": {
			"USD/2": ledger.NewVolumesInt64(10, 0),
		},
	}, tx3FromDatabase.PostCommitVolumes)
}

func TestListTransactions(t *testing.T) {
	t.Parallel()
	store := newLedgerStore(t)
	now := ledger.Now()

	tx1 := ledger.NewTransaction().
		WithIDUint64(0).
		WithPostings(
			ledger.NewPosting("world", "alice", "USD", big.NewInt(100)),
		).
		WithMetadata(metadata.Metadata{"category": "1"}).
		WithDate(now.Add(-3 * time.Hour))
	tx2 := ledger.NewTransaction().
		WithIDUint64(1).
		WithPostings(
			ledger.NewPosting("world", "bob", "USD", big.NewInt(100)),
		).
		WithMetadata(metadata.Metadata{"category": "2"}).
		WithDate(now.Add(-2 * time.Hour))
	tx3 := ledger.NewTransaction().
		WithIDUint64(2).
		WithPostings(
			ledger.NewPosting("world", "users:marley", "USD", big.NewInt(100)),
		).
		WithMetadata(metadata.Metadata{"category": "3"}).
		WithDate(now.Add(-time.Hour))

	require.NoError(t, insertTransactions(context.Background(), store, *tx1, *tx2, *tx3))

	type testCase struct {
		name     string
		query    ledgerstore.GetTransactionsQuery
		expected *api.Cursor[ledger.ExpandedTransaction]
	}
	testCases := []testCase{
		{
			name:  "nominal",
			query: ledgerstore.NewTransactionsQuery(),
			expected: &api.Cursor[ledger.ExpandedTransaction]{
				PageSize: 15,
				HasMore:  false,
				Data:     Reverse(ExpandTransactions(tx1, tx2, tx3)...),
			},
		},
		{
			name: "address filter",
			query: ledgerstore.NewTransactionsQuery().
				WithAccountFilter("bob"),
			expected: &api.Cursor[ledger.ExpandedTransaction]{
				PageSize: 15,
				HasMore:  false,
				Data:     ExpandTransactions(tx1, tx2)[1:],
			},
		},
		{
			name: "address filter using segment",
			query: ledgerstore.NewTransactionsQuery().
				WithAccountFilter("users:"),
			expected: &api.Cursor[ledger.ExpandedTransaction]{
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
			expected: &api.Cursor[ledger.ExpandedTransaction]{
				PageSize: 15,
				HasMore:  false,
				Data:     ExpandTransactions(tx1, tx2, tx3)[1:2],
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			tc.query.Options.ExpandVolumes = true
			tc.query.Options.ExpandEffectiveVolumes = false
			cursor, err := store.GetTransactions(context.Background(), tc.query)
			require.NoError(t, err)
			internaltesting.RequireEqual(t, *tc.expected, *cursor)

			count, err := store.CountTransactions(context.Background(), tc.query)
			require.NoError(t, err)
			require.EqualValues(t, len(tc.expected.Data), count)
		})
	}
}
