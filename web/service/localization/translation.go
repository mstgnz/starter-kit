package localization

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/mstgnz/starter-kit/web/config"
)

type Translation map[string]any

var Translations = map[string]Translation{}

func LoadTranslations() {
	for _, lang := range config.App().Langs {
		fileName := fmt.Sprintf("asset/lang/%s.json", lang)
		file, err := os.Open(fileName)
		if err != nil {
			log.Fatalf("Failed to open file: %v", err)
		}
		defer file.Close()

		bytes, err := io.ReadAll(file)
		if err != nil {
			log.Fatalf("Failed to read file: %v", err)
		}

		var translation Translation
		if err := json.Unmarshal(bytes, &translation); err != nil {
			log.Fatalf("Failed to unmarshal JSON: %v", err)
		}

		Translations[lang] = translation
	}
}

func GetLang(key string) string {
	translation, ok := Translations[config.App().Lang]
	if !ok {
		return ""
	}

	var result string
	var recursive func(key string, translation Translation)

	recursive = func(key string, translation Translation) {
		keys := strings.Split(key, ".")
		if len(keys) == 1 {
			if val, ok := translation[keys[0]].(string); ok {
				result = val
			}
			return
		}

		if nextTranslation, ok := translation[keys[0]].(map[string]any); ok {
			recursive(keys[1], nextTranslation)
		}
	}

	recursive(key, translation)

	return result
}
