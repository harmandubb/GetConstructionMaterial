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

	"googlemaps.github.io/maps"
)

type EmailFormInfo struct {
	Time  time.Time
	Email string
}

type ServerResponse struct {
	Success bool
}

const SUPPLIERCONTACTLIMIT = 3

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
		w.Header().Set("Access-Control-Allow-Origin", origin)
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

	http.HandleFunc("/materialForm", func(w http.ResponseWriter, r *http.Request) {
		setCORS(w, r)

		// Handle OPTIONS for preflight
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		if r.Method == http.MethodPost {
			w.Header().Set("Content-Type", "application/json")
			fmt.Println("I am in the material form branch")

			body, err := io.ReadAll(r.Body)
			if err != nil {
				w.WriteHeader(http.StatusExpectationFailed)
				return
			}

			var materialFormInfo g.MaterialFormInfo

			err = json.Unmarshal(body, &materialFormInfo)
			if err != nil {
				w.WriteHeader(http.StatusExpectationFailed)
				return
			}

			spreadsheetID := "1NXTK2G6sQOs0ZSQ1046ijoanPDNWPKOc0-I7dEMotQ8" //could make the storing of the id better. //Need to have the spread sheet id for the material form

			result := g.SendMaterialFormInfo(spreadsheetID, materialFormInfo)

			resp := ServerResponse{
				Success: result,
			}

			//Call chat gpt to catigorized the item
			catigorizationTemplate := "./material_catigorization_prompt.txt"

			catergory, err := PromptGPTMaterialCatogorization(catigorizationTemplate, materialFormInfo.Material)
			if err != nil {
				log.Fatalf("Catogirization Error: %v", err)
			}

			// Search for near by supplies for the category
			c, err := g.GetMapsClient()
			if err != nil {
				log.Fatalf("Map Client Connection Error: %v", err)
			}

			//TODO: get the location data from the user using the site

			loc := maps.LatLng{
				Lat: 49.05812,
				Lng: -122.81026,
			}

			searchResp, err := g.SearchSuppliers(c, catergory, &loc)
			if err != nil {
				log.Fatalf("Map Search Supplier Error: %v", err)
			}

			var supplierInfo []g.SupplierInfo

			for _, supplier := range searchResp.Results {
				supplier, _ := g.GetSupplierInfo(c, supplier)

				supplierInfo = append(supplierInfo, supplier)
			}

			//Get the supplier emails from the info that is found

			for _, supInfo := range supplierInfo {
				email, err := FindSupplierContactEmail(supInfo.Website)
				if err != nil {
					log.Fatalf("Supplier Email Get Error: %v", err)
					break
				}

				supInfo.Email = email
			}

			counter := 0

			srv := g.ConnectToGmailAPI()

			for _, supInfo := range supplierInfo {
				if counter < SUPPLIERCONTACTLIMIT {
					if len(supInfo.Email) != 0 {
						// get the email prompt from chat gpt
						subj, body, err := CreateEmailToSupplier("../email_prompt.txt", supInfo.Name, materialFormInfo.Material)
						if err != nil {
							log.Fatalf("GPT Email Create Error: %v", err)
						}

						// send the emal to the supplier
						g.SendEmail(srv, subj, body, supInfo.Email[0])
						counter = counter + 1
					}
				} else {
					break
				}
			}

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
	err := http.ListenAndServe("0.0.0.0:8080", nil)
	// err := http.ListenAndServeTLS("0.0.0.0:8080", getPath("cert.pem"), getPath("key.pem"), nil)
	// err := http.ListenAndServeTLS("0.0.0.0:8080", "cert.pem", "key.pem", nil)
	if err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
