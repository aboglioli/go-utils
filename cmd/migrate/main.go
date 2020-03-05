package main

import (
	"fmt"

	"github.com/aboglioli/go-utils/migrations"
)

const (
	MIGRATIONS_PATH = "./migrations"
)

func main() {
	// Get migrations and scripts to run
	ms, err := migrations.GetMigrations(MIGRATIONS_PATH, &migrations.MigrationOptions{
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
