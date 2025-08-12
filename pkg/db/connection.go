package db

import (
	"database/sql"
	"fmt"
	"path/filepath"
	"sync"

	_ "modernc.org/sqlite" // SQLite driver
)

// DBManager manages connections to potentially sharded databases.
// For simplicity, we'll use one connection per client app ID.
type DBManager struct {
	basePath string
	dbs      map[int8]*sql.DB
	mutex    sync.RWMutex
}

// NewDBManager creates a new database manager.
func NewDBManager(basePath string) *DBManager {
	return &DBManager{
		basePath: basePath,
		dbs:      make(map[int8]*sql.DB),
	}
}

// GetDBForClient retrieves or creates a database connection for a specific client app ID.
// This implements the sharding strategy: one SQLite file per client.
func (dm *DBManager) GetDBForClient(clientAppID int8) (*sql.DB, error) {
	dm.mutex.RLock()
	db, exists := dm.dbs[clientAppID]
	dm.mutex.RUnlock()

	if exists {
		return db, nil
	}

	// Need to create the connection
	dm.mutex.Lock()
	defer dm.mutex.Unlock()

	// Double-check after acquiring write lock
	db, exists = dm.dbs[clientAppID]
	if exists {
		return db, nil
	}

	// Construct database file path
	dbPath := filepath.Join(dm.basePath, fmt.Sprintf("client_%d.db", clientAppID))
	dataSourceName := fmt.Sprintf("file:%s?cache=shared&mode=rwc&_pragma=foreign_keys(1)", dbPath) // Enable FKs

	newDB, err := sql.Open("sqlite", dataSourceName)
	if err != nil {
		return nil, fmt.Errorf("failed to open database for client %d at %s: %w", clientAppID, dbPath, err)
	}

	// Configure connection pool if needed
	// newDB.SetMaxOpenConns(1) // SQLite is file-based, often 1 conn is sufficient for single writer

	// Run migrations for this new shard if necessary (or ensure they are run externally)
	// For now, assume migrations are run separately or on startup for default shards.

	dm.dbs[clientAppID] = newDB
	return newDB, nil
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
		delete(dm.dbs, clientID) // Clear reference
	}
	return firstErr
}

// GetSharedDB retrieves a connection to a shared database (e.g., for non-sharded tables like DictClientApps, ProcessingStatus, Sources, Summaries).
// This assumes a single shared database file.
func (dm *DBManager) GetSharedDB() (*sql.DB, error) {
	// For simplicity, use client ID 0 for the shared database
	// Or define a specific name like 'shared.db'
	return dm.GetDBForClient(0) // Assumes client 0 is used for shared tables
}
