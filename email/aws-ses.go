package email

import (
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ses"
	"github.com/aws/aws-sdk-go/service/ses/sesiface"
)

// Client for SES.
type SesClient struct {
	ses sesiface.SESAPI // Service implementation
}

// New client.
func New(ses sesiface.SESAPI) *SesClient {
	return &SesClient{
		ses: ses,
	}
}

func (c *SesClient) SendSesEmail(r *EmailData) error {
	if r.HTML == "" {
		r.HTML = r.Text
	}
	msg := &ses.Message{
		Subject: &ses.Content{
			Charset: aws.String("utf-8"),
			Data:    &r.Subject,
		},
		Body: &ses.Body{
			Html: &ses.Content{
				Charset: aws.String("utf-8"),
				Data:    &r.HTML,
			},
			Text: &ses.Content{
				Charset: aws.String("utf-8"),
				Data:    &r.Text,
			},
		},
	}

	dest := &ses.Destination{
		ToAddresses: aws.StringSlice(r.To),
	}

	output, err := c.ses.SendEmail(&ses.SendEmailInput{
		Source:           &r.From,
		Destination:      dest,
		Message:          msg,
		ReplyToAddresses: aws.StringSlice(r.ReplyTo),
	})
	log.Println(output)

	return err
}
