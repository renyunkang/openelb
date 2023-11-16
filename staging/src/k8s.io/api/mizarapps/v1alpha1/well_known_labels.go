package v1alpha1

const (
	ApplicationLabel = "app.plan.mizar-universe.io/name"
	ComponentLabel   = "app.plan.mizar-universe.io/component"
	DeployScopeLabel = "deployscope.plan.mizar-universe.io/name"
	NetworkScope     = "networkscope.plan.mizar-universe.io/name"
	WorkloadLabel    = "workload.plan.mizar-universe.io/name"

	ClusterLabel      = "topology.mizar-k8s.io/cluster"
	RegionLabel       = "topology.mizar-k8s.io/region"
	PhysicalZoneLabel = "topology.mizar-k8s.io/pzone"
	LogicalZoneLabel  = "topology.mizar-k8s.io/lzone"
	NodeLabel         = "topology.mizar-k8s.io/node"
)

var (
	WorkloadApplicationLabels = []string{
		ApplicationLabel,
		ComponentLabel,
		DeployScopeLabel,
		LogicalZoneLabel,
		NetworkScope,
		WorkloadLabel,
	}
	WorkloadTopologyLables = []string{
		ClusterLabel,
		RegionLabel,
		PhysicalZoneLabel,
	}
	WorkloadLabels = []string{
		ApplicationLabel,
		ComponentLabel,
		DeployScopeLabel,
		LogicalZoneLabel,
		NetworkScope,
		WorkloadLabel,
		ClusterLabel,
		RegionLabel,
		PhysicalZoneLabel,
	}
)
