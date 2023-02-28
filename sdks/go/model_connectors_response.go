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

// checks if the ConnectorsResponse type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &ConnectorsResponse{}

// ConnectorsResponse struct for ConnectorsResponse
type ConnectorsResponse struct {
	Data []ConnectorsResponseDataInner `json:"data"`
}

// NewConnectorsResponse instantiates a new ConnectorsResponse object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewConnectorsResponse(data []ConnectorsResponseDataInner) *ConnectorsResponse {
	this := ConnectorsResponse{}
	this.Data = data
	return &this
}

// NewConnectorsResponseWithDefaults instantiates a new ConnectorsResponse object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewConnectorsResponseWithDefaults() *ConnectorsResponse {
	this := ConnectorsResponse{}
	return &this
}

// GetData returns the Data field value
func (o *ConnectorsResponse) GetData() []ConnectorsResponseDataInner {
	if o == nil {
		var ret []ConnectorsResponseDataInner
		return ret
	}

	return o.Data
}

// GetDataOk returns a tuple with the Data field value
// and a boolean to check if the value has been set.
func (o *ConnectorsResponse) GetDataOk() ([]ConnectorsResponseDataInner, bool) {
	if o == nil {
		return nil, false
	}
	return o.Data, true
}

// SetData sets field value
func (o *ConnectorsResponse) SetData(v []ConnectorsResponseDataInner) {
	o.Data = v
}

func (o ConnectorsResponse) MarshalJSON() ([]byte, error) {
	toSerialize,err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o ConnectorsResponse) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	toSerialize["data"] = o.Data
	return toSerialize, nil
}

type NullableConnectorsResponse struct {
	value *ConnectorsResponse
	isSet bool
}

func (v NullableConnectorsResponse) Get() *ConnectorsResponse {
	return v.value
}

func (v *NullableConnectorsResponse) Set(val *ConnectorsResponse) {
	v.value = val
	v.isSet = true
}

func (v NullableConnectorsResponse) IsSet() bool {
	return v.isSet
}

func (v *NullableConnectorsResponse) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableConnectorsResponse(val *ConnectorsResponse) *NullableConnectorsResponse {
	return &NullableConnectorsResponse{value: val, isSet: true}
}

func (v NullableConnectorsResponse) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableConnectorsResponse) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


