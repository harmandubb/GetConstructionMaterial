package gcalls

import (
	"fmt"
	"testing"
)

func TestConnectToGmail(t *testing.T) {
	srv := ConnectToGmailAPI()

	fmt.Println(srv)
}

func TestSentEmail(t *testing.T) {
	srv := ConnectToGmailAPI()
	msg, err := SendEmail(srv, "Test", "This is a test sending to hdubb1.ubc@gmail.com", "hdubb1.ubc@gmail.com")
	if err != nil {
		t.Error(err)
	}

	fmt.Println(msg)

}

func TestGetLatestUnreadMessage(t *testing.T) {
	srv := ConnectToGmailAPI()
	user := "info@docstruction.com"
	messages, err := GetUnreadMessagesData(srv, user)
	if err != nil {
		fmt.Printf("Retrive error for unread Messages: %v\n", err)
		t.Error(err)
	}

	if len(messages.Messages) == 0 {
		t.Errorf("No unread Messages Found")
	}

	for _, message := range messages.Messages {
		emailInfo, _, err := GetMessage(srv, message, user)
		fmt.Println("Email Body:", emailInfo.Body)
		if err != nil {
			fmt.Printf("Specific Message retrive error: %v\n", err)
			t.Error(err)
		}
	}

}

func TestWatchPushNotification(t *testing.T) {
	srv := ConnectToGmailAPI()
	WatchPushNotification(srv)
}
