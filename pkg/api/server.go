package api

import (
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
