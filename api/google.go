package api

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"sync/atomic"
	"time"

	"cloud.google.com/go/pubsub"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"

	"github.com/agnivade/levenshtein"
)

type EmailInfo struct {
	Date         time.Time
	Subj         string
	From         string
	Body         string
	Body_size    int64
	attachmentID string
}

type ProductInfo struct {
	Date      time.Time
	Name      string
	Price     float64
	Currency  string
	DataSheet bool
}

// Retrieve a token, saves the token, then returns the generated client.
func getClient(config *oauth2.Config) *http.Client {
	// The file token.json stores the user's access and refresh tokens, and is
	// created automatically when the authorization flow completes for the first
	// time.
	tokFile := "token.json"
	tok, err := tokenFromFile(tokFile)
	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken(tokFile, tok)
	}
	return config.Client(context.Background(), tok)
}

// Request a token from the web, then returns the retrieved token.
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		log.Fatalf("Unable to read authorization code: %v", err)
	}

	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web: %v", err)
	}
	return tok
}

// Retrieves a token from a local file.
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

// Saves a token to a file path.
func saveToken(path string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

func ConnectToGmail() *gmail.Service {
	ctx := context.Background()
	// b, err := os.ReadFile("/Users/harmandeepdubb/Library/CloudStorage/OneDrive-Personal/Desktop/GetConstructionMaterial/Auth2/credentials.json")
	b, err := os.ReadFile("../Auth2/credentials.json")

	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	// If modifying these scopes, delete your previously saved token.json.
	config, err := google.ConfigFromJSON(b, gmail.GmailReadonlyScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	client := getClient(config)

	srv, err := gmail.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Fatalf("Unable to retrieve Gmail client: %v", err)
	}

	return srv

}

func checkMessage(srv *gmail.Service, subj string, loc string, body string) (bool, error) {
	user := "me"

	queryString := fmt.Sprintf("In:%s and Subject:%s", loc, subj)

	r, err := srv.Users.Messages.List(user).Q(queryString).Do()
	if err != nil {
		log.Fatalf("Unable to retrieve messages: %v", err)
	}

	for _, l := range r.Messages {
		msg, err := srv.Users.Messages.Get(user, l.Id).Do()
		if err != nil {
			log.Fatalln(err)
		}

		mesgBody, _ := base64.URLEncoding.DecodeString(msg.Payload.Body.Data)

		stringMesgBody := string(mesgBody)

		distance := levenshtein.ComputeDistance(body, stringMesgBody)

		fmt.Println(distance)

		if distance < 500 {
			return true, nil
		}

	}

	return false, err

}

func getLatestUnreadMessage(srv *gmail.Service) (EmailInfo, error) {
	user := "me"

	queryString := fmt.Sprintf("In:inbox and Is:unread")

	var empty EmailInfo

	r, err := srv.Users.Messages.List(user).Q(queryString).Do()
	if err != nil {
		log.Fatalf("Unable to retrieve messages: %v", err)
	}

	msg, err := srv.Users.Messages.Get(user, r.Messages[0].Id).Do()
	if err != nil {
		return empty, err
	}

	headers := msg.Payload.Headers // Check how the header structure looks like

	date, err := time.Parse("RFC1123", headers[16].Value)
	if err != nil {
		fmt.Printf("Time Parse Error: %s", err)
	}

	body, err := base64.URLEncoding.DecodeString(msg.Payload.Body.Data)
	if err != nil {
		fmt.Printf("No message body present: %s", err)
	}

	emailInfo := EmailInfo{
		Date:         date,
		Subj:         headers[22].Value,
		From:         headers[23].Value,
		Body:         string(body),
		Body_size:    msg.Payload.Body.Size,
		attachmentID: msg.Payload.Body.AttachmentId,
	}

	return emailInfo, nil

}

// TODO:
// 1. determine what product the email is related to
// 		- The subject of the email is likely to be the same as what is sent out
// 		- The name of the product would be be included in the email but may not be standard.
// 		- Could prompt Chat GPT to write the subject in a way that would encode what product we are looking for so I just need to use an algorithum to extract the product name
// 2. Does the sales person have the product
// 		- Just ask chat gpt if the emai shows a confrimation that the product is present and encode the information in a specific reply
// 3. What information did the sales person provide
// 		- Datasheet (attachement) --> Focus
// 		- price (in line or attachement) --> Focus
// 		- Brand (in line or attachement)
// 		- Model (in line or attachement)
// As an initial version of the app I can show the datasheet to the user for them to decide
// if what we have is what they are looking for.
// This information can be encoded in the response from chatgpt.
// 4. Who is the sales person or is it from the general email?
// 		- I can check the database if that company has that particular sales person present.
//		- this would need to be a seperate table to encode the information of the sales team at different locations.

func parseEmailResponseInfo(emailInfo EmailInfo) (bool, error) {

	// var product ProductInfo

	emailAnalysisiPrompt, err := createReceiceEmailAnalysisPrompt("../email_parse_prompt.txt", emailInfo.Body)
	if err != nil {
		return false, err
	}

	gptResp, err := promptGPT(emailAnalysisiPrompt)
	if err != nil {
		return false, err
	}

	// // Only continue with the insert if the stock is present

	emailProductInfo, err := parseGPTAnalysisResponse(gptResp)
	if err != nil {
		return false, err
	}

	if !emailProductInfo.Present {
		return false, err
	}

	name, err := extractProductName(emailInfo.Subj)
	if err != nil {
		return false, err
	}

	fmt.Println(time.Now())

	//extract datasheet if needed
	if emailProductInfo.Data_Sheet != false {
		//do something to get the datasheet from the email
	}

	product := ProductInfo{
		Date:      time.Now(),
		Name:      name,
		Price:     emailProductInfo.Price,
		Currency:  emailProductInfo.Currency,
		DataSheet: emailProductInfo.Data_Sheet,
	}

	//extract

	return true, nil

}

func publish(w io.Writer, projectID, topicID, msg string) error {
	projectID = "getconstructionmaterial"
	topicID = "getconstructionmaterial-topic"
	msg = "Hello World"
	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		return fmt.Errorf("pubsub: NewClient: %w", err)
	}
	defer client.Close()

	t := client.Topic(topicID)
	result := t.Publish(ctx, &pubsub.Message{
		Data: []byte(msg),
	})
	// Block until the result is returned and a server-generated
	// ID is returned for the published message.
	id, err := result.Get(ctx)
	if err != nil {
		return fmt.Errorf("pubsub: result.Get: %w", err)
	}
	fmt.Fprintf(w, "Published a message; msg ID: %v\n", id)
	return nil
}

