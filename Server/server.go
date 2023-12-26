package server

import (
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
	"time"
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
		"https://www.docstruction.com":            true,
		"https://www.getconstructionmaterial.com": true,
		"https://docstruction.com":                true,
		"https://getconstructionmaterial.com":     true,
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
			fmt.Println("In the post")
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

			spreadsheetID := "1NXTK2G6sQOs0ZSQ1046ijoanPDNWPKOc0-I7dEMotQ8" //could make the storing of the id better. //Need to have the spread sheet id for the material form

			result := g.SendMaterialFormInfo(spreadsheetID, materialFormInfo)

			p := d.ConnectToDataBase(os.Getenv("DB_NAME")) //need to set this in a environmental variabl

			inquiryID, err := d.AddBlankCustomerInquiry(p, materialFormInfo, os.Getenv("CUSTOMER_INQUIRY_TABLE"))
			if err != nil {
				log.Fatalf("Error when adding customer inquiry to database: %v", err)
			}

			resp := ServerResponse{
				Success: result,
			}

			// catigorizationTemplate := os.Getenv("CATIGORIZATION_TEMPLATE")
			// emailTemplate := os.Getenv("EMAIL_TEMPLATE")

			go api.ProcessCustomerInquiry(p, inquiryID, catigorizationTemplate, emailTemplate)

			jsonResp, err := json.Marshal(resp)
			if err != nil {
				w.WriteHeader(http.StatusExpectationFailed)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			w.Write(jsonResp)
		}
	})
	log.Println("Server is starting on port 8080...")
	err := http.ListenAndServe("0.0.0.0:8080", nil)
	if err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
