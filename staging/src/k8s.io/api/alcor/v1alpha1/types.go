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
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// StableModel is the Scheme for the stableModel API.
type StableModel struct {
	metav1.TypeMeta `json:",inline"`

	// Standard object's metadata.
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=objectMeta"`

	// Spec defines the specification of the desired behavior of the ReplicaSet.
	// +optional
	Spec StableModelSpec `json:"spec,omitempty" protobuf:"bytes,2,opt,name=spec"`

	// Status is the most recently observed status of the StableModel.
	// This data may be out of date by some window of time.
	// Populated by the system.
	// Read-only.
	// +optional
	Status StableModelStatus `json:"status,omitempty" protobuf:"bytes,3,opt,name=status"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// StableModelList is a collection of ReplicaSets.
type StableModelList struct {
	metav1.TypeMeta `json:",inline"`

	// Standard list metadata.
	// More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds
	// +optional
	metav1.ListMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=listMeta"`

	// List of StableModels.
	Items []StableModel `json:"items" protobuf:"bytes,2,rep,name=items"`
}

// StableModelSpec is the specification of a StableModel.
type StableModelSpec struct {

	// Selector is a label query over pods that should match the replica count.
	// Label keys and values that must match in order to be controlled by this replica set.
	// It must match the pod template's labels.
	// More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/labels/#label-selectors
	Selector *metav1.LabelSelector `json:"selector,omitempty" protobuf:"bytes,1,opt,name=selector"`

	Mode Mode `json:"mode,omitempty" protobuf:"bytes,2,opt,name=mode,casttype=Mode"`

	// Template is the object that describes the pod that will be created if
	// insufficient replicas are detected.
	// More info: https://kubernetes.io/docs/concepts/stableModels/controllers/replicationcontroller#pod-template
	// +optional
	PodTemplate v1.PodTemplateSpec `json:"podTemplate" protobuf:"bytes,3,opt,name=podTemplate"`

	Network NetworkRequest `json:"network,omitempty" protobuf:"bytes,4,opt,name=network"`

	Storage []StorageRequest `json:"storage,omitempty" protobuf:"bytes,5,rep,name=storage"`

	Feature Features `json:"feature,omitempty" protobuf:"bytes,6,opt,name=feature"`

	NodeName string `json:"nodeName,omitempty" protobuf:"bytes,7,opt,name=nodeName"`

	ResizePhase ResizePhase `json:"resizePhase,omitempty" protobuf:"bytes,12,opt,name=resizePhase,casttype=ResizePhase"`

	// ClusterSelector is a selector which must be true for the stableModel to fit on a cluster.
	// Selector which must match a cluster's labels for the stableModel to be scheduled on that cluster.
	// +optional
	ClusterSelector map[string]string `json:"clusterSelector,omitempty" protobuf:"bytes,8,rep,name=clusterSelector"`

	// ClusterName is a request to schedule this stableModel onto a specific cluster. If it is non-empty,
	// the scheduler simply schedules this stableModel onto that cluster, assuming that it fits resource
	// requirements.
	// +optional
	ClusterName string `json:"clusterName,omitempty" protobuf:"bytes,9,opt,name=clusterName"`

	// If specified, the stableModel's scheduling constraints
	// +optional
	Affinity *Affinity `json:"affinity,omitempty" protobuf:"bytes,10,opt,name=affinity"`

	// If specified, the stablemodel's tolerations.
	// +optional
	Tolerations []v1.Toleration `json:"tolerations,omitempty" protobuf:"bytes,13,rep,name=tolerations"`

	// If specified, the pod will be dispatched by specified scheduler.
	// If not specified, the pod will be dispatched by default scheduler.
	// +optional
	SchedulerName string `json:"schedulerName,omitempty" protobuf:"bytes,11,opt,name=schedulerName"`

	// TopologySpreadConstraints describes how a group of stablemodels ought to spread across topology
	// domains. Scheduler will schedule stablemodels in a way which abides by the constraints.
	// All topologySpreadConstraints are ANDed.
	// +optional
	// +patchMergeKey=topologyKey
	// +patchStrategy=merge
	// +listType=map
	// +listMapKey=topologyKey
	// +listMapKey=whenUnsatisfiable
	TopologySpreadConstraints []v1.TopologySpreadConstraint `json:"topologySpreadConstraints,omitempty" patchStrategy:"merge" patchMergeKey:"topologyKey" protobuf:"bytes,14,opt,name=topologySpreadConstraints"`
}

