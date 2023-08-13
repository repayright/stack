package engine

import (
	"context"

	"github.com/formancehq/ledger/internal"
	"github.com/formancehq/stack/libs/go-libs/metadata"
)

type Monitor interface {
	CommittedTransactions(ctx context.Context, res ...ledger.Transaction)
	SavedMetadata(ctx context.Context, targetType, id string, metadata metadata.Metadata)
	RevertedTransaction(ctx context.Context, reverted, revert *ledger.Transaction)
}

type noOpMonitor struct{}

func (n noOpMonitor) CommittedTransactions(ctx context.Context, res ...ledger.Transaction) {
}
func (n noOpMonitor) SavedMetadata(ctx context.Context, targetType string, id string, metadata metadata.Metadata) {
}
func (n noOpMonitor) RevertedTransaction(ctx context.Context, reverted, revert *ledger.Transaction) {
}

var _ Monitor = &noOpMonitor{}

func NewNoOpMonitor() *noOpMonitor {
	return &noOpMonitor{}
}
