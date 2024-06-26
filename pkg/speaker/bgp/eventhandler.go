package bgp

import (
	"context"

	"github.com/openelb/openelb/api/v1alpha2"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/util/workqueue"
	"k8s.io/klog/v2"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

type EnqueueRequestForNode struct {
	client.Client
	peer bool
}

func (e *EnqueueRequestForNode) getDefaultBgpConf() []v1alpha2.BgpConf {
	var def v1alpha2.BgpConf

	if err := e.Get(context.Background(), client.ObjectKey{Name: "default"}, &def); err != nil {
		return nil
	}

	return []v1alpha2.BgpConf{def}
}

func (e *EnqueueRequestForNode) getBgpPeers() []v1alpha2.BgpPeer {
	var peers v1alpha2.BgpPeerList

	if err := e.List(context.Background(), &peers); err != nil {
		return nil
	}

	return peers.Items
}

// Create implements EventHandler
func (e *EnqueueRequestForNode) Create(ctx context.Context, evt event.CreateEvent, q workqueue.RateLimitingInterface) {
}

// Update implements EventHandler
func (e *EnqueueRequestForNode) Update(ctx context.Context, evt event.UpdateEvent, q workqueue.RateLimitingInterface) {
	if evt.ObjectOld == nil {
		klog.Error("UpdateEvent received with no old metadata", "event", evt)
	}

	if evt.ObjectNew == nil {
		klog.Error("UpdateEvent received with no new metadata", "event", evt)
	}

	if !e.peer {
		for _, svc := range e.getDefaultBgpConf() {
			q.Add(reconcile.Request{NamespacedName: types.NamespacedName{
				Name:      svc.GetName(),
				Namespace: svc.GetNamespace(),
			}})
		}
	} else {
		for _, svc := range e.getBgpPeers() {
			q.Add(reconcile.Request{NamespacedName: types.NamespacedName{
				Name:      svc.GetName(),
				Namespace: svc.GetNamespace(),
			}})
		}
	}
}

// Delete implements EventHandler
func (e *EnqueueRequestForNode) Delete(ctx context.Context, evt event.DeleteEvent, q workqueue.RateLimitingInterface) {
}

// Generic implements EventHandler
func (e *EnqueueRequestForNode) Generic(ctx context.Context, evt event.GenericEvent, q workqueue.RateLimitingInterface) {

}
