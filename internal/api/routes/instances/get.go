package instances

import (
	"github.com/magic-artifactory/borrow/api/v1alpha1"
)

type Instance struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
}

func getInstanceInfo(MsList *v1alpha1.MicroserviceList, uid string) *Instance {
	for _, ms := range MsList.Items {
		if ms.Spec.Uid == uid {
			return &Instance{
				Name:      ms.GetName(),
				Namespace: ms.GetNamespace(),
			}
		}
	}

	return nil
}
