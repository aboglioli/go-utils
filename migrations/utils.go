package migrations

import (
	"encoding/json"
	"os"
)

func PopulateConfig(config interface{}, path string) {
	file, err := os.Open(path)
	if err == nil && file != nil {
		json.NewDecoder(file).Decode(config)
	}
}
