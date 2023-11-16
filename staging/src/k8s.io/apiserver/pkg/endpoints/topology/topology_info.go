package topology

import (
	"context"
	"errors"
	"fmt"
	"io"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/util/sets"
	"net/http"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apiserver/pkg/endpoints/request"
	"k8s.io/client-go/informers"
	"k8s.io/klog/v2"
)

const (
	topologyInfoCtxKey = "mizar-topology-info"

	regionLabelKey  = "topology.mizar-k8s.io/region"
	clusterLabelKey = "topology.mizar-k8s.io/cluster"
)

func WithTopologyInfo(parent context.Context, info TopologyInfo) context.Context {
	return context.WithValue(parent, topologyInfoCtxKey, info)
}

func TopologyInfoFrom(ctx context.Context) (TopologyInfo, bool) {
	if ctx == nil {
		return nil, false
	}
	value, ok := ctx.Value(topologyInfoCtxKey).(TopologyInfo)
	if ok {
		return value, true
	}

	return nil, false
}

func NewTopologyResolver(informers map[schema.GroupVersionResource]informers.GenericInformer, getter TopologyGetter) TopologyInfoResolver {
	return &TopologyInfoFactory{
		informers:      informers,
		TopologyGetter: getter,
	}
}

type TopologyInfo interface {
	GetRegionName() sets.String
	GetClusterName() sets.String
	GetRegionLabels() []map[string]string
	GetClusterLabels() []map[string]string

	IsTopologyResource() bool
}

type DefaultTopologyInfo struct {
	RegionNames   sets.String
	ClusterNames  sets.String
	RegionLabels  []map[string]string
	ClusterLabels []map[string]string

	TopologyResource bool
}

func NewDefaultTopologyInfo(isTopologyResource bool) *DefaultTopologyInfo {
	return &DefaultTopologyInfo{
		TopologyResource: isTopologyResource,
		ClusterNames:     sets.NewString(),
		RegionNames:      sets.NewString(),
		ClusterLabels:    make([]map[string]string, 0),
		RegionLabels:     make([]map[string]string, 0),
	}
}

func (t *DefaultTopologyInfo) GetRegionLabels() []map[string]string {
	return t.RegionLabels
}

func (t *DefaultTopologyInfo) GetClusterLabels() []map[string]string {
	return t.ClusterLabels
}

func (t *DefaultTopologyInfo) GetRegionName() sets.String {
	return t.RegionNames
}

func (t *DefaultTopologyInfo) GetClusterName() sets.String {
	return t.ClusterNames
}

func (t *DefaultTopologyInfo) IsTopologyResource() bool {
	return t.TopologyResource
}

type TopologyInfoResolver interface {
	NewTopologyInfo(req *http.Request) (TopologyInfo, error)
}

type TopologyInfoFactory struct {
	informers map[schema.GroupVersionResource]informers.GenericInformer
	TopologyGetter
}

func (t *TopologyInfoFactory) NewTopologyInfo(req *http.Request) (TopologyInfo, error) {
	var resourceLabels map[string]string
	var topologyInfo = NewDefaultTopologyInfo(false)

	requestInfo, exist := request.RequestInfoFrom(req.Context())
	if !exist {
		return nil, errors.New("request information not exist")
	}
	if !requestInfo.IsResourceRequest {
		return topologyInfo, nil
	}

	informer, exist := t.informers[schema.GroupVersionResource{
		Group:    requestInfo.APIGroup,
		Version:  requestInfo.APIVersion,
		Resource: requestInfo.Resource}]

	if !exist {
		klog.V(6).Infof("informer for resource %s not exist, skip request", requestInfo.Resource)
		return topologyInfo, nil
	}

	topologyInfo.TopologyResource = true

	switch requestInfo.Verb {
	// resource has been created, just get the resource labels
	case "get", "update", "patch", "delete":
		obj, err := informer.Lister().ByNamespace(requestInfo.Namespace).Get(requestInfo.Name)
		if err != nil {
			return nil, err
		}
		metaObj, err := meta.Accessor(obj)
		if err != nil {
			return nil, err
		}

		resourceLabels = metaObj.GetLabels()
	case "create":
		all, err := io.ReadAll(req.Body)
		if err != nil {
			return nil, err
		}
		obj := unstructured.Unstructured{}
		err = obj.UnmarshalJSON(all)
		if err != nil {
			return nil, err
		}

		resourceLabels = obj.GetLabels()
	case "list", "watch":
		var selector labels.Selector
		var err error
		// ?labelSelector=topology.mizar-k8s.io/region=xxx,topology.mizar-k8s.io/cluster=xxx
		rawSelector := req.FormValue("labelSelector")
		if rawSelector != "" {
			selector, err = labels.Parse(rawSelector)
			if err != nil {
				return nil, err
			}
		}

		if selector != nil {
			return t.getTopologyInfoByResourceSelector(req.Context(), informer, requestInfo.Namespace, selector)
		} else {
			return topologyInfo, nil
		}
	}

	clusterName, regionName := resourceLabels[clusterLabelKey], resourceLabels[regionLabelKey]
	regionLabel, clusterLabel, err := t.getTopologySelectorByTopologyName(req.Context(), regionName, clusterName)
	if err != nil {
		return nil, err
	}
	topologyInfo.Add(clusterName, regionName, clusterLabel, regionLabel)

	return topologyInfo, nil
}

