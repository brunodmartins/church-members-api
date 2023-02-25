package middleware

import (
	"context"
	"github.com/brunodmartins/church-members-api/platform/i18n"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/text/language"
)

const languageHeader = "Accept-language"

var I18NMiddleware = func(ctx *fiber.Ctx) error {
	localize := i18n.GetLocalize(language.MustParse(ctx.Get(languageHeader, "en")))
	ctx.SetUserContext(context.WithValue(ctx.UserContext(), "i18n", localize))
	return ctx.Next()
}
