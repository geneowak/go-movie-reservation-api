package testdb

import (
	"context"
	"fmt"
	"os"
	"sync"
	"testing"

	"github.com/geneowak/cinehold/internal/database"
	migrations "github.com/geneowak/cinehold/sql"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"

	"github.com/pressly/goose/v3"
)

var (
	openConn    *pgx.Conn
	migrateOnce sync.Once
	migrateErr  error
)

/**
* Connects to TEST_DB_URL and runs goose migrations once per test binary.
* Returns (nil, nil) when TEST_DB_URL is unset so callers can skip real-db tests gracefully
 */
func OpenAndMigrate(t *testing.T) *pgx.Conn {
	t.Helper()
	ctx := context.Background()

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
		openConn, err = pgx.Connect(ctx, dsn)
		if err != nil {
			migrateErr = fmt.Errorf("Error opening test db: %w", err)
			return
		}

		connConfig := openConn.Config()

		gooseDb := stdlib.OpenDB(*connConfig)

		goose.SetBaseFS(migrations.EmbedMigrations)
		if err := goose.SetDialect("postgres"); err != nil {
			migrateErr = fmt.Errorf("Error setting dialect: %w", err)
			return
		}
		if err := goose.Up(gooseDb, "schema"); err != nil {
			migrateErr = fmt.Errorf("Error migrating db: %w", err)
			return
		}
	})
	if migrateErr != nil {
		t.Fatalf("Test db setup failed: %v", migrateErr)
	}

	if err := openConn.Ping(context.Background()); err != nil {
		t.Fatalf("unable to reach database: %v", err)
	}

	return openConn
}

// bundles a transaction with a Queries instance bound to it
type Scope struct {
	Tx      pgx.Tx
	Queries *database.Queries
}

/**
* Starts a new transaction and returns a Scope
* Transaction is rolled back automatically when the test finishes
* So each test case runs in isolation on a clean slate, just like Laravel's Refresh Database
 */
func Begin(t *testing.T, conn *pgx.Conn) Scope {
	t.Helper()

	tx, err := conn.Begin(context.Background())
	if err != nil {
		t.Fatalf("Failed to begin transaction: %v", err)
	}

	t.Cleanup(func() {
		tx.Rollback(context.Background())
	})

	return Scope{
		Tx:      tx,
		Queries: database.New(tx),
	}
}
