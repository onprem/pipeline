package api

import (
	"net/http"

	"github.com/prmsrswt/pipeline/pkg/task"
)

type API struct {
	taskStore map[string]*task.Task
	uploadDir string
}

func NewAPI(uploadDir string) *API {
	return &API{
		taskStore: make(map[string]*task.Task),
		uploadDir: uploadDir,
	}
}

func (a *API) Register() {
	http.HandleFunc("/upload", a.handleUpload)
	http.HandleFunc("/status", a.handleStatus)
	http.HandleFunc("/pause", a.handlePause)
	http.HandleFunc("/resume", a.handleResume)
	http.HandleFunc("/terminate", a.handleTerminate)
}
