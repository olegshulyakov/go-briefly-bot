package db

import (
	"database/sql"
	"fmt"
)

// GetStatusID retrieves the ID for a given status display name from the shared database.
func GetStatusID(sharedDB *sql.DB, displayName string) (int8, error) {
	var id int8
	err := sharedDB.QueryRow("SELECT ID FROM ProcessingStatus WHERE DisplayName = ?", displayName).Scan(&id)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, fmt.Errorf("processing status '%s' not found", displayName)
		}
		return 0, fmt.Errorf("failed to query ProcessingStatus: %w", err)
	}
	return id, nil
}

// GetStatusDisplayName retrieves the display name for a given status ID from the shared database.
func GetStatusDisplayName(sharedDB *sql.DB, statusID int8) (string, error) {
	var displayName string
	err := sharedDB.QueryRow("SELECT DisplayName FROM ProcessingStatus WHERE ID = ?", statusID).Scan(&displayName)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("processing status ID %d not found", statusID)
		}
		return "", fmt.Errorf("failed to query ProcessingStatus: %w", err)
	}
	return displayName, nil
}

// GetAllProcessingStatuses retrieves all processing statuses from the shared database.
func GetAllProcessingStatuses(sharedDB *sql.DB) ([]ProcessingStatus, error) {
	rows, err := sharedDB.Query("SELECT ID, DisplayName FROM ProcessingStatus ORDER BY ID")
	if err != nil {
		return nil, fmt.Errorf("failed to query ProcessingStatus: %w", err)
	}
	defer rows.Close()

	var statuses []ProcessingStatus
	for rows.Next() {
		var status ProcessingStatus
		if err := rows.Scan(&status.ID, &status.DisplayName); err != nil {
			return nil, fmt.Errorf("failed to scan ProcessingStatus row: %w", err)
		}
		statuses = append(statuses, status)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating ProcessingStatus rows: %w", err)
	}

	return statuses, nil
}
