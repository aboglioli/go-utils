package migrations

import (
	"io/ioutil"
	"path/filepath"
)

type Migrator struct {
	db         Database
	logger     Logger
	migrations []*Migration
}

func NewMigrator(db Database, logger Logger) *Migrator {
	return &Migrator{
		db:         db,
		logger:     logger,
		migrations: make([]*Migration, 0),
	}
}

func (m *Migrator) SetMigrations(migrations []*Migration) {
	m.migrations = migrations
}

func (m *Migrator) Run() error {
	for _, migration := range m.migrations {
		if err := m.db.Connect(migration.Config.URL, migration.Config.Database, migration.Config.User, migration.Config.Password); err != nil {
			return err
		}

		m.logger.Log("# Migrating %s in %s:\n", migration.Config.Database, migration.Config.URL)

		oldScripts, err := m.db.GetScripts()
		if err != nil {
			return err
		}

		for _, script := range migration.Scripts {
			scriptPath := filepath.Join(migration.Path, migration.Dir, script)
			if _, ok := oldScripts[scriptPath]; ok {
				m.logger.Log("- %s\n", scriptPath)
				continue
			}

			m.logger.Log("+ %s: ", scriptPath)
			b, err := ioutil.ReadFile(scriptPath)
			if err != nil {
				return err
			}

			sql := string(b)
			err = m.db.Exec(sql)
			if err != nil {
				return err
			}

			m.db.AddScript(scriptPath)

			m.logger.Log("OK\n")
		}
		m.logger.Log("\n")

		if err := m.db.Close(); err != nil {
			return err
		}
	}

	return nil
}
