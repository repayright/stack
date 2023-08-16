package ledgerstore

import (
	"context"
	"database/sql/driver"
	"encoding/json"
	"math/big"
	"strings"

	ledger "github.com/formancehq/ledger/internal"
	"github.com/formancehq/ledger/internal/storage/paginate"
	"github.com/formancehq/stack/libs/go-libs/api"
	"github.com/formancehq/stack/libs/go-libs/metadata"
	"github.com/uptrace/bun"
)

const (
	MovesTableName = "moves"
)

type Transaction struct {
	bun.BaseModel `bun:"transactions,alias:transactions"`

	ID                         *paginate.BigInt             `bun:"id,type:numeric"`
	Timestamp                  ledger.Time                  `bun:"timestamp,type:timestamp without time zone"`
	Reference                  string                       `bun:"reference,type:varchar,unique,nullzero"`
	Postings                   []ledger.Posting             `bun:"postings,type:jsonb"`
	Metadata                   metadata.Metadata            `bun:"metadata,type:jsonb,default:'{}'"`
	PostCommitEffectiveVolumes ledger.AccountsAssetsVolumes `bun:"post_commit_effective_volumes,type:jsonb"`
	PostCommitVolumes          ledger.AccountsAssetsVolumes `bun:"post_commit_volumes,type:jsonb"`
	Reverted                   bool                         `bun:"reverted"`
	Revision                   int                          `bun:"revision"`
	LastUpdate                 *ledger.Time                 `bun:"last_update"`
}

func (t *Transaction) toCore() *ledger.ExpandedTransaction {
	var (
		preCommitEffectiveVolumes ledger.AccountsAssetsVolumes
		preCommitVolumes          ledger.AccountsAssetsVolumes
	)
	if t.PostCommitEffectiveVolumes != nil {
		preCommitEffectiveVolumes = t.PostCommitEffectiveVolumes.Copy()
		for _, posting := range t.Postings {
			preCommitEffectiveVolumes.AddOutput(posting.Source, posting.Asset, big.NewInt(0).Neg(posting.Amount))
			preCommitEffectiveVolumes.AddInput(posting.Destination, posting.Asset, big.NewInt(0).Neg(posting.Amount))
		}
	}
	if t.PostCommitVolumes != nil {
		preCommitVolumes = t.PostCommitVolumes.Copy()
		for _, posting := range t.Postings {
			preCommitVolumes.AddOutput(posting.Source, posting.Asset, big.NewInt(0).Neg(posting.Amount))
			preCommitVolumes.AddInput(posting.Destination, posting.Asset, big.NewInt(0).Neg(posting.Amount))
		}
	}
	return &ledger.ExpandedTransaction{
		Transaction: ledger.Transaction{
			TransactionData: ledger.TransactionData{
				Reference: t.Reference,
				Metadata:  t.Metadata,
				Timestamp: t.Timestamp,
				Postings:  t.Postings,
			},
			ID:       (*big.Int)(t.ID),
			Reverted: t.Reverted,
		},
		PreCommitEffectiveVolumes:  preCommitEffectiveVolumes,
		PostCommitEffectiveVolumes: t.PostCommitEffectiveVolumes,
		PreCommitVolumes:           preCommitVolumes,
		PostCommitVolumes:          t.PostCommitVolumes,
	}
}

type account string

var _ driver.Valuer = account("")

func (m1 account) Value() (driver.Value, error) {
	ret, err := json.Marshal(strings.Split(string(m1), ":"))
	if err != nil {
		return nil, err
	}
	return string(ret), nil
}

// Scan - Implement the database/sql scanner interface
func (m1 *account) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	v, err := driver.String.ConvertValue(value)
	if err != nil {
		return err
	}

	array := make([]string, 0)
	switch vv := v.(type) {
	case []uint8:
		err = json.Unmarshal(vv, &array)
	case string:
		err = json.Unmarshal([]byte(vv), &array)
	default:
		panic("not handled type")
	}
	if err != nil {
		return err
	}
	*m1 = account(strings.Join(array, ":"))
	return nil
}

