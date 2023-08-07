package ledgerstore

import (
	"context"
	"database/sql/driver"
	"encoding/json"

	"github.com/formancehq/ledger/pkg/core"
	storageerrors "github.com/formancehq/ledger/pkg/storage"
	"github.com/formancehq/ledger/pkg/storage/paginate"
	"github.com/formancehq/stack/libs/go-libs/api"
	"github.com/formancehq/stack/libs/go-libs/collectionutils"
	"github.com/lib/pq"
	"github.com/pkg/errors"
	"github.com/uptrace/bun"
)

const (
	LogTableName = "logs"
)

type Logs struct {
	bun.BaseModel `bun:"logs,alias:logs"`

	ID             uint64    `bun:"id,unique,type:bigint"`
	Type           string    `bun:"type,type:log_type"`
	Hash           []byte    `bun:"hash,type:bytea"`
	Date           core.Time `bun:"date,type:timestamptz"`
	Data           []byte    `bun:"data,type:jsonb"`
	IdempotencyKey string    `bun:"idempotency_key,type:varchar(256),unique"`
}

func (log *Logs) ToCore() *core.ChainedLog {
	payload, err := core.HydrateLog(core.LogTypeFromString(log.Type), log.Data)
	if err != nil {
		panic(errors.Wrap(err, "hydrating log data"))
	}

	return &core.ChainedLog{
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

func (s *Store) logsQueryBuilder(q LogsQueryFilters) func(*bun.SelectQuery) *bun.SelectQuery {
	return func(query *bun.SelectQuery) *bun.SelectQuery {
		query = query.Table(LogTableName)
		if !q.StartTime.IsZero() {
			query = query.Where("date >= ?", q.StartTime.UTC())
		}
		if !q.EndTime.IsZero() {
			query = query.Where("date < ?", q.EndTime.UTC())
		}
		return query
	}
}

func (s *Store) InsertLogs(ctx context.Context, activeLogs ...*core.ChainedLog) error {
	return s.withTransaction(ctx, func(tx bun.Tx) error {
		// Beware: COPY query is not supported by bun if the pgx driver is used.
		stmt, err := tx.Prepare(pq.CopyInSchema(
			s.name,
			LogTableName,
			"id", "type", "hash", "date", "data", "idempotency_key",
		))
		if err != nil {
			return storageerrors.PostgresError(err)
		}

		ls := make([]Logs, len(activeLogs))
		for i, chainedLogs := range activeLogs {
			data, err := json.Marshal(chainedLogs.Data)
			if err != nil {
				return errors.Wrap(err, "marshaling log data")
			}

			ls[i] = Logs{
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

		return stmt.Close()
	})
}

func (s *Store) GetLastLog(ctx context.Context) (*core.ChainedLog, error) {
	return fetchAndMap[*Logs, *core.ChainedLog](s, ctx, (*Logs).ToCore,
		func(query *bun.SelectQuery) *bun.SelectQuery {
			return query.OrderExpr("id desc").Limit(1)
		})
}

func (s *Store) GetLogs(ctx context.Context, q LogsQuery) (*api.Cursor[core.ChainedLog], error) {
	logs, err := paginateWithColumn[LogsQueryFilters, Logs](s, ctx,
		paginate.ColumnPaginatedQuery[LogsQueryFilters](q),
		s.logsQueryBuilder(q.Filters),
	)
	if err != nil {
		return nil, err
	}

	return api.MapCursor(logs, func(from Logs) core.ChainedLog {
		return *from.ToCore()
	}), nil
}

func (s *Store) ReadLastLogWithType(ctx context.Context, logTypes ...core.LogType) (*core.ChainedLog, error) {
	return fetchAndMap[*Logs](s, ctx, (*Logs).ToCore,
		func(query *bun.SelectQuery) *bun.SelectQuery {
			return query.
				Where("type IN (?)", bun.In(collectionutils.Map(logTypes, core.LogType.String))).
				OrderExpr("date DESC").
				Limit(1)
		})
}

func (s *Store) ReadLogWithIdempotencyKey(ctx context.Context, key string) (*core.ChainedLog, error) {
	return fetchAndMap[*Logs, *core.ChainedLog](s, ctx, (*Logs).ToCore,
		func(query *bun.SelectQuery) *bun.SelectQuery {
			return query.
				OrderExpr("id desc").
				Limit(1).
				Where("idempotency_key = ?", key)
		})
}

type LogsQueryFilters struct {
	EndTime   core.Time `json:"endTime"`
	StartTime core.Time `json:"startTime"`
}

type LogsQuery paginate.ColumnPaginatedQuery[LogsQueryFilters]

func NewLogsQuery() LogsQuery {
	return LogsQuery{
		PageSize: paginate.QueryDefaultPageSize,
		Column:   "id",
		Order:    paginate.OrderDesc,
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
