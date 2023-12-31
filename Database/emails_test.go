package database

import (
	g "docstruction/getconstructionmaterial/GCalls"
	"os"

	"testing"

	"googlemaps.github.io/maps"
)

func TestAddBlankEmailInquiryEntry(t *testing.T) {
	p := ConnectToDataBase("mynewdatabase")

	err := ResetTestDataBase(p, "emails")
	if err != nil {
		t.Error(err)
	}

	supDetail := g.SupplierEmailInfo{
		MapsID:  "IDTEST",
		Name:    "SupplierNameTest",
		Address: "Surrey Place",
		Geometry: maps.AddressGeometry{
			Location: maps.LatLng{
				Lat: 10.0,
				Lng: 12.0,
			},
			LocationType: "Test",
		},
		Website: "wwww.test.com",
		Email:   []string{"test@gmail.com"},
	}

	inquiryID := "Test_ID"
	material := "Fire Stopping Collars"
	client_email := "Client_test@gmail.com"
	sent_out := false

	err = AddBlankEmailInquiryEntry(p, inquiryID, client_email, material, supDetail, sent_out, "emails")
	if err != nil {
		t.Error(err)
	}

	emailInquiry, err := ReadEmailInquiryEntry(p, "emails", IDOption{Inquiry_ID: inquiryID})
	if err != nil {
		t.Error(err)
	}

	if inquiryID != emailInquiry.Inquiry_ID {
		t.Fail()
	}

	if client_email != emailInquiry.Client_Email {
		t.Fail()
	}

	if material != emailInquiry.Material {
		t.Fail()
	}

	if supDetail.MapsID != emailInquiry.Supplier_Map_ID {
		t.Fail()
	}

	if supDetail.Name != emailInquiry.Supplier_Name {
		t.Fail()
	}

	if supDetail.Geometry.Location.Lat != emailInquiry.Supplier_Lat {
		t.Fail()
	}

	if supDetail.Geometry.Location.Lng != emailInquiry.Supplier_Lng {
		t.Fail()
	}

	if supDetail.Email[0] != emailInquiry.Supplier_Email {
		t.Fail()
	}

	if sent_out != emailInquiry.Sent_Out {
		t.Fail()
	}

	if false != emailInquiry.Present {
		t.Fail()
	}

	if 0 != emailInquiry.Price {
		t.Fail()
	}

	if "" != emailInquiry.Currency {
		t.Fail()
	}

	if nil != emailInquiry.Data_Sheet {
		t.Fail()
	}
}

func TestPDFUploadAndDownload(t *testing.T) {
	p := ConnectToDataBase("mynewdatabase")

	inquiry_id := "TEST_INQUIRY_ID"
	thread_id := "TEST_THREAD_ID"

	supDetail := g.SupplierEmailInfo{
		MapsID:  "IDTEST",
		Name:    "SupplierNameTest",
		Address: "Surrey Place",
		Geometry: maps.AddressGeometry{
			Location: maps.LatLng{
				Lat: 10.0,
				Lng: 12.0,
			},
			LocationType: "Test",
		},
		Website:        "wwww.test.com",
		Email:          []string{"test@gmail.com"},
		Email_ThreadID: thread_id,
	}

	err := AddBlankEmailInquiryEntry(p,
		inquiry_id,
		"Test@gmail.com",
		"Fire Stopping Pipe Collars",
		supDetail,
		true,
		"emails",
	)
	if err != nil {
		t.Error(err)
	}

	// Read in a datasheet in a bte array
	file, err := os.ReadFile("./Attachment/1.pdf")
	if err != nil {
		t.Error(err)
	}

	err = UpdateEmailInquiryEntryDataSheet(p, inquiry_id, "emails", &file)
	if err != nil {
		t.Error(err)
	}

	// read the file to ensure that the pdf is being processed properly:
	emailInquiryInfo, err := ReadEmailInquiryEntry(p, "emails", IDOption{Inquiry_ID: inquiry_id})
	if err != nil {
		t.Error(err)
	}

	err = os.WriteFile("output.pdf", *emailInquiryInfo.Data_Sheet, 0644)
	if err != nil {
		t.Error(err)
	}

}