func pullMsgs(w io.Writer, projectID, subID string) error {
	// projectID := "my-project-id"
	// subID := "my-sub"
	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		return fmt.Errorf("pubsub.NewClient: %w", err)
	}
	defer client.Close()

	sub := client.Subscription(subID)

	// Receive messages for 10 seconds, which simplifies testing.
	// Comment this out in production, since `Receive` should
	// be used as a long running operation.
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	var received int32
	err = sub.Receive(ctx, func(_ context.Context, msg *pubsub.Message) {
		fmt.Fprintf(w, "Got message: %q\n", string(msg.Data))
		atomic.AddInt32(&received, 1)
		msg.Ack()
	})
	if err != nil {
		return fmt.Errorf("sub.Receive: %w", err)
	}
	fmt.Fprintf(w, "Received %d messages\n", received)

	return nil
}

func pushNotificationSetUp(srv *gmail.Service) (*gmail.WatchResponse, error) {
	usersrv := gmail.NewUsersService(srv)

	watchRqst := gmail.WatchRequest{
		TopicName: "projects/getconstructionmaterial/topics/getconstructionmaterial-topic",
		LabelIds:  []string{"INBOX"},
	}

	watchResponse, err := usersrv.Watch("me", &watchRqst).Do()
	if err != nil {
		return nil, err
	}

	return watchResponse, nil

}

func getDataSheet(srv gmail.Service, datasheetID string) error {

}
