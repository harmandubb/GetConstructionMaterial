package gcalls

import (
	"testing"

	"googlemaps.github.io/maps"
)

func TestMakeSearchNearByRequest(t *testing.T) {
	loc := maps.LatLng{
		Lat: 49.05812,
		Lng: -122.81026,
	}

	err := connectToMaps("Electrical", &loc)
	if err != nil {
		t.Error(err)
	}
}
