package gcalls

import (
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"strings"
	"time"

	"cloud.google.com/go/pubsub"
	"github.com/agnivade/levenshtein"
	"github.com/joho/godotenv"
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

	fmt.Printf("Sent Email to: %s\n", to)

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
// Error if present
func WatchPushNotification(srv *gmail.Service) (err error) {
	user_srv := gmail.NewUsersService(srv)

	watchRequest := gmail.WatchRequest{
		LabelFilterBehavior: "include",
		LabelIds:            []string{"INBOX"},
		TopicName:           "projects/getconstructionmaterial/topics/getconstructionmaterial-topic",
	}

	_, err = user_srv.Watch("info@docstruction.com", &watchRequest).Do()
	if err != nil {
		return err
	}

	return nil

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

// Purpose: get a list of the unread messages in the mail box
// Parameters:
// srv *gmail.Service --> gmail api access service pointer
// user string --> user email that you are checking for the messages
// Return:
// r *gmail.ListMessagesResponse --> list of undread message data
// Error if any present
func GetUnreadMessagesData(srv *gmail.Service, user string) (r *gmail.ListMessagesResponse, err error) {
	// queryString := "in:inbox and is:unread"
	// queryString := "in:inbox"
	// TEst
	queryString := "is:unread"

	r, err = srv.Users.Messages.List(user).Q(queryString).Do()
	if err != nil {
		return nil, err
	}

	fmt.Println("Checking the unread function", len(r.Messages))

	return r, nil
}

// Purpose: Get the contents of any message (intent is to be used to get the unread messages)
// Parmeters:
// srv *gmail.Service --> Gmail service access
// msg *gmail.Message --> gmail style metadata associated with the particular unreadMessage
// Return:
// emailINfo EmailInfo --> Struct to store the desired infromation from an email
// msgID string --> gmail api id that is used to get more informaiton around the email including attachement info
// error if present
func GetMessage(srv *gmail.Service, msg *gmail.Message, user string) (emailInfo EmailInfo, msgID string, err error) {
	msgData, err := srv.Users.Messages.Get(user, msg.Id).Do()
	if err != nil {
		return EmailInfo{}, "", err
	}

	headers := msgData.Payload.Headers // Check how the header structure looks like

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

	// I think I need to actually get the specific email.

	// for _, part := range msg.Payload.Parts {
	part := msgData.Payload

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

		attachment, err := srv.Users.Messages.Attachments.Get(user, msg.Id, part.Body.AttachmentId).Do()
		if err != nil {
			return EmailInfo{}, "", err
		}

		data, err := base64.URLEncoding.DecodeString(attachment.Data)
		if err != nil {
			return EmailInfo{}, "", err
		}

		// Save the attachment
		attachmentLoc := fmt.Sprintf("Attachment/%s", part.Filename)
		fmt.Println(attachmentLoc)
		os.WriteFile(attachmentLoc, data, 0644)

		attachementsLocations = append(attachementsLocations, attachmentLoc)
	}

	emailInfo = EmailInfo{
		Date:        time.Now(),
		Subj:        subj,
		From:        from,
		Body:        trimOriginalMessage(string(body)),
		Body_size:   bodySize,
		attachments: attachementsLocations,
	}

	return emailInfo, msg.Id, nil

}

// Purpose: The email body rep containes the original email sent below a tag. That is removed before chat gpt analysis
// Parameters:
// body string --> message body as is from the gmail api call
// Return:
// trimmedBody string --> Message without the original message.

func trimOriginalMessage(body string) (trimmedBody string) {
	i := strings.Index(body, "-----Original Message-----")
	if i > 0 {
		body = body[:i]
	}

	// Compile regular expression to match quoted text
	re := regexp.MustCompile(`(?m)^>.*$`)

	// Remove the quoted text
	trimmedBody = re.ReplaceAllString(body, "")

	// Trim whitespace
	trimmedBody = strings.TrimSpace(trimmedBody)

	return trimmedBody
}

// Purpose: Mark e read email by the api to be marked read once all of the processing has occured
// Parameters:
// srv *gmail.Service --> gmail connection service to api
// user string --> gmail for the user you want to impersinate to do the change
// emailID string --> id that repsersents the email you want to change
// Return:
// Error if present
func MarkEmailAsRead(srv *gmail.Service, user, emailID string) (err error) {
	userMessageService := gmail.NewUsersMessagesService(srv)

	modReq := gmail.ModifyMessageRequest{
		RemoveLabelIds: []string{"UNREAD"},
	}
	_, err = userMessageService.Modify(user, emailID, &modReq).Do()
	if err != nil {
		return err
	}

	return nil
}
