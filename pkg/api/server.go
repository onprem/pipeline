package api

import (
	"encoding/json"
	"net/http"

	"github.com/prmsrswt/pipeline/pkg/task"
)

// API represents the http API.
type API struct {
	taskStore map[string]*task.Task
	uploadDir string
}

// NewAPI returns an initialized instance of API.
func NewAPI(uploadDir string) *API {
	return &API{
		taskStore: make(map[string]*task.Task),
		uploadDir: uploadDir,
	}
}

// Register function registers the routes and handlers.
func (a *API) Register(mux *http.ServeMux) {
	mux.HandleFunc("/upload", a.handleUpload)
	mux.HandleFunc("/status", a.handleStatus)
	mux.HandleFunc("/pause", a.handlePause)
	mux.HandleFunc("/resume", a.handleResume)
	mux.HandleFunc("/terminate", a.handleTerminate)
}

type response struct {
	Status string      `json:"status"`
	Data   interface{} `json:"data"`
}

func respond(w http.ResponseWriter, resp response, code int) error {
	payload, err := json.Marshal(resp)
	if err != nil {
		return err
	}

	w.WriteHeader(code)
	_, err = w.Write(payload)
	return err
}

func respondError(w http.ResponseWriter, message string, code int) {
	respond(w, response{Status: "error", Data: map[string]string{"message": message}}, code)
}

func respondSuccess(w http.ResponseWriter, data interface{}) {
	respond(w, response{Status: "success", Data: data}, http.StatusOK)
}

func (a *API) getTaskFromReq(r *http.Request) (*task.Task, bool) {
	taskID := r.FormValue("id")
	if taskID == "" {
		return nil, false
	}

	t, ok := a.taskStore[taskID]
	return t, ok
}
