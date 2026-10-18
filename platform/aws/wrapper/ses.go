package wrapper

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/service/sesv2"
	"github.com/gofiber/fiber/v2/log"
)

//go:generate mockgen -source=./ses.go -destination=./mock/ses_mock.go
type SESAPI interface {
	SendEmail(ctx context.Context, input *sesv2.SendEmailInput, optFns ...func(*sesv2.Options)) (*sesv2.SendEmailOutput, error)
}

type mockSESAPI struct{}

func (m mockSESAPI) SendEmail(ctx context.Context, input *sesv2.SendEmailInput, optFns ...func(*sesv2.Options)) (*sesv2.SendEmailOutput, error) {
	log.Infof("Mocking email send: %v", input)
	return nil, nil
}

func NewMockSESAPI() SESAPI {
	return mockSESAPI{}
}
