# aws-ses-go

## AWS SES with Golang package. Parses HTML templates with dynamic variables and with attachment files.


### Usage example:

```go
package main

import (
	"log"
	"os"

	"github.com/evzpav/aws-ses-go/email"
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
		"info":         "123456",
		"supportEmail": senderEmail,
	}

	s := email.NewClient(awsRegion, awsAccessKeyId, awsSecretAccessKey)

	var baseLayoutVars = map[string]string{
		"supportEmail": senderEmail,
	}

	for k, v := range baseLayoutVars {
		vars[k] = v
	}

	//email based on HTML template
	var emailData = email.EmailData{
		From:           senderEmail,
		To:             []string{receiverEmail},
		ReplyTo:        []string{"noreply@domain.com"},
		Subject:        "My email subject",
		TemplateName:   "examples/layout-and-attachment/html-templates/registration.html",
		TemplateVars:   vars,
		BaseLayoutPath: "examples/layout-and-attachment/html-templates/base_layout.html",
	}
	err := s.Send(emailData)

	if err != nil {
		log.Println(err)
	}

	//email based on HTML template with attachment
	var emailWithAttachment = email.EmailData{
		From:         senderEmail,
		To:           []string{receiverEmail},
		ReplyTo:      []string{"noreply@domain.com"},
		Subject:      "My email subject",
		TemplateName: "email_template.html",
		TemplateVars: vars,
		AttachFiles: []string{"examples/layout-and-attachment/attachment-example/attachment.pdf",
			"examples/layout-and-attachment/attachment-example/attachment.txt"},
	}
	err = s.SendRaw(emailWithAttachment)

	if err != nil {
		log.Println(err)
	}
}


```

### To run example:
```bash
cp .env_example .env
# Fill credentials and emails on .env file

go build

./aws-ses-go

```