package ledgerstore

import (
	"context"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"math/big"
	"strings"

	"github.com/formancehq/ledger/pkg/core"
	storageerrors "github.com/formancehq/ledger/pkg/storage"
	"github.com/formancehq/stack/libs/go-libs/api"
	"github.com/formancehq/stack/libs/go-libs/metadata"
	"github.com/formancehq/stack/libs/go-libs/pointer"
	"github.com/uptrace/bun"
)

const (
	TransactionsTableName = "transactions"
	MovesTableName        = "moves"
)

type TransactionsQuery ColumnPaginatedQuery[TransactionsQueryFilters]

func NewTransactionsQuery() TransactionsQuery {
	return TransactionsQuery{
		PageSize: QueryDefaultPageSize,
		Column:   "id",
		Order:    OrderDesc,
		Filters: TransactionsQueryFilters{
			Metadata: metadata.Metadata{},
		},
	}
}

type TransactionsQueryFilters struct {
	AfterTxID   uint64            `json:"afterTxID,omitempty"`
	Reference   string            `json:"reference,omitempty"`
	Destination string            `json:"destination,omitempty"`
	Source      string            `json:"source,omitempty"`
	Account     string            `json:"account,omitempty"`
	EndTime     core.Time         `json:"endTime,omitempty"`
	StartTime   core.Time         `json:"startTime,omitempty"`
	Metadata    metadata.Metadata `json:"metadata,omitempty"`
}

func (a TransactionsQuery) WithPageSize(pageSize uint64) TransactionsQuery {
	if pageSize != 0 {
		a.PageSize = pageSize
	}

	return a
}

func (a TransactionsQuery) WithAfterTxID(after uint64) TransactionsQuery {
	a.Filters.AfterTxID = after

	return a
}

func (a TransactionsQuery) WithStartTimeFilter(start core.Time) TransactionsQuery {
	if !start.IsZero() {
		a.Filters.StartTime = start
	}

	return a
}

func (a TransactionsQuery) WithEndTimeFilter(end core.Time) TransactionsQuery {
	if !end.IsZero() {
		a.Filters.EndTime = end
	}

	return a
}

func (a TransactionsQuery) WithAccountFilter(account string) TransactionsQuery {
	a.Filters.Account = account

	return a
}

func (a TransactionsQuery) WithDestinationFilter(dest string) TransactionsQuery {
	a.Filters.Destination = dest

	return a
}

func (a TransactionsQuery) WithReferenceFilter(ref string) TransactionsQuery {
	a.Filters.Reference = ref

	return a
}

func (a TransactionsQuery) WithSourceFilter(source string) TransactionsQuery {
	a.Filters.Source = source

	return a
}

func (a TransactionsQuery) WithMetadataFilter(metadata metadata.Metadata) TransactionsQuery {
	a.Filters.Metadata = metadata

	return a
}

type Transaction struct {
	bun.BaseModel `bun:"transactions,alias:transactions"`

	ID                          uint64                     `bun:"id,type:bigint,pk"`
	Timestamp                   core.Time                  `bun:"date,type:timestamp without time zone"`
	Reference                   string                     `bun:"reference,type:varchar,unique,nullzero"`
	Postings                    []core.Posting             `bun:"postings,type:jsonb"`
	Metadata                    metadata.Metadata          `bun:"metadata,type:jsonb,default:'{}'"`
	PostCommitAggregatedVolumes core.AccountsAssetsVolumes `bun:"post_commit_aggregated_volumes,type:jsonb"`
}

