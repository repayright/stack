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

func TestGetAccounts(t *testing.T) {
	t.Parallel()
	store := newLedgerStore(t)
	now := core.Now()

	require.NoError(t, store.InsertLogs(context.Background(),
		core.ChainLogs(
			core.NewTransactionLog(
				core.NewTransaction().
					WithPostings(core.NewPosting("world", "account:1", "USD", big.NewInt(100))).
					WithDate(now),
				map[string]metadata.Metadata{
					"account:1": {
						"category": "4",
					},
				},
			).WithDate(now),
			core.NewSetMetadataOnAccountLog(core.Now(), "account:1", metadata.Metadata{"category": "1"}).WithDate(now.Add(time.Minute)),
			core.NewSetMetadataOnAccountLog(core.Now(), "account:2", metadata.Metadata{"category": "2"}).WithDate(now.Add(2*time.Minute)),
			core.NewSetMetadataOnAccountLog(core.Now(), "account:3", metadata.Metadata{"category": "3"}).WithDate(now.Add(3*time.Minute)),
			core.NewTransactionLog(
				core.NewTransaction().
					WithPostings(core.NewPosting("world", "account:1", "USD", big.NewInt(100))).
					WithIDUint64(1).
					WithDate(now.Add(4*time.Minute)),
				map[string]metadata.Metadata{},
			).WithDate(now.Add(100*time.Millisecond)),
			core.NewTransactionLog(
				core.NewTransaction().
					WithPostings(core.NewPosting("account:1", "bank", "USD", big.NewInt(50))).
					WithDate(now.Add(3*time.Minute)).
					WithIDUint64(2),
				map[string]metadata.Metadata{},
			).WithDate(now.Add(200*time.Millisecond)),
		)...,
	))

	t.Run("list all", func(t *testing.T) {
		accounts, err := store.GetAccountsWithVolumes(context.Background(), ledgerstore.NewGetAccountsQuery())
		require.NoError(t, err)
		require.Len(t, accounts.Data, 5)
	})

	t.Run("list using metadata", func(t *testing.T) {
		accounts, err := store.GetAccountsWithVolumes(context.Background(), ledgerstore.NewGetAccountsQuery().WithMetadataFilter(metadata.Metadata{
			"category": "1",
		}))
		require.NoError(t, err)
		require.Len(t, accounts.Data, 1)
	})

	t.Run("list before date", func(t *testing.T) {
		accounts, err := store.GetAccountsWithVolumes(context.Background(), ledgerstore.NewGetAccountsQuery().WithPIT(now))
		require.NoError(t, err)
		require.Len(t, accounts.Data, 2)
	})

	t.Run("list with volumes", func(t *testing.T) {
		accounts, err := store.GetAccountsWithVolumes(context.Background(), ledgerstore.NewGetAccountsQuery().
			WithAddress("account:1").
			WithExpandVolumes())
		require.NoError(t, err)
		require.Len(t, accounts.Data, 1)
		require.Equal(t, core.VolumesByAssets{
			"USD": core.NewVolumesInt64(200, 50),
		}, accounts.Data[0].Volumes)
	})

	t.Run("list with volumes using PIT", func(t *testing.T) {
		accounts, err := store.GetAccountsWithVolumes(context.Background(), ledgerstore.NewGetAccountsQuery().
			WithPIT(now).
			WithAddress("account:1").
			WithExpandVolumes())
		require.NoError(t, err)
		require.Len(t, accounts.Data, 1)
		require.Equal(t, core.VolumesByAssets{
			"USD": core.NewVolumesInt64(100, 0),
		}, accounts.Data[0].Volumes)
	})

	t.Run("list with effective volumes", func(t *testing.T) {
		accounts, err := store.GetAccountsWithVolumes(context.Background(), ledgerstore.NewGetAccountsQuery().
			WithAddress("account:1").
			WithExpandEffectiveVolumes())
		require.NoError(t, err)
		require.Len(t, accounts.Data, 1)
		require.Equal(t, core.VolumesByAssets{
			"USD": core.NewVolumesInt64(200, 50),
		}, accounts.Data[0].EffectiveVolumes)
	})

	t.Run("list with effective volumes using PIT", func(t *testing.T) {
		accounts, err := store.GetAccountsWithVolumes(context.Background(), ledgerstore.NewGetAccountsQuery().
			WithAddress("account:1").
			WithPIT(now).
			WithExpandEffectiveVolumes())
		require.NoError(t, err)
		require.Len(t, accounts.Data, 1)
		require.Equal(t, core.VolumesByAssets{
			"USD": core.NewVolumesInt64(100, 0),
		}, accounts.Data[0].EffectiveVolumes)
	})

	t.Run("list using filter on address", func(t *testing.T) {
		accounts, err := store.GetAccountsWithVolumes(context.Background(), ledgerstore.NewGetAccountsQuery().
			WithAddress("account:"))
		require.NoError(t, err)
		require.Len(t, accounts.Data, 3)
	})
}

