package v1alpha1

import (
	"context"

	"k8s.io/api/alcor/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/scheme"
)

type StableModelExpansion interface {
	Bind(ctx context.Context, binding *v1alpha1.Binding, opts metav1.CreateOptions) error
}

// Bind applies the provided binding to the named stablemodel in the current namespace (binding.Namespace is ignored).
func (c *stableModels) Bind(ctx context.Context, binding *v1alpha1.Binding, opts metav1.CreateOptions) error {
	return c.client.Post().Namespace(c.ns).Resource("stablemodels").Name(binding.Name).VersionedParams(&opts, scheme.ParameterCodec).SubResource("binding").Body(binding).Do(ctx).Error()
}
