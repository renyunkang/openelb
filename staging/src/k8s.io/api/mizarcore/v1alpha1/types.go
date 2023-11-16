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
// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Cluster is the Schema for the clusters API.
type Cluster struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	Spec   ClusterSpec   `json:"spec,omitempty" protobuf:"bytes,2,opt,name=spec"`
	Status ClusterStatus `json:"status,omitempty" protobuf:"bytes,3,opt,name=status"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ClusterList contains a list of Cluster.
type ClusterList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`
	Items           []Cluster `json:"items" protobuf:"bytes,2,rep,name=items"`
}

// ClusterSpec defines the desired state of Cluster.
type ClusterSpec struct {
	Taints []v1.Taint `json:"taints,omitempty" protobuf:"bytes,2,rep,name=taints"`

	// Addons describe the configuration of Cluster will be controlled by MizarGalaxy.
	Addons *Addons `json:"addons" protobuf:"bytes,1,opt,name=addons"`
}

type Addons struct {
	// InitializeWorkloadCluster descibes which plugin to establish a WorkloadCluster and configuration of the cluster
	// MizarGalaxy provides a default plugin, cluster-api-quick-start, to establish a WorkloadCluster.
	// +optional
	InitializeWorkloadCluster *InitializeWorkloadCluster `json:"initializeWorkloadCluster,omitempty" protobuf:"bytes,1,opt,name=initializeWorkloadCluster"`

	// PreRegisterWorkloadCluster describes the process before register WorkloadCluster to MizarGalax.
	// More specificly, it describes the configuration of kubelet-for-cluster deployed in the specific WorkloadCluster.
	// +optional
	PreRegisterWorkloadCluster *PreRegisterWorkloadCluster `json:"preRegisterWorkloadCluster" protobuf:"bytes,2,opt,name=preRegisterWorkloadCluster"`
}

type InitializeWorkloadCluster struct {
	// Name is the name of the plugin for creating a WorkloadCluster.
	// +required
	Name string `json:"name" protobuf:"bytes,1,opt,name=name"`

	// InitConfiguration descibes configuration of the WorkloadCluster that you want to establish.
	// +required
	InitConfiguration *InitConfiguration `json:"configurations" protobuf:"bytes,2,opt,name=configurations"`
}

type InitConfiguration struct {
	// KubernetesVersion represents version of the workloadCluster.
	// +optional
	KubernetesVersion string `json:"kubernetesVersion,omitempty" protobuf:"bytes,1,opt,name=kubernetesVersion"`

	// ControlPlaneMachineCount represents the numbers of controlPlaneMachine.
	// +optional
	ControlPlaneMachineCount int32 `json:"controlPlaneMachineCount,omitempty" protobuf:"varint,2,opt,name=controlPlaneMachineCount"`

	// WorkerMachineCount represents the numbers of WorkerMachine.
	// +optional
	WorkerMachineCount int32 `json:"workerMachineCount,omitempty" protobuf:"varint,3,opt,name=workerMachineCount"`
}

type PreRegisterWorkloadCluster struct {
	// Name is the name of the plugin for deploying kubelet-for-cluster.
	// +required
	Name string `json:"name" protobuf:"bytes,1,opt,name=name"`

	// Configurations descibes configrutions of kubelet-for-cluster
	// and which WhichCluster to deploy.
	// +required
	Configurations *KubeletForClusterConfiguration `json:"configurations" protobuf:"bytes,2,opt,name=configurations"`
}

