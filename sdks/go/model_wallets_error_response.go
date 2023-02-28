/*
Formance Stack API

Open, modular foundation for unique payments flows  # Introduction This API is documented in **OpenAPI format**.  # Authentication Formance Stack offers one forms of authentication:   - OAuth2 OAuth2 - an open protocol to allow secure authorization in a simple and standard method from web, mobile and desktop applications. <SecurityDefinitions /> 

API version: v1.0.20230228
Contact: support@formance.com
*/

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package formance

import (
	"encoding/json"
)

// checks if the WalletsErrorResponse type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &WalletsErrorResponse{}

// WalletsErrorResponse struct for WalletsErrorResponse
type WalletsErrorResponse struct {
	ErrorCode string `json:"errorCode"`
	ErrorMessage string `json:"errorMessage"`
}

// NewWalletsErrorResponse instantiates a new WalletsErrorResponse object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewWalletsErrorResponse(errorCode string, errorMessage string) *WalletsErrorResponse {
	this := WalletsErrorResponse{}
	this.ErrorCode = errorCode
	this.ErrorMessage = errorMessage
	return &this
}

// NewWalletsErrorResponseWithDefaults instantiates a new WalletsErrorResponse object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewWalletsErrorResponseWithDefaults() *WalletsErrorResponse {
	this := WalletsErrorResponse{}
	return &this
}

// GetErrorCode returns the ErrorCode field value
func (o *WalletsErrorResponse) GetErrorCode() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.ErrorCode
}

// GetErrorCodeOk returns a tuple with the ErrorCode field value
// and a boolean to check if the value has been set.
func (o *WalletsErrorResponse) GetErrorCodeOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.ErrorCode, true
}

// SetErrorCode sets field value
func (o *WalletsErrorResponse) SetErrorCode(v string) {
	o.ErrorCode = v
}

// GetErrorMessage returns the ErrorMessage field value
func (o *WalletsErrorResponse) GetErrorMessage() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.ErrorMessage
}

// GetErrorMessageOk returns a tuple with the ErrorMessage field value
// and a boolean to check if the value has been set.
func (o *WalletsErrorResponse) GetErrorMessageOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.ErrorMessage, true
}

// SetErrorMessage sets field value
func (o *WalletsErrorResponse) SetErrorMessage(v string) {
	o.ErrorMessage = v
}

func (o WalletsErrorResponse) MarshalJSON() ([]byte, error) {
	toSerialize,err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o WalletsErrorResponse) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	toSerialize["errorCode"] = o.ErrorCode
	toSerialize["errorMessage"] = o.ErrorMessage
	return toSerialize, nil
}

type NullableWalletsErrorResponse struct {
	value *WalletsErrorResponse
	isSet bool
}

func (v NullableWalletsErrorResponse) Get() *WalletsErrorResponse {
	return v.value
}

func (v *NullableWalletsErrorResponse) Set(val *WalletsErrorResponse) {
	v.value = val
	v.isSet = true
}

func (v NullableWalletsErrorResponse) IsSet() bool {
	return v.isSet
}

func (v *NullableWalletsErrorResponse) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableWalletsErrorResponse(val *WalletsErrorResponse) *NullableWalletsErrorResponse {
	return &NullableWalletsErrorResponse{value: val, isSet: true}
}

func (v NullableWalletsErrorResponse) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableWalletsErrorResponse) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


