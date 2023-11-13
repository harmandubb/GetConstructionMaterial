package api

import (
	"testing"
)

func TestSendEmail(t *testing.T) {

}

func TestdraftEmail(t *testing.T) {
	draftEmail("Meta Caulk Fire Stop Collar", "../email_prompt.txt", "")
}
