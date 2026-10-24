package operatools

import (
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func GetClient(conciliator *Conciliator, name string, namespace string, object client.Object) {
	resource := types.NamespacedName{
		Name:      name,
		Namespace: namespace,
	}

	if err := conciliator.Client.Get(conciliator.Ctx, resource, object); err != nil {
		conciliator.Log.Error(err, "")
		return
	}

	conciliator.Object = object
}
