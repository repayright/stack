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

// WebhooksCursor struct for WebhooksCursor
type WebhooksCursor struct {
	HasMore bool             `json:"hasMore"`
	Data    []WebhooksConfig `json:"data"`
}

// NewWebhooksCursor instantiates a new WebhooksCursor object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewWebhooksCursor(hasMore bool, data []WebhooksConfig) *WebhooksCursor {
	this := WebhooksCursor{}
	this.HasMore = hasMore
	this.Data = data
	return &this
}

// NewWebhooksCursorWithDefaults instantiates a new WebhooksCursor object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewWebhooksCursorWithDefaults() *WebhooksCursor {
	this := WebhooksCursor{}
	return &this
}

// GetHasMore returns the HasMore field value
func (o *WebhooksCursor) GetHasMore() bool {
	if o == nil {
		var ret bool
		return ret
	}

	return o.HasMore
}

// GetHasMoreOk returns a tuple with the HasMore field value
// and a boolean to check if the value has been set.
func (o *WebhooksCursor) GetHasMoreOk() (*bool, bool) {
	if o == nil {
		return nil, false
	}
	return &o.HasMore, true
}

// SetHasMore sets field value
func (o *WebhooksCursor) SetHasMore(v bool) {
	o.HasMore = v
}

// GetData returns the Data field value
func (o *WebhooksCursor) GetData() []WebhooksConfig {
	if o == nil {
		var ret []WebhooksConfig
		return ret
	}

	return o.Data
}

// GetDataOk returns a tuple with the Data field value
// and a boolean to check if the value has been set.
func (o *WebhooksCursor) GetDataOk() ([]WebhooksConfig, bool) {
	if o == nil {
		return nil, false
	}
	return o.Data, true
}

// SetData sets field value
func (o *WebhooksCursor) SetData(v []WebhooksConfig) {
	o.Data = v
}

func (o WebhooksCursor) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}
	if true {
		toSerialize["hasMore"] = o.HasMore
	}
	if true {
		toSerialize["data"] = o.Data
	}
	return json.Marshal(toSerialize)
}

type NullableWebhooksCursor struct {
	value *WebhooksCursor
	isSet bool
}

func (v NullableWebhooksCursor) Get() *WebhooksCursor {
	return v.value
}

func (v *NullableWebhooksCursor) Set(val *WebhooksCursor) {
	v.value = val
	v.isSet = true
}

func (v NullableWebhooksCursor) IsSet() bool {
	return v.isSet
}

func (v *NullableWebhooksCursor) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableWebhooksCursor(val *WebhooksCursor) *NullableWebhooksCursor {
	return &NullableWebhooksCursor{value: val, isSet: true}
}

func (v NullableWebhooksCursor) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableWebhooksCursor) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}