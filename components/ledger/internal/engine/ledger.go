package engine

import (
	"context"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/formancehq/ledger/internal"
	command2 "github.com/formancehq/ledger/internal/engine/command"
	ledgerstore2 "github.com/formancehq/ledger/internal/storage/ledgerstore"
	"github.com/formancehq/stack/libs/go-libs/api"
	"github.com/formancehq/stack/libs/go-libs/logging"
	"github.com/formancehq/stack/libs/go-libs/metadata"
	"github.com/pkg/errors"
)

type Ledger struct {
	commander              *command2.Commander
	store                  *ledgerstore2.Store
	updateVolumesPeriodic  *periodic
	updateEffectiveVolumes *periodic
}

func (l *Ledger) CreateTransaction(ctx context.Context, parameters command2.Parameters, data ledger.RunScript) (*ledger.Transaction, error) {
	return l.commander.CreateTransaction(ctx, parameters, data)
}

func (l *Ledger) RevertTransaction(ctx context.Context, parameters command2.Parameters, id uint64) (*ledger.Transaction, error) {
	return l.commander.RevertTransaction(ctx, parameters, id)
}

func (l *Ledger) SaveMeta(ctx context.Context, parameters command2.Parameters, targetType string, targetID any, m metadata.Metadata) error {
	return l.commander.SaveMeta(ctx, parameters, targetType, targetID, m)
}

func New(
	store *ledgerstore2.Store,
	publisher message.Publisher,
	compiler *command2.Compiler,
) *Ledger {
	//TODO: reimplements
	//var monitor Monitor = NewNoOpMonitor()
	//if publisher != nil {
	//	monitor = bus.NewLedgerMonitor(publisher, name)
	//}
	return &Ledger{
		commander: command2.New(
			store,
			command2.NewDefaultLocker(),
			compiler,
			command2.NewReferencer(),
		),
		store: store,
	}
}

func (l *Ledger) Start(ctx context.Context) {
	go l.commander.Run(logging.ContextWithField(ctx, "component", "commander"))
}

func (l *Ledger) Close(ctx context.Context) {
	logging.FromContext(ctx).Debugf("Close commander")
	l.commander.Close()
}

func (l *Ledger) GetTransactions(ctx context.Context, q ledgerstore2.GetTransactionsQuery) (*api.Cursor[ledger.ExpandedTransaction], error) {
	txs, err := l.store.GetTransactions(ctx, q)
	return txs, errors.Wrap(err, "getting transactions")
}

func (l *Ledger) CountTransactions(ctx context.Context, q ledgerstore2.GetTransactionsQuery) (uint64, error) {
	count, err := l.store.CountTransactions(ctx, q)
	return count, errors.Wrap(err, "counting transactions")
}

func (l *Ledger) GetTransactionWithVolumes(ctx context.Context, query ledgerstore2.GetTransactionQuery) (*ledger.ExpandedTransaction, error) {
	tx, err := l.store.GetTransactionWithVolumes(ctx, query)
	return tx, errors.Wrap(err, "getting transaction")
}

func (l *Ledger) CountAccounts(ctx context.Context, a ledgerstore2.GetAccountsQuery) (uint64, error) {
	count, err := l.store.CountAccounts(ctx, a)
	return count, errors.Wrap(err, "counting accounts")
}

func (l *Ledger) GetAccountsWithVolumes(ctx context.Context, a ledgerstore2.GetAccountsQuery) (*api.Cursor[ledger.ExpandedAccount], error) {
	accounts, err := l.store.GetAccountsWithVolumes(ctx, a)
	return accounts, errors.Wrap(err, "getting accounts")
}

func (l *Ledger) GetAccountWithVolumes(ctx context.Context, q ledgerstore2.GetAccountQuery) (*ledger.ExpandedAccount, error) {
	accounts, err := l.store.GetAccountWithVolumes(ctx, q)
	return accounts, errors.Wrap(err, "getting account")
}

func (l *Ledger) GetAggregatedBalances(ctx context.Context, q ledgerstore2.GetAggregatedBalancesQuery) (ledger.BalancesByAssets, error) {
	balances, err := l.store.GetAggregatedBalances(ctx, q)
	return balances, errors.Wrap(err, "getting balances aggregated")
}

func (l *Ledger) GetLogs(ctx context.Context, q ledgerstore2.GetLogsQuery) (*api.Cursor[ledger.ChainedLog], error) {
	logs, err := l.store.GetLogs(ctx, q)
	return logs, errors.Wrap(err, "getting logs")
}
