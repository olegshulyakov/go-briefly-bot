package db

import (
	"database/sql"
	"fmt"
	"path/filepath"
	"sync"
	"time"

	_ "modernc.org/sqlite"
)

// DBManager manages connections to potentially sharded databases.
type DBManager struct {
	basePath string
	dbs      map[int8]*sql.DB
	mutex    sync.RWMutex
}

// Config holds database configuration options.
type Config struct {
	BasePath          string
	MaxOpenConns      int
	MaxIdleConns      int
	ConnMaxLifetime   time.Duration
	ConnMaxIdleTime   time.Duration
	SharedDBClientID  int8 // Client ID used for the shared database file
}

// NewDBManager creates a new database manager with the given configuration.
func NewDBManager(config Config) *DBManager {
	// Set defaults if not provided
	if config.MaxOpenConns == 0 {
		config.MaxOpenConns = 10 // Reasonable default for SQLite
	}
	if config.MaxIdleConns == 0 {
		config.MaxIdleConns = 2
	}
	if config.ConnMaxLifetime == 0 {
		config.ConnMaxLifetime = 30 * time.Minute
	}
	if config.SharedDBClientID == 0 {
		config.SharedDBClientID = 0 // Default to client 0 for shared DB
	}

	return &DBManager{
		basePath: config.BasePath,
		dbs:      make(map[int8]*sql.DB),
	}
}

// getDBPath constructs the database file path for a given client ID.
func (dm *DBManager) getDBPath(clientAppID int8) string {
	// Use a specific name for the shared DB, otherwise shard by client ID
	if clientAppID == 0 { // Assuming 0 is reserved for shared
		return filepath.Join(dm.basePath, "shared.db")
	}
	return filepath.Join(dm.basePath, fmt.Sprintf("client_%d.db", clientAppID))
}

// GetDBForClient retrieves or creates a database connection for a specific client app ID.
// This implements the sharding strategy: one SQLite file per client.
func (dm *DBManager) GetDBForClient(clientAppID int8) (*sql.DB, error) {
	// Use shared DB client ID for shared tables
	actualClientID := clientAppID
	if clientAppID < 0 { // Convention: negative ID means use shared DB
		actualClientID = 0 // Or use config.SharedDBClientID
	}

	dm.mutex.RLock()
	db, exists := dm.dbs[actualClientID]
	dm.mutex.RUnlock()

	if exists {
		return db, nil
	}

	// Need to create the connection
	dm.mutex.Lock()
	defer dm.mutex.Unlock()

	// Double-check after acquiring write lock
	db, exists = dm.dbs[actualClientID]
	if exists {
		return db, nil
	}

	// Construct database file path
	dbPath := dm.getDBPath(actualClientID)
	dataSourceName := fmt.Sprintf("file:%s?cache=shared&mode=rwc&_pragma=foreign_keys(1)", dbPath)

	newDB, err := sql.Open("sqlite", dataSourceName)
	if err != nil {
		return nil, fmt.Errorf("failed to open database for client %d at %s: %w", actualClientID, dbPath, err)
	}

	// Configure connection pool
	newDB.SetMaxOpenConns(10) // Example values, should be configurable
	newDB.SetMaxIdleConns(2)
	newDB.SetConnMaxLifetime(30 * time.Minute)
	// newDB.SetConnMaxIdleTime(...) // Available in Go 1.15+

	dm.dbs[actualClientID] = newDB
	return newDB, nil
}

// GetSharedDB retrieves a connection to the shared database.
func (dm *DBManager) GetSharedDB() (*sql.DB, error) {
	return dm.GetDBForClient(0) // Use client ID 0 for shared DB
}

// CloseAll closes all managed database connections.
func (dm *DBManager) CloseAll() error {
	dm.mutex.Lock()
	defer dm.mutex.Unlock()

	var firstErr error
	for clientID, db := range dm.dbs {
		if err := db.Close(); err != nil && firstErr == nil {
			firstErr = fmt.Errorf("failed to close database for client %d: %w", clientID, err)
		}
		delete(dm.dbs, clientID)
	}
	return firstErr
}