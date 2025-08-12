package main

// # Create the data directory if it doesn't exist
// mkdir -p data

// # Apply migrations to the shared database
// # go run cmd/migrate/main.go -db-path=./data/shared.db -migrations-dir=./sql -action=up

// # Apply migrations to a client-specific database (if needed separately)
// # go run cmd/migrate/main.go -db-path=./data/client_1.db -migrations-dir=./sql -action=up

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"path/filepath"

	"github.com/pressly/goose/v3"
	_ "modernc.org/sqlite" // SQLite driver
)

const (
	driverName = "sqlite3" // Or check goose docs for modernc
)

func main() {
	// Define flags
	dbPath := flag.String("db-path", "./data/shared.db", "Path to the shared SQLite database file")
	migrationsDir := flag.String("migrations-dir", "./sql", "Directory containing migration files")
	action := flag.String("action", "up", "Migration action: up, down, status")
	flag.Parse()

	// Open database connection
	db, err := sql.Open(driverName, *dbPath)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	// Set goose driver
	goose.SetDialect("sqlite3") // Check goose docs for correct dialect

	// Resolve absolute path for migrations
	absMigrationsDir, err := filepath.Abs(*migrationsDir)
	if err != nil {
		log.Fatalf("Failed to resolve migrations directory path: %v", err)
	}

	// Perform migration action
	switch *action {
	case "up":
		if err := goose.Up(db, absMigrationsDir); err != nil {
			log.Fatalf("Failed to apply migrations: %v", err)
		}
		fmt.Println("Migrations applied successfully.")
	case "down":
		if err := goose.Down(db, absMigrationsDir); err != nil {
			log.Fatalf("Failed to rollback migrations: %v", err)
		}
		fmt.Println("Last migration rolled back.")
	case "status":
		if err := goose.Status(db, absMigrationsDir); err != nil {
			log.Fatalf("Failed to get migration status: %v", err)
		}
	default:
		log.Fatalf("Unknown action: %s", *action)
	}
}
