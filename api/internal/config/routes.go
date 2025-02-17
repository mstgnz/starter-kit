package config

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

func LoadRoutesFromJSON() error {
	for _, lang := range App().Langs {
		fileName := fmt.Sprintf("asset/lang/%s.json", lang)
		file, err := os.Open(fileName)
		if err != nil {
			return err
		}
		defer file.Close()

		byteValue, err := io.ReadAll(file)
		if err != nil {
			return err
		}

		var data map[string]any
		if err := json.Unmarshal(byteValue, &data); err != nil {
			return err
		}

		if routes, ok := data["routes"].(map[string]any); ok {
			for key, value := range routes {
				if strVal, ok := value.(string); ok {
					if _, exists := App().Routes[key]; !exists {
						App().Routes[key] = make(map[string]string)
					}
					App().Routes[key][lang] = strVal
				}
			}
		}
	}
	return nil
}
