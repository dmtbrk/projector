package postgres

import (
	"database/sql"
	"testing"

	_ "github.com/lib/pq"
)

func setupDB(t *testing.T) *sql.DB {
	db, err := sql.Open("postgres", "sslmode=disable user=postgres")
	if err != nil {
		t.Fatal(err)
	}

	if err := db.Ping(); err != nil {
		t.Fatal(err)
	}

	return db
}

func teardownDB(t *testing.T, db *sql.DB) {
	err := db.Close()
	if err != nil {
		t.Fatal(err)
	}
}
