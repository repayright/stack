package command

import (
	"context"

	"github.com/formancehq/ledger/pkg/core"
	"github.com/formancehq/ledger/pkg/machine/vm"
)

type Store interface {
	vm.Store
	InsertLogs(ctx context.Context, logs ...*core.ChainedLog) error
	GetLastLog(ctx context.Context) (*core.ChainedLog, error)
	ReadLogWithIdempotencyKey(ctx context.Context, key string) (*core.ChainedLog, error)
	ReadLastLogWithType(ctx context.Context, logType ...core.LogType) (*core.ChainedLog, error)
	GetTransactionByReference(ctx context.Context, ref string) (*core.ExpandedTransaction, error)
	GetTransaction(ctx context.Context, txID uint64) (*core.Transaction, error)
}
