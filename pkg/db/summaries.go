package db

import (
	"database/sql"
	"fmt"
	"time"
)

// InsertSummary inserts a new summary into the Summaries table.
func InsertSummary(manager *DBManager, summary Summary) error {
	db, err := manager.GetPrimaryDB()
	if err != nil {
		return fmt.Errorf("failed to get shared DB: %w", err)
	}

	query := `
		INSERT OR REPLACE INTO Summaries (Url, Language, Summary, CreatedAt)
		VALUES (?, ?, ?, ?)
	`

	createdAt := summary.CreatedAt
	if createdAt == "" {
		createdAt = time.Now().Format(time.RFC3339)
	}

	_, err = db.Exec(query, summary.Url, summary.Language, summary.Summary, createdAt)
	if err != nil {
		return fmt.Errorf("failed to insert into Summaries: %w", err)
	}
	return nil
}

// GetSummary retrieves a summary by URL and language.
func GetSummary(manager *DBManager, url string, language string) (*Summary, error) {
	db, err := manager.GetPrimaryDB()
	if err != nil {
		return nil, fmt.Errorf("failed to get shared DB: %w", err)
	}

	query := `
		SELECT Url, Language, Summary, CreatedAt
		FROM Summaries
		WHERE Url = ? AND Language = ?
	`

	row := db.QueryRow(query, url, language)
	var summary Summary
	err = row.Scan(&summary.Url, &summary.Language, &summary.Summary, &summary.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Not found
		}
		return nil, fmt.Errorf("failed to query Summaries: %w", err)
	}
	return &summary, nil
}

// GetSummariesForClient retrieves summaries for items in the ProcessingQueue with 'summarized' status for a specific client.
func GetSummariesForClient(manager *DBManager, clientAppID int8, limit int) ([]ProcessingQueue, error) {
	db, err := manager.GetPrimaryDB()
	if err != nil {
		return nil, fmt.Errorf("failed to get shared DB: %w", err)
	}

	// Get the 'summarized' status ID
	summarizedStatusID, err := GetStatusID(db, "summarized")
	if err != nil {
		return nil, fmt.Errorf("failed to get 'summarized' status ID: %w", err)
	}

	query := `
		SELECT pq.ClientAppID, pq.MessageID, pq.UserID, pq.Url, pq.Language, pq.StatusID, pq.CreatedAt, pq.ProcessedAt, pq.RetryCount, pq.ErrorMessage
		FROM ProcessingQueue pq
		WHERE pq.ClientAppID = ? AND pq.StatusID = ?
		ORDER BY pq.CreatedAt ASC
		LIMIT ?
	`

	rows, err := db.Query(query, clientAppID, summarizedStatusID, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to query ProcessingQueue for summarized items: %w", err)
	}
	defer rows.Close()

	var items []ProcessingQueue
	for rows.Next() {
		var item ProcessingQueue
		var processedAt sql.NullString
		var retryCount sql.NullInt16
		var errorMessage sql.NullString

		err := rows.Scan(
			&item.ClientAppID, &item.MessageID, &item.UserID,
			&item.Url, &item.Language, &item.StatusID,
			&item.CreatedAt, &processedAt, &retryCount, &errorMessage,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan ProcessingQueue row: %w", err)
		}

		if processedAt.Valid {
			item.ProcessedAt = &processedAt.String
		}
		if retryCount.Valid {
			val := int8(retryCount.Int16)
			item.RetryCount = &val
		}
		if errorMessage.Valid {
			item.ErrorMessage = &errorMessage.String
		}

		items = append(items, item)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating ProcessingQueue rows: %w", err)
	}

	return items, nil
}

// GetSummariesCount gets the total count of summaries.
func GetSummariesCount(manager *DBManager) (int, error) {
	db, err := manager.GetPrimaryDB()
	if err != nil {
		return 0, fmt.Errorf("failed to get shared DB: %w", err)
	}

	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM Summaries").Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count Summaries: %w", err)
	}
	return count, nil
}

// DeleteExpiredSummaries deletes summaries older than the specified number of days.
func DeleteExpiredSummaries(manager *DBManager, days int) error {
	db, err := manager.GetPrimaryDB()
	if err != nil {
		return fmt.Errorf("failed to get shared DB: %w", err)
	}

	// Calculate cutoff date in Go
	cutoffDate := time.Now().AddDate(0, 0, -days).Format(time.RFC3339)

	query := `DELETE FROM Summaries WHERE CreatedAt < ?`

	result, err := db.Exec(query, cutoffDate)
	if err != nil {
		return fmt.Errorf("failed to delete expired Summaries: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected for Summaries deletion: %w", err)
	}

	// Log the number of deleted rows if needed
	// log.Printf("Deleted %d expired Summaries entries", rowsAffected)
	_ = rowsAffected // Use or log as needed

	return nil
}
