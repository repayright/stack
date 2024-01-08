// Code generated by Speakeasy (https://speakeasyapi.dev). DO NOT EDIT.

package operations

import (
	"github.com/formancehq/formance-sdk-go/pkg/models/sdkerrors"
	"github.com/formancehq/formance-sdk-go/pkg/models/shared"
	"net/http"
)

type CreateTransactionRequest struct {
	// The request body must contain at least one of the following objects:
	//   - `postings`: suitable for simple transactions
	//   - `script`: enabling more complex transactions with Numscript
	//
	PostTransaction shared.PostTransaction `request:"mediaType=application/json"`
	// Name of the ledger.
	Ledger string `pathParam:"style=simple,explode=false,name=ledger"`
	// Set the preview mode. Preview mode doesn't add the logs to the database or publish a message to the message broker.
	Preview *bool `queryParam:"style=form,explode=true,name=preview"`
}

func (o *CreateTransactionRequest) GetPostTransaction() shared.PostTransaction {
	if o == nil {
		return shared.PostTransaction{}
	}
	return o.PostTransaction
}

func (o *CreateTransactionRequest) GetLedger() string {
	if o == nil {
		return ""
	}
	return o.Ledger
}

func (o *CreateTransactionRequest) GetPreview() *bool {
	if o == nil {
		return nil
	}
	return o.Preview
}

type CreateTransactionResponse struct {
	// HTTP response content type for this operation
	ContentType string
	// Error
	ErrorResponse *sdkerrors.ErrorResponse
	// HTTP response status code for this operation
	StatusCode int
	// Raw HTTP response; suitable for custom response parsing
	RawResponse *http.Response
	// OK
	TransactionsResponse *shared.TransactionsResponse
}

func (o *CreateTransactionResponse) GetContentType() string {
	if o == nil {
		return ""
	}
	return o.ContentType
}

func (o *CreateTransactionResponse) GetErrorResponse() *sdkerrors.ErrorResponse {
	if o == nil {
		return nil
	}
	return o.ErrorResponse
}

func (o *CreateTransactionResponse) GetStatusCode() int {
	if o == nil {
		return 0
	}
	return o.StatusCode
}

func (o *CreateTransactionResponse) GetRawResponse() *http.Response {
	if o == nil {
		return nil
	}
	return o.RawResponse
}

func (o *CreateTransactionResponse) GetTransactionsResponse() *shared.TransactionsResponse {
	if o == nil {
		return nil
	}
	return o.TransactionsResponse
}
