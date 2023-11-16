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

	v1alpha1 "k8s.io/api/mizarapps/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	mizarappsv1alpha1 "k8s.io/client-go/applyconfigurations/mizarapps/v1alpha1"
	scheme "k8s.io/client-go/kubernetes/scheme"
	rest "k8s.io/client-go/rest"
)

// GlobalQuotasGetter has a method to return a GlobalQuotaInterface.
// A group's client should implement this interface.
type GlobalQuotasGetter interface {
	GlobalQuotas(namespace string) GlobalQuotaInterface
}

// GlobalQuotaInterface has methods to work with GlobalQuota resources.
type GlobalQuotaInterface interface {
	Create(ctx context.Context, globalQuota *v1alpha1.GlobalQuota, opts v1.CreateOptions) (*v1alpha1.GlobalQuota, error)
	Update(ctx context.Context, globalQuota *v1alpha1.GlobalQuota, opts v1.UpdateOptions) (*v1alpha1.GlobalQuota, error)
	UpdateStatus(ctx context.Context, globalQuota *v1alpha1.GlobalQuota, opts v1.UpdateOptions) (*v1alpha1.GlobalQuota, error)
	Delete(ctx context.Context, name string, opts v1.DeleteOptions) error
	DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error
	Get(ctx context.Context, name string, opts v1.GetOptions) (*v1alpha1.GlobalQuota, error)
	List(ctx context.Context, opts v1.ListOptions) (*v1alpha1.GlobalQuotaList, error)
	Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error)
	Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1alpha1.GlobalQuota, err error)
	Apply(ctx context.Context, globalQuota *mizarappsv1alpha1.GlobalQuotaApplyConfiguration, opts v1.ApplyOptions) (result *v1alpha1.GlobalQuota, err error)
	ApplyStatus(ctx context.Context, globalQuota *mizarappsv1alpha1.GlobalQuotaApplyConfiguration, opts v1.ApplyOptions) (result *v1alpha1.GlobalQuota, err error)
	GlobalQuotaExpansion
}

// globalQuotas implements GlobalQuotaInterface
type globalQuotas struct {
	client rest.Interface
	ns     string
}

// newGlobalQuotas returns a GlobalQuotas
func newGlobalQuotas(c *MizarappsV1alpha1Client, namespace string) *globalQuotas {
	return &globalQuotas{
		client: c.RESTClient(),
		ns:     namespace,
	}
}

