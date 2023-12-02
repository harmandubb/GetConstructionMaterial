package gcalls

import (
	"context"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

func getPath(relativePath string) string {
	_, b, _, _ := runtime.Caller(0)
	// The directory of the file
	basepath := filepath.Dir(b)
	// Construct the path relative to the file
	return filepath.Join(basepath, relativePath)
}

func SendEmailInfo(time time.Time, email string, spreadSheetID string) bool {
	srv := ConnectToSheetsAPI()
	return appendEmailToSpreadSheet(srv, spreadSheetID, time, email)

}

// spreadsheet id: 1ZowyzJ008toPYNn0mFc2wG6YTAop9HfnbMPLIM4rRZw

func ConnectToSheetsAPI() *sheets.Service {
	ctx := context.Background()
	b, err := os.ReadFile(getPath("../Auth2/credentials.json"))
	if err != nil {
		log.Fatalf("Unable to read crednetials: %v", err)
	}

	config, err := google.ConfigFromJSON(b, "https://www.googleapis.com/auth/spreadsheets")
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	client := getClient(config)

	srv, err := sheets.NewService(ctx, option.WithHTTPClient(client))
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

	resp, err := srv.Spreadsheets.Values.Append(id, "Sheet1!A1", &values).ValueInputOption("USER_ENTERED").Do()
	if err != nil {
		log.Fatalf("Appending request Failed: %v", err)
	}

	if resp.ServerResponse.HTTPStatusCode == 200 {
		success = true
	}

	return success
}
