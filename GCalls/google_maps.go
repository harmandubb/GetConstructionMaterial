package gcalls

import (
	"context"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"googlemaps.github.io/maps"
)

type SupplierInfo struct {
	ID       string
	Name     string
	Address  string
	Location maps.LatLng
	Website  string
}

func GetMapsClient() (*maps.Client, error) {
	err := godotenv.Load()
	if err != nil {
		return &maps.Client{}, err
	}

	c, err := maps.NewClient(maps.WithAPIKey(os.Getenv("TESTING_MAPS_KEY")))
	if err != nil {
		return &maps.Client{}, err
	}

	return c, nil

}

func SearchSuppliers(c *maps.Client, category string, loc *maps.LatLng) (maps.PlacesSearchResponse, error) {
	empty := maps.PlacesSearchResponse{}

	ctx := context.Background()

	nearByString := fmt.Sprintf("%s supplier", category)

	nearBySearchReq := maps.NearbySearchRequest{
		Location: loc,
		Radius:   15000,
		Keyword:  nearByString,
		Language: "en",
		OpenNow:  false,
	}

	nearByResp, err := c.NearbySearch(ctx, &nearBySearchReq)
	if err != nil {
		return empty, err
	}

	// fmt.Println(nearByResp)
	// fmt.Println(nearByResp.Results)
	// for _, val := range nearByResp.Results {
	// 	fmt.Println(val)
	// }

	return nearByResp, nil

}

func GetSupplierInfo(c *maps.Client, placeResult maps.PlacesSearchResult) (SupplierInfo, error) {
	ctx := context.Background()

	id := placeResult.PlaceID

	detailsReq := maps.PlaceDetailsRequest{
		PlaceID:  id,
		Language: "en",
		Fields: []maps.PlaceDetailsFieldMask{
			maps.PlaceDetailsFieldMaskGeometryLocation,
			maps.PlaceDetailsFieldMaskWebsite,
		},
	}

	placeDetailsResp, err := c.PlaceDetails(ctx, &detailsReq)
	if err != nil {
		return SupplierInfo{}, err
	}

	var address string

	if placeResult.FormattedAddress != "" {
		address = placeResult.FormattedAddress
	} else {
		address = placeResult.Vicinity
	}

	supInfo := SupplierInfo{
		ID:      id,
		Name:    placeResult.Name,
		Address: address,
		Location: maps.LatLng{
			Lat: placeDetailsResp.Geometry.Location.Lat,
			Lng: placeDetailsResp.Geometry.Location.Lng,
		},
		Website: placeDetailsResp.Website,
	}

	return supInfo, nil
}
