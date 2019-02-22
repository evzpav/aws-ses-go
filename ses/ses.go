package ses

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

func (c *SesClient) SendSesEmail(r Request) error {
	if r.html == "" {
		r.html = r.text
	}

	msg := &ses.Message{
		Subject: &ses.Content{
			Charset: aws.String("utf-8"),
			Data:    &r.subject,
		},
		Body: &ses.Body{
			Html: &ses.Content{
				Charset: aws.String("utf-8"),
				Data:    &r.html,
			},
			Text: &ses.Content{
				Charset: aws.String("utf-8"),
				Data:    &r.text,
			},
		},
	}

	dest := &ses.Destination{
		ToAddresses: aws.StringSlice(r.to),
	}

	output, err := c.ses.SendEmail(&ses.SendEmailInput{
		Source:           &r.from,
		Destination:      dest,
		Message:          msg,
		ReplyToAddresses: aws.StringSlice(r.replyTo),
	})
	log.Println(output)

	return err
}
