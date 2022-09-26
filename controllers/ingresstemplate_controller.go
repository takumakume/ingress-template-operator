/*
Copyright 2022.

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

package controllers

import (
	"context"
	"fmt"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"

	ingresstemplatev1alpha1 "github.com/takumakume/ingress-template-operator/api/v1alpha1"
	"github.com/takumakume/ingress-template-operator/pkg/render"
	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	managedByValue = "ingress-template-controller"
)

// IngressTemplateReconciler reconciles a IngressTemplate object
type IngressTemplateReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=ingress-template.takumakume.github.io,resources=ingresstemplates,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=ingress-template.takumakume.github.io,resources=ingresstemplates/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=ingress-template.takumakume.github.io,resources=ingresstemplates/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the IngressTemplate object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.13.0/pkg/reconcile
func (r *IngressTemplateReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)

	log.Info("Run reconcile")
	ingresstemplate := &ingresstemplatev1alpha1.IngressTemplate{}
	if err := r.Get(ctx, req.NamespacedName, ingresstemplate); err != nil {
		if apierrors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}

		log.Error(err, "unable to fetch IngressTemplate")
		return ctrl.Result{}, err
	}

	ingress, err := ingressTemplateToIngress(ingresstemplate)
	if err != nil {
		return ctrl.Result{}, err
	}

	log.Info("Run create or update Ingress")
	op, err := ctrl.CreateOrUpdate(ctx, r.Client, ingress, func() error {
		if ingress.Labels == nil {
			ingress.Labels = map[string]string{}
		}
		ingress.Labels["app.kubernetes.io/managed-by"] = managedByValue
		return nil
	})
	if err != nil {
		log.Error(err, "unable to create or update Ingress")
		if statusUpdateErr := r.Update(ctx, ingresstemplate); statusUpdateErr != nil {
			return ctrl.Result{}, statusUpdateErr
		}
		return ctrl.Result{}, err
	}
	log.Info(fmt.Sprintf("Update success status:%s", op))

	if op != controllerutil.OperationResultNone {
		ingresstemplate.Status.Ready = corev1.ConditionTrue
		if statusUpdateErr := r.Update(ctx, ingresstemplate); statusUpdateErr != nil {
			return ctrl.Result{}, statusUpdateErr
		}
	}

	return ctrl.Result{Requeue: true}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *IngressTemplateReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&ingresstemplatev1alpha1.IngressTemplate{}).
		Owns(&networkingv1.Ingress{}).
		Complete(r)
}

func ingressTemplateToIngress(ingresstemplate *ingresstemplatev1alpha1.IngressTemplate) (*networkingv1.Ingress, error) {
	generated := &networkingv1.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			GenerateName: ingresstemplate.Name,
			Namespace:    ingresstemplate.Namespace,
			Annotations:  ingresstemplate.Spec.IngressAnnotations,
			Labels:       ingresstemplate.Spec.IngressLabels,
		},
		Spec: ingresstemplate.Spec.IngressSpecTemplate,
	}

	opt := render.Options{
		Namespace: ingresstemplate.Namespace,
	}

	generated, err := render.Render(generated, opt)
	if err != nil {
		return nil, err
	}

	return generated, nil
}
