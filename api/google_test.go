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
		t.FailNow()
	}

	err = pullMsgs(os.Stdout, "getconstructionmaterial", "getconstructionmaterial-sub")
	if err != nil {
		t.Error(err)
		t.FailNow()
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

func TestNewEmailReceive(t *testing.T) {
	srv := ConnectToGmail()
	_, err := pushNotificationSetUp(srv)
	if err != nil {
		t.Error(err)
		t.Fail()
	}

	// time.Sleep(time.Second * 30)

	err = pullMsgs(os.Stdout, "getconstructionmaterial", "getconstructionmaterial-sub")
	if err != nil {
		t.Error(err)
		t.Fail()
	}

}

func TestGetLatestUnreadMessage(t *testing.T) {
	srv := ConnectToGmail()
	getLatestUnreadMessage(srv)
}

func TestExtractProductName(t *testing.T) {
	product, err := extractProductName("Subject: Docstruction: Fire Stop Collars - Got Any in Stock?")
	if err != nil {
		t.Fail()
	}

	if product != "Fire Stop Collars" {
		t.Fail()
	}
}
