package api

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"net"
	"net/mail"
	"net/smtp"
	"os"
	"regexp"
	"strings"

	"github.com/joho/godotenv"
	"github.com/sashabaranov/go-openai"
)

type EmailOptions struct {
	OriginEmail string
	Password    string
	ToEmail     string
	Subj        string
	Body        string
}

type EmailContents struct {
	Product     string
	SalesPerson string
	CompanyName string
}

type EmailProductInfo struct {
	Present    bool
	Price      float64
	Currency   string
	Data_Sheet bool
}

func SendEmail(body string, subj string, toEmail string) error {
	err := godotenv.Load()
	if err != nil {
		return err
	}

	// Import the email vairables
	originEmail := os.Getenv("EMAIL")
	password := os.Getenv("APPPASSWORD")

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
		return err
	}

	c.StartTLS(tlsconfig)

	// Auth
	if err = c.Auth(auth); err != nil {
		return err
	}

	// Set the sender and recipient first
	if err := c.Mail(originEmail); err != nil {
		return err
	}

	if err := c.Rcpt(toEmail); err != nil {
		return err
	}

	w, err := c.Data()
	if err != nil {
		return err
	}

	_, err = w.Write([]byte(message))
	if err != nil {
		return err
	}

	err = w.Close()
	if err != nil {
		return err
	}

	// Send the QUIT command and close the connection.
	err = c.Quit()
	if err != nil {
		return err
	}

	return nil

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

func parseGPTAnalysisResponse(gptResponse string) (EmailProductInfo, error) {
	// // if present is avialable then we will continue the struct
	var emailProductInfo EmailProductInfo

	// gptResponse = strings.ToLower(gptResponse)
	// presentIndex := strings.Index(gptResponse, "present")

	// i := 0
	// for true {

	// 	if string(gptResponse[presentIndex+i]) == "n" || string(gptResponse[presentIndex+i]) == "y" {
	// 		break
	// 	}
	// }

	// if string(gptResponse[presentIndex+i]) != "y" {
	// 	emailProductInfo.Present = false
	// 	emailProductInfo.Data_Sheet = false
	// 	emailProductInfo.Price = 0

	// 	return emailProductInfo, nil
	// }

	// emailProductInfo.Present = true

	// re := regexp.MustCompile(`\d+\.\d+`)
	// matches := re.FindStringSubmatch(gptResponse)

	// price, err := strconv.ParseFloat(matches[0])
	// if err != nil {
	// 	return emailProductInfo, err
	// }

	// emailProductInfo.Price = price

	// // See if I can extract a currency?

	return emailProductInfo, nil

}

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
