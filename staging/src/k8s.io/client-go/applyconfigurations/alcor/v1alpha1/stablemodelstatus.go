/*
Copyright The Kubernetes Authors.

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

// Code generated by applyconfiguration-gen. DO NOT EDIT.

package v1alpha1

import (
	v1alpha1 "k8s.io/api/alcor/v1alpha1"
)

// StableModelStatusApplyConfiguration represents an declarative configuration of the StableModelStatus type for use
// with apply.
type StableModelStatusApplyConfiguration struct {
	Phase                  *v1alpha1.Phase                          `json:"phase,omitempty"`
	PodRef                 *PodStatusApplyConfiguration             `json:"podRef,omitempty"`
	IPClaimRef             *IPStatusApplyConfiguration              `json:"ipRef,omitempty"`
	PVCRefs                []PVCStatusApplyConfiguration            `json:"pvcRefs,omitempty"`
	Events                 []EventRecordApplyConfiguration          `json:"events,omitempty"`
	Conditions             []StableModelConditionApplyConfiguration `json:"conditions,omitempty"`
	ResizePhase            *v1alpha1.ResizePhase                    `json:"resizePhase,omitempty"`
	RebootPodTimestamp     *string                                  `json:"rebootPodTimestamp,omitempty"`
	ForceDriftPodTimestamp *string                                  `json:"forceDriftPodTimestamp,omitempty"`
	DriftPod               *bool                                    `json:"driftPod,omitempty"`
	NodeName               *string                                  `json:"nodeName,omitempty"`
}

// StableModelStatusApplyConfiguration constructs an declarative configuration of the StableModelStatus type for use with
// apply.
func StableModelStatus() *StableModelStatusApplyConfiguration {
	return &StableModelStatusApplyConfiguration{}
}

// WithPhase sets the Phase field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Phase field is set to the value of the last call.
func (b *StableModelStatusApplyConfiguration) WithPhase(value v1alpha1.Phase) *StableModelStatusApplyConfiguration {
	b.Phase = &value
	return b
}

// WithPodRef sets the PodRef field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the PodRef field is set to the value of the last call.
func (b *StableModelStatusApplyConfiguration) WithPodRef(value *PodStatusApplyConfiguration) *StableModelStatusApplyConfiguration {
	b.PodRef = value
	return b
}

// WithIPClaimRef sets the IPClaimRef field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the IPClaimRef field is set to the value of the last call.
func (b *StableModelStatusApplyConfiguration) WithIPClaimRef(value *IPStatusApplyConfiguration) *StableModelStatusApplyConfiguration {
	b.IPClaimRef = value
	return b
}

// WithPVCRefs adds the given value to the PVCRefs field in the declarative configuration
// and returns the receiver, so that objects can be build by chaining "With" function invocations.
// If called multiple times, values provided by each call will be appended to the PVCRefs field.
func (b *StableModelStatusApplyConfiguration) WithPVCRefs(values ...*PVCStatusApplyConfiguration) *StableModelStatusApplyConfiguration {
	for i := range values {
		if values[i] == nil {
			panic("nil value passed to WithPVCRefs")
		}
		b.PVCRefs = append(b.PVCRefs, *values[i])
	}
	return b
}

// WithEvents adds the given value to the Events field in the declarative configuration
// and returns the receiver, so that objects can be build by chaining "With" function invocations.
// If called multiple times, values provided by each call will be appended to the Events field.
func (b *StableModelStatusApplyConfiguration) WithEvents(values ...*EventRecordApplyConfiguration) *StableModelStatusApplyConfiguration {
	for i := range values {
		if values[i] == nil {
			panic("nil value passed to WithEvents")
		}
		b.Events = append(b.Events, *values[i])
	}
	return b
}

// WithConditions adds the given value to the Conditions field in the declarative configuration
// and returns the receiver, so that objects can be build by chaining "With" function invocations.
// If called multiple times, values provided by each call will be appended to the Conditions field.
func (b *StableModelStatusApplyConfiguration) WithConditions(values ...*StableModelConditionApplyConfiguration) *StableModelStatusApplyConfiguration {
	for i := range values {
		if values[i] == nil {
			panic("nil value passed to WithConditions")
		}
		b.Conditions = append(b.Conditions, *values[i])
	}
	return b
}

// WithResizePhase sets the ResizePhase field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the ResizePhase field is set to the value of the last call.
func (b *StableModelStatusApplyConfiguration) WithResizePhase(value v1alpha1.ResizePhase) *StableModelStatusApplyConfiguration {
	b.ResizePhase = &value
	return b
}

// WithRebootPodTimestamp sets the RebootPodTimestamp field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the RebootPodTimestamp field is set to the value of the last call.
func (b *StableModelStatusApplyConfiguration) WithRebootPodTimestamp(value string) *StableModelStatusApplyConfiguration {
	b.RebootPodTimestamp = &value
	return b
}

// WithForceDriftPodTimestamp sets the ForceDriftPodTimestamp field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the ForceDriftPodTimestamp field is set to the value of the last call.
func (b *StableModelStatusApplyConfiguration) WithForceDriftPodTimestamp(value string) *StableModelStatusApplyConfiguration {
	b.ForceDriftPodTimestamp = &value
	return b
}

// WithDriftPod sets the DriftPod field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the DriftPod field is set to the value of the last call.
func (b *StableModelStatusApplyConfiguration) WithDriftPod(value bool) *StableModelStatusApplyConfiguration {
	b.DriftPod = &value
	return b
}

// WithNodeName sets the NodeName field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the NodeName field is set to the value of the last call.
func (b *StableModelStatusApplyConfiguration) WithNodeName(value string) *StableModelStatusApplyConfiguration {
	b.NodeName = &value
	return b
}
