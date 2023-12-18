package server

import (
	g "docstruction/getconstructionmaterial/GCalls"
	"testing"
)

func TestIdle(t *testing.T) {
	Idle()
}

func TestClientTest(t *testing.T) {
	clientTest()
}

func TestContactSupplierForMaterial(t *testing.T) {
	catigorizationTemplate := "../material_catigorization_prompt.txt"
	emailTemplate := "../email_prompt.txt"

	matFormInfo := g.MaterialFormInfo{
		Email:    "info@gmail.com",
		Material: "Fire Stop Collars",
		Location: "Richmond BC",
	}

	err := ContactSupplierForMaterial(matFormInfo, catigorizationTemplate, emailTemplate)
	if err != nil {
		t.Error(err)
	}
}
