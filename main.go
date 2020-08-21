package main

import (
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"path"
	"time"
)

const (
	uploadDir = "uploads"
)

var taskStore map[string]*Task

func handleUpload(w http.ResponseWriter, r *http.Request) {
	// Maximum upload of 50 MB files.
	r.ParseMultipartForm(50 << 20)

	file, handler, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "file is required", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Create a file locally
	filePath := path.Join(uploadDir, handler.Filename)
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

	task := NewTask(filePath)
	taskStore[task.ID] = task

	task.Run()

	log.Println("[success] file uploaded: ", handler.Filename)
	fmt.Fprintf(w, "Upload successful\ntaskID: %s\n", task.ID)
}

func handleStatus(w http.ResponseWriter, r *http.Request) {
	taskID := r.FormValue("id")
	if taskID == "" {
		http.Error(w, "id is required", http.StatusBadRequest)
		return
	}

	task, ok := taskStore[taskID]
	if !ok {
		http.Error(w, "invalid task id", http.StatusBadRequest)
		return
	}

	fmt.Fprintf(w, "status: %s\n", task.State)
}

func main() {
	// Set up directory for uploads
	err := os.MkdirAll(uploadDir, 0755)
	if err != nil {
		log.Fatalln(err)
	}
	// Seed randon number generator
	rand.Seed(time.Now().UnixNano())

	taskStore = make(map[string]*Task)

	http.HandleFunc("/upload", handleUpload)
	http.HandleFunc("/status", handleStatus)

	log.Println("Web server started")
	log.Fatalln(http.ListenAndServe(":8080", nil))
}
