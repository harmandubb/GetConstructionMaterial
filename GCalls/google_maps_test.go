package gcalls

import (
	"fmt"
	"testing"

	"googlemaps.github.io/maps"
)

func TestSearchSuppliers(t *testing.T) {
	loc := maps.LatLng{
		Lat: 49.05812,
		Lng: -122.81026,
	}

	c, err := GetMapsClient()
	if err != nil {
		t.Error(err)
	}

	_, err = SearchSuppliers(c, "Electrical", &loc)
	if err != nil {
		t.Error(err)
	}
}

func TestGetSupplierInfo(t *testing.T) {
	loc := maps.LatLng{
		Lat: 49.05812,
		Lng: -122.81026,
	}

	c, err := GetMapsClient()
	if err != nil {
		t.Error(err)
	}

	supplierResponse, err := SearchSuppliers(c, "Electrical", &loc)
	if err != nil {
		t.Error(err)
	}

	supplierInfo, err := GetSupplierInfo(c, supplierResponse.Results[0])
	if err != nil {
		t.Error(err)
	}

	fmt.Println(supplierInfo.Name)
	fmt.Println(supplierInfo.Address)
	fmt.Println(supplierInfo.Geometry)
	fmt.Println(supplierInfo.Website)

}

func TestGeocdeGeneralLocation(t *testing.T) {
	c, err := GetMapsClient()
	if err != nil {
		t.Error(err)
	}

	result, err := GeocodeGeneralLocation(c, "Surrey BC")
	if err != nil {
		t.Error(err)
	}

	fmt.Println(result)
}
