package api

import (
	"io"
	"os"
	"testing"
)

func TestSendEmail(t *testing.T) {
	file, err := os.Open("../GPToutput.txt")
	if err != nil {
		t.Error(err)
	}

	emailByte, err := io.ReadAll(file)
	if err != nil {
		t.Error(err)
	}

	gptEmailString := string(emailByte)

	subj, body, err := parseGPTEmailResponse(gptEmailString)
	if err != nil {
		t.Error(err)
	}

	err = SendEmail(body, subj, "hdubb1.ubc@gmail.com")
	if err != nil {
		t.Error(err)
	}
}

func TestDraftEmail(t *testing.T) {
	// // emailDraft, err := draftEmail("Meta Caulk Fire Stop Collar", "../email_prompt.txt", "", "EECOL Electric")
	// if err != nil {
	// 	t.Error(err)
	// }

	// fmt.Println(emailDraft)

}
