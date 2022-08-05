package api

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	auth "github.com/numary/auth/pkg"
	"github.com/numary/auth/pkg/storage"
	"github.com/stretchr/testify/require"
	"github.com/zitadel/oidc/pkg/oidc"
	"github.com/zitadel/oidc/pkg/op"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func withDbAndClientRouter(t *testing.T, callback func(router *mux.Router, db *gorm.DB)) {
	db, err := storage.LoadGorm(sqlite.Open(":memory:"))
	require.NoError(t, err)
	require.NoError(t, storage.MigrateTables(context.Background(), db))

	router := mux.NewRouter()
	addClientRoutes(db, router)

	callback(router, db)
}

func TestCreateClient(t *testing.T) {

	type testCase struct {
		name           string
		options        auth.ClientOptions
		expectedClient auth.Client
	}
	for _, tc := range []testCase{
		{
			name: "confidential client",
			options: auth.ClientOptions{
				Name:                   "confidential client",
				RedirectUris:           []string{"http://localhost:8080"},
				Description:            "abc",
				PostLogoutRedirectUris: []string{"http://localhost:8080/logout"},
				Metadata: map[string]string{
					"foo": "bar",
				},
			},
			expectedClient: auth.Client{
				GrantTypes: auth.Array[oidc.GrantType]{
					oidc.GrantTypeCode,
					oidc.GrantTypeRefreshToken,
					oidc.GrantTypeClientCredentials,
				},
				AccessTokenType: op.AccessTokenTypeJWT,
				ResponseTypes: auth.Array[oidc.ResponseType]{
					oidc.ResponseTypeCode,
				},
				AuthMethod:             oidc.AuthMethodBasic,
				RedirectURIs:           []string{"http://localhost:8080"},
				Description:            "abc",
				PostLogoutRedirectUris: []string{"http://localhost:8080/logout"},
				Name:                   "confidential client",
				Metadata: map[string]string{
					"foo": "bar",
				},
			},
		},
		{
			name: "public client",
			options: auth.ClientOptions{
				Name:   "public client",
				Public: true,
			},
			expectedClient: auth.Client{
				GrantTypes: auth.Array[oidc.GrantType]{
					oidc.GrantTypeCode,
					oidc.GrantTypeRefreshToken,
				},
				AccessTokenType: op.AccessTokenTypeJWT,
				ResponseTypes: auth.Array[oidc.ResponseType]{
					oidc.ResponseTypeCode,
				},
				AuthMethod: oidc.AuthMethodNone,
				Name:       "public client",
			},
		},
	} {
		withDbAndClientRouter(t, func(router *mux.Router, db *gorm.DB) {
			req := httptest.NewRequest(http.MethodPost, "/clients", createJSONBuffer(t, tc.options))
			res := httptest.NewRecorder()

			router.ServeHTTP(res, req)

			require.Equal(t, http.StatusCreated, res.Code)

			createdClient := readTestResponse[client](t, res)
			require.NotEmpty(t, createdClient.ID)
			require.Equal(t, tc.options, createdClient.ClientOptions)

			tc.expectedClient.Id = createdClient.ID

			clientFromDatabase := auth.Client{}
			require.NoError(t, db.Find(&clientFromDatabase, "id = ?", createdClient.ID).Error)
			require.Equal(t, tc.expectedClient, clientFromDatabase)
		})
	}
}

func TestUpdateClient(t *testing.T) {

	type testCase struct {
		name           string
		options        auth.ClientOptions
		expectedClient auth.Client
	}
	for _, tc := range []testCase{
		{
			name: "confidential client",
			options: auth.ClientOptions{
				Name:                   "confidential client",
				RedirectUris:           []string{"http://localhost:8080"},
				Description:            "abc",
				PostLogoutRedirectUris: []string{"http://localhost:8080/logout"},
				Metadata: map[string]string{
					"foo": "bar",
				},
			},
			expectedClient: auth.Client{
				GrantTypes: auth.Array[oidc.GrantType]{
					oidc.GrantTypeCode,
					oidc.GrantTypeRefreshToken,
					oidc.GrantTypeClientCredentials,
				},
				AccessTokenType: op.AccessTokenTypeJWT,
				ResponseTypes: auth.Array[oidc.ResponseType]{
					oidc.ResponseTypeCode,
				},
				AuthMethod:             oidc.AuthMethodBasic,
				RedirectURIs:           []string{"http://localhost:8080"},
				Description:            "abc",
				PostLogoutRedirectUris: []string{"http://localhost:8080/logout"},
				Name:                   "confidential client",
				Metadata: map[string]string{
					"foo": "bar",
				},
			},
		},
		{
			name: "public client",
			options: auth.ClientOptions{
				Name:   "public client",
				Public: true,
			},
			expectedClient: auth.Client{
				GrantTypes: auth.Array[oidc.GrantType]{
					oidc.GrantTypeCode,
					oidc.GrantTypeRefreshToken,
				},
				AccessTokenType: op.AccessTokenTypeJWT,
				ResponseTypes: auth.Array[oidc.ResponseType]{
					oidc.ResponseTypeCode,
				},
				AuthMethod: oidc.AuthMethodNone,
				Name:       "public client",
			},
		},
	} {
		withDbAndClientRouter(t, func(router *mux.Router, db *gorm.DB) {

			initialClient := auth.NewClient(auth.ClientOptions{})
			require.NoError(t, db.Create(initialClient).Error)

			req := httptest.NewRequest(http.MethodPut, "/clients/"+initialClient.Id, createJSONBuffer(t, tc.options))
			res := httptest.NewRecorder()

			router.ServeHTTP(res, req)

			require.Equal(t, http.StatusOK, res.Code)

			updatedClient := readTestResponse[client](t, res)
			require.NotEmpty(t, updatedClient.ID)
			require.Equal(t, tc.options, updatedClient.ClientOptions)

			tc.expectedClient.Id = updatedClient.ID

			clientFromDatabase := auth.Client{}
			require.NoError(t, db.Find(&clientFromDatabase, "id = ?", updatedClient.ID).Error)
			require.Equal(t, tc.expectedClient, clientFromDatabase)
		})
	}
}

