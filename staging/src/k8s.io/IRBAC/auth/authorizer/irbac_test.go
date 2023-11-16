package authorizer

import (
	"context"
	"k8s.io/apimachinery/pkg/util/sets"

	"testing"
	"time"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/api/mizarrbac/v1alpha1"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apiserver/pkg/authentication/user"
	"k8s.io/apiserver/pkg/authorization/authorizer"
	"k8s.io/apiserver/pkg/endpoints/topology"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes/fake"
)

const TestNs = "test-ns"

type StaticRoles struct {
	roles        []*rbacv1.Role
	clusterRoles []*rbacv1.ClusterRole
	globalRoles  []*v1alpha1.GlobalRole
	namespaces   []*corev1.Namespace
}

type staticRoleBindings struct {
	roleBinding        *rbacv1.RoleBinding
	clusterRoleBinding *rbacv1.ClusterRoleBinding
	globalRoleBinding  *v1alpha1.GlobalRoleBinding
}

var staticRoles = StaticRoles{
	roles: []*rbacv1.Role{
		{
			ObjectMeta: metav1.ObjectMeta{
				Namespace: TestNs,
				Name:      "admin",
			},
			Rules: []rbacv1.PolicyRule{
				{
					Verbs:     []string{"*"},
					APIGroups: []string{"*"},
					Resources: []string{"*"},
				},
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Namespace: TestNs,
				Name:      "viewer",
			},
			Rules: []rbacv1.PolicyRule{
				{
					Verbs:     []string{"get", "list", "watch"},
					APIGroups: []string{"*"},
					Resources: []string{"*"},
				},
			},
		},
	},

	clusterRoles: []*rbacv1.ClusterRole{
		{
			ObjectMeta: metav1.ObjectMeta{
				Name: "cluster-viewer",
			},
			Rules: []rbacv1.PolicyRule{
				{
					Verbs:     []string{"get", "list", "watch"},
					APIGroups: []string{"*"},
					Resources: []string{"*"},
				},
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name: "cluster-admin",
			},
			Rules: []rbacv1.PolicyRule{
				{
					Verbs:     []string{"*"},
					APIGroups: []string{"*"},
					Resources: []string{"*"},
				},
			},
		},
	},

	globalRoles: []*v1alpha1.GlobalRole{
		{
			ObjectMeta: metav1.ObjectMeta{
				Name: "global-shanghai-admin",
			},
			RegionSelectors: []metav1.LabelSelector{
				{
					MatchLabels: map[string]string{
						"app": "shanghai",
					},
				},
			},
			ClusterSelectors: []metav1.LabelSelector{
				{
					MatchLabels: map[string]string{
						"example": "bigdata",
					},
				},
			},
			Rules: []rbacv1.PolicyRule{
				{
					Verbs:           []string{"*"},
					APIGroups:       []string{"*"},
					Resources:       []string{"*"},
					NonResourceURLs: []string{"*"},
				},
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name: "global-shanghai-viewer",
			},
			RegionSelectors: []metav1.LabelSelector{
				{
					MatchLabels: map[string]string{
						"app": "shanghai",
					},
				},
			},
			ClusterSelectors: []metav1.LabelSelector{
				{
					MatchLabels: map[string]string{
						"example": "bigdata",
					},
				},
			},
			Rules: []rbacv1.PolicyRule{
				{
					Verbs:           []string{"get", "list", "watch"},
					APIGroups:       []string{"*"},
					Resources:       []string{"*"},
					NonResourceURLs: []string{"*"},
				},
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name: "global-admin",
			},
			Rules: []rbacv1.PolicyRule{
				{
					Verbs:     []string{"*"},
					APIGroups: []string{"*"},
					Resources: []string{"*"},
				},
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name: "global-viewer",
			},
			Rules: []rbacv1.PolicyRule{
				{
					Verbs:     []string{"get", "list", "watch"},
					APIGroups: []string{"*"},
					Resources: []string{"*"},
				},
			},
		},
	},

	namespaces: []*corev1.Namespace{
		{
			ObjectMeta: metav1.ObjectMeta{
				Name: TestNs,
			},
		},
	},
}