// Get takes name of the globalQuota, and returns the corresponding globalQuota object, and an error if there is any.
func (c *globalQuotas) Get(ctx context.Context, name string, options v1.GetOptions) (result *v1alpha1.GlobalQuota, err error) {
	result = &v1alpha1.GlobalQuota{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("globalquotas").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do(ctx).
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of GlobalQuotas that match those selectors.
func (c *globalQuotas) List(ctx context.Context, opts v1.ListOptions) (result *v1alpha1.GlobalQuotaList, err error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	result = &v1alpha1.GlobalQuotaList{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("globalquotas").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Do(ctx).
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested globalQuotas.
func (c *globalQuotas) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	opts.Watch = true
	return c.client.Get().
		Namespace(c.ns).
		Resource("globalquotas").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Watch(ctx)
}

// Create takes the representation of a globalQuota and creates it.  Returns the server's representation of the globalQuota, and an error, if there is any.
func (c *globalQuotas) Create(ctx context.Context, globalQuota *v1alpha1.GlobalQuota, opts v1.CreateOptions) (result *v1alpha1.GlobalQuota, err error) {
	result = &v1alpha1.GlobalQuota{}
	err = c.client.Post().
		Namespace(c.ns).
		Resource("globalquotas").
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(globalQuota).
		Do(ctx).
		Into(result)
	return
}

// Update takes the representation of a globalQuota and updates it. Returns the server's representation of the globalQuota, and an error, if there is any.
func (c *globalQuotas) Update(ctx context.Context, globalQuota *v1alpha1.GlobalQuota, opts v1.UpdateOptions) (result *v1alpha1.GlobalQuota, err error) {
	result = &v1alpha1.GlobalQuota{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("globalquotas").
		Name(globalQuota.Name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(globalQuota).
		Do(ctx).
		Into(result)
	return
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *globalQuotas) UpdateStatus(ctx context.Context, globalQuota *v1alpha1.GlobalQuota, opts v1.UpdateOptions) (result *v1alpha1.GlobalQuota, err error) {
	result = &v1alpha1.GlobalQuota{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("globalquotas").
		Name(globalQuota.Name).
		SubResource("status").
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(globalQuota).
		Do(ctx).
		Into(result)
	return
}

// Delete takes name of the globalQuota and deletes it. Returns an error if one occurs.
func (c *globalQuotas) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("globalquotas").
		Name(name).
		Body(&opts).
		Do(ctx).
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *globalQuotas) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	var timeout time.Duration
	if listOpts.TimeoutSeconds != nil {
		timeout = time.Duration(*listOpts.TimeoutSeconds) * time.Second
	}
	return c.client.Delete().
		Namespace(c.ns).
		Resource("globalquotas").
		VersionedParams(&listOpts, scheme.ParameterCodec).
		Timeout(timeout).
		Body(&opts).
		Do(ctx).
		Error()
}

// Patch applies the patch and returns the patched globalQuota.
func (c *globalQuotas) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1alpha1.GlobalQuota, err error) {
	result = &v1alpha1.GlobalQuota{}
	err = c.client.Patch(pt).
		Namespace(c.ns).
		Resource("globalquotas").
		Name(name).
		SubResource(subresources...).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(data).
		Do(ctx).
		Into(result)
	return
}

// Apply takes the given apply declarative configuration, applies it and returns the applied globalQuota.
func (c *globalQuotas) Apply(ctx context.Context, globalQuota *mizarappsv1alpha1.GlobalQuotaApplyConfiguration, opts v1.ApplyOptions) (result *v1alpha1.GlobalQuota, err error) {
	if globalQuota == nil {
		return nil, fmt.Errorf("globalQuota provided to Apply must not be nil")
	}
	patchOpts := opts.ToPatchOptions()
	data, err := json.Marshal(globalQuota)
	if err != nil {
		return nil, err
	}
	name := globalQuota.Name
	if name == nil {
		return nil, fmt.Errorf("globalQuota.Name must be provided to Apply")
	}
	result = &v1alpha1.GlobalQuota{}
	err = c.client.Patch(types.ApplyPatchType).
		Namespace(c.ns).
		Resource("globalquotas").
		Name(*name).
		VersionedParams(&patchOpts, scheme.ParameterCodec).
		Body(data).
		Do(ctx).
		Into(result)
	return
}

// ApplyStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating ApplyStatus().
func (c *globalQuotas) ApplyStatus(ctx context.Context, globalQuota *mizarappsv1alpha1.GlobalQuotaApplyConfiguration, opts v1.ApplyOptions) (result *v1alpha1.GlobalQuota, err error) {
	if globalQuota == nil {
		return nil, fmt.Errorf("globalQuota provided to Apply must not be nil")
	}
	patchOpts := opts.ToPatchOptions()
	data, err := json.Marshal(globalQuota)
	if err != nil {
		return nil, err
	}

	name := globalQuota.Name
	if name == nil {
		return nil, fmt.Errorf("globalQuota.Name must be provided to Apply")
	}

	result = &v1alpha1.GlobalQuota{}
	err = c.client.Patch(types.ApplyPatchType).
		Namespace(c.ns).
		Resource("globalquotas").
		Name(*name).
		SubResource("status").
		VersionedParams(&patchOpts, scheme.ParameterCodec).
		Body(data).
		Do(ctx).
		Into(result)
	return
}
