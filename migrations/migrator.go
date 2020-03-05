package migrations

import "io/ioutil"

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
		if err := m.db.Connect("localhost:5432", migration.Database, "admin", "admin"); err != nil {
			return err
		}

		m.logger.Log("# Migrating %s:\n", migration.Database)

		oldScripts, err := m.db.GetScripts()
		if err != nil {
			return err
		}

		for _, script := range migration.Scripts {
			if _, ok := oldScripts[script]; ok {
				m.logger.Log("- %s\n", script)
				continue
			}

			m.logger.Log("+ %s: ", script)
			b, err := ioutil.ReadFile(script)
			if err != nil {
				return err
			}

			sql := string(b)

			err = m.db.Exec(sql)
			if err != nil {
				return err
			}

			m.db.AddScript(script)

			m.logger.Log("OK\n")
		}
		m.logger.Log("\n")

		if err := m.db.Close(); err != nil {
			return err
		}
	}

	return nil
}
