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

package v1alpha2

import (
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DataLayerConfig) DeepCopyInto(out *DataLayerConfig) {
	*out = *in
	if in.Native != nil {
		in, out := &in.Native, &out.Native
		*out = new(NativeConfig)
		**out = **in
	}
	if in.LargeScale != nil {
		in, out := &in.LargeScale, &out.LargeScale
		*out = new(LargeScaleConfig)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DataLayerConfig.
func (in *DataLayerConfig) DeepCopy() *DataLayerConfig {
	if in == nil {
		return nil
	}
	out := new(DataLayerConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *KafkaConfig) DeepCopyInto(out *KafkaConfig) {
	*out = *in
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
func (in *LargeScaleConfig) DeepCopyInto(out *LargeScaleConfig) {
	*out = *in
	if in.Kafka != nil {
		in, out := &in.Kafka, &out.Kafka
		*out = new(KafkaConfig)
		**out = **in
	}
	out.Postgres = in.Postgres
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new LargeScaleConfig.
func (in *LargeScaleConfig) DeepCopy() *LargeScaleConfig {
	if in == nil {
		return nil
	}
	out := new(LargeScaleConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *MulticlusterGlobalHub) DeepCopyInto(out *MulticlusterGlobalHub) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new MulticlusterGlobalHub.
func (in *MulticlusterGlobalHub) DeepCopy() *MulticlusterGlobalHub {
	if in == nil {
		return nil
	}
	out := new(MulticlusterGlobalHub)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *MulticlusterGlobalHub) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *MulticlusterGlobalHubList) DeepCopyInto(out *MulticlusterGlobalHubList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]MulticlusterGlobalHub, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new MulticlusterGlobalHubList.
func (in *MulticlusterGlobalHubList) DeepCopy() *MulticlusterGlobalHubList {
	if in == nil {
		return nil
	}
	out := new(MulticlusterGlobalHubList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *MulticlusterGlobalHubList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *MulticlusterGlobalHubSpec) DeepCopyInto(out *MulticlusterGlobalHubSpec) {
	*out = *in
	if in.NodeSelector != nil {
		in, out := &in.NodeSelector, &out.NodeSelector
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	if in.Tolerations != nil {
		in, out := &in.Tolerations, &out.Tolerations
		*out = make([]v1.Toleration, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.DataLayer != nil {
		in, out := &in.DataLayer, &out.DataLayer
		*out = new(DataLayerConfig)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new MulticlusterGlobalHubSpec.
func (in *MulticlusterGlobalHubSpec) DeepCopy() *MulticlusterGlobalHubSpec {
	if in == nil {
		return nil
	}
	out := new(MulticlusterGlobalHubSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *MulticlusterGlobalHubStatus) DeepCopyInto(out *MulticlusterGlobalHubStatus) {
	*out = *in
	if in.Conditions != nil {
		in, out := &in.Conditions, &out.Conditions
		*out = make([]metav1.Condition, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new MulticlusterGlobalHubStatus.
func (in *MulticlusterGlobalHubStatus) DeepCopy() *MulticlusterGlobalHubStatus {
	if in == nil {
		return nil
	}
	out := new(MulticlusterGlobalHubStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *NativeConfig) DeepCopyInto(out *NativeConfig) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new NativeConfig.
func (in *NativeConfig) DeepCopy() *NativeConfig {
	if in == nil {
		return nil
	}
	out := new(NativeConfig)
	in.DeepCopyInto(out)
	return out
}