func (store *Store) transactionQuery(p PITFilter) func(query *bun.SelectQuery) *bun.SelectQuery {
	return func(query *bun.SelectQuery) *bun.SelectQuery {
		subQuery := query.NewSelect().
			Table("transactions").
			ColumnExpr("distinct on(transactions.id) transactions.*").
			Apply(filterPIT(p.PIT, "transactions.timestamp")).
			OrderExpr("transactions.id desc, revision desc")

		query = query.
			TableExpr("(" + subQuery.String() + ") transactions").
			ColumnExpr("transactions.*")
		if p.ExpandEffectiveVolumes {
			query = query.ColumnExpr("get_aggregated_effective_volumes_for_transaction(transactions) as post_commit_effective_volumes")
		}
		if p.ExpandVolumes {
			query = query.ColumnExpr("get_aggregated_volumes_for_transaction(transactions) as post_commit_volumes")
		}
		return query
	}
}

func (store *Store) transactionListBuilder(p GetTransactionsQuery) func(query *bun.SelectQuery) *bun.SelectQuery {
	return func(query *bun.SelectQuery) *bun.SelectQuery {
		query = store.transactionQuery(p.Options.PITFilter)(query).
			Apply(filterMetadata(p.Options.Metadata))

		if p.Options.Reference != "" {
			query.Where("transactions.reference = ?", p.Options.Reference)
		}
		if !p.Options.StartTime.IsZero() {
			query.Where("transactions.timestamp >= ?", p.Options.StartTime)
		}
		if !p.Options.EndTime.IsZero() {
			query.Where("transactions.timestamp < ?", p.Options.EndTime)
		}
		if p.Options.AfterTxID != 0 {
			query.Where("transactions.id > ?", p.Options.AfterTxID)
		}
		if p.Options.Source != "" || p.Options.Destination != "" || p.Options.Account != "" {
			query.Join("join moves m on transactions.id = m.transaction_id")
			if p.Options.Source != "" {
				query = query.
					Where("m.is_source").
					Apply(filterAccountAddressBuilder(p.Options.Source, "account_address"))

			}
			if p.Options.Destination != "" {
				query = query.
					Where("not m.is_source").
					Apply(filterAccountAddressBuilder(p.Options.Destination, "account_address"))
			}
			if p.Options.Account != "" {
				query = query.Apply(filterAccountAddressBuilder(p.Options.Account, "account_address"))
			}
		}
		return query
	}
}

func (store *Store) GetTransactions(ctx context.Context, q GetTransactionsQuery) (*api.Cursor[ledger.ExpandedTransaction], error) {
	transactions, err := paginateWithColumn[TransactionsQueryOptions, Transaction](store, ctx,
		paginate.ColumnPaginatedQuery[TransactionsQueryOptions](q),
		store.transactionListBuilder(q),
	)
	if err != nil {
		return nil, err
	}
	return api.MapCursor(transactions, func(from Transaction) ledger.ExpandedTransaction {
		return *from.toCore()
	}), nil
}

func (store *Store) CountTransactions(ctx context.Context, q GetTransactionsQuery) (uint64, error) {
	return count(store, ctx, store.transactionListBuilder(q))
}

func (store *Store) GetTransactionWithVolumes(ctx context.Context, filter GetTransactionQuery) (*ledger.ExpandedTransaction, error) {
	return fetchAndMap[*Transaction, *ledger.ExpandedTransaction](store, ctx,
		(*Transaction).toCore,
		func(query *bun.SelectQuery) *bun.SelectQuery {
			return store.transactionQuery(filter.PITFilter)(query).
				Where("id = ?", filter.ID).
				Limit(1)
		})
}

func (store *Store) GetTransaction(ctx context.Context, txId *big.Int) (*ledger.Transaction, error) {
	return fetch[*ledger.Transaction](store, ctx,
		func(query *bun.SelectQuery) *bun.SelectQuery {
			return query.
				Table("transactions").
				ColumnExpr(`transactions.id, transactions.reference, transactions.metadata, transactions.postings, transactions.timestamp, transactions.reverted`).
				Where("id = ?", (*paginate.BigInt)(txId)).
				Order("revision desc").
				Limit(1)
		})
}

