package dbconfig

import (
	"github.com/arangodb/go-driver"
)

type Database struct {
	Name string
	User string
	Pass string
}

func ptr[T any](t T) *T {
	return &t
}

func (db *Database) Create(client driver.Client) error {
	exists, err := client.DatabaseExists(nil, db.Name)
	if err != nil {
		return err
	}
	if !exists {
		_, err := client.CreateDatabase(nil, db.Name, &driver.CreateDatabaseOptions{
			Users: []driver.CreateDatabaseUserOptions{{
				UserName: db.User,
				Password: db.Pass,
				Active:   ptr(true),
			}},
			Options: driver.CreateDatabaseDefaultOptions{},
		})
		if err != nil {
			return err
		}
	}
	return nil
}
