package db

import (
	"database/sql"
	"fmt"
)

// GetClientAppID retrieves the ID for a given client application name from the shared database.
func GetClientAppID(sharedDB *sql.DB, appName string) (int8, error) {
	var id int8
	err := sharedDB.QueryRow("SELECT ID FROM DictClientApps WHERE App = ?", appName).Scan(&id)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, fmt.Errorf("client app '%s' not found", appName)
		}
		return 0, fmt.Errorf("failed to query DictClientApps: %w", err)
	}
	return id, nil
}

// ValidateClientApp checks if a client app ID exists in the shared database.
func ValidateClientApp(sharedDB *sql.DB, clientAppID int8) (bool, error) {
	var exists bool
	// Using EXISTS for efficiency
	err := sharedDB.QueryRow("SELECT EXISTS(SELECT 1 FROM DictClientApps WHERE ID = ?)", clientAppID).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to validate client app ID %d: %w", clientAppID, err)
	}
	return exists, nil
}

// GetAllClientApps retrieves all client applications from the shared database.
func GetAllClientApps(sharedDB *sql.DB) ([]DictClientApp, error) {
	rows, err := sharedDB.Query("SELECT ID, App FROM DictClientApps")
	if err != nil {
		return nil, fmt.Errorf("failed to query DictClientApps: %w", err)
	}
	defer rows.Close()

	var apps []DictClientApp
	for rows.Next() {
		var app DictClientApp
		if err := rows.Scan(&app.ID, &app.App); err != nil {
			return nil, fmt.Errorf("failed to scan DictClientApp row: %w", err)
		}
		apps = append(apps, app)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating DictClientApps rows: %w", err)
	}

	return apps, nil
}
