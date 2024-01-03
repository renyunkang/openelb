package lb

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/util/workqueue"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

var deAndDsEnqueueLog = ctrl.Log.WithName("eventhandler").WithName("EnqueueRequestForDeAndDs")

// Enqueue requests for Deployments and DaemonSets type
// Only OpenELB NodeProxy needs this
type EnqueueRequestForDeAndDs struct {
	client.Client
}

// Get all OpenELB NodeProxy Services to reconcile them later
// These Services will be exposed by Proxy Pod
func (e *EnqueueRequestForDeAndDs) getServices() []corev1.Service {
	var svcs corev1.ServiceList

	if err := e.List(context.Background(), &svcs); err != nil {
		deAndDsEnqueueLog.Error(err, "Failed to list services")
		return nil
	}

	var result []corev1.Service
	for _, svc := range svcs.Items {
		if IsOpenELBNPService(&svc) {
			result = append(result, svc)
		}
	}

	return result
}

// Create implements EventHandler
func (e *EnqueueRequestForDeAndDs) Create(evt event.CreateEvent, q workqueue.RateLimitingInterface) {
	if evt.Object == nil {
		deAndDsEnqueueLog.Error(nil, "CreateEvent received with no metadata", "event", evt)
		return
	}

	for _, svc := range e.getServices() {
		q.Add(reconcile.Request{NamespacedName: types.NamespacedName{
			Name:      svc.GetName(),
			Namespace: svc.GetNamespace(),
		}})
	}
}

// Update implements EventHandler
func (e *EnqueueRequestForDeAndDs) Update(evt event.UpdateEvent, q workqueue.RateLimitingInterface) {
	if evt.ObjectOld == nil {
		deAndDsEnqueueLog.Error(nil, "UpdateEvent received with no old metadata", "event", evt)
	}

	if evt.ObjectNew == nil {
		deAndDsEnqueueLog.Error(nil, "UpdateEvent received with no new metadata", "event", evt)
	}

	for _, svc := range e.getServices() {
		q.Add(reconcile.Request{NamespacedName: types.NamespacedName{
			Name:      svc.GetName(),
			Namespace: svc.GetNamespace(),
		}})
	}
}

// Delete implements EventHandler
func (e *EnqueueRequestForDeAndDs) Delete(evt event.DeleteEvent, q workqueue.RateLimitingInterface) {
	if evt.Object == nil {
		deAndDsEnqueueLog.Error(nil, "DeleteEvent received with no metadata", "event", evt)
		return
	}
	for _, svc := range e.getServices() {
		q.Add(reconcile.Request{NamespacedName: types.NamespacedName{
			Name:      svc.GetName(),
			Namespace: svc.GetNamespace(),
		}})
	}
}

// Generic implements EventHandler
func (e *EnqueueRequestForDeAndDs) Generic(evt event.GenericEvent, q workqueue.RateLimitingInterface) {

}