const (
	// "default-scheduler" is the name of default scheduler.
	DefaultSchedulerName = "default-scheduler"
)

type Mode string

const (
	DirectMode Mode = "direct"
	LegacyMode Mode = "legacy"
)

type ResizePhase string

const (
	PreResize       ResizePhase = "PreResize"
	PreResizing     ResizePhase = "PreResizing"
	PreResized      ResizePhase = "PreResized"
	PreResizeFailed ResizePhase = "PreResizeFailed"
	ResizeRollBack  ResizePhase = "ResizeRollBack"
	Resize          ResizePhase = "Resize"
	Resizing        ResizePhase = "Resizing"
	Resized         ResizePhase = "Resized"
	ResizeFailed    ResizePhase = "ResizeFailed"
)

// 网络申请
type NetworkConfig struct {
	IP     string `json:"ip,omitempty" protobuf:"bytes,1,opt,name=iP"`
	NIC    string `json:"nic,omitempty" protobuf:"bytes,2,opt,name=nIC"`
	NICMAC string `json:"nicmac,omitempty" protobuf:"bytes,3,opt,name=nICMAC"`
	Host   string `json:"host,omitempty" protobuf:"bytes,4,opt,name=host"`
	Pool   string `json:"pool,omitempty" protobuf:"bytes,5,opt,name=pool"`
	Mbps   string `json:"mbps,omitempty" protobuf:"bytes,6,opt,name=mbps"`
	VNI    bool   `json:"vni,omitempty" protobuf:"varint,7,opt,name=vNI"`
}

type NetworkRequest struct {
	Name             string        `json:"name,omitempty" protobuf:"bytes,1,opt,name=name"`
	NetworkClassName string        `json:"networkClassName,omitempty" protobuf:"bytes,2,opt,name=networkClassName"`
	Config           NetworkConfig `json:"config,omitempty" protobuf:"bytes,3,opt,name=config"`
}

type StorageRequest struct {
	Name string                       `json:"name,omitempty" protobuf:"bytes,1,opt,name=name"`
	Spec v1.PersistentVolumeClaimSpec `json:"spec,omitempty" protobuf:"bytes,2,opt,name=spec"`
}

type ResourceLimit struct {
	Cpu    resource.Quantity `json:"cpu,omitempty" protobuf:"bytes,1,opt,name=cpu"`
	Memory resource.Quantity `json:"memory,omitempty" protobuf:"bytes,2,opt,name=memory"`
}

// ContainersResource
type ContainersResource map[string]ResourceLimit

// ResizeResource
type ResizeResource struct {
	PodResource        ResourceLimit      `json:"podResource,omitempty" protobuf:"bytes,1,opt,name=podResource"`
	ContainersResource ContainersResource `json:"containersResource,omitempty" protobuf:"bytes,2,rep,name=containersResource"`
}

