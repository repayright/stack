/*
Formance Stack API

Open, modular foundation for unique payments flows  # Introduction This API is documented in **OpenAPI format**.  # Authentication Formance Stack offers one forms of authentication:   - OAuth2 OAuth2 - an open protocol to allow secure authorization in a simple and standard method from web, mobile and desktop applications. <SecurityDefinitions /> 

API version: develop
Contact: support@formance.com
*/

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package formance

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"net/url"
	"strings"
)


type ScopesApi interface {

	/*
	AddTransientScope Add a transient scope to a scope

	Add a transient scope to a scope

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param scopeId Scope ID
	@param transientScopeId Transient scope ID
	@return ApiAddTransientScopeRequest
	*/
	AddTransientScope(ctx context.Context, scopeId string, transientScopeId string) ApiAddTransientScopeRequest

	// AddTransientScopeExecute executes the request
	AddTransientScopeExecute(r ApiAddTransientScopeRequest) (*http.Response, error)

	/*
	CreateScope Create scope

	Create scope

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@return ApiCreateScopeRequest
	*/
	CreateScope(ctx context.Context) ApiCreateScopeRequest

	// CreateScopeExecute executes the request
	//  @return CreateScopeResponse
	CreateScopeExecute(r ApiCreateScopeRequest) (*CreateScopeResponse, *http.Response, error)

	/*
	DeleteScope Delete scope

	Delete scope

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param scopeId Scope ID
	@return ApiDeleteScopeRequest
	*/
	DeleteScope(ctx context.Context, scopeId string) ApiDeleteScopeRequest

	// DeleteScopeExecute executes the request
	DeleteScopeExecute(r ApiDeleteScopeRequest) (*http.Response, error)

	/*
	DeleteTransientScope Delete a transient scope from a scope

	Delete a transient scope from a scope

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param scopeId Scope ID
	@param transientScopeId Transient scope ID
	@return ApiDeleteTransientScopeRequest
	*/
	DeleteTransientScope(ctx context.Context, scopeId string, transientScopeId string) ApiDeleteTransientScopeRequest

	// DeleteTransientScopeExecute executes the request
	DeleteTransientScopeExecute(r ApiDeleteTransientScopeRequest) (*http.Response, error)

	/*
	ListScopes List scopes

	List Scopes

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@return ApiListScopesRequest
	*/
	ListScopes(ctx context.Context) ApiListScopesRequest

	// ListScopesExecute executes the request
	//  @return ListScopesResponse
	ListScopesExecute(r ApiListScopesRequest) (*ListScopesResponse, *http.Response, error)

	/*
	ReadScope Read scope

	Read scope

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param scopeId Scope ID
	@return ApiReadScopeRequest
	*/
	ReadScope(ctx context.Context, scopeId string) ApiReadScopeRequest

	// ReadScopeExecute executes the request
	//  @return CreateScopeResponse
	ReadScopeExecute(r ApiReadScopeRequest) (*CreateScopeResponse, *http.Response, error)

	/*
	UpdateScope Update scope

	Update scope

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param scopeId Scope ID
	@return ApiUpdateScopeRequest
	*/
	UpdateScope(ctx context.Context, scopeId string) ApiUpdateScopeRequest

	// UpdateScopeExecute executes the request
	//  @return CreateScopeResponse
	UpdateScopeExecute(r ApiUpdateScopeRequest) (*CreateScopeResponse, *http.Response, error)
}

// ScopesApiService ScopesApi service
type ScopesApiService service

type ApiAddTransientScopeRequest struct {
	ctx context.Context
	ApiService ScopesApi
	scopeId string
	transientScopeId string
}

func (r ApiAddTransientScopeRequest) Execute() (*http.Response, error) {
	return r.ApiService.AddTransientScopeExecute(r)
}

/*
AddTransientScope Add a transient scope to a scope

Add a transient scope to a scope

 @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 @param scopeId Scope ID
 @param transientScopeId Transient scope ID
 @return ApiAddTransientScopeRequest
*/
func (a *ScopesApiService) AddTransientScope(ctx context.Context, scopeId string, transientScopeId string) ApiAddTransientScopeRequest {
	return ApiAddTransientScopeRequest{
		ApiService: a,
		ctx: ctx,
		scopeId: scopeId,
		transientScopeId: transientScopeId,
	}
}

