package gcalls

import (
	"context"
	"fmt"
	"io"
	"sync/atomic"
	"time"

	"cloud.google.com/go/pubsub"
	"google.golang.org/api/gmail/v1"
)

type ProductInfo struct {
	Date      time.Time
	Name      string
	Category  string
	Price     float64
	Currency  string
	DataSheet []string //multiple different types of datasheets can be added
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

// func MarkEmailAsRead(srv *gmail.Service, userID string, messageID string) error {
// 	// Create a ModifyMessageRequest to remove the UNREAD label
// 	modReq := &gmail.ModifyMessageRequest{
// 		RemoveLabelIds: []string{"UNREAD"},
// 	}

// 	// Call the Gmail API to modify the message
// 	_, err := srv.Users.Messages.Modify(userID, messageID, modReq).Do()
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }
