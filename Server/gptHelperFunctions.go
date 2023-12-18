package server

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
	"github.com/sashabaranov/go-openai"
)

// Purpose: Convert an int from gpt to a category that is defined in the prompt template
// Parameters:
// num int --> num to corrolate to the category
// Return:
// categroy string
func numToCategory(num int) string {
	switch num {
	case 1:
		return "Electrical"
	case 2:
		return "Plumbing"
	case 3:
		return "Fire Prevention"
	case 4:
		return "HVAC"
	case 5:
		return "Welding"
	case 6:
		return "Roofing"
	default:
		return "Unknown"
	}
}

// Purpose: General function that creates a rpompt in the form that gpt understands
// Parameters:
// templatePath string --> to the template that you want to use
// promptVals ...string --> needed amount of strings that are needed in the prompt to complete
// Return:
// prompt string --> prompt in the format that is sccepeted by gpt
// error if present

func createGPTPrompt(promptTemplatePath string, promptVals ...string) (string, error) {
	file, err := os.Open(promptTemplatePath)
	if err != nil {
		return "", err
	}

	prompt_Bytes, err := io.ReadAll(file)
	if err != nil {
		return "", err
	}

	prompt_String := string(prompt_Bytes)

	// Convert []string to []interface{}
	args := make([]interface{}, len(promptVals))
	for i, v := range promptVals {
		args[i] = v
	}

	prompt := fmt.Sprintf(prompt_String, args...)

	return prompt, nil

}

// Purpose: creste the prompt for gpt to catogirize the materila
// Parameters:
// promptTemplatePath string --> file path name to the template prompt
// material string --> material that you want to catogorize
// Return:
// prompt string --> completed prompt based on the template and material name
func createMaterialCategorizationPrompt(promptTemplatePath string, material string) (string, error) {
	return createGPTPrompt(promptTemplatePath, material)
}

// Purpose: General funciton that prompts gpt
// Parameters:
// prompt string --> the prompt that you want to send to gpt
// Return:
// resp string --> answer from gpt
// error if there are any present
func promptGPT(prompt string) (string, error) {
	err := godotenv.Load()
	if err != nil {
		return "", err
	}
	key := os.Getenv("OPEN_AI_KEY")
	client := openai.NewClient(key)

	resp, err := client.CreateChatCompletion(
		context.TODO(),
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: prompt,
				},
			},
		},
	)
	if err != nil {
		fmt.Printf("ChatCompletion error: %v\n", err)
		return "", err
	}

	return resp.Choices[0].Message.Content, nil

}

// Purpose: get a category name for material
// Parameters:
// promptTemplatePath string --> file path to the prompt template that you want to use
// material string --> material you are trying to categorize
// Return:
// category string --> returns category, this is the categories that are specified in the prompt and should match the parsing function used to convert an int to the caterogy
// error if present

func PromptGPTMaterialCatogorization(promptTemplatePath string, material string) (string, error) {
	prompt, err := createMaterialCategorizationPrompt(promptTemplatePath, material)
	if err != nil {
		return "", err
	}

	resp, err := promptGPT(prompt)
	if err != nil {
		return "", err
	}

	num, err := strconv.Atoi(resp)
	if err != nil {
		return "", err
	}

	category := numToCategory(num)

	return category, nil

}

func CreateEmailToSupplier(promptTemplatePath string, supplier string, material string) (string, string, error) {
	emailMaterialRequestPrompt, err := createEmailMaterialRequestPrompt(promptTemplatePath, "", supplier, material)
	if err != nil {
		return "", "", err
	}

	gptEmailStructure, err := promptGPT(emailMaterialRequestPrompt)
	if err != nil {
		return "", "", err
	}

	subj, body, err := parseGPTEmailResponse(gptEmailStructure)
	if err != nil {
		return "", "", err
	}

	return subj, body, nil
}

func createEmailMaterialRequestPrompt(promptTemplatePath string, salesPersonName string, companyName string, product string) (string, error) {
	file, err := os.Open(promptTemplatePath)
	if err != nil {
		return "", err
	}

	emailByte, err := io.ReadAll(file)
	if err != nil {
		return "", err
	}

	emailString := string(emailByte)

	if salesPersonName == "" {
		salesPersonName = "the sales team"
	}

	prompt := fmt.Sprintf(emailString, salesPersonName, companyName, product)

	return prompt, nil

}

func createReceiceEmailAnalysisPrompt(receiveAnalysisTemplatePath string, body string) (string, error) {
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

func gptAnalysisPresent(str string) bool {

	re := regexp.MustCompile(`present: ([yn])`)
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
