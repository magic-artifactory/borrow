package instances

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/magic-artifactory/borrow/api/v1alpha1"
	"github.com/magic-artifactory/borrow/pkg/operatools"
)

type Endpoints struct {
	Conciliator *operatools.Conciliator
}

func NewEndpoints(conciliator *operatools.Conciliator) *Endpoints {
	return &Endpoints{
		Conciliator: conciliator,
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

	pendingInstance, err := addInstance(e.Conciliator, body)
	if err != nil {
		internalServerError(response, err)
		return
	}

	response.Header().Set("Content-type", "application/json")
	response.WriteHeader(http.StatusOK)

	if pendingInstance != nil {
		state = true
		status = "add to pending..."

		if err := json.NewEncoder(response).Encode(map[string]string{"creation": strconv.FormatBool(state), "status": status, "uuid": pendingInstance.Uid}); err != nil {
			internalServerError(response, err)
		}
	} else {
		status = "another instance already exist with same project name"
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

	if err := e.Conciliator.Client.List(e.Conciliator.Ctx, &msList); err != nil {
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
	response.WriteHeader(http.StatusInternalServerError)

	if err := json.NewEncoder(response).Encode(map[string]string{"error": err.Error()}); err != nil {
		log.Printf("failed to encode error response: %v", err)
	}
}
