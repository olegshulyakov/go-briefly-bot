package briefly

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

// bundle is a global instance of the i18n bundle used for managing translations.
var bundle *i18n.Bundle

// getBundle initializes the localization bundle by loading translation files
// from the `locales` directory.
//
// The function registers JSON as the unmarshal function for translation files and
// loads all `.json` files from the `locales` directory.
//
// Example:
//
//	getBundle()
func getBundle(localesDir string) *i18n.Bundle {
	if localesDir == "" {
		localesDir = "locales"
	}

	// Create a new bundle with the default language (English)
	bundle = i18n.NewBundle(language.English)

	bundle.RegisterUnmarshalFunc("json", json.Unmarshal)

	// Load translations from the locales directory
	files, err := os.ReadDir(localesDir)
	if err != nil {
		Error("Failed to read locales directory", "error", err)
	}

	for _, file := range files {
		if filepath.Ext(file.Name()) == ".json" {
			bundle.MustLoadMessageFile(filepath.Join(localesDir, file.Name()))
		}
	}

	return bundle
}

// GetLocalizer returns a localizer for the specified language.
//
// If the language is not specified, the default language (English) is used.
//
// Parameters:
//   - lang: The language code (e.g. "en") for which to create the localizer.
//
// Returns:
//   - A pointer to an i18n.Localizer instance for the specified language.
//
// Example:
//
//	localizer := GetLocalizer("ru")
//	message := localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "welcome"})
func GetLocalizer(lang string) *i18n.Localizer {
	if lang == "" {
		lang = language.English.String()
	}
	if bundle == nil {
		bundle = getBundle("")
	}

	localizer := i18n.NewLocalizer(bundle, lang)
	return localizer
}
