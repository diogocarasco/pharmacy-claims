package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

// ApplyMigrations applies the database migrations.
func ApplyMigrations(db *sql.DB) error {
	log.Println("Applying database migrations...")
	schema := `
	CREATE TABLE IF NOT EXISTS pharmacies (
		chain TEXT NOT NULL,
		npi TEXT PRIMARY KEY UNIQUE
	);
	CREATE TABLE IF NOT EXISTS claims (
		id TEXT PRIMARY KEY,
		ndc TEXT NOT NULL,
		npi TEXT NOT NULL,
		quantity REAL NOT NULL,
		price REAL NOT NULL,
		timestamp TEXT NOT NULL,
		reverted BOOLEAN NOT NULL DEFAULT FALSE
	);
	CREATE TABLE IF NOT EXISTS reverts (
		id TEXT PRIMARY KEY,
		claim_id TEXT NOT NULL,
		timestamp TEXT NOT NULL,
		FOREIGN KEY (claim_id) REFERENCES claims(id)
	);
	`
	_, err := db.Exec(schema)
	if err != nil {
		return fmt.Errorf("error applying migrations: %w", err)
	}
	log.Println("Migrations applied successfully.")
	return nil
}
