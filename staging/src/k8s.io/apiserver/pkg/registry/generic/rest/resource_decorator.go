/*
Copyright 2014 The Kubernetes Authors.

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

package rest

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/klog"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

// Resource decoder is a resource that can encode itself from kubelet for cluster
type ResourceDecoder struct {
	Location  *url.URL
	Transport http.RoundTripper
	Resource  runtime.Object
	// ResponseChecker HttpResponseChecker
	// RedirectChecker func(req *http.Request, via []*http.Request) error
	// TLSVerificationErrorCounter is an optional value that will Inc every time a TLS error is encountered.  This can
	// be wired a single prometheus counter instance to get counts overall.
	// TLSVerificationErrorCounter CounterMetric
}

func (obj *ResourceDecoder) GetObjectKind() schema.ObjectKind {
	return schema.EmptyObjectKind
}
func (obj *ResourceDecoder) DeepCopyObject() runtime.Object {
	panic("rest.LocationStreamer does not implement DeepCopyObject")
}

func NewResourceDecoder(location *url.URL, ctx context.Context, obj runtime.Object) (*ResourceDecoder, error) {
	if location == nil {
		return nil, fmt.Errorf("location is nil")
	}
	transport := http.DefaultTransport

	client := &http.Client{
		Transport: transport,
	}

	req, err := http.NewRequest("GET", location.String(), nil)

	if err != nil {
		return nil, fmt.Errorf("failed to construct request for %s, got %v", location.String(), err)
	}
	// Pass the parent context down to the request to ensure that the resources
	// will be release properly.
	req = req.WithContext(ctx)

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get response for %s, got %v", location.String(), err)
	}

	bodyBytes, err := readFromResp(resp)
	if err != nil {
		return nil, fmt.Errorf("failed to get bytes from request for %s, got %v", location.String(), err)
	}

	scheme := runtime.NewScheme()
	v1.SchemeBuilder.AddToScheme(scheme) //nolint:errcheck
	codecs := serializer.NewCodecFactory(scheme)
	decoder := codecs.UniversalDecoder(v1.SchemeGroupVersion)
	resource, gvk, err := decoder.Decode(bodyBytes, nil, obj)

	klog.V(4).Infof("decode output is %s %s %s", gvk.Group, gvk.Version, gvk.Kind)

	if err != nil {
		return nil, fmt.Errorf("failed to decode request for %s, got %v", location.String(), err)
	}

	return &ResourceDecoder{
		Location:  location,
		Transport: transport,
		Resource:  resource,
	}, nil
}

func readFromResp(resp *http.Response) ([]byte, error) {
	defer resp.Body.Close()
	bodyBytes, err := ioutil.ReadAll(resp.Body)

	return bodyBytes, err
}
