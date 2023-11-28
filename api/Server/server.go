package server

import (
	"log"
	"net/http"
)

func idle() {

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, HTTPS!"))
	})

	err := http.ListenAndServeTLS(":443", "cert.pem", "key.pem", nil)
	if err != nil {
		log.Fatalf("Sever Error: %v", err)
	}

}
