package postgresql

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

var migrations = []struct {
	name string
	stmt string
}{
	{
		name: "enable-UUID-extension",
		stmt: enableUUIDExtension,
	},
	{
		name: "create-users-table",
		stmt: createTableUsers,
	},
	{
		name: "create-key_values-table",
		stmt: createTableKeyValues,
	},
}

// Migrate performs the database migration. If the migration fails
// and error is returned.
func Migrate(db *sql.DB) error {
	if err := createMigrationHistoryTable(db); err != nil {
		return err
	}
	completed, err := selectCompletedMigrations(db)
	if err != nil && err != sql.ErrNoRows {
		return err
	}
	for _, migration := range migrations {
		log.Print(migration.name)
		if _, ok := completed[migration.name]; ok {
			log.Println(" skipped")
			continue
		}
		log.Println(" executing")
		if _, err := db.Exec(migration.stmt); err != nil {
			return err
		}
		if err := addMigration(db, migration.name); err != nil {
			return err
		}

	}
	return nil
}

func createMigrationHistoryTable(db *sql.DB) error {
	_, err := db.Exec(migrationTableCreate)
	return err
}

func addMigration(db *sql.DB, name string) error {
	_, err := db.Exec(migrationInsert, name)
	return err
}

func selectCompletedMigrations(db *sql.DB) (map[string]struct{}, error) {
	migrations := map[string]struct{}{}
	rows, err := db.Query(migrationSelect)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, err
		}
		migrations[name] = struct{}{}
	}
	return migrations, nil
}

//
// migration table ddl and sql
//

var migrationTableCreate = `
CREATE TABLE IF NOT EXISTS migration_history (
name VARCHAR(255),
UNIQUE(name)
)
`

var migrationInsert = `
INSERT INTO migration_history (name) VALUES ($1)
`

var migrationSelect = `
SELECT name FROM migration_history
`

var enableUUIDExtension = `
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
`

var createTableUsers = `
CREATE TABLE IF NOT EXISTS users (
	user_name VARCHAR(255) NOT NULL,
	password  VARCHAR(255) NOT NULL,
	urls	  JSONB,
	UNIQUE(user_name)
)
`

var createTableKeyValues = `
CREATE TABLE IF NOT EXISTS key_values (
	user_name VARCHAR(255) NOT NULL,
	key   VARCHAR(255) NOT NULL,
	value   VARCHAR(255) NOT NULL,
	UNIQUE(user_name, key)
)
`
