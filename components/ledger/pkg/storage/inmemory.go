package storage

import (
	"context"
	"math/big"

	"github.com/formancehq/ledger/pkg/core"
	"github.com/formancehq/stack/libs/go-libs/collectionutils"
	"github.com/formancehq/stack/libs/go-libs/metadata"
)

type InMemoryStore struct {
	logs         []*core.ChainedLog
	transactions []*core.ExpandedTransaction
	accounts     []*core.Account
}

func (m *InMemoryStore) GetTransactionByReference(ctx context.Context, ref string) (*core.ExpandedTransaction, error) {
	filtered := collectionutils.Filter(m.transactions, func(transaction *core.ExpandedTransaction) bool {
		return transaction.Reference == ref
	})
	if len(filtered) == 0 {
		return nil, ErrNotFound
	}
	return filtered[0], nil
}

func (m *InMemoryStore) GetTransaction(ctx context.Context, txID uint64) (*core.Transaction, error) {
	filtered := collectionutils.Filter(m.transactions, func(transaction *core.ExpandedTransaction) bool {
		return transaction.ID == txID
	})
	if len(filtered) == 0 {
		return nil, ErrNotFound
	}
	return &filtered[0].Transaction, nil
}

func (m *InMemoryStore) GetLastLog(ctx context.Context) (*core.ChainedLog, error) {
	if len(m.logs) == 0 {
		return nil, nil
	}
	return m.logs[len(m.logs)-1], nil
}

func (m *InMemoryStore) GetBalance(ctx context.Context, address, asset string) (*big.Int, error) {
	balance := new(big.Int)
	for _, log := range m.logs {
		switch payload := log.Data.(type) {
		case core.NewTransactionLogPayload:
			postings := payload.Transaction.Postings
			for _, posting := range postings {
				if posting.Asset != asset {
					continue
				}
				if posting.Source == address {
					balance = balance.Sub(balance, posting.Amount)
				}
				if posting.Destination == address {
					balance = balance.Add(balance, posting.Amount)
				}
			}
		}
	}
	return balance, nil
}

func (m *InMemoryStore) GetAccount(ctx context.Context, address string) (*core.Account, error) {
	account := collectionutils.Filter(m.accounts, func(account *core.Account) bool {
		return account.Address == address
	})
	if len(account) == 0 {
		return &core.Account{
			Address:  address,
			Metadata: metadata.Metadata{},
		}, nil
	}
	return account[0], nil
}

func (m *InMemoryStore) ReadLogWithIdempotencyKey(ctx context.Context, key string) (*core.ChainedLog, error) {
	first := collectionutils.First(m.logs, func(log *core.ChainedLog) bool {
		return log.IdempotencyKey == key
	})
	if first == nil {
		return nil, ErrNotFound
	}
	return first, nil
}

func (m *InMemoryStore) ReadLastLogWithType(background context.Context, logType ...core.LogType) (*core.ChainedLog, error) {
	first := collectionutils.First(m.logs, func(log *core.ChainedLog) bool {
		return collectionutils.Contains(logType, log.Type)
	})
	if first == nil {
		return nil, ErrNotFound
	}
	return first, nil
}

func (m *InMemoryStore) InsertLogs(ctx context.Context, logs ...*core.ChainedLog) error {
	m.logs = append(m.logs, logs...)
	for _, log := range logs {
		switch payload := log.Data.(type) {
		case core.NewTransactionLogPayload:
			m.transactions = append(m.transactions, &core.ExpandedTransaction{
				Transaction: *payload.Transaction,
				// TODO
				PreCommitVolumes:  nil,
				PostCommitVolumes: nil,
			})
		case core.RevertedTransactionLogPayload:
			tx := collectionutils.Filter(m.transactions, func(transaction *core.ExpandedTransaction) bool {
				return transaction.ID == payload.RevertedTransactionID
			})[0]
			tx.Reverted = true
		case core.SetMetadataLogPayload:
		}
	}

	return nil
}

func NewInMemoryStore() *InMemoryStore {
	return &InMemoryStore{
		logs: []*core.ChainedLog{},
	}
}
