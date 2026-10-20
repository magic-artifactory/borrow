package instances

import (
	"encoding/json"

	"github.com/google/uuid"
	"github.com/magic-artifactory/borrow/api/v1alpha1"
	"github.com/magic-artifactory/borrow/pkg/operatools"

	"github.com/magic-artifactory/borrow/pkg/utils"
)

func checkDuplicateMsResource(MsList *v1alpha1.MicroserviceList, data *utils.PendingResource) bool {
	for _, ms := range MsList.Items {
		if ms.Status.Spec.Name == data.ProjectName || ms.Status.Spec.Uid == data.Uid {
			return true
		}
	}
	return false
}

func checkDuplicatePendingResource(items []utils.PendingResource, data *utils.PendingResource) bool {
	for _, item := range items {
		if item.ProjectName == data.ProjectName || item.Uid == data.Uid {
			return true
		}
	}

	return false
}

func addInstance(conciliotor *operatools.Conciliator, MsList *v1alpha1.MicroserviceList, body []byte) (*utils.PendingResource, error) {
	var data utils.PendingResource

	microservicePending := conciliotor.Object.(*v1alpha1.PendingMicroservice)

	if err := utils.ObjectNil(microservicePending); err != nil {
		return nil, err
	}

	if err := json.Unmarshal(body, &data); err != nil {
		return nil, err
	}

	data.Uid = uuid.New().String()

	if checkDuplicatePendingResource(microservicePending.Status.Resources, &data) {
		return nil, nil
	}

	if checkDuplicateMsResource(MsList, &data) {
		return nil, nil
	}

	patch := microservicePending.DeepCopy()
	microservicePending.Status.Resources = append(microservicePending.Status.Resources, data)

	if err := operatools.PatchStatusResource(
		conciliotor.Ctx,
		conciliotor.Client,
		microservicePending,
		patch,
	); err != nil {
		return nil, err
	}

	return &data, nil
}