type KubeletForClusterConfiguration struct {
	// Auth describes authentication of specific WorkloadCluster.
	// If selected initializeWorkloadCluster, this filed can be ignored.
	// +optional
	Auth *Auth `json:"auth,omitempty" protobuf:"bytes,1,opt,name=auth"`

	// Target represents which WorkloadCluster kubelet-for-cluster will be deployed on.
	// If selected initializeWorkloadCluster, this filed can be ignored.
	// +optional
	Target *Target `json:"target,omitempty" protobuf:"bytes,2,opt,name=target"`

	// Image represents which image to create kubelet-for-cluster.
	// +optional
	Image string `json:"image,omitempty" protobuf:"bytes,3,opt,name=image"`

	// Template represents configuraton template of kubelet-for-cluster,
	// Parameters that are not specified will be filled by the configuration in the template.
	// +optional
	Template string `json:"template,omitempty" protobuf:"bytes,4,opt,name=template"`

	// ReservedCPU represents WorkloadCluster reserved CPU.
	// +optional
	ReservedCPU resource.Quantity `json:"reservedCPU,omitempty" protobuf:"bytes,5,opt,name=reservedCPU"`

	// ReservedCPU represents WorkloadCluster reserved Memory.
	// +optional
	ReservedMemory resource.Quantity `json:"reservedMemory,omitempty" protobuf:"bytes,6,opt,name=reservedMemory"`

	// TopNodeN represents cluster resource top node n.
	// +optional
	TopNodeN int32 `json:"topNodeN,omitempty" protobuf:"varint,7,opt,name=topNodeN"`

	// HeartbeatFrequency represents the heartbeat frequency with the Federation.
	// +optional
	HeartbeatFrequency int64 `json:"heartbeatFrequency,omitempty" protobuf:"varint,8,opt,name=heartbeatFrequency"`

	// LeaseDurationSeconds represents the lease period of the Workload cluster.
	// +optional
	LeaseDurationSeconds int64 `json:"leaseDurationSeconds,omitempty" protobuf:"varint,9,opt,name=leaseDurationSeconds"`

	// ForceSyncFrequency represents forced synchronous cycle with Federation
	// +optional
	ForceSyncFrequency int64 `json:"forceSyncFrequency,omitempty" protobuf:"varint,10,opt,name=forceSyncFrequency"`

	// NamespacePrefix represents mapping from Federated namespaces to WorkLoadClusters
	// +optional
	NamespacePrefix string `json:"namespacePrefix,omitempty" protobuf:"bytes,11,opt,name=namespacePrefix"`
}
type Auth struct {
	// Type defines type of authentication, "KubeConig" or "Token"
	// +required
	Type string `json:"type" protobuf:"bytes,1,opt,name=type"`

	// Kubeconfig represents the path of kubernetes configuration
	// +optional
	Kubeconfig string `json:"kubeconfig,omitempty" protobuf:"bytes,2,opt,name=kubeconfig"`

	// Token represents the token to access WorkloadCluster
	// +optional
	Token string `json:"token,omitempty" protobuf:"bytes,3,opt,name=token"`
}
type Target struct {
	// Apiserver represents the IP address and port of WorkloadCLuster apiserver
	// +optional
	Apiserver string `json:"apiserver,omitempty" protobuf:"bytes,1,opt,name=apiserver"`
}

// ClusterPhase describes the status of WorkloadCluster
type ClusterPhase string

const (
	// ONLINE represents kubelet-for-cluster is running.
	ONLINE ClusterPhase = "Online"
	// OFFLINE represents kubelet-for-cluster connection with manager cluster apiserver failed.
	OFFLINE ClusterPhase = "Offline"
	// PENDING represents Feduration does not receive first heartbeat from kubelet-for-cluster
	PENDING ClusterPhase = "Pending"
	// TERMINATED represents cluster object is deleting
	TERMINATED ClusterPhase = "Terminated"
	// NodeReady means kubelet is healthy and ready to accept pods.
	ClusterReady string = "Ready"
)

