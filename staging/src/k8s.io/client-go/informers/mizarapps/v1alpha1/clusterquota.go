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

// Code generated by informer-gen. DO NOT EDIT.

package v1alpha1

import (
	"context"
	time "time"

	mizarappsv1alpha1 "k8s.io/api/mizarapps/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
	watch "k8s.io/apimachinery/pkg/watch"
	internalinterfaces "k8s.io/client-go/informers/internalinterfaces"
	kubernetes "k8s.io/client-go/kubernetes"
	v1alpha1 "k8s.io/client-go/listers/mizarapps/v1alpha1"
	cache "k8s.io/client-go/tools/cache"
)

// ClusterQuotaInformer provides access to a shared informer and lister for
// ClusterQuotas.
type ClusterQuotaInformer interface {
	Informer() cache.SharedIndexInformer
	Lister() v1alpha1.ClusterQuotaLister
}

type clusterQuotaInformer struct {
	factory          internalinterfaces.SharedInformerFactory
	tweakListOptions internalinterfaces.TweakListOptionsFunc
	namespace        string
}

// NewClusterQuotaInformer constructs a new informer for ClusterQuota type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewClusterQuotaInformer(client kubernetes.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers) cache.SharedIndexInformer {
	return NewFilteredClusterQuotaInformer(client, namespace, resyncPeriod, indexers, nil)
}

// NewFilteredClusterQuotaInformer constructs a new informer for ClusterQuota type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewFilteredClusterQuotaInformer(client kubernetes.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers, tweakListOptions internalinterfaces.TweakListOptionsFunc) cache.SharedIndexInformer {
	return cache.NewSharedIndexInformer(
		&cache.ListWatch{
			ListFunc: func(options v1.ListOptions) (runtime.Object, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.MizarappsV1alpha1().ClusterQuotas(namespace).List(context.TODO(), options)
			},
			WatchFunc: func(options v1.ListOptions) (watch.Interface, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.MizarappsV1alpha1().ClusterQuotas(namespace).Watch(context.TODO(), options)
			},
		},
		&mizarappsv1alpha1.ClusterQuota{},
		resyncPeriod,
		indexers,
	)
}

func (f *clusterQuotaInformer) defaultInformer(client kubernetes.Interface, resyncPeriod time.Duration) cache.SharedIndexInformer {
	return NewFilteredClusterQuotaInformer(client, f.namespace, resyncPeriod, cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc}, f.tweakListOptions)
}

func (f *clusterQuotaInformer) Informer() cache.SharedIndexInformer {
	return f.factory.InformerFor(&mizarappsv1alpha1.ClusterQuota{}, f.defaultInformer)
}

func (f *clusterQuotaInformer) Lister() v1alpha1.ClusterQuotaLister {
	return v1alpha1.NewClusterQuotaLister(f.Informer().GetIndexer())
}
