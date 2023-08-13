package controllers

import (
	"net/http"

	"github.com/formancehq/ledger/pkg/api/apierrors"
	"github.com/formancehq/ledger/pkg/storage/ledgerstore"
	sharedapi "github.com/formancehq/stack/libs/go-libs/api"
)

func GetBalancesAggregated(w http.ResponseWriter, r *http.Request) {
	l := LedgerFromContext(r.Context())

	balancesQuery := ledgerstore.NewBalancesQuery().WithAddressFilter(r.URL.Query().Get("address"))
	balances, err := l.GetAggregatedBalances(r.Context(), balancesQuery)
	if err != nil {
		apierrors.ResponseError(w, r, err)
		return
	}

	sharedapi.Ok(w, balances)
}
