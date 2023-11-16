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

// StableModelsGetter has a method to return a StableModelInterface.
// A group's client should implement this interface.
type StableModelsGetter interface {
	StableModels(namespace string) StableModelInterface
}

// StableModelInterface has methods to work with StableModel resources.
type StableModelInterface interface {
	Create(ctx context.Context, stableModel *v1alpha1.StableModel, opts v1.CreateOptions) (*v1alpha1.StableModel, error)
	Update(ctx context.Context, stableModel *v1alpha1.StableModel, opts v1.UpdateOptions) (*v1alpha1.StableModel, error)
	UpdateStatus(ctx context.Context, stableModel *v1alpha1.StableModel, opts v1.UpdateOptions) (*v1alpha1.StableModel, error)
	Delete(ctx context.Context, name string, opts v1.DeleteOptions) error
	DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error
	Get(ctx context.Context, name string, opts v1.GetOptions) (*v1alpha1.StableModel, error)
	List(ctx context.Context, opts v1.ListOptions) (*v1alpha1.StableModelList, error)
	Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error)
	Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1alpha1.StableModel, err error)
	Apply(ctx context.Context, stableModel *alcorv1alpha1.StableModelApplyConfiguration, opts v1.ApplyOptions) (result *v1alpha1.StableModel, err error)
	ApplyStatus(ctx context.Context, stableModel *alcorv1alpha1.StableModelApplyConfiguration, opts v1.ApplyOptions) (result *v1alpha1.StableModel, err error)
	StableModelExpansion
}

// stableModels implements StableModelInterface
type stableModels struct {
	client rest.Interface
	ns     string
}

// newStableModels returns a StableModels
func newStableModels(c *AlcorV1alpha1Client, namespace string) *stableModels {
	return &stableModels{
		client: c.RESTClient(),
		ns:     namespace,
	}
}

// Get takes name of the stableModel, and returns the corresponding stableModel object, and an error if there is any.
func (c *stableModels) Get(ctx context.Context, name string, options v1.GetOptions) (result *v1alpha1.StableModel, err error) {
	result = &v1alpha1.StableModel{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("stablemodels").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do(ctx).
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of StableModels that match those selectors.
func (c *stableModels) List(ctx context.Context, opts v1.ListOptions) (result *v1alpha1.StableModelList, err error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	result = &v1alpha1.StableModelList{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("stablemodels").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Do(ctx).
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested stableModels.
func (c *stableModels) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	opts.Watch = true
	return c.client.Get().
		Namespace(c.ns).
		Resource("stablemodels").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Watch(ctx)
}

// Create takes the representation of a stableModel and creates it.  Returns the server's representation of the stableModel, and an error, if there is any.
func (c *stableModels) Create(ctx context.Context, stableModel *v1alpha1.StableModel, opts v1.CreateOptions) (result *v1alpha1.StableModel, err error) {
	result = &v1alpha1.StableModel{}
	err = c.client.Post().
		Namespace(c.ns).
		Resource("stablemodels").
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(stableModel).
		Do(ctx).
		Into(result)
	return
}

// Update takes the representation of a stableModel and updates it. Returns the server's representation of the stableModel, and an error, if there is any.
func (c *stableModels) Update(ctx context.Context, stableModel *v1alpha1.StableModel, opts v1.UpdateOptions) (result *v1alpha1.StableModel, err error) {
	result = &v1alpha1.StableModel{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("stablemodels").
		Name(stableModel.Name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(stableModel).
		Do(ctx).
		Into(result)
	return
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *stableModels) UpdateStatus(ctx context.Context, stableModel *v1alpha1.StableModel, opts v1.UpdateOptions) (result *v1alpha1.StableModel, err error) {
	result = &v1alpha1.StableModel{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("stablemodels").
		Name(stableModel.Name).
		SubResource("status").
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(stableModel).
		Do(ctx).
		Into(result)
	return
}

// Delete takes name of the stableModel and deletes it. Returns an error if one occurs.
func (c *stableModels) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("stablemodels").
		Name(name).
		Body(&opts).
		Do(ctx).
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *stableModels) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	var timeout time.Duration
	if listOpts.TimeoutSeconds != nil {
		timeout = time.Duration(*listOpts.TimeoutSeconds) * time.Second
	}
	return c.client.Delete().
		Namespace(c.ns).
		Resource("stablemodels").
		VersionedParams(&listOpts, scheme.ParameterCodec).
		Timeout(timeout).
		Body(&opts).
		Do(ctx).
		Error()
}

// Patch applies the patch and returns the patched stableModel.
func (c *stableModels) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1alpha1.StableModel, err error) {
	result = &v1alpha1.StableModel{}
	err = c.client.Patch(pt).
		Namespace(c.ns).
		Resource("stablemodels").
		Name(name).
		SubResource(subresources...).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(data).
		Do(ctx).
		Into(result)
	return
}

// Apply takes the given apply declarative configuration, applies it and returns the applied stableModel.
func (c *stableModels) Apply(ctx context.Context, stableModel *alcorv1alpha1.StableModelApplyConfiguration, opts v1.ApplyOptions) (result *v1alpha1.StableModel, err error) {
	if stableModel == nil {
		return nil, fmt.Errorf("stableModel provided to Apply must not be nil")
	}
	patchOpts := opts.ToPatchOptions()
	data, err := json.Marshal(stableModel)
	if err != nil {
		return nil, err
	}
	name := stableModel.Name
	if name == nil {
		return nil, fmt.Errorf("stableModel.Name must be provided to Apply")
	}
	result = &v1alpha1.StableModel{}
	err = c.client.Patch(types.ApplyPatchType).
		Namespace(c.ns).
		Resource("stablemodels").
		Name(*name).
		VersionedParams(&patchOpts, scheme.ParameterCodec).
		Body(data).
		Do(ctx).
		Into(result)
	return
}

// ApplyStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating ApplyStatus().
func (c *stableModels) ApplyStatus(ctx context.Context, stableModel *alcorv1alpha1.StableModelApplyConfiguration, opts v1.ApplyOptions) (result *v1alpha1.StableModel, err error) {
	if stableModel == nil {
		return nil, fmt.Errorf("stableModel provided to Apply must not be nil")
	}
	patchOpts := opts.ToPatchOptions()
	data, err := json.Marshal(stableModel)
	if err != nil {
		return nil, err
	}

	name := stableModel.Name
	if name == nil {
		return nil, fmt.Errorf("stableModel.Name must be provided to Apply")
	}

	result = &v1alpha1.StableModel{}
	err = c.client.Patch(types.ApplyPatchType).
		Namespace(c.ns).
		Resource("stablemodels").
		Name(*name).
		SubResource("status").
		VersionedParams(&patchOpts, scheme.ParameterCodec).
		Body(data).
		Do(ctx).
		Into(result)
	return
}
