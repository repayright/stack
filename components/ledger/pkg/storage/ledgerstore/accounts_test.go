package ledgerstore_test

import (
	"context"
	"math/big"
	"testing"

	"github.com/formancehq/ledger/pkg/core"
	"github.com/formancehq/ledger/pkg/storage/ledgerstore"
	"github.com/formancehq/stack/libs/go-libs/metadata"
	"github.com/stretchr/testify/require"
)

func TestUpdateAccountsMetadata(t *testing.T) {
	t.Parallel()
	store := newLedgerStore(t)

	t.Run("update metadata", func(t *testing.T) {
		metadata := metadata.Metadata{
			"foo": "bar",
		}

		require.NoError(t, store.InsertLogs(context.Background(),
			core.NewSetMetadataOnAccountLog(core.Now(), "bank", metadata).ChainLog(nil),
		), "account insertion should not fail")

		account, err := store.GetAccount(context.Background(), "bank")
		require.NoError(t, err, "account retrieval should not fail")

		require.Equal(t, "bank", account.Address, "account address should match")
		require.Equal(t, metadata, account.Metadata, "account metadata should match")
	})

	t.Run("success updating multiple account metadata", func(t *testing.T) {
		accounts := []core.Account{
			{
				Address:  "test:account1",
				Metadata: metadata.Metadata{"foo1": "bar1"},
			},
			{
				Address:  "test:account2",
				Metadata: metadata.Metadata{"foo2": "bar2"},
			},
			{
				Address:  "test:account3",
				Metadata: metadata.Metadata{"foo3": "bar3"},
			},
		}

		err := store.InsertLogs(context.Background(),
			core.NewSetMetadataOnAccountLog(core.Now(), "test:account1", metadata.Metadata{"foo1": "bar1"}).ChainLog(nil),
			core.NewSetMetadataOnAccountLog(core.Now(), "test:account2", metadata.Metadata{"foo2": "bar2"}).ChainLog(nil),
			core.NewSetMetadataOnAccountLog(core.Now(), "test:account3", metadata.Metadata{"foo3": "bar3"}).ChainLog(nil),
		)
		require.NoError(t, err, "account insertion should not fail")

		for _, account := range accounts {
			acc, err := store.GetAccount(context.Background(), account.Address)
			require.NoError(t, err, "account retrieval should not fail")

			require.Equal(t, account.Address, acc.Address, "account address should match")
			require.Equal(t, account.Metadata, acc.Metadata, "account metadata should match")
		}
	})
}

func TestGetAccount(t *testing.T) {
	t.Parallel()
	store := newLedgerStore(t)

	require.NoError(t, insertTransactions(context.Background(), store,
		*core.NewTransaction().WithPostings(
			core.NewPosting("world", "multi", "USD/2", big.NewInt(100)),
		),
	))

	account, err := store.GetAccount(context.Background(), "multi")
	require.NoError(t, err)
	require.Equal(t, core.Account{
		Address:  "multi",
		Metadata: metadata.Metadata{},
	}, *account)
}

func TestGetAccounts(t *testing.T) {
	t.Parallel()
	store := newLedgerStore(t)

	require.NoError(t, store.InsertLogs(context.Background(),
		core.NewSetMetadataOnAccountLog(core.Now(), "account1", metadata.Metadata{"category": "1"}).ChainLog(nil),
		core.NewSetMetadataOnAccountLog(core.Now(), "account2", metadata.Metadata{"category": "2"}).ChainLog(nil),
		core.NewSetMetadataOnAccountLog(core.Now(), "account3", metadata.Metadata{"category": "3"}).ChainLog(nil),
	))

	accounts, err := store.GetAccounts(context.Background(), ledgerstore.NewAccountsQuery())
	require.NoError(t, err)
	require.Len(t, accounts.Data, 3)

	accounts, err = store.GetAccounts(context.Background(), ledgerstore.NewAccountsQuery().WithMetadataFilter(metadata.Metadata{
		"category": "1",
	}))
	require.NoError(t, err)
	require.Len(t, accounts.Data, 1)
}

func TestGetAccountWithVolumes(t *testing.T) {
	t.Parallel()
	store := newLedgerStore(t)

	require.NoError(t, insertTransactions(context.Background(), store,
		*core.NewTransaction().WithPostings(
			core.NewPosting("world", "multi", "USD/2", big.NewInt(100)),
		),
	))

	accountWithVolumes, err := store.GetAccountWithVolumes(context.Background(), "multi", true, false)
	require.NoError(t, err)
	require.Equal(t, &core.AccountWithVolumes{
		Account: core.Account{
			Address:  "multi",
			Metadata: metadata.Metadata{},
		},
		Volumes: map[string]*core.Volumes{
			"USD/2": core.NewEmptyVolumes().WithInputInt64(100),
		},
	}, accountWithVolumes)
}

func TestUpdateAccountMetadata(t *testing.T) {
	t.Parallel()
	store := newLedgerStore(t)

	require.NoError(t, store.InsertLogs(context.Background(),
		core.NewSetMetadataOnAccountLog(core.Now(), "central_bank", metadata.Metadata{
			"foo": "bar",
		}).ChainLog(nil),
	))

	account, err := store.GetAccount(context.Background(), "central_bank")
	require.NoError(t, err)
	require.EqualValues(t, "bar", account.Metadata["foo"])
}

func TestGetAccountWithAccountNotExisting(t *testing.T) {
	t.Parallel()
	store := newLedgerStore(t)

	account, err := store.GetAccount(context.Background(), "account_not_existing")
	require.NoError(t, err)
	require.NotNil(t, account)
}

func TestCountAccounts(t *testing.T) {
	t.Parallel()
	store := newLedgerStore(t)

	require.NoError(t, insertTransactions(context.Background(), store,
		*core.NewTransaction().WithPostings(
			core.NewPosting("world", "central_bank", "USD/2", big.NewInt(100)),
		),
	))

	countAccounts, err := store.CountAccounts(context.Background(), ledgerstore.AccountsQuery{})
	require.NoError(t, err)
	require.EqualValues(t, 2, countAccounts) // world + central_bank
}
