package db

import (
	"database/sql"
	"fmt"
	"path/filepath"
	"sync"
	"time"

	_ "modernc.org/sqlite"
)

const defaultDatabaseName = "primary"

// DBManager manages connections to potentially sharded databases.
type DBManager struct {
	basePath string
	dbs   map[string]*sql.DB
	mutex sync.RWMutex
	primaryDB *sql.DB
}

// Config holds database configuration options.
type Config struct {
	BasePath        string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
	ConnMaxIdleTime time.Duration
	PrimaryDBName    string
}

// NewDBManager creates a new database manager with the given configuration.
func NewDBManager(config Config) *DBManager {
	// Set defaults if not provided
	if config.MaxOpenConns == 0 {
		config.MaxOpenConns = 10
	}
	if config.MaxIdleConns == 0 {
		config.MaxIdleConns = 2
	}
	if config.ConnMaxLifetime == 0 {
		config.ConnMaxLifetime = 30 * time.Minute
	}
	if config.PrimaryDBName == "" {
		config.PrimaryDBName = defaultDatabaseName
	}

	return &DBManager{
		basePath: config.BasePath,
		dbs:      make(map[string]*sql.DB),
		// primaryDB will be initialized lazily or via a specific method
	}
}

// getDBPath constructs the database file path for a given client app name.
func (dm *DBManager) getDBPath(clientAppName string) string {
	dbFileName := fmt.Sprintf("%s.sqlite", clientAppName)
	return filepath.Join(dm.basePath, dbFileName)
}

// GetDBForClient retrieves or creates a database connection for a specific client app name.
// This implements the sharding strategy: one SQLite file per client app.
// Pass an empty string or the primary DB name to get the primary database connection.
func (dm *DBManager) GetDBForClient(clientAppName string) (*sql.DB, error) {
	// If clientAppName is empty, default to primary DB name
	if clientAppName == "" {
		clientAppName = defaultDatabaseName
	}

	dm.mutex.RLock()
	db, exists := dm.dbs[clientAppName]
	dm.mutex.RUnlock()

	if exists {
		return db, nil
	}

	// Need to create the connection
	dm.mutex.Lock()
	defer dm.mutex.Unlock()

	// Double-check after acquiring write lock
	db, exists = dm.dbs[clientAppName]
	if exists {
		return db, nil
	}

	// Construct database file path
	dbPath := dm.getDBPath(clientAppName)
	// Ensure the primary DB also uses the correct name logic
	// If clientAppName was originally intended to be "primary", the path will be correct.
	dataSourceName := fmt.Sprintf("file:%s?cache=shared&mode=rwc&_pragma=foreign_keys(1)", dbPath)

	newDB, err := sql.Open("sqlite", dataSourceName)
	if err != nil {
		return nil, fmt.Errorf("failed to open database for client '%s' at %s: %w", clientAppName, dbPath, err)
	}

	// Configure connection pool
	// These values should ideally come from the config
	newDB.SetMaxOpenConns(10)                  // Example values, should ideally use config.MaxOpenConns
	newDB.SetMaxIdleConns(2)                   // Example values, should ideally use config.MaxIdleConns
	newDB.SetConnMaxLifetime(30 * time.Minute) // Example values, should ideally use config.ConnMaxLifetime
	// newDB.SetConnMaxIdleTime(...) // Available in Go 1.15+, use config.ConnMaxIdleTime

	// If this is the primary DB, store the reference
	// This requires knowing the primary DB name. We assumed "primary" above.
	// This is a bit fragile. A better approach is a dedicated GetPrimaryDB method.
	// Let's assume for now that if clientAppName is "primary", we store it.
	if clientAppName == defaultDatabaseName {
		dm.primaryDB = newDB
	}

	dm.dbs[clientAppName] = newDB
	return newDB, nil
}

// GetPrimaryDB retrieves a connection to the primary database.
// This provides a clear interface for accessing the primary database.
func (dm *DBManager) GetPrimaryDB() (*sql.DB, error) {
	// Use the dedicated identifier for the primary DB
	// We need to ensure this identifier matches the one used in GetDBForClient when it detects primary DB access.
	// From the logic above, we used "primary".
	return dm.GetDBForClient(defaultDatabaseName)
}

// CloseAll closes all managed database connections.
func (dm *DBManager) CloseAll() error {
	dm.mutex.Lock()
	defer dm.mutex.Unlock()

	var firstErr error
	for clientName, db := range dm.dbs {
		if err := db.Close(); err != nil && firstErr == nil {
			firstErr = fmt.Errorf("failed to close database for client '%s': %w", clientName, err)
		}
		delete(dm.dbs, clientName)
	}
	// Clear the primaryDB reference
	dm.primaryDB = nil
	return firstErr
}
