package ledgerstore_test

import (
	"context"
	"testing"

	"github.com/formancehq/ledger/pkg/core"
	"github.com/formancehq/ledger/pkg/storage/ledgerstore"
	"github.com/formancehq/stack/libs/go-libs/collectionutils"
	"github.com/formancehq/stack/libs/go-libs/metadata"
	"github.com/stretchr/testify/require"
)

func TestInitializeStore(t *testing.T) {
	t.Parallel()
	store := newLedgerStore(t)

	modified, err := store.Migrate(context.Background())
	require.NoError(t, err)
	require.False(t, modified)
}

func insertTransactions(ctx context.Context, s *ledgerstore.Store, txs ...core.Transaction) error {
	logs := collectionutils.Map(txs, func(from core.Transaction) *core.ChainedLog {
		return core.NewTransactionLog(&from, map[string]metadata.Metadata{}).ChainLog(nil)
	})
	return s.InsertLogs(ctx, logs...)
}
