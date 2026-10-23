package instances

import (
	"encoding/json"

	"github.com/google/uuid"
	"github.com/magic-artifactory/borrow/api/v1alpha1"
	"github.com/magic-artifactory/borrow/pkg/operatools"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/magic-artifactory/borrow/pkg/utils"
)

func checkDuplicateMsResource(msList *v1alpha1.MicroserviceList, data *utils.SpecPendingResource) bool {
	for _, ms := range msList.Items {
		if ms.Spec.ProjectName == data.ProjectName || ms.Spec.Uid == data.Uid {
			return true
		}
	}
	return false
}

func checkDuplicatePendingResource(pendingMsList *v1alpha1.PendingMicroserviceList, data *utils.SpecPendingResource) bool {
	for _, item := range pendingMsList.Items {
		if item.GetName() == data.ProjectName || item.Spec.Uuid == data.Uid {
			return true
		}
	}
	return false
}

func addInstance(conciliator *operatools.Conciliator, body []byte) (*utils.SpecPendingResource, error) {
	var data utils.SpecPendingResource
	var pendingMsList v1alpha1.PendingMicroserviceList
	var msList v1alpha1.MicroserviceList

	if err := conciliator.Client.List(conciliator.Ctx, &msList); err != nil {
		return nil, err
	}

	if err := conciliator.Client.List(conciliator.Ctx, &pendingMsList); err != nil {
		return nil, err
	}

	if err := json.Unmarshal(body, &data); err != nil {
		return nil, err
	}

	data.Uid = uuid.New().String()

	if checkDuplicatePendingResource(&pendingMsList, &data) || checkDuplicateMsResource(&msList, &data) {
		return nil, nil
	}

	newPending := v1alpha1.PendingMicroservice{
		ObjectMeta: metav1.ObjectMeta{
			Name:      data.ProjectName,
			Namespace: utils.NamespaceSRE,
		},
	}

	newPending.Spec.Charts = data.Charts
	newPending.Spec.Uuid = data.Uid

	if err := conciliator.Client.Create(conciliator.Ctx, &newPending); err != nil {
		return nil, err
	}

	return &data, nil
}
