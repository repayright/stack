//go:build !ignore_autogenerated
// +build !ignore_autogenerated

/*
Copyright 2022.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Code generated by controller-gen. DO NOT EDIT.

package v1beta1

import (
	"encoding/json"

	auth_componentsv1beta1 "github.com/formancehq/operator/apis/auth.components/v1beta1"
	v2 "k8s.io/api/autoscaling/v2"
	v1 "k8s.io/api/core/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Auth) DeepCopyInto(out *Auth) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Auth.
func (in *Auth) DeepCopy() *Auth {
	if in == nil {
		return nil
	}
	out := new(Auth)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *Auth) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AuthClientConfiguration) DeepCopyInto(out *AuthClientConfiguration) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AuthClientConfiguration.
func (in *AuthClientConfiguration) DeepCopy() *AuthClientConfiguration {
	if in == nil {
		return nil
	}
	out := new(AuthClientConfiguration)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AuthConfigSpec) DeepCopyInto(out *AuthConfigSpec) {
	*out = *in
	if in.OAuth2 != nil {
		in, out := &in.OAuth2, &out.OAuth2
		*out = new(OAuth2ConfigSpec)
		(*in).DeepCopyInto(*out)
	}
	if in.HTTPBasic != nil {
		in, out := &in.HTTPBasic, &out.HTTPBasic
		*out = new(HTTPBasicConfigSpec)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AuthConfigSpec.
func (in *AuthConfigSpec) DeepCopy() *AuthConfigSpec {
	if in == nil {
		return nil
	}
	out := new(AuthConfigSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AuthList) DeepCopyInto(out *AuthList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]Auth, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AuthList.
func (in *AuthList) DeepCopy() *AuthList {
	if in == nil {
		return nil
	}
	out := new(AuthList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *AuthList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AuthSpec) DeepCopyInto(out *AuthSpec) {
	*out = *in
	in.Scalable.DeepCopyInto(&out.Scalable)
	in.ImageHolder.DeepCopyInto(&out.ImageHolder)
	in.Postgres.DeepCopyInto(&out.Postgres)
	if in.Ingress != nil {
		in, out := &in.Ingress, &out.Ingress
		*out = new(IngressSpec)
		(*in).DeepCopyInto(*out)
	}
	in.DelegatedOIDCServer.DeepCopyInto(&out.DelegatedOIDCServer)
	if in.Monitoring != nil {
		in, out := &in.Monitoring, &out.Monitoring
		*out = new(MonitoringSpec)
		(*in).DeepCopyInto(*out)
	}
	if in.StaticClients != nil {
		in, out := &in.StaticClients, &out.StaticClients
		*out = make([]auth_componentsv1beta1.StaticClient, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AuthSpec.
func (in *AuthSpec) DeepCopy() *AuthSpec {
	if in == nil {
		return nil
	}
	out := new(AuthSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Batching) DeepCopyInto(out *Batching) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Batching.
func (in *Batching) DeepCopy() *Batching {
	if in == nil {
		return nil
	}
	out := new(Batching)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *CollectorConfig) DeepCopyInto(out *CollectorConfig) {
	*out = *in
	in.KafkaConfig.DeepCopyInto(&out.KafkaConfig)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new CollectorConfig.
func (in *CollectorConfig) DeepCopy() *CollectorConfig {
	if in == nil {
		return nil
	}
	out := new(CollectorConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Condition) DeepCopyInto(out *Condition) {
	*out = *in
	in.LastTransitionTime.DeepCopyInto(&out.LastTransitionTime)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Condition.
func (in *Condition) DeepCopy() *Condition {
	if in == nil {
		return nil
	}
	out := new(Condition)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in Conditions) DeepCopyInto(out *Conditions) {
	{
		in := &in
		*out = make(Conditions, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Conditions.
func (in Conditions) DeepCopy() Conditions {
	if in == nil {
		return nil
	}
	out := new(Conditions)
	in.DeepCopyInto(out)
	return *out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ConfigSource) DeepCopyInto(out *ConfigSource) {
	*out = *in
	if in.ConfigMapKeyRef != nil {
		in, out := &in.ConfigMapKeyRef, &out.ConfigMapKeyRef
		*out = new(v1.ConfigMapKeySelector)
		(*in).DeepCopyInto(*out)
	}
	if in.SecretKeyRef != nil {
		in, out := &in.SecretKeyRef, &out.SecretKeyRef
		*out = new(v1.SecretKeySelector)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ConfigSource.
func (in *ConfigSource) DeepCopy() *ConfigSource {
	if in == nil {
		return nil
	}
	out := new(ConfigSource)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Control) DeepCopyInto(out *Control) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Control.
func (in *Control) DeepCopy() *Control {
	if in == nil {
		return nil
	}
	out := new(Control)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *Control) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ControlList) DeepCopyInto(out *ControlList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]Control, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ControlList.
func (in *ControlList) DeepCopy() *ControlList {
	if in == nil {
		return nil
	}
	out := new(ControlList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *ControlList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ControlSpec) DeepCopyInto(out *ControlSpec) {
	*out = *in
	in.Scalable.DeepCopyInto(&out.Scalable)
	in.ImageHolder.DeepCopyInto(&out.ImageHolder)
	if in.Ingress != nil {
		in, out := &in.Ingress, &out.Ingress
		*out = new(IngressSpec)
		(*in).DeepCopyInto(*out)
	}
	if in.Monitoring != nil {
		in, out := &in.Monitoring, &out.Monitoring
		*out = new(MonitoringSpec)
		(*in).DeepCopyInto(*out)
	}
	if in.AuthClientConfiguration != nil {
		in, out := &in.AuthClientConfiguration, &out.AuthClientConfiguration
		*out = new(AuthClientConfiguration)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ControlSpec.
func (in *ControlSpec) DeepCopy() *ControlSpec {
	if in == nil {
		return nil
	}
	out := new(ControlSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DelegatedOIDCServerConfiguration) DeepCopyInto(out *DelegatedOIDCServerConfiguration) {
	*out = *in
	if in.IssuerFrom != nil {
		in, out := &in.IssuerFrom, &out.IssuerFrom
		*out = new(ConfigSource)
		(*in).DeepCopyInto(*out)
	}
	if in.ClientIDFrom != nil {
		in, out := &in.ClientIDFrom, &out.ClientIDFrom
		*out = new(ConfigSource)
		(*in).DeepCopyInto(*out)
	}
	if in.ClientSecretFrom != nil {
		in, out := &in.ClientSecretFrom, &out.ClientSecretFrom
		*out = new(ConfigSource)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DelegatedOIDCServerConfiguration.
func (in *DelegatedOIDCServerConfiguration) DeepCopy() *DelegatedOIDCServerConfiguration {
	if in == nil {
		return nil
	}
	out := new(DelegatedOIDCServerConfiguration)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ElasticSearchBasicAuthConfig) DeepCopyInto(out *ElasticSearchBasicAuthConfig) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ElasticSearchBasicAuthConfig.
func (in *ElasticSearchBasicAuthConfig) DeepCopy() *ElasticSearchBasicAuthConfig {
	if in == nil {
		return nil
	}
	out := new(ElasticSearchBasicAuthConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ElasticSearchConfig) DeepCopyInto(out *ElasticSearchConfig) {
	*out = *in
	if in.HostFrom != nil {
		in, out := &in.HostFrom, &out.HostFrom
		*out = new(ConfigSource)
		(*in).DeepCopyInto(*out)
	}
	if in.PortFrom != nil {
		in, out := &in.PortFrom, &out.PortFrom
		*out = new(ConfigSource)
		(*in).DeepCopyInto(*out)
	}
	out.TLS = in.TLS
	if in.BasicAuth != nil {
		in, out := &in.BasicAuth, &out.BasicAuth
		*out = new(ElasticSearchBasicAuthConfig)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ElasticSearchConfig.
func (in *ElasticSearchConfig) DeepCopy() *ElasticSearchConfig {
	if in == nil {
		return nil
	}
	out := new(ElasticSearchConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ElasticSearchTLSConfig) DeepCopyInto(out *ElasticSearchTLSConfig) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ElasticSearchTLSConfig.
func (in *ElasticSearchTLSConfig) DeepCopy() *ElasticSearchTLSConfig {
	if in == nil {
		return nil
	}
	out := new(ElasticSearchTLSConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *HTTPBasicConfigSpec) DeepCopyInto(out *HTTPBasicConfigSpec) {
	*out = *in
	if in.Credentials != nil {
		in, out := &in.Credentials, &out.Credentials
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new HTTPBasicConfigSpec.
func (in *HTTPBasicConfigSpec) DeepCopy() *HTTPBasicConfigSpec {
	if in == nil {
		return nil
	}
	out := new(HTTPBasicConfigSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ImageHolder) DeepCopyInto(out *ImageHolder) {
	*out = *in
	if in.ImagePullSecrets != nil {
		in, out := &in.ImagePullSecrets, &out.ImagePullSecrets
		*out = make([]v1.LocalObjectReference, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ImageHolder.
func (in *ImageHolder) DeepCopy() *ImageHolder {
	if in == nil {
		return nil
	}
	out := new(ImageHolder)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *IngressSpec) DeepCopyInto(out *IngressSpec) {
	*out = *in
	if in.Annotations != nil {
		in, out := &in.Annotations, &out.Annotations
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	if in.TLS != nil {
		in, out := &in.TLS, &out.TLS
		*out = new(IngressTLS)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new IngressSpec.
func (in *IngressSpec) DeepCopy() *IngressSpec {
	if in == nil {
		return nil
	}
	out := new(IngressSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *IngressTLS) DeepCopyInto(out *IngressTLS) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new IngressTLS.
func (in *IngressTLS) DeepCopy() *IngressTLS {
	if in == nil {
		return nil
	}
	out := new(IngressTLS)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *KafkaConfig) DeepCopyInto(out *KafkaConfig) {
	*out = *in
	if in.Brokers != nil {
		in, out := &in.Brokers, &out.Brokers
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.BrokersFrom != nil {
		in, out := &in.BrokersFrom, &out.BrokersFrom
		*out = new(ConfigSource)
		(*in).DeepCopyInto(*out)
	}
	if in.SASL != nil {
		in, out := &in.SASL, &out.SASL
		*out = new(KafkaSASLConfig)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new KafkaConfig.
func (in *KafkaConfig) DeepCopy() *KafkaConfig {
	if in == nil {
		return nil
	}
	out := new(KafkaConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *KafkaSASLConfig) DeepCopyInto(out *KafkaSASLConfig) {
	*out = *in
	if in.UsernameFrom != nil {
		in, out := &in.UsernameFrom, &out.UsernameFrom
		*out = new(ConfigSource)
		(*in).DeepCopyInto(*out)
	}
	if in.PasswordFrom != nil {
		in, out := &in.PasswordFrom, &out.PasswordFrom
		*out = new(ConfigSource)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new KafkaSASLConfig.
func (in *KafkaSASLConfig) DeepCopy() *KafkaSASLConfig {
	if in == nil {
		return nil
	}
	out := new(KafkaSASLConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Ledger) DeepCopyInto(out *Ledger) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Ledger.
func (in *Ledger) DeepCopy() *Ledger {
	if in == nil {
		return nil
	}
	out := new(Ledger)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *Ledger) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *LedgerList) DeepCopyInto(out *LedgerList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]Ledger, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new LedgerList.
func (in *LedgerList) DeepCopy() *LedgerList {
	if in == nil {
		return nil
	}
	out := new(LedgerList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *LedgerList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *LedgerSpec) DeepCopyInto(out *LedgerSpec) {
	*out = *in
	in.Scalable.DeepCopyInto(&out.Scalable)
	in.ImageHolder.DeepCopyInto(&out.ImageHolder)
	if in.Ingress != nil {
		in, out := &in.Ingress, &out.Ingress
		*out = new(IngressSpec)
		(*in).DeepCopyInto(*out)
	}
	in.Postgres.DeepCopyInto(&out.Postgres)
	if in.Auth != nil {
		in, out := &in.Auth, &out.Auth
		*out = new(AuthConfigSpec)
		(*in).DeepCopyInto(*out)
	}
	if in.Monitoring != nil {
		in, out := &in.Monitoring, &out.Monitoring
		*out = new(MonitoringSpec)
		(*in).DeepCopyInto(*out)
	}
	if in.Collector != nil {
		in, out := &in.Collector, &out.Collector
		*out = new(CollectorConfig)
		(*in).DeepCopyInto(*out)
	}
	in.LockingStrategy.DeepCopyInto(&out.LockingStrategy)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new LedgerSpec.
func (in *LedgerSpec) DeepCopy() *LedgerSpec {
	if in == nil {
		return nil
	}
	out := new(LedgerSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *LockingStrategy) DeepCopyInto(out *LockingStrategy) {
	*out = *in
	if in.Redis != nil {
		in, out := &in.Redis, &out.Redis
		*out = new(LockingStrategyRedisConfig)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new LockingStrategy.
func (in *LockingStrategy) DeepCopy() *LockingStrategy {
	if in == nil {
		return nil
	}
	out := new(LockingStrategy)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *LockingStrategyRedisConfig) DeepCopyInto(out *LockingStrategyRedisConfig) {
	*out = *in
	if in.UriFrom != nil {
		in, out := &in.UriFrom, &out.UriFrom
		*out = new(ConfigSource)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new LockingStrategyRedisConfig.
func (in *LockingStrategyRedisConfig) DeepCopy() *LockingStrategyRedisConfig {
	if in == nil {
		return nil
	}
	out := new(LockingStrategyRedisConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *MongoDBConfig) DeepCopyInto(out *MongoDBConfig) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new MongoDBConfig.
func (in *MongoDBConfig) DeepCopy() *MongoDBConfig {
	if in == nil {
		return nil
	}
	out := new(MongoDBConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *MonitoringSpec) DeepCopyInto(out *MonitoringSpec) {
	*out = *in
	if in.Traces != nil {
		in, out := &in.Traces, &out.Traces
		*out = new(TracesSpec)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new MonitoringSpec.
func (in *MonitoringSpec) DeepCopy() *MonitoringSpec {
	if in == nil {
		return nil
	}
	out := new(MonitoringSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *NatsConfig) DeepCopyInto(out *NatsConfig) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new NatsConfig.
func (in *NatsConfig) DeepCopy() *NatsConfig {
	if in == nil {
		return nil
	}
	out := new(NatsConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *OAuth2ConfigSpec) DeepCopyInto(out *OAuth2ConfigSpec) {
	*out = *in
	if in.Audiences != nil {
		in, out := &in.Audiences, &out.Audiences
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new OAuth2ConfigSpec.
func (in *OAuth2ConfigSpec) DeepCopy() *OAuth2ConfigSpec {
	if in == nil {
		return nil
	}
	out := new(OAuth2ConfigSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Payments) DeepCopyInto(out *Payments) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Payments.
func (in *Payments) DeepCopy() *Payments {
	if in == nil {
		return nil
	}
	out := new(Payments)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *Payments) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PaymentsList) DeepCopyInto(out *PaymentsList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]Payments, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PaymentsList.
func (in *PaymentsList) DeepCopy() *PaymentsList {
	if in == nil {
		return nil
	}
	out := new(PaymentsList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *PaymentsList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PaymentsSpec) DeepCopyInto(out *PaymentsSpec) {
	*out = *in
	in.ImageHolder.DeepCopyInto(&out.ImageHolder)
	if in.Ingress != nil {
		in, out := &in.Ingress, &out.Ingress
		*out = new(IngressSpec)
		(*in).DeepCopyInto(*out)
	}
	if in.Auth != nil {
		in, out := &in.Auth, &out.Auth
		*out = new(AuthConfigSpec)
		(*in).DeepCopyInto(*out)
	}
	if in.Monitoring != nil {
		in, out := &in.Monitoring, &out.Monitoring
		*out = new(MonitoringSpec)
		(*in).DeepCopyInto(*out)
	}
	if in.Collector != nil {
		in, out := &in.Collector, &out.Collector
		*out = new(CollectorConfig)
		(*in).DeepCopyInto(*out)
	}
	out.MongoDB = in.MongoDB
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PaymentsSpec.
func (in *PaymentsSpec) DeepCopy() *PaymentsSpec {
	if in == nil {
		return nil
	}
	out := new(PaymentsSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PostgresConfig) DeepCopyInto(out *PostgresConfig) {
	*out = *in
	if in.PortFrom != nil {
		in, out := &in.PortFrom, &out.PortFrom
		*out = new(ConfigSource)
		(*in).DeepCopyInto(*out)
	}
	if in.HostFrom != nil {
		in, out := &in.HostFrom, &out.HostFrom
		*out = new(ConfigSource)
		(*in).DeepCopyInto(*out)
	}
	if in.UsernameFrom != nil {
		in, out := &in.UsernameFrom, &out.UsernameFrom
		*out = new(ConfigSource)
		(*in).DeepCopyInto(*out)
	}
	if in.PasswordFrom != nil {
		in, out := &in.PasswordFrom, &out.PasswordFrom
		*out = new(ConfigSource)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PostgresConfig.
func (in *PostgresConfig) DeepCopy() *PostgresConfig {
	if in == nil {
		return nil
	}
	out := new(PostgresConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PostgresConfigCreateDatabase) DeepCopyInto(out *PostgresConfigCreateDatabase) {
	*out = *in
	in.PostgresConfigWithDatabase.DeepCopyInto(&out.PostgresConfigWithDatabase)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PostgresConfigCreateDatabase.
func (in *PostgresConfigCreateDatabase) DeepCopy() *PostgresConfigCreateDatabase {
	if in == nil {
		return nil
	}
	out := new(PostgresConfigCreateDatabase)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PostgresConfigWithDatabase) DeepCopyInto(out *PostgresConfigWithDatabase) {
	*out = *in
	in.PostgresConfig.DeepCopyInto(&out.PostgresConfig)
	if in.DatabaseFrom != nil {
		in, out := &in.DatabaseFrom, &out.DatabaseFrom
		*out = new(ConfigSource)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PostgresConfigWithDatabase.
func (in *PostgresConfigWithDatabase) DeepCopy() *PostgresConfigWithDatabase {
	if in == nil {
		return nil
	}
	out := new(PostgresConfigWithDatabase)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ReplicationStatus) DeepCopyInto(out *ReplicationStatus) {
	*out = *in
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ReplicationStatus.
func (in *ReplicationStatus) DeepCopy() *ReplicationStatus {
	if in == nil {
		return nil
	}
	out := new(ReplicationStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Scalable) DeepCopyInto(out *Scalable) {
	*out = *in
	if in.Replicas != nil {
		in, out := &in.Replicas, &out.Replicas
		*out = new(int32)
		**out = **in
	}
	if in.MinReplicas != nil {
		in, out := &in.MinReplicas, &out.MinReplicas
		*out = new(int32)
		**out = **in
	}
	if in.Metrics != nil {
		in, out := &in.Metrics, &out.Metrics
		*out = make([]v2.MetricSpec, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Scalable.
func (in *Scalable) DeepCopy() *Scalable {
	if in == nil {
		return nil
	}
	out := new(Scalable)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Search) DeepCopyInto(out *Search) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Search.
func (in *Search) DeepCopy() *Search {
	if in == nil {
		return nil
	}
	out := new(Search)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *Search) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *SearchIngester) DeepCopyInto(out *SearchIngester) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new SearchIngester.
func (in *SearchIngester) DeepCopy() *SearchIngester {
	if in == nil {
		return nil
	}
	out := new(SearchIngester)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *SearchIngester) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *SearchIngesterList) DeepCopyInto(out *SearchIngesterList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]SearchIngester, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new SearchIngesterList.
func (in *SearchIngesterList) DeepCopy() *SearchIngesterList {
	if in == nil {
		return nil
	}
	out := new(SearchIngesterList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *SearchIngesterList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *SearchIngesterSpec) DeepCopyInto(out *SearchIngesterSpec) {
	*out = *in
	if in.Pipeline != nil {
		in, out := &in.Pipeline, &out.Pipeline
		*out = make(json.RawMessage, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new SearchIngesterSpec.
func (in *SearchIngesterSpec) DeepCopy() *SearchIngesterSpec {
	if in == nil {
		return nil
	}
	out := new(SearchIngesterSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *SearchList) DeepCopyInto(out *SearchList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]Search, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new SearchList.
func (in *SearchList) DeepCopy() *SearchList {
	if in == nil {
		return nil
	}
	out := new(SearchList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *SearchList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *SearchSpec) DeepCopyInto(out *SearchSpec) {
	*out = *in
	in.Scalable.DeepCopyInto(&out.Scalable)
	in.ImageHolder.DeepCopyInto(&out.ImageHolder)
	if in.Ingress != nil {
		in, out := &in.Ingress, &out.Ingress
		*out = new(IngressSpec)
		(*in).DeepCopyInto(*out)
	}
	if in.Auth != nil {
		in, out := &in.Auth, &out.Auth
		*out = new(AuthConfigSpec)
		(*in).DeepCopyInto(*out)
	}
	if in.Monitoring != nil {
		in, out := &in.Monitoring, &out.Monitoring
		*out = new(MonitoringSpec)
		(*in).DeepCopyInto(*out)
	}
	in.ElasticSearch.DeepCopyInto(&out.ElasticSearch)
	in.KafkaConfig.DeepCopyInto(&out.KafkaConfig)
	out.Batching = in.Batching
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new SearchSpec.
func (in *SearchSpec) DeepCopy() *SearchSpec {
	if in == nil {
		return nil
	}
	out := new(SearchSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Status) DeepCopyInto(out *Status) {
	*out = *in
	if in.Conditions != nil {
		in, out := &in.Conditions, &out.Conditions
		*out = make(Conditions, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Status.
func (in *Status) DeepCopy() *Status {
	if in == nil {
		return nil
	}
	out := new(Status)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *TracesOtlpSpec) DeepCopyInto(out *TracesOtlpSpec) {
	*out = *in
	if in.EndpointFrom != nil {
		in, out := &in.EndpointFrom, &out.EndpointFrom
		*out = new(v1.EnvVarSource)
		(*in).DeepCopyInto(*out)
	}
	if in.PortFrom != nil {
		in, out := &in.PortFrom, &out.PortFrom
		*out = new(v1.EnvVarSource)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new TracesOtlpSpec.
func (in *TracesOtlpSpec) DeepCopy() *TracesOtlpSpec {
	if in == nil {
		return nil
	}
	out := new(TracesOtlpSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *TracesSpec) DeepCopyInto(out *TracesSpec) {
	*out = *in
	if in.Otlp != nil {
		in, out := &in.Otlp, &out.Otlp
		*out = new(TracesOtlpSpec)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new TracesSpec.
func (in *TracesSpec) DeepCopy() *TracesSpec {
	if in == nil {
		return nil
	}
	out := new(TracesSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Webhooks) DeepCopyInto(out *Webhooks) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Webhooks.
func (in *Webhooks) DeepCopy() *Webhooks {
	if in == nil {
		return nil
	}
	out := new(Webhooks)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *Webhooks) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *WebhooksList) DeepCopyInto(out *WebhooksList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]Webhooks, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new WebhooksList.
func (in *WebhooksList) DeepCopy() *WebhooksList {
	if in == nil {
		return nil
	}
	out := new(WebhooksList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *WebhooksList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *WebhooksSpec) DeepCopyInto(out *WebhooksSpec) {
	*out = *in
	in.ImageHolder.DeepCopyInto(&out.ImageHolder)
	if in.Ingress != nil {
		in, out := &in.Ingress, &out.Ingress
		*out = new(IngressSpec)
		(*in).DeepCopyInto(*out)
	}
	if in.Auth != nil {
		in, out := &in.Auth, &out.Auth
		*out = new(AuthConfigSpec)
		(*in).DeepCopyInto(*out)
	}
	if in.Monitoring != nil {
		in, out := &in.Monitoring, &out.Monitoring
		*out = new(MonitoringSpec)
		(*in).DeepCopyInto(*out)
	}
	if in.Collector != nil {
		in, out := &in.Collector, &out.Collector
		*out = new(CollectorConfig)
		(*in).DeepCopyInto(*out)
	}
	out.MongoDB = in.MongoDB
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new WebhooksSpec.
func (in *WebhooksSpec) DeepCopy() *WebhooksSpec {
	if in == nil {
		return nil
	}
	out := new(WebhooksSpec)
	in.DeepCopyInto(out)
	return out
}
