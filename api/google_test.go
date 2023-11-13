package api

import (
	"testing"
)

func TestReadGmailEmails(t *testing.T) {
	err := connectToGmail()
	if err != nil {
		t.Error(err)
	}
}
