package main

import (
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/prmsrswt/pipeline/pkg/api"
)

const (
	uploadDir = "uploads"
)

func main() {
	// Set up directory for uploads
	err := os.MkdirAll(uploadDir, 0755)
	if err != nil {
		log.Fatalln(err)
	}
	// Seed randon number generator
	rand.Seed(time.Now().UnixNano())

	mux := http.NewServeMux()

	pipelineAPI := api.NewAPI(uploadDir)
	pipelineAPI.Register(mux)

	log.Println("Web server started")
	log.Fatalln(http.ListenAndServe(":8080", mux))
}
