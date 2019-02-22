package ses

import (
	"bytes"
	"html/template"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
)

type Request struct {
	awsRegion          string
	awsAccessKeyId     string
	awsSecretAccessKey string
	from               string
	to                 []string
	subject            string
	body               string
	text               string
	html               string
	replyTo            []string
}

func NewRequest(sender, noReplyEmail string, to []string, subject string) *Request {
	return &Request{
		from:    sender,
		to:      to,
		subject: subject,
		replyTo: []string{noReplyEmail},
	}
}

func (r *Request) SetAwsCredentials(awsRegion, awsAccessKeyId, awsSecretAccessKey string) *Request {
	r.awsRegion = awsRegion
	r.awsAccessKeyId = awsAccessKeyId
	r.awsSecretAccessKey = awsSecretAccessKey
	return r
}

func (r *Request) parseTemplate(fileName string, data interface{}) error {
	filePath := fileName
	t, err := template.ParseFiles(filePath)
	if err != nil {
		return err
	}
	buffer := new(bytes.Buffer)
	if err = t.Execute(buffer, data); err != nil {
		return err
	}
	r.html = buffer.String()
	return nil
}

func (r *Request) Send(templateName string, items interface{}) error {
	err := r.parseTemplate(templateName, items)
	if err != nil {
		log.Printf("Could not parse email template: %s\n", templateName)
		return err
	}
	err = r.sendMail()
	if err != nil {
		log.Printf("Failed to send '%s' the email to %s\n", r.subject, r.to)
		return err
	}
	log.Printf("Email '%s' has been sent to %s\n", r.subject, r.to)
	return nil
}

func (r *Request) newSesClient() *SesClient {
	return New(ses.New(session.New(&aws.Config{
		Region: aws.String(r.awsRegion),
		Credentials: credentials.NewStaticCredentials(
			r.awsAccessKeyId,
			r.awsSecretAccessKey,
			"",
		),
	})))

}

func (r *Request) sendMail() error {
	return r.newSesClient().SendSesEmail(*r)
}
