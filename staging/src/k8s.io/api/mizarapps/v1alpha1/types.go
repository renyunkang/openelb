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

package v1alpha1

import (
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

var _ = resource.Quantity{}

const (
	// annotation value "{\"cpu\":\"250m\",\"memory\":\"1024Mi\",\"pods\":\"1\"}"
	ResourceAnnotation string = "mizargalaxy.mizar-k8s.io/resources"

	CRDManifestsMetaAnnotation string = "mizargalaxy.mizar-k8s.io/manifestmetas"

	WorkloadAnnotation string = "mizargalaxy.mizar-k8s.io/workload"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Workload is the Scheme for the workload API.
type Workload struct {
	metav1.TypeMeta `json:",inline"`

	// Standard object's metadata.
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	// Spec defines the specification of the desired behavior of the ReplicaSet.
	// +optional
	Spec WorkloadSpec `json:"spec,omitempty" protobuf:"bytes,2,opt,name=spec"`

	// Status is the most recently observed status of the Workload.
	// This data may be out of date by some window of time.
	// Populated by the system.
	// Read-only.
	// +optional
	Status WorkloadStatus `json:"status,omitempty" protobuf:"bytes,3,opt,name=status"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// WorkloadList is a collection of ReplicaSets.
type WorkloadList struct {
	metav1.TypeMeta `json:",inline"`

	// Standard list metadata.
	// More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds
	// +optional
	metav1.ListMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	// List of Workloads.
	Items []Workload `json:"items" protobuf:"bytes,2,rep,name=items"`
}

// WorkloadSpec is the specification of a Workload.
type WorkloadSpec struct {

	// Types describe the type of Workload, could be stateless|stateful|stable
	// +optional
	Type WorkloadType `json:"type,omitempty" protobuf:"bytes,1,opt,name=type,casttype=WorkloadType"`

	// K8SObjects containers k8s object raw definitons if Type is k8s-object.
	K8SObjects []runtime.RawExtension `json:"k8sObjects" protobuf:"bytes,2,opt,name=k8sObjects"`

	// scheudler use DestAffinity to choose cluster and set location, kubelet-for-cluster reconcile workload
	// if location.NextCluster equal self cluster name.
	// +optional
	Location *Location `json:"location,omitempty" protobuf:"bytes,3,opt,name=location"`

	// If specified, the workload's scheduling constraints
	// +optional
	DestAffinity *Affinity `json:"destAffinity,omitempty" protobuf:"bytes,4,opt,name=destAffinity"`

	// If specified, the pod will be dispatched by specified scheduler.
	// If not specified, the pod will be dispatched by default scheduler.
	// +optional
	SchedulerName string `json:"schedulerName,omitempty" protobuf:"bytes,5,opt,name=schedulerName"`

	// CheckFilters special objects
	CheckFilters []K8SObject `json:"checkFilters,omitempty" protobuf:"bytes,6,rep,name=checkFilters"`
}

// K8SObject
type K8SObject struct {
	metav1.GroupVersionKind `json:",inline" protobuf:"bytes,1,opt,name=groupVersionKind"`
	// Name
	Name string `json:"name,omitempty" protobuf:"bytes,2,opt,name=name"`
}

const (
	// "default-scheduler" is the name of default scheduler.
	DefaultSchedulerName = "default-scheduler"
)

type K8SObjectStatus struct {
	Reference v1.ObjectReference `json:"reference" protobuf:"bytes,1,opt,name=reference"`

	// Scheduled specifies whether the resource is scheduled.
	Scheduled bool `json:"scheduled" protobuf:"varint,2,opt,name=scheduled"`

	// Ready specifies whether the resource is ready.
	Ready bool `json:"ready" protobuf:"varint,3,opt,name=ready"`

	// Deleted
	Deleted bool `json:"deleted" protobuf:"varint,4,opt,name=deleted"`

	// Replicas
	Replicas int32 `json:"replica" protobuf:"varint,5,opt,name=replica"`

	// ReadyReplicas
	ReadyReplicas int32 `json:"readyReplicas" protobuf:"varint,6,opt,name=readyReplicas"`
}

// WorkloadStatus represents the current status of a Workload.
type WorkloadStatus struct {
	// Phase is one of Running, Pending, Unknown, Terminating
	Phase WorkloadPhase `json:"phase" protobuf:"varint,1,opt,name=phase"`

	// observedGeneration is the most recent generation observed for this Workload. It corresponds to the
	// StatefulSet's generation, which is updated on mutation by the API Server.
	// +optional
	ObservedGeneration int64 `json:"observedGeneration,omitempty" protobuf:"varint,2,opt,name=observedGeneration"`

	// Represents the latest available observations of a replica set's current state.
	// +optional
	// +patchMergeKey=type
	// +patchStrategy=merge
	Conditions []WorkloadCondition `json:"conditions,omitempty" patchStrategy:"merge" patchMergeKey:"type" protobuf:"bytes,3,rep,name=conditions"`

	// K8SObjectStatus record the related resources created by this workload.
	K8SObjectStatus []K8SObjectStatus `json:"k8sObjectStatus" protobuf:"bytes,4,opt,name=k8sObjectStatus"`

	// AllocatedResources
	AllocatedResources v1.ResourceList `json:"allocatedResources" protobuf:"bytes,5,opt,name=allocatedResources"`

	// Region
	Region string `json:"region" protobuf:"bytes,6,opt,name=region"`

	// PhysicalZone
	PhysicalZone []string `json:"physicalZone" protobuf:"bytes,7,opt,name=physicalZone"`

	// LogicZone
	LogicZone []string `json:"logicZone" protobuf:"bytes,8,opt,name=logicZone"`
}

type WorkloadConditionType string

// These are valid conditions of a workload.
const (
	WorkloadScheduled WorkloadConditionType = "WorkloadScheduled"
)

// WorkloadCondition describes the state of a workload at a certain point.
type WorkloadCondition struct {
	// Type of replica set condition.
	Type WorkloadConditionType `json:"type" protobuf:"bytes,1,opt,name=type,casttype=WorkloadConditionType"`
	// Status of the condition, one of True, False, Unknown.
	Status v1.ConditionStatus `json:"status" protobuf:"bytes,2,opt,name=status,casttype=k8s.io/api/core/v1.ConditionStatus"`
	// The last time the condition transitioned from one status to another.
	// +optional
	LastTransitionTime metav1.Time `json:"lastTransitionTime,omitempty" protobuf:"bytes,3,opt,name=lastTransitionTime"`
	// The reason for the condition's last transition.
	// +optional
	Reason string `json:"reason,omitempty" protobuf:"bytes,4,opt,name=reason"`
	// A human readable message indicating details about the transition.
	// +optional
	Message string `json:"message,omitempty" protobuf:"bytes,5,opt,name=message"`
}

// Location
type Location struct {
	// NextCluster is a request to schedule this workload onto a specific cluster.
	// If NextCluster match affinity, then unpackage workload in this cluster. Otherwise, set NextCluster to empty and schedule
	// this workload to sub cluster.
	// +optional
	NextCluster string `json:"nextCluster,omitempty" protobuf:"bytes,1,opt,name=nextCluster"`
}

// Affinity is a group of affinity scheduling rules.
type Affinity struct {
	// Describes cluster affinity scheduling rules for the workload.
	// +optional
	ClusterAffinity *ClusterAffinity `json:"clusterAffinity,omitempty" protobuf:"bytes,1,opt,name=clusterAffinity"`
	// Describes workload affinity scheduling rules (e.g. co-locate this workload in the same cluster, node, zone, etc. as some other workload(s)).
	// +optional
	WorkloadAffinity *WorkloadAffinity `json:"workloadAffinity,omitempty" protobuf:"bytes,2,opt,name=workloadAffinity"`
	// Describes workload anti-affinity scheduling rules (e.g. avoid putting this workload in the same cluster, node, zone, etc. as some other workload(s)).
	// +optional
	WorkloadAntiAffinity *WorkloadAntiAffinity `json:"workloadAntiAffinity,omitempty" protobuf:"bytes,3,opt,name=workloadAntiAffinity"`
}

// ClusterAffinity is a group of cluster affinity scheduling rules.
type ClusterAffinity struct {

	// If the affinity requirements specified by this field are not met at
	// scheduling time, the workload will not be scheduled onto the cluster.
	// If the affinity requirements specified by this field cease to be met
	// at some point during workload execution (e.g. due to an update), the system
	// may or may not try to eventually evict the workload from its cluster.
	// +optional
	RequiredDuringSchedulingIgnoredDuringExecution *ClusterSelector `json:"requiredDuringSchedulingIgnoredDuringExecution,omitempty" protobuf:"bytes,1,opt,name=requiredDuringSchedulingIgnoredDuringExecution"`
	// The scheduler will prefer to schedule workloads to clusters that satisfy
	// the affinity expressions specified by this field, but it may choose
	// a cluster that violates one or more of the expressions. The cluster that is
	// most preferred is the one with the greatest sum of weights, i.e.
	// for each cluster that meets all of the scheduling requirements (resource
	// request, requiredDuringScheduling affinity expressions, etc.),
	// compute a sum by iterating through the elements of this field and adding
	// "weight" to the sum if the cluster matches the corresponding matchExpressions; the
	// cluster(s) with the highest sum are the most preferred.
	// +optional
	PreferredDuringSchedulingIgnoredDuringExecution []PreferredSchedulingTerm `json:"preferredDuringSchedulingIgnoredDuringExecution,omitempty" protobuf:"bytes,2,rep,name=preferredDuringSchedulingIgnoredDuringExecution"`
}

// WorkloadAffinity is a group of inter worklaod affinity scheduling rules.
type WorkloadAffinity struct {

	// If the affinity requirements specified by this field are not met at
	// scheduling time, the workload will not be scheduled onto the cluster.
	// If the affinity requirements specified by this field cease to be met
	// at some point during workload execution (e.g. due to a workload label update), the
	// system may or may not try to eventually evict the workload from its cluster.
	// When there are multiple elements, the lists of cluster corresponding to each
	// workloadAffinityTerm are intersected, i.e. all terms must be satisfied.
	// +optional
	RequiredDuringSchedulingIgnoredDuringExecution []WorkloadAffinityTerm `json:"requiredDuringSchedulingIgnoredDuringExecution,omitempty" protobuf:"bytes,1,rep,name=requiredDuringSchedulingIgnoredDuringExecution"`
	// The scheduler will prefer to schedule workloads to clusters that satisfy
	// the affinity expressions specified by this field, but it may choose
	// a cluster that violates one or more of the expressions. The cluster that is
	// most preferred is the one with the greatest sum of weights, i.e.
	// for each cluster that meets all of the scheduling requirements (resource
	// request, requiredDuringScheduling affinity expressions, etc.),
	// compute a sum by iterating through the elements of this field and adding
	// "weight" to the sum if the cluster has workloads which matches the corresponding workloadAffinityTerm; the
	// cluster(s) with the highest sum are the most preferred.
	// +optional
	PreferredDuringSchedulingIgnoredDuringExecution []WeightedWorkloadAffinityTerm `json:"preferredDuringSchedulingIgnoredDuringExecution,omitempty" protobuf:"bytes,2,rep,name=preferredDuringSchedulingIgnoredDuringExecution"`
}

// WorkloadAntiAffinity is a group of inter worklaod anti affinity scheduling rules.
type WorkloadAntiAffinity struct {

	// If the anti-affinity requirements specified by this field are not met at
	// scheduling time, the workload will not be scheduled onto the cluster.
	// If the anti-affinity requirements specified by this field cease to be met
	// at some point during workload execution (e.g. due to a worklaod label update), the
	// system may or may not try to eventually evict the workload from its cluster.
	// When there are multiple elements, the lists of clusters corresponding to each
	// workloadAffinityTerm are intersected, i.e. all terms must be satisfied.
	// +optional
	RequiredDuringSchedulingIgnoredDuringExecution []WorkloadAffinityTerm `json:"requiredDuringSchedulingIgnoredDuringExecution,omitempty" protobuf:"bytes,1,rep,name=requiredDuringSchedulingIgnoredDuringExecution"`
	// The scheduler will prefer to schedule workloads to clusters that satisfy
	// the anti-affinity expressions specified by this field, but it may choose
	// a cluster that violates one or more of the expressions. The cluster that is
	// most preferred is the one with the greatest sum of weights, i.e.
	// for each cluster that meets all of the scheduling requirements (resource
	// request, requiredDuringScheduling anti-affinity expressions, etc.),
	// compute a sum by iterating through the elements of this field and adding
	// "weight" to the sum if the cluster has workloads which matches the corresponding workloadAffinityTerm; the
	// workload(s) with the highest sum are the most preferred.
	// +optional
	PreferredDuringSchedulingIgnoredDuringExecution []WeightedWorkloadAffinityTerm `json:"preferredDuringSchedulingIgnoredDuringExecution,omitempty" protobuf:"bytes,2,rep,name=preferredDuringSchedulingIgnoredDuringExecution"`
}

// ClusterSelector A cluster selector represents the union of the results of one or more label queries
// over a set of clusters; that is, it represents the OR of the selectors represented
// by the cluster selector terms.
type ClusterSelector struct {
	//Required. A list of cluster selector terms. The terms are ORed.
	ClusterSelectorTerms []ClusterSelectorTerm `json:"clusterSelectorTerms" protobuf:"bytes,1,rep,name=clusterSelectorTerms"`
}

// ClusterSelectorTerm A null or empty cluster selector term matches no objects. The requirements of
// them are ANDed.
// The TopologySelectorTerm type implements a subset of the ClusterSelectorTerm.
type ClusterSelectorTerm struct {
	// A list of cluster selector requirements by cluster's labels.
	// +optional
	MatchExpressions []ClusterSelectorRequirement `json:"matchExpressions,omitempty" protobuf:"bytes,1,rep,name=matchExpressions"`
	// A list of cluster selector requirements by cluster's fields.
	// +optional
	MatchFields []ClusterSelectorRequirement `json:"matchFields,omitempty" protobuf:"bytes,2,rep,name=matchFields"`
}

// ClusterSelectorRequirement A cluster selector requirement is a selector that contains values, a key, and an operator
// that relates the key and values.
type ClusterSelectorRequirement struct {
	// The label key that the selector applies to.
	Key string `json:"key" protobuf:"bytes,1,opt,name=key"`
	// Represents a key's relationship to a set of values.
	// Valid operators are In, NotIn, Exists, DoesNotExist. Gt, and Lt.
	Operator ClusterSelectorOperator `json:"operator" protobuf:"bytes,2,opt,name=operator,casttype=ClusterSelectorOperator"`
	// An array of string values. If the operator is In or NotIn,
	// the values array must be non-empty. If the operator is Exists or DoesNotExist,
	// the values array must be empty. If the operator is Gt or Lt, the values
	// array must have a single element, which will be interpreted as an integer.
	// This array is replaced during a strategic merge patch.
	// +optional
	Values []string `json:"values,omitempty" protobuf:"bytes,3,rep,name=values"`
}

// ClusterSelectorOperator A cluster selector operator is the set of operators that can be used in
// a cluster selector requirement.
type ClusterSelectorOperator string

const (
	ClusterSelectorOpIn           ClusterSelectorOperator = "In"
	ClusterSelectorOpNotIn        ClusterSelectorOperator = "NotIn"
	ClusterSelectorOpExists       ClusterSelectorOperator = "Exists"
	ClusterSelectorOpDoesNotExist ClusterSelectorOperator = "DoesNotExist"
	ClusterSelectorOpGt           ClusterSelectorOperator = "Gt"
	ClusterSelectorOpLt           ClusterSelectorOperator = "Lt"
)

// PreferredSchedulingTerm An empty preferred scheduling term matches all objects with implicit weight 0
// (i.e. it's a no-op). A null preferred scheduling term matches no objects (i.e. is also a no-op).
type PreferredSchedulingTerm struct {
	// Weight associated with matching the corresponding clusterSelectorTerm, in the range 1-100.
	Weight int32 `json:"weight" protobuf:"varint,1,opt,name=weight"`
	// A cluster selector term, associated with the corresponding weight.
	Preference ClusterSelectorTerm `json:"preference" protobuf:"bytes,2,opt,name=preference"`
}

// WorkloadAffinityTerm Defines a set of workloads (namely those matching the labelSelector
// relative to the given namespace(s)) that this workload should be
// co-located (affinity) or not co-located (anti-affinity) with,
// where co-located is defined as running on a cluster whose value of
// the label with key <topologyKey> matches that of any cluster on which
// a cluster of the set of clusters is running
type WorkloadAffinityTerm struct {
	// A label query over a set of resources, in this case workloads.
	// +optional
	LabelSelector *metav1.LabelSelector `json:"labelSelector,omitempty" protobuf:"bytes,1,opt,name=labelSelector"`
	// namespaces specifies which namespaces the labelSelector applies to (matches against);
	// null or empty list means "this workload's namespace"
	// +optional
	Namespaces []string `json:"namespaces,omitempty" protobuf:"bytes,2,rep,name=namespaces"`
	// This workload should be co-located (affinity) or not co-located (anti-affinity) with the workloads matching
	// the labelSelector in the specified namespaces, where co-located is defined as running on a cluster
	// whose value of the label with key topologyKey matches that of any cluster on which any of the
	// selected workloads is running.
	// Empty topologyKey is not allowed.
	TopologyKey string `json:"topologyKey" protobuf:"bytes,3,opt,name=topologyKey"`
	// A label query over the set of namespaces that the term applies to.
	// The term is applied to the union of the namespaces selected by this field
	// and the ones listed in the namespaces field.
	// null selector and null or empty namespaces list means "this pod's namespace".
	// An empty selector ({}) matches all namespaces.
	// +optional
	NamespaceSelector *metav1.LabelSelector `json:"namespaceSelector,omitempty" protobuf:"bytes,4,opt,name=namespaceSelector"`
}

// WeightedWorkloadAffinityTerm The weights of all of the matched WeightedWorkloadAffinityTerm fields are added per-cluster to find the most preferred cluster(s)
type WeightedWorkloadAffinityTerm struct {
	// weight associated with matching the corresponding workloadAffinityTerm,
	// in the range 1-100.
	Weight int32 `json:"weight" protobuf:"varint,1,opt,name=weight"`
	// Required. A workload affinity term, associated with the corresponding weight.
	WorkloadAffinityTerm WorkloadAffinityTerm `json:"workloadAffinityTerm" protobuf:"bytes,2,opt,name=workloadAffinityTerm"`
}

type WorkloadType string

const (
	K8SOBJECT WorkloadType = "k8s-object"
)

type WorkloadPhase string

const (
	WorkloadPending  WorkloadPhase = "Pending"
	WorkloadCreating WorkloadPhase = "Creating"
	WorkloadUnpacked WorkloadPhase = "Unpacked"
	WorkloadRunning  WorkloadPhase = "Running"
	WorkloadUnknown  WorkloadPhase = "Unknown"
)

type WorkloadUpdateStrategy struct {
	// Type indicates the type of the WorkloadUpdateStrategyType.
	// Default is RollingUpdate.
	// +optional
	Type WorkloadUpdateStrategyType `json:"type,omitempty" protobuf:"bytes,1,opt,name=type,casttype=WorkloadUpdateStrategyType"`

	// RollingUpdate is used to communicate parameters when Type is RollingUpdateWorkloadStrategyType.
	// +optional
	RollingUpdate *RollingUpdateWorkloadStrategy `json:"rollingUpdate,omitempty" protobuf:"bytes,2,opt,name=rollingUpdate"`
}

// WorkloadUpdateStrategyType is a string enumeration type that enumerates
// all possible update strategies for the Workload controller.
type WorkloadUpdateStrategyType string

const (
	// RollingUpdateWorkloadStrategyType indicates that update will be
	// applied to all Pods in the Workload with respect to the Workload
	// ordering constraints. When a scale operation is performed with this
	// strategy, new Pods will be created from the specification version indicated
	// by the Workload's updateRevision.
	RollingUpdateWorkloadStrategyType WorkloadUpdateStrategyType = "RollingUpdate"
	// OnDeleteWorkloadStrategyType triggers the legacy behavior. Version
	// tracking and ordered rolling restarts are disabled. Pods are recreated
	// from the WorkloadSpec when they are manually deleted. When a scale
	// operation is performed with this strategy,specification version indicated
	// by the Workload's currentRevision.
	OnDeleteWorkloadStrategyType WorkloadUpdateStrategyType = "OnDelete"
)

// RollingUpdateWorkloadStrategy is used to communicate parameter for RollingUpdateWorkloadStrategyType.
type RollingUpdateWorkloadStrategy struct {
	// Partition indicates the ordinal at which the Workload should be
	// partitioned.
	Partition *int32 `json:"partition,omitempty" protobuf:"varint,1,opt,name=partition"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Binding ties one object to another; for example, a pod is bound to a node by a scheduler.
// Deprecated in 1.7, please use the bindings subresource of pods instead.
type Binding struct {
	metav1.TypeMeta `json:",inline"`
	// ObjectMeta describes the object that is being bound.
	// +optional
	metav1.ObjectMeta `json:"objectMeta,omitempty" protobuf:"bytes,1,opt,name=objectMeta"`

	// Target is the object to bind to.
	Target v1.ObjectReference `json:"target" protobuf:"bytes,2,opt,name=target"`
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// GlobalQuota sets aggregate quota restrictions enforced per namespace
type GlobalQuota struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	// Spec defines the desired quota.
	// +optional
	Spec v1.ResourceQuotaSpec `json:"spec,omitempty" protobuf:"bytes,2,opt,name=spec"`

	// Status defines the actual enforced quota and its current usage.
	// +optional
	Status v1.ResourceQuotaStatus `json:"status,omitempty" protobuf:"bytes,3,opt,name=status"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// GlobalQuotaList is a list of GlobalQuota resources.
type GlobalQuotaList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	Items []GlobalQuota `json:"items" protobuf:"bytes,2,rep,name=items"`
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// GlobalQuota sets aggregate quota restrictions enforced per namespace
type ClusterQuota struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	// Spec defines the desired quota.
	// +optional
	Spec v1.ClusterResourceQuotaSpec `json:"spec,omitempty" protobuf:"bytes,2,opt,name=spec"`

	// Status defines the actual enforced quota and its current usage.
	// +optional
	Status v1.ResourceQuotaStatus `json:"status,omitempty" protobuf:"bytes,3,opt,name=status"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// GlobalQuotaList is a list of GlobalQuota resources.
type ClusterQuotaList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	Items []ClusterQuota `json:"items" protobuf:"bytes,2,rep,name=items"`
}

// +genclient
// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// CRDResourceHook
type CRDResourceHook struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	Spec   CRDResourceHookSpec   `json:"spec,omitempty" protobuf:"bytes,2,opt,name=spec"`
	Status CRDResourceHookStatus `json:"status,omitempty" protobuf:"bytes,3,opt,name=status"`
}

// CRDResourceHookSpec
type CRDResourceHookSpec struct {
	// Hook
	Hook ResourceCalculateHook `json:"hook,omitempty" protobuf:"bytes,1,opt,name=hook"`

	// CRDReference
	CRDReference metav1.GroupVersionKind `json:"crdReference,omitempty" protobuf:"bytes,2,opt,name=crdReference"`
}

type CRDResourceHookStatus struct {
	Active bool `json:"active,omitempty" protobuf:"varint,1,opt,name=active"`
}

// ResourceCalculateHook
type ResourceCalculateHook struct {
	// `url` gives the location of the webhook, in standard URL form
	// (`scheme://host:port/path`). Exactly one of `url` or `service`
	// must be specified.
	//
	// The `host` should not refer to a service running in the cluster; use
	// the `service` field instead. The host might be resolved via external
	// DNS in some apiservers (e.g., `kube-apiserver` cannot resolve
	// in-cluster DNS as that would be a layering violation). `host` may
	// also be an IP address.
	//
	// Please note that using `localhost` or `127.0.0.1` as a `host` is
	// risky unless you take great care to run this webhook on all hosts
	// which run an apiserver which might need to make calls to this
	// webhook. Such installs are likely to be non-portable, i.e., not easy
	// to turn up in a new cluster.
	//
	// The scheme must be "https"; the URL must begin with "https://".
	//
	// A path is optional, and if present may be any string permissible in
	// a URL. You may use the path to pass an arbitrary string to the
	// webhook, for example, a cluster identifier.
	//
	// Attempting to use a user or basic auth e.g. "user:password@" is not
	// allowed. Fragments ("#...") and query parameters ("?...") are not
	// allowed, either.
	//
	// +optional
	URL *string `json:"url,omitempty" protobuf:"bytes,1,opt,name=url"`

	// `caBundle` is a PEM encoded CA bundle which will be used to validate the webhook's server certificate.
	// If unspecified, system trust roots on the apiserver are used.
	// +optional
	CABase64 string `json:"caBase64,omitempty" protobuf:"bytes,2,opt,name=caBase64"`

	// CertBase64
	CertBase64 string `json:"certBase64,omitempty" protobuf:"bytes,3,opt,name=certBase64"`

	// KeyBase64
	KeyBase64 string `json:"keyBase64,omitempty" protobuf:"bytes,4,opt,name=keyBase64"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// CRDResourceHookList is a list of GlobalQuota resources.
type CRDResourceHookList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	Items []CRDResourceHook `json:"items" protobuf:"bytes,2,rep,name=items"`
}
