package ledgerstore

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	"github.com/formancehq/ledger/pkg/storage"
	"github.com/formancehq/ledger/pkg/storage/paginate"
	"github.com/formancehq/stack/libs/go-libs/api"
	"github.com/formancehq/stack/libs/go-libs/metadata"
	"github.com/uptrace/bun"
)

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