// 特性
type Features struct {
	RebootPod              bool           `json:"rebootPod,omitempty" protobuf:"varint,1,opt,name=rebootPod"`
	ForceDriftPod          bool           `json:"forceDriftPod,omitempty" protobuf:"varint,2,opt,name=forceDriftPod"`
	PausePod               bool           `json:"pausePod" protobuf:"varint,3,opt,name=pausePod"`
	DriftPod               bool           `json:"driftPod,omitempty" protobuf:"varint,4,opt,name=driftPod"`
	RebootPodTimestamp     string         `json:"rebootPodTimestamp,omitempty" protobuf:"bytes,5,opt,name=rebootPodTimestamp"`
	ForceDriftPodTimestamp string         `json:"forceDriftPodTimestamp,omitempty" protobuf:"bytes,6,opt,name=forceDriftPodTimestamp"`
	PreResizePod           ResizeResource `json:"preResizePod,omitempty" protobuf:"bytes,7,opt,name=preResizePod"`
	ResizePod              ResizeResource `json:"resizePod,omitempty" protobuf:"bytes,8,opt,name=resizePod"`
	DelayUpdatePod         bool           `json:"delayUpdatePod" protobuf:"varint,9,opt,name=delayUpdatePod"`
}

// PVC状态
type PVCStatus struct {
	Name     string                        `json:"name" protobuf:"bytes,1,opt,name=name"`
	Status   v1.PersistentVolumeClaimPhase `json:"status" protobuf:"bytes,2,opt,name=status,casttype=k8s.io/api/core/v1.PersistentVolumeClaimPhase"`
	Capacity v1.ResourceList               `json:"capacity,omitempty" protobuf:"bytes,3,rep,name=capacity,casttype=k8s.io/api/core/v1.ResourceList,castkey=k8s.io/api/core/v1.ResourceName"`
}

// Pod状态
type PodStatus struct {
	Name   string `json:"name" protobuf:"bytes,1,opt,name=name"`
	Status Phase  `json:"status,omitempty" protobuf:"bytes,2,opt,name=status,casttype=Phase"`
}

// IPC状态
type IPStatus struct {
	Name    string `json:"name" protobuf:"bytes,1,opt,name=name"`
	Status  string `json:"status,omitempty" protobuf:"bytes,2,opt,name=status"`
	Ip      string `json:"ip,omitempty" protobuf:"bytes,3,opt,name=ip"`
	Backend string `json:"backend,omitempty" protobuf:"bytes,4,opt,name=backend"`
}

// 事件记录
type EventRecord struct {
	EventType   string `json:"eventType" protobuf:"bytes,1,opt,name=eventType"`
	CreateTime  string `json:"createTime" protobuf:"bytes,2,opt,name=createTime"`
	LastMessage string `json:"lastMessage" protobuf:"bytes,3,opt,name=lastMessage"`
}

type ResizeStatus string

// StableModelStatus defines the observed state of StableModel
type StableModelStatus struct {
	Phase      Phase         `json:"phase" protobuf:"bytes,1,opt,name=phase,casttype=Phase"`
	PodRef     PodStatus     `json:"podRef,omitempty" protobuf:"bytes,2,opt,name=podRef"`
	IPClaimRef IPStatus      `json:"ipRef,omitempty" protobuf:"bytes,3,opt,name=ipRef"`
	PVCRefs    []PVCStatus   `json:"pvcRefs,omitempty" protobuf:"bytes,4,rep,name=pvcRefs"`
	Events     []EventRecord `json:"events,omitempty" protobuf:"bytes,5,rep,name=events"`
	// Represents the latest available observations of a replica set's current state.
	// +optional
	// +patchMergeKey=type
	// +patchStrategy=merge
	Conditions []StableModelCondition `json:"conditions,omitempty" patchStrategy:"merge" patchMergeKey:"type" protobuf:"bytes,6,rep,name=conditions"`

	ResizePhase            ResizePhase `json:"resizePhase,omitempty" protobuf:"bytes,7,opt,name=resizePhase,casttype=ResizePhase"`
	RebootPodTimestamp     string      `json:"rebootPodTimestamp,omitempty" protobuf:"bytes,8,opt,name=rebootPodTimestamp"`
	ForceDriftPodTimestamp string      `json:"forceDriftPodTimestamp,omitempty" protobuf:"bytes,9,opt,name=forceDriftPodTimestamp"`
	DriftPod               bool        `json:"driftPod,omitempty" protobuf:"varint,10,opt,name=driftPod"`
	NodeName               string      `json:"nodeName,omitempty" protobuf:"bytes,11,opt,name=nodeName"`
}

