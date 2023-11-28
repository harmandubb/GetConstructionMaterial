package server

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"
)

type EmailFormInfo struct {
	Time  time.Time
	Email string
}

func idle() {

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, HTTPS!"))
	})

	http.HandleFunc("/emailForm", func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			log.Fatalf("Cannot read request body: %v")
		}

		var emailFormInfo EmailFormInfo

		err = json.Unmarshal(body, &emailFormInfo)
		if err != nil {
			log.Fatalf("Cannot convert data into struct: %v")
		}

	})

	err := http.ListenAndServeTLS(":443", "cert.pem", "key.pem", nil)
	if err != nil {
		log.Fatalf("Sever Error: %v", err)
	}

}
