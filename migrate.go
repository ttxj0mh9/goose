package goose

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"path/filepath"
	"time"
)

// MigrationRecord represents a single migration record stored in the database.
type MigrationRecord struct {
	VersionID int64
	TStamp    time.Time
	IsApplied bool
}

// Migration represents a migration file with its version, source path, and
// optional Go migration functions.
type Migration struct {
	Version    int64
	Next       int64  // next version or -1 if none
	Previous   int64  // previous version or -1 if none
	Source     string // path to .sql script or go file
	Registered bool

	// Go migration functions (nil for SQL migrations)
	UpFn   func(tx *sql.Tx) error
	DownFn func(tx *sql.Tx) error

	// Go migration functions with context support
	UpFnContext   func(ctx context.Context, tx *sql.Tx) error
	DownFnContext func(ctx context.Context, tx *sql.Tx) error

	// NoVersioning disables versioning for this migration.
	NoVersioning bool
}

// String returns a human-readable representation of a Migration.
func (m *Migration) String() string {
	return fmt.Sprintf(m.Source)
}

// Label returns the migration label, which is the base name of the source file.
func (m *Migration) Label() string {
	return filepath.Base(m.Source)
}

// Migrations is a sortable list of migrations.
type Migrations []*Migration

func (ms Migrations) Len() int      { return len(ms) }
func (ms Migrations) Swap(i, j int) { ms[i], ms[j] = ms[j], ms[i] }
func (ms Migrations) Less(i, j int) bool {
	if ms[i].Version == ms[j].Version {
		panic(fmt.Sprintf("goose: duplicate version %v detected:\n%v\n%v", ms[i].Version, ms[i].Source, ms[j].Source))
	}
	return ms[i].Version < ms[j].Version
}

// Current returns the current migration in the list, or an error if not found.
func (ms Migrations) Current(current int64) (*Migration, error) {
	for i, migration := range ms {
		if migration.Version == current {
			return ms[i], nil
		}
	}
	return nil, ErrNoCurrentVersion
}

// Last returns the last migration in the list.
func (ms Migrations) Last() (*Migration, error) {
	if len(ms) == 0 {
		return nil, ErrNoNextVersion
	}
	return ms[len(ms)-1], nil
}

// Versioned returns only versioned migrations (non-unversioned).
func (ms Migrations) Versioned() (Migrations, error) {
	versioned := Migrations{}
	for _, m := range ms {
		if !m.NoVersioning {
			versioned = append(versioned, m)
		}
	}
	return versioned, nil
}

// Unversioned returns only unversioned migrations.
func (ms Migrations) Unversioned() (Migrations, error) {
	unversioned := Migrations{}
	for _, m := range ms {
		if m.NoVersioning {
			unversioned = append(unversioned, m)
		}
	}
	return unversioned, nil
}

var (
	// ErrNoCurrentVersion is returned when no current migration version is found.
	ErrNoCurrentVersion = errors.New("goose: no current version found")
	// ErrNoNextVersion is returned when no next migration version is found.
	ErrNoNextVersion = errors.New("goose: no next version found")
)
