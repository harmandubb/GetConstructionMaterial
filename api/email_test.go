package api

import (
	"fmt"
	"testing"
)

// func TestSendEmail(t *testing.T) {

// }

func TestDraftEmail(t *testing.T) {
	emailDraft, err := draftEmail("Meta Caulk Fire Stop Collar", "../email_prompt.txt", "", "EECOL Electric")
	if err != nil {
		t.Error(err)
	}

	fmt.Println(emailDraft)

}
