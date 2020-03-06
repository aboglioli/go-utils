package migrations

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const (
	SCRIPTS_PATTERN = "/*/*.sql"
	CONFIG_FILE     = "db.json"
)

type Config struct {
	URL      string `json:"url"`
	User     string `json:"user"`
	Password string `json:"password"`
	Database string `json:"database"`
}

type Migration struct {
	Config  Config
	Path    string
	Dir     string
	Scripts []string
}

func GetMigrations(path string) ([]*Migration, error) {
	glob := filepath.Join(path, SCRIPTS_PATTERN)
	files, err := filepath.Glob(glob)
	if err != nil {
		return nil, err
	}

	dirs := make(map[string][]string)
	for _, file := range files {
		parts := strings.Split(file, "/")
		dir := parts[len(parts)-2]
		file := parts[len(parts)-1]
		dirs[dir] = append(dirs[dir], file)
	}

	ms := make([]*Migration, 0)
	for dir, files := range dirs {
		// Get configuration file
		configPath := filepath.Join(path, dir, CONFIG_FILE)
		config := Config{
			URL:      "localhost:5432",
			User:     "",
			Password: "",
			Database: "database",
		}

		file, err := os.Open(configPath)
		fmt.Println(configPath)
		if err == nil && file != nil {
			if err := json.NewDecoder(file).Decode(&config); err != nil {
				return nil, err
			}
		}

		m := &Migration{
			Config:  config,
			Path:    path,
			Dir:     dir,
			Scripts: files,
		}

		ms = append(ms, m)
	}

	return ms, nil
}
