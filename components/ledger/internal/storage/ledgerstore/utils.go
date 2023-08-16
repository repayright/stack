package ledgerstore

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	ledger "github.com/formancehq/ledger/internal"
	"github.com/formancehq/ledger/internal/storage"
	paginate2 "github.com/formancehq/ledger/internal/storage/paginate"
	"github.com/formancehq/stack/libs/go-libs/api"
	"github.com/formancehq/stack/libs/go-libs/metadata"
	"github.com/uptrace/bun"
)

func fetch[T any](s *Store, ctx context.Context, builders ...func(query *bun.SelectQuery) *bun.SelectQuery) (T, error) {
	var ret T
	ret = reflect.New(reflect.TypeOf(ret).Elem()).Interface().(T)
	err := s.withTransaction(ctx, func(tx bun.Tx) error {
		query := s.db.NewSelect().Conn(tx)
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
	mapper func(T) TO,
	builders ...func(query *bun.SelectQuery) *bun.SelectQuery) (TO, error) {
	ret, err := fetch[T](s, ctx, builders...)
	if err != nil {
		var zero TO
		return zero, err
	}
	return mapper(ret), nil
}

func paginateWithOffset[FILTERS any, RETURN any](s *Store, ctx context.Context,
	q paginate2.OffsetPaginatedQuery[FILTERS], builders ...func(query *bun.SelectQuery) *bun.SelectQuery) (*api.Cursor[RETURN], error) {
	tx, err := s.prepareTransaction(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	var ret RETURN
	query := s.db.NewSelect().Conn(tx)
	for _, builder := range builders {
		query = query.Apply(builder)
	}
	if query.GetModel() == nil && query.GetTableName() == "" {
		query = query.Model(ret)
	}

	return paginate2.UsingOffset[FILTERS, RETURN](ctx, query, q)
}

func paginateWithColumn[FILTERS any, RETURN any](s *Store, ctx context.Context, q paginate2.ColumnPaginatedQuery[FILTERS], builders ...func(query *bun.SelectQuery) *bun.SelectQuery) (*api.Cursor[RETURN], error) {
	tx, err := s.prepareTransaction(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	var ret RETURN
	query := s.db.NewSelect().Conn(tx)
	for _, builder := range builders {
		query = query.Apply(builder)
	}
	if query.GetModel() == nil && query.GetTableName() == "" {
		query = query.Model(ret)
	}

	return paginate2.UsingColumn[FILTERS, RETURN](ctx, query, q)
}

func count(s *Store, ctx context.Context, builders ...func(query *bun.SelectQuery) *bun.SelectQuery) (uint64, error) {
	var (
		count int
		err   error
	)
	if err := s.withTransaction(ctx, func(tx bun.Tx) error {
		query := s.db.NewSelect().Conn(tx)
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

func filterAccountAddress(address, key string) string {
	parts := make([]string, 0)
	src := strings.Split(address, ":")

	needSegmentCheck := false
	for _, segment := range src {
		needSegmentCheck = segment == ""
		if needSegmentCheck {
			break
		}
	}

	if needSegmentCheck {
		parts = append(parts, fmt.Sprintf("jsonb_array_length(%s_array) = %d", key, len(src)))

		for i, segment := range src {
			if len(segment) == 0 {
				continue
			}
			parts = append(parts, fmt.Sprintf("%s_array @@ ('$[%d] == \"%s\"')::jsonpath", key, i, segment))
		}
	} else {
		parts = append(parts, fmt.Sprintf("%s = '%s'", key, address))
	}

	return strings.Join(parts, " and ")
}

func filterAccountAddressBuilder(address, key string) func(query *bun.SelectQuery) *bun.SelectQuery {
	return func(query *bun.SelectQuery) *bun.SelectQuery {
		if address != "" {
			return query.Where(filterAccountAddress(address, key))
		}
		return query
	}
}

func filterPIT(pit *ledger.Time, column string) func(query *bun.SelectQuery) *bun.SelectQuery {
	return func(query *bun.SelectQuery) *bun.SelectQuery {
		if pit == nil || pit.IsZero() {
			return query
		}
		return query.Where(fmt.Sprintf("%s <= ?", column), pit)
	}
}
