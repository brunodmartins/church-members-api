package i18n

import (
	"fmt"
	"github.com/spf13/viper"

	"github.com/BurntSushi/toml"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/sirupsen/logrus"
	"golang.org/x/text/language"
)

//Localizer var
var Localizer *i18n.Localizer

func init() {
	lang := language.English
	if envLang := viper.GetString("lang"); envLang != "" {
		lang = language.MustParse(envLang)
	}
	bundle := loadBundle(lang)
	Localizer = i18n.NewLocalizer(bundle)
}

func loadBundle(language language.Tag) *i18n.Bundle {
	bundle := i18n.NewBundle(language)
	bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)
	logrus.Infof("Loading bundle for language %s", language)
	path := fmt.Sprintf("./bundles/%s.toml", language)
	bundle.MustLoadMessageFile(path)
	logrus.Infof("Bundle %s loaded", language)
	return bundle
}
