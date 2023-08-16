// Code generated by MockGen. DO NOT EDIT.
// Source: api.go

// Package api_test is a generated GoMock package.
package api_test

import (
	context "context"
	reflect "reflect"

	internal "github.com/formancehq/ledger/internal"
	api "github.com/formancehq/ledger/internal/api"
	engine "github.com/formancehq/ledger/internal/engine"
	command "github.com/formancehq/ledger/internal/engine/command"
	ledgerstore "github.com/formancehq/ledger/internal/storage/ledgerstore"
	api0 "github.com/formancehq/stack/libs/go-libs/api"
	metadata "github.com/formancehq/stack/libs/go-libs/metadata"
	migrations "github.com/formancehq/stack/libs/go-libs/migrations"
	gomock "github.com/golang/mock/gomock"
)

// MockLedger is a mock of Ledger interface.
type MockLedger struct {
	ctrl     *gomock.Controller
	recorder *MockLedgerMockRecorder
}

// MockLedgerMockRecorder is the mock recorder for MockLedger.
type MockLedgerMockRecorder struct {
	mock *MockLedger
}

// NewMockLedger creates a new mock instance.
func NewMockLedger(ctrl *gomock.Controller) *MockLedger {
	mock := &MockLedger{ctrl: ctrl}
	mock.recorder = &MockLedgerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockLedger) EXPECT() *MockLedgerMockRecorder {
	return m.recorder
}

