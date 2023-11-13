package api

import (
	"context"

	"google.golang.org/api/gmail/v1"
)

func connectToGmail() error {
	_, err := gmail.NewService(context.Background())

	if err != nil {
		return err
	}

	return nil
}
