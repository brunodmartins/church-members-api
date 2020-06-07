package i18n

import (
	"fmt"

	"github.com/BurntSushi/toml"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/sirupsen/logrus"
	"golang.org/x/text/language"
)

//Localizer var
var Localizer *i18n.Localizer

func init() {
	bundle := loadBundle(language.English)
	Localizer = i18n.NewLocalizer(bundle)
}

func loadBundle(language language.Tag) *i18n.Bundle {
	bundle := i18n.NewBundle(language)
	bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)
	languageStr := language.String()
	logrus.Infof("Loading bundle for language %s", language)
	bundle.MustLoadMessageFile(fmt.Sprintf("./bundles/%s.toml", languageStr))
	logrus.Infof("Bundle %s loaded", language)
	return bundle
}
