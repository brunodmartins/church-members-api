package email

import (
	"errors"
	mock_wrapper "github.com/brunodmartins/church-members-api/platform/aws/wrapper/mock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"testing"
)

func TestEmailSender_SendEmail(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ses := mock_wrapper.NewMockSESAPI(ctrl)
	emailService := NewEmailService(ses, "")
	t.Run("Success", func(t *testing.T) {
		ses.EXPECT().SendEmail(gomock.Any()).Return(nil, nil)
		assert.Nil(t, emailService.SendEmail(buildCommand()))
	})
	t.Run("Fail", func(t *testing.T) {
		ses.EXPECT().SendEmail(gomock.Any()).Return(nil, errors.New("error"))
		assert.NotNil(t, emailService.SendEmail(buildCommand()))
	})
}

func TestEmailSender_buildEmailInput(t *testing.T) {
	emailInput := buildEmailInput(buildCommand(), "from@")
	assert.Equal(t, "from@", *emailInput.FromEmailAddress)
	assert.Len(t, emailInput.Destination.ToAddresses, 2)
	assert.NotNil(t, emailInput.Content.Simple.Subject)
	assert.NotNil(t, emailInput.Content.Simple.Body)
}

func buildCommand() Command {
	return Command{
		Recipients: []string{"test.com", "test1.com"},
		Body:       "test email",
		Subject:    "test Subject",
	}
}
