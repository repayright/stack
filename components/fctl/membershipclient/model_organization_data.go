/*
Membership API

No description provided (generated by Openapi Generator https://github.com/openapitools/openapi-generator)

API version: 0.1.0
*/

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package membershipclient

import (
	"encoding/json"
)

// checks if the OrganizationData type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &OrganizationData{}

// OrganizationData struct for OrganizationData
type OrganizationData struct {
	// Organization name
	Name string `json:"name"`
}

// NewOrganizationData instantiates a new OrganizationData object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewOrganizationData(name string) *OrganizationData {
	this := OrganizationData{}
	this.Name = name
	return &this
}

// NewOrganizationDataWithDefaults instantiates a new OrganizationData object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewOrganizationDataWithDefaults() *OrganizationData {
	this := OrganizationData{}
	return &this
}

// GetName returns the Name field value
func (o *OrganizationData) GetName() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Name
}

// GetNameOk returns a tuple with the Name field value
// and a boolean to check if the value has been set.
func (o *OrganizationData) GetNameOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Name, true
}

// SetName sets field value
func (o *OrganizationData) SetName(v string) {
	o.Name = v
}

func (o OrganizationData) MarshalJSON() ([]byte, error) {
	toSerialize,err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o OrganizationData) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	toSerialize["name"] = o.Name
	return toSerialize, nil
}

type NullableOrganizationData struct {
	value *OrganizationData
	isSet bool
}

func (v NullableOrganizationData) Get() *OrganizationData {
	return v.value
}

func (v *NullableOrganizationData) Set(val *OrganizationData) {
	v.value = val
	v.isSet = true
}

func (v NullableOrganizationData) IsSet() bool {
	return v.isSet
}

func (v *NullableOrganizationData) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableOrganizationData(val *OrganizationData) *NullableOrganizationData {
	return &NullableOrganizationData{value: val, isSet: true}
}

func (v NullableOrganizationData) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableOrganizationData) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