func TestUpdateAccountsMetadata(t *testing.T) {
	t.Parallel()
	store := newLedgerStore(t)

	metadata := metadata.Metadata{
		"foo": "bar",
	}

	require.NoError(t, store.InsertLogs(context.Background(),
		core.NewSetMetadataOnAccountLog(core.Now(), "bank", metadata).ChainLog(nil),
	), "account insertion should not fail")

	account, err := store.GetAccountWithVolumes(context.Background(), ledgerstore.NewGetAccountQuery("bank"))
	require.NoError(t, err, "account retrieval should not fail")

	require.Equal(t, "bank", account.Address, "account address should match")
	require.Equal(t, metadata, account.Metadata, "account metadata should match")
}

func TestGetAccount(t *testing.T) {
	t.Parallel()
	store := newLedgerStore(t)
	now := core.Now()

	require.NoError(t, store.InsertLogs(context.Background(),
		core.ChainLogs(
			core.NewTransactionLog(core.NewTransaction().WithPostings(
				core.NewPosting("world", "multi", "USD/2", big.NewInt(100)),
			).WithDate(now), map[string]metadata.Metadata{}),
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
		account, err := store.GetAccountWithVolumes(context.Background(), ledgerstore.NewGetAccountQuery("multi"))
		require.NoError(t, err)
		require.Equal(t, core.ExpandedAccount{
			Account: core.Account{
				Address: "multi",
				Metadata: metadata.Metadata{
					"category": "gold",
				},
			},
		}, *account)
	})

	t.Run("find account with volumes", func(t *testing.T) {
		account, err := store.GetAccountWithVolumes(context.Background(), ledgerstore.
			NewGetAccountQuery("multi").
			WithExpandVolumes())
		require.NoError(t, err)
		require.Equal(t, core.ExpandedAccount{
			Account: core.Account{
				Address: "multi",
				Metadata: metadata.Metadata{
					"category": "gold",
				},
			},
			Volumes: core.VolumesByAssets{
				"USD/2": core.NewVolumesInt64(100, 0),
			},
		}, *account)
	})

	t.Run("find account with effective volumes", func(t *testing.T) {
		account, err := store.GetAccountWithVolumes(context.Background(), ledgerstore.
			NewGetAccountQuery("multi").
			WithExpandEffectiveVolumes())
		require.NoError(t, err)
		require.Equal(t, core.ExpandedAccount{
			Account: core.Account{
				Address: "multi",
				Metadata: metadata.Metadata{
					"category": "gold",
				},
			},
			EffectiveVolumes: core.VolumesByAssets{
				"USD/2": core.NewVolumesInt64(100, 0),
			},
		}, *account)
	})

	t.Run("find account using pit", func(t *testing.T) {
		account, err := store.GetAccountWithVolumes(context.Background(), ledgerstore.NewGetAccountQuery("multi").WithPIT(now))
		require.NoError(t, err)
		require.Equal(t, core.ExpandedAccount{
			Account: core.Account{
				Address:  "multi",
				Metadata: metadata.Metadata{},
			},
			Volumes: core.VolumesByAssets{},
		}, *account)
	})

	t.Run("not existent account", func(t *testing.T) {
		account, err := store.GetAccountWithVolumes(context.Background(), ledgerstore.NewGetAccountQuery("account_not_existing"))
		require.NoError(t, err)
		require.NotNil(t, account)
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

	accountWithVolumes, err := store.GetAccountWithVolumes(context.Background(),
		ledgerstore.NewGetAccountQuery("multi").WithExpandVolumes())
	require.NoError(t, err)
	require.Equal(t, &core.ExpandedAccount{
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

	account, err := store.GetAccountWithVolumes(context.Background(), ledgerstore.NewGetAccountQuery("central_bank"))
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

	countAccounts, err := store.CountAccounts(context.Background(), ledgerstore.GetAccountsQuery{})
	require.NoError(t, err)
	require.EqualValues(t, 2, countAccounts) // world + central_bank
}
