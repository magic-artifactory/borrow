package instances

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/lukas-sgx/borrow/api/v1alpha1"
	"github.com/lukas-sgx/borrow/pkg/operatools"
)

type Endpoints struct {
	ConciliatorMicroservice *operatools.Conciliator
	ConciliatorPending      *operatools.Conciliator
}

func NewEndpoints(conciliator *operatools.Conciliator) *Endpoints {
	conciliatorPending := new(operatools.Conciliator)
	*conciliatorPending = *conciliator

	operatools.GetClient(conciliator, "microservice-sample", "sre-system", &v1alpha1.Microservice{})
	operatools.GetClient(conciliatorPending, "pendingmicroservice-sample", "sre-system", &v1alpha1.PendingMicroservice{})

	return &Endpoints{
		ConciliatorPending:      conciliatorPending,
		ConciliatorMicroservice: conciliator,
	}
}

func (e *Endpoints) Add(response http.ResponseWriter, req *http.Request) {
	var status string

	state := false

	body, err := io.ReadAll(req.Body)
	if err != nil {
		internalServerError(response, err)
		return
	}

	ok, err := addInstance(e.ConciliatorPending, body)
	if err != nil {
		internalServerError(response, err)
		return
	}

	response.Header().Set("Content-type", "application/json")
	response.WriteHeader(http.StatusOK)

	if ok {
		state = true
		status = "add to pending"
	} else {
		status = "duplicate"
	}

	if err := json.NewEncoder(response).Encode(map[string]string{"creation": strconv.FormatBool(state), "status": status}); err != nil {
		internalServerError(response, err)
	}
}

func (e *Endpoints) Delete(response http.ResponseWriter, req *http.Request) {
	response.Header().Set("Content-type", "application/json")
	response.WriteHeader(http.StatusOK)
}

func (e *Endpoints) Get(response http.ResponseWriter, req *http.Request) {
	path := strings.Split(req.URL.EscapedPath(), "/")

	if len(path) != 2 {
		internalServerError(response, errors.New("malformed request"))
		return
	}
	instanceName := path[1]

	obj, ok := e.ConciliatorMicroservice.Object.(*v1alpha1.Microservice)
	if !ok {
		internalServerError(response, errors.New("v1alpha1.Microservice: object nil"))
		return
	}

	data, err := json.Marshal(getInstanceInfo(obj, instanceName))
	if err != nil {
		internalServerError(response, err)
		return
	}

	response.Header().Set("Content-type", "application/json")
	response.WriteHeader(http.StatusOK)

	if _, err := response.Write(data); err != nil {
		internalServerError(response, err)
	}
}

func internalServerError(response http.ResponseWriter, err error) {
	response.Header().Set("Content-type", "application/json")
	response.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(response).Encode(map[string]string{"error": err.Error()}); err != nil {
		internalServerError(response, err)
	}
}
