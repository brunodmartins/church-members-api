package i18n

import (
	"context"
	"embed"
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/sirupsen/logrus"
	"golang.org/x/text/language"
)

//go:embed languages/*.toml
var LocaleFS embed.FS

var bundles = make(map[string]*i18n.Localizer)

func GetMessage(ctx context.Context, key string) string {
	localize := GetLocalize(language.English)
	if ctx.Value("i18n") != nil {
		localize = ctx.Value("i18n").(*i18n.Localizer)
	}
	return localize.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID: key,
		},
	})
}

func loadBundle(language language.Tag) *i18n.Bundle {
	bundle := i18n.NewBundle(language)
	bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)
	logrus.Infof("Loading bundle for language %s", language)
	_, _ = bundle.LoadMessageFileFS(LocaleFS, fmt.Sprintf("languages/%s.toml", language.String()))
	logrus.Infof("Bundle %s loaded", language)
	return bundle
}

// GetLocalize returns an i18n.Localize based on a language tag
func GetLocalize(language language.Tag) *i18n.Localizer {
	if bundles[language.String()] != nil {
		return bundles[language.String()]
	}
	result := i18n.NewLocalizer(loadBundle(language))
	bundles[language.String()] = result
	return result
}