func TestListClients(t *testing.T) {
	withDbAndClientRouter(t, func(router *mux.Router, db *gorm.DB) {
		client1 := auth.NewClient(auth.ClientOptions{})
		require.NoError(t, db.Create(client1).Error)

		client2 := auth.NewClient(auth.ClientOptions{
			Metadata: map[string]string{
				"foo": "bar",
			},
		})
		require.NoError(t, db.Create(client2).Error)

		req := httptest.NewRequest(http.MethodGet, "/clients", nil)
		res := httptest.NewRecorder()

		router.ServeHTTP(res, req)

		require.Equal(t, http.StatusOK, res.Code)

		clients := readTestResponse[[]client](t, res)
		require.Len(t, clients, 2)
		require.Len(t, clients[1].Metadata, 1)
		require.Equal(t, clients[1].Metadata["foo"], "bar")
	})
}

func TestReadClient(t *testing.T) {
	withDbAndClientRouter(t, func(router *mux.Router, db *gorm.DB) {

		scope1 := auth.NewScope(auth.ScopeOptions{Label: "XXX"})
		require.NoError(t, db.Create(scope1).Error)

		opts := auth.ClientOptions{
			Metadata: map[string]string{
				"foo": "bar",
			},
		}
		client1 := auth.NewClient(opts)
		client1.Scopes = append(client1.Scopes, *scope1)
		require.NoError(t, db.Create(client1).Error)

		req := httptest.NewRequest(http.MethodGet, "/clients/"+client1.Id, nil)
		res := httptest.NewRecorder()

		router.ServeHTTP(res, req)

		require.Equal(t, http.StatusOK, res.Code)

		ret := readTestResponse[client](t, res)
		require.Equal(t, client{
			ClientOptions: opts,
			ID:            client1.Id,
			Scopes:        []string{scope1.ID},
		}, ret)
	})
}

func TestGenerateNewSecret(t *testing.T) {
	withDbAndClientRouter(t, func(router *mux.Router, db *gorm.DB) {
		client := auth.NewClient(auth.ClientOptions{})
		require.NoError(t, db.Create(client).Error)

		req := httptest.NewRequest(http.MethodPost, "/clients/"+client.Id+"/secrets", createJSONBuffer(t, secretCreate{
			Name: "secret1",
		}))
		res := httptest.NewRecorder()

		router.ServeHTTP(res, req)

		result := readTestResponse[secretCreateResult](t, res)
		require.NotEmpty(t, result.Clear)
		require.Equal(t, result.LastDigits, result.Clear[len(result.Clear)-4:])
		require.Equal(t, result.Name, "secret1")

		require.Equal(t, http.StatusOK, res.Code)
		require.NoError(t, db.First(client, "id = ?", client.Id).Error)
		require.Len(t, client.Secrets, 1)
		require.True(t, client.Secrets[0].Check(result.Clear))
	})
}

func TestDeleteSecret(t *testing.T) {
	withDbAndClientRouter(t, func(router *mux.Router, db *gorm.DB) {
		client := auth.NewClient(auth.ClientOptions{})
		secret, _ := client.GenerateNewSecret("testing")
		require.NoError(t, db.Create(client).Error)

		req := httptest.NewRequest(http.MethodDelete, "/clients/"+client.Id+"/secrets/"+secret.ID, nil)
		res := httptest.NewRecorder()

		router.ServeHTTP(res, req)

		require.Equal(t, http.StatusNoContent, res.Code)
		require.NoError(t, db.First(client, "id = ?", client.Id).Error)
		require.Len(t, client.Secrets, 0)
	})
}

func TestAddScope(t *testing.T) {
	withDbAndClientRouter(t, func(router *mux.Router, db *gorm.DB) {
		client := auth.NewClient(auth.ClientOptions{})
		require.NoError(t, db.Create(client).Error)

		scope := auth.NewScope(auth.ScopeOptions{Label: "XXX"})
		require.NoError(t, db.Create(scope).Error)

		req := httptest.NewRequest(http.MethodPut, "/clients/"+client.Id+"/scopes/"+scope.ID, nil)
		res := httptest.NewRecorder()

		router.ServeHTTP(res, req)

		require.Equal(t, http.StatusNoContent, res.Code)

		require.NoError(t, db.Preload("Scopes").First(client).Error)
		require.Len(t, client.Scopes, 1)
		require.Equal(t, *scope, client.Scopes[0])
	})
}

func TestRemoveScope(t *testing.T) {
	withDbAndClientRouter(t, func(router *mux.Router, db *gorm.DB) {

		scope := auth.NewScope(auth.ScopeOptions{Label: "XXX"})
		require.NoError(t, db.Create(scope).Error)

		client := auth.NewClient(auth.ClientOptions{})
		client.Scopes = append(client.Scopes, *scope)
		require.NoError(t, db.Create(client).Error)

		req := httptest.NewRequest(http.MethodDelete, "/clients/"+client.Id+"/scopes/"+scope.ID, nil)
		res := httptest.NewRecorder()

		router.ServeHTTP(res, req)

		require.Equal(t, http.StatusNoContent, res.Code)

		require.NoError(t, db.Preload("Scopes").First(client).Error)
		require.Len(t, client.Scopes, 0)
	})
}
