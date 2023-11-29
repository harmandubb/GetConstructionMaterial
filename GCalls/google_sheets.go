package gcalls

import (
	"context"
	"log"
	"os"
	"time"

	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

func sendEmailInfo(time time.Time, email string, spreadSheetID string) bool {
	srv := connectToSheetsAPI()
	return appendEmailToSpreadSheet(srv, spreadSheetID, time, email)

}

// spreadsheet id: 1ZowyzJ008toPYNn0mFc2wG6YTAop9HfnbMPLIM4rRZw

func connectToSheetsAPI() *sheets.Service {
	ctx := context.Background()

	b, err := os.ReadFile("../Auth2/credentials.json")

	if err != nil {
		log.Fatalf("Unable to read crednetials: %v", err)
	}

	srv, err := sheets.NewService(ctx, option.WithCredentialsJSON(b))
	if err != nil {
		log.Fatalf("Unable to connnect to service %v", err)
	}

	return srv

}

func appendEmailToSpreadSheet(srv *sheets.Service, id string, time time.Time, email string) bool {
	success := false

	values := sheets.ValueRange{
		Values: [][]interface{}{
			{time, email},
		},
	}

	resp, err := srv.Spreadsheets.Values.Append(id, "Sheet1!A1", &values).Do()
	if err != nil {
		log.Fatalf("Appending request Failed: %v", err)
	}

	if resp.Updates.ServerResponse.HTTPStatusCode == 200 {
		success = true
	}

	return success
}
