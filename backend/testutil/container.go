// Package testutil provides helpers for integration tests that require a real
// Postgres instance. Each test isolates its writes via transaction rollback.
package testutil

import (
	"context"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	tcpostgres "github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"

	interndb "github.com/todoist/backend/internal/db"
)

// NewPostgresPool starts a throwaway Postgres 17 container, runs all
// migrations, and returns a ready pool. The container is terminated via
// t.Cleanup so callers don't need to manage lifecycle manually.
func NewPostgresPool(t *testing.T) *pgxpool.Pool {
	t.Helper()
	ctx := context.Background()

	container, err := tcpostgres.Run(ctx,
		"postgres:17-alpine",
		tcpostgres.WithDatabase("testdb"),
		tcpostgres.WithUsername("test"),
		tcpostgres.WithPassword("test"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(60*time.Second),
		),
	)
	require.NoError(t, err, "start postgres container")
	t.Cleanup(func() { _ = container.Terminate(ctx) })

	dsn, err := container.ConnectionString(ctx, "sslmode=disable")
	require.NoError(t, err, "get connection string")

	pool, err := pgxpool.New(ctx, dsn)
	require.NoError(t, err, "create pool")
	t.Cleanup(pool.Close)

	applyMigrations(t, pool)
	return pool
}

// applyMigrations reads and executes each up-migration file in sequence.
// Using raw SQL execution keeps testutil independent of the migrate CLI binary.
func applyMigrations(t *testing.T, pool *pgxpool.Pool) {
	t.Helper()
	migrationsDir := migrationsPath()
	ctx := context.Background()

	files := []string{
		"00001_create_users.sql",
		"00002_create_refresh_tokens.sql",
		"00003_create_projects.sql",
		"00004_create_labels.sql",
		"00005_create_tasks.sql",
		"00006_create_task_labels.sql",
		"00007_create_comments.sql",
		"00008_harden_multi_tenant_constraints.sql",
	}

	for _, name := range files {
		path := filepath.Join(migrationsDir, name)
		sql, err := os.ReadFile(path)
		require.NoError(t, err, "read migration %s", name)
		_, err = pool.Exec(ctx, extractUpMigration(string(sql)))
		require.NoError(t, err, "apply migration %s", name)
	}
}

func extractUpMigration(contents string) string {
	upSection, _, _ := strings.Cut(contents, "-- +goose Down")
	upSection = strings.ReplaceAll(upSection, "-- +goose Up", "")
	return strings.TrimSpace(upSection)
}

// migrationsPath resolves the absolute path of the migrations/ directory
// relative to this source file, so it works regardless of the working
// directory the test binary was invoked from.
func migrationsPath() string {
	_, file, _, _ := runtime.Caller(0)
	// file is .../backend/testutil/container.go
	// migrations/ is at .../backend/migrations/
	return filepath.Join(filepath.Dir(file), "..", "migrations")
}

// TxContext starts a transaction, injects it into a new context via
// db.InjectTx, and registers a rollback in t.Cleanup. The returned context
// is passed to repository methods so they execute inside the transaction.
//
// All writes made through the returned context are rolled back when the test
// finishes — no need to truncate tables between tests.
func TxContext(t *testing.T, pool *pgxpool.Pool) context.Context {
	t.Helper()
	ctx := context.Background()

	tx, err := pool.Begin(ctx)
	require.NoError(t, err, "begin test transaction")
	t.Cleanup(func() { _ = tx.Rollback(ctx) })

	return interndb.InjectTx(ctx, pgx.Tx(tx))
}
