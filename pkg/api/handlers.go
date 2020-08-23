package api

import (
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
		respondError(w, "file is required", http.StatusBadRequest)
		return
	}
	defer file.Close()

	id := uuid.New().String()

	// Create a file locally
	filePath := path.Join(a.uploadDir, id+handler.Filename)
	dst, err := os.Create(filePath)
	if err != nil {
		respondError(w, "error creating file", http.StatusInternalServerError)
		log.Println("[error] creating file: ", err)
		return
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		respondError(w, "error saving file", http.StatusInternalServerError)
		log.Println("[error] saving file: ", err)
		return
	}

	t := task.NewTask(id, filePath)
	a.taskStore[t.ID] = t

	t.Run()
	respondSuccess(w, map[string]string{"id": t.ID})

	log.Println("[success] file uploaded: ", handler.Filename)
}

func (a *API) handleStatus(w http.ResponseWriter, r *http.Request) {
	t, ok := a.getTaskFromReq(r)
	if !ok {
		respondError(w, "invalid task id", http.StatusBadRequest)
		return
	}

	respondSuccess(w, map[string]task.Status{"status": t.State})
}

func (a *API) handlePause(w http.ResponseWriter, r *http.Request) {
	t, ok := a.getTaskFromReq(r)
	if !ok {
		respondError(w, "invalid task id", http.StatusBadRequest)
		return
	}

	t.Pause()
	respondSuccess(w, map[string]string{"message": "task paused"})
}

func (a *API) handleResume(w http.ResponseWriter, r *http.Request) {
	t, ok := a.getTaskFromReq(r)
	if !ok {
		respondError(w, "invalid task id", http.StatusBadRequest)
		return
	}

	t.Resume()
	respondSuccess(w, map[string]string{"message": "task resumed"})
}

func (a *API) handleTerminate(w http.ResponseWriter, r *http.Request) {
	t, ok := a.getTaskFromReq(r)
	if !ok {
		respondError(w, "invalid task id", http.StatusBadRequest)
		return
	}

	t.Terminate()
	respondSuccess(w, map[string]string{"message": "task terminated"})
}
