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

func Idle() {
	fmt.Println("Starting Server")

	// TODO: Implement the serveMUX

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("I am in the default form branch")
		// Set CORS headers for all responses
		w.Header().Set("Access-Control-Allow-Origin", "*") // Adjust in production
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")
		w.Header().Set("Content-Type", "application/json")

		// Handle OPTIONS for preflight
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		if r.Method == http.MethodPost {

			w.Write([]byte("Hello, HTTPS!"))
		}
	})

	http.HandleFunc("/emailForm", func(w http.ResponseWriter, r *http.Request) {

		// Set CORS headers for all responses
		w.Header().Set("Access-Control-Allow-Origin", "*") // Adjust in production
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")

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
				log.Fatalf("Cannot read request body: %v", err)
			}

			var emailFormInfo EmailFormInfo

			err = json.Unmarshal(body, &emailFormInfo)
			if err != nil {
				log.Fatalf("Cannot convert data into struct: %v", err)
			}

			spreadsheetID := "1ZowyzJ008toPYNn0mFc2wG6YTAop9HfnbMPLIM4rRZw" //could make the storing of the id better.

			result := g.SendEmailInfo(emailFormInfo.Time, emailFormInfo.Email, spreadsheetID)

			resp := ServerResponse{
				Success: result,
			}

			fmt.Println(resp.Success)

			jsonResp, err := json.Marshal(resp)
			if err != nil {
				log.Fatal("Was not able to encode struct to json")
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			fmt.Println(jsonResp)

			fmt.Println("Sending response")

			w.Write(jsonResp)
		}

	})

	err := http.ListenAndServeTLS(":443", getPath("cert.pem"), getPath("key.pem"), nil)
	if err != nil {
		log.Fatalf("Sever Error: %v", err)
	}

}
