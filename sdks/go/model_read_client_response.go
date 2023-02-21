/*
Formance Stack API

Open, modular foundation for unique payments flows  # Introduction This API is documented in **OpenAPI format**.  # Authentication Formance Stack offers one forms of authentication:   - OAuth2 OAuth2 - an open protocol to allow secure authorization in a simple and standard method from web, mobile and desktop applications. <SecurityDefinitions /> 

API version: develop
Contact: support@formance.com
*/

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package formance

import (
	"encoding/json"
)

// checks if the ReadClientResponse type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &ReadClientResponse{}

// ReadClientResponse struct for ReadClientResponse
type ReadClientResponse struct {
	Data *Client `json:"data,omitempty"`
}

// NewReadClientResponse instantiates a new ReadClientResponse object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewReadClientResponse() *ReadClientResponse {
	this := ReadClientResponse{}
	return &this
}

// NewReadClientResponseWithDefaults instantiates a new ReadClientResponse object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewReadClientResponseWithDefaults() *ReadClientResponse {
	this := ReadClientResponse{}
	return &this
}

// GetData returns the Data field value if set, zero value otherwise.
func (o *ReadClientResponse) GetData() Client {
	if o == nil || IsNil(o.Data) {
		var ret Client
		return ret
	}
	return *o.Data
}

// GetDataOk returns a tuple with the Data field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ReadClientResponse) GetDataOk() (*Client, bool) {
	if o == nil || IsNil(o.Data) {
		return nil, false
	}
	return o.Data, true
}

// HasData returns a boolean if a field has been set.
func (o *ReadClientResponse) HasData() bool {
	if o != nil && !IsNil(o.Data) {
		return true
	}

	return false
}

// SetData gets a reference to the given Client and assigns it to the Data field.
func (o *ReadClientResponse) SetData(v Client) {
	o.Data = &v
}

func (o ReadClientResponse) MarshalJSON() ([]byte, error) {
	toSerialize,err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o ReadClientResponse) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !IsNil(o.Data) {
		toSerialize["data"] = o.Data
	}
	return toSerialize, nil
}

type NullableReadClientResponse struct {
	value *ReadClientResponse
	isSet bool
}

func (v NullableReadClientResponse) Get() *ReadClientResponse {
	return v.value
}

func (v *NullableReadClientResponse) Set(val *ReadClientResponse) {
	v.value = val
	v.isSet = true
}

func (v NullableReadClientResponse) IsSet() bool {
	return v.isSet
}

func (v *NullableReadClientResponse) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableReadClientResponse(val *ReadClientResponse) *NullableReadClientResponse {
	return &NullableReadClientResponse{value: val, isSet: true}
}

func (v NullableReadClientResponse) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableReadClientResponse) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


