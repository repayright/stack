package paginate_test

import (
	"context"
	"testing"

	"github.com/formancehq/ledger/pkg/storage"
	"github.com/formancehq/ledger/pkg/storage/paginate"
	"github.com/formancehq/stack/libs/go-libs/pgtesting"
	"github.com/formancehq/stack/libs/go-libs/pointer"
	"github.com/stretchr/testify/require"
)

func TestColumnPagination(t *testing.T) {
	t.Parallel()

	pgServer := pgtesting.NewPostgresDatabase(t)
	db, err := storage.OpenSQLDB(storage.ConnectionOptions{
		DatabaseSourceName: pgServer.ConnString(),
		Debug:              testing.Verbose(),
		Trace:              testing.Verbose(),
	})
	require.NoError(t, err)

	_, err = db.Exec(`
		CREATE TABLE "models" (id int, pair boolean);
	`)
	require.NoError(t, err)

	type model struct {
		ID   uint64 `bun:"id"`
		Pair bool   `bun:"pair"`
	}

	models := make([]model, 0)
	for i := 0; i < 100; i++ {
		models = append(models, model{
			ID:   uint64(i),
			Pair: i%2 == 0,
		})
	}

	_, err = db.NewInsert().
		Model(&models).
		Exec(context.Background())
	require.NoError(t, err)

	type testCase struct {
		name                  string
		query                 paginate.ColumnPaginatedQuery[bool]
		expectedNext          *paginate.ColumnPaginatedQuery[bool]
		expectedPrevious      *paginate.ColumnPaginatedQuery[bool]
		expectedNumberOfItems uint64
	}
	testCases := []testCase{
		{
			name: "asc first page",
			query: paginate.ColumnPaginatedQuery[bool]{
				PageSize: 10,
				Column:   "id",
				Order:    paginate.OrderAsc,
			},
			expectedNext: &paginate.ColumnPaginatedQuery[bool]{
				PageSize:     10,
				Column:       "id",
				PaginationID: pointer.For(uint64(10)),
				Order:        paginate.OrderAsc,
				Bottom:       pointer.For(uint64(0)),
			},
			expectedNumberOfItems: 10,
		},
		{
			name: "asc second page using next cursor",
			query: paginate.ColumnPaginatedQuery[bool]{
				PageSize:     10,
				Column:       "id",
				PaginationID: pointer.For(uint64(10)),
				Order:        paginate.OrderAsc,
				Bottom:       pointer.For(uint64(0)),
			},
			expectedPrevious: &paginate.ColumnPaginatedQuery[bool]{
				PageSize:     10,
				Column:       "id",
				Order:        paginate.OrderAsc,
				Bottom:       pointer.For(uint64(0)),
				PaginationID: pointer.For(uint64(10)),
				Reverse:      true,
			},
			expectedNext: &paginate.ColumnPaginatedQuery[bool]{
				PageSize:     10,
				Column:       "id",
				PaginationID: pointer.For(uint64(20)),
				Order:        paginate.OrderAsc,
				Bottom:       pointer.For(uint64(0)),
			},
			expectedNumberOfItems: 10,
		},
		{
			name: "asc last page using next cursor",
			query: paginate.ColumnPaginatedQuery[bool]{
				PageSize:     10,
				Column:       "id",
				PaginationID: pointer.For(uint64(90)),
				Order:        paginate.OrderAsc,
				Bottom:       pointer.For(uint64(0)),
			},
			expectedPrevious: &paginate.ColumnPaginatedQuery[bool]{
				PageSize:     10,
				Column:       "id",
				Order:        paginate.OrderAsc,
				PaginationID: pointer.For(uint64(90)),
				Bottom:       pointer.For(uint64(0)),
				Reverse:      true,
			},
			expectedNumberOfItems: 10,
		},
		{
			name: "desc first page",
			query: paginate.ColumnPaginatedQuery[bool]{
				PageSize: 10,
				Column:   "id",
				Order:    paginate.OrderDesc,
			},
			expectedNext: &paginate.ColumnPaginatedQuery[bool]{
				PageSize:     10,
				Bottom:       pointer.For(uint64(99)),
				Column:       "id",
				PaginationID: pointer.For(uint64(89)),
				Order:        paginate.OrderDesc,
			},
			expectedNumberOfItems: 10,
		},
		{
			name: "desc second page using next cursor",
			query: paginate.ColumnPaginatedQuery[bool]{
				PageSize:     10,
				Bottom:       pointer.For(uint64(99)),
				Column:       "id",
				PaginationID: pointer.For(uint64(89)),
				Order:        paginate.OrderDesc,
			},
			expectedPrevious: &paginate.ColumnPaginatedQuery[bool]{
				PageSize:     10,
				Bottom:       pointer.For(uint64(99)),
				Column:       "id",
				PaginationID: pointer.For(uint64(89)),
				Order:        paginate.OrderDesc,
				Reverse:      true,
			},
			expectedNext: &paginate.ColumnPaginatedQuery[bool]{
				PageSize:     10,
				Bottom:       pointer.For(uint64(99)),
				Column:       "id",
				PaginationID: pointer.For(uint64(79)),
				Order:        paginate.OrderDesc,
			},
			expectedNumberOfItems: 10,
		},
		{
			name: "desc last page using next cursor",
			query: paginate.ColumnPaginatedQuery[bool]{
				PageSize:     10,
				Bottom:       pointer.For(uint64(99)),
				Column:       "id",
				PaginationID: pointer.For(uint64(9)),
				Order:        paginate.OrderDesc,
			},
			expectedPrevious: &paginate.ColumnPaginatedQuery[bool]{
				PageSize:     10,
				Bottom:       pointer.For(uint64(99)),
				Column:       "id",
				PaginationID: pointer.For(uint64(9)),
				Order:        paginate.OrderDesc,
				Reverse:      true,
			},
			expectedNumberOfItems: 10,
		},
		{
			name: "asc first page using previous cursor",
			query: paginate.ColumnPaginatedQuery[bool]{
				PageSize:     10,
				Bottom:       pointer.For(uint64(0)),
				Column:       "id",
				PaginationID: pointer.For(uint64(10)),
				Order:        paginate.OrderAsc,
				Reverse:      true,
			},
			expectedNext: &paginate.ColumnPaginatedQuery[bool]{
				PageSize:     10,
				Bottom:       pointer.For(uint64(0)),
				Column:       "id",
				PaginationID: pointer.For(uint64(10)),
				Order:        paginate.OrderAsc,
			},
			expectedNumberOfItems: 10,
		},
		{
			name: "desc first page using previous cursor",
			query: paginate.ColumnPaginatedQuery[bool]{
				PageSize:     10,
				Bottom:       pointer.For(uint64(99)),
				Column:       "id",
				PaginationID: pointer.For(uint64(89)),
				Order:        paginate.OrderDesc,
				Reverse:      true,
			},
			expectedNext: &paginate.ColumnPaginatedQuery[bool]{
				PageSize:     10,
				Bottom:       pointer.For(uint64(99)),
				Column:       "id",
				PaginationID: pointer.For(uint64(89)),
				Order:        paginate.OrderDesc,
			},
			expectedNumberOfItems: 10,
		},
		{
			name: "asc first page with filter",
			query: paginate.ColumnPaginatedQuery[bool]{
				PageSize: 10,
				Column:   "id",
				Order:    paginate.OrderAsc,
				Filters:  true,
			},
			expectedNext: &paginate.ColumnPaginatedQuery[bool]{
				PageSize:     10,
				Column:       "id",
				PaginationID: pointer.For(uint64(20)),
				Order:        paginate.OrderAsc,
				Filters:      true,
				Bottom:       pointer.For(uint64(0)),
			},
			expectedNumberOfItems: 10,
		},
		{
			name: "asc second page with filter",
			query: paginate.ColumnPaginatedQuery[bool]{
				PageSize:     10,
				Column:       "id",
				PaginationID: pointer.For(uint64(20)),
				Order:        paginate.OrderAsc,
				Filters:      true,
				Bottom:       pointer.For(uint64(0)),
			},
			expectedNext: &paginate.ColumnPaginatedQuery[bool]{
				PageSize:     10,
				Column:       "id",
				PaginationID: pointer.For(uint64(40)),
				Order:        paginate.OrderAsc,
				Filters:      true,
				Bottom:       pointer.For(uint64(0)),
			},
			expectedPrevious: &paginate.ColumnPaginatedQuery[bool]{
				PageSize:     10,
				Column:       "id",
				PaginationID: pointer.For(uint64(20)),
				Order:        paginate.OrderAsc,
				Filters:      true,
				Bottom:       pointer.For(uint64(0)),
				Reverse:      true,
			},
			expectedNumberOfItems: 10,
		},
		{
			name: "desc first page with filter",
			query: paginate.ColumnPaginatedQuery[bool]{
				PageSize: 10,
				Column:   "id",
				Order:    paginate.OrderDesc,
				Filters:  true,
			},
			expectedNext: &paginate.ColumnPaginatedQuery[bool]{
				PageSize:     10,
				Column:       "id",
				PaginationID: pointer.For(uint64(78)),
				Order:        paginate.OrderDesc,
				Filters:      true,
				Bottom:       pointer.For(uint64(98)),
			},
			expectedNumberOfItems: 10,
		},
		{
			name: "desc second page with filter",
			query: paginate.ColumnPaginatedQuery[bool]{
				PageSize:     10,
				Column:       "id",
				PaginationID: pointer.For(uint64(78)),
				Order:        paginate.OrderDesc,
				Filters:      true,
				Bottom:       pointer.For(uint64(98)),
			},
			expectedNext: &paginate.ColumnPaginatedQuery[bool]{
				PageSize:     10,
				Column:       "id",
				PaginationID: pointer.For(uint64(58)),
				Order:        paginate.OrderDesc,
				Filters:      true,
				Bottom:       pointer.For(uint64(98)),
			},
			expectedPrevious: &paginate.ColumnPaginatedQuery[bool]{
				PageSize:     10,
				Column:       "id",
				PaginationID: pointer.For(uint64(78)),
				Order:        paginate.OrderDesc,
				Filters:      true,
				Bottom:       pointer.For(uint64(98)),
				Reverse:      true,
			},
			expectedNumberOfItems: 10,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			models := make([]model, 0)
			query := db.NewSelect().Model(&models).Column("id")
			if tc.query.Filters {
				query = query.Where("pair = ?", true)
			}
			cursor, err := paginate.UsingColumn[bool, model](context.Background(), query, tc.query)
			require.NoError(t, err)

			if tc.expectedNext == nil {
				require.Empty(t, cursor.Next)
			} else {
				require.NotEmpty(t, cursor.Next)

				q := paginate.ColumnPaginatedQuery[bool]{}
				require.NoError(t, paginate.UnmarshalCursor(cursor.Next, &q))
				require.EqualValues(t, *tc.expectedNext, q)
			}

			if tc.expectedPrevious == nil {
				require.Empty(t, cursor.Previous)
			} else {
				require.NotEmpty(t, cursor.Previous)

				q := paginate.ColumnPaginatedQuery[bool]{}
				require.NoError(t, paginate.UnmarshalCursor(cursor.Previous, &q))
				require.EqualValues(t, *tc.expectedPrevious, q)
			}
		})
	}
}