func (t Transaction) toCore() core.ExpandedTransaction {
	preCommitVolumes := t.PostCommitAggregatedVolumes.Copy()
	for _, posting := range t.Postings {
		preCommitVolumes.AddOutput(posting.Source, posting.Asset, big.NewInt(0).Neg(posting.Amount))
		preCommitVolumes.AddInput(posting.Destination, posting.Asset, big.NewInt(0).Neg(posting.Amount))
	}
	return core.ExpandedTransaction{
		Transaction: core.Transaction{
			TransactionData: core.TransactionData{
				Reference: t.Reference,
				Metadata:  t.Metadata,
				Date:      t.Timestamp,
				Postings:  t.Postings,
			},
			ID: t.ID,
		},
		PreCommitVolumes:  preCommitVolumes,
		PostCommitVolumes: t.PostCommitAggregatedVolumes,
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

func (s *Store) buildTransactionsQuery(p TransactionsQueryFilters, models *[]Transaction) *bun.SelectQuery {

	selectMatchingTransactions := s.schema.NewSelect(TransactionsTableName).
		Model(models).
		ColumnExpr("distinct on(transactions.id) transactions.id").
		ColumnExpr("transactions.reference").
		ColumnExpr("transactions.metadata").
		ColumnExpr("transactions.postings").
		ColumnExpr("transactions.date").
		ColumnExpr("get_aggregated_volumes_for_transaction(transactions) as post_commit_aggregated_volumes")
	if p.Reference != "" {
		selectMatchingTransactions.Where("transactions.reference = ?", p.Reference)
	}
	if !p.StartTime.IsZero() {
		selectMatchingTransactions.Where("transactions.date >= ?", p.StartTime)
	}
	if !p.EndTime.IsZero() {
		selectMatchingTransactions.Where("transactions.date < ?", p.EndTime)
	}
	if p.AfterTxID != 0 {
		selectMatchingTransactions.Where("transactions.id > ?", p.AfterTxID)
	}
	if p.Metadata != nil && len(p.Metadata) > 0 {
		selectMatchingTransactions.Where("transactions.metadata @> ?", p.Metadata)
	}
	if p.Source != "" || p.Destination != "" || p.Account != "" {
		selectMatchingTransactions.Join(fmt.Sprintf("join %s m on transactions.id = m.transaction_id", s.schema.Table("moves")))
		if p.Source != "" {
			parts := strings.Split(p.Source, ":")
			selectMatchingTransactions.Where(fmt.Sprintf("m.is_source and jsonb_array_length(m.account_address_array) = %d", len(parts)))
			for index, segment := range parts {
				if len(segment) == 0 {
					continue
				}
				selectMatchingTransactions.Where(fmt.Sprintf(`m.account_address_array @@ ('$[%d] == "%s"')`, index, segment))
			}
		}
		if p.Destination != "" {
			parts := strings.Split(p.Destination, ":")
			selectMatchingTransactions.Where(fmt.Sprintf("not m.is_source and jsonb_array_length(m.account_address_array) = %d", len(parts)))
			for index, segment := range parts {
				if len(segment) == 0 {
					continue
				}
				selectMatchingTransactions.Where(fmt.Sprintf(`m.account_address_array @@ ('$[%d] == "%s"')`, index, segment))
			}
		}
		if p.Account != "" {
			parts := strings.Split(p.Account, ":")
			selectMatchingTransactions.Where(fmt.Sprintf("jsonb_array_length(m.account_address_array) = %d", len(parts)))
			for index, segment := range parts {
				if len(segment) == 0 {
					continue
				}
				selectMatchingTransactions.Where(fmt.Sprintf(`m.account_address_array @@ ('$[%d] == "%s"')`, index, segment))
			}
		}
	}

	return selectMatchingTransactions
}

func (s *Store) GetTransactions(ctx context.Context, q TransactionsQuery) (*api.Cursor[core.ExpandedTransaction], error) {
	cursor, err := UsingColumn[TransactionsQueryFilters, Transaction](ctx,
		s.buildTransactionsQuery, ColumnPaginatedQuery[TransactionsQueryFilters](q),
	)
	if err != nil {
		return nil, err
	}

	return api.MapCursor(cursor, Transaction.toCore), nil
}

func (s *Store) CountTransactions(ctx context.Context, q TransactionsQuery) (uint64, error) {
	models := make([]Transaction, 0)
	count, err := s.buildTransactionsQuery(q.Filters, &models).Count(ctx)

	return uint64(count), storageerrors.PostgresError(err)
}

func (s *Store) GetTransaction(ctx context.Context, txId uint64) (*core.ExpandedTransaction, error) {

	tx := &Transaction{}
	err := s.schema.NewSelect(TransactionsTableName).
		Model(tx).
		ColumnExpr(`transactions.id, transactions.reference, transactions.metadata, transactions.postings,
			transactions.date, get_aggregated_volumes_for_transaction(transactions) as post_commit_aggregated_volumes`).
		Where("id = ?", txId).
		Order("revision desc").
		Limit(1).
		Scan(ctx)
	if err != nil {
		return nil, storageerrors.PostgresError(err)
	}

	return pointer.For(tx.toCore()), nil
}

func (s *Store) GetTransactionByReference(ctx context.Context, ref string) (*core.ExpandedTransaction, error) {
	tx := &Transaction{}
	err := s.schema.NewSelect(TransactionsTableName).
		Model(tx).
		Where("reference = ?", ref).
		OrderExpr("id DESC").
		Scan(ctx)
	if err != nil {
		return nil, storageerrors.PostgresError(err)
	}

	return pointer.For(tx.toCore()), nil
}
