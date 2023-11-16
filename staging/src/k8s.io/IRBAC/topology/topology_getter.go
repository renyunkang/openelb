package topology

import (
	"context"

	"k8s.io/apimachinery/pkg/api/meta"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apiserver/pkg/endpoints/topology"
	"k8s.io/client-go/dynamic"
)

const (
	regionResourceName  = "regions"
	clusterResourceName = "clusters"
)

var (
	regionGVR  = schema.GroupVersionResource{Group: "core.plan.mizar-universe.io", Version: "v1alpha1", Resource: "regions"}
	clusterGVR = schema.GroupVersionResource{Group: "core.mizargalaxy.mizar-k8s.io", Version: "v1alpha1", Resource: "clusters"}
)

func NewTopologyGetter(client dynamic.Interface) topology.TopologyGetter {
	return topologyGetter{client: client}
}

type topologyGetter struct {
	client dynamic.Interface
}

func (t topologyGetter) GetRegion(ctx context.Context, name string) (v1.Object, error) {
	region, err := t.client.
		Resource(regionGVR).
		Get(ctx, name, v1.GetOptions{})
	if err != nil {
		return nil, err
	}
	return region, nil
}

func (t topologyGetter) GetCluster(ctx context.Context, name string) (v1.Object, error) {
	cluster, err := t.client.
		Resource(clusterGVR).
		Get(ctx, name, v1.GetOptions{})
	if err != nil {
		return nil, err
	}
	return cluster, nil
}

func (t topologyGetter) ListClusters(ctx context.Context, opt v1.ListOptions) ([]v1.Object, error) {
	var objs []v1.Object

	clusters, err := t.client.
		Resource(clusterGVR).List(ctx, opt)
	if err != nil {
		return nil, err
	}

	list, err := meta.ExtractList(clusters)
	if err != nil {
		return nil, err
	}

	for _, item := range list {
		accessor, err := meta.Accessor(item)
		if err != nil {
			return nil, err
		}
		objs = append(objs, accessor)
	}
	return objs, nil
}

func (t topologyGetter) ListRegions(ctx context.Context, opt v1.ListOptions) ([]v1.Object, error) {
	var objs []v1.Object

	regions, err := t.client.
		Resource(regionGVR).List(ctx, opt)
	if err != nil {
		return nil, err
	}

	list, err := meta.ExtractList(regions)
	if err != nil {
		return nil, err
	}

	for _, item := range list {
		accessor, err := meta.Accessor(item)
		if err != nil {
			return nil, err
		}
		objs = append(objs, accessor)
	}
	return objs, nil
}
