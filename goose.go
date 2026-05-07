// Package goose provides database migration tooling.
// It is a fork of pressly/goose with additional features and improvements.
package goose

import (
	"database/sql"
	"fmt"
)

// DefaultTableName is the default name for the goose migrations table.
const DefaultTableName = "goose_db_version"

// TableName returns the name of the migrations table.
func TableName() string {
	return tableName
}

// SetTableName sets a custom name for the goose migrations table.
func SetTableName(name string) {
	tableName = name
}

var tableName = DefaultTableName

// Status prints the migration status for the current DB.
func Status(db *sql.DB, dir string) error {
	return status(db, dir)
}

// Up applies all available migrations.
func Up(db *sql.DB, dir string) error {
	return up(db, dir, 0)
}

// UpByOne migrates up by a single version.
func UpByOne(db *sql.DB, dir string) error {
	return upByOne(db, dir)
}

// UpTo migrates up to a specific version.
func UpTo(db *sql.DB, dir string, version int64) error {
	return up(db, dir, version)
}

// Down rolls back a single migration from the current version.
func Down(db *sql.DB, dir string) error {
	return down(db, dir)
}

// DownTo rolls back migrations to a specific version.
func DownTo(db *sql.DB, dir string, version int64) error {
	return downTo(db, dir, version)
}

// Redo re-runs the latest migration.
func Redo(db *sql.DB, dir string) error {
	return redo(db, dir)
}

// Reset rolls back all migrations.
func Reset(db *sql.DB, dir string) error {
	return reset(db, dir)
}

// Version prints the current version of the database.
func Version(db *sql.DB, dir string) error {
	current, err := GetDBVersion(db)
	if err != nil {
		return err
	}
	fmt.Printf("goose: version %d\n", current)
	return nil
}

// GetDBVersion returns the current version of the database.
func GetDBVersion(db *sql.DB) (int64, error) {
	return getDBVersion(db)
}

// Create writes a new blank migration file.
func Create(db *sql.DB, dir, name, migrationType string) error {
	return create(db, dir, name, migrationType)
}

// Fix rewrites migration filenames to use the timestamp format.
func Fix(dir string) error {
	return fix(dir)
}
