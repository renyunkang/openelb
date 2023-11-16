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

// AffinityApplyConfiguration represents an declarative configuration of the Affinity type for use
// with apply.
type AffinityApplyConfiguration struct {
	ClusterAffinity      *ClusterAffinityApplyConfiguration      `json:"clusterAffinity,omitempty"`
	WorkloadAffinity     *WorkloadAffinityApplyConfiguration     `json:"workloadAffinity,omitempty"`
	WorkloadAntiAffinity *WorkloadAntiAffinityApplyConfiguration `json:"workloadAntiAffinity,omitempty"`
}

// AffinityApplyConfiguration constructs an declarative configuration of the Affinity type for use with
// apply.
func Affinity() *AffinityApplyConfiguration {
	return &AffinityApplyConfiguration{}
}

// WithClusterAffinity sets the ClusterAffinity field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the ClusterAffinity field is set to the value of the last call.
func (b *AffinityApplyConfiguration) WithClusterAffinity(value *ClusterAffinityApplyConfiguration) *AffinityApplyConfiguration {
	b.ClusterAffinity = value
	return b
}

// WithWorkloadAffinity sets the WorkloadAffinity field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the WorkloadAffinity field is set to the value of the last call.
func (b *AffinityApplyConfiguration) WithWorkloadAffinity(value *WorkloadAffinityApplyConfiguration) *AffinityApplyConfiguration {
	b.WorkloadAffinity = value
	return b
}

// WithWorkloadAntiAffinity sets the WorkloadAntiAffinity field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the WorkloadAntiAffinity field is set to the value of the last call.
func (b *AffinityApplyConfiguration) WithWorkloadAntiAffinity(value *WorkloadAntiAffinityApplyConfiguration) *AffinityApplyConfiguration {
	b.WorkloadAntiAffinity = value
	return b
}
