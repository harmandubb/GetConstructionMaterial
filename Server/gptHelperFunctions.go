package server

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/joho/godotenv"
	"github.com/sashabaranov/go-openai"
)

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

func CreateMaterialCategorizationPrompt(promptTemplatePath string, material string) (string, error) {
	return createGPTPrompt(promptTemplatePath, material)
}

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
