package main

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
)

func sendEmail(from, to, subject, body string) error {
	// Create a new AWS session
	session := session.Must(session.NewSession())
	// Create a new SES client
	svc := ses.New(session)

	// Create the email request
	input := &ses.SendEmailInput{
		Destination: &ses.Destination{
			ToAddresses: []*string{
				aws.String(to),
			},
		},
		Message: &ses.Message{
			Body: &ses.Body{
				Html: &ses.Content{
					Data: aws.String(body),
				},
			},
			Subject: &ses.Content{
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