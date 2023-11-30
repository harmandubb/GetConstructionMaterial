package server

import (
	g "docstruction/getconstructionmaterial/GCalls"
	"encoding/json"
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

	// TODO: Implement the serveMUX

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, HTTPS!"))
	})

	http.HandleFunc("/emailForm", func(w http.ResponseWriter, r *http.Request) {
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

		jsonResp, err := json.Marshal(resp)
		if err != nil {
			log.Fatal("Was not able to encode struct to json")
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		w.Header().Set("Content-Type", "application/json") //need to google http headers for all of the headers that can be used

		w.Write(jsonResp)

	})

	err := http.ListenAndServeTLS(":443", getPath("cert.pem"), getPath("key.pem"), nil)
	if err != nil {
		log.Fatalf("Sever Error: %v", err)
	}

}
