package server

import (
	"context"
	api "docstruction/getconstructionmaterial/API"
	d "docstruction/getconstructionmaterial/Database"
	g "docstruction/getconstructionmaterial/GCalls"
	_ "embed"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/api/gmail/v1"
)

type EmailFormInfo struct {
	Time  time.Time
	Email string
}

type ServerResponse struct {
	Success bool
}

//go:embed GPT_Prompts/material_catigorization_prompt.txt
var catigorizationTemplate string

//go:embed GPT_Prompts/email_prompt.txt
var emailTemplate string

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
		"https://www.docstruction.com":                                            true,
		"https://www.getconstructionmaterial.com":                                 true,
		"https://docstruction.com":                                                true,
		"https://getconstructionmaterial.com":                                     true,
		"getconstructionmaterial@getconstructionmaterial.iam.gserviceaccount.com": true,
	}

	fmt.Println("Origin Request:", origin)
	fmt.Println("Server present:", allowedOrigins[origin])

	_, ok := allowedOrigins[origin]

	fmt.Print("OK:", ok)

	if _, ok := allowedOrigins[origin]; ok {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")
	}

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

			body, err := io.ReadAll(r.Body)
			if err != nil {
				w.WriteHeader(http.StatusExpectationFailed)
				fmt.Println(err)
				return
			}

			var emailFormInfo EmailFormInfo

			err = json.Unmarshal(body, &emailFormInfo)
			if err != nil {
				w.WriteHeader(http.StatusExpectationFailed)
				fmt.Println(err)
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
				fmt.Println(err)
				return
			}

			w.Write(jsonResp)
		}

	})

	http.HandleFunc("/materialForm", func(w http.ResponseWriter, r *http.Request) {
		setCORS(w, r)

		// err := godotenv.Load()
		// if err != nil {
		// 	log.Fatalf("Error loading .env file: %v", err)
		// }

		// Handle OPTIONS for preflight
		if r.Method == http.MethodOptions {
			fmt.Println("In the options")
			w.WriteHeader(http.StatusOK)
			return
		}

		if r.Method == http.MethodPost {
			w.Header().Set("Content-Type", "application/json")

			body, err := io.ReadAll(r.Body)
			if err != nil {
				fmt.Println("Error in read function")
				w.WriteHeader(http.StatusExpectationFailed)
				return
			}

			var materialFormInfo g.MaterialFormInfo

			err = json.Unmarshal(body, &materialFormInfo)
			if err != nil {
				fmt.Println("Error in json function")
				w.WriteHeader(http.StatusExpectationFailed)
				return
			}

			dataBaseConnectionPool := &sync.Pool{
				New: func() interface{} {
					return d.ConnectToDataBase(os.Getenv("DB_NAME"))
				},
			}

			gmailServicePool := sync.Pool{
				New: func() interface{} {
					return g.ConnectToGmailAPI()
				},
			}

			p := dataBaseConnectionPool.Get().(*pgxpool.Pool)
			defer dataBaseConnectionPool.Put(p)

			errStream := make(chan error)
			inquiryIDStream := make(chan string)

			ctx, cancel := context.WithCancel(context.Background())

			go d.ConcurrentAddBlankCustomerInquiry(inquiryIDStream, errStream, ctx, p, materialFormInfo, os.Getenv("CUSTOMER_INQUIRY_TABLE"))

			//save result in the spread sheet for a back up
			spreadsheetID := "1NXTK2G6sQOs0ZSQ1046ijoanPDNWPKOc0-I7dEMotQ8"

			result := g.SendMaterialFormInfo(spreadsheetID, materialFormInfo)

			resp := ServerResponse{
				Success: result,
			}

			var inquiryID string

			// Want to assess what the outcome of the of the blacnk add of the customer inquiry looks like
			select {
			//error stream is full
			case err := <-errStream:
				fmt.Println("Error Occured when making a entry of the inquiry: %v", err)
				w.WriteHeader(http.StatusExpectationFailed)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			//uniqueID is present
			case inquiryID = <-inquiryIDStream:

			case <-time.After(5 * time.Second):
				cancel()
				fmt.Println("Adding blank inquiry has timed out")
				w.WriteHeader(http.StatusExpectationFailed)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			close(errStream)
			close(inquiryIDStream)

			errStream = make(chan error)

			// TODO: How should the cancelling of the goroutine occur and under what circumstances.

			// catigorizationTemplate := os.Getenv("CATIGORIZATION_TEMPLATE")
			// emailTemplate := os.Getenv("EMAIL_TEMPLATE")
			var wg sync.WaitGroup
			wg.Add(1)

			srv := gmailServicePool.Get().(*gmail.Service)
			defer gmailServicePool.Put(srv)

			go api.ConcurrentProcessCustomerInquiry(&wg, errStream, srv, p, inquiryID, catigorizationTemplate, emailTemplate)

			jsonResp, err := json.Marshal(resp)
			if err != nil {
				w.WriteHeader(http.StatusExpectationFailed)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			w.Write(jsonResp)

			wg.Wait() //Wait for the Customer Inquiry process to finish for this thread to move onto another problem.

		}
	})

	http.HandleFunc("/emailNotification", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("New Email Received")

		setCORS(w, r)

		// Handle OPTIONS for preflight
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

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
