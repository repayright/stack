package api_test

import (
	"math/big"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	ledger "github.com/formancehq/ledger/internal"
	"github.com/formancehq/ledger/internal/api"
	"github.com/formancehq/ledger/internal/opentelemetry/metrics"
	"github.com/formancehq/ledger/internal/storage/ledgerstore"
	sharedapi "github.com/formancehq/stack/libs/go-libs/api"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestGetBalancesAggregated(t *testing.T) {
	t.Parallel()

	type testCase struct {
		name        string
		queryParams url.Values
		expectQuery ledgerstore.GetAggregatedBalancesQuery
	}

	testCases := []testCase{
		{
			name:        "nominal",
			expectQuery: ledgerstore.NewGetAggregatedBalancesQuery(),
		},
		{
			name: "using address",
			queryParams: url.Values{
				"address": []string{"foo"},
			},
			expectQuery: ledgerstore.NewGetAggregatedBalancesQuery().WithAddressFilter("foo"),
		},
	}
	for _, testCase := range testCases {
		testCase := testCase
		t.Run(testCase.name, func(t *testing.T) {

			expectedBalances := ledger.BalancesByAssets{
				"world": big.NewInt(-100),
			}
			backend, mock := newTestingBackend(t)
			mock.EXPECT().
				GetAggregatedBalances(gomock.Any(), testCase.expectQuery).
				Return(expectedBalances, nil)

			router := api.NewRouter(backend, nil, metrics.NewNoOpRegistry())

			req := httptest.NewRequest(http.MethodGet, "/xxx/aggregate/balances", nil)
			rec := httptest.NewRecorder()
			req.URL.RawQuery = testCase.queryParams.Encode()

			router.ServeHTTP(rec, req)

			require.Equal(t, http.StatusOK, rec.Code)
			balances, ok := sharedapi.DecodeSingleResponse[ledger.BalancesByAssets](t, rec.Body)
			require.True(t, ok)
			require.Equal(t, expectedBalances, balances)
		})
	}
}
