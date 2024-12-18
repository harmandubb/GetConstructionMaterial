package gcalls

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"golang.org/x/oauth2/jwt"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

type MaterialFormInfo struct {
	Email    string
	Material string
	Loc      string
}

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

func SendMaterialFormInfo(spreadSheetID string, materialFormInfo MaterialFormInfo) bool {
	return AppendDataToSpreadSheet(spreadSheetID, time.Now(), materialFormInfo.Email, materialFormInfo.Material, materialFormInfo.Loc)
}

func AppendDataToSpreadSheet(spreadSheetID string, time time.Time, vals ...string) bool {
	success := false
	srv := ConnectToSheetsAPI()

	// Prepare the data for the ValueRange
	var data []interface{}

	data = append(data, time)
	for _, val := range vals {
		data = append(data, val)
	}

	values := sheets.ValueRange{
		Values: [][]interface{}{
			data,
		},
	}

	resp, err := srv.Spreadsheets.Values.Append(spreadSheetID, "Sheet1!A1", &values).ValueInputOption("USER_ENTERED").Do()
	if err != nil {
		log.Fatalf("Appending request Failed: %v", err)
	}

	if resp.ServerResponse.HTTPStatusCode == 200 {
		success = true
	}

	return success

}

// spreadsheet id: 1ZowyzJ008toPYNn0mFc2wG6YTAop9HfnbMPLIM4rRZw

func ConnectToSheetsAPI() *sheets.Service {
	ctx := context.Background()

	// err := godotenv.Load() // This will load your .env file
	// if err != nil {
	// 	log.Fatalf("Error loading .env file: %v", err)
	// }
	key := os.Getenv("PRIVATE_KEY")

	newkey := strings.Replace(key, "\\n", "\n", -1)

	pKey := []byte(newkey)

	conf := &jwt.Config{
		Email:        os.Getenv("CLIENT_EMAIL"),
		PrivateKeyID: os.Getenv("PRIVATE_KEY_ID"),
		PrivateKey:   pKey,
		Scopes: []string{
			"https://www.googleapis.com/auth/spreadsheets",
		},
		TokenURL: os.Getenv("TOKEN_URL"),
	}

	client := conf.Client(ctx)

	srv, err := sheets.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Fatalf("Unable to connnect to service %v", err)
		fmt.Println("Error after service", err)
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
