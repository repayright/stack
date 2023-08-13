package api

import (
	"net/http"

	ledger "github.com/formancehq/ledger/internal"
	"github.com/formancehq/ledger/internal/engine/command"
	"github.com/formancehq/ledger/internal/storage/ledgerstore"
	"github.com/formancehq/ledger/internal/storage/paginate"
	sharedapi "github.com/formancehq/stack/libs/go-libs/api"
	"github.com/formancehq/stack/libs/go-libs/errorsutil"
	"github.com/formancehq/stack/libs/go-libs/migrations"
	"github.com/go-chi/chi/v5"
	"github.com/pkg/errors"
)

type Info struct {
	Name    string      `json:"name"`
	Storage StorageInfo `json:"storage"`
}

type StorageInfo struct {
	Migrations []migrations.Info `json:"migrations"`
}

func getLedgerInfo(w http.ResponseWriter, r *http.Request) {
	ledger := LedgerFromContext(r.Context())

	var err error
	res := Info{
		Name:    chi.URLParam(r, "ledger"),
		Storage: StorageInfo{},
	}
	res.Storage.Migrations, err = ledger.GetMigrationsInfo(r.Context())
	if err != nil {
		ResponseError(w, r, err)
		return
	}

	sharedapi.Ok(w, res)
}

func getStats(w http.ResponseWriter, r *http.Request) {
	l := LedgerFromContext(r.Context())

	stats, err := l.Stats(r.Context())
	if err != nil {
		ResponseError(w, r, err)
		return
	}

	sharedapi.Ok(w, stats)
}

func getLogs(w http.ResponseWriter, r *http.Request) {
	l := LedgerFromContext(r.Context())

	logsQuery := ledgerstore.NewLogsQuery()

	if r.URL.Query().Get(QueryKeyCursor) != "" {
		if r.URL.Query().Get(QueryKeyStartTime) != "" ||
			r.URL.Query().Get(QueryKeyEndTime) != "" ||
			r.URL.Query().Get(QueryKeyPageSize) != "" {
			ResponseError(w, r, errorsutil.NewError(command.ErrValidation,
				errors.Errorf("no other query params can be set with '%s'", QueryKeyCursor)))
			return
		}

		err := paginate.UnmarshalCursor(r.URL.Query().Get(QueryKeyCursor), &logsQuery)
		if err != nil {
			ResponseError(w, r, errorsutil.NewError(command.ErrValidation,
				errors.Errorf("invalid '%s' query param", QueryKeyCursor)))
			return
		}
	} else {
		var err error

		var startTimeParsed, endTimeParsed ledger.Time
		if r.URL.Query().Get(QueryKeyStartTime) != "" {
			startTimeParsed, err = ledger.ParseTime(r.URL.Query().Get(QueryKeyStartTime))
			if err != nil {
				ResponseError(w, r, errorsutil.NewError(command.ErrValidation, ErrInvalidStartTime))
				return
			}
		}

		if r.URL.Query().Get(QueryKeyEndTime) != "" {
			endTimeParsed, err = ledger.ParseTime(r.URL.Query().Get(QueryKeyEndTime))
			if err != nil {
				ResponseError(w, r, errorsutil.NewError(command.ErrValidation, ErrInvalidEndTime))
				return
			}
		}

		pageSize, err := getPageSize(r)
		if err != nil {
			ResponseError(w, r, err)
			return
		}

		logsQuery = logsQuery.
			WithStartTimeFilter(startTimeParsed).
			WithEndTimeFilter(endTimeParsed).
			WithPageSize(pageSize)
	}

	cursor, err := l.GetLogs(r.Context(), logsQuery)
	if err != nil {
		ResponseError(w, r, err)
		return
	}

	sharedapi.RenderCursor(w, *cursor)
}
