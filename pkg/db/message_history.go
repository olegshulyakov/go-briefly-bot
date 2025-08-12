package db

import "database/sql"

func InsertMessageHistory(db *sql.DB, message MessageHistory) error {
	// Insert message into MessageHistory table
	return nil
}

func GetMessageHistoryCount(db *sql.DB) (int, error) {
	// Get count of messages in MessageHistory table
	return 0, nil
}
