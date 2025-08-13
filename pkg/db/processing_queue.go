package db

import (
	"database/sql"
	"fmt"
	"strings"
	"time"
)

// InsertProcessingQueueItem inserts a new item into the ProcessingQueue.
func InsertProcessingQueueItem(manager *DBManager, item ProcessingQueue) error {
	// ProcessingQueue is likely in a shared database or a specific shard.
	// Based on the design, it seems to be in a shared context.
	// Let's use the shared DB.
	db, err := manager.GetPrimaryDB()
	if err != nil {
		return fmt.Errorf("failed to get shared DB: %w", err)
	}

	query := `
		INSERT INTO ProcessingQueue (
			ClientAppID, MessageID, UserID, Url, Language, StatusID, CreatedAt, ProcessedAt, RetryCount, ErrorMessage
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	// Handle nullable fields
	var processedAt interface{}
	if item.ProcessedAt != nil {
		processedAt = *item.ProcessedAt
	} else {
		processedAt = nil
	}

	var retryCount interface{}
	if item.RetryCount != nil {
		retryCount = *item.RetryCount
	} else {
		retryCount = nil
	}

	var errorMessage interface{}
	if item.ErrorMessage != nil {
		errorMessage = *item.ErrorMessage
	} else {
		errorMessage = nil
	}

	createdAt := item.CreatedAt
	if createdAt == "" {
		createdAt = time.Now().Format(time.RFC3339)
	}

	_, err = db.Exec(query,
		item.ClientAppID, item.MessageID, item.UserID,
		item.Url, item.Language, item.StatusID,
		createdAt, processedAt, retryCount, errorMessage,
	)
	if err != nil {
		return fmt.Errorf("failed to insert into ProcessingQueue: %w", err)
	}
	return nil
}

// GetProcessingQueueItemsByStatus retrieves items from the ProcessingQueue with a specific status.
func GetProcessingQueueItemsByStatus(manager *DBManager, statusID int8, limit int) ([]ProcessingQueue, error) {
	db, err := manager.GetPrimaryDB()
	if err != nil {
		return nil, fmt.Errorf("failed to get shared DB: %w", err)
	}

	query := `
		SELECT ClientAppID, MessageID, UserID, Url, Language, StatusID, CreatedAt, ProcessedAt, RetryCount, ErrorMessage
		FROM ProcessingQueue
		WHERE StatusID = ?
		ORDER BY CreatedAt ASC
		LIMIT ?
	`

	rows, err := db.Query(query, statusID, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to query ProcessingQueue by status %d: %w", statusID, err)
	}
	defer rows.Close()

	var items []ProcessingQueue
	for rows.Next() {
		var item ProcessingQueue
		var processedAt sql.NullString
		var retryCount sql.NullInt16 // TINYINT maps to int16 in Go
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


// UpdateProcessingQueueStatus updates the status of an item in the ProcessingQueue.
func UpdateProcessingQueueStatus(manager *DBManager, clientAppID int8, messageID int64, userID int64, url string, statusID int8) error {
	return UpdateProcessingQueueStatusWithDetails(manager, clientAppID, messageID, userID, url, statusID, nil, nil)
}

// UpdateProcessingQueueStatusWithDetails updates the status and potentially other fields of an item.
func UpdateProcessingQueueStatusWithDetails(manager *DBManager, clientAppID int8, messageID int64, userID int64, url string, statusID int8, processedAt *string, errorMessage *string) error {
	db, err := manager.GetPrimaryDB()
	if err != nil {
		return fmt.Errorf("failed to get shared DB: %w", err)
	}

	// Build dynamic query based on provided fields
	setParts := []string{"StatusID = ?"}
	args := []interface{}{statusID}

	if processedAt != nil {
		setParts = append(setParts, "ProcessedAt = ?")
		args = append(args, *processedAt)
	} else if statusID >= 50 { // Assuming statuses >= 50 are terminal
		setParts = append(setParts, "ProcessedAt = ?")
		args = append(args, time.Now().Format(time.RFC3339))
	}

	if errorMessage != nil {
		setParts = append(setParts, "ErrorMessage = ?")
		args = append(args, *errorMessage)
		// Increment retry count if there's an error
		setParts = append(setParts, "RetryCount = COALESCE(RetryCount, 0) + 1")
	}

	query := fmt.Sprintf(`
		UPDATE ProcessingQueue
		SET %s
		WHERE ClientAppID = ? AND MessageID = ? AND UserID = ? AND Url = ?
	`, strings.Join(setParts, ", "))

	args = append(args, clientAppID, messageID, userID, url)

	result, err := db.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("failed to update ProcessingQueue status: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		// This might be okay if the item was already processed by another worker
		// Consider logging this at debug level
		// log.Printf("Warning: No rows updated for ProcessingQueue item (ClientAppID=%d, MessageID=%d, UserID=%d, Url=%s)", clientAppID, messageID, userID, url)
	}

	return nil
}


// GetProcessingQueueCountByStatus gets counts of items grouped by status.
func GetProcessingQueueCountByStatus(manager *DBManager) (map[string]int, error) {
	db, err := manager.GetPrimaryDB()
	if err != nil {
		return nil, fmt.Errorf("failed to get shared DB: %w", err)
	}

	query := `
		SELECT ps.DisplayName, COUNT(pq.StatusID) as Count
		FROM ProcessingStatus ps
		LEFT JOIN ProcessingQueue pq ON ps.ID = pq.StatusID
		GROUP BY ps.ID, ps.DisplayName
	`

	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query ProcessingQueue counts: %w", err)
	}
	defer rows.Close()

	counts := make(map[string]int)
	for rows.Next() {
		var status string
		var count int
		err := rows.Scan(&status, &count)
		if err != nil {
			return nil, fmt.Errorf("failed to scan ProcessingQueue count row: %w", err)
		}
		counts[status] = count
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating ProcessingQueue count rows: %w", err)
	}

	return counts, nil
}

// GetProcessingQueueItemsByStatusAndCondition retrieves items based on status and additional conditions.
func GetProcessingQueueItemsByStatusAndCondition(manager *DBManager, statusID int8, condition string, limit int) ([]ProcessingQueue, error) {
	db, err := manager.GetPrimaryDB()
	if err != nil {
		return nil, fmt.Errorf("failed to get shared DB: %w", err)
	}

	baseQuery := `
		SELECT ClientAppID, MessageID, UserID, Url, Language, StatusID, CreatedAt, ProcessedAt, RetryCount, ErrorMessage
		FROM ProcessingQueue
		WHERE StatusID = ?
	`

	var query string
	var args []interface{}
	args = append(args, statusID)

	if condition != "" {
		query = baseQuery + " AND " + condition
	} else {
		query = baseQuery
	}
	query += " ORDER BY CreatedAt ASC LIMIT ?"
	args = append(args, limit)

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query ProcessingQueue by status %d with condition '%s': %w", statusID, condition, err)
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

// MarkProcessingQueueItemsAsCompleted updates items to 'completed' status.
// This function assumes items are already summarized.
func MarkProcessingQueueItemsAsCompleted(manager *DBManager, items []ProcessingQueue) error {
	if len(items) == 0 {
		return nil
	}

	db, err := manager.GetPrimaryDB()
	if err != nil {
		return fmt.Errorf("failed to get shared DB: %w", err)
	}

	// Use a transaction for atomicity
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// Get the 'completed' status ID
	completedStatusID, err := GetStatusID(db, "completed")
	if err != nil {
		return fmt.Errorf("failed to get 'completed' status ID: %w", err)
	}

	for _, item := range items {
		// Update status to 'completed'
		err = UpdateProcessingQueueStatusWithDetails(manager, item.ClientAppID, item.MessageID, item.UserID, item.Url, completedStatusID, nil, nil)
		if err != nil {
			return fmt.Errorf("failed to update item to 'completed': %w", err)
		}
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}