// Code generated by Speakeasy (https://speakeasyapi.dev). DO NOT EDIT.

package operations

import (
	"net/http"
)

type V2ImportLogsRequest struct {
	RequestBody *string `request:"mediaType=application/octet-stream"`
	// Name of the ledger.
	Ledger string `pathParam:"style=simple,explode=false,name=ledger"`
}

func (o *V2ImportLogsRequest) GetRequestBody() *string {
	if o == nil {
		return nil
	}
	return o.RequestBody
}

func (o *V2ImportLogsRequest) GetLedger() string {
	if o == nil {
		return ""
	}
	return o.Ledger
}

type V2ImportLogsResponse struct {
	// HTTP response content type for this operation
	ContentType string
	// HTTP response status code for this operation
	StatusCode int
	// Raw HTTP response; suitable for custom response parsing
	RawResponse *http.Response
}

func (o *V2ImportLogsResponse) GetContentType() string {
	if o == nil {
		return ""
	}
	return o.ContentType
}

func (o *V2ImportLogsResponse) GetStatusCode() int {
	if o == nil {
		return 0
	}
	return o.StatusCode
}

func (o *V2ImportLogsResponse) GetRawResponse() *http.Response {
	if o == nil {
		return nil
	}
	return o.RawResponse
}
