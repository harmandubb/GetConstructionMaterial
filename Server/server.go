package server

import (
	g "docstruction/getconstructionmaterial/GCalls"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"path/filepath"
	"runtime"
	"time"
)

type EmailFormInfo struct {
	Time  time.Time
	Email string
}

type ServerResponse struct {
	Success bool
}

func getPath(relativePath string) string {
	_, b, _, _ := runtime.Caller(0)
	// The directory of the file
	basepath := filepath.Dir(b)
	// Construct the path relative to the file
	return filepath.Join(basepath, relativePath)
}

func setCORS(w http.ResponseWriter, r *http.Request) {
	origin := r.Header.Get("Origin")

	// List of allowed origins
	allowedOrigins := map[string]bool{
		"https://www.docstruction.com":            true,
		"https://www.getconstructionmaterial.com": true,
	}

	// Check if the origin is in the list of allowed origins
	if _, ok := allowedOrigins[origin]; ok {
		w.Header().Set("Access-Control-Allow-Origin", origin)
	}

	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")

}

func Idle() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// setCORS(w, r)
		log.Println("Health check request received")

		// Write an HTTP 200 OK status
		w.WriteHeader(http.StatusOK)

		// Send a response body
		w.Write([]byte("OK"))
	})

	http.HandleFunc("/emailForm", func(w http.ResponseWriter, r *http.Request) {
		setCORS(w, r)

		// Handle OPTIONS for preflight
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		if r.Method == http.MethodPost {
			w.Header().Set("Content-Type", "application/json")
			fmt.Println("I am in the email form branch")

			body, err := io.ReadAll(r.Body)
			if err != nil {
				w.WriteHeader(http.StatusExpectationFailed)
				return
			}

			var emailFormInfo EmailFormInfo

			err = json.Unmarshal(body, &emailFormInfo)
			if err != nil {
				w.WriteHeader(http.StatusExpectationFailed)
				return
			}

			spreadsheetID := "1ZowyzJ008toPYNn0mFc2wG6YTAop9HfnbMPLIM4rRZw" //could make the storing of the id better.

			result := g.SendEmailInfo(emailFormInfo.Time, emailFormInfo.Email, spreadsheetID)

			resp := ServerResponse{
				Success: result,
			}

			fmt.Println(resp.Success)

			jsonResp, err := json.Marshal(resp)
			if err != nil {
				w.WriteHeader(http.StatusExpectationFailed)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			fmt.Println(jsonResp)

			fmt.Println("Sending response")

			w.Write(jsonResp)
		}

	})

	log.Println("Server is starting on port 8080...")
	// err := http.ListenAndServe("0.0.0.0:8080", nil)
	err := http.ListenAndServeTLS("0.0.0.0:8080", getPath("cert.pem"), getPath("key.pem"), nil)
	if err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
