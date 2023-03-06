package storagetest

import (
	"database/sql"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func makeSqliteTestDB(t *testing.T) *sql.DB {
	t.Helper()

	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("unable to open database: %v", err)
	}

	return db
}