func TestIRBACAuthorizer_Authorize(t *testing.T) {

	testData := []struct {
		StaticRoles
		RoleBindings     staticRoleBindings
		Description      string
		Ctx              context.Context
		Request          authorizer.AttributesRecord
		ExpectedDecision authorizer.Decision
	}{
		{
			Description: "allow admin to create workloads by globalRole global-shanghai-admin",
			StaticRoles: staticRoles,
			RoleBindings: staticRoleBindings{
				globalRoleBinding: &v1alpha1.GlobalRoleBinding{
					ObjectMeta: metav1.ObjectMeta{
						Name: "admin",
					},
					RoleRef: rbacv1.RoleRef{
						APIGroup: v1alpha1.SchemeGroupVersion.Group,
						Kind:     "GlobalRole",
						Name:     "global-shanghai-admin",
					},
					Subjects: []rbacv1.Subject{
						{
							Kind:     rbacv1.UserKind,
							APIGroup: rbacv1.GroupName,
							Name:     "admin",
						},
					},
				},
			},
			Ctx: topology.WithTopologyInfo(context.Background(), &topology.DefaultTopologyInfo{
				RegionLabels:     []map[string]string{{"app": "shanghai"}},
				ClusterLabels:    []map[string]string{{"example": "bigdata"}},
				RegionNames:      sets.NewString("example-region"),
				ClusterNames:     sets.NewString("example-cluster"),
				TopologyResource: true,
			}),
			Request: authorizer.AttributesRecord{
				User: &user.DefaultInfo{
					Name: "admin",
				},
				Verb:            "create",
				Namespace:       TestNs,
				APIGroup:        "apps.mizargalaxy.mizar-k8s.io",
				APIVersion:      "v1alpha1",
				Resource:        "workloads",
				ResourceRequest: true,
			},
			ExpectedDecision: authorizer.DecisionAllow,
		},

		{
			Description: "user1 is not allowed to create workload due to topology mismatch",
			StaticRoles: staticRoles,
			RoleBindings: staticRoleBindings{
				globalRoleBinding: &v1alpha1.GlobalRoleBinding{
					ObjectMeta: metav1.ObjectMeta{
						Name: "user1",
					},
					RoleRef: rbacv1.RoleRef{
						APIGroup: v1alpha1.SchemeGroupVersion.Group,
						Kind:     "GlobalRole",
						Name:     "global-shanghai-admin",
					},
					Subjects: []rbacv1.Subject{
						{
							Kind:     rbacv1.UserKind,
							APIGroup: rbacv1.GroupName,
							Name:     "user1",
						},
					},
				},
			},
			Ctx: topology.WithTopologyInfo(context.Background(), &topology.DefaultTopologyInfo{
				RegionLabels:     []map[string]string{{"app": "beijing"}},
				RegionNames:      sets.NewString("beijing"),
				TopologyResource: true,
			}),
			Request: authorizer.AttributesRecord{
				User: &user.DefaultInfo{
					Name: "user1",
				},
				Verb:            "create",
				Namespace:       TestNs,
				APIGroup:        "apps.mizargalaxy.mizar-k8s.io",
				APIVersion:      "v1alpha1",
				Resource:        "workloads",
				ResourceRequest: true,
			},
			ExpectedDecision: authorizer.DecisionNoOpinion,
		},

		{
			Description: "user1 is allowed to create workload due to cluster match",
			StaticRoles: staticRoles,
			RoleBindings: staticRoleBindings{
				globalRoleBinding: &v1alpha1.GlobalRoleBinding{
					ObjectMeta: metav1.ObjectMeta{
						Name: "user1",
					},
					RoleRef: rbacv1.RoleRef{
						APIGroup: v1alpha1.SchemeGroupVersion.Group,
						Kind:     "GlobalRole",
						Name:     "global-shanghai-admin",
					},
					Subjects: []rbacv1.Subject{
						{
							Kind:     rbacv1.UserKind,
							APIGroup: rbacv1.GroupName,
							Name:     "user1",
						},
					},
				},
			},
			Ctx: topology.WithTopologyInfo(context.Background(), &topology.DefaultTopologyInfo{
				RegionLabels:     []map[string]string{{"app": "beijing"}},
				ClusterLabels:    []map[string]string{{"example": "bigdata"}},
				RegionNames:      sets.NewString("beijing"),
				ClusterNames:     sets.NewString("example-cluster"),
				TopologyResource: true,
			}),
			Request: authorizer.AttributesRecord{
				User: &user.DefaultInfo{
					Name: "user1",
				},
				Verb:            "create",
				Namespace:       TestNs,
				APIGroup:        "apps.mizargalaxy.mizar-k8s.io",
				APIVersion:      "v1alpha1",
				Resource:        "workloads",
				ResourceRequest: true,
			},
			ExpectedDecision: authorizer.DecisionAllow,
		},

		{
			Description: "viewer is not allowed to create workloads",
			StaticRoles: staticRoles,
			RoleBindings: staticRoleBindings{
				globalRoleBinding: &v1alpha1.GlobalRoleBinding{
					ObjectMeta: metav1.ObjectMeta{
						Name: "user1",
					},
					RoleRef: rbacv1.RoleRef{
						APIGroup: v1alpha1.SchemeGroupVersion.Group,
						Kind:     "GlobalRole",
						Name:     "global-shanghai-viewer",
					},
					Subjects: []rbacv1.Subject{
						{
							Kind:     rbacv1.UserKind,
							APIGroup: rbacv1.GroupName,
							Name:     "user1",
						},
					},
				},
			},
			Ctx: topology.WithTopologyInfo(context.Background(), &topology.DefaultTopologyInfo{
				RegionLabels:     []map[string]string{{"app": "shanghai"}},
				RegionNames:      sets.NewString("shanghai"),
				TopologyResource: true,
			}),
			Request: authorizer.AttributesRecord{
				User: &user.DefaultInfo{
					Name: "viewer",
				},
				Verb:            "create",
				Namespace:       TestNs,
				APIGroup:        "apps.mizargalaxy.mizar-k8s.io",
				APIVersion:      "v1alpha1",
				Resource:        "workloads",
				ResourceRequest: true,
			},
			ExpectedDecision: authorizer.DecisionNoOpinion,
		},

		{
			Description: "admin is allowed to create deployment",
			StaticRoles: staticRoles,
			RoleBindings: staticRoleBindings{
				globalRoleBinding: &v1alpha1.GlobalRoleBinding{
					ObjectMeta: metav1.ObjectMeta{
						Name: "admin",
					},
					RoleRef: rbacv1.RoleRef{
						APIGroup: v1alpha1.SchemeGroupVersion.Group,
						Kind:     "GlobalRole",
						Name:     "global-admin",
					},
					Subjects: []rbacv1.Subject{
						{
							Kind:     rbacv1.UserKind,
							APIGroup: rbacv1.GroupName,
							Name:     "admin",
						},
					},
				},
			},
			Ctx: topology.WithTopologyInfo(context.Background(), &topology.DefaultTopologyInfo{
				TopologyResource: false,
			}),
			Request: authorizer.AttributesRecord{
				User:            &user.DefaultInfo{Name: "admin"},
				Verb:            "create",
				Namespace:       TestNs,
				APIGroup:        "apps",
				APIVersion:      "v1",
				Resource:        "deployments",
				ResourceRequest: true,
			},
			ExpectedDecision: authorizer.DecisionAllow,
		},

		{
			Description: "viewer is not allowed to create deployment",
			StaticRoles: staticRoles,
			RoleBindings: staticRoleBindings{
				globalRoleBinding: &v1alpha1.GlobalRoleBinding{
					ObjectMeta: metav1.ObjectMeta{
						Name: "viewer",
					},
					RoleRef: rbacv1.RoleRef{
						APIGroup: v1alpha1.SchemeGroupVersion.Group,
						Kind:     "GlobalRole",
						Name:     "global-viewer",
					},
					Subjects: []rbacv1.Subject{
						{
							Kind:     rbacv1.UserKind,
							APIGroup: rbacv1.GroupName,
							Name:     "viewer",
						},
					},
				},
			},
			Ctx: topology.WithTopologyInfo(context.Background(), &topology.DefaultTopologyInfo{
				TopologyResource: false,
			}),
			Request: authorizer.AttributesRecord{
				User:            &user.DefaultInfo{Name: "viewer"},
				Verb:            "create",
				Namespace:       TestNs,
				APIGroup:        "apps",
				APIVersion:      "v1",
				Resource:        "deployments",
				ResourceRequest: true,
			},
			ExpectedDecision: authorizer.DecisionNoOpinion,
		},

		{
			Description: "cluster viewer is not allowed to create deployment",
			StaticRoles: staticRoles,
			RoleBindings: staticRoleBindings{
				clusterRoleBinding: &rbacv1.ClusterRoleBinding{
					ObjectMeta: metav1.ObjectMeta{
						Name: "viewer",
					},
					RoleRef: rbacv1.RoleRef{
						APIGroup: rbacv1.SchemeGroupVersion.Group,
						Kind:     "ClusterRole",
						Name:     "cluster-viewer",
					},
					Subjects: []rbacv1.Subject{
						{
							Kind:     rbacv1.UserKind,
							APIGroup: rbacv1.GroupName,
							Name:     "cluster-viewer",
						},
					},
				},
			},
			Ctx: topology.WithTopologyInfo(context.Background(), &topology.DefaultTopologyInfo{
				TopologyResource: false,
			}),
			Request: authorizer.AttributesRecord{
				User:            &user.DefaultInfo{Name: "cluster-viewer"},
				Verb:            "create",
				Namespace:       TestNs,
				APIGroup:        "apps",
				APIVersion:      "v1",
				Resource:        "deployments",
				ResourceRequest: true,
			},
			ExpectedDecision: authorizer.DecisionNoOpinion,
		},

		{
			Description: "sa is allowed to create deployment",
			StaticRoles: staticRoles,
			RoleBindings: staticRoleBindings{
				clusterRoleBinding: &rbacv1.ClusterRoleBinding{
					ObjectMeta: metav1.ObjectMeta{
						Name: "example-sa",
					},
					RoleRef: rbacv1.RoleRef{
						APIGroup: rbacv1.SchemeGroupVersion.Group,
						Kind:     "ClusterRole",
						Name:     "cluster-admin",
					},
					Subjects: []rbacv1.Subject{
						{
							Kind:      rbacv1.ServiceAccountKind,
							APIGroup:  rbacv1.GroupName,
							Namespace: TestNs,
							Name:      "example-sa",
						},
					},
				},
			},
			Ctx: topology.WithTopologyInfo(context.Background(), &topology.DefaultTopologyInfo{
				TopologyResource: false,
			}),
			Request: authorizer.AttributesRecord{
				User:            &user.DefaultInfo{Name: "system:serviceaccount:test-ns:example-sa"},
				Verb:            "create",
				Namespace:       TestNs,
				APIGroup:        "apps",
				APIVersion:      "v1",
				Resource:        "deployments",
				ResourceRequest: true,
			},
			ExpectedDecision: authorizer.DecisionAllow,
		},

		{
			Description: "namespace admin is allowed to create deployment",
			StaticRoles: staticRoles,
			RoleBindings: staticRoleBindings{
				clusterRoleBinding: &rbacv1.ClusterRoleBinding{
					ObjectMeta: metav1.ObjectMeta{
						Name: "user01",
					},
					RoleRef: rbacv1.RoleRef{
						APIGroup: rbacv1.SchemeGroupVersion.Group,
						Kind:     "ClusterRole",
						Name:     "cluster-viewer",
					},
					Subjects: []rbacv1.Subject{
						{
							Kind:     rbacv1.UserKind,
							APIGroup: rbacv1.GroupName,
							Name:     "user01",
						},
					},
				},
				roleBinding: &rbacv1.RoleBinding{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "user01",
						Namespace: TestNs,
					},
					RoleRef: rbacv1.RoleRef{
						APIGroup: rbacv1.SchemeGroupVersion.Group,
						Kind:     "Role",
						Name:     "admin",
					},
					Subjects: []rbacv1.Subject{
						{
							Kind:     rbacv1.UserKind,
							APIGroup: rbacv1.GroupName,
							Name:     "user01",
						},
					},
				},
			},
			Ctx: topology.WithTopologyInfo(context.Background(), &topology.DefaultTopologyInfo{
				TopologyResource: false,
			}),
			Request: authorizer.AttributesRecord{
				User:            &user.DefaultInfo{Name: "user01"},
				Verb:            "create",
				Namespace:       TestNs,
				APIGroup:        "apps",
				APIVersion:      "v1",
				Resource:        "deployments",
				ResourceRequest: true,
			},
			ExpectedDecision: authorizer.DecisionAllow,
		},
		{
			Description: "namespace sa is allowed to create deployment",
			StaticRoles: staticRoles,
			RoleBindings: staticRoleBindings{
				clusterRoleBinding: &rbacv1.ClusterRoleBinding{
					ObjectMeta: metav1.ObjectMeta{
						Name: "example-sa",
					},
					RoleRef: rbacv1.RoleRef{
						APIGroup: rbacv1.SchemeGroupVersion.Group,
						Kind:     "ClusterRole",
						Name:     "cluster-viewer",
					},
					Subjects: []rbacv1.Subject{
						{
							Kind:     rbacv1.UserKind,
							APIGroup: rbacv1.GroupName,
							Name:     "user01",
						},
					},
				},
				roleBinding: &rbacv1.RoleBinding{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "example-sa",
						Namespace: TestNs,
					},
					RoleRef: rbacv1.RoleRef{
						APIGroup: rbacv1.SchemeGroupVersion.Group,
						Kind:     "Role",
						Name:     "admin",
					},
					Subjects: []rbacv1.Subject{
						{
							Kind:      rbacv1.ServiceAccountKind,
							APIGroup:  rbacv1.GroupName,
							Name:      "example-sa",
							Namespace: TestNs,
						},
					},
				},
			},
			Ctx: topology.WithTopologyInfo(context.Background(), &topology.DefaultTopologyInfo{
				TopologyResource: false,
			}),
			Request: authorizer.AttributesRecord{
				User:            &user.DefaultInfo{Name: "system:serviceaccount:test-ns:example-sa"},
				Verb:            "create",
				Namespace:       TestNs,
				APIGroup:        "apps",
				APIVersion:      "v1",
				Resource:        "deployments",
				ResourceRequest: true,
			},
			ExpectedDecision: authorizer.DecisionAllow,
		},

		{
			Description: "user01 is allowed to create workloads by cluster admin",
			StaticRoles: staticRoles,
			RoleBindings: staticRoleBindings{
				clusterRoleBinding: &rbacv1.ClusterRoleBinding{
					ObjectMeta: metav1.ObjectMeta{
						Name: "user01",
					},
					RoleRef: rbacv1.RoleRef{
						APIGroup: rbacv1.SchemeGroupVersion.Group,
						Kind:     "ClusterRole",
						Name:     "cluster-admin",
					},
					Subjects: []rbacv1.Subject{
						{
							Kind:     rbacv1.UserKind,
							APIGroup: rbacv1.GroupName,
							Name:     "user01",
						},
					},
				},
			},
			Ctx: topology.WithTopologyInfo(context.Background(), &topology.DefaultTopologyInfo{
				TopologyResource: true,
			}),
			Request: authorizer.AttributesRecord{
				User:            &user.DefaultInfo{Name: "user01"},
				Verb:            "create",
				Namespace:       TestNs,
				APIGroup:        "apps.mizargalaxy.mizar-k8s.io",
				APIVersion:      "v1alpha1",
				Resource:        "workloads",
				ResourceRequest: true,
			},
			ExpectedDecision: authorizer.DecisionAllow,
		},
	}

	for i, tc := range testData {
		t.Run(tc.Description, func(t *testing.T) {

			authorize, err := MockClient(&tc.StaticRoles, tc.RoleBindings)
			if err != nil {
				t.Errorf("failed to create mock client, error: %s", err.Error())

			}
			decision, message, err := authorize.Authorize(tc.Ctx, &tc.Request)

			if err != nil {
				t.Errorf("case %d: %v: %s", i, err, message)
			}

			if decision != tc.ExpectedDecision {
				t.Errorf("case %d: %d != %d", i, decision, tc.ExpectedDecision)
			}
		})
	}
}

