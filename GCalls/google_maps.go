package gcalls

import (
	"context"
	"fmt"
	"os"
	"strings"

	"googlemaps.github.io/maps"
)

type SupplierEmailInfo struct {
	MapsID         string
	Name           string
	Address        string
	Geometry       maps.AddressGeometry
	Website        string
	Email          []string
	Email_ThreadID string
}

var countriesToCurrencies = map[string]string{
	"USA":    "USD",
	"UK":     "GBP",
	"CANADA": "CAD",
	"JPY":    "Yen",
	// other entries...
}

func GetMapsClient() (*maps.Client, error) {
	// err := godotenv.Load()
	// if err != nil {
	// 	return &maps.Client{}, err
	// }

	c, err := maps.NewClient(maps.WithAPIKey(os.Getenv("TESTING_MAPS_KEY")))
	if err != nil {
		return &maps.Client{}, err
	}

	return c, nil

}

// Purpose: get the near by suppliers based on the category that is identified for the material
// Parameters:
// c *maps.Client --> client that is etablish the api service for
// category string --> Supplier category to fine
// loc *maps.LatLng --> lat and longintue infromation to define the place you are finding suppliers near
// Return:
// resp maps.PlaceSearchResponse --> returns the respons from the maps api, would need to get details
// error if present.

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

func GetSupplierInfo(c *maps.Client, placeResult maps.PlacesSearchResult) (SupplierEmailInfo, error) {
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
		return SupplierEmailInfo{}, err
	}

	var address string

	if placeResult.FormattedAddress != "" {
		address = placeResult.FormattedAddress
	} else {
		address = placeResult.Vicinity
	}

	supInfo := SupplierEmailInfo{
		MapsID:   id,
		Name:     placeResult.Name,
		Address:  address,
		Geometry: placeDetailsResp.Geometry,
		Website:  placeDetailsResp.Website,
	}

	return supInfo, nil
}

func GeocodeGeneralLocation(c *maps.Client, loc string) (geometry maps.AddressGeometry, address string, err error) {
	fmt.Printf("Location String inputed: %s\n", loc)
	ctx := context.Background()

	geoReq := maps.GeocodingRequest{
		Address: loc,
	}

	geoResp, err := c.Geocode(ctx, &geoReq)
	if err != nil {
		return maps.AddressGeometry{}, "", err
	}

	return geoResp[0].Geometry, geoResp[0].FormattedAddress, err

}

func GetCurrency(c *maps.Client, loc string) (currency string) {
	_, address, err := GeocodeGeneralLocation(c, loc)
	if err != nil {
		currency = ""
	} else {
		currency = returnCurrencyCode(address)
	}

	return currency
}

func returnCurrencyCode(address string) (currency string) {
	words := strings.Fields(strings.ToUpper(address))
	country := words[len(words)-1]
	return countriesToCurrencies[country]
}
