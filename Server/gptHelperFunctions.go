package server

import (
	"context"
	"fmt"
	"io"
	"os"
	"strconv"

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
