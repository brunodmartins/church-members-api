package i18n

import (
	"embed"
	"fmt"
	"github.com/spf13/viper"
	"sync"

	"github.com/BurntSushi/toml"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/sirupsen/logrus"
	"golang.org/x/text/language"
)

//go:embed languages/*.toml
var LocaleFS embed.FS

type MessageService struct {
	localize *i18n.Localizer
}

var bundles = make(map[string]*i18n.Localizer)

func (service *MessageService) GetMessage(key, defaultValue string) string {
	return service.localize.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    key,
			Other: defaultValue,
		},
	})
}

var (
	service     *MessageService
	serviceOnce sync.Once
)

// GetMessageService builds a singleton instance for MessageService
func GetMessageService() *MessageService {
	serviceOnce.Do(func() {
		buildMessageService()
	})
	return service
}

func buildMessageService() {
	lang := language.English
	if envLang := viper.GetString("lang"); envLang != "" {
		lang = language.MustParse(envLang)
	}
	bundle := loadBundle(lang)
	service = &MessageService{
		localize: i18n.NewLocalizer(bundle),
	}
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
