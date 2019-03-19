package email

import (
	"bytes"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ses"
	"github.com/aws/aws-sdk-go/service/ses/sesiface"
	gomail "gopkg.in/gomail.v2"
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

func (s *SesClient) SendSesEmail(mail *EmailData) error {
	input, err := createInput(mail)
	if err != nil {
		return err
	}

	_, err = s.ses.SendEmail(input)
	if err != nil {
		return err

	}
	return nil
}

func (s *SesClient) SendSesRawEmail(mail *EmailData) error {
	input, err := createRawInput(mail)
	if err != nil {
		return err
	}

	_, err = s.ses.SendRawEmail(input)
	if err != nil {
		return err

	}
	return nil
}

func createInput(mail *EmailData) (*ses.SendEmailInput, error) {
	if mail.HTML == "" {
		mail.HTML = mail.Text
	}
	msg := &ses.Message{
		Subject: &ses.Content{
			Charset: aws.String("utf-8"),
			Data:    &mail.Subject,
		},
		Body: &ses.Body{
			Html: &ses.Content{
				Charset: aws.String("utf-8"),
				Data:    &mail.HTML,
			},
			Text: &ses.Content{
				Charset: aws.String("utf-8"),
				Data:    &mail.Text,
			},
		},
	}

	dest := &ses.Destination{
		ToAddresses:  aws.StringSlice(mail.To),
		CcAddresses:  aws.StringSlice(mail.CC),
		BccAddresses: aws.StringSlice(mail.BCC),
	}

	return &ses.SendEmailInput{
		Source:               &mail.From,
		Destination:          dest,
		Message:              msg,
		ReplyToAddresses:     aws.StringSlice(mail.ReplyTo),
		ConfigurationSetName: aws.String(mail.ConfigSet),
	}, nil

}

func createRawInput(mail *EmailData) (*ses.SendRawEmailInput, error) {
	gm := gomail.NewMessage()
	gm.SetHeader("From", mail.From)
	gm.SetHeader("To", mail.To...)
	if len(mail.CC) > 0 {
		gm.SetHeader("Cc", mail.CC...)
	}
	if len(mail.CC) > 0 {
		gm.SetHeader("Bcc", mail.BCC...)
	}
	gm.SetHeader("Subject", mail.Subject)

	if mail.ConfigSet != "" {
		gm.SetHeader("X-SES-CONFIGURATION-SET", mail.ConfigSet)
	}

	var contentType string
	if mail.HTML != "" {
		contentType = "text/html;charset=UTF-8"
	} else {
		mail.HTML = mail.Text
		contentType = "text/plain;charset=UTF-8"
	}

	gm.SetBody(contentType, mail.HTML)

	for _, attachment := range mail.AttachFiles {
		gm.Attach(attachment)
	}

	var rawData bytes.Buffer
	if _, err := gm.WriteTo(&rawData); err != nil {
		return nil, err
	}

	return &ses.SendRawEmailInput{
		RawMessage: &ses.RawMessage{
			Data: rawData.Bytes(),
		},
		Source: aws.String(mail.From),
	}, nil
}
