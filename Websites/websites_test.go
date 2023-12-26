package websites

import (
	"fmt"
	"testing"
)

func TestFindEmailsOnPage(t *testing.T) {
	emails, err := FindEmailsOnPage("https://westernpacificelectrical.com/contact-us-2/")
	if err != nil {
		t.Error(err)
	}

	for _, value := range emails {
		fmt.Println(value)
	}
}

func TestFindContactURLOnPage(t *testing.T) {
	contactLinks, err := FindContactURLOnPage("https://www.gescan.com/content/surrey-branch?utm_source=Google_My_Business&utm_medium=google&utm_campaign=google-local")
	if err != nil {
		t.Error(err)
	}

	fmt.Println(contactLinks)
}

// func TestFindContactFormInputs(t *testing.T) {
// 	formInputs, err := FindContactFormInputs("https://ortechindustries.ca/contact")
// 	if err != nil {
// 		t.Error(err)
// 	}

// 	fmt.Println(formInputs)
// }
