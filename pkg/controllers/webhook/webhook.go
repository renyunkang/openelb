package webhook

import (
	"context"
	"fmt"
	"reflect"

	"github.com/openelb/openelb/api/v1alpha2"
	"github.com/openelb/openelb/pkg/constant"
	"github.com/openelb/openelb/pkg/validate"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

var _ admission.CustomValidator = &EipWebhookHandler{}

type EipWebhookHandler struct {
	client.Client
}

func SetupWithManager(mgr manager.Manager) error {
	e := &EipWebhookHandler{Client: mgr.GetClient()}

	return builder.WebhookManagedBy(mgr).
		For(&v1alpha2.Eip{}).
		WithValidator(e).
		Complete()
}

func (e *EipWebhookHandler) validate(ctx context.Context, eip *v1alpha2.Eip, overlap bool) (admission.Warnings, error) {
	eips := &v1alpha2.EipList{}
	if err := e.Client.List(ctx, eips); err != nil {
		return nil, err
	}

	if overlap {
		for _, e := range eips.Items {
			if eip.Name == e.Name {
				continue
			}

			if eip.IsOverlap(e) {
				return nil, fmt.Errorf("eip address overlap with %s", eip.Name)
			}
		}
	}

	// validate default eip
	return nil, e.validateDefault(eips, eip)
}

func (e *EipWebhookHandler) validateDefault(eips *v1alpha2.EipList, eip *v1alpha2.Eip) error {
	if !validate.HasOpenELBDefaultEipAnnotation(eip.Annotations) {
		return nil
	}

	for _, e := range eips.Items {
		if e.Name == eip.Name {
			continue
		}

		if validate.HasOpenELBDefaultEipAnnotation(e.Annotations) {
			return fmt.Errorf("already exists a default EIP")
		}
	}

	return nil
}

// ValidateCreate implements admission.Validator so a webhook can intercept requests and reject them if they do not pass validation.
// validate the create request
func (e *EipWebhookHandler) ValidateCreate(ctx context.Context, obj runtime.Object) (admission.Warnings, error) {
	eip := obj.(*v1alpha2.Eip)
	_, _, err := eip.GetSize()
	if err != nil {
		return nil, err
	}

	if eip.Spec.Protocol == constant.OpenELBProtocolLayer2 && eip.Spec.Interface == "" {
		return nil, fmt.Errorf("when the protocol is layer2, the spec.interface should not be empty")
	}

	return e.validate(ctx, eip, true)
}

// ValidateUpdate implements admission.Validator so a webhook can intercept requests and reject them if they do not pass validation.
// validate the update request
func (e *EipWebhookHandler) ValidateUpdate(ctx context.Context, oldObj, newObj runtime.Object) (admission.Warnings, error) {
	newE := newObj.(*v1alpha2.Eip)
	oldE := oldObj.(*v1alpha2.Eip)

	if !reflect.DeepEqual(newE.Spec, oldE.Spec) {
		if newE.Spec.Address != oldE.Spec.Address {
			return nil, fmt.Errorf("the address field is not allowed to be modified")
		}
	}

	if !reflect.DeepEqual(newE.Annotations, oldE.Annotations) {
		return e.validate(ctx, newE, false)
	}

	return nil, nil
}

// ValidateDelete implements admission.Validator so a webhook can intercept requests and reject them if they do not pass validation.
// validate the delete request
func (e *EipWebhookHandler) ValidateDelete(ctx context.Context, obj runtime.Object) (admission.Warnings, error) {
	return nil, nil
}
