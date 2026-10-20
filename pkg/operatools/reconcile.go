package operatools

import (
	"context"

	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type Reconciler interface {
	client.Reader
	Process(conciliator *Conciliator) error
}

func processRessource(conciliator *Conciliator, req ctrl.Request, r Reconciler) error {
	conciliator.Log.Info("Processing Ressource",
		"namespace", req.Namespace,
		"name", req.Name,
	)

	return r.Process(conciliator)
}

func Reconcile(r Reconciler, c client.Client, req ctrl.Request,
	ctx context.Context, object client.Object) (ctrl.Result, error) {
	if err := r.Get(ctx, req.NamespacedName, object); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	conciliator := NewConciliator(ctrl.Log, object, c, ctx)

	if err := processRessource(conciliator, req, r); err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

func PatchStatusRessource(ctx context.Context, c client.Client, obj client.Object, patch client.Object) error {
	mergePatch := client.MergeFrom(patch)

	if err := c.Status().Patch(ctx, obj, mergePatch); err != nil {
		return err
	}
	return nil
}