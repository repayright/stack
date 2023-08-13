package api

import (
	"net/http"

	"github.com/formancehq/ledger/internal/storage/ledgerstore"
	sharedapi "github.com/formancehq/stack/libs/go-libs/api"
)

func getBalancesAggregated(w http.ResponseWriter, r *http.Request) {
	l := LedgerFromContext(r.Context())

	balancesQuery := ledgerstore.NewGetAggregatedBalancesQuery().WithAddressFilter(r.URL.Query().Get("address"))
	balances, err := l.GetAggregatedBalances(r.Context(), balancesQuery)
	if err != nil {
		ResponseError(w, r, err)
		return
	}

	sharedapi.Ok(w, balances)
}
