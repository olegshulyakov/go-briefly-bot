// Package config provides functionality for localization and internationalization (i18n).
package config

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

// bundle is a global instance of the i18n bundle used for managing translations.
var bundle *i18n.Bundle

// SetupLocalizer initializes the localization bundle by loading translation files
// from the `locales` directory.
//
// The function registers JSON as the unmarshal function for translation files and
// loads all `.json` files from the `locales` directory.
//
// Example:
//
//	SetupLocalizer()
func SetupLocalizer() {
	// Create a new bundle with the default language (English)
	bundle = i18n.NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("json", json.Unmarshal) // Use json.Unmarshal for JSON files

	// Load translations from the locales directory
	localesDir := "locales"
	files, err := os.ReadDir(localesDir)
	if err != nil {
		panic(err)
	}

	for _, file := range files {
		if filepath.Ext(file.Name()) == ".json" {
			bundle.MustLoadMessageFile(filepath.Join(localesDir, file.Name()))
		}
	}
}

// GetLocalizer returns a localizer for the specified language.
//
// If the language is not specified, the default language (English) is used.
//
// Parameters:
//   - lang: The language code (e.g., "en", "ru") for which to create the localizer.
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

	localizer := i18n.NewLocalizer(bundle, lang)
	return localizer
}
