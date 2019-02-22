# aws-ses-go-example

## Example how to use AWS SES with Golang


### Usage example:

```go
package main

import (
	"github.com/joho/godotenv"
	"log"
	"os"
	"github.com/evzpav/aws-ses-go-example/ses"
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

	subject := "My email subject"
	noReplyEmail := "noreply@domain.com"
	r := ses.NewRequest(senderEmail, noReplyEmail, []string{receiverEmail}, subject)
	r = r.SetAwsCredentials(awsRegion, awsAccessKeyId, awsSecretAccessKey)
	err := r.Send("email_template.html", vars)

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

./aws-ses-go-example

```