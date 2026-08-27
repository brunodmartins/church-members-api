package i18n

import (
	"testing"

	goi18n "github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

func TestLoadBundleFallsBackToParentLanguage(t *testing.T) {
	localize := GetLocalize(language.MustParse("en-US"))

	message, err := localize.Localize(&goi18n.LocalizeConfig{
		MessageID: "Domain.Name",
	})
	if err != nil {
		t.Fatalf("localizing parent bundle message: %v", err)
	}

	if message != "Name" {
		t.Fatalf("expected parent bundle message %q, got %q", "Name", message)
	}
}

func TestLoadBundleFallsBackToEnglish(t *testing.T) {
	localize := GetLocalize(language.MustParse("fr-CA"))

	message, err := localize.Localize(&goi18n.LocalizeConfig{
		MessageID: "Domain.Name",
	})
	if err != nil {
		t.Fatalf("localizing English fallback message: %v", err)
	}

	if message != "Name" {
		t.Fatalf("expected English fallback message %q, got %q", "Name", message)
	}
}