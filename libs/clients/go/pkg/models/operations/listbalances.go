// Code generated by Speakeasy (https://speakeasyapi.dev). DO NOT EDIT.

package operations

import (
	"github.com/formancehq/formance-sdk-go/pkg/models/shared"
	"net/http"
)

type ListBalancesRequest struct {
	ID string `pathParam:"style=simple,explode=false,name=id"`
}

func (o *ListBalancesRequest) GetID() string {
	if o == nil {
		return ""
	}
	return o.ID
}

type ListBalancesResponse struct {
	// HTTP response content type for this operation
	ContentType string
	// Balances list
	ListBalancesResponse *shared.ListBalancesResponse
	// HTTP response status code for this operation
	StatusCode int
	// Raw HTTP response; suitable for custom response parsing
	RawResponse *http.Response
}

func (o *ListBalancesResponse) GetContentType() string {
	if o == nil {
		return ""
	}
	return o.ContentType
}

func (o *ListBalancesResponse) GetListBalancesResponse() *shared.ListBalancesResponse {
	if o == nil {
		return nil
	}
	return o.ListBalancesResponse
}

func (o *ListBalancesResponse) GetStatusCode() int {
	if o == nil {
		return 0
	}
	return o.StatusCode
}

func (o *ListBalancesResponse) GetRawResponse() *http.Response {
	if o == nil {
		return nil
	}
	return o.RawResponse
}