func (store *Store) GetTransactionByReference(ctx context.Context, ref string) (*ledger.ExpandedTransaction, error) {
	return fetchAndMap[*Transaction, *ledger.ExpandedTransaction](store, ctx,
		(*Transaction).toCore,
		func(query *bun.SelectQuery) *bun.SelectQuery {
			return query.
				Table("transactions").
				ColumnExpr(`transactions.id, transactions.reference, transactions.metadata, transactions.postings, transactions.timestamp, transactions.reverted`).
				Where("reference = ?", ref).
				Order("revision desc").
				Limit(1)
		})
}

type GetTransactionsQuery paginate.ColumnPaginatedQuery[TransactionsQueryOptions]

func NewTransactionsQuery() GetTransactionsQuery {
	return GetTransactionsQuery{
		PageSize: paginate.QueryDefaultPageSize,
		Column:   "id",
		Order:    paginate.OrderDesc,
		Options: TransactionsQueryOptions{
			Metadata: metadata.Metadata{},
		},
	}
}

type TransactionsQueryOptions struct {
	PITFilter
	AfterTxID   uint64            `json:"afterTxID,omitempty"`
	Reference   string            `json:"reference,omitempty"`
	Destination string            `json:"destination,omitempty"`
	Source      string            `json:"source,omitempty"`
	Account     string            `json:"account,omitempty"`
	EndTime     ledger.Time       `json:"endTime,omitempty"`
	StartTime   ledger.Time       `json:"startTime,omitempty"`
	Metadata    metadata.Metadata `json:"metadata,omitempty"`
}

func (a GetTransactionsQuery) WithPageSize(pageSize uint64) GetTransactionsQuery {
	if pageSize != 0 {
		a.PageSize = pageSize
	}

	return a
}

func (a GetTransactionsQuery) WithAfterTxID(after uint64) GetTransactionsQuery {
	a.Options.AfterTxID = after

	return a
}

func (a GetTransactionsQuery) WithStartTimeFilter(start ledger.Time) GetTransactionsQuery {
	if !start.IsZero() {
		a.Options.StartTime = start
	}

	return a
}

func (a GetTransactionsQuery) WithEndTimeFilter(end ledger.Time) GetTransactionsQuery {
	if !end.IsZero() {
		a.Options.EndTime = end
	}

	return a
}

func (a GetTransactionsQuery) WithAccountFilter(account string) GetTransactionsQuery {
	a.Options.Account = account

	return a
}

func (a GetTransactionsQuery) WithDestinationFilter(dest string) GetTransactionsQuery {
	a.Options.Destination = dest

	return a
}

func (a GetTransactionsQuery) WithReferenceFilter(ref string) GetTransactionsQuery {
	a.Options.Reference = ref

	return a
}

func (a GetTransactionsQuery) WithSourceFilter(source string) GetTransactionsQuery {
	a.Options.Source = source

	return a
}

func (a GetTransactionsQuery) WithMetadataFilter(metadata metadata.Metadata) GetTransactionsQuery {
	a.Options.Metadata = metadata

	return a
}

func (a GetTransactionsQuery) WithExpandEffectiveVolumes(v bool) GetTransactionsQuery {
	a.Options.ExpandEffectiveVolumes = v

	return a
}

func (a GetTransactionsQuery) WithExpandVolumes(v bool) GetTransactionsQuery {
	a.Options.ExpandVolumes = v

	return a
}

type GetTransactionQuery struct {
	PITFilter
	ID *big.Int
}

func (q GetTransactionQuery) WithExpandVolumes() GetTransactionQuery {
	q.ExpandVolumes = true

	return q
}

func (q GetTransactionQuery) WithExpandEffectiveVolumes() GetTransactionQuery {
	q.ExpandEffectiveVolumes = true

	return q
}

func NewGetTransactionQuery(id *big.Int) GetTransactionQuery {
	return GetTransactionQuery{
		PITFilter: PITFilter{},
		ID:        id,
	}
}
