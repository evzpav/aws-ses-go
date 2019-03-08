package email

import (
	"bytes"
	"html/template"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
)

type Client struct {
	awsRegion          string
	awsAccessKeyId     string
	awsSecretAccessKey string
}

type EmailData struct {
	From         string
	To           []string
	Subject      string
	Body         string
	Text         string
	HTML         string
	ReplyTo      []string
	TemplateName string
	TemplateVars interface{}
}

func NewClient(awsRegion, awsAccessKeyId, awsSecretAccessKey string) *Client {
	return &Client{
		awsRegion:          awsRegion,
		awsAccessKeyId:     awsAccessKeyId,
		awsSecretAccessKey: awsSecretAccessKey,
	}
}

func (s *Client) Send(e EmailData) error {
	err := e.parseTemplate()
	if err != nil {
		log.Printf("Could not parse email template: %s\n", e.TemplateName)
		return err
	}
	err = s.sendMail(&e)
	if err != nil {
		log.Printf("Failed to send '%s' the email to %s\n", e.Subject, e.To)
		return err
	}
	log.Printf("Email '%s' has been sent to %s\n", e.Subject, e.To)
	return nil
}

func (e *EmailData) parseTemplate() error {
	t, err := template.ParseFiles(e.TemplateName)
	if err != nil {
		return err
	}
	buffer := new(bytes.Buffer)
	if err = t.Execute(buffer, e.TemplateVars); err != nil {
		return err
	}
	e.HTML = buffer.String()
	return nil
}

func (s *Client) newSesClient() *SesClient {
	return New(ses.New(session.New(&aws.Config{
		Region: aws.String(s.awsRegion),
		Credentials: credentials.NewStaticCredentials(
			s.awsAccessKeyId,
			s.awsSecretAccessKey,
			"",
		),
	})))

}

func (s *Client) sendMail(e *EmailData) error {
	return s.newSesClient().SendSesEmail(e)
}
