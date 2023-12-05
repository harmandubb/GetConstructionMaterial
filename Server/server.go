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
	// origin := r.Header.Get("Origin")

	// // List of allowed origins
	// allowedOrigins := map[string]bool{
	// 	"https://www.docstruction.com":            true,
	// 	"https://www.getconstructionmaterial.com": true,
	// }

	// // Check if the origin is in the list of allowed origins
	// if _, ok := allowedOrigins[origin]; ok {
	// 	w.Header().Set("Access-Control-Allow-Origin", origin)
	// }
	w.Header().Set("Access-Control-Allow-Origin", "https://www.getconstructionmaterial.com")

	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")

}

func Idle() {
	fmt.Println("Starting Server")

	// TODO: Implement the serveMUX

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Println("You are visiting the foot backend")
		setCORS(w, r)

		// Handle OPTIONS for preflight
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}
		fmt.Println("in the main")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "OK")

		responseText := "Hello, you are in the main"
		w.Write([]byte(responseText)) // Convert string to []byte
	})

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		setCORS(w, r)

		// Handle OPTIONS for preflight
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "OK")
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

	// err := http.ListenAndServeTLS("0.0.0.0:8080", getPath("cert.pem"), getPath("key.pem"), nil)
	err := http.ListenAndServe("0.0.0.0:8080", nil)
	if err != nil {
		log.Fatalf("Sever Error: %v", err)
	}

}
