/*
Copyright The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Code generated by client-gen. DO NOT EDIT.

package v1alpha1

import (
	"context"
	json "encoding/json"
	"fmt"
	"time"

	v1alpha1 "k8s.io/api/alcor/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	alcorv1alpha1 "k8s.io/client-go/applyconfigurations/alcor/v1alpha1"
	scheme "k8s.io/client-go/kubernetes/scheme"
	rest "k8s.io/client-go/rest"
)

// IPSetsGetter has a method to return a IPSetInterface.
// A group's client should implement this interface.
type IPSetsGetter interface {
	IPSets(namespace string) IPSetInterface
}

// IPSetInterface has methods to work with IPSet resources.
type IPSetInterface interface {
	Create(ctx context.Context, iPSet *v1alpha1.IPSet, opts v1.CreateOptions) (*v1alpha1.IPSet, error)
	Update(ctx context.Context, iPSet *v1alpha1.IPSet, opts v1.UpdateOptions) (*v1alpha1.IPSet, error)
	UpdateStatus(ctx context.Context, iPSet *v1alpha1.IPSet, opts v1.UpdateOptions) (*v1alpha1.IPSet, error)
	Delete(ctx context.Context, name string, opts v1.DeleteOptions) error
	DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error
	Get(ctx context.Context, name string, opts v1.GetOptions) (*v1alpha1.IPSet, error)
	List(ctx context.Context, opts v1.ListOptions) (*v1alpha1.IPSetList, error)
	Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error)
	Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1alpha1.IPSet, err error)
	Apply(ctx context.Context, iPSet *alcorv1alpha1.IPSetApplyConfiguration, opts v1.ApplyOptions) (result *v1alpha1.IPSet, err error)
	ApplyStatus(ctx context.Context, iPSet *alcorv1alpha1.IPSetApplyConfiguration, opts v1.ApplyOptions) (result *v1alpha1.IPSet, err error)
	IPSetExpansion
}

// iPSets implements IPSetInterface
type iPSets struct {
	client rest.Interface
	ns     string
}

// newIPSets returns a IPSets
func newIPSets(c *AlcorV1alpha1Client, namespace string) *iPSets {
	return &iPSets{
		client: c.RESTClient(),
		ns:     namespace,
	}
}

// Get takes name of the iPSet, and returns the corresponding iPSet object, and an error if there is any.
func (c *iPSets) Get(ctx context.Context, name string, options v1.GetOptions) (result *v1alpha1.IPSet, err error) {
	result = &v1alpha1.IPSet{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("ipsets").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do(ctx).
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of IPSets that match those selectors.
func (c *iPSets) List(ctx context.Context, opts v1.ListOptions) (result *v1alpha1.IPSetList, err error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	result = &v1alpha1.IPSetList{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("ipsets").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Do(ctx).
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested iPSets.
func (c *iPSets) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	opts.Watch = true
	return c.client.Get().
		Namespace(c.ns).
		Resource("ipsets").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Watch(ctx)
}

// Create takes the representation of a iPSet and creates it.  Returns the server's representation of the iPSet, and an error, if there is any.
func (c *iPSets) Create(ctx context.Context, iPSet *v1alpha1.IPSet, opts v1.CreateOptions) (result *v1alpha1.IPSet, err error) {
	result = &v1alpha1.IPSet{}
	err = c.client.Post().
		Namespace(c.ns).
		Resource("ipsets").
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(iPSet).
		Do(ctx).
		Into(result)
	return
}

// Update takes the representation of a iPSet and updates it. Returns the server's representation of the iPSet, and an error, if there is any.
func (c *iPSets) Update(ctx context.Context, iPSet *v1alpha1.IPSet, opts v1.UpdateOptions) (result *v1alpha1.IPSet, err error) {
	result = &v1alpha1.IPSet{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("ipsets").
		Name(iPSet.Name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(iPSet).
		Do(ctx).
		Into(result)
	return
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *iPSets) UpdateStatus(ctx context.Context, iPSet *v1alpha1.IPSet, opts v1.UpdateOptions) (result *v1alpha1.IPSet, err error) {
	result = &v1alpha1.IPSet{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("ipsets").
		Name(iPSet.Name).
		SubResource("status").
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(iPSet).
		Do(ctx).
		Into(result)
	return
}

// Delete takes name of the iPSet and deletes it. Returns an error if one occurs.
func (c *iPSets) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("ipsets").
		Name(name).
		Body(&opts).
		Do(ctx).
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *iPSets) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	var timeout time.Duration
	if listOpts.TimeoutSeconds != nil {
		timeout = time.Duration(*listOpts.TimeoutSeconds) * time.Second
	}
	return c.client.Delete().
		Namespace(c.ns).
		Resource("ipsets").
		VersionedParams(&listOpts, scheme.ParameterCodec).
		Timeout(timeout).
		Body(&opts).
		Do(ctx).
		Error()
}

// Patch applies the patch and returns the patched iPSet.
func (c *iPSets) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1alpha1.IPSet, err error) {
	result = &v1alpha1.IPSet{}
	err = c.client.Patch(pt).
		Namespace(c.ns).
		Resource("ipsets").
		Name(name).
		SubResource(subresources...).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(data).
		Do(ctx).
		Into(result)
	return
}

// Apply takes the given apply declarative configuration, applies it and returns the applied iPSet.
func (c *iPSets) Apply(ctx context.Context, iPSet *alcorv1alpha1.IPSetApplyConfiguration, opts v1.ApplyOptions) (result *v1alpha1.IPSet, err error) {
	if iPSet == nil {
		return nil, fmt.Errorf("iPSet provided to Apply must not be nil")
	}
	patchOpts := opts.ToPatchOptions()
	data, err := json.Marshal(iPSet)
	if err != nil {
		return nil, err
	}
	name := iPSet.Name
	if name == nil {
		return nil, fmt.Errorf("iPSet.Name must be provided to Apply")
	}
	result = &v1alpha1.IPSet{}
	err = c.client.Patch(types.ApplyPatchType).
		Namespace(c.ns).
		Resource("ipsets").
		Name(*name).
		VersionedParams(&patchOpts, scheme.ParameterCodec).
		Body(data).
		Do(ctx).
		Into(result)
	return
}

// ApplyStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating ApplyStatus().
func (c *iPSets) ApplyStatus(ctx context.Context, iPSet *alcorv1alpha1.IPSetApplyConfiguration, opts v1.ApplyOptions) (result *v1alpha1.IPSet, err error) {
	if iPSet == nil {
		return nil, fmt.Errorf("iPSet provided to Apply must not be nil")
	}
	patchOpts := opts.ToPatchOptions()
	data, err := json.Marshal(iPSet)
	if err != nil {
		return nil, err
	}

	name := iPSet.Name
	if name == nil {
		return nil, fmt.Errorf("iPSet.Name must be provided to Apply")
	}

	result = &v1alpha1.IPSet{}
	err = c.client.Patch(types.ApplyPatchType).
		Namespace(c.ns).
		Resource("ipsets").
		Name(*name).
		SubResource("status").
		VersionedParams(&patchOpts, scheme.ParameterCodec).
		Body(data).
		Do(ctx).
		Into(result)
	return
}
