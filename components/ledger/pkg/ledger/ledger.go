package ledger

import (
	"context"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/formancehq/ledger/pkg/core"
	"github.com/formancehq/ledger/pkg/ledger/command"
	"github.com/formancehq/ledger/pkg/opentelemetry/metrics"
	"github.com/formancehq/ledger/pkg/storage/ledgerstore"
	"github.com/formancehq/stack/libs/go-libs/api"
	"github.com/formancehq/stack/libs/go-libs/logging"
	"github.com/formancehq/stack/libs/go-libs/metadata"
	"github.com/pkg/errors"
)

type Ledger struct {
	commander              *command.Commander
	store                  *ledgerstore.Store
	updateVolumesPeriodic  *periodic
	updateEffectiveVolumes *periodic
}

func (l *Ledger) CreateTransaction(ctx context.Context, parameters command.Parameters, data core.RunScript) (*core.Transaction, error) {
	return l.commander.CreateTransaction(ctx, parameters, data)
}

func (l *Ledger) RevertTransaction(ctx context.Context, parameters command.Parameters, id uint64) (*core.Transaction, error) {
	return l.commander.RevertTransaction(ctx, parameters, id)
}

func (l *Ledger) SaveMeta(ctx context.Context, parameters command.Parameters, targetType string, targetID any, m metadata.Metadata) error {
	return l.commander.SaveMeta(ctx, parameters, targetType, targetID, m)
}

func New(
	name string,
	store *ledgerstore.Store,
	publisher message.Publisher,
	compiler *command.Compiler,
) *Ledger {
	//TODO: reimplements
	//var monitor Monitor = NewNoOpMonitor()
	//if publisher != nil {
	//	monitor = bus.NewLedgerMonitor(publisher, name)
	//}
	metricsRegistry, err := metrics.RegisterPerLedgerMetricsRegistry(name)
	if err != nil {
		panic(err)
	}
	return &Ledger{
		commander: command.New(
			store,
			command.NewDefaultLocker(),
			compiler,
			command.NewReferencer(),
			metricsRegistry,
		),
		store:                  store,
		updateVolumesPeriodic:  newPeriodic(store.UpdateVolumes),
		updateEffectiveVolumes: newPeriodic(store.UpdateEffectiveVolumes),
	}
}

func (l *Ledger) Start(ctx context.Context) {
	go l.commander.Run(logging.ContextWithField(ctx, "component", "commander"))
	go l.updateVolumesPeriodic.Run(logging.ContextWithField(ctx, "component", "volumes updater"))
	go l.updateEffectiveVolumes.Run(logging.ContextWithField(ctx, "component", "effective volumes updater"))
}

func (l *Ledger) Close(ctx context.Context) {
	logging.FromContext(ctx).Debugf("Close commander")
	l.commander.Close()

	logging.FromContext(ctx).Debugf("Close volumes updater")
	l.updateVolumesPeriodic.Stop()

	logging.FromContext(ctx).Debugf("Close effective volumes updater")
	l.updateEffectiveVolumes.Stop()
}

func (l *Ledger) GetTransactions(ctx context.Context, q ledgerstore.TransactionsQuery) (*api.Cursor[core.ExpandedTransaction], error) {
	txs, err := l.store.GetTransactions(ctx, q)
	return txs, errors.Wrap(err, "getting transactions")
}

func (l *Ledger) CountTransactions(ctx context.Context, q ledgerstore.TransactionsQuery) (uint64, error) {
	count, err := l.store.CountTransactions(ctx, q)
	return count, errors.Wrap(err, "counting transactions")
}

func (l *Ledger) GetTransactionWithVolumes(ctx context.Context, id uint64, expandVolumes, expandEffectiveVolumes bool) (*core.ExpandedTransaction, error) {
	tx, err := l.store.GetTransactionWithVolumes(ctx, id, expandVolumes, expandEffectiveVolumes)
	return tx, errors.Wrap(err, "getting transaction")
}

func (l *Ledger) CountAccounts(ctx context.Context, a ledgerstore.AccountsQuery) (uint64, error) {
	count, err := l.store.CountAccounts(ctx, a)
	return count, errors.Wrap(err, "counting accounts")
}

func (l *Ledger) GetAccounts(ctx context.Context, a ledgerstore.AccountsQuery) (*api.Cursor[core.Account], error) {
	accounts, err := l.store.GetAccounts(ctx, a)
	return accounts, errors.Wrap(err, "getting accounts")
}

func (l *Ledger) GetAccountWithVolumes(ctx context.Context, address string, expandVolumes, expandEffectiveVolumes bool) (*core.AccountWithVolumes, error) {
	accounts, err := l.store.GetAccountWithVolumes(ctx, address, expandVolumes, expandEffectiveVolumes)
	return accounts, errors.Wrap(err, "getting account")
}

func (l *Ledger) GetBalances(ctx context.Context, q ledgerstore.BalancesQuery) (*api.Cursor[core.BalancesByAssetsByAccounts], error) {
	balances, err := l.store.GetBalances(ctx, q)
	return balances, errors.Wrap(err, "getting balances")
}

func (l *Ledger) GetBalancesAggregated(ctx context.Context, q ledgerstore.BalancesQuery) (core.BalancesByAssets, error) {
	balances, err := l.store.GetAggregatedBalances(ctx, q)
	return balances, errors.Wrap(err, "getting balances aggregated")
}

func (l *Ledger) GetLogs(ctx context.Context, q ledgerstore.LogsQuery) (*api.Cursor[core.ChainedLog], error) {
	logs, err := l.store.GetLogs(ctx, q)
	return logs, errors.Wrap(err, "getting logs")
}
