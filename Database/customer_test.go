package database

import (
	"bytes"
	g "docstruction/getconstructionmaterial/GCalls"
	"testing"
)

func TestAddBlankCustomerInquiry(t *testing.T) {
	p := ConnectToDataBase("mynewdatabase")
	c, _ := g.GetMapsClient()

	matForm := g.MaterialFormInfo{
		Email:    "harmand1999@gmail.com",
		Loc:      "Surrey BC",
		Material: "Fire Stop Fire Collars 2 in",
	}

	currency := g.GetCurrency(c, matForm.Loc)

	_, err := AddBlankCustomerInquiry(p, matForm, currency, "customer_inquiry")
	if err != nil {
		t.Error(err)
	}

	cust, err := ReadCustomerInquiry(p, "customer_inquiry", "harmand1999@gmail.com")
	if err != nil {
		t.Error(err)
	}

	correctInquiry := CustomerInquiry{
		Email:    matForm.Email,
		Material: matForm.Material,
		Loc:      matForm.Loc,
		Present:  false,
	}

	if correctInquiry.Email != cust.Email {
		t.Fail()
	}

	if correctInquiry.Material != cust.Material {
		t.Fail()
	}

	if correctInquiry.Loc != cust.Loc {
		t.Fail()
	}

	if correctInquiry.Present != cust.Present {
		t.Fail()
	}

	if correctInquiry.Price != cust.Price {
		t.Fail()
	}

	if correctInquiry.Currency != cust.Currency {
		t.Fail()
	}

	if correctInquiry.Data_Sheet != cust.Data_Sheet {
		t.Fail()
	}

}

func TestGenerateInquiryNumber(t *testing.T) {
	generateInquiryID()

}

func TestUpadteCustomerInquiryMaterial(t *testing.T) {
	p := ConnectToDataBase("mynewdatabase")
	c, _ := g.GetMapsClient()

	matForm := g.MaterialFormInfo{
		Email:    "harmand1999@gmail.com",
		Loc:      "Surrey BC",
		Material: "Fire Stop Fire Collars 2 in",
	}

	tableName := "Customer_Inquiry"

	currency := g.GetCurrency(c, matForm.Loc)

	inquiry_id, err := AddBlankCustomerInquiry(p, matForm, currency, tableName)
	if err != nil {
		t.Error(err)
	}

	supplier_email_id_thread := "SupID"
	price := 10.0
	datasheet := []byte("THISISAPLACEHOLDER")

	err = UpdateCustomerInquiryMaterial(p, tableName, inquiry_id, supplier_email_id_thread, price, currency, &datasheet)
	if err != nil {
		t.Error(err)
	}

	cust, err := ReadCustomerInquiry(p, tableName, inquiry_id)
	if err != nil {
		t.Error(err)
	}

	if true != cust.Present {
		t.Fail()
	}

	if price != cust.Price {
		t.Fail()
	}

	if currency != cust.Currency {
		t.Fail()
	}

	if !bytes.Equal(datasheet, *cust.Data_Sheet) {
		t.Fail()
	}

}
