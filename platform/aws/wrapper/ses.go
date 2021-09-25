package wrapper

import "github.com/aws/aws-sdk-go/service/sesv2"

//go:generate mockgen -source=./ses.go -destination=./mock/ses_mock.go
type SESAPI interface {
	SendEmail(input *sesv2.SendEmailInput) (*sesv2.SendEmailOutput, error)
}
