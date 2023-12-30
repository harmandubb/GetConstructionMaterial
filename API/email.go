package api

import (
	"errors"
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func CreateReceiceEmailAnalysisPrompt(receiveAnalysisTemplatePath string, body string) (string, error) {
	file, err := os.Open(receiveAnalysisTemplatePath)
	if err != nil {
		return "", err
	}

	emailByte, err := io.ReadAll(file)
	if err != nil {
		return "", err
	}

	emailString := string(emailByte)

	prompt := fmt.Sprintf(emailString, body)

	return prompt, nil

}

func parseGPTEmailResponse(gptResponse string) (string, string, error) {
	subStartIndex := strings.Index(gptResponse, "Subject:")

	if subStartIndex == -1 {
		return "", "", errors.New("subject line not found")
	}

	subEndIndex := strings.Index(gptResponse[subStartIndex:], "\n")

	flufLen := len("Subject: ")

	subj := gptResponse[subStartIndex+flufLen : subEndIndex]

	//The body of the paragraph should be the rest of the stirng
	body := gptResponse[subEndIndex+2:]

	return subj, body, nil

}

// func parseGPTAnalysisResponse(gptResponse string) (EmailProductInfo, error) {
// 	// // if present is avialable then we will continue the struct
// 	var emailProductInfo EmailProductInfo

// 	gptResponse = strings.ToLower(gptResponse)

// 	present := gptAnalysisPresent(gptResponse)

// 	emailProductInfo.Present = present

// 	if !present {
// 		emailProductInfo.Currency = ""
// 		emailProductInfo.Data_Sheet = false
// 		emailProductInfo.Price = 0

// 		return emailProductInfo, nil
// 	}

// 	price, currency := gptAnalysisPrice(gptResponse)
// 	emailProductInfo.Price = price
// 	emailProductInfo.Currency = currency

// 	datasheet := gptAnalysisDataSheet(gptResponse)
// 	emailProductInfo.Data_Sheet = datasheet

// 	return emailProductInfo, nil

// }

func gptAnalysisPrice(str string) (float64, string) {
	// Regular expression to find a pattern like "6.16 CAD"
	// and capture the currency code as well
	defaultCurrency := "CAD"

	re := regexp.MustCompile(`price: (\d+\.\d+)(?: ([A-Z]{3}))?`)
	matches := re.FindStringSubmatch(str)

	if len(matches) >= 2 {
		priceStr := matches[1]
		currency := defaultCurrency // Default currency value

		// Check if currency is captured
		if len(matches) > 2 && matches[2] != "" {
			currency = matches[2]
		}

		// Convert the extracted string to a float
		price, err := strconv.ParseFloat(priceStr, 64)
		if err != nil {
			fmt.Println("Error parsing price:", err)
			return 0, ""
		}

		// fmt.Printf("Extracted Price: %.2f %s\n", price, currency)

		if currency == "" {
			currency = defaultCurrency
		}

		return price, currency
	} else {
		fmt.Println("Price or currency not found")
	}

	return 0, ""
}

func gptAnalysisDataSheet(str string) bool {

	re := regexp.MustCompile(`datasheet: ([yn])`)
	matches := re.FindStringSubmatch(str)

	if len(matches) > 1 {
		presentValue := matches[1]

		switch presentValue {
		case "y":
			return true
		case "n":
			return false
		default:
			return false
		}
	} else {
		fmt.Println("Pattern not found")
		return false
	}
}

func extractProductName(str string) (string, error) {

	// str = strings.ToLower(str)
	re := regexp.MustCompile(`[Dd]ocstruction:\s*(.*?)\s*-`)
	matches := re.FindStringSubmatch(str)

	if len(matches) > 1 {
		return matches[1], nil
	} else {
		return "", errors.New("no product name found")
	}

}
