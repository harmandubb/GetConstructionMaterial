package database

import (
	g "docstruction/getconstructionmaterial/GCalls"
	"testing"
)

func TestAddBlankCustomerInquiry(t *testing.T) {
	p := connectToDataBase("mynewdatabase")

	matForm := g.MaterialFormInfo{
		Email:    "harmand1999@gmail",
		Loc:      "Surrey BC",
		Material: "Fire Stop Fire Collars 2 in",
	}

	err := AddBlankCustomerInquiry(p, matForm, "mynewdatabase", "Customer_Inquiry")
	if err != nil {
		t.Error(err)
	}
}
