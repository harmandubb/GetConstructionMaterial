package gcalls

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"googlemaps.github.io/maps"
)

func connectToMaps(category string, loc *maps.LatLng) error {
	err := godotenv.Load()
	if err != nil {
		return err
	}

	c, err := maps.NewClient(maps.WithAPIKey(os.Getenv("TESTING_MAPS_KEY")))
	if err != nil {
		log.Fatalf("fatal error: %s", err)
	}

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
		return err
	}

	// fmt.Println(nearByResp)
	// fmt.Println(nearByResp.Results)
	for _, val := range nearByResp.Results {
		fmt.Println(val)
	}

	return nil

}
