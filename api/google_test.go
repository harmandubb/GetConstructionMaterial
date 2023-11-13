package api

import (
	"testing"
)

func TestReadGmailEmails(t *testing.T) {
	err := ReadGmailEmails()
	if err != nil {
		t.Error(err)
	}
}
