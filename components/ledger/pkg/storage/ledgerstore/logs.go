package ledgerstore

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"math/big"

	"github.com/formancehq/ledger/pkg/core"
	storageerrors "github.com/formancehq/ledger/pkg/storage"
	"github.com/formancehq/stack/libs/go-libs/api"
	"github.com/formancehq/stack/libs/go-libs/collectionutils"
	"github.com/formancehq/stack/libs/go-libs/metadata"
	"github.com/lib/pq"
	"github.com/pkg/errors"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/extra/bunbig"
)

const (
	LogTableName = "logs"
)

type LogsQueryFilters struct {
	EndTime   core.Time `json:"endTime"`
	StartTime core.Time `json:"startTime"`
}

type LogsQuery ColumnPaginatedQuery[LogsQueryFilters]

func NewLogsQuery() LogsQuery {
	return LogsQuery{
		PageSize: QueryDefaultPageSize,
		Column:   "id",
		Order:    OrderDesc,
		Filters:  LogsQueryFilters{},
	}
}

func (a LogsQuery) WithPaginationID(id uint64) LogsQuery {
	a.PaginationID = &id
	return a
}

func (l LogsQuery) WithPageSize(pageSize uint64) LogsQuery {
	if pageSize != 0 {
		l.PageSize = pageSize
	}

	return l
}

func (l LogsQuery) WithStartTimeFilter(start core.Time) LogsQuery {
	if !start.IsZero() {
		l.Filters.StartTime = start
	}

	return l
}

func (l LogsQuery) WithEndTimeFilter(end core.Time) LogsQuery {
	if !end.IsZero() {
		l.Filters.EndTime = end
	}

	return l
}

type AccountWithBalances struct {
	bun.BaseModel `bun:"accounts,alias:accounts"`

	Address  string              `bun:"address,type:varchar,unique,notnull"`
	Metadata metadata.Metadata   `bun:"metadata,type:bytea,default:'{}'"`
	Balances map[string]*big.Int `bun:"balances,type:bytea,default:'{}'"`
}

type LogsV2 struct {
	bun.BaseModel `bun:"logs,alias:logs"`

	ID             uint64    `bun:"id,unique,type:bigint"`
	Type           string    `bun:"type,type:log_type"`
	Hash           []byte    `bun:"hash,type:bytea"`
	Date           core.Time `bun:"date,type:timestamptz"`
	Data           []byte    `bun:"data,type:jsonb"`
	IdempotencyKey string    `bun:"idempotency_key,type:varchar(256),unique"`
}

func (log LogsV2) ToCore() core.ChainedLog {
	payload, err := core.HydrateLog(core.LogTypeFromString(log.Type), log.Data)
	if err != nil {
		panic(errors.Wrap(err, "hydrating log data"))
	}

	return core.ChainedLog{
		Log: core.Log{
			Type:           core.LogTypeFromString(log.Type),
			Data:           payload,
			Date:           log.Date.UTC(),
			IdempotencyKey: log.IdempotencyKey,
		},
		ID:   log.ID,
		Hash: log.Hash,
	}
}

type RawMessage json.RawMessage

func (j RawMessage) Value() (driver.Value, error) {
	if j == nil {
		return nil, nil
	}
	return string(j), nil
}

func (s *Store) InsertLogs(ctx context.Context, activeLogs ...*core.ChainedLog) error {

	txn, err := s.schema.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return storageerrors.PostgresError(err)
	}

	// Beware: COPY query is not supported by bun if the pgx driver is used.
	stmt, err := txn.Prepare(pq.CopyInSchema(
		s.schema.Name(),
		LogTableName,
		"id", "type", "hash", "date", "data", "idempotency_key",
	))
	if err != nil {
		return storageerrors.PostgresError(err)
	}

	ls := make([]LogsV2, len(activeLogs))
	for i, chainedLogs := range activeLogs {
		data, err := json.Marshal(chainedLogs.Data)
		if err != nil {
			return errors.Wrap(err, "marshaling log data")
		}

		ls[i] = LogsV2{
			ID:             chainedLogs.ID,
			Type:           chainedLogs.Type.String(),
			Hash:           chainedLogs.Hash,
			Date:           chainedLogs.Date,
			Data:           data,
			IdempotencyKey: chainedLogs.IdempotencyKey,
		}

		_, err = stmt.Exec(ls[i].ID, ls[i].Type, ls[i].Hash, ls[i].Date, RawMessage(ls[i].Data), chainedLogs.IdempotencyKey)
		if err != nil {
			return storageerrors.PostgresError(err)
		}
	}

	_, err = stmt.Exec()
	if err != nil {
		return storageerrors.PostgresError(err)
	}

	err = stmt.Close()
	if err != nil {
		return storageerrors.PostgresError(err)
	}

	return storageerrors.PostgresError(txn.Commit())
}

