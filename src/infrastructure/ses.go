package infrastructure

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	AWSses "github.com/aws/aws-sdk-go/service/ses"
)

type SESEmailRepository struct {
	session *session.Session
}

func NewSESEmailRepository(sess *session.Session) *SESEmailRepository {
	return &SESEmailRepository{session: sess}
}

func (ses *SESEmailRepository) SendEmail(from, to, subject, body string) error {
	// Create a new SES client
	svc := AWSses.New(ses.session)

	// Create the email request
	input := &AWSses.SendEmailInput{
		Destination: &AWSses.Destination{
			ToAddresses: []*string{
				aws.String(to),
			},
		},
		Message: &AWSses.Message{
			Body: &AWSses.Body{
				Html: &AWSses.Content{
					Data: aws.String(body),
				},
			},
			Subject: &AWSses.Content{
				Data: aws.String(subject),
			},
		},
		Source: aws.String(from),
	}

	// Send the email
	_, err := svc.SendEmail(input)
	if err != nil {
		return fmt.Errorf("Error sending email from SES client, err %#v", err)
	}

	return nil
}
