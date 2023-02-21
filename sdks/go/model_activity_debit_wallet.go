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

// checks if the ActivityDebitWallet type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &ActivityDebitWallet{}

// ActivityDebitWallet struct for ActivityDebitWallet
type ActivityDebitWallet struct {
	Id *string `json:"id,omitempty"`
	Data *DebitWalletRequest `json:"data,omitempty"`
}

// NewActivityDebitWallet instantiates a new ActivityDebitWallet object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewActivityDebitWallet() *ActivityDebitWallet {
	this := ActivityDebitWallet{}
	return &this
}

// NewActivityDebitWalletWithDefaults instantiates a new ActivityDebitWallet object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewActivityDebitWalletWithDefaults() *ActivityDebitWallet {
	this := ActivityDebitWallet{}
	return &this
}

// GetId returns the Id field value if set, zero value otherwise.
func (o *ActivityDebitWallet) GetId() string {
	if o == nil || IsNil(o.Id) {
		var ret string
		return ret
	}
	return *o.Id
}

// GetIdOk returns a tuple with the Id field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ActivityDebitWallet) GetIdOk() (*string, bool) {
	if o == nil || IsNil(o.Id) {
		return nil, false
	}
	return o.Id, true
}

// HasId returns a boolean if a field has been set.
func (o *ActivityDebitWallet) HasId() bool {
	if o != nil && !IsNil(o.Id) {
		return true
	}

	return false
}

// SetId gets a reference to the given string and assigns it to the Id field.
func (o *ActivityDebitWallet) SetId(v string) {
	o.Id = &v
}

// GetData returns the Data field value if set, zero value otherwise.
func (o *ActivityDebitWallet) GetData() DebitWalletRequest {
	if o == nil || IsNil(o.Data) {
		var ret DebitWalletRequest
		return ret
	}
	return *o.Data
}

// GetDataOk returns a tuple with the Data field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ActivityDebitWallet) GetDataOk() (*DebitWalletRequest, bool) {
	if o == nil || IsNil(o.Data) {
		return nil, false
	}
	return o.Data, true
}

// HasData returns a boolean if a field has been set.
func (o *ActivityDebitWallet) HasData() bool {
	if o != nil && !IsNil(o.Data) {
		return true
	}

	return false
}

// SetData gets a reference to the given DebitWalletRequest and assigns it to the Data field.
func (o *ActivityDebitWallet) SetData(v DebitWalletRequest) {
	o.Data = &v
}

func (o ActivityDebitWallet) MarshalJSON() ([]byte, error) {
	toSerialize,err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o ActivityDebitWallet) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !IsNil(o.Id) {
		toSerialize["id"] = o.Id
	}
	if !IsNil(o.Data) {
		toSerialize["data"] = o.Data
	}
	return toSerialize, nil
}

type NullableActivityDebitWallet struct {
	value *ActivityDebitWallet
	isSet bool
}

func (v NullableActivityDebitWallet) Get() *ActivityDebitWallet {
	return v.value
}

func (v *NullableActivityDebitWallet) Set(val *ActivityDebitWallet) {
	v.value = val
	v.isSet = true
}

func (v NullableActivityDebitWallet) IsSet() bool {
	return v.isSet
}

func (v *NullableActivityDebitWallet) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableActivityDebitWallet(val *ActivityDebitWallet) *NullableActivityDebitWallet {
	return &NullableActivityDebitWallet{value: val, isSet: true}
}

func (v NullableActivityDebitWallet) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableActivityDebitWallet) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


