package email

import (
	"github.com/BrunoDM2943/church-members-api/platform/aws/wrapper"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go/service/sesv2"
)

type Service interface {
	SendEmail(Command) error
}

type Command struct {
	recipients []string
	body       string
	subject    string
}

type emailService struct {
	ses       wrapper.SESAPI
	fromEmail string
}

func NewEmailService(ses wrapper.SESAPI, fromEmail string) Service {
	return emailService{ses: ses, fromEmail: fromEmail}
}

func (e emailService) SendEmail(command Command) error {
	_, err := e.ses.SendEmail(buildEmailInput(command, e.fromEmail))
	return err
}

func buildEmailInput(command Command, fromEmail string) *sesv2.SendEmailInput {
	return &sesv2.SendEmailInput{
		Content: &sesv2.EmailContent{
			Simple: &sesv2.Message{
				Body: &sesv2.Body{
					Text: &sesv2.Content{
						Charset: aws.String("ISO-8859-1"),
						Data:    aws.String(command.body),
					},
				},
				Subject: &sesv2.Content{
					Charset: aws.String("ISO-8859-1"),
					Data:    aws.String(command.subject),
				},
			},
		},
		Destination: &sesv2.Destination{
			ToAddresses: aws.StringSlice(command.recipients),
		},
		FromEmailAddress: aws.String(fromEmail),
	}
}
