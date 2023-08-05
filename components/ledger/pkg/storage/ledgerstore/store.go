package ledgerstore

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"
	"strings"
	"sync"

	"github.com/formancehq/ledger/pkg/storage"
	"github.com/formancehq/ledger/pkg/storage/paginate"
	"github.com/formancehq/stack/libs/go-libs/api"
	"github.com/formancehq/stack/libs/go-libs/metadata"
	"github.com/formancehq/stack/libs/go-libs/migrations"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pkg/errors"
	"github.com/uptrace/bun"
)

type Store struct {
	schema   storage.Schema
	onDelete func(ctx context.Context) error

	once sync.Once

	isInitialized bool
}

func (s *Store) Schema() storage.Schema {
	return s.schema
}

func (s *Store) Name() string {
	return s.schema.Name()
}

func (s *Store) Delete(ctx context.Context) error {
	if err := s.schema.Delete(ctx); err != nil {
		return err
	}
	return errors.Wrap(s.onDelete(ctx), "deleting ledger store")
}

func (s *Store) getMigrator() *migrations.Migrator {
	migrator := migrations.NewMigrator(migrations.WithSchema(s.Name(), true))
	registerMigrations(migrator, s.schema)
	return migrator
}

func (s *Store) Migrate(ctx context.Context) (bool, error) {
	migrator := s.getMigrator()

	if err := migrator.Up(ctx, s.schema.IDB); err != nil {
		return false, err
	}

	// TODO: Update migrations package to return modifications
	return false, nil
}

func (s *Store) GetMigrationsInfo(ctx context.Context) ([]migrations.Info, error) {
	return s.getMigrator().GetMigrations(ctx, s.schema.IDB)
}

func (s *Store) IsInitialized() bool {
	return s.isInitialized
}

func (s *Store) prepareTransaction(ctx context.Context) (*storage.Tx, error) {
	tx, err := s.schema.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return tx, err
	}
	if _, err := tx.Exec(fmt.Sprintf(`set search_path = "%s"`, s.Name())); err != nil {
		return tx, err
	}
	return tx, nil
}

func (s *Store) withTransaction(ctx context.Context, callback func(tx *storage.Tx) error) error {
	tx, err := s.prepareTransaction(ctx)
	if err != nil {
		return err
	}
	if err := callback(tx); err != nil {
		_ = tx.Rollback()
		return storage.PostgresError(err)
	}
	return tx.Commit()
}

func New(
	schema storage.Schema,
	onDelete func(ctx context.Context) error,
) (*Store, error) {
	return &Store{
		schema:   schema,
		onDelete: onDelete,
	}, nil
}

func fetch[T any](s *Store, ctx context.Context, builders ...func(query *bun.SelectQuery) *bun.SelectQuery) (T, error) {
	var ret T
	ret = reflect.New(reflect.TypeOf(ret).Elem()).Interface().(T)
	err := s.withTransaction(ctx, func(tx *storage.Tx) error {
		query := s.schema.IDB.NewSelect().Conn(tx)
		for _, builder := range builders {
			query = query.Apply(builder)
		}
		if query.GetTableName() == "" && query.GetModel() == nil {
			query = query.Model(ret)
		}

		return storage.PostgresError(query.Scan(ctx, ret))
	})
	return ret, err
}

func fetchAndMap[T any, TO any](s *Store, ctx context.Context,
	mapper func(T) (TO),
	builders ...func(query *bun.SelectQuery) *bun.SelectQuery) (TO, error) {
	ret, err := fetch[T](s, ctx, builders...)
	if err != nil {
		var zero TO
		return zero, err
	}
	return mapper(ret), nil
}

func paginateWithOffset[FILTERS any, RETURN any](s *Store, ctx context.Context,
	q paginate.OffsetPaginatedQuery[FILTERS], builders ...func(query *bun.SelectQuery) *bun.SelectQuery) (*api.Cursor[RETURN], error) {
	tx, err := s.prepareTransaction(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	var ret RETURN
	query := s.schema.IDB.NewSelect().Conn(tx)
	for _, builder := range builders {
		query = query.Apply(builder)
	}
	if query.GetModel() == nil && query.GetTableName() == "" {
		query = query.Model(ret)
	}

	return paginate.UsingOffset[FILTERS, RETURN](ctx, query, q)
}

func paginateWithColumn[FILTERS any, RETURN any](s *Store, ctx context.Context, q paginate.ColumnPaginatedQuery[FILTERS], builders ...func(query *bun.SelectQuery) *bun.SelectQuery) (*api.Cursor[RETURN], error) {
	tx, err := s.prepareTransaction(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	var ret RETURN
	query := s.schema.IDB.NewSelect().Conn(tx)
	for _, builder := range builders {
		query = query.Apply(builder)
	}
	if query.GetModel() == nil && query.GetTableName() == "" {
		query = query.Model(ret)
	}

	return paginate.UsingColumn[FILTERS, RETURN](ctx, query, q)
}

func count(s *Store, ctx context.Context, builders ...func(query *bun.SelectQuery) *bun.SelectQuery) (uint64, error) {
	var (
		count int
		err   error
	)
	if err := s.withTransaction(ctx, func(tx *storage.Tx) error {
		query := s.schema.IDB.NewSelect().Conn(tx)
		for _, builder := range builders {
			query = query.Apply(builder)
		}
		count, err = query.Conn(tx).Count(ctx)
		return err
	}); err != nil {
		return 0, err
	}
	return uint64(count), nil
}

func filterMetadata(metadata metadata.Metadata) func(query *bun.SelectQuery) *bun.SelectQuery {
	return func(query *bun.SelectQuery) *bun.SelectQuery {
		if len(metadata) > 0 {
			return query.Where("metadata @> ?", metadata)
		}
		return query
	}
}

func filterAccountAddress(address, key string) func(query *bun.SelectQuery) *bun.SelectQuery {
	return func(query *bun.SelectQuery) *bun.SelectQuery {
		// todo: add check if we really need to filter on segments
		if address != "" {
			src := strings.Split(address, ":")
			query.Where(fmt.Sprintf("jsonb_array_length(%s_array) = %d", key, len(src)))

			for i, segment := range src {
				if len(segment) == 0 {
					continue
				}
				query.Where(fmt.Sprintf("%s_array @@ ('$[%d] == \"' || ?::text || '\"')::jsonpath", key, i), segment)
			}
		}
		return query
	}
}
