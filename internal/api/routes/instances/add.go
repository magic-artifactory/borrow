package instances

import (
	"context"
	"encoding/json"

	"github.com/lukas-sgx/borrow/api/v1alpha1"
	"github.com/lukas-sgx/borrow/pkg/operatools"

	"github.com/lukas-sgx/borrow/pkg/utils"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func patchStatusRessource(ctx context.Context, c client.Client, obj client.Object, patch client.Object) error {
	mergePatch := client.MergeFrom(patch)

	if err := c.Status().Patch(ctx, obj, mergePatch); err != nil {
		return err
	}
	return nil
}

func checkDuplicatePendingRessource(items []utils.PendingRessource, data *utils.PendingRessource) bool {
	for _, item := range items {
		if item.ProjectName == data.ProjectName {
			return true
		}
	}

	return false
}

func addInstance(conciliotor *operatools.Conciliator, body []byte) (bool, error) {
	var data utils.PendingRessource

	microservicePending := conciliotor.Object.(*v1alpha1.PendingMicroservice)

	if err := utils.ObjectNil(microservicePending); err != nil {
		return false, err
	}

	if err := json.Unmarshal(body, &data); err != nil {
		return false, err
	}

	if checkDuplicatePendingRessource(microservicePending.Status.Ressources, &data) {
		return false, nil
	}

	patch := microservicePending.DeepCopy()
	microservicePending.Status.Ressources = append(microservicePending.Status.Ressources, data)

	if err := patchStatusRessource(
		conciliotor.Ctx,
		conciliotor.Client,
		microservicePending,
		patch,
	); err != nil {
		return false, err
	}

	return true, nil
}
