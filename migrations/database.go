package migrations

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type Database interface {
	Connect(url, database, user, passwd string) error
	Close() error

	CreateMigrationTable() error
	AddScript(script string) error
	GetScripts() (map[string]bool, error)

	Exec(sql string) error
}

type postgresDB struct {
	db *sql.DB
}

func NewPostgresDB() *postgresDB {
	return &postgresDB{}
}

func (postgres *postgresDB) Connect(url, database, user, passwd string) error {
	connStr := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable", user, passwd, url, database)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return err
	}

	if err := db.Ping(); err != nil {
		return err
	}

	postgres.db = db
	return nil
}

func (postgres *postgresDB) Close() error {
	return postgres.db.Close()
}

func (postgres *postgresDB) CreateMigrationTable() error {
	_, err := postgres.db.Exec(`
		CREATE TABLE IF NOT EXISTS migrations (
			script TEXT NOT NULL
		)
	`)
	return err
}

func (postgres *postgresDB) AddScript(script string) error {
	_, err := postgres.db.Exec(`INSERT INTO migrations(script) VALUES($1)`, script)
	return err
}

func (postgres *postgresDB) GetScripts() (map[string]bool, error) {
	if err := postgres.CreateMigrationTable(); err != nil {
		return nil, err
	}

	rows, err := postgres.db.Query("SELECT script FROM migrations")
	if err != nil {
		return nil, err
	}

	scripts := make(map[string]bool)
	for rows.Next() {
		var s string
		err := rows.Scan(&s)
		if err != nil {
			return nil, err
		}
		scripts[s] = true
	}

	return scripts, nil
}

func (postgres *postgresDB) Exec(sql string) error {
	_, err := postgres.db.Exec(sql)
	return err
}
