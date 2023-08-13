package ledgerstore

import (
	"context"
	"database/sql/driver"
	"encoding/json"
	"math/big"

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

	ID             *paginate.BigInt `bun:"id,unique,type:numeric"`
	Type           string           `bun:"type,type:log_type"`
	Hash           []byte           `bun:"hash,type:bytea"`
	Date           core.Time        `bun:"date,type:timestamptz"`
	Data           []byte           `bun:"data,type:jsonb"`
	IdempotencyKey string           `bun:"idempotency_key,type:varchar(256),unique"`
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
		ID:   (*big.Int)(log.ID),
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

func (store *Store) logsQueryBuilder(q LogsQueryOptions) func(*bun.SelectQuery) *bun.SelectQuery {
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

func (store *Store) InsertLogs(ctx context.Context, activeLogs ...*core.ChainedLog) error {
	return store.withTransaction(ctx, func(tx bun.Tx) error {
		// Beware: COPY query is not supported by bun if the pgx driver is used.
		stmt, err := tx.Prepare(pq.CopyInSchema(
			store.name,
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
				ID:             (*paginate.BigInt)(chainedLogs.ID),
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

func (store *Store) GetLastLog(ctx context.Context) (*core.ChainedLog, error) {
	return fetchAndMap[*Logs, *core.ChainedLog](store, ctx, (*Logs).ToCore,
		func(query *bun.SelectQuery) *bun.SelectQuery {
			return query.OrderExpr("id desc").Limit(1)
		})
}

func (store *Store) GetLogs(ctx context.Context, q GetLogsQuery) (*api.Cursor[core.ChainedLog], error) {
	logs, err := paginateWithColumn[LogsQueryOptions, Logs](store, ctx,
		paginate.ColumnPaginatedQuery[LogsQueryOptions](q),
		store.logsQueryBuilder(q.Options),
	)
	if err != nil {
		return nil, err
	}

	return api.MapCursor(logs, func(from Logs) core.ChainedLog {
		return *from.ToCore()
	}), nil
}

func (store *Store) ReadLastLogWithType(ctx context.Context, logTypes ...core.LogType) (*core.ChainedLog, error) {
	return fetchAndMap[*Logs](store, ctx, (*Logs).ToCore,
		func(query *bun.SelectQuery) *bun.SelectQuery {
			return query.
				Where("type IN (?)", bun.In(collectionutils.Map(logTypes, core.LogType.String))).
				OrderExpr("date DESC").
				Limit(1)
		})
}

func (store *Store) ReadLogWithIdempotencyKey(ctx context.Context, key string) (*core.ChainedLog, error) {
	return fetchAndMap[*Logs, *core.ChainedLog](store, ctx, (*Logs).ToCore,
		func(query *bun.SelectQuery) *bun.SelectQuery {
			return query.
				OrderExpr("id desc").
				Limit(1).
				Where("idempotency_key = ?", key)
		})
}

type LogsQueryOptions struct {
	EndTime   core.Time `json:"endTime"`
	StartTime core.Time `json:"startTime"`
}

type GetLogsQuery paginate.ColumnPaginatedQuery[LogsQueryOptions]

func NewLogsQuery() GetLogsQuery {
	return GetLogsQuery{
		PageSize: paginate.QueryDefaultPageSize,
		Column:   "id",
		Order:    paginate.OrderDesc,
		Options:  LogsQueryOptions{},
	}
}

func (a GetLogsQuery) WithPaginationID(id uint64) GetLogsQuery {
	a.PaginationID = big.NewInt(int64(id))
	return a
}

func (l GetLogsQuery) WithPageSize(pageSize uint64) GetLogsQuery {
	if pageSize != 0 {
		l.PageSize = pageSize
	}

	return l
}

func (l GetLogsQuery) WithStartTimeFilter(start core.Time) GetLogsQuery {
	if !start.IsZero() {
		l.Options.StartTime = start
	}

	return l
}

func (l GetLogsQuery) WithEndTimeFilter(end core.Time) GetLogsQuery {
	if !end.IsZero() {
		l.Options.EndTime = end
	}

	return l
}
