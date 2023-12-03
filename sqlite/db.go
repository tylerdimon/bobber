package sqlite

import (
	"context"
	"embed"
	"fmt"
	"github.com/jmoiron/sqlx"
	"io/fs"
	_ "modernc.org/sqlite"
	"os"
	"path/filepath"
	"sort"
)

//go:embed migration/*.sql
var migrationFS embed.FS

type DB struct {
	conn   *sqlx.DB
	ctx    context.Context
	cancel func()
	DSN    string
}

func NewDB(dsn string) *DB {
	db := &DB{
		DSN: dsn,
	}
	db.ctx, db.cancel = context.WithCancel(context.Background())
	return db
}

func (db *DB) Open() (err error) {
	if db.DSN == "" {
		return fmt.Errorf("dsn required")
	}

	if db.DSN != ":memory:" {
		if err := os.MkdirAll(filepath.Dir(db.DSN), 0700); err != nil {
			return err
		}
	}

	if db.conn, err = sqlx.Open("sqlite", db.DSN); err != nil {
		return err
	}

	if _, err := db.conn.Exec(`PRAGMA journal_mode = wal;`); err != nil {
		return fmt.Errorf("enable wal: %w", err)
	}

	if _, err := db.conn.Exec(`PRAGMA foreign_keys = ON;`); err != nil {
		return fmt.Errorf("foreign keys pragma: %w", err)
	}

	if err := db.migrate(); err != nil {
		return fmt.Errorf("migrate: %w", err)
	}

	return nil
}

func (db *DB) migrate() error {
	if _, err := db.conn.Exec(`CREATE TABLE IF NOT EXISTS migration (name TEXT PRIMARY KEY);`); err != nil {
		return fmt.Errorf("cannot create migration table: %w", err)
	}

	// Read migration files from our embedded file system.
	// This uses Go 1.16's 'embed' package.
	names, err := fs.Glob(migrationFS, "migration/*.sql")
	if err != nil {
		return err
	}
	sort.Strings(names)

	// Loop over all migration files and execute them in order.
	for _, name := range names {
		if err := db.migrateFile(name); err != nil {
			return fmt.Errorf("migration error: name=%q err=%w", name, err)
		}
	}
	return nil
}

// migrate runs a single migration file within a transaction. On success, the
// migration file name is saved to the "migration" table to prevent re-running.
func (db *DB) migrateFile(name string) error {
	tx, err := db.conn.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Ensure migration has not already been run.
	var n int
	if err := tx.QueryRow(`SELECT COUNT(*) FROM migration WHERE name = ?`, name).Scan(&n); err != nil {
		return err
	} else if n != 0 {
		return nil // already run migration, skip
	}

	// Read and execute migration file.
	if buf, err := fs.ReadFile(migrationFS, name); err != nil {
		return err
	} else if _, err := tx.Exec(string(buf)); err != nil {
		return err
	}

	// Insert record into migration to prevent re-running migration.
	if _, err := tx.Exec(`INSERT INTO migration (name) VALUES (?)`, name); err != nil {
		return err
	}

	return tx.Commit()
}

// Close closes the database connection.
func (db *DB) Close() error {
	// Cancel background context.
	db.cancel()

	// Close database.
	if db.conn != nil {
		return db.conn.Close()
	}
	return nil
}
