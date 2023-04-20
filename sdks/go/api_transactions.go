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
	"time"
)


type TransactionsApi interface {

	/*
	AddMetadataOnTransaction Set the metadata of a transaction by its ID

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param ledger Name of the ledger.
	@param txid Transaction ID.
	@return ApiAddMetadataOnTransactionRequest
	*/
	AddMetadataOnTransaction(ctx context.Context, ledger string, txid int64) ApiAddMetadataOnTransactionRequest

	// AddMetadataOnTransactionExecute executes the request
	AddMetadataOnTransactionExecute(r ApiAddMetadataOnTransactionRequest) (*http.Response, error)

	/*
	CountTransactions Count the transactions from a ledger

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param ledger Name of the ledger.
	@return ApiCountTransactionsRequest
	*/
	CountTransactions(ctx context.Context, ledger string) ApiCountTransactionsRequest

	// CountTransactionsExecute executes the request
	CountTransactionsExecute(r ApiCountTransactionsRequest) (*http.Response, error)

	/*
	CreateTransaction Create a new transaction to a ledger

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param ledger Name of the ledger.
	@return ApiCreateTransactionRequest
	*/
	CreateTransaction(ctx context.Context, ledger string) ApiCreateTransactionRequest

	// CreateTransactionExecute executes the request
	//  @return CreateTransactionResponse
	CreateTransactionExecute(r ApiCreateTransactionRequest) (*CreateTransactionResponse, *http.Response, error)

	/*
	GetTransaction Get transaction from a ledger by its ID

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param ledger Name of the ledger.
	@param txid Transaction ID.
	@return ApiGetTransactionRequest
	*/
	GetTransaction(ctx context.Context, ledger string, txid int64) ApiGetTransactionRequest

	// GetTransactionExecute executes the request
	//  @return GetTransactionResponse
	GetTransactionExecute(r ApiGetTransactionRequest) (*GetTransactionResponse, *http.Response, error)

	/*
	ListTransactions List transactions from a ledger

	List transactions from a ledger, sorted by txid in descending order.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param ledger Name of the ledger.
	@return ApiListTransactionsRequest
	*/
	ListTransactions(ctx context.Context, ledger string) ApiListTransactionsRequest

	// ListTransactionsExecute executes the request
	//  @return TransactionsCursorResponse
	ListTransactionsExecute(r ApiListTransactionsRequest) (*TransactionsCursorResponse, *http.Response, error)

	/*
	RevertTransaction Revert a ledger transaction by its ID

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param ledger Name of the ledger.
	@param txid Transaction ID.
	@return ApiRevertTransactionRequest
	*/
	RevertTransaction(ctx context.Context, ledger string, txid int64) ApiRevertTransactionRequest

	// RevertTransactionExecute executes the request
	//  @return CreateTransactionResponse
	RevertTransactionExecute(r ApiRevertTransactionRequest) (*CreateTransactionResponse, *http.Response, error)
}

// TransactionsApiService TransactionsApi service
type TransactionsApiService service

type ApiAddMetadataOnTransactionRequest struct {
	ctx context.Context
	ApiService TransactionsApi
	ledger string
	txid int64
	dryRun *bool
	async *bool
	idempotencyKey *string
	requestBody *map[string]string
}

// Set the dryRun mode. Dry run mode doesn&#39;t add the logs to the database or publish a message to the message broker.
func (r ApiAddMetadataOnTransactionRequest) DryRun(dryRun bool) ApiAddMetadataOnTransactionRequest {
	r.dryRun = &dryRun
	return r
}

// Set async mode.
func (r ApiAddMetadataOnTransactionRequest) Async(async bool) ApiAddMetadataOnTransactionRequest {
	r.async = &async
	return r
}

// Use an idempotency key
func (r ApiAddMetadataOnTransactionRequest) IdempotencyKey(idempotencyKey string) ApiAddMetadataOnTransactionRequest {
	r.idempotencyKey = &idempotencyKey
	return r
}

