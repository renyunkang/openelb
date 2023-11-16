package fake

import (
	"context"
	"k8s.io/api/mizarapps/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	core "k8s.io/client-go/testing"
)

func (c *FakeWorkloads) Bind(ctx context.Context, binding *v1alpha1.Binding, opts metav1.CreateOptions) error {
	action := core.CreateActionImpl{}
	action.Verb = "create"
	action.Namespace = binding.Namespace
	action.Resource = workloadsResource
	action.Subresource = "binding"
	action.Object = binding

	_, err := c.Fake.Invokes(action, binding)
	return err
}

func (c *FakeWorkloads) GetBinding(name string) (result *v1alpha1.Binding, err error) {
	obj, err := c.Fake.
		Invokes(core.NewGetSubresourceAction(workloadsResource, c.ns, "binding", name), &v1alpha1.Binding{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.Binding), err
}
