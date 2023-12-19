package server

import (
	g "docstruction/getconstructionmaterial/GCalls"
	_ "embed"
	"testing"
)

func TestIdle(t *testing.T) {
	Idle()
}

func TestClientTest(t *testing.T) {
	clientTest()
}

// //go:embed GPT_Prompts/material_catigorization_prompt.txt
// var catigorizationTemplate string

// //go:embed GPT_Prompts/email_prompt.txt
// var emailTemplate string

func TestContactSupplierForMaterial(t *testing.T) {

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