type StableModelConditionType string

// These are valid conditions of a stablemodel.
const (
	StableModelReplicaFailureType StableModelConditionType = "ReplicaFailure"

	StableModelScheduled StableModelConditionType = "StableModelScheduled"
)

// StableModelCondition describes the state of a stablemodel at a certain point.
type StableModelCondition struct {
	// Type of replica set condition.
	Type StableModelConditionType `json:"type" protobuf:"bytes,1,opt,name=type,casttype=StableModelConditionType"`
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

// Affinity is a group of affinity scheduling rules.
type Affinity struct {
	// Describes cluster affinity scheduling rules for the stablemodel.
	// +optional
	ClusterAffinity *ClusterAffinity `json:"clusterAffinity,omitempty" protobuf:"bytes,1,opt,name=clusterAffinity"`
	// Describes stablemodel affinity scheduling rules (e.g. co-locate this stablemodel in the same cluster, node, zone, etc. as some other stablemodel(s)).
	// +optional
	StableModelAffinity *StableModelAffinity `json:"stablemodelAffinity,omitempty" protobuf:"bytes,2,opt,name=stablemodelAffinity"`
	// Describes stablemodel anti-affinity scheduling rules (e.g. avoid putting this stablemodel in the same cluster, node, zone, etc. as some other stablemodel(s)).
	// +optional
	StableModelAntiAffinity *StableModelAntiAffinity `json:"stablemodelAntiAffinity,omitempty" protobuf:"bytes,3,opt,name=stablemodelAntiAffinity"`
}

// ClusterAffinity is a group of cluster affinity scheduling rules.
type ClusterAffinity struct {

	// If the affinity requirements specified by this field are not met at
	// scheduling time, the stablemodel will not be scheduled onto the cluster.
	// If the affinity requirements specified by this field cease to be met
	// at some point during stablemodel execution (e.g. due to an update), the system
	// may or may not try to eventually evict the stablemodel from its cluster.
	// +optional
	RequiredDuringSchedulingIgnoredDuringExecution *ClusterSelector `json:"requiredDuringSchedulingIgnoredDuringExecution,omitempty" protobuf:"bytes,1,opt,name=requiredDuringSchedulingIgnoredDuringExecution"`
	// The scheduler will prefer to schedule stablemodels to clusters that satisfy
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

// StableModelAffinity is a group of inter worklaod affinity scheduling rules.
type StableModelAffinity struct {

	// If the affinity requirements specified by this field are not met at
	// scheduling time, the stablemodel will not be scheduled onto the cluster.
	// If the affinity requirements specified by this field cease to be met
	// at some point during stablemodel execution (e.g. due to a stablemodel label update), the
	// system may or may not try to eventually evict the stablemodel from its cluster.
	// When there are multiple elements, the lists of cluster corresponding to each
	// stablemodelAffinityTerm are intersected, i.e. all terms must be satisfied.
	// +optional
	RequiredDuringSchedulingIgnoredDuringExecution []StableModelAffinityTerm `json:"requiredDuringSchedulingIgnoredDuringExecution,omitempty" protobuf:"bytes,1,rep,name=requiredDuringSchedulingIgnoredDuringExecution"`
	// The scheduler will prefer to schedule stablemodels to clusters that satisfy
	// the affinity expressions specified by this field, but it may choose
	// a cluster that violates one or more of the expressions. The cluster that is
	// most preferred is the one with the greatest sum of weights, i.e.
	// for each cluster that meets all of the scheduling requirements (resource
	// request, requiredDuringScheduling affinity expressions, etc.),
	// compute a sum by iterating through the elements of this field and adding
	// "weight" to the sum if the cluster has stablemodels which matches the corresponding stablemodelAffinityTerm; the
	// cluster(s) with the highest sum are the most preferred.
	// +optional
	PreferredDuringSchedulingIgnoredDuringExecution []WeightedStableModelAffinityTerm `json:"preferredDuringSchedulingIgnoredDuringExecution,omitempty" protobuf:"bytes,2,rep,name=preferredDuringSchedulingIgnoredDuringExecution"`
}

// StableModelAntiAffinity is a group of inter worklaod anti affinity scheduling rules.
type StableModelAntiAffinity struct {

	// If the anti-affinity requirements specified by this field are not met at
	// scheduling time, the stablemodel will not be scheduled onto the cluster.
	// If the anti-affinity requirements specified by this field cease to be met
	// at some point during stablemodel execution (e.g. due to a worklaod label update), the
	// system may or may not try to eventually evict the stablemodel from its cluster.
	// When there are multiple elements, the lists of clusters corresponding to each
	// stablemodelAffinityTerm are intersected, i.e. all terms must be satisfied.
	// +optional
	RequiredDuringSchedulingIgnoredDuringExecution []StableModelAffinityTerm `json:"requiredDuringSchedulingIgnoredDuringExecution,omitempty" protobuf:"bytes,1,rep,name=requiredDuringSchedulingIgnoredDuringExecution"`
	// The scheduler will prefer to schedule stablemodels to clusters that satisfy
	// the anti-affinity expressions specified by this field, but it may choose
	// a cluster that violates one or more of the expressions. The cluster that is
	// most preferred is the one with the greatest sum of weights, i.e.
	// for each cluster that meets all of the scheduling requirements (resource
	// request, requiredDuringScheduling anti-affinity expressions, etc.),
	// compute a sum by iterating through the elements of this field and adding
	// "weight" to the sum if the cluster has stablemodels which matches the corresponding stablemodelAffinityTerm; the
	// stablemodel(s) with the highest sum are the most preferred.
	// +optional
	PreferredDuringSchedulingIgnoredDuringExecution []WeightedStableModelAffinityTerm `json:"preferredDuringSchedulingIgnoredDuringExecution,omitempty" protobuf:"bytes,2,rep,name=preferredDuringSchedulingIgnoredDuringExecution"`
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

// StableModelAffinityTerm Defines a set of stablemodels (namely those matching the labelSelector
// relative to the given namespace(s)) that this stablemodel should be
// co-located (affinity) or not co-located (anti-affinity) with,
// where co-located is defined as running on a cluster whose value of
// the label with key <topologyKey> matches that of any cluster on which
// a cluster of the set of clusters is running
type StableModelAffinityTerm struct {
	// A label query over a set of resources, in this case stablemodels.
	// +optional
	LabelSelector *metav1.LabelSelector `json:"labelSelector,omitempty" protobuf:"bytes,1,opt,name=labelSelector"`
	// namespaces specifies which namespaces the labelSelector applies to (matches against);
	// null or empty list means "this stablemodel's namespace"
	// +optional
	Namespaces []string `json:"namespaces,omitempty" protobuf:"bytes,2,rep,name=namespaces"`
	// This stablemodel should be co-located (affinity) or not co-located (anti-affinity) with the stablemodels matching
	// the labelSelector in the specified namespaces, where co-located is defined as running on a cluster
	// whose value of the label with key topologyKey matches that of any cluster on which any of the
	// selected stablemodels is running.
	// Empty topologyKey is not allowed.
	TopologyKey string `json:"topologyKey" protobuf:"bytes,3,opt,name=topologyKey"`
}

// WeightedStableModelAffinityTerm The weights of all of the matched WeightedStableModelAffinityTerm fields are added per-cluster to find the most preferred cluster(s)
type WeightedStableModelAffinityTerm struct {
	// weight associated with matching the corresponding stablemodelAffinityTerm,
	// in the range 1-100.
	Weight int32 `json:"weight" protobuf:"varint,1,opt,name=weight"`
	// Required. A stablemodel affinity term, associated with the corresponding weight.
	StableModelAffinityTerm StableModelAffinityTerm `json:"stablemodelAffinityTerm" protobuf:"bytes,2,opt,name=stablemodelAffinityTerm"`
}

type StableModelType string

const (
	Stateless StableModelType = "stateless"
	Stateful  StableModelType = "stateful"
	Stable    StableModelType = "stable"
)

type Phase string

const (
	StableModelRunning     Phase = "Running"
	StableModelPending     Phase = "Pending"
	StableModelUnknown     Phase = "Unknown"
	StableModelTerminating Phase = "Terminating"
)

type StableModelUpdateStrategy struct {
	// Type indicates the type of the StableModelUpdateStrategyType.
	// Default is RollingUpdate.
	// +optional
	Type StableModelUpdateStrategyType `json:"type,omitempty" protobuf:"bytes,1,opt,name=type,casttype=StableModelUpdateStrategyType"`

	// RollingUpdate is used to communicate parameters when Type is RollingUpdateStableModelStrategyType.
	// +optional
	RollingUpdate *RollingUpdateStableModelStrategy `json:"rollingUpdate,omitempty" protobuf:"bytes,2,opt,name=rollingUpdate"`
}

// StableModelUpdateStrategyType is a string enumeration type that enumerates
// all possible update strategies for the StableModel controller.
type StableModelUpdateStrategyType string

const (
	// RollingUpdateStableModelStrategyType indicates that update will be
	// applied to all Pods in the StableModel with respect to the StableModel
	// ordering constraints. When a scale operation is performed with this
	// strategy, new Pods will be created from the specification version indicated
	// by the StableModel's updateRevision.
	RollingUpdateStableModelStrategyType StableModelUpdateStrategyType = "RollingUpdate"
	// OnDeleteStableModelStrategyType triggers the legacy behavior. Version
	// tracking and ordered rolling restarts are disabled. Pods are recreated
	// from the StableModelSpec when they are manually deleted. When a scale
	// operation is performed with this strategy,specification version indicated
	// by the StableModel's currentRevision.
	OnDeleteStableModelStrategyType StableModelUpdateStrategyType = "OnDelete"
)

// RollingUpdateStableModelStrategy is used to communicate parameter for RollingUpdateStableModelStrategyType.
type RollingUpdateStableModelStrategy struct {
	// Partition indicates the ordinal at which the StableModel should be
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

// IPSetSpec defines the desired state of IPSet
type IPSetSpec struct {
	// Number of IPs to claim. Should not less than replicas of deployment or statefulset.
	Replicas int64 `json:"replicas" protobuf:"varint,1,opt,name=replicas"`
	// From which IPPool IPSet is going to claim IP from. Can be omited for overlay scenario.
	Pool string `json:"pool,omitempty" protobuf:"bytes,2,opt,name=pool"`
	// For VNI case, only for overlay scenario.
	VNI bool `json:"vni,omitempty" protobuf:"varint,3,opt,name=vni"`
	// Deprecated, ignored.
	Sync bool `json:"sync,omitempty" protobuf:"varint,4,opt,name=sync"`
	// Deprecated, ignored.
	Deprecated bool `json:"deprecated,omitempty" protobuf:"varint,5,opt,name=deprecated"`
}

// IPSetStatus defines the observed state of IPSet
type IPSetStatus struct {
	Status         string            `json:"Status" protobuf:"bytes,1,opt,name=Status,json=status"`
	Message        string            `json:"Message,omitempty" protobuf:"bytes,2,opt,name=Message,json=message"`
	IPs            []string          `json:"IPs,omitempty" protobuf:"bytes,3,rep,name=IPs,json=iPs"`
	DirtyIPs       []string          `json:"DirtyIPs,omitempty" protobuf:"bytes,4,rep,name=DirtyIPs,json=dirtyIPs"`
	IPMap          map[string]string `json:"IPMap,omitempty" protobuf:"bytes,5,rep,name=IPMap,json=iPMap"`
	LastApplied    map[string]string `json:"LastApplied,omitempty" protobuf:"bytes,6,rep,name=LastApplied,json=lastApplied"`
	CallbackReason string            `json:"CallbackReason,omitempty" protobuf:"bytes,7,opt,name=CallbackReason,json=callbackReason"`
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type IPSet struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	Spec   IPSetSpec   `json:"spec,omitempty" protobuf:"bytes,2,opt,name=spec"`
	Status IPSetStatus `json:"status,omitempty" protobuf:"bytes,3,opt,name=status"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// IPSetList contains a list of IPSet
type IPSetList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`
	Items           []IPSet `json:"items" protobuf:"bytes,2,rep,name=items"`
}

// IPPoolSpec defines the desired state of IPPool
type IPPoolSpec struct {
	// CopyFrom an IPPool with some spec overwrited. The main target for "copy" is
	// .Status.AllocatedIPs, so admin can set spec to overwrite like
	// Route/Mask/Vlan, but not FirstIP/LastIP which can cause IP pool range changing.
	CopyFrom string `json:"copyFrom,omitempty" protobuf:"bytes,1,opt,name=copyFrom"`
	// FirstIP is the first IP of IPPool, when the pool does not start with subnet first IP
	FirstIP string `json:"firstIP,omitempty" protobuf:"bytes,2,opt,name=firstIP"`
	// FirstIP is the first IP of IPPool, when the pool does not start with subnet last IP
	LastIP string `json:"lastIP,omitempty" protobuf:"bytes,3,opt,name=lastIP"`
	// AdminStatus marks IPPool as maintained by admin
	AdminStatus bool `json:"adminStatus,omitempty" protobuf:"varint,4,opt,name=adminStatus"`
	// Deprecated marks IPPool as deprecated, it only accept release IP, no more allocate IP
	Deprecated bool `json:"deprecated,omitempty" protobuf:"varint,5,opt,name=deprecated"`
	// Detect is a trigger will make operator do detect to find any dirty IP exists
	Detect bool `json:"detect,omitempty" protobuf:"varint,6,opt,name=detect"`
	// Grid will let IPClaim and IPClaim allocate IPs in different order, enabled by default.
	Grid bool `json:"grid,omitempty" protobuf:"varint,7,opt,name=grid"`
	// SpecialClaimedQuota limits count to create IPClaim with lb/cfs as user. Emit means no limit, 0 means not allow.
	SpecialClaimedQuota map[string]int64 `json:"specialClaimedQuota,omitempty" protobuf:"bytes,8,rep,name=specialClaimedQuota"`
	// Defautl route/gateway in IPPool, only for underlay scenario
	Route string `json:"route,omitempty" protobuf:"bytes,9,opt,name=route"`
	// Mask for IP in IPPool, e.g. IP/MASK: 10.10.10.2/23, only for underlay scenario
	Mask int64 `json:"mask,omitempty" protobuf:"varint,10,opt,name=mask"`
	// Vlan for IP in IPPool, only for underlay scenario
	Vlan int64 `json:"vlan,omitempty" protobuf:"varint,11,opt,name=vlan"`
	// AdminAssign is a map with owner as key, with a string of IPs and IP-ranges with comma as seperator as value.
	// IP-range likes: 10.0.0.10-20 which contains 11 IPs, and only d part of IP(a.b.c.d) is
	// allowed to use for IP-range, which means one like 10.0.0.10-1.20 is invalid.
	// This allows admin to assign owner of IPs in .Status.AllocatedIPs, like to
	// - do pre-allocate IPs for ipset or ipclaim;
	// - mark IPs as "EXCLUDED" or revert;
	// - or anything to fix in .Status.AllocatedIPs.
	AdminAssign map[string]string `json:"adminAssign,omitempty" protobuf:"bytes,12,rep,name=adminAssign"`
	// TORID is logical ID of tor which ippool associated with. Nodes under tor must have label
	// "alcor.io/saishang-ipam.torID" to indicate the torID. // lv: TODO 可用做subnetID，用于实现vpc下隔离子网
	TORID string `json:"torID,omitempty" protobuf:"bytes,13,opt,name=torID"`

	// just for backward compatibility reason, use fisrtIP in new version
	StartIP string `json:"startIP,omitempty" protobuf:"bytes,14,opt,name=startIP"`
	// just for backward compatibility reason, use lastIP in new version
	EndIP string `json:"endIP,omitempty" protobuf:"bytes,15,opt,name=endIP"`
	// just for backward compatibility reason, no longer used
	DisabledIP []string `json:"disabledIP,omitempty" protobuf:"bytes,16,rep,name=disabledIP"`
	// just for backward compatibility reason, no longer used
	Backend string `json:"backend,omitempty" protobuf:"bytes,17,opt,name=backend"`
}

// IPPoolStatus defines the observed state of IPPool
type IPPoolStatus struct {
	// IPCount is number of total available IPs in pool.
	IPCount int64 `json:"IPCount,omitempty" protobuf:"varint,1,opt,name=IPCount,json=iPCount"`
	// AvailableCount is number of current available IPs in pool.
	AvailableCount int64 `json:"AvailableCount,omitempty" protobuf:"varint,2,opt,name=AvailableCount,json=availableCount"`
	// AllocatedIPs is map of allocated IPs.
	AllocatedIPs map[string]string `json:"Allocated,omitempty" protobuf:"bytes,3,rep,name=Allocated,json=allocated"`
	// Status is status of IPPool.
	Status string `json:"Status,omitempty" protobuf:"bytes,4,opt,name=Status,json=status"`
	// LastApplied is what applied in last time to IPPool.
	LastApplied map[string]string `json:"LastApplied,omitempty" protobuf:"bytes,5,rep,name=LastApplied,json=lastApplied"`
	// Message is operator processing message about IPPool.
	Message             string            `json:"Message,omitempty" protobuf:"bytes,6,opt,name=Message,json=message"`
	BlockedIPs          map[string]string `json:"BlockedIPs,omitempty" protobuf:"bytes,7,rep,name=BlockedIPs,json=blockedIPs"`
	SpecialClaimedUsage map[string]int64  `json:"SpecialClaimedUsage,omitempty" protobuf:"bytes,8,rep,name=SpecialClaimedUsage,json=specialClaimedUsage"`
	// Whether this IPPool has migrated from old one.
	Migrated bool `json:"Migrated,omitempty" protobuf:"varint,9,opt,name=Migrated,json=migrated"`

	// just for backward compatibility reason, no longer used
	OldAllocatedIPs map[string]string `json:"oldAllocatedIPs,omitempty" protobuf:"bytes,10,rep,name=oldAllocatedIPs"`
	// just for backward compatibility reason, no longer used
	OldAvailableCount int64 `json:"oldAvailableCount,omitempty" protobuf:"varint,11,opt,name=oldAvailableCount"`
	// just for backward compatibility reason, no longer used
	OldIPCount int64 `json:"oldIPCount,omitempty" protobuf:"varint,12,opt,name=oldIPCount"`
	// just for backward compatibility reason, no longer used
	OldStatus string `json:"oldStatus,omitempty" protobuf:"bytes,13,opt,name=oldStatus"`
	// just for backward compatibility reason, no longer used
	OldLastAppliedStartIP string `json:"lastAppliedStartIP,omitempty" protobuf:"bytes,14,opt,name=lastAppliedStartIP"`
	// just for backward compatibility reason, no longer used
	OldLastAppliedEndIP string `json:"lastAppliedEndIP,omitempty" protobuf:"bytes,15,opt,name=lastAppliedEndIP"`
}

// +genclient
// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type IPPool struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	Spec   IPPoolSpec   `json:"spec,omitempty" protobuf:"bytes,2,opt,name=spec"`
	Status IPPoolStatus `json:"status,omitempty" protobuf:"bytes,3,opt,name=status"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// IPPoolList contains a list of IPPool
type IPPoolList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`
	Items           []IPPool `json:"items" protobuf:"bytes,2,rep,name=items"`
}
