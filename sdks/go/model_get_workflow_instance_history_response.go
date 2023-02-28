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

// checks if the GetWorkflowInstanceHistoryResponse type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &GetWorkflowInstanceHistoryResponse{}

// GetWorkflowInstanceHistoryResponse struct for GetWorkflowInstanceHistoryResponse
type GetWorkflowInstanceHistoryResponse struct {
	Data []WorkflowInstanceHistory `json:"data"`
}

// NewGetWorkflowInstanceHistoryResponse instantiates a new GetWorkflowInstanceHistoryResponse object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewGetWorkflowInstanceHistoryResponse(data []WorkflowInstanceHistory) *GetWorkflowInstanceHistoryResponse {
	this := GetWorkflowInstanceHistoryResponse{}
	this.Data = data
	return &this
}

// NewGetWorkflowInstanceHistoryResponseWithDefaults instantiates a new GetWorkflowInstanceHistoryResponse object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewGetWorkflowInstanceHistoryResponseWithDefaults() *GetWorkflowInstanceHistoryResponse {
	this := GetWorkflowInstanceHistoryResponse{}
	return &this
}

// GetData returns the Data field value
func (o *GetWorkflowInstanceHistoryResponse) GetData() []WorkflowInstanceHistory {
	if o == nil {
		var ret []WorkflowInstanceHistory
		return ret
	}

	return o.Data
}

// GetDataOk returns a tuple with the Data field value
// and a boolean to check if the value has been set.
func (o *GetWorkflowInstanceHistoryResponse) GetDataOk() ([]WorkflowInstanceHistory, bool) {
	if o == nil {
		return nil, false
	}
	return o.Data, true
}

// SetData sets field value
func (o *GetWorkflowInstanceHistoryResponse) SetData(v []WorkflowInstanceHistory) {
	o.Data = v
}

func (o GetWorkflowInstanceHistoryResponse) MarshalJSON() ([]byte, error) {
	toSerialize,err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o GetWorkflowInstanceHistoryResponse) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	toSerialize["data"] = o.Data
	return toSerialize, nil
}

type NullableGetWorkflowInstanceHistoryResponse struct {
	value *GetWorkflowInstanceHistoryResponse
	isSet bool
}

func (v NullableGetWorkflowInstanceHistoryResponse) Get() *GetWorkflowInstanceHistoryResponse {
	return v.value
}

func (v *NullableGetWorkflowInstanceHistoryResponse) Set(val *GetWorkflowInstanceHistoryResponse) {
	v.value = val
	v.isSet = true
}

func (v NullableGetWorkflowInstanceHistoryResponse) IsSet() bool {
	return v.isSet
}

func (v *NullableGetWorkflowInstanceHistoryResponse) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableGetWorkflowInstanceHistoryResponse(val *GetWorkflowInstanceHistoryResponse) *NullableGetWorkflowInstanceHistoryResponse {
	return &NullableGetWorkflowInstanceHistoryResponse{value: val, isSet: true}
}

func (v NullableGetWorkflowInstanceHistoryResponse) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableGetWorkflowInstanceHistoryResponse) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


