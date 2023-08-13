package api

import (
	"context"

	"github.com/formancehq/ledger/internal"
	ledger2 "github.com/formancehq/ledger/internal/engine"
	"github.com/formancehq/ledger/internal/engine/command"
	"github.com/formancehq/ledger/internal/storage/driver"
	ledgerstore2 "github.com/formancehq/ledger/internal/storage/ledgerstore"
	"github.com/formancehq/stack/libs/go-libs/api"
	"github.com/formancehq/stack/libs/go-libs/metadata"
	"github.com/formancehq/stack/libs/go-libs/migrations"
)

//go:generate mockgen -source api.go -destination api_test.go -package controllers_test . Ledger

type Ledger interface {
	GetAccountWithVolumes(ctx context.Context, query ledgerstore2.GetAccountQuery) (*ledger.ExpandedAccount, error)
	GetAccountsWithVolumes(ctx context.Context, query ledgerstore2.GetAccountsQuery) (*api.Cursor[ledger.ExpandedAccount], error)
	CountAccounts(ctx context.Context, query ledgerstore2.GetAccountsQuery) (uint64, error)
	GetAggregatedBalances(ctx context.Context, q ledgerstore2.GetAggregatedBalancesQuery) (ledger.BalancesByAssets, error)
	GetMigrationsInfo(ctx context.Context) ([]migrations.Info, error)
	Stats(ctx context.Context) (ledger2.Stats, error)
	GetLogs(ctx context.Context, query ledgerstore2.GetLogsQuery) (*api.Cursor[ledger.ChainedLog], error)
	CountTransactions(ctx context.Context, query ledgerstore2.GetTransactionsQuery) (uint64, error)
	GetTransactions(ctx context.Context, query ledgerstore2.GetTransactionsQuery) (*api.Cursor[ledger.ExpandedTransaction], error)
	GetTransactionWithVolumes(ctx context.Context, query ledgerstore2.GetTransactionQuery) (*ledger.ExpandedTransaction, error)

	CreateTransaction(ctx context.Context, parameters command.Parameters, data ledger.RunScript) (*ledger.Transaction, error)
	RevertTransaction(ctx context.Context, parameters command.Parameters, id uint64) (*ledger.Transaction, error)
	SaveMeta(ctx context.Context, parameters command.Parameters, targetType string, targetID any, m metadata.Metadata) error
}

type Backend interface {
	GetLedger(ctx context.Context, name string) (Ledger, error)
	ListLedgers(ctx context.Context) ([]string, error)
	GetVersion() string
}

type DefaultBackend struct {
	storageDriver *driver.Driver
	resolver      *ledger2.Resolver
	version       string
}

func (d DefaultBackend) GetLedger(ctx context.Context, name string) (Ledger, error) {
	return d.resolver.GetLedger(ctx, name)
}

func (d DefaultBackend) ListLedgers(ctx context.Context) ([]string, error) {
	return d.storageDriver.GetSystemStore().ListLedgers(ctx)
}

func (d DefaultBackend) GetVersion() string {
	return d.version
}

var _ Backend = (*DefaultBackend)(nil)

func NewDefaultBackend(driver *driver.Driver, version string, resolver *ledger2.Resolver) *DefaultBackend {
	return &DefaultBackend{
		storageDriver: driver,
		resolver:      resolver,
		version:       version,
	}
}
