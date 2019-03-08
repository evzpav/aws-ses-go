package main

import (
	"log"
	"os"

	"github.com/evzpav/aws-ses-go/ses"
	"github.com/joho/godotenv"
)

var senderEmail string
var receiverEmail string
var awsRegion string
var awsAccessKeyId string
var awsSecretAccessKey string

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	senderEmail = os.Getenv("SENDER_EMAIL")
	receiverEmail = os.Getenv("RECEIVER_EMAIL")
	awsRegion = os.Getenv("AWS_REGION")
	awsAccessKeyId = os.Getenv("AWS_ACCESS_KEY_ID")
	awsSecretAccessKey = os.Getenv("AWS_SECRET_ACCESS_KEY")
}

func main() {

	vars := map[string]string{ //variables that will go to HTML template
		"name":         "evzpav",
		"userID":       "123456",
		"supportEmail": senderEmail,
	}

	s := ses.NewClient(awsRegion, awsAccessKeyId, awsSecretAccessKey)
	var emailData = ses.EmailData{
		From:         senderEmail,
		To:           []string{receiverEmail},
		ReplyTo:      []string{"noreply@domain.com"},
		Subject:      "My email subject",
		TemplateName: "email_template.html",
		TemplateVars: vars,
	}
	err := s.Send(emailData)

	if err != nil {
		log.Println(err)
	}
}
