package i18n

import (
	"fmt"
	"os"

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
	logrus.Infof("Loading bundle for language %s", language)
	base := os.Getenv("GOPATH")
	appPath := "github.com/BrunoDM2943/church-members-api"
	bundle.MustLoadMessageFile(fmt.Sprintf("%s/src/%s/bundles/%s.toml", base, appPath, language))
	logrus.Infof("Bundle %s loaded", language)
	return bundle
}
