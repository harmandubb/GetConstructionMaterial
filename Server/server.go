package server

import (
	"log"
	"net/http"
)

func Idle() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Health check request received")

		// Write an HTTP 200 OK status
		w.WriteHeader(http.StatusOK)

		// Send a response body
		w.Write([]byte("OK"))
	})

	log.Println("Server is starting on port 8080...")
	err := http.ListenAndServe("0.0.0.0:8080", nil)
	if err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
