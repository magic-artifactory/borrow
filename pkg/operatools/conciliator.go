package operatools

import (
	"context"

	"github.com/go-logr/logr"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type Conciliator struct {
	Log    logr.Logger
	Object client.Object
	Client client.Client
	Ctx    context.Context
}

func NewConciliator(log logr.Logger, object client.Object, c client.Client, ctx context.Context) *Conciliator {
	return &Conciliator{
		Log:    log,
		Object: object,
		Client: c,
		Ctx:    ctx,
	}
}