// Execute executes the request
func (a *ScopesApiService) AddTransientScopeExecute(r ApiAddTransientScopeRequest) (*http.Response, error) {
	var (
		localVarHTTPMethod   = http.MethodPut
		localVarPostBody     interface{}
		formFiles            []formFile
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "ScopesApiService.AddTransientScope")
	if err != nil {
		return nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/auth/scopes/{scopeId}/transient/{transientScopeId}"
	localVarPath = strings.Replace(localVarPath, "{"+"scopeId"+"}", url.PathEscape(parameterValueToString(r.scopeId, "scopeId")), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"transientScopeId"+"}", url.PathEscape(parameterValueToString(r.transientScopeId, "transientScopeId")), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}

	// to determine the Content-Type header
	localVarHTTPContentTypes := []string{}

	// set Content-Type header
	localVarHTTPContentType := selectHeaderContentType(localVarHTTPContentTypes)
	if localVarHTTPContentType != "" {
		localVarHeaderParams["Content-Type"] = localVarHTTPContentType
	}

	// to determine the Accept header
	localVarHTTPHeaderAccepts := []string{}

	// set Accept header
	localVarHTTPHeaderAccept := selectHeaderAccept(localVarHTTPHeaderAccepts)
	if localVarHTTPHeaderAccept != "" {
		localVarHeaderParams["Accept"] = localVarHTTPHeaderAccept
	}
	req, err := a.client.prepareRequest(r.ctx, localVarPath, localVarHTTPMethod, localVarPostBody, localVarHeaderParams, localVarQueryParams, localVarFormParams, formFiles)
	if err != nil {
		return nil, err
	}

	localVarHTTPResponse, err := a.client.callAPI(req)
	if err != nil || localVarHTTPResponse == nil {
		return localVarHTTPResponse, err
	}

	localVarBody, err := io.ReadAll(localVarHTTPResponse.Body)
	localVarHTTPResponse.Body.Close()
	localVarHTTPResponse.Body = io.NopCloser(bytes.NewBuffer(localVarBody))
	if err != nil {
		return localVarHTTPResponse, err
	}

	if localVarHTTPResponse.StatusCode >= 300 {
		newErr := &GenericOpenAPIError{
			body:  localVarBody,
			error: localVarHTTPResponse.Status,
		}
		return localVarHTTPResponse, newErr
	}

	return localVarHTTPResponse, nil
}

type ApiCreateScopeRequest struct {
	ctx context.Context
	ApiService ScopesApi
	body *ScopeOptions
}

func (r ApiCreateScopeRequest) Body(body ScopeOptions) ApiCreateScopeRequest {
	r.body = &body
	return r
}

func (r ApiCreateScopeRequest) Execute() (*CreateScopeResponse, *http.Response, error) {
	return r.ApiService.CreateScopeExecute(r)
}

/*
CreateScope Create scope

Create scope

 @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 @return ApiCreateScopeRequest
*/
func (a *ScopesApiService) CreateScope(ctx context.Context) ApiCreateScopeRequest {
	return ApiCreateScopeRequest{
		ApiService: a,
		ctx: ctx,
	}
}

// Execute executes the request
//  @return CreateScopeResponse
func (a *ScopesApiService) CreateScopeExecute(r ApiCreateScopeRequest) (*CreateScopeResponse, *http.Response, error) {
	var (
		localVarHTTPMethod   = http.MethodPost
		localVarPostBody     interface{}
		formFiles            []formFile
		localVarReturnValue  *CreateScopeResponse
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "ScopesApiService.CreateScope")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/auth/scopes"

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}

	// to determine the Content-Type header
	localVarHTTPContentTypes := []string{"application/json"}

	// set Content-Type header
	localVarHTTPContentType := selectHeaderContentType(localVarHTTPContentTypes)
	if localVarHTTPContentType != "" {
		localVarHeaderParams["Content-Type"] = localVarHTTPContentType
	}

	// to determine the Accept header
	localVarHTTPHeaderAccepts := []string{"application/json"}

	// set Accept header
	localVarHTTPHeaderAccept := selectHeaderAccept(localVarHTTPHeaderAccepts)
	if localVarHTTPHeaderAccept != "" {
		localVarHeaderParams["Accept"] = localVarHTTPHeaderAccept
	}
	// body params
	localVarPostBody = r.body
	req, err := a.client.prepareRequest(r.ctx, localVarPath, localVarHTTPMethod, localVarPostBody, localVarHeaderParams, localVarQueryParams, localVarFormParams, formFiles)
	if err != nil {
		return localVarReturnValue, nil, err
	}

	localVarHTTPResponse, err := a.client.callAPI(req)
	if err != nil || localVarHTTPResponse == nil {
		return localVarReturnValue, localVarHTTPResponse, err
	}

	localVarBody, err := io.ReadAll(localVarHTTPResponse.Body)
	localVarHTTPResponse.Body.Close()
	localVarHTTPResponse.Body = io.NopCloser(bytes.NewBuffer(localVarBody))
	if err != nil {
		return localVarReturnValue, localVarHTTPResponse, err
	}

	if localVarHTTPResponse.StatusCode >= 300 {
		newErr := &GenericOpenAPIError{
			body:  localVarBody,
			error: localVarHTTPResponse.Status,
		}
		return localVarReturnValue, localVarHTTPResponse, newErr
	}

	err = a.client.decode(&localVarReturnValue, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
	if err != nil {
		newErr := &GenericOpenAPIError{
			body:  localVarBody,
			error: err.Error(),
		}
		return localVarReturnValue, localVarHTTPResponse, newErr
	}

	return localVarReturnValue, localVarHTTPResponse, nil
}

type ApiDeleteScopeRequest struct {
	ctx context.Context
	ApiService ScopesApi
	scopeId string
}

func (r ApiDeleteScopeRequest) Execute() (*http.Response, error) {
	return r.ApiService.DeleteScopeExecute(r)
}

/*
DeleteScope Delete scope

Delete scope

 @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 @param scopeId Scope ID
 @return ApiDeleteScopeRequest
*/
func (a *ScopesApiService) DeleteScope(ctx context.Context, scopeId string) ApiDeleteScopeRequest {
	return ApiDeleteScopeRequest{
		ApiService: a,
		ctx: ctx,
		scopeId: scopeId,
	}
}

// Execute executes the request
func (a *ScopesApiService) DeleteScopeExecute(r ApiDeleteScopeRequest) (*http.Response, error) {
	var (
		localVarHTTPMethod   = http.MethodDelete
		localVarPostBody     interface{}
		formFiles            []formFile
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "ScopesApiService.DeleteScope")
	if err != nil {
		return nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/auth/scopes/{scopeId}"
	localVarPath = strings.Replace(localVarPath, "{"+"scopeId"+"}", url.PathEscape(parameterValueToString(r.scopeId, "scopeId")), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}

	// to determine the Content-Type header
	localVarHTTPContentTypes := []string{}

	// set Content-Type header
	localVarHTTPContentType := selectHeaderContentType(localVarHTTPContentTypes)
	if localVarHTTPContentType != "" {
		localVarHeaderParams["Content-Type"] = localVarHTTPContentType
	}

	// to determine the Accept header
	localVarHTTPHeaderAccepts := []string{}

	// set Accept header
	localVarHTTPHeaderAccept := selectHeaderAccept(localVarHTTPHeaderAccepts)
	if localVarHTTPHeaderAccept != "" {
		localVarHeaderParams["Accept"] = localVarHTTPHeaderAccept
	}
	req, err := a.client.prepareRequest(r.ctx, localVarPath, localVarHTTPMethod, localVarPostBody, localVarHeaderParams, localVarQueryParams, localVarFormParams, formFiles)
	if err != nil {
		return nil, err
	}

	localVarHTTPResponse, err := a.client.callAPI(req)
	if err != nil || localVarHTTPResponse == nil {
		return localVarHTTPResponse, err
	}

	localVarBody, err := io.ReadAll(localVarHTTPResponse.Body)
	localVarHTTPResponse.Body.Close()
	localVarHTTPResponse.Body = io.NopCloser(bytes.NewBuffer(localVarBody))
	if err != nil {
		return localVarHTTPResponse, err
	}

	if localVarHTTPResponse.StatusCode >= 300 {
		newErr := &GenericOpenAPIError{
			body:  localVarBody,
			error: localVarHTTPResponse.Status,
		}
		return localVarHTTPResponse, newErr
	}

	return localVarHTTPResponse, nil
}

type ApiDeleteTransientScopeRequest struct {
	ctx context.Context
	ApiService ScopesApi
	scopeId string
	transientScopeId string
}

func (r ApiDeleteTransientScopeRequest) Execute() (*http.Response, error) {
	return r.ApiService.DeleteTransientScopeExecute(r)
}

/*
DeleteTransientScope Delete a transient scope from a scope

Delete a transient scope from a scope

 @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 @param scopeId Scope ID
 @param transientScopeId Transient scope ID
 @return ApiDeleteTransientScopeRequest
*/
func (a *ScopesApiService) DeleteTransientScope(ctx context.Context, scopeId string, transientScopeId string) ApiDeleteTransientScopeRequest {
	return ApiDeleteTransientScopeRequest{
		ApiService: a,
		ctx: ctx,
		scopeId: scopeId,
		transientScopeId: transientScopeId,
	}
}

// Execute executes the request
func (a *ScopesApiService) DeleteTransientScopeExecute(r ApiDeleteTransientScopeRequest) (*http.Response, error) {
	var (
		localVarHTTPMethod   = http.MethodDelete
		localVarPostBody     interface{}
		formFiles            []formFile
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "ScopesApiService.DeleteTransientScope")
	if err != nil {
		return nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/auth/scopes/{scopeId}/transient/{transientScopeId}"
	localVarPath = strings.Replace(localVarPath, "{"+"scopeId"+"}", url.PathEscape(parameterValueToString(r.scopeId, "scopeId")), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"transientScopeId"+"}", url.PathEscape(parameterValueToString(r.transientScopeId, "transientScopeId")), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}

	// to determine the Content-Type header
	localVarHTTPContentTypes := []string{}

	// set Content-Type header
	localVarHTTPContentType := selectHeaderContentType(localVarHTTPContentTypes)
	if localVarHTTPContentType != "" {
		localVarHeaderParams["Content-Type"] = localVarHTTPContentType
	}

	// to determine the Accept header
	localVarHTTPHeaderAccepts := []string{}

	// set Accept header
	localVarHTTPHeaderAccept := selectHeaderAccept(localVarHTTPHeaderAccepts)
	if localVarHTTPHeaderAccept != "" {
		localVarHeaderParams["Accept"] = localVarHTTPHeaderAccept
	}
	req, err := a.client.prepareRequest(r.ctx, localVarPath, localVarHTTPMethod, localVarPostBody, localVarHeaderParams, localVarQueryParams, localVarFormParams, formFiles)
	if err != nil {
		return nil, err
	}

	localVarHTTPResponse, err := a.client.callAPI(req)
	if err != nil || localVarHTTPResponse == nil {
		return localVarHTTPResponse, err
	}

	localVarBody, err := io.ReadAll(localVarHTTPResponse.Body)
	localVarHTTPResponse.Body.Close()
	localVarHTTPResponse.Body = io.NopCloser(bytes.NewBuffer(localVarBody))
	if err != nil {
		return localVarHTTPResponse, err
	}

	if localVarHTTPResponse.StatusCode >= 300 {
		newErr := &GenericOpenAPIError{
			body:  localVarBody,
			error: localVarHTTPResponse.Status,
		}
		return localVarHTTPResponse, newErr
	}

	return localVarHTTPResponse, nil
}

type ApiListScopesRequest struct {
	ctx context.Context
	ApiService ScopesApi
}

func (r ApiListScopesRequest) Execute() (*ListScopesResponse, *http.Response, error) {
	return r.ApiService.ListScopesExecute(r)
}

/*
ListScopes List scopes

List Scopes

 @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 @return ApiListScopesRequest
*/
func (a *ScopesApiService) ListScopes(ctx context.Context) ApiListScopesRequest {
	return ApiListScopesRequest{
		ApiService: a,
		ctx: ctx,
	}
}

// Execute executes the request
//  @return ListScopesResponse
func (a *ScopesApiService) ListScopesExecute(r ApiListScopesRequest) (*ListScopesResponse, *http.Response, error) {
	var (
		localVarHTTPMethod   = http.MethodGet
		localVarPostBody     interface{}
		formFiles            []formFile
		localVarReturnValue  *ListScopesResponse
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "ScopesApiService.ListScopes")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/auth/scopes"

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}

	// to determine the Content-Type header
	localVarHTTPContentTypes := []string{}

	// set Content-Type header
	localVarHTTPContentType := selectHeaderContentType(localVarHTTPContentTypes)
	if localVarHTTPContentType != "" {
		localVarHeaderParams["Content-Type"] = localVarHTTPContentType
	}

	// to determine the Accept header
	localVarHTTPHeaderAccepts := []string{"application/json"}

	// set Accept header
	localVarHTTPHeaderAccept := selectHeaderAccept(localVarHTTPHeaderAccepts)
	if localVarHTTPHeaderAccept != "" {
		localVarHeaderParams["Accept"] = localVarHTTPHeaderAccept
	}
	req, err := a.client.prepareRequest(r.ctx, localVarPath, localVarHTTPMethod, localVarPostBody, localVarHeaderParams, localVarQueryParams, localVarFormParams, formFiles)
	if err != nil {
		return localVarReturnValue, nil, err
	}

	localVarHTTPResponse, err := a.client.callAPI(req)
	if err != nil || localVarHTTPResponse == nil {
		return localVarReturnValue, localVarHTTPResponse, err
	}

	localVarBody, err := io.ReadAll(localVarHTTPResponse.Body)
	localVarHTTPResponse.Body.Close()
	localVarHTTPResponse.Body = io.NopCloser(bytes.NewBuffer(localVarBody))
	if err != nil {
		return localVarReturnValue, localVarHTTPResponse, err
	}

	if localVarHTTPResponse.StatusCode >= 300 {
		newErr := &GenericOpenAPIError{
			body:  localVarBody,
			error: localVarHTTPResponse.Status,
		}
		return localVarReturnValue, localVarHTTPResponse, newErr
	}

	err = a.client.decode(&localVarReturnValue, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
	if err != nil {
		newErr := &GenericOpenAPIError{
			body:  localVarBody,
			error: err.Error(),
		}
		return localVarReturnValue, localVarHTTPResponse, newErr
	}

	return localVarReturnValue, localVarHTTPResponse, nil
}

type ApiReadScopeRequest struct {
	ctx context.Context
	ApiService ScopesApi
	scopeId string
}

func (r ApiReadScopeRequest) Execute() (*CreateScopeResponse, *http.Response, error) {
	return r.ApiService.ReadScopeExecute(r)
}

/*
ReadScope Read scope

Read scope

 @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 @param scopeId Scope ID
 @return ApiReadScopeRequest
*/
func (a *ScopesApiService) ReadScope(ctx context.Context, scopeId string) ApiReadScopeRequest {
	return ApiReadScopeRequest{
		ApiService: a,
		ctx: ctx,
		scopeId: scopeId,
	}
}

// Execute executes the request
//  @return CreateScopeResponse
func (a *ScopesApiService) ReadScopeExecute(r ApiReadScopeRequest) (*CreateScopeResponse, *http.Response, error) {
	var (
		localVarHTTPMethod   = http.MethodGet
		localVarPostBody     interface{}
		formFiles            []formFile
		localVarReturnValue  *CreateScopeResponse
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "ScopesApiService.ReadScope")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/auth/scopes/{scopeId}"
	localVarPath = strings.Replace(localVarPath, "{"+"scopeId"+"}", url.PathEscape(parameterValueToString(r.scopeId, "scopeId")), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}

	// to determine the Content-Type header
	localVarHTTPContentTypes := []string{}

	// set Content-Type header
	localVarHTTPContentType := selectHeaderContentType(localVarHTTPContentTypes)
	if localVarHTTPContentType != "" {
		localVarHeaderParams["Content-Type"] = localVarHTTPContentType
	}

	// to determine the Accept header
	localVarHTTPHeaderAccepts := []string{"application/json"}

	// set Accept header
	localVarHTTPHeaderAccept := selectHeaderAccept(localVarHTTPHeaderAccepts)
	if localVarHTTPHeaderAccept != "" {
		localVarHeaderParams["Accept"] = localVarHTTPHeaderAccept
	}
	req, err := a.client.prepareRequest(r.ctx, localVarPath, localVarHTTPMethod, localVarPostBody, localVarHeaderParams, localVarQueryParams, localVarFormParams, formFiles)
	if err != nil {
		return localVarReturnValue, nil, err
	}

	localVarHTTPResponse, err := a.client.callAPI(req)
	if err != nil || localVarHTTPResponse == nil {
		return localVarReturnValue, localVarHTTPResponse, err
	}

	localVarBody, err := io.ReadAll(localVarHTTPResponse.Body)
	localVarHTTPResponse.Body.Close()
	localVarHTTPResponse.Body = io.NopCloser(bytes.NewBuffer(localVarBody))
	if err != nil {
		return localVarReturnValue, localVarHTTPResponse, err
	}

	if localVarHTTPResponse.StatusCode >= 300 {
		newErr := &GenericOpenAPIError{
			body:  localVarBody,
			error: localVarHTTPResponse.Status,
		}
		return localVarReturnValue, localVarHTTPResponse, newErr
	}

	err = a.client.decode(&localVarReturnValue, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
	if err != nil {
		newErr := &GenericOpenAPIError{
			body:  localVarBody,
			error: err.Error(),
		}
		return localVarReturnValue, localVarHTTPResponse, newErr
	}

	return localVarReturnValue, localVarHTTPResponse, nil
}

type ApiUpdateScopeRequest struct {
	ctx context.Context
	ApiService ScopesApi
	scopeId string
	body *ScopeOptions
}

func (r ApiUpdateScopeRequest) Body(body ScopeOptions) ApiUpdateScopeRequest {
	r.body = &body
	return r
}

func (r ApiUpdateScopeRequest) Execute() (*CreateScopeResponse, *http.Response, error) {
	return r.ApiService.UpdateScopeExecute(r)
}

/*
UpdateScope Update scope

Update scope

 @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 @param scopeId Scope ID
 @return ApiUpdateScopeRequest
*/
func (a *ScopesApiService) UpdateScope(ctx context.Context, scopeId string) ApiUpdateScopeRequest {
	return ApiUpdateScopeRequest{
		ApiService: a,
		ctx: ctx,
		scopeId: scopeId,
	}
}

// Execute executes the request
//  @return CreateScopeResponse
func (a *ScopesApiService) UpdateScopeExecute(r ApiUpdateScopeRequest) (*CreateScopeResponse, *http.Response, error) {
	var (
		localVarHTTPMethod   = http.MethodPut
		localVarPostBody     interface{}
		formFiles            []formFile
		localVarReturnValue  *CreateScopeResponse
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "ScopesApiService.UpdateScope")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/auth/scopes/{scopeId}"
	localVarPath = strings.Replace(localVarPath, "{"+"scopeId"+"}", url.PathEscape(parameterValueToString(r.scopeId, "scopeId")), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}

	// to determine the Content-Type header
	localVarHTTPContentTypes := []string{"application/json"}

	// set Content-Type header
	localVarHTTPContentType := selectHeaderContentType(localVarHTTPContentTypes)
	if localVarHTTPContentType != "" {
		localVarHeaderParams["Content-Type"] = localVarHTTPContentType
	}

	// to determine the Accept header
	localVarHTTPHeaderAccepts := []string{"application/json"}

	// set Accept header
	localVarHTTPHeaderAccept := selectHeaderAccept(localVarHTTPHeaderAccepts)
	if localVarHTTPHeaderAccept != "" {
		localVarHeaderParams["Accept"] = localVarHTTPHeaderAccept
	}
	// body params
	localVarPostBody = r.body
	req, err := a.client.prepareRequest(r.ctx, localVarPath, localVarHTTPMethod, localVarPostBody, localVarHeaderParams, localVarQueryParams, localVarFormParams, formFiles)
	if err != nil {
		return localVarReturnValue, nil, err
	}

	localVarHTTPResponse, err := a.client.callAPI(req)
	if err != nil || localVarHTTPResponse == nil {
		return localVarReturnValue, localVarHTTPResponse, err
	}

	localVarBody, err := io.ReadAll(localVarHTTPResponse.Body)
	localVarHTTPResponse.Body.Close()
	localVarHTTPResponse.Body = io.NopCloser(bytes.NewBuffer(localVarBody))
	if err != nil {
		return localVarReturnValue, localVarHTTPResponse, err
	}

	if localVarHTTPResponse.StatusCode >= 300 {
		newErr := &GenericOpenAPIError{
			body:  localVarBody,
			error: localVarHTTPResponse.Status,
		}
		return localVarReturnValue, localVarHTTPResponse, newErr
	}

	err = a.client.decode(&localVarReturnValue, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
	if err != nil {
		newErr := &GenericOpenAPIError{
			body:  localVarBody,
			error: err.Error(),
		}
		return localVarReturnValue, localVarHTTPResponse, newErr
	}

	return localVarReturnValue, localVarHTTPResponse, nil
}
