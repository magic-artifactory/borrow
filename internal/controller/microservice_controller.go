/*
Copyright 2026 magic-artifactory.

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

package controller

import (
	"context"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	brwv1alpha1 "github.com/magic-artifactory/borrow/api/v1alpha1"
	"github.com/magic-artifactory/borrow/pkg/operatools"
	"github.com/magic-artifactory/borrow/pkg/utils"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
)

// MicroserviceReconciler reconciles a Microservice object
type MicroserviceReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=brw.soigneux.works,resources=microservices,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=brw.soigneux.works,resources=microservices/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=brw.soigneux.works,resources=microservices/finalizers,verbs=update
// +kubebuilder:rbac:groups=brw.soigneux.works,resources=pendingmicroservices,verbs=get;list;watch;delete

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Microservice object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.21.0/pkg/reconcile
func (r *MicroserviceReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	return operatools.Reconcile(
		r,
		r.Client,
		req,
		ctx,
		&brwv1alpha1.Microservice{},
	)
}

// SetupWithManager sets up the controller with the Manager.
func (r *MicroserviceReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&brwv1alpha1.Microservice{}).
		Named("microservice").
		Complete(r)
}

func (r *MicroserviceReconciler) addDeployment(conciliator *operatools.Conciliator, ms *brwv1alpha1.Microservice) error {
	image := ms.Spec.Image

	replicas := int32(1)
	if ms.Spec.Replicas != nil {
		replicas = *ms.Spec.Replicas
	}

	labels := map[string]string{
		"app.kubernetes.io/name": ms.GetName(),
	}

	deploy := appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      ms.GetName(),
			Namespace: ms.GetNamespace(),
			Labels:    labels,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  ms.GetName(),
							Image: image,
						},
					},
				},
			},
		},
	}

	if err := r.Create(conciliator.Ctx, &deploy); err != nil {
		if apierrors.IsAlreadyExists(err) {
			if err := r.Update(conciliator.Ctx, &deploy); err != nil {
				return err
			}
		}
		return err
	}

	return nil
}

func (r *MicroserviceReconciler) Process(conciliator *operatools.Conciliator) error {
	ms := conciliator.Object.(*brwv1alpha1.Microservice)
	if err := utils.ObjectNil(ms); err != nil {
		return err
	}

	if err := r.addDeployment(conciliator, ms); err != nil {
		return err
	}

	return nil
}
