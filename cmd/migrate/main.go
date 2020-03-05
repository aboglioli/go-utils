package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/aboglioli/go-utils/migrations"
)

type Config struct {
	MigrationPath string `json:"path"`
}

func ReadConfig() *Config {
	config := &Config{
		MigrationPath: "migrations",
	}

	file, err := os.Open("migrations.json")
	defer file.Close()
	if err == nil && file != nil {
		json.NewDecoder(file).Decode(config)
	}

	return config
}

func main() {
	config := ReadConfig()
	fmt.Printf("[Config]\n- path: %s\n\n", config.MigrationPath)

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

	if err := migrator.Run(); err != nil {
		fmt.Printf("# ERROR\n%s", err)
	}
}
