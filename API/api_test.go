package api

import (
	"bytes"
	d "docstruction/getconstructionmaterial/Database"
	g "docstruction/getconstructionmaterial/GCalls"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"testing"
)

func TestContactSupplierForMaterial(t *testing.T) {

	matFormInfo := g.MaterialFormInfo{
		Email:    "info@gmail.com",
		Material: "Fire Stop Collars",
		Loc:      "Richmond BC",
	}

	catigorizationFilePath := "../Server/GPT_Prompts/material_catigorization_prompt.txt"

	file, err := os.Open(catigorizationFilePath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// Read the entire file
	data, err := io.ReadAll(file)
	if err != nil {
		log.Fatal(err)
	}

	// Convert the bytes to a string
	catigorizationTemplate := string(data)

	emailFilePath := "../Server/GPT_Prompts/email_prompt.txt"

	file, err = os.Open(emailFilePath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// Read the entire file
	data, err = io.ReadAll(file)
	if err != nil {
		log.Fatal(err)
	}

	// Convert the bytes to a string
	emailTemplate := string(data)

	srv := g.ConnectToGmailAPI()

	_, err = ContactSupplierForMaterial(srv, matFormInfo, catigorizationTemplate, emailTemplate)
	if err != nil {
		t.Error(err)
	}
}

func TestAlertAdmin(t *testing.T) {
	matFormInfo := g.MaterialFormInfo{
		Email:    "info@gmail.com",
		Material: "Fire Stop Collars",
		Loc:      "Richmond BC",
	}

	srv := g.ConnectToGmailAPI()

	err := AlertAdmin(srv, matFormInfo, []string{"test1@example.com", "test2@example.com", "test3@example.com", "test4@example.com"})
	if err != nil {
		t.Fail()
	}
}

func TestMaterialFormHandler(t *testing.T) {
	//Want to clean the database that is local
	p := d.ConnectToDataBase("mynewdatabase")

	err := d.ResetTestDataBase(p, "customer_inquiry")
	if err != nil {
		t.Error(err)
	}
	url := "http://localhost:8080/materialForm"
	contentType := "application/json"
	matInfo := g.MaterialFormInfo{
		Email:    "test@gmail.com",
		Material: "Fire Stopping Pipe Collars 2 in",
		Loc:      "Las Angeles California",
	}

	content, err := json.Marshal(matInfo)
	if err != nil {
		t.Error(err)
	}

	reader := bytes.NewReader(content)

	resp, err := http.Post(url, contentType, reader)

	if err != nil {
		t.Error(err)
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}
	str := string(body)
	fmt.Println(str)
}
