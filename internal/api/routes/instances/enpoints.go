package instances

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/magic-artifactory/borrow/api/v1alpha1"
	"github.com/magic-artifactory/borrow/pkg/operatools"
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
	var msList v1alpha1.MicroserviceList

	state := false

	body, err := io.ReadAll(req.Body)
	if err != nil {
		internalServerError(response, err)
		return
	}

	if err := e.ConciliatorMicroservice.Client.List(e.ConciliatorMicroservice.Ctx, &msList); err != nil {
		internalServerError(response, err)
		return
	}

	pendingInstance, err := addInstance(e.ConciliatorPending, &msList, body)
	if err != nil {
		internalServerError(response, err)
		return
	}

	response.Header().Set("Content-type", "application/json")
	response.WriteHeader(http.StatusOK)

	if pendingInstance != nil {
		state = true
		status = "add to pending"

		if err := json.NewEncoder(response).Encode(map[string]string{"creation": strconv.FormatBool(state), "status": status, "uuid": pendingInstance.Uid}); err != nil {
			internalServerError(response, err)
		}
	} else {
		status = "not created"
		if err := json.NewEncoder(response).Encode(map[string]string{"creation": strconv.FormatBool(state), "status": status}); err != nil {
			internalServerError(response, err)
		}
	}
}

func (e *Endpoints) Delete(response http.ResponseWriter, req *http.Request) {
	response.Header().Set("Content-type", "application/json")
	response.WriteHeader(http.StatusOK)
}

func (e *Endpoints) Get(response http.ResponseWriter, req *http.Request) {
	var msList v1alpha1.MicroserviceList

	path := strings.Split(req.URL.EscapedPath(), "/")

	if len(path) != 2 {
		internalServerError(response, errors.New("malformed request"))
		return
	}
	instanceName := path[1]

	if err := e.ConciliatorMicroservice.Client.List(e.ConciliatorMicroservice.Ctx, &msList); err != nil {
		internalServerError(response, err)
		return
	}

	data, err := json.Marshal(getInstanceInfo(&msList, instanceName))
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

	if err := json.NewEncoder(response).Encode(map[string]error{"error": err}); err != nil {
		internalServerError(response, err)
	}
}
