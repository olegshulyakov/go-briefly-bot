package config

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

var bundle *i18n.Bundle

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

func GetLocalizer(lang string) (*i18n.Localizer){
	if lang == "" {
		lang = language.English.String()
	}

	localizer := i18n.NewLocalizer(bundle, lang)
	return localizer
}
