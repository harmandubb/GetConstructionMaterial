package server

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"log"
	"net/http"
	"os"
)

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
		log.Fatal("Failed to make request: %v", err)
	}

	fmt.Println(resp.Body)

	defer resp.Body.Close()

	return nil
}
