package api

import (
	"context"
	"fmt"

	"github.com/joho/godotenv"
	"google.golang.org/api/gmail/v1"
)

func connectToGmail() error {
	err := godotenv.Load()
	if err != nil {
		return err
	}
	gmailService, err := gmail.NewService(context.Background())
	if err != nil {
		return err
	}

	messageList, err := gmailService.Users.Messages.List("me").Do()
	if err != nil {
		return err
	}

	fmt.Println(messageList)

	return nil
}
