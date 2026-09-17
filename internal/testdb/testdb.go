package testdb

import (
	"database/sql"
	"fmt"
	"os"
	"sync"
	"testing"

	"github.com/geneowak/go-expense-tracker/internal/database"
	migrations "github.com/geneowak/go-expense-tracker/sql"
	"github.com/joho/godotenv"

	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
)

var (
	openDb      *sql.DB
	migrateOnce sync.Once
	migrateErr  error
)

/**
* Connects to TEST_DB_URL and runs goose migrations once per test binary.
* Returns (nil, nil) when TEST_DB_URL is unset so callers can skip real-db tests gracefully
 */
func OpenAndMigrate(t *testing.T) *sql.DB {
	t.Helper()

	err := godotenv.Load("../.env.testing")
	if err != nil {
		t.Fatalf("Error loading .env file: %v", err)
	}

	dsn := os.Getenv("TEST_DB_URL")
	if dsn == "" {
		return nil
	}

	migrateOnce.Do(func() {
		var err error
		openDb, err = sql.Open("postgres", dsn)
		if err != nil {
			migrateErr = fmt.Errorf("Error opening test db: %w", err)
			return
		}

		goose.SetBaseFS(migrations.EmbedMigrations)
		if err := goose.SetDialect("postgres"); err != nil {
			migrateErr = fmt.Errorf("Error setting dialect: %w", err)
			return
		}
		if err := goose.Up(openDb, "schema"); err != nil {
			migrateErr = fmt.Errorf("Error migrating db: %w", err)
			return
		}
	})
	if migrateErr != nil {
		t.Fatalf("Test db setup failed: %v", migrateErr)
	}

	if err := openDb.Ping(); err != nil {
		t.Fatalf("unable to reach database: %v", err)
	}

	return openDb
}

// bundles a transaction with a Queries instance bound to it
type Scope struct {
	Tx      *sql.Tx
	Queries *database.Queries
}

/**
* Starts a new transaction and returns a Scope
* Transaction is rolled back automatically when the test finishes
* So each test case runs in isolation on a clean slate, just like Laravel's Refresh Database
 */
func Begin(t *testing.T, db *sql.DB) Scope {
	t.Helper()

	tx, err := db.Begin()
	if err != nil {
		t.Fatalf("Failed to begin transaction: %v", err)
	}

	t.Cleanup(func() {
		tx.Rollback()
	})

	return Scope{
		Tx:      tx,
		Queries: database.New(tx),
	}
}
