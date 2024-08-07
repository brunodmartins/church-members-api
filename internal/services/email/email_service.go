package email

import (
	"bytes"
	"embed"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go/service/sesv2"
	"github.com/brunodmartins/church-members-api/platform/aws/wrapper"
	"html/template"
)

//go:generate mockgen -source=./email_service.go -destination=./mock/email_service_mock.go
type Service interface {
	SendTemplateEmail(template string, data any, subject string, recipients ...string) error
}

//go:embed templates/*
var fs embed.FS

const WeeklyBirthTemplate = "weekly_birth_template"
const TestTemplate = "test_template"
const ConfirmEmailTemplate = "confirm_email_template"

type emailService struct {
	ses       wrapper.SESAPI
	fromEmail string
	templates map[string]*template.Template
}

func NewEmailService(ses wrapper.SESAPI, fromEmail string) Service {
	var templates = make(map[string]*template.Template)
	templates[WeeklyBirthTemplate] = createTemplate(WeeklyBirthTemplate, string(loadEmailHTML("templates/weekly_birthdays_template.html")))
	templates[TestTemplate] = createTemplate(TestTemplate, string(loadEmailHTML("templates/test.html")))
	templates[ConfirmEmailTemplate] = createTemplate(TestTemplate, string(loadEmailHTML("templates/confirm_email_template.html")))
	return emailService{ses: ses, fromEmail: fromEmail, templates: templates}
}

func (service emailService) SendTemplateEmail(template string, data any, subject string, recipients ...string) error {
	emailTemplate, found := service.templates[template]
	if !found {
		return fmt.Errorf("template %s not found", template)
	}
	body, err := parseTemplate(emailTemplate, data)
	if err != nil {
		return err
	}
	_, err = service.ses.SendEmail(buildEmailInput(service.fromEmail, subject, body, recipients...))
	return err
}

func buildEmailInput(fromEmail string, subject string, body string, recipients ...string) *sesv2.SendEmailInput {
	return &sesv2.SendEmailInput{
		Content: &sesv2.EmailContent{
			Simple: &sesv2.Message{
				Body: &sesv2.Body{
					Html: &sesv2.Content{
						Charset: aws.String("ISO-8859-1"),
						Data:    aws.String(body),
					},
				},
				Subject: &sesv2.Content{
					Charset: aws.String("ISO-8859-1"),
					Data:    aws.String(subject),
				},
			},
		},
		Destination: &sesv2.Destination{
			ToAddresses: aws.StringSlice(recipients),
		},
		FromEmailAddress: aws.String(fromEmail),
	}
}

func loadEmailHTML(fileName string) []byte {
	data, err := fs.ReadFile(fileName)
	if err != nil {
		panic(err)
	}
	return data
}

func createTemplate(name string, html string) *template.Template {
	return template.Must(template.New(name).Parse(html))
}

func parseTemplate(emailTemplate *template.Template, data any) (string, error) {
	buf := new(bytes.Buffer)
	err := emailTemplate.Execute(buf, data)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}
