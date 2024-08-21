package wrapper

import (
	"github.com/aws/aws-sdk-go/service/sesv2"
	"github.com/gofiber/fiber/v2/log"
)

//go:generate mockgen -source=./ses.go -destination=./mock/ses_mock.go
type SESAPI interface {
	SendEmail(input *sesv2.SendEmailInput) (*sesv2.SendEmailOutput, error)
}

type mockSESAPI struct{}

func (m mockSESAPI) SendEmail(input *sesv2.SendEmailInput) (*sesv2.SendEmailOutput, error) {
	log.Infof("Mocking email send: %v", input)
	return nil, nil
}

func NewMockSESAPI() SESAPI {
	return mockSESAPI{}
}
