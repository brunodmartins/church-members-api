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

func GetMessage(ctx context.Context, key string, value ...any) string {
	localize := GetLocalize(language.English)
	if ctx.Value("i18n") != nil {
		localize = ctx.Value("i18n").(*i18n.Localizer)
	}
	return fmt.Sprintf(localize.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID: key,
		},
	}), value...)
}

func loadBundle(tag language.Tag) *i18n.Bundle {
	for _, candidate := range languageCandidates(tag) {
		bundle := i18n.NewBundle(candidate)
		bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)
		path := fmt.Sprintf("languages/%s.toml", candidate.String())
		logrus.Infof("Loading bundle for language %s from %s", candidate, path)
		if _, err := bundle.LoadMessageFileFS(LocaleFS, path); err == nil {
			logrus.Infof("Bundle for language %s loaded from %s", tag, path)
			return bundle
		} else {
			logrus.Warnf("Unable to load bundle for language %s from %s: %v", candidate, path, err)
		}
	}

	logrus.Warnf("No bundle found for language %s; using English bundle", tag)
	bundle := i18n.NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)
	return bundle
}

func languageCandidates(tag language.Tag) []language.Tag {
	candidates := []language.Tag{tag}
	for parent := tag.Parent(); parent != language.Und && parent != tag; parent = parent.Parent() {
		candidates = append(candidates, parent)
	}
	if tag != language.English {
		candidates = append(candidates, language.English)
	}
	return candidates
}

// GetLocalize returns an i18n.Localize based on a language tag
func GetLocalize(tag language.Tag) *i18n.Localizer {
	if bundles[tag.String()] != nil {
		return bundles[tag.String()]
	}
	languages := languageCandidates(tag)
	languageStrings := make([]string, len(languages))
	for index, candidate := range languages {
		languageStrings[index] = candidate.String()
	}
	result := i18n.NewLocalizer(loadBundle(tag), languageStrings...)
	bundles[tag.String()] = result
	return result
}
