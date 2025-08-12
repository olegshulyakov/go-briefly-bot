package db

import "database/sql"

func InsertSummary(db *sql.DB, summary Summary) error {
	// Insert summary into Summaries table
	return nil
}

func GetSummariesForClient(db *sql.DB, clientAppID int8, limit int) ([]ProcessingQueue, error) {
	// Get summaries for a specific client
	return nil, nil
}

func GetSummariesCount(db *sql.DB) (int, error) {
	// Get count of summaries in Summaries table
	return 0, nil
}

func DeleteExpiredSummaries(db *sql.DB, days int) error {
	// Delete summaries older than specified days
	return nil
}
