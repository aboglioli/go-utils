package migrations

import (
	"path/filepath"
	"strings"
)

const (
	SCRIPTS_PATTERN = "/*/*.sql"
)

type MigrationOptions struct {
	TestDB string
}

type Migration struct {
	Database string
	Scripts  []string
}

func GetMigrations(path string, opts *MigrationOptions) ([]*Migration, error) {
	glob := filepath.Join(path, SCRIPTS_PATTERN)
	files, err := filepath.Glob(glob)
	if err != nil {
		return nil, err
	}

	m := make(map[string][]string)
	for _, file := range files {
		parts := strings.Split(file, "/")
		dbName := parts[1]
		m[dbName] = append(m[dbName], file)
	}

	migrations := make([]*Migration, 0)
	test := &Migration{Database: opts.TestDB}
	for dbName, scripts := range m {
		migrations = append(migrations, &Migration{
			Database: dbName,
			Scripts:  scripts,
		})
		test.Scripts = append(test.Scripts, scripts...)
	}

	if opts.TestDB != "" {
		migrations = append(migrations, test)
	}

	return migrations, nil
}
