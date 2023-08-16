package api

import (
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
	"regexp"

	ledger "github.com/formancehq/ledger/internal"
	"github.com/formancehq/ledger/internal/engine/command"
	"github.com/formancehq/ledger/internal/storage/ledgerstore"
	"github.com/formancehq/ledger/internal/storage/paginate"
	sharedapi "github.com/formancehq/stack/libs/go-libs/api"
	"github.com/formancehq/stack/libs/go-libs/collectionutils"
	"github.com/formancehq/stack/libs/go-libs/errorsutil"
	"github.com/formancehq/stack/libs/go-libs/metadata"
	"github.com/go-chi/chi/v5"
	"github.com/pkg/errors"
)

var (
	balanceRegexp = regexp.MustCompile("(" + ledger.AssetPattern + ")(<=|<|=|>|>=)([0-9]+)")
)

func countAccounts(w http.ResponseWriter, r *http.Request) {
	l := LedgerFromContext(r.Context())

	accountsQuery := ledgerstore.NewGetAccountsQuery().
		WithAddress(r.URL.Query()["address"]...).
		WithMetadataFilter(sharedapi.GetQueryMap(r.URL.Query(), "metadata"))

	count, err := l.CountAccounts(r.Context(), accountsQuery)
	if err != nil {
		ResponseError(w, r, err)
		return
	}

	w.Header().Set("Count", fmt.Sprint(count))
	sharedapi.NoContent(w)
}

func getAccounts(w http.ResponseWriter, r *http.Request) {
	l := LedgerFromContext(r.Context())

	accountsQuery := ledgerstore.NewGetAccountsQuery()

	if r.URL.Query().Get(QueryKeyCursor) != "" {
		err := paginate.UnmarshalCursor(r.URL.Query().Get(QueryKeyCursor), &accountsQuery)
		if err != nil {
			ResponseError(w, r, errorsutil.NewError(command.ErrValidation,
				errors.Errorf("invalid '%s' query param", QueryKeyCursor)))
			return
		}
	} else {
		balances := r.URL.Query()["balance"]
		balanceFilter := make(map[string]map[string]*big.Int)

		for _, balance := range balances {
			if !balanceRegexp.Match([]byte(balance)) {
				panic("balance regexp does not match")
			}

			matches := balanceRegexp.FindAllStringSubmatch(balance, 3)

			asset := matches[0][1]
			operator := matches[0][3]
			value, ok := big.NewInt(0).SetString(matches[0][4], 10)
			if !ok {
				panic("should not happen")
			}

			if _, ok := balanceFilter[asset]; !ok {
				balanceFilter[asset] = map[string]*big.Int{}
			}
			balanceFilter[asset][operator] = value
		}

		pageSize, err := getPageSize(r)
		if err != nil {
			ResponseError(w, r, err)
			return
		}

		accountsQuery = accountsQuery.
			WithAfterAddress(r.URL.Query().Get("after")).
			WithAddress(r.URL.Query()["address"]...).
			WithMetadataFilter(sharedapi.GetQueryMap(r.URL.Query(), "metadata")).
			WithBalances(balanceFilter).
			WithPageSize(pageSize)
	}

	cursor, err := l.GetAccountsWithVolumes(r.Context(), accountsQuery)
	if err != nil {
		ResponseError(w, r, err)
		return
	}

	sharedapi.RenderCursor(w, *cursor)
}

func getAccount(w http.ResponseWriter, r *http.Request) {
	l := LedgerFromContext(r.Context())

	query := ledgerstore.NewGetAccountQuery(chi.URLParam(r, "address"))
	if collectionutils.Contains(r.URL.Query()["expand"], "volumes") {
		query = query.WithExpandVolumes()
	}
	if collectionutils.Contains(r.URL.Query()["expand"], "effectiveVolumes") {
		query = query.WithExpandEffectiveVolumes()
	}

	acc, err := l.GetAccountWithVolumes(r.Context(), query)
	if err != nil {
		ResponseError(w, r, err)
		return
	}

	sharedapi.Ok(w, acc)
}

func postAccountMetadata(w http.ResponseWriter, r *http.Request) {
	l := LedgerFromContext(r.Context())

	if !ledger.ValidateAddress(chi.URLParam(r, "address")) {
		ResponseError(w, r, errorsutil.NewError(command.ErrValidation,
			errors.New("invalid account address format")))
		return
	}

	var m metadata.Metadata
	if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
		ResponseError(w, r, errorsutil.NewError(command.ErrValidation,
			errors.New("invalid metadata format")))
		return
	}

	err := l.SaveMeta(r.Context(), getCommandParameters(r), ledger.MetaTargetTypeAccount, chi.URLParam(r, "address"), m)
	if err != nil {
		ResponseError(w, r, err)
		return
	}

	sharedapi.NoContent(w)
}

func deleteAccountMetadata(w http.ResponseWriter, r *http.Request) {
	if err := LedgerFromContext(r.Context()).
		DeleteMetadata(
			r.Context(),
			getCommandParameters(r),
			ledger.MetaTargetTypeAccount,
			chi.URLParam(r, "address"),
			chi.URLParam(r, "key"),
		); err != nil {
		ResponseError(w, r, err)
		return
	}

	sharedapi.NoContent(w)
}
