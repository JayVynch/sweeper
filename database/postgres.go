package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"

	"github.com/JayVynch/sweeper/config"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file" // File source driver
	"github.com/jackc/pgx/v5/pgxpool"
)

type DB struct {
	Pool *pgxpool.Pool
	conf *config.Config
}

func New(ctx context.Context, conf config.Config) *DB {
	if conf.Database.URL == "" {
		log.Fatalf("Database URL is empties")
	}

	fmt.Printf("Database URL: %s\n", conf.Database.URL)

	dbConfig, err := pgxpool.ParseConfig(conf.Database.URL)

	if err != nil {
		log.Fatalf("Cannot parse database config %v", err)
	}

	pool, err := pgxpool.NewWithConfig(ctx, dbConfig)
	if err != nil {
		log.Fatalf("Cannot connect database %v", err)
	}

	db := &DB{Pool: pool, conf: &conf}

	db.Ping(ctx)
	return db
}

func (db *DB) Ping(ctx context.Context) {
	if err := db.Pool.Ping(ctx); err != nil {
		log.Fatalf("cannot ping database %v", err)
	}
}

func (db *DB) Open(cxt context.Context) {
	if err := db.Pool.Ping(cxt); err != nil {
		log.Fatalf("Could not ping Postgres: %v", err)
	}

	log.Println("Postgres pinged")
}

func (db *DB) Close() {
	db.Pool.Close()
}

func (db *DB) Migrate() error {
	_, b, _, _ := runtime.Caller(0)

	baseDir := filepath.Dir(b)
	migrationPath := filepath.Join(baseDir, "migrations")

	// Debug: Print paths
	fmt.Printf("Base directory: %s\n", baseDir)
	fmt.Printf("Migration directory: %s\n", migrationPath)

	// Check if migrations directory exists
	if _, err := os.Stat(migrationPath); os.IsNotExist(err) {
		return fmt.Errorf("migrations directory does not exist: %s", migrationPath)
	}

	// List migration files
	files, err := os.ReadDir(migrationPath)
	if err != nil {
		return fmt.Errorf("error reading migrations directory: %v", err)
	}

	fmt.Printf("Found %d files in migrations directory:\n", len(files))
	for _, file := range files {
		fmt.Printf("  - %s\n", file.Name())
	}

	for i, file := range files {
		fullPath := filepath.Join(migrationPath, file.Name())
		fileInfo, _ := file.Info()
		fmt.Printf("  [%d] %s (size: %d bytes)\n", i+1, fullPath, fileInfo.Size())
	}

	if len(files) == 0 {
		return fmt.Errorf("no migration files found in %s", migrationPath)
	}

	// Convert to file:// URL format
	migrationURL := fmt.Sprintf("file://%s", migrationPath)
	fmt.Printf("Migration URL: %s\n", migrationURL)
	fmt.Printf("Database URL: %s\n", db.conf.Database.URL)

	m, err := migrate.New(migrationURL, db.conf.Database.URL)

	if err != nil {
		return fmt.Errorf("error creating migrate instance: %v", err)
	}
	defer m.Close()

	// Get current version and dirty state
	version, dirty, err := m.Version()

	if err != nil && err != migrate.ErrNilVersion {
		return fmt.Errorf("error getting migration version: %v", err)
	}

	if dirty {
		log.Printf("Database is in dirty state at version %d, forcing clean state", version)
		log.Printf("Database is in dirty state at version %d, attempting recovery", version)

		// Option 1: Force to previous version and re-apply
		if version > 1 {
			log.Printf("Forcing database to previous version: %d", version-1)
			if err := m.Force(int(version - 1)); err != nil {
				return fmt.Errorf("error forcing to version %d: %v", version-1, err)
			}
		} else {
			// Option 2: If at version 1 or 0, force to clean slate
			log.Printf("Forcing database to version 0 (clean slate)")
			if err := m.Force(0); err != nil {
				return fmt.Errorf("error forcing to version 0: %v", err)
			}
		}
	}

	if err = m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("error while migration up: %v", err)
	}

	newVersion, newDirty, _ := m.Version()
	fmt.Printf("Migration completed. New version: %d, dirty: %t\n", newVersion, newDirty)

	log.Println("migration complete")

	return nil

}

func (db *DB) Drop() error {
	_, b, _, _ := runtime.Caller(0)

	baseDir := filepath.Dir(b)
	migrationPath := filepath.Join(baseDir, "migrations")

	// Debug: Print paths
	fmt.Printf("Base directory: %s\n", baseDir)
	fmt.Printf("Migration directory: %s\n", migrationPath)

	// Check if migrations directory exists
	if _, err := os.Stat(migrationPath); os.IsNotExist(err) {
		return fmt.Errorf("migrations directory does not exist: %s", migrationPath)
	}

	// List migration files
	files, err := os.ReadDir(migrationPath)
	if err != nil {
		return fmt.Errorf("error reading migrations directory: %v", err)
	}

	fmt.Printf("Found %d files in migrations directory:\n", len(files))
	for _, file := range files {
		fmt.Printf("  - %s\n", file.Name())
	}

	if len(files) == 0 {
		return fmt.Errorf("no migration files found in %s", migrationPath)
	}

	// Convert to file:// URL format
	migrationURL := fmt.Sprintf("file://%s", migrationPath)
	fmt.Printf("Migration URL: %s\n", migrationURL)
	fmt.Printf("Database URL: %s\n", db.conf.Database.URL)

	m, err := migrate.New(migrationURL, db.conf.Database.URL)
	if err != nil {
		return fmt.Errorf("error creating migrate instance: %v", err)
	}
	defer m.Close()

	// Get current version and dirty state
	version, dirty, err := m.Version()
	if err != nil && err != migrate.ErrNilVersion {
		return fmt.Errorf("error getting migration version: %v", err)
	}

	if dirty {
		log.Printf("Database is in dirty state at version %d, forcing clean state", version)
		if err := m.Force(int(version)); err != nil {
			return fmt.Errorf("error forcing version %d: %v", version, err)
		}
	}

	if err = m.Drop(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("error while dropping: %v", err)
	}

	log.Println("migration Drop")

	return nil

}

func (db *DB) Truncate(ctx context.Context) error {
	if _, err := db.Pool.Exec(ctx,
		`DELETE FROM users;
	`); err != nil {
		return fmt.Errorf("error truncating: %v", err)
	}

	return nil
}
