package server

import (
	g "docstruction/getconstructionmaterial/GCalls"
	"testing"

	"googlemaps.github.io/maps"
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

	loc := maps.LatLng{
		Lat: 49.05812,
		Lng: -122.81026,
	}

	matFormInfo := g.MaterialFormInfo{
		Email:    "info@gmail.com",
		Material: "Fire Stop Collars",
	}

	err := ContactSupplierForMaterial(matFormInfo, catigorizationTemplate, emailTemplate, &loc)
	if err != nil {
		t.Error(err)
	}
}
