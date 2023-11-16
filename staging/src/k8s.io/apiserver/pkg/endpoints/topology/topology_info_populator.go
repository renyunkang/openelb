package topology

import (
	"context"
	"fmt"
	"k8s.io/api/mizarrbac/v1alpha1"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apiserver/pkg/endpoints/request"
	mizarrbaclisters "k8s.io/client-go/listers/mizarrbac/v1alpha1"
	"net/http"
	"strings"
)

type Populator interface {
	PopulateTopologyInfo(req *http.Request) (TopologyInfo, error)
}

func NewTopologyInfoPopulator(getter TopologyGetter, globalRoleLister mizarrbaclisters.GlobalRoleLister,
	globalRoleBindingLister mizarrbaclisters.GlobalRoleBindingLister) Populator {
	return &populator{
		globalRoleBindingLister: globalRoleBindingLister,
		globalRoleLister:        globalRoleLister,
		TopologyGetter:          getter,
	}
}

type populator struct {
	globalRoleBindingLister mizarrbaclisters.GlobalRoleBindingLister
	globalRoleLister        mizarrbaclisters.GlobalRoleLister
	TopologyGetter          TopologyGetter
}

func (p *populator) PopulateTopologyInfo(req *http.Request) (TopologyInfo, error) {
	topologyInfo, exist := TopologyInfoFrom(req.Context())
	if !exist {
		return nil, fmt.Errorf("can not get TopologyInfo")
	}
	if !topologyInfo.IsTopologyResource() {
		return topologyInfo, nil
	}
	requestInfo, exist := request.RequestInfoFrom(req.Context())
	if !exist {
		return topologyInfo, nil
	}

	if requestInfo.Verb != "list" && requestInfo.Verb != "watch" {
		return topologyInfo, nil
	}

	// Populate only if request query is empty
	if value := req.FormValue("labelSelector"); value != "" {
		return topologyInfo, nil
	}

	userFrom, exist := request.UserFrom(req.Context())
	if !exist {
		return topologyInfo, nil
	}
	username := userFrom.GetName()
	globalRoleBindings, err := p.globalRoleBindingLister.List(labels.Everything())
	if err != nil {
		return nil, err
	}

	for _, binding := range globalRoleBindings {
		for _, subject := range binding.Subjects {
			if subject.Name != username || subject.Kind != rbacv1.UserKind {
				continue
			}
			globalRole, err := p.globalRoleLister.Get(binding.RoleRef.Name)
			if err != nil {
				return nil, err
			}

			topologyInfo, err = p.getTopologyInfoFromGlobalRole(req.Context(), globalRole, topologyInfo)
			if err != nil {
				return nil, err
			}
		}
	}

	req.Form.Set("labelSelector", createTopologyQuery(topologyInfo))

	return topologyInfo, nil
}

func (p *populator) getTopologyInfoFromGlobalRole(ctx context.Context, globalRole *v1alpha1.GlobalRole, info TopologyInfo) (TopologyInfo, error) {
	var topologyInfo *DefaultTopologyInfo
	if info == nil {
		topologyInfo = NewDefaultTopologyInfo(true)
	} else {
		topologyInfo = info.(*DefaultTopologyInfo)
	}

	if globalRole.RegionSelectors != nil {
		for _, rs := range globalRole.RegionSelectors {
			selector, err := metav1.LabelSelectorAsSelector(&rs)
			if err != nil {
				return nil, err
			}
			regions, err := p.TopologyGetter.ListRegions(ctx, metav1.ListOptions{LabelSelector: selector.String()})
			if err != nil {
				return nil, err
			}
			for _, region := range regions {
				topologyInfo.Add("", region.GetName(), nil, region.GetLabels())
			}
		}
	}
	if globalRole.ClusterSelectors != nil {
		for _, cs := range globalRole.ClusterSelectors {

			selector, err := metav1.LabelSelectorAsSelector(&cs)
			if err != nil {
				return nil, err
			}
			clusters, err := p.TopologyGetter.ListClusters(ctx, metav1.ListOptions{LabelSelector: selector.String()})
			if err != nil {
				return nil, err
			}
			for _, cluster := range clusters {
				topologyInfo.Add(cluster.GetName(), "", cluster.GetLabels(), nil)
			}
		}
	}
	return topologyInfo, nil
}

func createTopologyQuery(topologyInfo TopologyInfo) string {
	var regions, clusters, regionQuery, clusterQuery string
	if topologyInfo.GetRegionName().Len() != 0 {
		regions = strings.Join(topologyInfo.GetRegionName().List(), ",")
	}
	if regions != "" {
		regionQuery = fmt.Sprintf("%s in (%s)", regionLabelKey, regions)
	}

	if topologyInfo.GetClusterName().Len() != 0 {
		clusters = strings.Join(topologyInfo.GetClusterName().List(), ",")
	}
	if clusters != "" {
		clusterQuery = fmt.Sprintf("%s in (%s)", clusterLabelKey, clusters)
	}

	if regionQuery != "" && clusterQuery == "" {
		return regionQuery
	}
	if clusterQuery != "" && regionQuery == "" {
		return clusterQuery
	}
	return fmt.Sprintf("%s,%s", regionQuery, clusterQuery)
}
