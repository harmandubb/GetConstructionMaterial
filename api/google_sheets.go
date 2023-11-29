package api

import (
	"context"
	"docstruction/getconstructionmaterial/server"
	"log"
	"os"

	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

func sendEmailInfo(emailFormInfo server.EmailFormInfo, spreadSheetID string) bool {
	srv := connectToSheetsAPI()
	return appendToSpreadSheet(srv, spreadSheetID, emailFormInfo)

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

func appendToSpreadSheet(srv *sheets.Service, id string, emaillFormInfo server.EmailFormInfo) bool {
	success := false

	values := sheets.ValueRange{
		MajorDimension: "ROWS",
		Values: [][]interface{}{
			{emaillFormInfo.Time, emaillFormInfo.Email},
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
