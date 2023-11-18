package api

import (
	"fmt"
	"os"
	"testing"
)

// func TestReadGmailEmails(t *testing.T) {
// 	srv := ConnectToGmail()
// 	result, err := retrieveEmail(srv, "Docstruction", "Test", "sent")

// }

func TestPublish(t *testing.T) {
	err := publish(os.Stdout, "getconstructionmaterial", "getconstructionmaterial-topic", "Hello World")
	if err != nil {
		t.Error(err)
	}

	err = pullMsgs(os.Stdout, "getconstructionmaterial", "getconstructionmaterial-sub")
	if err != nil {
		t.Error(err)
	}

}

func TestPushNotificationSetUp(t *testing.T) {
	srv := ConnectToGmail()
	watchResponse, err := pushNotificationSetUp(srv)
	if err != nil {
		t.Error(err)
		t.Fail()
	}
	fmt.Println(watchResponse.Expiration)

	//TODO: Set up renewing of the watch daily
	//TODO: Set up an endpoint to receive the push notification updates.
}
