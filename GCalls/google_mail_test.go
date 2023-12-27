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
	emailInfo, _, err := GetLatestUnreadMessage(srv)
	if err != nil {
		t.Error(err)
	}

	fmt.Println(emailInfo)
}

func TestWatchPushNotification(t *testing.T) {
	srv := ConnectToGmailAPI()
	WatchPushNotification(srv)
}
