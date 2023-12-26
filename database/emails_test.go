package database

import (
	g "docstruction/getconstructionmaterial/GCalls"

	"testing"

	"googlemaps.github.io/maps"
)

func TestAddBlankEmailInquiryEntry(t *testing.T) {
	p := ConnectToDataBase("mynewdatabase")

	err := ResetTestDataBase(p, "emails")
	if err != nil {
		t.Error(err)
	}

	supDetail := g.SupplierInfo{
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

	emailInquiry, err := ReadEmailInquiryEntry(p, "emails", inquiryID)
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
