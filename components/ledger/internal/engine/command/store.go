package command

import (
	"context"

	"github.com/formancehq/ledger/internal"
	"github.com/formancehq/ledger/internal/machine/vm"
)

type Store interface {
	vm.Store
	InsertLogs(ctx context.Context, logs ...*ledger.ChainedLog) error
	GetLastLog(ctx context.Context) (*ledger.ChainedLog, error)
	ReadLogWithIdempotencyKey(ctx context.Context, key string) (*ledger.ChainedLog, error)
	ReadLastLogWithType(ctx context.Context, logType ...ledger.LogType) (*ledger.ChainedLog, error)
	GetTransactionByReference(ctx context.Context, ref string) (*ledger.ExpandedTransaction, error)
	GetTransaction(ctx context.Context, txID uint64) (*ledger.Transaction, error)
}
