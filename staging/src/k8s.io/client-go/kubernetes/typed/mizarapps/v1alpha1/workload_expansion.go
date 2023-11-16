package v1alpha1

import (
	"context"

	"k8s.io/api/mizarapps/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/scheme"
)

type WorkloadExpansion interface {
	Bind(ctx context.Context, binding *v1alpha1.Binding, opts metav1.CreateOptions) error
}

// Bind applies the provided binding to the named workload in the current namespace (binding.Namespace is ignored).
func (c *workloads) Bind(ctx context.Context, binding *v1alpha1.Binding, opts metav1.CreateOptions) error {
	return c.client.Post().Namespace(c.ns).Resource("workloads").Name(binding.Name).VersionedParams(&opts, scheme.ParameterCodec).SubResource("binding").Body(binding).Do(ctx).Error()
}
