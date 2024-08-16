package email

import (
	"errors"
	mock_wrapper "github.com/brunodmartins/church-members-api/platform/aws/wrapper/mock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"testing"
)

const subject = "Test Subject"
const receiver1 = "receiver1@email.com"
const receiver2 = "receiver2@email.com"

func TestEmailSender_SendEmail(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ses := mock_wrapper.NewMockSESAPI(ctrl)
	service := NewEmailService(ses, "")
	data := map[string]string{
		"Title": "Test",
		"Body":  "Body",
	}
	t.Run("Given a valid email template, when send it with valid data, then succeeds", func(t *testing.T) {
		ses.EXPECT().SendEmail(gomock.Any()).Return(nil, nil)
		assert.NoError(t, service.SendTemplateEmail(TestTemplate, data, subject, receiver1, receiver2))
	})
	t.Run("Given a invalid email template, when send it with valida data, then fails", func(t *testing.T) {
		assert.Error(t, service.SendTemplateEmail("", data, subject, receiver1, receiver2))
	})
	t.Run("Given a valid email template, when send it with valid data, then fails due to SES", func(t *testing.T) {
		ses.EXPECT().SendEmail(gomock.Any()).Return(nil, errors.New("error"))
		assert.Error(t, service.SendTemplateEmail(TestTemplate, data, subject, receiver1, receiver2))
	})
}

func TestEmailSender_buildEmailInput(t *testing.T) {
	const emailBody = "<html></html>"
	const sender = "sender@email.com"

	emailInput := buildEmailInput(sender, subject, emailBody, receiver1, receiver2)
	assert.Equal(t, sender, *emailInput.FromEmailAddress)
	assert.Len(t, emailInput.Destination.ToAddresses, 2)
	assert.Equal(t, receiver1, *emailInput.Destination.ToAddresses[0])
	assert.Equal(t, receiver2, *emailInput.Destination.ToAddresses[1])
	assert.NotNil(t, emailInput.Content.Simple.Subject)
	assert.Equal(t, subject, *emailInput.Content.Simple.Subject.Data)
	assert.NotNil(t, emailInput.Content.Simple.Body)
	assert.Equal(t, emailBody, *emailInput.Content.Simple.Body.Html.Data)
}
