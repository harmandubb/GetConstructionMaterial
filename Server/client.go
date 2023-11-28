package server

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

type Text struct {
	text string
}

func clientTest() error {
	cert, err := os.ReadFile("cert.pem")
	if err != nil {
		log.Fatalf("Unable to read cert.pem: %v", err)
	}

	certPool := x509.NewCertPool()
	if !certPool.AppendCertsFromPEM(cert) {
		log.Fatalf("Failed to append cert to pool")
	}

	tlsConfig := &tls.Config{
		RootCAs: certPool,
	}

	client := http.Client{
		Transport: &http.Transport{TLSClientConfig: tlsConfig},
	}

	//use the client to make a request:
	resp, err := client.Get("https://localhost:443")
	if err != nil {
		log.Fatalf("Failed to make request: %v", err)
	}

	fmt.Println(resp.Body)

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Unable to read the https response: %v", err)
	}

	// var text Text

	// err = json.Unmarshal(body, &text)
	// if err != nil {
	// 	log.Fatalf("Unable to decode the https response: %v", err)
	// }

	fmt.Println(string(body))

	return nil
}
