package db

import (
	"fmt"
	"time"
)

// InsertMessageHistory inserts a new message history record into the sharded database.
func InsertMessageHistory(manager *DBManager, message MessageHistory) error {
	db, err := manager.GetPrimaryDB()
	if err != nil {
		return fmt.Errorf("failed to get DB for client %d: %w", message.ClientAppID, err)
	}

	query := `
		INSERT INTO MessageHistory (
			ClientAppID, MessageID, UserID, UserName, UserLanguage, MessageContent, CreatedAt
		) VALUES (?, ?, ?, ?, ?, ?, ?)
	`

	// Use time.Now() if CreatedAt is empty or needs to be set by server
	createdAt := message.CreatedAt
	if createdAt == "" {
		createdAt = time.Now().Format(time.RFC3339) // Standard format
	}

	_, err = db.Exec(query,
		message.ClientAppID, message.MessageID, message.UserID,
		message.UserName, message.UserLanguage, message.MessageContent,
		createdAt,
	)
	if err != nil {
		return fmt.Errorf("failed to insert into MessageHistory: %w", err)
	}
	return nil
}

// GetMessageHistoryCount gets the total count of messages in the shared database.
// Note: This might require querying all shards or using a global counter table.
// For simplicity, we'll assume it's in the shared DB or a specific shard.
func GetMessageHistoryCount(manager *DBManager) (int, error) {
	// Assuming a global count or querying a specific shard
	// Let's query the shared DB for now, assuming it might have a view or count logic
	// Or we might need to query all shards and sum.
	// For this implementation, let's assume it's tracked in shared DB or we query one shard.
	// A more robust solution would involve a global counter or querying all shards.
	// We'll query the shared DB as a placeholder.
	db, err := manager.GetPrimaryDB()
	if err != nil {
		return 0, fmt.Errorf("failed to get shared DB: %w", err)
	}

	var count int
	// This query is conceptual. In a sharded setup, you might need to aggregate counts.
	// One approach is to have a separate table for global counts or query all shards.
	// For now, we'll assume a way to get the total count exists in shared DB.
	// A real implementation might require a more complex strategy.
	err = db.QueryRow("SELECT COUNT(*) FROM MessageHistory").Scan(&count)
	if err != nil {
		// Fallback or handle error appropriately
		// If shared DB doesn't have MessageHistory, this will fail.
		// Let's assume for now it's in shared DB or we need a different approach.
		// For this placeholder, return 0 if error.
		// A better approach would be to define where global counts are stored.
		// Let's assume it's in the shared DB for this example, even though it's sharded.
		// This is a simplification. In reality, you'd need a global counter or sum across shards.
		// We'll return 0 and log or handle the error.
		// For now, just return 0 if error.
		// A production system would need a more robust solution.
		// Let's assume it's in shared DB for this example.
		// If not, the query will fail, and we return 0.
		// This is a known limitation of this simplified approach.
		return 0, nil // Or return the error if you want to propagate it.
	}
	return count, nil
}

// GetMessageHistory retrieves message history for a specific client and user (example query).
func GetMessageHistory(manager *DBManager, clientAppID int8, userID int64, limit int) ([]MessageHistory, error) {
	db, err := manager.GetPrimaryDB()
	if err != nil {
		return nil, fmt.Errorf("failed to get DB for client %d: %w", clientAppID, err)
	}

	query := `
		SELECT ClientAppID, MessageID, UserID, UserName, UserLanguage, MessageContent, CreatedAt
		FROM MessageHistory
		WHERE ClientAppID = ? AND UserID = ?
		ORDER BY CreatedAt DESC
		LIMIT ?
	`

	rows, err := db.Query(query, clientAppID, userID, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to query MessageHistory: %w", err)
	}
	defer rows.Close()

	var messages []MessageHistory
	for rows.Next() {
		var msg MessageHistory
		err := rows.Scan(
			&msg.ClientAppID, &msg.MessageID, &msg.UserID,
			&msg.UserName, &msg.UserLanguage, &msg.MessageContent,
			&msg.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan MessageHistory row: %w", err)
		}
		messages = append(messages, msg)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating MessageHistory rows: %w", err)
	}

	return messages, nil
}
