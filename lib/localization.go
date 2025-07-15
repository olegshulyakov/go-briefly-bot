package lib

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"sync/atomic"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
	"gopkg.in/yaml.v3"
)

const (
	// localesFolder is the name of the directory containing translation files.
	localesFolder = "locales"
	// extension is the file extension for translation files.
	extension = "yml"
)

// defaultBundle holds the atomic pointer to the i18n bundle.
var defaultBundle atomic.Pointer[i18n.Bundle]

func init() {
	defaultBundle.Store(bundle())
}

// bundle initializes and returns a new i18n bundle with translations loaded
// from the locales directory. The bundle uses English as default language
// and YAML as the translation file format.
func bundle() *i18n.Bundle {
	// Create a new bundle with the default language (English)
	bundle := i18n.NewBundle(language.English)

	bundle.RegisterUnmarshalFunc(extension, yaml.Unmarshal)

	// Load translations from the locales directory
	files, err := os.ReadDir(localesFolder)
	if err != nil {
		slog.Error("Failed to read locales directory", "error", err)
	}

	for _, file := range files {
		if filepath.Ext(file.Name()) == fmt.Sprintf(".%s", extension) {
			bundle.MustLoadMessageFile(filepath.Join(localesFolder, file.Name()))
		}
	}

	return bundle
}

// localizer creates and returns a new Localizer for the specified language.
// If language is empty, defaults to English.
func localizer(lang string) *i18n.Localizer {
	if lang == "" {
		lang = language.English.String()
	}

	localizer := i18n.NewLocalizer(defaultBundle.Load(), lang)
	return localizer
}

// Localize translates a message ID to the specified language.
// Returns the translated string or an error if translation fails.
// If languageCode is empty, defaults to English.
func Localize(languageCode string, messageID string) (string, error) {
	return localizer(languageCode).Localize(&i18n.LocalizeConfig{MessageID: messageID})
}

// MustLocalize translates a message ID to the specified language.
// Panics if translation fails. If languageCode is empty, defaults to English.
func MustLocalize(languageCode string, messageID string) string {
	return localizer(languageCode).MustLocalize(&i18n.LocalizeConfig{MessageID: messageID})
}

// MustLocalizeTemplate translates a message ID with template data to the specified language.
// Panics if translation fails. If languageCode is empty, defaults to English.
// The templateData parameter is used to interpolate values in the translated message.
func MustLocalizeTemplate(languageCode string, messageID string, templateData interface{}) string {
	return localizer(languageCode).MustLocalize(&i18n.LocalizeConfig{
		MessageID:    messageID,
		TemplateData: templateData,
	})
}
