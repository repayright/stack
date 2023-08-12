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

func TestGetBalancesAggregated(t *testing.T) {
	t.Parallel()
	store := newLedgerStore(t)
	now := core.Now()

	tx1 := core.NewTransaction().WithPostings(
		core.NewPosting("world", "users:1", "USD", big.NewInt(1)),
		core.NewPosting("world", "users:2", "USD", big.NewInt(199)),
	).WithTimestamp(now)

	tx2 := core.NewTransaction().WithPostings(
		core.NewPosting("world", "users:1", "USD", big.NewInt(1)),
		core.NewPosting("world", "users:2", "USD", big.NewInt(199)),
	).WithTimestamp(now.Add(time.Minute)).WithID(1)

	require.NoError(t, store.InsertLogs(context.Background(),
		core.ChainLogs(
			core.NewTransactionLog(tx1, map[string]metadata.Metadata{}),
			core.NewTransactionLog(tx2, map[string]metadata.Metadata{}),
		)...))

	t.Run("aggregate on all", func(t *testing.T) {
		q := ledgerstore.NewBalancesQuery().WithPageSize(10)
		cursor, err := store.GetAggregatedBalances(context.Background(), q)
		require.NoError(t, err)
		RequireEqual(t, core.BalancesByAssets{
			"USD": big.NewInt(0),
		}, cursor)
	})
	t.Run("filter on address", func(t *testing.T) {
		ret, err := store.GetAggregatedBalances(context.Background(), ledgerstore.NewBalancesQuery().WithPageSize(10).WithAddressFilter("users:"))
		require.NoError(t, err)
		require.Equal(t, core.BalancesByAssets{
			"USD": big.NewInt(400),
		}, ret)
	})
	t.Run("using pit", func(t *testing.T) {
		ret, err := store.GetAggregatedBalances(context.Background(), ledgerstore.NewBalancesQuery().
			WithPageSize(10).WithAddressFilter("users:").WithPIT(now))
		require.NoError(t, err)
		require.Equal(t, core.BalancesByAssets{
			"USD": big.NewInt(200),
		}, ret)
	})
}
