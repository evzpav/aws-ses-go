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
	From           string
	To             []string
	CC             []string
	BCC            []string
	ReplyTo        []string
	Subject        string
	Text           string
	HTML           string
	TemplateName   string
	TemplateVars   interface{}
	AttachFiles    []string
	ConfigSet      string
	BaseLayoutPath string
}

func NewClient(awsRegion, awsAccessKeyId, awsSecretAccessKey string) *Client {
	return &Client{
		awsRegion:          awsRegion,
		awsAccessKeyId:     awsAccessKeyId,
		awsSecretAccessKey: awsSecretAccessKey,
	}
}

//Send email based on HTML template
func (s *Client) Send(mail EmailData) error {
	err := mail.parseTemplate()
	if err != nil {
		log.Printf("Could not parse email template: %s\n", mail.TemplateName)
		return err
	}
	err = s.sendMail(&mail)
	if err != nil {
		log.Printf("Failed to send '%s' the email to %s\n", mail.Subject, mail.To)
		return err
	}
	log.Printf("Email '%s' has been sent to %s\n", mail.Subject, mail.To)
	return nil
}

//SendRaw email based on HTML template with attachment
func (s *Client) SendRaw(mail EmailData) error {
	err := mail.parseTemplate()
	if err != nil {
		log.Printf("Could not parse raw email template: %s\n", mail.TemplateName)
		return err
	}

	err = s.sendRawMail(&mail)
	if err != nil {
		log.Printf("Failed to send '%s' the raw email to %s\n", mail.Subject, mail.To)
		return err
	}
	log.Printf("Email Raw '%s' has been sent to %s\n", mail.Subject, mail.To)
	return nil
}

// parseTemplate HTML with provided variables
func (mail *EmailData) parseTemplate() error {
	var t *template.Template
	var err error

	if mail.BaseLayoutPath != "" {
		t, err = template.ParseFiles(mail.BaseLayoutPath, mail.TemplateName)
	} else if mail.TemplateName != "" {
		t, err = template.ParseFiles(mail.TemplateName)
	} else {
		return nil
	}

	if err != nil {
		return err
	}
	buffer := new(bytes.Buffer)
	if err = t.Execute(buffer, mail.TemplateVars); err != nil {
		return err
	}
	mail.HTML = buffer.String()
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

func (s *Client) sendMail(mail *EmailData) error {
	return s.newSesClient().SendSesEmail(mail)
}

func (s *Client) sendRawMail(mail *EmailData) error {
	return s.newSesClient().SendSesRawEmail(mail)
}
