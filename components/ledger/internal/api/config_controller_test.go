package api_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/formancehq/ledger/internal/api"
	"github.com/formancehq/ledger/internal/opentelemetry/metrics"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestGetInfo(t *testing.T) {
	t.Parallel()

	backend, _ := newTestingBackend(t)
	router := api.NewRouter(backend, nil, metrics.NewNoOpRegistry())

	backend.
		EXPECT().
		ListLedgers(gomock.Any()).
		Return([]string{"a", "b"}, nil)

	backend.
		EXPECT().
		GetVersion().
		Return("latest")

	req := httptest.NewRequest(http.MethodGet, "/_info", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)

	info := api.ConfigInfo{}
	require.NoError(t, json.NewDecoder(rec.Body).Decode(&info))

	require.EqualValues(t, api.ConfigInfo{
		Server:  "ledger",
		Version: "latest",
		Config: &api.LedgerConfig{
			LedgerStorage: &api.LedgerStorage{
				Driver:  "postgres",
				Ledgers: []string{"a", "b"},
			},
		},
	}, info)
}