// metadata
func (r ApiAddMetadataOnTransactionRequest) RequestBody(requestBody map[string]string) ApiAddMetadataOnTransactionRequest {
	r.requestBody = &requestBody
	return r
}

func (r ApiAddMetadataOnTransactionRequest) Execute() (*http.Response, error) {
	return r.ApiService.AddMetadataOnTransactionExecute(r)
}

/*
AddMetadataOnTransaction Set the metadata of a transaction by its ID

 @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 @param ledger Name of the ledger.
 @param txid Transaction ID.
 @return ApiAddMetadataOnTransactionRequest
*/
func (a *TransactionsApiService) AddMetadataOnTransaction(ctx context.Context, ledger string, txid int64) ApiAddMetadataOnTransactionRequest {
	return ApiAddMetadataOnTransactionRequest{
		ApiService: a,
		ctx: ctx,
		ledger: ledger,
		txid: txid,
	}
}

// Execute executes the request
func (a *TransactionsApiService) AddMetadataOnTransactionExecute(r ApiAddMetadataOnTransactionRequest) (*http.Response, error) {
	var (
		localVarHTTPMethod   = http.MethodPost
		localVarPostBody     interface{}
		formFiles            []formFile
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "TransactionsApiService.AddMetadataOnTransaction")
	if err != nil {
		return nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/ledger/{ledger}/transactions/{txid}/metadata"
	localVarPath = strings.Replace(localVarPath, "{"+"ledger"+"}", url.PathEscape(parameterValueToString(r.ledger, "ledger")), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"txid"+"}", url.PathEscape(parameterValueToString(r.txid, "txid")), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}
	if r.txid < 0 {
		return nil, reportError("txid must be greater than 0")
	}

	if r.dryRun != nil {
		parameterAddToHeaderOrQuery(localVarQueryParams, "dryRun", r.dryRun, "")
	}
	if r.async != nil {
		parameterAddToHeaderOrQuery(localVarQueryParams, "async", r.async, "")
	}
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
	if r.idempotencyKey != nil {
		parameterAddToHeaderOrQuery(localVarHeaderParams, "Idempotency-Key", r.idempotencyKey, "")
	}
	// body params
	localVarPostBody = r.requestBody
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
			var v ErrorResponse
			err = a.client.decode(&v, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
			if err != nil {
				newErr.error = err.Error()
				return localVarHTTPResponse, newErr
			}
					newErr.error = formatErrorMessage(localVarHTTPResponse.Status, &v)
					newErr.model = v
		return localVarHTTPResponse, newErr
	}

	return localVarHTTPResponse, nil
}

type ApiCountTransactionsRequest struct {
	ctx context.Context
	ApiService TransactionsApi
	ledger string
	reference *string
	account *string
	source *string
	destination *string
	startTime *time.Time
	endTime *time.Time
	metadata *map[string]string
}

// Filter transactions by reference field.
func (r ApiCountTransactionsRequest) Reference(reference string) ApiCountTransactionsRequest {
	r.reference = &reference
	return r
}

// Filter transactions with postings involving given account, either as source or destination (regular expression placed between ^ and $).
func (r ApiCountTransactionsRequest) Account(account string) ApiCountTransactionsRequest {
	r.account = &account
	return r
}

// Filter transactions with postings involving given account at source (regular expression placed between ^ and $).
func (r ApiCountTransactionsRequest) Source(source string) ApiCountTransactionsRequest {
	r.source = &source
	return r
}

// Filter transactions with postings involving given account at destination (regular expression placed between ^ and $).
func (r ApiCountTransactionsRequest) Destination(destination string) ApiCountTransactionsRequest {
	r.destination = &destination
	return r
}

// Filter transactions that occurred after this timestamp. The format is RFC3339 and is inclusive (for example, \&quot;2023-01-02T15:04:01Z\&quot; includes the first second of 4th minute). 
func (r ApiCountTransactionsRequest) StartTime(startTime time.Time) ApiCountTransactionsRequest {
	r.startTime = &startTime
	return r
}

// Filter transactions that occurred before this timestamp. The format is RFC3339 and is exclusive (for example, \&quot;2023-01-02T15:04:01Z\&quot; excludes the first second of 4th minute). 
func (r ApiCountTransactionsRequest) EndTime(endTime time.Time) ApiCountTransactionsRequest {
	r.endTime = &endTime
	return r
}

// Filter transactions by metadata key value pairs. Nested objects can be used as seen in the example below.
func (r ApiCountTransactionsRequest) Metadata(metadata map[string]string) ApiCountTransactionsRequest {
	r.metadata = &metadata
	return r
}

func (r ApiCountTransactionsRequest) Execute() (*http.Response, error) {
	return r.ApiService.CountTransactionsExecute(r)
}

/*
CountTransactions Count the transactions from a ledger

 @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 @param ledger Name of the ledger.
 @return ApiCountTransactionsRequest
*/
func (a *TransactionsApiService) CountTransactions(ctx context.Context, ledger string) ApiCountTransactionsRequest {
	return ApiCountTransactionsRequest{
		ApiService: a,
		ctx: ctx,
		ledger: ledger,
	}
}

// Execute executes the request
func (a *TransactionsApiService) CountTransactionsExecute(r ApiCountTransactionsRequest) (*http.Response, error) {
	var (
		localVarHTTPMethod   = http.MethodHead
		localVarPostBody     interface{}
		formFiles            []formFile
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "TransactionsApiService.CountTransactions")
	if err != nil {
		return nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/ledger/{ledger}/transactions"
	localVarPath = strings.Replace(localVarPath, "{"+"ledger"+"}", url.PathEscape(parameterValueToString(r.ledger, "ledger")), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}

	if r.reference != nil {
		parameterAddToHeaderOrQuery(localVarQueryParams, "reference", r.reference, "")
	}
	if r.account != nil {
		parameterAddToHeaderOrQuery(localVarQueryParams, "account", r.account, "")
	}
	if r.source != nil {
		parameterAddToHeaderOrQuery(localVarQueryParams, "source", r.source, "")
	}
	if r.destination != nil {
		parameterAddToHeaderOrQuery(localVarQueryParams, "destination", r.destination, "")
	}
	if r.startTime != nil {
		parameterAddToHeaderOrQuery(localVarQueryParams, "startTime", r.startTime, "")
	}
	if r.endTime != nil {
		parameterAddToHeaderOrQuery(localVarQueryParams, "endTime", r.endTime, "")
	}
	if r.metadata != nil {
		parameterAddToHeaderOrQuery(localVarQueryParams, "metadata", r.metadata, "")
	}
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
			var v ErrorResponse
			err = a.client.decode(&v, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
			if err != nil {
				newErr.error = err.Error()
				return localVarHTTPResponse, newErr
			}
					newErr.error = formatErrorMessage(localVarHTTPResponse.Status, &v)
					newErr.model = v
		return localVarHTTPResponse, newErr
	}

	return localVarHTTPResponse, nil
}

type ApiCreateTransactionRequest struct {
	ctx context.Context
	ApiService TransactionsApi
	ledger string
	postTransaction *PostTransaction
	dryRun *bool
	async *bool
	idempotencyKey *string
}

// The request body must contain at least one of the following objects:   - &#x60;postings&#x60;: suitable for simple transactions   - &#x60;script&#x60;: enabling more complex transactions with Numscript 
func (r ApiCreateTransactionRequest) PostTransaction(postTransaction PostTransaction) ApiCreateTransactionRequest {
	r.postTransaction = &postTransaction
	return r
}

// Set the dryRun mode. dry run mode doesn&#39;t add the logs to the database or publish a message to the message broker.
func (r ApiCreateTransactionRequest) DryRun(dryRun bool) ApiCreateTransactionRequest {
	r.dryRun = &dryRun
	return r
}

// Set async mode.
func (r ApiCreateTransactionRequest) Async(async bool) ApiCreateTransactionRequest {
	r.async = &async
	return r
}

// Use an idempotency key
func (r ApiCreateTransactionRequest) IdempotencyKey(idempotencyKey string) ApiCreateTransactionRequest {
	r.idempotencyKey = &idempotencyKey
	return r
}

func (r ApiCreateTransactionRequest) Execute() (*CreateTransactionResponse, *http.Response, error) {
	return r.ApiService.CreateTransactionExecute(r)
}

/*
CreateTransaction Create a new transaction to a ledger

 @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 @param ledger Name of the ledger.
 @return ApiCreateTransactionRequest
*/
func (a *TransactionsApiService) CreateTransaction(ctx context.Context, ledger string) ApiCreateTransactionRequest {
	return ApiCreateTransactionRequest{
		ApiService: a,
		ctx: ctx,
		ledger: ledger,
	}
}

// Execute executes the request
//  @return CreateTransactionResponse
func (a *TransactionsApiService) CreateTransactionExecute(r ApiCreateTransactionRequest) (*CreateTransactionResponse, *http.Response, error) {
	var (
		localVarHTTPMethod   = http.MethodPost
		localVarPostBody     interface{}
		formFiles            []formFile
		localVarReturnValue  *CreateTransactionResponse
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "TransactionsApiService.CreateTransaction")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/ledger/{ledger}/transactions"
	localVarPath = strings.Replace(localVarPath, "{"+"ledger"+"}", url.PathEscape(parameterValueToString(r.ledger, "ledger")), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}
	if r.postTransaction == nil {
		return localVarReturnValue, nil, reportError("postTransaction is required and must be specified")
	}

	if r.dryRun != nil {
		parameterAddToHeaderOrQuery(localVarQueryParams, "dryRun", r.dryRun, "")
	}
	if r.async != nil {
		parameterAddToHeaderOrQuery(localVarQueryParams, "async", r.async, "")
	}
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
	if r.idempotencyKey != nil {
		parameterAddToHeaderOrQuery(localVarHeaderParams, "Idempotency-Key", r.idempotencyKey, "")
	}
	// body params
	localVarPostBody = r.postTransaction
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
			var v ErrorResponse
			err = a.client.decode(&v, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
			if err != nil {
				newErr.error = err.Error()
				return localVarReturnValue, localVarHTTPResponse, newErr
			}
					newErr.error = formatErrorMessage(localVarHTTPResponse.Status, &v)
					newErr.model = v
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

type ApiGetTransactionRequest struct {
	ctx context.Context
	ApiService TransactionsApi
	ledger string
	txid int64
}

func (r ApiGetTransactionRequest) Execute() (*GetTransactionResponse, *http.Response, error) {
	return r.ApiService.GetTransactionExecute(r)
}

/*
GetTransaction Get transaction from a ledger by its ID

 @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 @param ledger Name of the ledger.
 @param txid Transaction ID.
 @return ApiGetTransactionRequest
*/
func (a *TransactionsApiService) GetTransaction(ctx context.Context, ledger string, txid int64) ApiGetTransactionRequest {
	return ApiGetTransactionRequest{
		ApiService: a,
		ctx: ctx,
		ledger: ledger,
		txid: txid,
	}
}

// Execute executes the request
//  @return GetTransactionResponse
func (a *TransactionsApiService) GetTransactionExecute(r ApiGetTransactionRequest) (*GetTransactionResponse, *http.Response, error) {
	var (
		localVarHTTPMethod   = http.MethodGet
		localVarPostBody     interface{}
		formFiles            []formFile
		localVarReturnValue  *GetTransactionResponse
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "TransactionsApiService.GetTransaction")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/ledger/{ledger}/transactions/{txid}"
	localVarPath = strings.Replace(localVarPath, "{"+"ledger"+"}", url.PathEscape(parameterValueToString(r.ledger, "ledger")), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"txid"+"}", url.PathEscape(parameterValueToString(r.txid, "txid")), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}
	if r.txid < 0 {
		return localVarReturnValue, nil, reportError("txid must be greater than 0")
	}

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
			var v ErrorResponse
			err = a.client.decode(&v, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
			if err != nil {
				newErr.error = err.Error()
				return localVarReturnValue, localVarHTTPResponse, newErr
			}
					newErr.error = formatErrorMessage(localVarHTTPResponse.Status, &v)
					newErr.model = v
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

type ApiListTransactionsRequest struct {
	ctx context.Context
	ApiService TransactionsApi
	ledger string
	pageSize *int64
	reference *string
	account *string
	source *string
	destination *string
	startTime *time.Time
	endTime *time.Time
	cursor *string
	metadata *map[string]string
}

// The maximum number of results to return per page. 
func (r ApiListTransactionsRequest) PageSize(pageSize int64) ApiListTransactionsRequest {
	r.pageSize = &pageSize
	return r
}

// Find transactions by reference field.
func (r ApiListTransactionsRequest) Reference(reference string) ApiListTransactionsRequest {
	r.reference = &reference
	return r
}

// Filter transactions with postings involving given account, either as source or destination (regular expression placed between ^ and $).
func (r ApiListTransactionsRequest) Account(account string) ApiListTransactionsRequest {
	r.account = &account
	return r
}

// Filter transactions with postings involving given account at source (regular expression placed between ^ and $).
func (r ApiListTransactionsRequest) Source(source string) ApiListTransactionsRequest {
	r.source = &source
	return r
}

// Filter transactions with postings involving given account at destination (regular expression placed between ^ and $).
func (r ApiListTransactionsRequest) Destination(destination string) ApiListTransactionsRequest {
	r.destination = &destination
	return r
}

// Filter transactions that occurred after this timestamp. The format is RFC3339 and is inclusive (for example, \&quot;2023-01-02T15:04:01Z\&quot; includes the first second of 4th minute). 
func (r ApiListTransactionsRequest) StartTime(startTime time.Time) ApiListTransactionsRequest {
	r.startTime = &startTime
	return r
}

// Filter transactions that occurred before this timestamp. The format is RFC3339 and is exclusive (for example, \&quot;2023-01-02T15:04:01Z\&quot; excludes the first second of 4th minute). 
func (r ApiListTransactionsRequest) EndTime(endTime time.Time) ApiListTransactionsRequest {
	r.endTime = &endTime
	return r
}

// Parameter used in pagination requests. Maximum page size is set to 15. Set to the value of next for the next page of results. Set to the value of previous for the previous page of results. No other parameters can be set when this parameter is set. 
func (r ApiListTransactionsRequest) Cursor(cursor string) ApiListTransactionsRequest {
	r.cursor = &cursor
	return r
}

// Filter transactions by metadata key value pairs.
func (r ApiListTransactionsRequest) Metadata(metadata map[string]string) ApiListTransactionsRequest {
	r.metadata = &metadata
	return r
}

func (r ApiListTransactionsRequest) Execute() (*TransactionsCursorResponse, *http.Response, error) {
	return r.ApiService.ListTransactionsExecute(r)
}

/*
ListTransactions List transactions from a ledger

List transactions from a ledger, sorted by txid in descending order.

 @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 @param ledger Name of the ledger.
 @return ApiListTransactionsRequest
*/
func (a *TransactionsApiService) ListTransactions(ctx context.Context, ledger string) ApiListTransactionsRequest {
	return ApiListTransactionsRequest{
		ApiService: a,
		ctx: ctx,
		ledger: ledger,
	}
}

// Execute executes the request
//  @return TransactionsCursorResponse
func (a *TransactionsApiService) ListTransactionsExecute(r ApiListTransactionsRequest) (*TransactionsCursorResponse, *http.Response, error) {
	var (
		localVarHTTPMethod   = http.MethodGet
		localVarPostBody     interface{}
		formFiles            []formFile
		localVarReturnValue  *TransactionsCursorResponse
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "TransactionsApiService.ListTransactions")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/ledger/{ledger}/transactions"
	localVarPath = strings.Replace(localVarPath, "{"+"ledger"+"}", url.PathEscape(parameterValueToString(r.ledger, "ledger")), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}

	if r.pageSize != nil {
		parameterAddToHeaderOrQuery(localVarQueryParams, "pageSize", r.pageSize, "")
	}
	if r.reference != nil {
		parameterAddToHeaderOrQuery(localVarQueryParams, "reference", r.reference, "")
	}
	if r.account != nil {
		parameterAddToHeaderOrQuery(localVarQueryParams, "account", r.account, "")
	}
	if r.source != nil {
		parameterAddToHeaderOrQuery(localVarQueryParams, "source", r.source, "")
	}
	if r.destination != nil {
		parameterAddToHeaderOrQuery(localVarQueryParams, "destination", r.destination, "")
	}
	if r.startTime != nil {
		parameterAddToHeaderOrQuery(localVarQueryParams, "startTime", r.startTime, "")
	}
	if r.endTime != nil {
		parameterAddToHeaderOrQuery(localVarQueryParams, "endTime", r.endTime, "")
	}
	if r.cursor != nil {
		parameterAddToHeaderOrQuery(localVarQueryParams, "cursor", r.cursor, "")
	}
	if r.metadata != nil {
		parameterAddToHeaderOrQuery(localVarQueryParams, "metadata", r.metadata, "")
	}
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
			var v ErrorResponse
			err = a.client.decode(&v, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
			if err != nil {
				newErr.error = err.Error()
				return localVarReturnValue, localVarHTTPResponse, newErr
			}
					newErr.error = formatErrorMessage(localVarHTTPResponse.Status, &v)
					newErr.model = v
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

type ApiRevertTransactionRequest struct {
	ctx context.Context
	ApiService TransactionsApi
	ledger string
	txid int64
}

func (r ApiRevertTransactionRequest) Execute() (*CreateTransactionResponse, *http.Response, error) {
	return r.ApiService.RevertTransactionExecute(r)
}

/*
RevertTransaction Revert a ledger transaction by its ID

 @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 @param ledger Name of the ledger.
 @param txid Transaction ID.
 @return ApiRevertTransactionRequest
*/
func (a *TransactionsApiService) RevertTransaction(ctx context.Context, ledger string, txid int64) ApiRevertTransactionRequest {
	return ApiRevertTransactionRequest{
		ApiService: a,
		ctx: ctx,
		ledger: ledger,
		txid: txid,
	}
}

// Execute executes the request
//  @return CreateTransactionResponse
func (a *TransactionsApiService) RevertTransactionExecute(r ApiRevertTransactionRequest) (*CreateTransactionResponse, *http.Response, error) {
	var (
		localVarHTTPMethod   = http.MethodPost
		localVarPostBody     interface{}
		formFiles            []formFile
		localVarReturnValue  *CreateTransactionResponse
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "TransactionsApiService.RevertTransaction")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/ledger/{ledger}/transactions/{txid}/revert"
	localVarPath = strings.Replace(localVarPath, "{"+"ledger"+"}", url.PathEscape(parameterValueToString(r.ledger, "ledger")), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"txid"+"}", url.PathEscape(parameterValueToString(r.txid, "txid")), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}
	if r.txid < 0 {
		return localVarReturnValue, nil, reportError("txid must be greater than 0")
	}

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
			var v ErrorResponse
			err = a.client.decode(&v, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
			if err != nil {
				newErr.error = err.Error()
				return localVarReturnValue, localVarHTTPResponse, newErr
			}
					newErr.error = formatErrorMessage(localVarHTTPResponse.Status, &v)
					newErr.model = v
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
