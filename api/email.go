package api

import (
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"net"
	"net/mail"
	"net/smtp"
	"os"

	"github.com/joho/godotenv"
	"github.com/sashabaranov/go-openai"
)

func SendEmail(body string, subj string, toEmail string) {
	// Import the email vairables
	originEmail := os.Getenv("EMAIL")
	password := os.Getenv("PASSWORD")

	from := mail.Address{"", originEmail}
	to := mail.Address{"", toEmail}

	// Setup headers
	headers := make(map[string]string) //creates an emptry map (dictionary that can be populated)
	headers["From"] = from.String()
	headers["To"] = to.String()
	headers["Subject"] = subj

	// Setup message
	message := ""
	for k, v := range headers {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + body

	//Connect to the remote SMTP server
	servername := "smtp.gmail.com:587"
	host, _, _ := net.SplitHostPort(servername)

	auth := smtp.PlainAuth("", originEmail, password, host)

	// TLS Config
	tlsconfig := &tls.Config{
		InsecureSkipVerify: true, //TODO:This should be changed for production implementation
		ServerName:         host,
	}

	c, err := smtp.Dial(servername)
	if err != nil {
		log.Fatal("Error:", err)
	}

	c.StartTLS(tlsconfig)

	// Auth
	if err = c.Auth(auth); err != nil {
		log.Panic(err)
	}

	// Set the sender and recipient first
	if err := c.Mail(originEmail); err != nil {
		log.Fatal(err)
	}

	if err := c.Rcpt("hdubb1.ubc@gmail.com"); err != nil {
		log.Fatal(err)
	}

	w, err := c.Data()
	if err != nil {
		log.Fatal(err)
	}

	_, err = w.Write([]byte(message))
	if err != nil {
		log.Panic(err)
	}

	err = w.Close()
	if err != nil {
		log.Fatal(err)
	}

	// Send the QUIT command and close the connection.
	err = c.Quit()
	if err != nil {
		log.Fatal(err)
	}

}

func draftEmail(product string, promptTemplatePath string, salesPersonName string, companyName string) (string, error) {
	err := godotenv.Load()
	if err != nil {
		return "", err
	}
	key := os.Getenv("OPEN_AI_KEY")
	client := openai.NewClient(key)

	prompt, err := createEmailRequestPrompt(promptTemplatePath, salesPersonName, companyName, product)
	if err != nil {
		return "", err
	}

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

func createEmailRequestPrompt(promptTemplatePath string, salesPersonName string, companyName string, product string) (string, error) {
	file, err := os.Open(promptTemplatePath)
	if err != nil {
		return "", err
	}

	emailByte, err := io.ReadAll(file)
	if err != nil {
		return "", err
	}

	emailString := string(emailByte)

	prompt := fmt.Sprintf(emailString, salesPersonName, companyName, product)

	return prompt, nil

}
