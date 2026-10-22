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
	"fmt"
	"io"
	"net/http"
	"time"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	brwv1alpha1 "github.com/magic-artifactory/borrow/api/v1alpha1"
	"github.com/magic-artifactory/borrow/pkg/operatools"
	"github.com/magic-artifactory/borrow/pkg/utils"
)

// PendingMicroserviceReconciler reconciles a PendingMicroservice object
type PendingMicroserviceReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=brw.soigneux.works,resources=pendingmicroservices,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=brw.soigneux.works,resources=pendingmicroservices/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=brw.soigneux.works,resources=pendingmicroservices/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the PendingMicroservice object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.21.0/pkg/reconcile
func (r *PendingMicroserviceReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	return operatools.Reconcile(
		r,
		r.Client,
		req,
		ctx,
		&brwv1alpha1.PendingMicroservice{},
	)
}

// SetupWithManager sets up the controller with the Manager.
func (r *PendingMicroserviceReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&brwv1alpha1.PendingMicroservice{}).
		Named("pendingmicroservice").
		Complete(r)
}

func (r *PendingMicroserviceReconciler) processPendingMS(conciliator *operatools.Conciliator, pendingMs *brwv1alpha1.PendingMicroservice) error {
	clientHttp := http.Client{
		Timeout: 10 * time.Second,
	}

	response, err := clientHttp.Get(pendingMs.Spec.Charts)
	if err != nil {
		return err
	}

	defer func() {
		if err := response.Body.Close(); err != nil {
			conciliator.Log.Error(err, "")
		}
	}()

	if response.StatusCode < 200 || response.StatusCode >= 300 {
		body, _ := io.ReadAll(io.LimitReader(response.Body, 2048))
		return fmt.Errorf("unexpected status %d (%s): %s",
			response.StatusCode, http.StatusText(response.StatusCode), body)
	}

	return nil
}

func (r *PendingMicroserviceReconciler) Process(conciliator *operatools.Conciliator) error {
	object := conciliator.Object.(*brwv1alpha1.PendingMicroservice)
	if err := utils.ObjectNil(object); err != nil {
		return err
	}

	if err := r.processPendingMS(conciliator, object); err != nil {
		for range 3 {
			if err := r.Delete(conciliator.Ctx, object); err == nil {
				break
			} else {
				time.Sleep(5 * time.Second)
				continue
			}
		}
	}

	return nil
}