// CountAccounts mocks base method.
func (m *MockLedger) CountAccounts(ctx context.Context, query ledgerstore.GetAccountsQuery) (uint64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CountAccounts", ctx, query)
	ret0, _ := ret[0].(uint64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CountAccounts indicates an expected call of CountAccounts.
func (mr *MockLedgerMockRecorder) CountAccounts(ctx, query interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CountAccounts", reflect.TypeOf((*MockLedger)(nil).CountAccounts), ctx, query)
}

// CountTransactions mocks base method.
func (m *MockLedger) CountTransactions(ctx context.Context, query ledgerstore.GetTransactionsQuery) (uint64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CountTransactions", ctx, query)
	ret0, _ := ret[0].(uint64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CountTransactions indicates an expected call of CountTransactions.
func (mr *MockLedgerMockRecorder) CountTransactions(ctx, query interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CountTransactions", reflect.TypeOf((*MockLedger)(nil).CountTransactions), ctx, query)
}

// CreateTransaction mocks base method.
func (m *MockLedger) CreateTransaction(ctx context.Context, parameters command.Parameters, data internal.RunScript) (*internal.Transaction, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateTransaction", ctx, parameters, data)
	ret0, _ := ret[0].(*internal.Transaction)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateTransaction indicates an expected call of CreateTransaction.
func (mr *MockLedgerMockRecorder) CreateTransaction(ctx, parameters, data interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateTransaction", reflect.TypeOf((*MockLedger)(nil).CreateTransaction), ctx, parameters, data)
}

// GetAccountWithVolumes mocks base method.
func (m *MockLedger) GetAccountWithVolumes(ctx context.Context, query ledgerstore.GetAccountQuery) (*internal.ExpandedAccount, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAccountWithVolumes", ctx, query)
	ret0, _ := ret[0].(*internal.ExpandedAccount)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAccountWithVolumes indicates an expected call of GetAccountWithVolumes.
func (mr *MockLedgerMockRecorder) GetAccountWithVolumes(ctx, query interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAccountWithVolumes", reflect.TypeOf((*MockLedger)(nil).GetAccountWithVolumes), ctx, query)
}

// GetAccountsWithVolumes mocks base method.
func (m *MockLedger) GetAccountsWithVolumes(ctx context.Context, query ledgerstore.GetAccountsQuery) (*api0.Cursor[internal.ExpandedAccount], error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAccountsWithVolumes", ctx, query)
	ret0, _ := ret[0].(*api0.Cursor[internal.ExpandedAccount])
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAccountsWithVolumes indicates an expected call of GetAccountsWithVolumes.
func (mr *MockLedgerMockRecorder) GetAccountsWithVolumes(ctx, query interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAccountsWithVolumes", reflect.TypeOf((*MockLedger)(nil).GetAccountsWithVolumes), ctx, query)
}

// GetAggregatedBalances mocks base method.
func (m *MockLedger) GetAggregatedBalances(ctx context.Context, q ledgerstore.GetAggregatedBalancesQuery) (internal.BalancesByAssets, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAggregatedBalances", ctx, q)
	ret0, _ := ret[0].(internal.BalancesByAssets)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAggregatedBalances indicates an expected call of GetAggregatedBalances.
func (mr *MockLedgerMockRecorder) GetAggregatedBalances(ctx, q interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAggregatedBalances", reflect.TypeOf((*MockLedger)(nil).GetAggregatedBalances), ctx, q)
}

// GetLogs mocks base method.
func (m *MockLedger) GetLogs(ctx context.Context, query ledgerstore.GetLogsQuery) (*api0.Cursor[internal.ChainedLog], error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetLogs", ctx, query)
	ret0, _ := ret[0].(*api0.Cursor[internal.ChainedLog])
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetLogs indicates an expected call of GetLogs.
func (mr *MockLedgerMockRecorder) GetLogs(ctx, query interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetLogs", reflect.TypeOf((*MockLedger)(nil).GetLogs), ctx, query)
}

// GetMigrationsInfo mocks base method.
func (m *MockLedger) GetMigrationsInfo(ctx context.Context) ([]migrations.Info, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetMigrationsInfo", ctx)
	ret0, _ := ret[0].([]migrations.Info)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetMigrationsInfo indicates an expected call of GetMigrationsInfo.
func (mr *MockLedgerMockRecorder) GetMigrationsInfo(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetMigrationsInfo", reflect.TypeOf((*MockLedger)(nil).GetMigrationsInfo), ctx)
}

// GetTransactionWithVolumes mocks base method.
func (m *MockLedger) GetTransactionWithVolumes(ctx context.Context, query ledgerstore.GetTransactionQuery) (*internal.ExpandedTransaction, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTransactionWithVolumes", ctx, query)
	ret0, _ := ret[0].(*internal.ExpandedTransaction)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTransactionWithVolumes indicates an expected call of GetTransactionWithVolumes.
func (mr *MockLedgerMockRecorder) GetTransactionWithVolumes(ctx, query interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTransactionWithVolumes", reflect.TypeOf((*MockLedger)(nil).GetTransactionWithVolumes), ctx, query)
}

// GetTransactions mocks base method.
func (m *MockLedger) GetTransactions(ctx context.Context, query ledgerstore.GetTransactionsQuery) (*api0.Cursor[internal.ExpandedTransaction], error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTransactions", ctx, query)
	ret0, _ := ret[0].(*api0.Cursor[internal.ExpandedTransaction])
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTransactions indicates an expected call of GetTransactions.
func (mr *MockLedgerMockRecorder) GetTransactions(ctx, query interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTransactions", reflect.TypeOf((*MockLedger)(nil).GetTransactions), ctx, query)
}

// RevertTransaction mocks base method.
func (m *MockLedger) RevertTransaction(ctx context.Context, parameters command.Parameters, id uint64) (*internal.Transaction, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RevertTransaction", ctx, parameters, id)
	ret0, _ := ret[0].(*internal.Transaction)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// RevertTransaction indicates an expected call of RevertTransaction.
func (mr *MockLedgerMockRecorder) RevertTransaction(ctx, parameters, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RevertTransaction", reflect.TypeOf((*MockLedger)(nil).RevertTransaction), ctx, parameters, id)
}

// SaveMeta mocks base method.
func (m_2 *MockLedger) SaveMeta(ctx context.Context, parameters command.Parameters, targetType string, targetID any, m metadata.Metadata) error {
	m_2.ctrl.T.Helper()
	ret := m_2.ctrl.Call(m_2, "SaveMeta", ctx, parameters, targetType, targetID, m)
	ret0, _ := ret[0].(error)
	return ret0
}

// SaveMeta indicates an expected call of SaveMeta.
func (mr *MockLedgerMockRecorder) SaveMeta(ctx, parameters, targetType, targetID, m interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveMeta", reflect.TypeOf((*MockLedger)(nil).SaveMeta), ctx, parameters, targetType, targetID, m)
}

// Stats mocks base method.
func (m *MockLedger) Stats(ctx context.Context) (engine.Stats, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Stats", ctx)
	ret0, _ := ret[0].(engine.Stats)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Stats indicates an expected call of Stats.
func (mr *MockLedgerMockRecorder) Stats(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Stats", reflect.TypeOf((*MockLedger)(nil).Stats), ctx)
}

// MockBackend is a mock of Backend interface.
type MockBackend struct {
	ctrl     *gomock.Controller
	recorder *MockBackendMockRecorder
}

// MockBackendMockRecorder is the mock recorder for MockBackend.
type MockBackendMockRecorder struct {
	mock *MockBackend
}

// NewMockBackend creates a new mock instance.
func NewMockBackend(ctrl *gomock.Controller) *MockBackend {
	mock := &MockBackend{ctrl: ctrl}
	mock.recorder = &MockBackendMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockBackend) EXPECT() *MockBackendMockRecorder {
	return m.recorder
}

// GetLedger mocks base method.
func (m *MockBackend) GetLedger(ctx context.Context, name string) (api.Ledger, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetLedger", ctx, name)
	ret0, _ := ret[0].(api.Ledger)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetLedger indicates an expected call of GetLedger.
func (mr *MockBackendMockRecorder) GetLedger(ctx, name interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetLedger", reflect.TypeOf((*MockBackend)(nil).GetLedger), ctx, name)
}

// GetVersion mocks base method.
func (m *MockBackend) GetVersion() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetVersion")
	ret0, _ := ret[0].(string)
	return ret0
}

// GetVersion indicates an expected call of GetVersion.
func (mr *MockBackendMockRecorder) GetVersion() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetVersion", reflect.TypeOf((*MockBackend)(nil).GetVersion))
}

// ListLedgers mocks base method.
func (m *MockBackend) ListLedgers(ctx context.Context) ([]string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListLedgers", ctx)
	ret0, _ := ret[0].([]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListLedgers indicates an expected call of ListLedgers.
func (mr *MockBackendMockRecorder) ListLedgers(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListLedgers", reflect.TypeOf((*MockBackend)(nil).ListLedgers), ctx)
}
