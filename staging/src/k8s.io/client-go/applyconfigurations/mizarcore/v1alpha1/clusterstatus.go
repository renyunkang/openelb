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
	v1 "k8s.io/api/core/v1"
	mizarcorev1alpha1 "k8s.io/api/mizarcore/v1alpha1"
)

// ClusterStatusApplyConfiguration represents an declarative configuration of the ClusterStatus type for use
// with apply.
type ClusterStatusApplyConfiguration struct {
	Configuration            *KubeletForClusterConfigurationApplyConfiguration `json:"configuration,omitempty"`
	Addresses                []AddressApplyConfiguration                       `json:"addresses,omitempty"`
	Phase                    *mizarcorev1alpha1.ClusterPhase                   `json:"phase,omitempty"`
	Allocatable              *v1.ResourceList                                  `json:"allocatable,omitempty"`
	Usage                    *v1.ResourceList                                  `json:"usage,omitempty"`
	Nodes                    []NodeLeftResourceApplyConfiguration              `json:"nodes,omitempty"`
	PlanetClusterAllocatable *v1.ResourceList                                  `json:"planetClusterAllocatable,omitempty"`
	PlanetClusterUsage       *v1.ResourceList                                  `json:"planetClusterUsage,omitempty"`
	PlanetClusterNodes       []NodeLeftResourceApplyConfiguration              `json:"planetClusterNodes,omitempty"`
	Namespaces               []NamespaceUsageApplyConfiguration                `json:"namespaces,omitempty"`
	Condition                []ClusterConditionApplyConfiguration              `json:"conditions,omitempty"`
	ClusterInfo              *InfoApplyConfiguration                           `json:"clusterInfo,omitempty"`
	DaemonEndpoints          *ClusterDaemonEndpointsApplyConfiguration         `json:"daemonEndpoints,omitempty"`
	Partitions               []string                                          `json:"partitions,omitempty"`
	SecretRef                *ClusterSecretRefApplyConfiguration               `json:"secretRef,omitempty"`
	Storage                  []string                                          `json:"storage,omitempty"`
}

// ClusterStatusApplyConfiguration constructs an declarative configuration of the ClusterStatus type for use with
// apply.
func ClusterStatus() *ClusterStatusApplyConfiguration {
	return &ClusterStatusApplyConfiguration{}
}

// WithConfiguration sets the Configuration field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Configuration field is set to the value of the last call.
func (b *ClusterStatusApplyConfiguration) WithConfiguration(value *KubeletForClusterConfigurationApplyConfiguration) *ClusterStatusApplyConfiguration {
	b.Configuration = value
	return b
}

// WithAddresses adds the given value to the Addresses field in the declarative configuration
// and returns the receiver, so that objects can be build by chaining "With" function invocations.
// If called multiple times, values provided by each call will be appended to the Addresses field.
func (b *ClusterStatusApplyConfiguration) WithAddresses(values ...*AddressApplyConfiguration) *ClusterStatusApplyConfiguration {
	for i := range values {
		if values[i] == nil {
			panic("nil value passed to WithAddresses")
		}
		b.Addresses = append(b.Addresses, *values[i])
	}
	return b
}

// WithPhase sets the Phase field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Phase field is set to the value of the last call.
func (b *ClusterStatusApplyConfiguration) WithPhase(value mizarcorev1alpha1.ClusterPhase) *ClusterStatusApplyConfiguration {
	b.Phase = &value
	return b
}

// WithAllocatable sets the Allocatable field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Allocatable field is set to the value of the last call.
func (b *ClusterStatusApplyConfiguration) WithAllocatable(value v1.ResourceList) *ClusterStatusApplyConfiguration {
	b.Allocatable = &value
	return b
}

// WithUsage sets the Usage field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Usage field is set to the value of the last call.
func (b *ClusterStatusApplyConfiguration) WithUsage(value v1.ResourceList) *ClusterStatusApplyConfiguration {
	b.Usage = &value
	return b
}

// WithNodes adds the given value to the Nodes field in the declarative configuration
// and returns the receiver, so that objects can be build by chaining "With" function invocations.
// If called multiple times, values provided by each call will be appended to the Nodes field.
func (b *ClusterStatusApplyConfiguration) WithNodes(values ...*NodeLeftResourceApplyConfiguration) *ClusterStatusApplyConfiguration {
	for i := range values {
		if values[i] == nil {
			panic("nil value passed to WithNodes")
		}
		b.Nodes = append(b.Nodes, *values[i])
	}
	return b
}

