package gcalls

import (
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"

	"cloud.google.com/go/pubsub"
	"github.com/agnivade/levenshtein"
	"golang.org/x/oauth2/jwt"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
)

type EmailInfo struct {
	Date        time.Time
	Subj        string
	From        string
	Body        string
	Body_size   int64
	attachments []string
}

// Purpose: Provides a server connection to the gmail api that is set in the environmental variables
// Parameters: Non
// Return: Returns gmail service pointer to be used with outher functions
func ConnectToGmailAPI() *gmail.Service {
	ctx := context.Background()

	// err := godotenv.Load()
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
			"https://www.googleapis.com/auth/gmail.modify",
		},
		TokenURL: os.Getenv("TOKEN_URL"),
	}

	conf.Subject = "info@docstruction.com"

	client := conf.Client(ctx)

	srv, err := gmail.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Fatalf("Unable to connect to service %v", err)
	}

	return srv
}

// Puspose: Send an email through gmail api
// Parameters:
// srv *gmail.Service --> pointer to the gmail api enabled service
// subj string --> subject line
// msg string --> message body
// to string --> email to send to
// Return:
// *gmail.Message --> Pointer to either an emptry message struct or the message response form server
// error if any present
func SendEmail(srv *gmail.Service, subj, msg, to string) (*gmail.Message, error) {
	message, err := gmail.NewUsersMessagesService(srv).Send(
		"info@docstruction.com",
		prepMessage(subj, msg, to),
	).Do()

	if err != nil {
		return message, err
	}

	fmt.Printf("Sent Email to: %s", to)

	return message, nil

}

// Purpose: Create the message datatype to send an email with the gmail api
// Parameters:
// srv *gmail.Service --> pointer to the gmail api enabled service
// subj string --> subject line
// msg string --> message body
// to string --> email to send to
// Return:
// *gmail.Message --> message struct is returned to send an email
func prepMessage(subj, msg, to string) *gmail.Message {
	header := make(map[string]string)
	header["To"] = to
	header["Subject"] = subj
	header["Content-Type"] = "text/plain; charset=utf-8"
	header["Date"] = time.Now().Format(time.RFC1123Z)

	var headers []*gmail.MessagePartHeader
	for k, v := range header {
		headers = append(headers, &gmail.MessagePartHeader{Name: k, Value: v})
	}

	messagePart := &gmail.MessagePart{
		Body: &gmail.MessagePartBody{
			Data: base64.URLEncoding.EncodeToString([]byte(msg)),
		},
		Headers: headers,
	}

	message := gmail.Message{
		Payload: messagePart,
		Raw:     rawMessage(msg, header),
	}

	return &message
}

// Purpose: Create a raw message for the send of the gmail function becuase that is the only thing that is needed
// Parameters:
// msg string --> message body
// headers map[string]string --> contains the needed headers for the send transmission (to, Subject, Content-type)
// Return:
// raw string --> RFC 2822 format and base64url encoded string 2
func rawMessage(msg string, headers map[string]string) string {
	fullString := ""
	for key, value := range headers {
		fullString = fullString + fmt.Sprintf("%s: %s\r\n", key, value)
	}

	fullString = fullString + "\r\n" + msg

	// Convert the RFC822 formatted message to a base64URL encoded string
	base64RawMessage := base64.URLEncoding.EncodeToString([]byte(fullString))

	return base64RawMessage

}

// Purpose: test function to test the publish subscription functionality in google api for push notifcation testing
// Parameters:
//
//	w io.Writier --> Could not figure out in the time of writing
//	projectID string --> Google api related parmater
//	topicID string --> Google api related parameter, speciifc topic you want to send something to
//	msg string --> what you want to publish to the topic
//
// Return: error if there are any
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

// Purpose: request that the watch is enabled on the particular email that you want to received push notifications for
// Parameters:
// srv *gmail.Service --> connection to the gmail api service
// Return:
// NON
func WatchPushNotification(srv *gmail.Service) {
	user_srv := gmail.NewUsersService(srv)

	watchRequest := gmail.WatchRequest{
		LabelFilterBehavior: "include",
		LabelIds:            []string{"INBOX"},
		TopicName:           "projects/getconstructionmaterial/topics/getconstructionmaterial-topic",
	}

	user_srv.Watch("info@docstruction.com", &watchRequest)

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

// Purpose: Check the most recent unread message
// Parmeters:
// srv *gmail.Service --> Gmail service access
// Return:
// emailINfo EmailInfo --> Struct to store the desired infromation from an email
// msgID string --> gmail api id that is used to get more informaiton around the email including attachement info
// error if present
func GetLatestUnreadMessage(srv *gmail.Service) (emailInfo EmailInfo, msgID string, err error) {
	user := "me"

	queryString := "In:inbox and Is:unread"

	var empty EmailInfo

	r, err := srv.Users.Messages.List(user).Q(queryString).Do()
	if err != nil {
		log.Fatalf("Unable to retrieve messages: %v", err)
	}

	msgID = r.Messages[0].Id

	msg, err := srv.Users.Messages.Get(user, msgID).Do()
	if err != nil {
		return empty, "", err
	}

	headers := msg.Payload.Headers // Check how the header structure looks like

	var subj, from, currentHeader string
	var body []byte
	var bodySize int64

	for i := range headers {
		currentHeader = headers[i].Name

		switch currentHeader {
		case "Subject":
			subj = headers[i].Value
		// case "Date":
		// 	date = headers[i].Value
		case "From":
			from = headers[i].Value
		}
	}

	var attachementsLocations []string

	// for _, part := range msg.Payload.Parts {
	part := msg.Payload

	if len(part.Body.Data) > 0 {
		body, err = base64.URLEncoding.DecodeString(part.Body.Data)
		bodySize = part.Body.Size
		if err != nil {
			fmt.Printf("No message body present: %s", err)
		}

	}

	if err != nil {
		fmt.Println("Error")
	}
	if part.Filename != "" && part.Body.AttachmentId != "" {

		attachment, err := srv.Users.Messages.Attachments.Get(user, msgID, part.Body.AttachmentId).Do()
		if err != nil {
			return empty, "", err
		}

		data, err := base64.URLEncoding.DecodeString(attachment.Data)
		if err != nil {
			return empty, "", err
		}

		// Save the attachment
		attachmentLoc := fmt.Sprintf("Attachment/%s", part.Filename)
		fmt.Println(attachmentLoc)
		os.WriteFile(attachmentLoc, data, 0644)

		attachementsLocations = append(attachementsLocations, attachmentLoc)
	}
	// }

	emailInfo = EmailInfo{
		Date:        time.Now(),
		Subj:        subj,
		From:        from,
		Body:        trimOriginalMessage(string(body)),
		Body_size:   bodySize,
		attachments: attachementsLocations,
	}

	return emailInfo, msgID, nil

}

func trimOriginalMessage(body string) (trimmedBody string) {
	trimmedBody = body[:strings.Index(body, "-----Original Message-----")]

	return trimmedBody
}
