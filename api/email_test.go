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

func TestGptAnalysisPresent(t *testing.T) {
	if gptAnalysisPresent("ajdflkasdlkfjalkdjflajdfl;jasd;lfjasdlfjsldkfj present: n") == true {
		t.Fail()
	}

	if gptAnalysisPresent("jfokajd present: y jdkfaodkfh;oaid") == false {
		t.Fail()
	}

}

func TestGptAnalysisPrice(t *testing.T) {
	price, currency := gptAnalysisPrice("Based on the provided response: present: y price: 6.16 CAD per unit datasheet: n")

	if price != 6.16 {
		t.Fail()
	}
	if currency != "CAD" {
		t.Fail()
	}

	price, currency = gptAnalysisPrice("Based on the provided response: present: y price: 6.16 per unit datasheet: n")

	if price != 6.16 {
		t.Fail()
	}
	if currency != "CAD" {
		t.Fail()
	}

	price, currency = gptAnalysisPrice("present: y price: 12.32 USD datasheet: n")

	if price != 12.32 {
		t.Fail()
	}
	if currency != "USD" {
		t.Fail()
	}

	price, currency = gptAnalysisPrice("present: y price: 12.32 datasheet: n")

	if price != 12.32 {
		t.Fail()
	}
	if currency != "CAD" {
		t.Fail()
	}

}

func TestGptAnalysisDataSheet(t *testing.T) {
	if gptAnalysisDataSheet(";jasd;lfjasdlfjsldkfj present: n") != false {
		t.Fail()
	}

	if gptAnalysisDataSheet("present: y price: 12.32 datasheet: n") != false {
		t.Fail()
	}

	if gptAnalysisDataSheet("present: y price: 12.32 datasheet: y") != true {
		t.Fail()
	}
}
