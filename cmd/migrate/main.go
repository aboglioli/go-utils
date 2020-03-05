package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/aboglioli/go-utils/migrations"
)

type Config struct {
	MigrationPath string `json:"path"`
	URL           string `json:"db_url"`
	User          string `json:"db_user"`
	Password      string `json:"db_password"`
}

func ReadConfig() *Config {
	config := &Config{
		MigrationPath: "migrations",
		URL:           "localhost:5432",
		User:          "",
		Password:      "",
	}

	file, err := os.Open("migrations.json")
	defer file.Close()
	if err == nil && file != nil {
		json.NewDecoder(file).Decode(config)
	}

	return config
}

func main() {
	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	config := ReadConfig()
	fmt.Printf("[Config from %s]\n%#v\n\n", dir, config)

	// Get migrations and scripts to run
	ms, err := migrations.GetMigrations(config.MigrationPath, &migrations.MigrationOptions{
		TestDB: "test",
	})
	if err != nil {
		panic(err)
	}

	// Initialize migrator
	db := migrations.NewPostgresDB()
	logger := migrations.DefaultLogger()
	migrator := migrations.NewMigrator(db, logger)

	migrator.SetMigrations(ms)

	if err := migrator.Run(config.URL, config.User, config.Password); err != nil {
		fmt.Printf("# ERROR\n%s", err)
	}
}
