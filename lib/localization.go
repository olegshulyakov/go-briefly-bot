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

// localesFolder specifies where localization files are placed
const (
	localesFolder = "locales"
	extension     = "yml"
)

// bundle is a global instance of the i18n bundle used for managing translations.
var defaultBundle atomic.Pointer[i18n.Bundle]

func init() {
	defaultBundle.Store(bundle())
}

// bundle initializes the localization bundle by loading translation files
// from the `locales` directory.
//
// The function registers JSON as the unmarshal function for translation files and
// loads all `.json` files from the `locales` directory.
//
// Example:
//
//	bundle()
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

// localizer returns a localizer for the specified language.
//
// If the language is not specified, the default language (English) is used.
//
// Parameters:
//   - lang: The language code (e.g. "en") for which to create the localizer.
//
// Returns:
//   - A pointer to an i18n.localizer instance for the specified language.
//
// Example:
//
//	localizer := localizer("ru")
//	message := localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "welcome"})
func localizer(lang string) *i18n.Localizer {
	if lang == "" {
		lang = language.English.String()
	}

	localizer := i18n.NewLocalizer(defaultBundle.Load(), lang)
	return localizer
}

func Localize(languageCode string, messageID string) (string, error) {
	return localizer(languageCode).Localize(&i18n.LocalizeConfig{MessageID: messageID})
}

func MustLocalize(languageCode string, messageID string) string {
	return localizer(languageCode).MustLocalize(&i18n.LocalizeConfig{MessageID: messageID})
}

func MustLocalizeTemplate(languageCode string, messageID string, templateData interface{}) string {
	return localizer(languageCode).MustLocalize(&i18n.LocalizeConfig{MessageID: messageID, TemplateData: templateData})
}
