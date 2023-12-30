package server

import (
	"bytes"
	g "docstruction/getconstructionmaterial/GCalls"
	_ "embed"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
	"testing"
)

func TestIdle(t *testing.T) {
	Idle()
}

func TestClientTest(t *testing.T) {
	clientTest()
}

func TestMaterialEndPointConcurrently(t *testing.T) {
	matFormInfo := g.MaterialFormInfo{
		Email:    "info@gmail.com",
		Material: "Fire Stop Collars",
		Loc:      "Richmond BC",
	}

	jsonData, err := json.Marshal(matFormInfo)
	if err != nil {
		log.Fatalf("Error marshaling JSON: %v", err)
	}

	// URL of the endpoint you want to test
	url := "http://localhost:8080/materialForm"

	// WaitGroup to wait for all goroutines to finish
	var wg sync.WaitGroup

	// Number of concurrent requests
	const numRequests = 20

	for i := 0; i < numRequests; i++ {
		wg.Add(1)

		// Start a goroutine
		go func(i int) {
			defer wg.Done()

			// Send a GET request to the server
			req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
			if err != nil {
				log.Printf("Request %d failed: %v", i, err)
				return
			}

			// Set the Content-Type header
			req.Header.Set("Content-Type", "application/json")

			// Create an HTTP client and send the request
			client := &http.Client{}
			response, err := client.Do(req)
			if err != nil {
				log.Printf("Error sending request %d: %v", i, err)
				return
			}

			defer response.Body.Close()

			// Read the response body
			responseBody, err := io.ReadAll(response.Body)
			if err != nil {
				log.Printf("Error reading response for request %d: %v", i, err)
				return
			}

			// Print the response status and body
			fmt.Printf("Request %d - Status Code: %d, Response Body: %s\n", i, response.StatusCode, string(responseBody))
		}(i)
	}

	// Wait for all requests to complete
	wg.Wait()
}
