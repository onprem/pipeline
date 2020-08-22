package api

import (
	"encoding/json"
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

	t := task.NewTask(id, filePath)
	a.taskStore[t.ID] = t

	t.Run()

	log.Println("[success] file uploaded: ", handler.Filename)

	resp, err := json.Marshal(response{Status: "success", Data: map[string]string{"id": t.ID}})
	if err != nil {
		http.Error(w, "error encoding response", http.StatusInternalServerError)
		log.Println("[error] encoding response: ", err)
		return
	}

	_, err = w.Write(resp)
	if err != nil {
		log.Println("[error] responding: ", err)
	}
}

func (a *API) handleStatus(w http.ResponseWriter, r *http.Request) {
	taskID := r.FormValue("id")
	if taskID == "" {
		http.Error(w, "id is required", http.StatusBadRequest)
		return
	}

	t, ok := a.taskStore[taskID]
	if !ok {
		http.Error(w, "invalid task id", http.StatusBadRequest)
		return
	}

	resp, err := json.Marshal(response{Status: "success", Data: map[string]task.Status{"status": t.State}})
	if err != nil {
		http.Error(w, "error encoding response", http.StatusInternalServerError)
		log.Println("[error] encoding response: ", err)
		return
	}

	_, err = w.Write(resp)
	if err != nil {
		log.Println("[error] responding: ", err)
	}
}

func (a *API) handlePause(w http.ResponseWriter, r *http.Request) {
	taskID := r.FormValue("id")
	if taskID == "" {
		http.Error(w, "id is required", http.StatusBadRequest)
		return
	}

	t, ok := a.taskStore[taskID]
	if !ok {
		http.Error(w, "invalid task id", http.StatusBadRequest)
		return
	}

	t.Pause()
	resp, err := json.Marshal(response{Status: "success", Data: map[string]string{"message": "task paused"}})
	if err != nil {
		http.Error(w, "error encoding response", http.StatusInternalServerError)
		log.Println("[error] encoding response: ", err)
		return
	}

	_, err = w.Write(resp)
	if err != nil {
		log.Println("[error] responding: ", err)
	}
}

func (a *API) handleResume(w http.ResponseWriter, r *http.Request) {
	taskID := r.FormValue("id")
	if taskID == "" {
		http.Error(w, "id is required", http.StatusBadRequest)
		return
	}

	t, ok := a.taskStore[taskID]
	if !ok {
		http.Error(w, "invalid task id", http.StatusBadRequest)
		return
	}

	t.Resume()
	resp, err := json.Marshal(response{Status: "success", Data: map[string]string{"message": "task resumed"}})
	if err != nil {
		http.Error(w, "error encoding response", http.StatusInternalServerError)
		log.Println("[error] encoding response: ", err)
		return
	}

	_, err = w.Write(resp)
	if err != nil {
		log.Println("[error] responding: ", err)
	}
}

func (a *API) handleTerminate(w http.ResponseWriter, r *http.Request) {
	taskID := r.FormValue("id")
	if taskID == "" {
		http.Error(w, "id is required", http.StatusBadRequest)
		return
	}

	t, ok := a.taskStore[taskID]
	if !ok {
		http.Error(w, "invalid task id", http.StatusBadRequest)
		return
	}

	t.Terminate()
	resp, err := json.Marshal(response{Status: "success", Data: map[string]string{"message": "task terminated"}})
	if err != nil {
		http.Error(w, "error encoding response", http.StatusInternalServerError)
		log.Println("[error] encoding response: ", err)
		return
	}

	_, err = w.Write(resp)
	if err != nil {
		log.Println("[error] responding: ", err)
	}
}
