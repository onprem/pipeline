package api

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"

	"github.com/prmsrswt/pipeline/pkg/task"

	"github.com/google/uuid"
)

func (a *API) handleUpload(w http.ResponseWriter, r *http.Request) {
	// Maximum upload of 50 MB files.
	r.ParseMultipartForm(50 << 20)

	file, handler, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "file is required", http.StatusBadRequest)
		return
	}
	defer file.Close()

	id := uuid.New().String()

	// Create a file locally
	filePath := path.Join(a.uploadDir, id+handler.Filename)
	dst, err := os.Create(filePath)
	if err != nil {
		http.Error(w, "error creating file", http.StatusInternalServerError)
		log.Println("[error] creating file: ", err)
		return
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		http.Error(w, "error saving file", http.StatusInternalServerError)
		log.Println("[error] saving file: ", err)
		return
	}

	task := task.NewTask(id, filePath)
	a.taskStore[task.ID] = task

	task.Run()

	log.Println("[success] file uploaded: ", handler.Filename)
	fmt.Fprintf(w, "Upload successful\ntaskID: %s\n", task.ID)
}

func (a *API) handleStatus(w http.ResponseWriter, r *http.Request) {
	taskID := r.FormValue("id")
	if taskID == "" {
		http.Error(w, "id is required", http.StatusBadRequest)
		return
	}

	task, ok := a.taskStore[taskID]
	if !ok {
		http.Error(w, "invalid task id", http.StatusBadRequest)
		return
	}

	fmt.Fprintf(w, "status: %s\n", task.State)
}

func (a *API) handlePause(w http.ResponseWriter, r *http.Request) {
	taskID := r.FormValue("id")
	if taskID == "" {
		http.Error(w, "id is required", http.StatusBadRequest)
		return
	}

	task, ok := a.taskStore[taskID]
	if !ok {
		http.Error(w, "invalid task id", http.StatusBadRequest)
		return
	}

	task.Pause()
	fmt.Fprintf(w, "task %s paused\n", task.ID)
}

func (a *API) handleResume(w http.ResponseWriter, r *http.Request) {
	taskID := r.FormValue("id")
	if taskID == "" {
		http.Error(w, "id is required", http.StatusBadRequest)
		return
	}

	task, ok := a.taskStore[taskID]
	if !ok {
		http.Error(w, "invalid task id", http.StatusBadRequest)
		return
	}

	task.Resume()
	fmt.Fprintf(w, "task %s resumed\n", task.ID)
}

func (a *API) handleTerminate(w http.ResponseWriter, r *http.Request) {
	taskID := r.FormValue("id")
	if taskID == "" {
		http.Error(w, "id is required", http.StatusBadRequest)
		return
	}

	task, ok := a.taskStore[taskID]
	if !ok {
		http.Error(w, "invalid task id", http.StatusBadRequest)
		return
	}

	task.Terminate()
	fmt.Fprintf(w, "task %s terminated\n", task.ID)
}
