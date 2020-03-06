package migrations

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetMigrations(t *testing.T) {
	assert := assert.New(t)
	ms, err := GetMigrations("../migrations_example")
	assert.Nil(err)

	expected := []*Migration{
		&Migration{
			Config: Config{ // default
				URL:      "localhost:5432",
				User:     "",
				Password: "",
				Database: "database",
			},
			Path:    "../migrations_example",
			Dir:     "products",
			Scripts: []string{"0001_tables.sql", "0002_populate.sql", "0003_add_column.sql"},
		},
		&Migration{
			Config: Config{ // from db.json
				URL:      "remote.com:5432",
				User:     "admin",
				Password: "passwd",
				Database: "users",
			},
			Path:    "../migrations_example",
			Dir:     "users",
			Scripts: []string{"0001_tables.sql", "0002_populate.sql"},
		},
	}
	// if !reflect.DeepEqual(expected, ms) {
	// 	t.Errorf("Error:\n-expected:%#v\n-actual  :%#v", expected, ms)
	// }
	assert.Equal(expected, ms)
}
