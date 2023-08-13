package ledgerstore

import (
	"context"
	"database/sql/driver"
	"encoding/json"
	"math/big"
	"strings"

	"github.com/formancehq/ledger/pkg/core"
	"github.com/formancehq/ledger/pkg/storage/paginate"
	"github.com/formancehq/stack/libs/go-libs/api"
	"github.com/formancehq/stack/libs/go-libs/metadata"
	"github.com/uptrace/bun"
)

const (
	MovesTableName = "moves"
)

type Transaction struct {
	bun.BaseModel `bun:"transactions,alias:transactions"`

	ID                         uint64                     `bun:"id,type:bigint"`
	Timestamp                  core.Time                  `bun:"date,type:timestamp without time zone"`
	Reference                  string                     `bun:"reference,type:varchar,unique,nullzero"`
	Postings                   []core.Posting             `bun:"postings,type:jsonb"`
	Metadata                   metadata.Metadata          `bun:"metadata,type:jsonb,default:'{}'"`
	PostCommitEffectiveVolumes core.AccountsAssetsVolumes `bun:"post_commit_effective_volumes,type:jsonb"`
	PostCommitVolumes          core.AccountsAssetsVolumes `bun:"post_commit_volumes,type:jsonb"`
}

func (t *Transaction) toCore() *core.ExpandedTransaction {
	var (
		preCommitEffectiveVolumes core.AccountsAssetsVolumes
		preCommitVolumes          core.AccountsAssetsVolumes
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
	return &core.ExpandedTransaction{
		Transaction: core.Transaction{
			TransactionData: core.TransactionData{
				Reference: t.Reference,
				Metadata:  t.Metadata,
				Date:      t.Timestamp,
				Postings:  t.Postings,
			},
			ID: t.ID,
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

func (s *Store) listTransactionsBuilder(p TransactionsQueryOptions) func(query *bun.SelectQuery) *bun.SelectQuery {
	return func(query *bun.SelectQuery) *bun.SelectQuery {
		query = query.
			Table("transactions").
			ColumnExpr("distinct on(transactions.id) transactions.id").
			ColumnExpr("transactions.reference").
			ColumnExpr("transactions.metadata").
			ColumnExpr("transactions.postings").
			ColumnExpr("transactions.date").
			Apply(filterMetadata(p.Metadata))
		if p.Reference != "" {
			query.Where("transactions.reference = ?", p.Reference)
		}
		if !p.StartTime.IsZero() {
			query.Where("transactions.date >= ?", p.StartTime)
		}
		if !p.EndTime.IsZero() {
			query.Where("transactions.date < ?", p.EndTime)
		}
		if p.AfterTxID != 0 {
			query.Where("transactions.id > ?", p.AfterTxID)
		}
		if p.Source != "" || p.Destination != "" || p.Account != "" {
			query.Join("join moves m on transactions.id = m.transaction_id")
			if p.Source != "" {
				query = query.
					Where("m.is_source").
					Apply(filterAccountAddress(p.Source, "account_address"))

			}
			if p.Destination != "" {
				query = query.
					Where("not m.is_source").
					Apply(filterAccountAddress(p.Destination, "account_address"))
			}
			if p.Account != "" {
				query = query.Apply(filterAccountAddress(p.Account, "account_address"))
			}
		}
		if p.ExpandEffectiveVolumes {
			query = query.ColumnExpr("get_aggregated_effective_volumes_for_transaction(transactions) as post_commit_effective_volumes")
		}
		if p.ExpandVolumes {
			query = query.ColumnExpr("get_aggregated_volumes_for_transaction(transactions) as post_commit_volumes")
		}
		return query
	}
}

func (s *Store) GetTransactions(ctx context.Context, q GetTransactionsQuery) (*api.Cursor[core.ExpandedTransaction], error) {
	transactions, err := paginateWithColumn[TransactionsQueryOptions, Transaction](s, ctx,
		paginate.ColumnPaginatedQuery[TransactionsQueryOptions](q),
		s.listTransactionsBuilder(q.Options),
	)
	if err != nil {
		return nil, err
	}
	return api.MapCursor(transactions, func(from Transaction) core.ExpandedTransaction {
		return *from.toCore()
	}), nil
}

func (s *Store) CountTransactions(ctx context.Context, q GetTransactionsQuery) (uint64, error) {
	return count(s, ctx, s.listTransactionsBuilder(q.Options))
}

func (s *Store) GetTransactionWithVolumes(ctx context.Context, txId uint64, expandVolumes, expandEffectiveVolumes bool) (*core.ExpandedTransaction, error) {
	return fetchAndMap[*Transaction, *core.ExpandedTransaction](s, ctx,
		(*Transaction).toCore,
		func(query *bun.SelectQuery) *bun.SelectQuery {
			query = query.
				ColumnExpr(`transactions.id, transactions.reference, transactions.metadata, transactions.postings, transactions.date`).
				Where("id = ?", txId).
				Order("revision desc").
				Limit(1)
			if expandEffectiveVolumes {
				query = query.ColumnExpr(`get_aggregated_effective_volumes_for_transaction(transactions) as post_commit_effective_volumes`)
			}
			if expandVolumes {
				query = query.ColumnExpr(`get_aggregated_volumes_for_transaction(transactions) as post_commit_volumes`)
			}
			return query
		})
}

func (s *Store) GetTransaction(ctx context.Context, txId uint64) (*core.ExpandedTransaction, error) {
	return fetchAndMap[*Transaction, *core.ExpandedTransaction](s, ctx,
		(*Transaction).toCore,
		func(query *bun.SelectQuery) *bun.SelectQuery {
			return query.
				ColumnExpr(`transactions.id, transactions.reference, transactions.metadata, transactions.postings, transactions.date`).
				Where("id = ?", txId).
				Order("revision desc").
				Limit(1)
		})
}

func (s *Store) GetTransactionByReference(ctx context.Context, ref string) (*core.ExpandedTransaction, error) {
	return fetchAndMap[*Transaction, *core.ExpandedTransaction](s, ctx,
		(*Transaction).toCore,
		func(query *bun.SelectQuery) *bun.SelectQuery {
			return query.
				ColumnExpr(`transactions.id, transactions.reference, transactions.metadata, transactions.postings, transactions.date`).
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
	AfterTxID              uint64            `json:"afterTxID,omitempty"`
	Reference              string            `json:"reference,omitempty"`
	Destination            string            `json:"destination,omitempty"`
	Source                 string            `json:"source,omitempty"`
	Account                string            `json:"account,omitempty"`
	EndTime                core.Time         `json:"endTime,omitempty"`
	StartTime              core.Time         `json:"startTime,omitempty"`
	Metadata               metadata.Metadata `json:"metadata,omitempty"`
	ExpandVolumes          bool              `json:"expandVolumes"`
	ExpandEffectiveVolumes bool              `json:"expandEffectiveVolumes"`
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

func (a GetTransactionsQuery) WithStartTimeFilter(start core.Time) GetTransactionsQuery {
	if !start.IsZero() {
		a.Options.StartTime = start
	}

	return a
}

func (a GetTransactionsQuery) WithEndTimeFilter(end core.Time) GetTransactionsQuery {
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
