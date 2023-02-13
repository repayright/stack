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

// checks if the GetHoldsResponseCursorAllOf type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &GetHoldsResponseCursorAllOf{}

// GetHoldsResponseCursorAllOf struct for GetHoldsResponseCursorAllOf
type GetHoldsResponseCursorAllOf struct {
	Data []Hold `json:"data"`
}

// NewGetHoldsResponseCursorAllOf instantiates a new GetHoldsResponseCursorAllOf object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewGetHoldsResponseCursorAllOf(data []Hold) *GetHoldsResponseCursorAllOf {
	this := GetHoldsResponseCursorAllOf{}
	this.Data = data
	return &this
}

// NewGetHoldsResponseCursorAllOfWithDefaults instantiates a new GetHoldsResponseCursorAllOf object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewGetHoldsResponseCursorAllOfWithDefaults() *GetHoldsResponseCursorAllOf {
	this := GetHoldsResponseCursorAllOf{}
	return &this
}

// GetData returns the Data field value
func (o *GetHoldsResponseCursorAllOf) GetData() []Hold {
	if o == nil {
		var ret []Hold
		return ret
	}

	return o.Data
}

// GetDataOk returns a tuple with the Data field value
// and a boolean to check if the value has been set.
func (o *GetHoldsResponseCursorAllOf) GetDataOk() ([]Hold, bool) {
	if o == nil {
		return nil, false
	}
	return o.Data, true
}

// SetData sets field value
func (o *GetHoldsResponseCursorAllOf) SetData(v []Hold) {
	o.Data = v
}

func (o GetHoldsResponseCursorAllOf) MarshalJSON() ([]byte, error) {
	toSerialize,err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o GetHoldsResponseCursorAllOf) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	toSerialize["data"] = o.Data
	return toSerialize, nil
}

type NullableGetHoldsResponseCursorAllOf struct {
	value *GetHoldsResponseCursorAllOf
	isSet bool
}

func (v NullableGetHoldsResponseCursorAllOf) Get() *GetHoldsResponseCursorAllOf {
	return v.value
}

func (v *NullableGetHoldsResponseCursorAllOf) Set(val *GetHoldsResponseCursorAllOf) {
	v.value = val
	v.isSet = true
}

func (v NullableGetHoldsResponseCursorAllOf) IsSet() bool {
	return v.isSet
}

func (v *NullableGetHoldsResponseCursorAllOf) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableGetHoldsResponseCursorAllOf(val *GetHoldsResponseCursorAllOf) *NullableGetHoldsResponseCursorAllOf {
	return &NullableGetHoldsResponseCursorAllOf{value: val, isSet: true}
}

func (v NullableGetHoldsResponseCursorAllOf) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableGetHoldsResponseCursorAllOf) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}

