package v1alpha1

import (
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +genclient:noStatus
// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// GlobalRole is a global level, logical grouping of PolicyRules
type GlobalRole struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,6,opt,name=metadata"`

	// RegionSelectors is a selector list of regions this rule applies to.
	// +optional
	RegionSelectors []metav1.LabelSelector `json:"regionSelectors,omitempty" protobuf:"bytes,2,rep,name=regionSelectors"`

	// ClusterSelectors is a selector list of clusters this rule applies to.
	// +optional
	ClusterSelectors []metav1.LabelSelector `json:"clusterSelectors,omitempty" protobuf:"bytes,3,rep,name=clusterSelectors"`

	// Rules holds all the PolicyRules for this ClusterRole
	// +optional
	Rules []rbacv1.PolicyRule `json:"rules,omitempty" protobuf:"bytes,4,rep,name=rules"`

	// AggregationRule is an optional field that describes how to build the Rules for this GlobalRole.
	// If AggregationRule is set, then the Rules are controller managed and direct changes to Rules will be
	// stomped by the controller.
	// +optional
	AggregationRule *AggregationRule `json:"aggregationRule,omitempty" protobuf:"bytes,5,opt,name=aggregationRule"`
}

// AggregationRule describes how to locate GlobalRoles to aggregate into the GlobalRole
type AggregationRule struct {
	// GlobalRoleSelectors holds a list of selectors which will be used to find GlobalRoles and create the rules.
	// If any of the selectors match, then the GlobalRole's permissions will be added
	// +optional
	GlobalRoleSelectors []metav1.LabelSelector `json:"globalRoleSelectors,omitempty" protobuf:"bytes,1,rep,name=globalRoleSelectors"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// GlobalRoleList contains a list of GlobalRole
type GlobalRoleList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`
	Items           []GlobalRole `json:"items" protobuf:"bytes,2,rep,name=items"`
}

// +genclient
// +genclient:noStatus
// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// GlobalRoleBinding references a GlobalRole, but not contain it.  It can reference a GlobalRole in the global namespace,
// and adds who information via Subject.
type GlobalRoleBinding struct {
	metav1.TypeMeta `json:",inline"`
	// Standard object's metadata.
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	// Subjects holds references to the objects the role applies to.
	// +optional
	Subjects []rbacv1.Subject `json:"subjects,omitempty" protobuf:"bytes,2,rep,name=subjects"`

	// RoleRef can only reference a ClusterRole in the global namespace.
	// If the RoleRef cannot be resolved, the Authorizer must return an error.
	RoleRef rbacv1.RoleRef `json:"roleRef" protobuf:"bytes,3,opt,name=roleRef"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// GlobalRoleBindingList contains a list of GlobalRoleBinding
type GlobalRoleBindingList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`
	Items           []GlobalRoleBinding `json:"items" protobuf:"bytes,2,rep,name=items"`
}
