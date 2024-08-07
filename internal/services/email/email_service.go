package email

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go/service/sesv2"
	"github.com/brunodmartins/church-members-api/platform/aws/wrapper"
)

//go:generate mockgen -source=./email_service.go -destination=./mock/email_service_mock.go
type Service interface {
	SendEmail(Command) error
}

type Command struct {
	Recipients []string
	Body       string
	Subject    string
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
					Html: &sesv2.Content{
						Charset: aws.String("ISO-8859-1"),
						Data:    aws.String(command.Body),
					},
				},
				Subject: &sesv2.Content{
					Charset: aws.String("ISO-8859-1"),
					Data:    aws.String(command.Subject),
				},
			},
		},
		Destination: &sesv2.Destination{
			ToAddresses: aws.StringSlice(command.Recipients),
		},
		FromEmailAddress: aws.String(fromEmail),
	}
}
