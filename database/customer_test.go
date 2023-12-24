package database

import (
	g "docstruction/getconstructionmaterial/GCalls"
	"testing"
)

func TestAddBlankCustomerInquiry(t *testing.T) {
	p := connectToDataBase("mynewdatabase")

	matForm := g.MaterialFormInfo{
		Email:    "harmand1999@gmail.com",
		Loc:      "Surrey BC",
		Material: "Fire Stop Fire Collars 2 in",
	}

	err := AddBlankCustomerInquiry(p, matForm, "mynewdatabase", "Customer_Inquiry")
	if err != nil {
		t.Error(err)
	}

	cust, err := readCustomerInquiry("customer_inquiry", "harmand1999@gmail.com")
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