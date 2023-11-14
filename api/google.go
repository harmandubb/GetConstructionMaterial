package api

import (
	"context"
	"fmt"

	"github.com/joho/godotenv"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
)

func connectToGmail() error {
	err := godotenv.Load()
	if err != nil {
		return err
	}
	gmailService, err := gmail.NewService(context.Background(), option.WithScopes(gmail.GmailReadonlyScope))
	if err != nil {
		return err
	}

	resp, err := gmailService.Users.Messages.List("harmandubb@docstruction.com").Do()
	if err != nil {
		return err
	}

	fmt.Println(resp)

	return nil
}
