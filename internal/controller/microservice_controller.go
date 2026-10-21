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

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	brwv1alpha1 "github.com/magic-artifactory/borrow/api/v1alpha1"
	"github.com/magic-artifactory/borrow/pkg/operatools"
	"github.com/magic-artifactory/borrow/pkg/utils"
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
		Complete(r)
}

// func (r *MicroserviceReconciler) ProcessPendingMs(conciliator *operatools.Conciliator, pendingObject *brwv1alpha1.PendingMicroservice) error {
// 	var req = &http.Client{Timeout: 10 * time.Second}
// 	patch := pendingObject.DeepCopy()

// 	done := make(map[string]bool)
// 	for _, res := range slices.Clone(pendingObject.Status.Resources) {
// 		resp, err := req.Get(res.Charts)
// 		fmt.Println(res.Charts)
// 		if err != nil {
// 			continue
// 		}
// 		resp.Body.Close()
// 		done[res.Uid] = true
// 	}

// 	if len(done) == 0 {
// 		return nil
// 	}

// 	pendingObject.Status.Resources = slices.DeleteFunc(pendingObject.Status.Resources, func(res utils.PendingResource) bool {
// 		return done[res.Uid]
// 	})

// 	return operatools.PatchStatusResource(conciliator.Ctx, conciliator.Client, pendingObject, patch)
// }

func (r *MicroserviceReconciler) Process(conciliator *operatools.Conciliator) error {
	msObject := conciliator.Object.(*brwv1alpha1.Microservice)
	if err := utils.ObjectNil(msObject); err != nil {
		return err
	}

	// var pendingList brwv1alpha1.PendingMicroserviceList
	// if err := r.Client.List(conciliator.Ctx, &pendingList, client.InNamespace(msObject.GetNamespace())); err != nil {
	// 	return nil
	// }

	// for _, pendingItem := range pendingList.Items {
	// 	if err := r.ProcessPendingMs(conciliator, &pendingItem); err != nil {
	// 		return err
	// 	}
	// }

	return nil
}