func MockClient(staticRoles *StaticRoles, roleBindings staticRoleBindings) (authorizer.Authorizer, error) {
	informerFactory := informers.NewSharedInformerFactory(fake.NewSimpleClientset(), time.Second)
	stop := make(chan struct{})
	informerFactory.Start(stop)
	informerFactory.WaitForCacheSync(stop)

	for _, role := range staticRoles.roles {
		err := informerFactory.Rbac().V1().Roles().Informer().GetIndexer().Add(role)
		if err != nil {
			return nil, err
		}
	}

	for _, role := range staticRoles.clusterRoles {
		err := informerFactory.Rbac().V1().ClusterRoles().Informer().GetIndexer().Add(role)
		if err != nil {
			return nil, err
		}
	}

	for _, globalRole := range staticRoles.globalRoles {
		err := informerFactory.Mizarrbac().V1alpha1().GlobalRoles().Informer().GetIndexer().Add(globalRole)
		if err != nil {
			return nil, err
		}
	}

	for _, ns := range staticRoles.namespaces {
		err := informerFactory.Core().V1().Namespaces().Informer().GetIndexer().Add(ns)
		if err != nil {
			return nil, err
		}
	}

	if roleBindings.roleBinding != nil {
		err := informerFactory.Rbac().V1().RoleBindings().Informer().GetIndexer().Add(roleBindings.roleBinding)
		if err != nil {
			return nil, err
		}
	}

	if roleBindings.clusterRoleBinding != nil {
		err := informerFactory.Rbac().V1().ClusterRoleBindings().Informer().GetIndexer().Add(roleBindings.clusterRoleBinding)
		if err != nil {
			return nil, err
		}
	}

	if roleBindings.globalRoleBinding != nil {
		err := informerFactory.Mizarrbac().V1alpha1().GlobalRoleBindings().Informer().GetIndexer().Add(roleBindings.globalRoleBinding)
		if err != nil {
			return nil, err
		}
	}

	return NewIRBACAuthorizer(
		informerFactory.Rbac().V1().Roles().Lister(),
		informerFactory.Rbac().V1().RoleBindings().Lister(),
		informerFactory.Rbac().V1().ClusterRoles().Lister(),
		informerFactory.Rbac().V1().ClusterRoleBindings().Lister(),
		informerFactory.Mizarrbac().V1alpha1().GlobalRoles().Lister(),
		informerFactory.Mizarrbac().V1alpha1().GlobalRoleBindings().Lister(),
	), nil

}
