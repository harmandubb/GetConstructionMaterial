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

	"googlemaps.github.io/maps"
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

func TestAddressPushNotification(t *testing.T) {

	// Sample Email used:
	// Hello Docstruction,

	// We do have fire stop collars.

	// For 2 in the price is $3.48 per collar.

	// Let me know if you beed anything else.

	// Thanks,

	// Harman

	p := d.ConnectToDataBase("mynewdatabase")
	srv := g.ConnectToGmailAPI()
	c, err := g.GetMapsClient()
	if err != nil {
		t.Error(err)
	}

	user := "info@docstruction.com"

	// construct a dumby table entery to work on
	// Get the email thread ID that is needed from the email that you want to work with
	messages, err := g.GetUnreadMessagesData(srv, user)
	if err != nil {
		t.Error(err)
	}

	sup_thread_id := messages.Messages[0].ThreadId

	matFormInfo := g.MaterialFormInfo{
		Email:    "harmand1999@gmail.com",
		Material: "Fire Stop Collars",
		Loc:      "Richmond BC",
	}

	currency := g.GetCurrency(c, matFormInfo.Loc)

	inquiry_id, err := //write a customer_inquiry line as well
		d.AddBlankCustomerInquiry(
			p,
			matFormInfo,
			"customer_inquiry",
			currency,
		)

	if err != nil {
		t.Error(err)
	}

	//Create a emails entery with the email info above.
	err = d.AddBlankEmailInquiryEntry(
		p,
		inquiry_id,
		"test_client@gmail.com",
		"Fire Stopping Collars 2 in",
		g.SupplierEmailInfo{
			MapsID:  "TEST_MAPS_ID",
			Name:    "TEST_SUPPLIER",
			Address: "TEST BC",
			Geometry: maps.AddressGeometry{
				Location: maps.LatLng{
					Lat: 1.0,
					Lng: 2.0,
				},
				LocationType: "NON",
				Bounds: maps.LatLngBounds{
					NorthEast: maps.LatLng{
						Lat: 1.0,
						Lng: 1.0,
					},
				},
				Viewport: maps.LatLngBounds{
					NorthEast: maps.LatLng{
						Lat: 1.0,
						Lng: 1.0,
					},
				},
				Types: []string{"NON"},
			},
			Website:        "TEST.com",
			Email:          []string{"SUPPLIER_TEST@gmail.com"},
			Email_ThreadID: sup_thread_id,
		},
		true,
		"emails",
	)

	if err != nil {
		t.Error(err)
	}

	file, err := os.ReadFile("../Server/GPT_Prompts/email_receive_prompt.txt")
	if err != nil {
		t.Error(err)
	}

	prompt := string(file)

	err = AddressPushNotification(p, srv, user, prompt, "emails", "customer_inquiry")
	if err != nil {
		t.Error(err)
	}

	// Read the customer section and comapre if it is what is expected
	custInquiry, err := d.ReadCustomerInquiry(p, "customer_inquiry", inquiry_id)
	if err != nil {
		t.Error(err)
	}

	if custInquiry.Price != 5.95 {
		t.Fail()
	}

}