func (s *Store) GetLastLog(ctx context.Context) (*core.ChainedLog, error) {
	raw := &LogsV2{}
	err := s.schema.NewSelect(LogTableName).
		Model(raw).
		OrderExpr("id desc").
		Limit(1).
		Scan(ctx)
	if err != nil {
		return nil, storageerrors.PostgresError(err)
	}

	l := raw.ToCore()
	return &l, nil
}

func (s *Store) GetLogs(ctx context.Context, q LogsQuery) (*api.Cursor[core.ChainedLog], error) {
	cursor, err := UsingColumn[LogsQueryFilters, LogsV2](ctx,
		s.buildLogsQuery,
		ColumnPaginatedQuery[LogsQueryFilters](q),
	)
	if err != nil {
		return nil, err
	}

	return api.MapCursor(cursor, LogsV2.ToCore), nil
}

func (s *Store) buildLogsQuery(q LogsQueryFilters, models *[]LogsV2) *bun.SelectQuery {
	sb := s.schema.NewSelect(LogTableName).
		Model(models).
		Column("id", "type", "hash", "date", "data", "idempotency_key")

	if !q.StartTime.IsZero() {
		sb.Where("date >= ?", q.StartTime.UTC())
	}
	if !q.EndTime.IsZero() {
		sb.Where("date < ?", q.EndTime.UTC())
	}

	return sb
}

func (s *Store) GetNextLogID(ctx context.Context) (uint64, error) {
	var logID uint64
	err := s.schema.
		NewSelect(LogTableName).
		ColumnExpr("min(id)").
		Where("projected = FALSE").
		Limit(1).
		Scan(ctx, &logID)
	if err != nil {
		return 0, storageerrors.PostgresError(err)
	}
	return logID, nil
}

func (s *Store) ReadLogsRange(ctx context.Context, idMin, idMax uint64) ([]core.ChainedLog, error) {
	rawLogs := make([]LogsV2, 0)
	err := s.schema.
		NewSelect(LogTableName).
		Where("id >= ?", idMin).
		Where("id < ?", idMax).
		Model(&rawLogs).
		Scan(ctx)
	if err != nil {
		return nil, storageerrors.PostgresError(err)
	}

	return collectionutils.Map(rawLogs, LogsV2.ToCore), nil
}

func (s *Store) ReadLastLogWithType(ctx context.Context, logTypes ...core.LogType) (*core.ChainedLog, error) {
	raw := &LogsV2{}
	err := s.schema.
		NewSelect(LogTableName).
		Where("type IN (?)", bun.In(collectionutils.Map(logTypes, core.LogType.String))).
		OrderExpr("date DESC").
		Model(raw).
		Limit(1).
		Scan(ctx)
	if err != nil {
		return nil, storageerrors.PostgresError(err)
	}
	ret := raw.ToCore()

	return &ret, nil
}

func (s *Store) ReadLogWithIdempotencyKey(ctx context.Context, key string) (*core.ChainedLog, error) {
	raw := &LogsV2{}
	err := s.schema.NewSelect(LogTableName).
		Model(raw).
		OrderExpr("id desc").
		Limit(1).
		Where("idempotency_key = ?", key).
		Scan(ctx)
	if err != nil {
		return nil, storageerrors.PostgresError(err)
	}

	l := raw.ToCore()
	return &l, nil
}

func (s *Store) GetBalance(ctx context.Context, address, asset string) (*big.Int, error) {
	selectLogsForExistingAccount := s.schema.
		NewSelect(LogTableName).
		Model(&LogsV2{}).
		Where(fmt.Sprintf(`data->'transaction'->'postings' @> '[{"destination": "%s", "asset": "%s"}]' OR data->'transaction'->'postings' @> '[{"source": "%s", "asset": "%s"}]'`, address, asset, address, asset))

	selectPostings := s.schema.IDB.NewSelect().
		TableExpr(`(` + selectLogsForExistingAccount.String() + `) as logs`).
		ColumnExpr("jsonb_array_elements(logs.data->'transaction'->'postings') as postings")

	selectBalances := s.schema.IDB.NewSelect().
		TableExpr(`(` + selectPostings.String() + `) as postings`).
		ColumnExpr(fmt.Sprintf("SUM(CASE WHEN (postings.postings::jsonb)->>'source' = '%s' THEN -((((postings.postings::jsonb)->'amount')::numeric)) ELSE ((postings.postings::jsonb)->'amount')::numeric END)", address))

	row := s.schema.IDB.QueryRowContext(ctx, selectBalances.String())
	if row.Err() != nil {
		return nil, row.Err()
	}

	var balance *bunbig.Int
	if err := row.Scan(&balance); err != nil {
		return nil, err
	}
	return (*big.Int)(balance), nil
}
