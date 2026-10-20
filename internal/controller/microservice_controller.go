/*
Copyright 2026 lukas-sgx.

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
	"net/http"
	"slices"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	"github.com/lukas-sgx/borrow/api/v1alpha1"
	brwv1alpha1 "github.com/lukas-sgx/borrow/api/v1alpha1"
	"github.com/lukas-sgx/borrow/pkg/operatools"
	"github.com/lukas-sgx/borrow/pkg/utils"
)

// MicroserviceReconciler reconciles a Microservice object
type MicroserviceReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=brw.soigneux.works,resources=microservices,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=brw.soigneux.works,resources=microservices/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=brw.soigneux.works,resources=microservices/finalizers,verbs=update
// +kubebuilder:rbac:groups=brw.soigneux.works,resources=pendingmicroservices,verbs=get;list;watch
// +kubebuilder:rbac:groups=brw.soigneux.works,resources=pendingmicroservices/status,verbs=get;update;patch

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
		Watches(
			&v1alpha1.PendingMicroservice{},
			handler.EnqueueRequestsFromMapFunc(r.pendingToMicroservice),
		).
		Complete(r)
}

func (r *MicroserviceReconciler) pendingToMicroservice(ctx context.Context, obj client.Object) []reconcile.Request {
	pending, ok := obj.(*v1alpha1.PendingMicroservice)
	if !ok {
		return nil
	}

	return []reconcile.Request{{
		NamespacedName: types.NamespacedName{
			Name:      "microservice-sample",
			Namespace: pending.Namespace,
		},
	}}
}

func deletePendingObject(conciliator *operatools.Conciliator, pendingObject *brwv1alpha1.PendingMicroservice, uid string) error {
	patch := pendingObject.DeepCopy()
	pendingObject.Status.Resources = slices.DeleteFunc(pendingObject.Status.Resources, func(res utils.PendingResource) bool {
		return uid == res.Uid
	})

	if err := operatools.PatchStatusResource(
		conciliator.Ctx,
		conciliator.Client,
		pendingObject,
		patch,
	); err != nil {
		return err
	}

	return nil
}

func (r *MicroserviceReconciler) ProcessPendingMs(conciliator *operatools.Conciliator, pendingObject *brwv1alpha1.PendingMicroservice) error {
	for _, object := range pendingObject.Status.Resources {
		_, err := http.Get(object.Charts)
		if err != nil {
			deletePendingObject(conciliator, pendingObject, object.Uid)
		}
	}
	return nil
}

func (r *MicroserviceReconciler) Process(conciliator *operatools.Conciliator) error {
	msObject := conciliator.Object.(*brwv1alpha1.Microservice)

	pendingObject := v1alpha1.PendingMicroservice{}
	Resource := types.NamespacedName{
		Name:      "pendingmicroservice-sample",
		Namespace: utils.NamespaceSRE,
	}

	return utils.NewRunner().
		Run(func() error { return utils.ObjectNil(msObject) }).
		Run(func() error { return conciliator.Client.Get(conciliator.Ctx, Resource, &pendingObject) }).
		Run(func() error { return r.ProcessPendingMs(conciliator, &pendingObject) }).
		Err()
}