func (t *TopologyInfoFactory) getTopologySelectorByTopologyName(
	ctx context.Context, regionName, clusterName string) (
	regionLabel, clusterLabel map[string]string, err error) {
	if regionName != "" {
		region, err := t.TopologyGetter.GetRegion(ctx, regionName)
		if err != nil {
			return nil, nil, err
		}
		regionLabel = region.GetLabels()
	}

	if clusterName != "" {
		if clusterName != regionName {
			clusterName = fmt.Sprintf("%s.%s", regionName, clusterName)
		}
		cluster, err := t.TopologyGetter.GetCluster(ctx, clusterName)
		if err != nil {
			return nil, nil, err
		}
		clusterLabel = cluster.GetLabels()
	}
	return
}

func (t *TopologyInfoFactory) getTopologyInfoByResourceSelector(ctx context.Context, informer informers.GenericInformer, namespace string, selector labels.Selector) (TopologyInfo, error) {
	var topologyInfo = NewDefaultTopologyInfo(true)

	if selector != nil {
		// handle the list request when the selector is not regionLabelKey or clusterLabelKey
		resources, err := informer.Lister().ByNamespace(namespace).List(selector)
		if err != nil {
			return nil, err
		}
		for _, rs := range resources {
			accessor, err := meta.Accessor(rs)
			if err != nil {
				return nil, err
			}
			resourceLabels := accessor.GetLabels()
			clusterName, regionName := resourceLabels[clusterLabelKey], resourceLabels[regionLabelKey]
			cl, rl, err := t.getTopologySelectorByTopologyName(ctx, regionName, clusterName)
			if err != nil {
				return nil, err
			}
			topologyInfo.Add(clusterName, regionName, cl, rl)
		}
	}
	return topologyInfo, nil
}

func (t *DefaultTopologyInfo) Add(clusterName, regionName string, clusterLabel, regionLabel map[string]string) {
	if regionName != "" && !t.RegionNames.Has(regionName) {
		t.RegionNames = t.RegionNames.Insert(regionName)
		if regionLabel != nil {
			t.RegionLabels = append(t.RegionLabels, regionLabel)
		}
	}

	if clusterName != "" && !t.ClusterNames.Has(clusterName) {
		t.ClusterNames = t.ClusterNames.Insert(clusterName)
		if clusterLabel != nil {
			t.ClusterLabels = append(t.ClusterLabels, clusterLabel)
		}
	}
}

type TopologyGetter interface {
	GetRegion(ctx context.Context, name string) (v1.Object, error)
	GetCluster(ctx context.Context, name string) (v1.Object, error)
	ListRegions(ctx context.Context, opt v1.ListOptions) ([]v1.Object, error)
	ListClusters(ctx context.Context, opt v1.ListOptions) ([]v1.Object, error)
}

type SimpleTopologyGetter struct {
	Regions  map[string]v1.Object
	Clusters map[string]v1.Object
}

var _ TopologyGetter = &SimpleTopologyGetter{}

func (s *SimpleTopologyGetter) GetRegion(_ context.Context, name string) (v1.Object, error) {
	return s.Regions[name], nil
}

func (s *SimpleTopologyGetter) GetCluster(_ context.Context, name string) (v1.Object, error) {
	return s.Clusters[name], nil
}

func (s *SimpleTopologyGetter) ListRegions(_ context.Context, opt v1.ListOptions) ([]v1.Object, error) {
	var result []v1.Object
	selector := opt.LabelSelector
	for _, region := range s.Regions {
		l := region.GetLabels()
		parse, err := labels.Parse(selector)
		if err != nil {
			return nil, err
		}

		if parse.Matches(labels.Set(l)) {
			result = append(result, region)
		}
	}
	return result, nil
}

func (s *SimpleTopologyGetter) ListClusters(ctx context.Context, opt v1.ListOptions) ([]v1.Object, error) {
	var result []v1.Object
	selector := opt.LabelSelector
	for _, cluster := range s.Clusters {
		l := cluster.GetLabels()
		parse, err := labels.Parse(selector)
		if err != nil {
			return nil, err
		}

		if parse.Matches(labels.Set(l)) {
			result = append(result, cluster)
		}
	}
	return result, nil
}
