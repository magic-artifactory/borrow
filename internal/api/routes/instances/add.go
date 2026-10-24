package instances

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/magic-artifactory/borrow/api/v1alpha1"
	"github.com/magic-artifactory/borrow/pkg/operatools"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/magic-artifactory/borrow/pkg/utils"
)

func checkDuplicateMsResource(msList *v1alpha1.MicroserviceList, data *utils.SpecPendingResource) bool {
	for _, ms := range msList.Items {
		if ms.Spec.ProjectName == data.ProjectName || ms.Spec.Uid == data.Uuid {
			return true
		}
	}
	return false
}

func checkDuplicatePendingResource(pendingMsList *v1alpha1.PendingMicroserviceList, data *utils.SpecPendingResource) bool {
	for _, item := range pendingMsList.Items {
		if item.GetName() == data.ProjectName || item.Spec.Uuid == data.Uuid {
			return true
		}
	}
	return false
}

func addInstance(conciliator *operatools.Conciliator, body []byte) (*utils.SpecPendingResource, error) {
	var data utils.SpecPendingResource
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, fmt.Errorf("invalid request body: %w", err)
	}
	if data.ProjectName == "" {
		return nil, errors.New("projectName is required")
	}

	var msList v1alpha1.MicroserviceList
	if err := conciliator.Client.List(conciliator.Ctx, &msList); err != nil {
		return nil, fmt.Errorf("listing microservices: %w", err)
	}

	var pendingMsList v1alpha1.PendingMicroserviceList
	if err := conciliator.Client.List(conciliator.Ctx, &pendingMsList); err != nil {
		return nil, fmt.Errorf("listing pending microservices: %w", err)
	}

	if checkDuplicatePendingResource(&pendingMsList, &data) || checkDuplicateMsResource(&msList, &data) {
		return nil, errors.New(utils.ErrDuplicateResource)
	}

	data.Uuid = uuid.New().String()

	newPending := v1alpha1.PendingMicroservice{
		ObjectMeta: metav1.ObjectMeta{
			Name:      data.ProjectName,
			Namespace: utils.NamespaceSRE,
		},
		Spec: v1alpha1.PendingMicroserviceSpec{
			Charts: data.Charts,
			Uuid:   data.Uuid,
		},
	}

	if err := conciliator.Client.Create(conciliator.Ctx, &newPending); err != nil {
		if apierrors.IsAlreadyExists(err) {
			return nil, errors.New(utils.ErrDuplicateResource)
		}
		return nil, fmt.Errorf("creating pending microservice %q: %w", data.ProjectName, err)
	}

	return &data, nil
}