// ClusterStatus defines the observed state of Cluster
type ClusterStatus struct {
	// Configuration describes the configuration of kubelet-for-cluster
	// +optional
	Configuration *KubeletForClusterConfiguration `json:"configuration,omitempty" protobuf:"bytes,1,opt,name=configuration"`

	// Addresses shows workload cluster apiserver addresses
	// +optional
	Addresses []Address `json:"addresses,omitempty" protobuf:"bytes,2,rep,name=addresses"`

	// Phase represents the status of the WorkCluster.
	// +optional
	Phase ClusterPhase `json:"phase,omitempty" default:"pending" protobuf:"bytes,3,opt,name=phase,casttype=ClusterPhase"`

	// Allocatable represents the resources of a cluster that are available for scheduling.
	// +optional
	Allocatable v1.ResourceList `json:"allocatable,omitempty" protobuf:"bytes,4,opt,name=allocatable"`

	// Usage represents the resources of a cluster that has already used.
	// +optional
	Usage v1.ResourceList `json:"usage,omitempty" protobuf:"bytes,5,opt,name=usage"`

	// Nodes describes remained resource of top N nodes based on remaining resources
	// +optional
	Nodes []NodeLeftResource `json:"nodes,omitempty" protobuf:"bytes,6,rep,name=nodes"`

	// PlanetClusterAllocatable represents the resources of a cluster that are available for scheduling.
	// +optional
	PlanetClusterAllocatable v1.ResourceList `json:"planetClusterAllocatable,omitempty" protobuf:"bytes,7,opt,name=planetClusterAllocatable"`

	// PlanetClusterUsage represents the resources of a cluster that has already used.
	// +optional
	PlanetClusterUsage v1.ResourceList `json:"planetClusterUsage,omitempty" protobuf:"bytes,8,opt,name=planetClusterUsage"`

	// PlanetClusterNodes describes remained resource of top N nodes based on remaining resources
	// +optional
	PlanetClusterNodes []NodeLeftResource `json:"planetClusterNodes,omitempty" protobuf:"bytes,9,rep,name=planetClusterNodes"`

	// Namespaces describes resource occupation of a federal namespace in the WorkloadCluster.
	// +optional
	Namespaces []NamespaceUsage `json:"namespaces,omitempty" protobuf:"bytes,10,rep,name=namespaces"`

	// Conditions is an array of current cluster conditions.
	// +optional
	Condition []ClusterCondition `json:"conditions,omitempty" protobuf:"bytes,11,rep,name=conditions"`

	// ClusterInfo describes the specific information of WorkloadCluster
	// +optional
	ClusterInfo Info `json:"clusterInfo,omitempty" protobuf:"bytes,12,opt,name=clusterInfo"`

	// Endpoints of daemons running on the Cluster.
	// +optional
	DaemonEndpoints ClusterDaemonEndpoints `json:"daemonEndpoints,omitempty" protobuf:"bytes,13,opt,name=daemonEndpoints"`

	// Partitions describes which partitions cluster contains
	// +optional
	Partitions []string `json:"partitions,omitempty" protobuf:"bytes,14,rep,name=partitions"`

	// SecretRef of planet cluster
	// +optional
	SecretRef ClusterSecretRef `json:"secretRef,omitempty" protobuf:"bytes,15,opt,name=secretRef"`

	// Storage is an array of csi plugins installed in the planet cluster
	// +optional
	Storage []string `json:"storage,omitempty" protobuf:"bytes,16,rep,name=storage"`
}

const (
	PartitionLabel   = "mizar-k8s.io/partition"
	StarClusterSplit = "."

	ClusterLevelLabel         = "mizargalaxy.mizar-k8s.io/cluster-level"
	ClusterLevelStarCluster   = "1"
	ClusterLevelPlanetCluster = "2"

	ClusterProxyBlackAnnotation = "mizargalaxy.mizar-k8s.io/cluster-proxy-disabled"
)

// ClusterSecretRef contains information about a credential for accessing planet cluster
type ClusterSecretRef struct {
	// Name of secret
	// +optional
	Name string `json:"name,omitempty" protobuf:"bytes,1,opt,name=name"`

	// Namespace of secret
	// +optional
	Namespace string `json:"namespace,omitempty" protobuf:"bytes,2,opt,name=namespace"`
}

// DaemonEndpoint contains information about a single Daemon endpoint.
type DaemonEndpoint struct {

	// Port number of the given endpoint.
	// +required
	Port int32 `json:"port,omitempty" protobuf:"bytes,1,opt,name=port"`

	// Protocol of the given endpoint, http/https
	// +required
	Protocol string `json:"protocol,omitempty" protobuf:"bytes,2,opt,name=protocol"`

	// Address of the given endpoint
	// +required
	Address string `json:"address,omitempty" protobuf:"bytes,3,opt,name=address"`
}

// ClusterDaemonEndpoints lists ports opened by daemons running on the Node.
type ClusterDaemonEndpoints struct {
	// Endpoint on which KubeletForCluster is listening.
	// +optional
	KubeletForClusterEndpoint DaemonEndpoint `json:"kubeletForClusterEndpoint,omitempty" protobuf:"bytes,1,opt,name=kubeletForClusterEndpoint"`

	// Endpoint on which KubeApiServer is listening in planet cluster
	// +optional
	ApiServerEndpoint DaemonEndpoint `json:"apiServerEndpoint,omitempty" protobuf:"bytes,2,opt,name=apiServerEndpoint"`
}

// AddressType defines the type of Address, one of InternalIP, ExternalIP
type AddressType string

