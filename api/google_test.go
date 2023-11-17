package api

import (
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
