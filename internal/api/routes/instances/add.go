package instances

import (
	"context"
	"encoding/json"

	"github.com/google/uuid"
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

func checkDuplicateMsRessource(MsList *v1alpha1.MicroserviceList, data *utils.PendingRessource) bool {
	for _, ms := range MsList.Items {
		if ms.Status.Spec.Name == data.ProjectName || ms.Status.Spec.Uid == data.Uid {
			return true
		}
	}
	return false
}

func checkDuplicatePendingRessource(items []utils.PendingRessource, data *utils.PendingRessource) bool {
	for _, item := range items {
		if item.ProjectName == data.ProjectName || item.Uid == data.Uid {
			return true
		}
	}

	return false
}

func addInstance(conciliotor *operatools.Conciliator, MsList *v1alpha1.MicroserviceList, body []byte) (*utils.PendingRessource, error) {
	var data utils.PendingRessource

	microservicePending := conciliotor.Object.(*v1alpha1.PendingMicroservice)

	if err := utils.ObjectNil(microservicePending); err != nil {
		return nil, err
	}

	if err := json.Unmarshal(body, &data); err != nil {
		return nil, err
	}

	data.Uid = uuid.New().String()

	if checkDuplicatePendingRessource(microservicePending.Status.Ressources, &data) {
		return nil, nil
	}

	if checkDuplicateMsRessource(MsList, &data) {
		return nil, nil
	}

	patch := microservicePending.DeepCopy()
	microservicePending.Status.Ressources = append(microservicePending.Status.Ressources, data)

	if err := patchStatusRessource(
		conciliotor.Ctx,
		conciliotor.Client,
		microservicePending,
		patch,
	); err != nil {
		return nil, err
	}

	return &data, nil
}
