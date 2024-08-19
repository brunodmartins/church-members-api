package email

import (
	"context"
	"github.com/brunodmartins/church-members-api/platform/i18n"
)

type WeeklyBirthTemplateDTO struct {
	Title           string
	BirthTitle      string
	MarriageTitle   string
	NameColumn      string
	DateColumn      string
	MembersBirth    []MemberDTO
	MembersMarriage []MemberDTO
}

type MemberDTO struct {
	Name string
	Date string
}

type ConfirmEmailTemplateDTO struct {
	Title          string
	Message        string
	User           string
	Link           string
	UnaskedMessage string
	ConfirmButton  string
}

func NewWeeklyBirthTemplateDTO(ctx context.Context) WeeklyBirthTemplateDTO {
	return WeeklyBirthTemplateDTO{
		Title:         i18n.GetMessage(ctx, "Emails.WeeklyBirth.Title"),
		BirthTitle:    i18n.GetMessage(ctx, "Emails.WeeklyBirth.Birth"),
		MarriageTitle: i18n.GetMessage(ctx, "Emails.WeeklyBirth.Marriage"),
		NameColumn:    i18n.GetMessage(ctx, "Domain.Name"),
		DateColumn:    i18n.GetMessage(ctx, "Domain.Date"),
	}
}
func NewConfirmEmailTemplateDTO(ctx context.Context) ConfirmEmailTemplateDTO {
	return ConfirmEmailTemplateDTO{
		Title:          i18n.GetMessage(ctx, "Emails.ConfirmEmail.Title"),
		Message:        i18n.GetMessage(ctx, "Emails.ConfirmEmail.Message"),
		UnaskedMessage: i18n.GetMessage(ctx, "Emails.ConfirmEmail.UnaskedMessage"),
		ConfirmButton:  i18n.GetMessage(ctx, "Emails.ConfirmEmail.Button"),
	}
}
