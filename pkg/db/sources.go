package db

import (
	"database/sql"
	"fmt"
	"time"
)

// InsertSource inserts a new source into the Sources table.
func InsertSource(manager *DBManager, source Source) error {
	db, err := manager.GetPrimaryDB()
	if err != nil {
		return fmt.Errorf("failed to get shared DB: %w", err)
	}

	query := `
		INSERT OR REPLACE INTO Sources (Url, Language, Title, Text, CreatedAt)
		VALUES (?, ?, ?, ?, ?)
	`

	createdAt := source.CreatedAt
	if createdAt == "" {
		createdAt = time.Now().Format(time.RFC3339)
	}

	_, err = db.Exec(query, source.Url, source.Language, source.Title, source.Text, createdAt)
	if err != nil {
		return fmt.Errorf("failed to insert into Sources: %w", err)
	}
	return nil
}

// GetSource retrieves a source by URL and language.
func GetSource(manager *DBManager, url string, language string) (*Source, error) {
	db, err := manager.GetPrimaryDB()
	if err != nil {
		return nil, fmt.Errorf("failed to get shared DB: %w", err)
	}

	query := `
		SELECT Url, Language, Title, Text, CreatedAt
		FROM Sources
		WHERE Url = ? AND Language = ?
	`

	row := db.QueryRow(query, url, language)
	var source Source
	err = row.Scan(&source.Url, &source.Language, &source.Title, &source.Text, &source.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Not found
		}
		return nil, fmt.Errorf("failed to query Sources: %w", err)
	}
	return &source, nil
}

// GetSourcesCount gets the total count of sources.
func GetSourcesCount(manager *DBManager) (int, error) {
	db, err := manager.GetPrimaryDB()
	if err != nil {
		return 0, fmt.Errorf("failed to get shared DB: %w", err)
	}

	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM Sources").Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count Sources: %w", err)
	}
	return count, nil
}

// DeleteExpiredSources deletes sources older than the specified number of days.
func DeleteExpiredSources(manager *DBManager, days int) error {
	db, err := manager.GetPrimaryDB()
	if err != nil {
		return fmt.Errorf("failed to get shared DB: %w", err)
	}

	query := `
		DELETE FROM Sources
		WHERE CreatedAt < datetime('now', '-${days} days')
	`
	// Note: SQLite date functions. Using string replacement for days.
	// A better approach is to use a parameter, but SQLite's date functions are tricky with parameters for modifiers.
	// Let's use a safer approach with a parameter for the date calculation.
	// Calculate the cutoff date in Go and pass it as a parameter.
	cutoffDate := time.Now().AddDate(0, 0, -days).Format(time.RFC3339)

	query = `
		DELETE FROM Sources
		WHERE CreatedAt < ?
	`

	result, err := db.Exec(query, cutoffDate)
	if err != nil {
		return fmt.Errorf("failed to delete expired Sources: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected for Sources deletion: %w", err)
	}

	// Log the number of deleted rows if needed
	// log.Printf("Deleted %d expired Sources entries", rowsAffected)
	_ = rowsAffected // Use or log as needed

	return nil
}
