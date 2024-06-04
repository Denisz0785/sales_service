package databasetest

import (
	"sales_service/internal/platform/database"
	"sales_service/internal/schema"
	"testing"
	"time"

	"github.com/jmoiron/sqlx"
)

func Setup(t *testing.T) (*sqlx.DB, func()) {
	t.Helper()

	c := startContainer(t)

	db, err := database.OpenDB(database.Config{
		Host:       c.Host,
		Name:       "postgres",
		User:       "postgres",
		Password:   "postgres",
		DisableTLS: "disable",
	})
	if err != nil {
		t.Fatalf("failed to open database: %v", err)
	}
	t.Log("waiting for database to be ready...")

	var pingError error
	maxAttempts := 20
	for attempts := 1; attempts <= maxAttempts; attempts++ {
		pingError = db.Ping()
		if pingError == nil {
			break
		}

		time.Sleep(time.Duration(attempts) * 100 * time.Millisecond)
	}
	if pingError != nil {
		stopContainer(t, c)
		t.Fatalf("database is not ready after %d attempts: %v", maxAttempts, pingError)
	}

	if err := schema.Migrate(db); err != nil {
		stopContainer(t, c)
		t.Fatalf("failed to apply migrations: %v", err)
	}

	teardown := func() {
		t.Helper()
		db.Close()
		stopContainer(t, c)
	}

	return db, teardown
}
