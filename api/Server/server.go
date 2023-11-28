package server

import (
	"log"
	"net/http"
)

func idle() {

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, HTTPS!"))
	})

	err := http.ListenAndServeTLS(":443", "./Server/cert.pem", "./Server/key.pem", nil)
	if err != nil {
		log.Fatal("Sever Error")
	}

}