// WithPlanetClusterAllocatable sets the PlanetClusterAllocatable field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the PlanetClusterAllocatable field is set to the value of the last call.
func (b *ClusterStatusApplyConfiguration) WithPlanetClusterAllocatable(value v1.ResourceList) *ClusterStatusApplyConfiguration {
	b.PlanetClusterAllocatable = &value
	return b
}

// WithPlanetClusterUsage sets the PlanetClusterUsage field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the PlanetClusterUsage field is set to the value of the last call.
func (b *ClusterStatusApplyConfiguration) WithPlanetClusterUsage(value v1.ResourceList) *ClusterStatusApplyConfiguration {
	b.PlanetClusterUsage = &value
	return b
}

// WithPlanetClusterNodes adds the given value to the PlanetClusterNodes field in the declarative configuration
// and returns the receiver, so that objects can be build by chaining "With" function invocations.
// If called multiple times, values provided by each call will be appended to the PlanetClusterNodes field.
func (b *ClusterStatusApplyConfiguration) WithPlanetClusterNodes(values ...*NodeLeftResourceApplyConfiguration) *ClusterStatusApplyConfiguration {
	for i := range values {
		if values[i] == nil {
			panic("nil value passed to WithPlanetClusterNodes")
		}
		b.PlanetClusterNodes = append(b.PlanetClusterNodes, *values[i])
	}
	return b
}

// WithNamespaces adds the given value to the Namespaces field in the declarative configuration
// and returns the receiver, so that objects can be build by chaining "With" function invocations.
// If called multiple times, values provided by each call will be appended to the Namespaces field.
func (b *ClusterStatusApplyConfiguration) WithNamespaces(values ...*NamespaceUsageApplyConfiguration) *ClusterStatusApplyConfiguration {
	for i := range values {
		if values[i] == nil {
			panic("nil value passed to WithNamespaces")
		}
		b.Namespaces = append(b.Namespaces, *values[i])
	}
	return b
}

// WithCondition adds the given value to the Condition field in the declarative configuration
// and returns the receiver, so that objects can be build by chaining "With" function invocations.
// If called multiple times, values provided by each call will be appended to the Condition field.
func (b *ClusterStatusApplyConfiguration) WithCondition(values ...*ClusterConditionApplyConfiguration) *ClusterStatusApplyConfiguration {
	for i := range values {
		if values[i] == nil {
			panic("nil value passed to WithCondition")
		}
		b.Condition = append(b.Condition, *values[i])
	}
	return b
}

// WithClusterInfo sets the ClusterInfo field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the ClusterInfo field is set to the value of the last call.
func (b *ClusterStatusApplyConfiguration) WithClusterInfo(value *InfoApplyConfiguration) *ClusterStatusApplyConfiguration {
	b.ClusterInfo = value
	return b
}

// WithDaemonEndpoints sets the DaemonEndpoints field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the DaemonEndpoints field is set to the value of the last call.
func (b *ClusterStatusApplyConfiguration) WithDaemonEndpoints(value *ClusterDaemonEndpointsApplyConfiguration) *ClusterStatusApplyConfiguration {
	b.DaemonEndpoints = value
	return b
}

// WithPartitions adds the given value to the Partitions field in the declarative configuration
// and returns the receiver, so that objects can be build by chaining "With" function invocations.
// If called multiple times, values provided by each call will be appended to the Partitions field.
func (b *ClusterStatusApplyConfiguration) WithPartitions(values ...string) *ClusterStatusApplyConfiguration {
	for i := range values {
		b.Partitions = append(b.Partitions, values[i])
	}
	return b
}

// WithSecretRef sets the SecretRef field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the SecretRef field is set to the value of the last call.
func (b *ClusterStatusApplyConfiguration) WithSecretRef(value *ClusterSecretRefApplyConfiguration) *ClusterStatusApplyConfiguration {
	b.SecretRef = value
	return b
}

// WithStorage adds the given value to the Storage field in the declarative configuration
// and returns the receiver, so that objects can be build by chaining "With" function invocations.
// If called multiple times, values provided by each call will be appended to the Storage field.
func (b *ClusterStatusApplyConfiguration) WithStorage(values ...string) *ClusterStatusApplyConfiguration {
	for i := range values {
		b.Storage = append(b.Storage, values[i])
	}
	return b
}
