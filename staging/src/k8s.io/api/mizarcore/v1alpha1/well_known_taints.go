package v1alpha1

var (
	// TaintClusterNotReady will be added when cluster is not offline（ Cluster Condition Ready = False）
	// and removed when cluster becomes ready.
	TaintClusterNotReady = "cluster.kubernetes.io/not-ready"

	// TaintClusterNotReady will be added when cluster Cluster Condition is UNKNOWN
	// and removed when cluster becomes ready.
	TaintClusterUnreachable = "cluster.kubernetes.io/unreachable"
)
