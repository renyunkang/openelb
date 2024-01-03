/*
Copyright 2020 The Kubesphere Authors.

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

package speaker

import (
	"context"
	"reflect"

	"github.com/openelb/openelb/api/v1alpha2"
	"github.com/openelb/openelb/pkg/constant"
	"github.com/openelb/openelb/pkg/util"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

// LBReconciler reconciles a service object
type LBReconciler struct {
	client.Client
	record.EventRecorder

	Handler func(*corev1.Service) error
}

func SetupLBReconciler(mgr ctrl.Manager, svchandler func(*corev1.Service) error) error {
	lb := &LBReconciler{
		Handler:       svchandler,
		Client:        mgr.GetClient(),
		EventRecorder: mgr.GetEventRecorderFor("lb"),
	}

	return lb.SetupWithManager(mgr)
}

func IsOpenELBService(obj runtime.Object) bool {
	svc, ok := obj.(*corev1.Service)
	if !ok {
		return false
	}

	if svc.Annotations == nil {
		return false
	}

	if svc.Spec.Type != corev1.ServiceTypeLoadBalancer {
		return false
	}

	if value, ok := svc.Annotations[constant.OpenELBAnnotationKey]; ok && value == constant.OpenELBAnnotationValue {
		return true
	}

	return false
}

func (r *LBReconciler) shouldReconcileEP(e metav1.Object) bool {
	if e.GetAnnotations() != nil {
		if e.GetAnnotations()["control-plane.alpha.kubernetes.io/leader"] != "" {
			return false
		}
	}

	svc := &corev1.Service{}
	err := r.Get(context.TODO(), types.NamespacedName{Namespace: e.GetNamespace(), Name: e.GetName()}, svc)
	if err != nil {
		if errors.IsNotFound(err) {
			return false
		}
		return true
	}

	return IsOpenELBService(svc)
}

func (r *LBReconciler) SetupWithManager(mgr ctrl.Manager) error {
	p := predicate.Funcs{
		CreateFunc: func(e event.CreateEvent) bool {
			return IsOpenELBService(e.Object)
		},
		DeleteFunc: func(e event.DeleteEvent) bool {
			return IsOpenELBService(e.Object)
		},
	}

	ctl, err := ctrl.NewControllerManagedBy(mgr).
		For(&corev1.Service{}).
		Owns(&corev1.Endpoints{}).
		WithEventFilter(p).
		Named("LBController").
		Build(r)
	if err != nil {
		return err
	}

	bp := predicate.Funcs{
		CreateFunc: func(e event.CreateEvent) bool {
			return false
		},
		DeleteFunc: func(e event.DeleteEvent) bool {
			return false
		},
		UpdateFunc: func(e event.UpdateEvent) bool {
			old := e.ObjectOld.(*v1alpha2.BgpConf)
			new := e.ObjectNew.(*v1alpha2.BgpConf)

			return !reflect.DeepEqual(old.Annotations, new.Annotations)
		},
	}
	err = ctl.Watch(&source.Kind{Type: &v1alpha2.BgpConf{}}, &EnqueueRequestForNode{Client: r.Client}, bp)
	if err != nil {
		return err
	}

	np := predicate.Funcs{
		UpdateFunc: func(e event.UpdateEvent) bool {
			if util.NodeReady(e.ObjectOld) != util.NodeReady(e.ObjectNew) {
				return true
			}
			if nodeAddrChange(e.ObjectOld, e.ObjectNew) {
				return true
			}
			return false
		},
		CreateFunc: func(e event.CreateEvent) bool {
			return false
		},
		DeleteFunc: func(e event.DeleteEvent) bool {
			return false
		},
	}
	err = ctl.Watch(&source.Kind{Type: &corev1.Node{}}, &EnqueueRequestForNode{Client: r.Client}, np)
	if err != nil {
		return err
	}

	//endpoints
	ep := predicate.Funcs{
		UpdateFunc: func(e event.UpdateEvent) bool {
			return r.shouldReconcileEP(e.ObjectNew)
		},
		CreateFunc: func(e event.CreateEvent) bool {
			return r.shouldReconcileEP(e.Object)
		},
		DeleteFunc: func(e event.DeleteEvent) bool {
			return r.shouldReconcileEP(e.Object)
		},
	}

	return ctl.Watch(&source.Kind{Type: &corev1.Endpoints{}}, &handler.EnqueueRequestForObject{}, ep)
}

//+kubebuilder:rbac:groups=network.kubesphere.io,resources=bgpconfs,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=network.kubesphere.io,resources=bgpconfs/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=network.kubesphere.io,resources=bgpconfs/finalizers,verbs=update
//+kubebuilder:rbac:groups="",resources=configmaps,verbs=get;list;watch

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster BgpConf CRD closer to the desired state.
func (l *LBReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := ctrl.Log.WithValues("service", req.NamespacedName)
	log.V(1).Info("start setup openelb service")
	defer log.V(1).Info("finish reconcile openelb service")

	svc := &corev1.Service{}
	if err := l.Client.Get(ctx, req.NamespacedName, svc); err != nil {
		if errors.IsNotFound(err) {
			log.Info("NotFound", " req.NamespacedName", req.NamespacedName)
			return ctrl.Result{}, nil
		}
		return ctrl.Result{Requeue: true}, err
	}

	if err := l.Handler(svc); err != nil {
		return ctrl.Result{Requeue: true}, err
	}
	return ctrl.Result{}, nil
}