const (
	InternalIP AddressType = "InternalIP"
	ExternalIP AddressType = "ExternalIP"
)

type Info struct {
	Major        string `json:"major" protobuf:"bytes,1,opt,name=major"`
	Minor        string `json:"minor" protobuf:"bytes,2,opt,name=minor"`
	GitVersion   string `json:"gitVersion" protobuf:"bytes,3,opt,name=gitVersion"`
	GitCommit    string `json:"gitCommit" protobuf:"bytes,4,opt,name=gitCommit"`
	GitTreeState string `json:"gitTreeState" protobuf:"bytes,5,opt,name=gitTreeState"`
	BuildDate    string `json:"buildDate" protobuf:"bytes,6,opt,name=buildDate"`
	GoVersion    string `json:"goVersion" protobuf:"bytes,7,opt,name=goVersion"`
	Compiler     string `json:"compiler" protobuf:"bytes,8,opt,name=compiler"`
	Platform     string `json:"platform" protobuf:"bytes,9,opt,name=platform"`
}

// Addresses shows workload cluster apiserver addresses.
type Address struct {
	// AddressIP represents IP of the address
	// +required
	AddressIP string `json:"address" protobuf:"bytes,1,opt,name=address"`

	// Type represents the type of Address.
	// +required
	Type AddressType `json:"type" protobuf:"bytes,2,opt,name=type,casttype=AddressType"`
}

// NamespaceUsage describes requests and limits resource of a federal namespace in the WorkloadCluster.
type NamespaceUsage struct {
	// Name represents the name of namespace.
	// +required
	Name string `json:"name" protobuf:"bytes,1,opt,name=name"`

	// Usage describes the requests and limits resource of a federal namespace in the WorkloadCluster.
	// +required
	Usage v1.ResourceRequirements `json:"usage" protobuf:"bytes,2,opt,name=usage"`
}

// NodeLeftResource describes a node's remained resource.
type NodeLeftResource struct {
	// Name represents the name of the node.
	// +required
	Name string `json:"name" protobuf:"bytes,1,opt,name=name"`

	// Left describes remained resource of the node.
	// +required
	Left v1.ResourceList `json:"left" protobuf:"bytes,2,opt,name=left"`
}

type ConditionStatus string

const (
	TrueCondition    ConditionStatus = "True"
	FalseCondition   ConditionStatus = "False"
	UnknownCondition ConditionStatus = "Unknown"
)

type ClusterCondition struct {
	// Type of cluster condition.
	Type string `json:"type" protobuf:"bytes,1,opt,name=type"`

	// Status of the condition, one of True, False, Unknown.
	Status ConditionStatus `json:"status" protobuf:"bytes,2,opt,name=status,casttype=ConditionStatus"`

	// Last time we got an update on a given condition.
	// +optional
	LastHeartbeatTime metav1.Time `json:"lastHeartbeatTime,omitempty" protobuf:"bytes,3,opt,name=lastHeartbeatTime"`

	// Last time the condition transit from one status to another.
	// +optional
	LastTransitionTime metav1.Time `json:"lastTransitionTime,omitempty" protobuf:"bytes,4,opt,name=lastTransitionTime"`

	// (brief) reason for the condition's last transition.
	// +optional
	Reason string `json:"reason,omitempty" protobuf:"bytes,5,opt,name=reason"`

	// Human readable message indicating details about last transition.
	// +optional
	Message string `json:"message,omitempty" protobuf:"bytes,6,opt,name=message"`
}

const (
	// NamespaceClusterLease is the namespace where we place cluster lease objects (used for cluster heartbeats)
	NamespaceClusterLease string = "kube-cluster-lease"
)

// +k8s:conversion-gen:explicit-from=net/url.Values
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ClusterProxyOptions is the query options to a Node's proxy call
type ClusterProxyOptions struct {
	metav1.TypeMeta `json:",inline"`

	// Path is the URL path to use for the current proxy request
	Path string `json:"path" protobuf:"bytes,1,opt,name=path"`
}

// +genclient
// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Cluster is the Schema for the clusters API.
type ClusterSet struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	Clusters []string `json:"clusters,omitempty" protobuf:"bytes,2,opt,name=clusters"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ClusterList contains a list of Cluster.
type ClusterSetList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`
	Items           []ClusterSet `json:"items" protobuf:"bytes,2,rep,name=items"`
}
