package db

import "database/sql"

func InsertSource(db *sql.DB, source Source) error {
	// Insert source into Sources table
	return nil
}

func GetSource(db *sql.DB, url string, language string) (*Source, error) {
	// Get source from Sources table
	return nil, nil
}

func GetSourcesCount(db *sql.DB) (int, error) {
	// Get count of sources in Sources table
	return 0, nil
}

func DeleteExpiredSources(db *sql.DB, days int) error {
	// Delete sources older than specified days
	return nil
}
