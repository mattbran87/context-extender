package database

import (
	"time"
)

// Configuration types moved here from types.go to avoid circular imports

// BackendType represents the type of database backend
type BackendType string

const (
	BackendPureGoSQLite BackendType = "pure_go_sqlite"
	BackendCGOSQLite    BackendType = "cgo_sqlite"
	BackendAuto         BackendType = "auto"
)

// DatabaseConfig represents database configuration
type DatabaseConfig struct {
	Backend           BackendType            `json:"backend"`
	DatabasePath      string                 `json:"database_path"`
	ConnectionTimeout time.Duration          `json:"connection_timeout"`
	QueryTimeout      time.Duration          `json:"query_timeout"`
	BackendOptions    map[string]interface{} `json:"backend_options,omitempty"`
}

// Configuration types only - other types moved to types.go

// DefaultDatabaseConfig returns a default database configuration
func DefaultDatabaseConfig() *DatabaseConfig {
	return &DatabaseConfig{
		Backend:           BackendAuto,
		DatabasePath:      getDefaultDatabasePath(),
		ConnectionTimeout: 30 * time.Second,
		QueryTimeout:      30 * time.Second,
	}
}


