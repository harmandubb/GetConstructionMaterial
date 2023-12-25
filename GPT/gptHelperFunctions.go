package GPT

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/sashabaranov/go-openai"
)

type EmailMaterialInfo struct {
	Present    bool
	Price      float64
	Currency   string
	Data_Sheet bool
}

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
// promptTemplate string --> actual template that you want to use
// promptVals ...string --> needed amount of strings that are needed in the prompt to complete
// Return:
// prompt string --> prompt in the format that is sccepeted by gpt
// error if present

func createGPTPrompt(promptTemplate string, promptVals ...string) (string, error) {
	// file, err := os.Open(promptTemplatePath)
	// if err != nil {
	// 	return "", err
	// }

	// prompt_Bytes, err := io.ReadAll(file)
	// if err != nil {
	// 	return "", err
	// }

	// prompt_String := string(prompt_Bytes)
	prompt_String := promptTemplate
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
// promptTemplatePath string --> prompt in string form
// material string --> material that you want to catogorize
// Return:
// prompt string --> completed prompt based on the template and material name
func createMaterialCategorizationPrompt(promptTemplate string, material string) (string, error) {
	return createGPTPrompt(promptTemplate, material)
}

// Purpose: General funciton that prompts gpt
// Parameters:
// prompt string --> the prompt that you want to send to gpt
// Return:
// resp string --> answer from gpt
// error if there are any present
func promptGPT(prompt string) (string, error) {
	// err := godotenv.Load()
	// if err != nil {
	// 	return "", err
	// }
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
		fmt.Printf("Chat Completion error: %v\n", err)
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

// Purpose: Create an email using key words such as supplier name and material to a supplier
// Parameters:
// emailPrompt string --> The template that is used to write the email
// supplier string --> Supplier name to include in the email
// material string --> Material name to include in the email
// Return:
// subj string --> Subject of the email
// body string --> main messahe body of the email
// error if present

func CreateEmailToSupplier(emailPrompt string, supplier string, material string) (subj string, body string, err error) {
	emailMaterialRequestPrompt, err := createEmailMaterialRequestPrompt(emailPrompt, "", supplier, material)
	if err != nil {
		return "", "", err
	}

	gptEmailStructure, err := promptGPT(emailMaterialRequestPrompt)
	if err != nil {
		return "", "", err
	}

	subj, body, err = parseGPTEmailResponse(gptEmailStructure)
	if err != nil {
		return "", "", err
	}

	return subj, body, nil
}

// Purpose: creates the filled in prompt based on inputs and template
// Parameters:
// promptTemplate string --> The string of the prompt
// salesPersonName string --> name of the sales person if we know.
// companyName string --> Name of the supplier that is being contact in this emial to make it more personal
// Return:
// prompt string --> Filled in prompt ready to submit to chat gpt to write the email
// error if present
func createEmailMaterialRequestPrompt(promptTemplate string, salesPersonName string, companyName string, product string) (string, error) {
	// file, err := os.Open(promptTemplatePath)
	// if err != nil {
	// 	return "", err
	// }

	// emailByte, err := io.ReadAll(file)
	// if err != nil {
	// 	return "", err
	// }

	// emailString := string(emailByte)

	if salesPersonName == "" {
		salesPersonName = "the sales team"
	}

	prompt := fmt.Sprintf(promptTemplate, salesPersonName, companyName, product)

	return prompt, nil

}

// Purpose: Create a prompt for gpt that checks what the outcome of the supplier message
// Can have outcomes in the following catigory:
// 1. Success --> They have the product and they are providing some information
// 2. Fail --> They do not have the product
// TODO: can implement more behavours such as no but we have these products for future use.
// TODO: Please contact this person for further inquiry (either internally or external of the corporation)
// Parameters:
// AnalysisTemplate string --> template that is to be used in the situation
// Body string --> message body that is to be checked to see what the suplier is trying to say
// Return:
// prompt string --> the prompt that is suppose to be sent to chat gpt
// error if present
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

// Purpose: Breaks the subject and the body from the gpt response that is used in sending out an email
// Parameter:
// gptResponse string --> response form gpt that is the subject and body for the email to be written
// Return:
// subj string --> subject to email
// body string --> body message of an email
// error if any present
func parseGPTEmailResponse(gptResponse string) (subj string, body string, err error) {
	subStartIndex := strings.Index(gptResponse, "Subject:")

	if subStartIndex == -1 {
		return "", "", errors.New("subject line not found")
	}

	subEndIndex := strings.Index(gptResponse[subStartIndex:], "\n")

	flufLen := len("Subject: ")

	subj = gptResponse[subStartIndex+flufLen : subEndIndex]

	//The body of the paragraph should be the rest of the stirng
	body = gptResponse[subEndIndex+2:]

	return subj, body, nil

}

func parseGPTAnalysisMaterialResponse(gptResponse string) (EmailMaterialInfo, error) {
	// // if present is avialable then we will continue the struct
	var emailProductInfo EmailMaterialInfo

	gptResponse = strings.ToLower(gptResponse)

	values := strings.Split(gptResponse, ",")

	present := gptAnalysisMaterialPresent(values[0])

	emailProductInfo.Present = present
	emailProductInfo.Currency = ""
	emailProductInfo.Data_Sheet = false
	emailProductInfo.Price = 0

	if !present {
		return emailProductInfo, nil
	}

	if gptAnalysisSwitchStatement(values[1]) {
		price, currency := gptAnalysisPrice(gptResponse)
		emailProductInfo.Price = price
		emailProductInfo.Currency = currency
	}

	if gptAnalysisSwitchStatement(values[2]) {
		emailProductInfo.Data_Sheet = true
	}

	return emailProductInfo, nil

}

// Purpose: Function to quicly call to see if a field is present or not in an email
// Parameter: Value string --> value you are tyring to see if something is present
// Return: bool --> if a value is present

func gptAnalysisSwitchStatement(presentValue string) bool {
	switch presentValue {
	case "y", "yes":
		return true
	case "n", "no":
		return false
	default:
		return false
	}
}

// Purpose: Checks the gpt analysis to see if the material is deemed to be present.
// Parameter:
// str string --> analysis of the gpt response
// Return:
// bool --> if the product is present as analyized by chat gpt.

func gptAnalysisMaterialPresent(presentValue string) bool {
	return gptAnalysisSwitchStatement(presentValue)
}

// Purpose: given a EmailMaterialInfo struct what should I do to response to a user
// Parameter:
// emailMaterialInfo struct --> struct that contains the bools that shows certains fields present (price, datasheet)
func ReactToEmailInfoContents(emailMaterialInfo EmailMaterialInfo) (emailPrompt string, err error) {
	if emailMaterialInfo.Present {
		if emailMaterialInfo.Price > 0 {

		}
	}

	return "", nil
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
