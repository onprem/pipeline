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
func (a *API) Register() {
	http.HandleFunc("/upload", a.handleUpload)
	http.HandleFunc("/status", a.handleStatus)
	http.HandleFunc("/pause", a.handlePause)
	http.HandleFunc("/resume", a.handleResume)
	http.HandleFunc("/terminate", a.handleTerminate)
}
