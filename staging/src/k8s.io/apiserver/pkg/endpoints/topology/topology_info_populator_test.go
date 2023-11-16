package topology

import (
	"k8s.io/api/mizarrbac/v1alpha1"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/apiserver/pkg/authentication/user"
	request "k8s.io/apiserver/pkg/endpoints/request"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes/fake"
	"net/url"
	"reflect"
	"time"

	"net/http"
	"testing"
)

const TestNs = "test-ns"

type StaticRoles struct {
	globalRoles []*v1alpha1.GlobalRole
}

type StaticRoleBindings struct {
	globalRoleBindings []*v1alpha1.GlobalRoleBinding
}

var staticRoles = StaticRoles{
	globalRoles: []*v1alpha1.GlobalRole{
		{
			ObjectMeta: metav1.ObjectMeta{
				Name: "global-shanghai-admin",
			},
			ClusterSelectors: []metav1.LabelSelector{
				{
					MatchLabels: map[string]string{
						"topology.mizar-k8s.io/region":      "shanghai",
						"topology.mizar-k8s.io/clusterType": "private",
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
				Name: "global-beijing-admin",
			},
			ClusterSelectors: []metav1.LabelSelector{
				{
					MatchLabels: map[string]string{
						"topology.mizar-k8s.io/region":      "beijing",
						"topology.mizar-k8s.io/clusterType": "private",
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
				Name: "global-admin",
			},
			RegionSelectors: []metav1.LabelSelector{
				{
					MatchLabels: map[string]string{
						"app": "example",
					},
				},
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
}

var simpleTopologyGetter = &SimpleTopologyGetter{
	Regions: map[string]metav1.Object{
		"beijing": &unstructured.Unstructured{
			Object: map[string]interface{}{
				"apiVersion": "core.plan.mizar-universe.io/v1alpha1",
				"kind":       "Region",
				"metadata": map[string]interface{}{
					"name": "beijing",
					"labels": map[string]interface{}{
						"topology.mizar-k8s.io/region": "beijing",
						"app":                          "example",
					},
				},
			},
		},
		"shanghai": &unstructured.Unstructured{
			Object: map[string]interface{}{
				"apiVersion": "core.plan.mizar-universe.io/v1alpha1",
				"kind":       "Region",
				"metadata": map[string]interface{}{
					"name": "shanghai",
					"labels": map[string]interface{}{
						"topology.mizar-k8s.io/region": "shanghai",
					},
				},
			},
		},
	},
	Clusters: map[string]metav1.Object{
		"shanghai": &unstructured.Unstructured{
			Object: map[string]interface{}{
				"apiVersion": "core.mizargalaxy.mizar-k8s.io/v1alpha1",
				"kind":       "Cluster",
				"metadata": map[string]interface{}{
					"name": "shanghai",
					"labels": map[string]interface{}{
						"topology.mizar-k8s.io/region":      "shanghai",
						"topology.mizar-k8s.io/cluster":     "shanghai",
						"topology.mizar-k8s.io/clusterType": "public",
					},
				},
			},
		},
		"shanghai.planet-cluster-0": &unstructured.Unstructured{
			Object: map[string]interface{}{
				"apiVersion": "core.mizargalaxy.mizar-k8s.io/v1alpha1",
				"kind":       "Cluster",
				"metadata": map[string]interface{}{
					"name": "shanghai.planet-cluster-0",
					"labels": map[string]interface{}{
						"topology.mizar-k8s.io/region":      "shanghai",
						"topology.mizar-k8s.io/cluster":     "planet-cluster-0",
						"topology.mizar-k8s.io/clusterType": "private",
					},
				},
			},
		},
		"shanghai.planet-cluster-1": &unstructured.Unstructured{
			Object: map[string]interface{}{
				"apiVersion": "core.mizargalaxy.mizar-k8s.io/v1alpha1",
				"kind":       "Cluster",
				"metadata": map[string]interface{}{
					"name": "shanghai.planet-cluster-1",
					"labels": map[string]interface{}{
						"topology.mizar-k8s.io/region":      "shanghai",
						"topology.mizar-k8s.io/cluster":     "planet-cluster-1",
						"topology.mizar-k8s.io/clusterType": "private",
					},
				},
			},
		},
		"beijing": &unstructured.Unstructured{
			Object: map[string]interface{}{
				"apiVersion": "core.mizargalaxy.mizar-k8s.io/v1alpha1",
				"kind":       "Cluster",
				"metadata": map[string]interface{}{
					"name": "beijing",
					"labels": map[string]interface{}{
						"topology.mizar-k8s.io/region":      "beijing",
						"topology.mizar-k8s.io/cluster":     "beijing",
						"topology.mizar-k8s.io/clusterType": "public",
					},
				},
			},
		},
		"beijing.planet-cluster-0": &unstructured.Unstructured{
			Object: map[string]interface{}{
				"apiVersion": "core.mizargalaxy.mizar-k8s.io/v1alpha1",
				"kind":       "Cluster",
				"metadata": map[string]interface{}{
					"name": "beijing.planet-cluster-0",
					"labels": map[string]interface{}{
						"topology.mizar-k8s.io/region":      "beijing",
						"topology.mizar-k8s.io/cluster":     "planet-cluster-0",
						"topology.mizar-k8s.io/clusterType": "public",
					},
				},
			},
		},
		"beijing.planet-cluster-1": &unstructured.Unstructured{
			Object: map[string]interface{}{
				"apiVersion": "core.mizargalaxy.mizar-k8s.io/v1alpha1",
				"kind":       "Cluster",
				"metadata": map[string]interface{}{
					"name": "beijing.planet-cluster-1",
					"labels": map[string]interface{}{
						"topology.mizar-k8s.io/region":      "beijing",
						"topology.mizar-k8s.io/cluster":     "planet-cluster-1",
						"topology.mizar-k8s.io/clusterType": "private",
					},
				},
			},
		},
	},
}

func TestPopulator_PopulateTopologyInfo(t *testing.T) {
	testcases := []struct {
		Description                  string
		Role                         StaticRoles
		RoleBinding                  StaticRoleBindings
		Request                      *http.Request
		User                         user.Info
		RequestInfo                  *request.RequestInfo
		TopologyInfo                 TopologyInfo
		ExpectedTopologyInfo         TopologyInfo
		ExpectedRequestLabelSelector string
	}{
		{
			Description: "user bind the role global-shanghai-admin",
			Role:        staticRoles,
			RoleBinding: StaticRoleBindings{
				globalRoleBindings: []*v1alpha1.GlobalRoleBinding{
					{
						ObjectMeta: metav1.ObjectMeta{
							Name: "shanghai-user01",
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
								Name:     "user01",
							},
						},
					},
				},
			},
			Request:      &http.Request{Method: http.MethodGet, URL: &url.URL{Scheme: "https", Host: "www.example.com"}},
			User:         &user.DefaultInfo{Name: "user01"},
			RequestInfo:  &request.RequestInfo{Verb: "list"},
			TopologyInfo: NewDefaultTopologyInfo(true),
			ExpectedTopologyInfo: &DefaultTopologyInfo{
				ClusterNames: sets.NewString("shanghai.planet-cluster-0", "shanghai.planet-cluster-1"),
				ClusterLabels: []map[string]string{
					{
						"topology.mizar-k8s.io/region":      "shanghai",
						"topology.mizar-k8s.io/cluster":     "planet-cluster-0",
						"topology.mizar-k8s.io/clusterType": "private",
					},
					{
						"topology.mizar-k8s.io/region":      "shanghai",
						"topology.mizar-k8s.io/cluster":     "planet-cluster-1",
						"topology.mizar-k8s.io/clusterType": "private",
					},
				},
				TopologyResource: true,
			},
			ExpectedRequestLabelSelector: "topology.mizar-k8s.io/cluster in (shanghai.planet-cluster-0,shanghai.planet-cluster-1)",
		},
	}

	for i, tc := range testcases {
		t.Run(tc.Description, func(t *testing.T) {
			mockPopulator, err := NewMockPopulator(&tc.Role, &tc.RoleBinding, simpleTopologyGetter)
			if err != nil {
				t.Errorf("failed to get populator: %s", err)
			}

			mockRequest := NewMockRequest(tc.Request, tc.User, tc.TopologyInfo, tc.RequestInfo)
			topologyInfo, err := mockPopulator.PopulateTopologyInfo(mockRequest)
			if err != nil {
				t.Fatalf("falied to get populated topology info: %s", err.Error())
			}

			if !topologyInfoEqual(topologyInfo, tc.ExpectedTopologyInfo) {
				t.Errorf("case %d: expeced: %v\n got: %v", i, tc.ExpectedTopologyInfo, topologyInfo)
			}

			if mockRequest.FormValue("labelSelector") != tc.ExpectedRequestLabelSelector {
				t.Errorf("case %d: expeced: %s, got: %s", i, tc.ExpectedRequestLabelSelector, mockRequest.FormValue("labelSelector"))
			}

		})
	}

}

func topologyInfoEqual(a, b TopologyInfo) bool {
	if a.IsTopologyResource() != b.IsTopologyResource() {
		return false
	}
	if !a.GetRegionName().HasAll(b.GetRegionName().List()...) {
		return false
	}
	if !a.GetClusterName().HasAll(b.GetClusterName().List()...) {
		return false
	}

	if len(a.GetRegionLabels()) != len(b.GetRegionLabels()) {
		return false
	}
	for _, al := range a.GetRegionLabels() {
		bls := b.GetRegionLabels()
		if !hasLabel(bls, al) {
			return false
		}
	}

	if len(a.GetClusterLabels()) != len(b.GetClusterLabels()) {
		return false
	}
	for _, al := range a.GetClusterLabels() {
		bls := b.GetClusterLabels()
		if !hasLabel(bls, al) {
			return false
		}
	}

	return true
}

func hasLabel(set []map[string]string, needle map[string]string) bool {
	for _, item := range set {
		if reflect.DeepEqual(item, needle) {
			return true
		}
	}
	return false
}

func NewMockPopulator(staticRoles *StaticRoles, roleBindings *StaticRoleBindings, topologyGetter TopologyGetter) (Populator, error) {
	informerFactory := informers.NewSharedInformerFactory(fake.NewSimpleClientset(), time.Second)
	stop := make(chan struct{})
	informerFactory.Start(stop)
	informerFactory.WaitForCacheSync(stop)

	for _, globalRole := range staticRoles.globalRoles {
		err := informerFactory.Mizarrbac().V1alpha1().GlobalRoles().Informer().GetIndexer().Add(globalRole)
		if err != nil {
			return nil, err
		}
	}

	for _, roleBinding := range roleBindings.globalRoleBindings {
		err := informerFactory.Mizarrbac().V1alpha1().GlobalRoleBindings().Informer().GetIndexer().Add(roleBinding)
		if err != nil {
			return nil, err
		}
	}

	return NewTopologyInfoPopulator(topologyGetter,
		informerFactory.Mizarrbac().V1alpha1().GlobalRoles().Lister(),
		informerFactory.Mizarrbac().V1alpha1().GlobalRoleBindings().Lister()), nil
}

func NewMockRequest(req *http.Request, user user.Info, topologyInfo TopologyInfo, requestInfo *request.RequestInfo) *http.Request {
	req = req.WithContext(request.WithUser(req.Context(), user))
	req = req.WithContext(WithTopologyInfo(req.Context(), topologyInfo))
	req = req.WithContext(request.WithRequestInfo(req.Context(), requestInfo))
	return req
}
