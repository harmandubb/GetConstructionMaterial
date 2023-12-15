package gcalls

import (
	"context"
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"golang.org/x/oauth2/jwt"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
)

func ConnectToGmailAPI() *gmail.Service {
	ctx := context.Background()

	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

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

	srv, err := gmail.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Fatalf("Unable to connect to service %v", err)
	}

	return srv
}
