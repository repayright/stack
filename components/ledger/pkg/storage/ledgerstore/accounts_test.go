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

		account, err := store.GetAccount(context.Background(), ledgerstore.NewGetAccountQuery("bank"))
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
			acc, err := store.GetAccount(context.Background(), ledgerstore.NewGetAccountQuery(account.Address))
			require.NoError(t, err, "account retrieval should not fail")

			require.Equal(t, account.Address, acc.Address, "account address should match")
			require.Equal(t, account.Metadata, acc.Metadata, "account metadata should match")
		}
	})
}

func TestGetAccount(t *testing.T) {
	t.Parallel()
	store := newLedgerStore(t)
	now := core.Now()

	require.NoError(t, store.InsertLogs(context.Background(),
		core.ChainLogs(
			core.NewTransactionLog(core.NewTransaction().WithPostings(
				core.NewPosting("world", "multi", "USD/2", big.NewInt(100)),
			).WithTimestamp(now), map[string]metadata.Metadata{}),
			core.NewSetMetadataLog(now.Add(time.Minute), core.SetMetadataLogPayload{
				TargetType: core.MetaTargetTypeAccount,
				TargetID:   "multi",
				Metadata: metadata.Metadata{
					"category": "gold",
				},
			}),
		)...,
	))

	t.Run("find account", func(t *testing.T) {
		account, err := store.GetAccount(context.Background(), ledgerstore.NewGetAccountQuery("multi"))
		require.NoError(t, err)
		require.Equal(t, core.Account{
			Address: "multi",
			Metadata: metadata.Metadata{
				"category": "gold",
			},
		}, *account)
	})

	t.Run("find account using pit", func(t *testing.T) {
		account, err := store.GetAccount(context.Background(), ledgerstore.NewGetAccountQuery("multi").WithPIT(now))
		require.NoError(t, err)
		require.Equal(t, core.Account{
			Address:  "multi",
			Metadata: metadata.Metadata{},
		}, *account)
	})

	t.Run("not existent account", func(t *testing.T) {
		account, err := store.GetAccount(context.Background(), ledgerstore.NewGetAccountQuery("account_not_existing"))
		require.NoError(t, err)
		require.NotNil(t, account)
	})
}

func TestGetAccounts(t *testing.T) {
	t.Parallel()
	store := newLedgerStore(t)
	now := core.Now()

	require.NoError(t, store.InsertLogs(context.Background(),
		core.ChainLogs(
			core.NewTransactionLog(
				core.NewTransaction().WithPostings(core.NewPosting("world", "account1", "USD", big.NewInt(100))),
				map[string]metadata.Metadata{
					"account1": {
						"category": "4",
					},
				},
			).WithDate(now),
			core.NewSetMetadataOnAccountLog(core.Now(), "account1", metadata.Metadata{"category": "1"}).WithDate(now.Add(time.Minute)),
			core.NewSetMetadataOnAccountLog(core.Now(), "account2", metadata.Metadata{"category": "2"}).WithDate(now.Add(2*time.Minute)),
			core.NewSetMetadataOnAccountLog(core.Now(), "account3", metadata.Metadata{"category": "3"}).WithDate(now.Add(3*time.Minute)),
		)...,
	))

	t.Run("list all", func(t *testing.T) {
		accounts, err := store.GetAccounts(context.Background(), ledgerstore.NewAccountsQuery())
		require.NoError(t, err)
		require.Len(t, accounts.Data, 4)
	})

	t.Run("list using metadata", func(t *testing.T) {
		accounts, err := store.GetAccounts(context.Background(), ledgerstore.NewAccountsQuery().WithMetadataFilter(metadata.Metadata{
			"category": "1",
		}))
		require.NoError(t, err)
		require.Len(t, accounts.Data, 1)
	})

	t.Run("list before date", func(t *testing.T) {
		accounts, err := store.GetAccounts(context.Background(), ledgerstore.NewAccountsQuery().WithPIT(now))
		require.NoError(t, err)
		require.Len(t, accounts.Data, 2)
	})
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

	account, err := store.GetAccount(context.Background(), ledgerstore.NewGetAccountQuery("central_bank"))
	require.NoError(t, err)
	require.EqualValues(t, "bar", account.Metadata["foo"])
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
