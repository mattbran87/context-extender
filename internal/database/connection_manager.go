package database

import (
	"context"
	"database/sql"
	"fmt"
	"sync"
	"time"

	_ "modernc.org/sqlite"
)

// ConnectionManager handles database connections with pooling and optimization
type ConnectionManager struct {
	db              *sql.DB
	config          *ConnectionConfig
	stats           *ConnectionStats
	mu              sync.RWMutex
	initialized     bool
	lastPingTime    time.Time
	pingInterval    time.Duration
}

// ConnectionConfig holds connection configuration
type ConnectionConfig struct {
	DatabasePath      string        `json:"database_path"`
	MaxOpenConns      int           `json:"max_open_conns"`
	MaxIdleConns      int           `json:"max_idle_conns"`
	ConnMaxLifetime   time.Duration `json:"conn_max_lifetime"`
	ConnMaxIdleTime   time.Duration `json:"conn_max_idle_time"`
	EnableWAL         bool          `json:"enable_wal"`
	CacheSize         int           `json:"cache_size"`         // In pages
	MMAPSize          int64         `json:"mmap_size"`          // In bytes
	SynchronousMode   string        `json:"synchronous_mode"`   // OFF, NORMAL, FULL
	TempStore         string        `json:"temp_store"`         // DEFAULT, FILE, MEMORY
	BusyTimeout       time.Duration `json:"busy_timeout"`
	EnableForeignKeys bool          `json:"enable_foreign_keys"`
	QueryTimeout      time.Duration `json:"query_timeout"`
}

// ConnectionStats tracks connection performance metrics
type ConnectionStats struct {
	mu                sync.RWMutex
	OpenConnections   int           `json:"open_connections"`
	IdleConnections   int           `json:"idle_connections"`
	TotalQueries      int64         `json:"total_queries"`
	FailedQueries     int64         `json:"failed_queries"`
	AverageQueryTime  time.Duration `json:"average_query_time"`
	LastQueryTime     time.Time     `json:"last_query_time"`
	ConnectionsOpened int64         `json:"connections_opened"`
	ConnectionsClosed int64         `json:"connections_closed"`
	MaxOpenReached    int           `json:"max_open_reached"`
}

// DefaultConnectionConfig returns optimized default configuration
func DefaultConnectionConfig(databasePath string) *ConnectionConfig {
	return &ConnectionConfig{
		DatabasePath:      databasePath,
		MaxOpenConns:      25,                // Allow reasonable concurrency
		MaxIdleConns:      5,                 // Keep some connections warm
		ConnMaxLifetime:   time.Hour,         // Rotate connections hourly
		ConnMaxIdleTime:   15 * time.Minute,  // Close idle connections
		EnableWAL:         true,              // Better concurrency
		CacheSize:         10000,             // 10K pages (~40MB with 4K pages)
		MMAPSize:          268435456,         // 256MB memory mapping
		SynchronousMode:   "NORMAL",          // Good balance of speed/safety
		TempStore:         "memory",          // Use memory for temp tables
		BusyTimeout:       30 * time.Second,  // Wait for locks
		EnableForeignKeys: true,              // Enforce referential integrity
		QueryTimeout:      30 * time.Second,  // Prevent hanging queries
	}
}

// NewConnectionManager creates a new connection manager
func NewConnectionManager(config *ConnectionConfig) *ConnectionManager {
	if config == nil {
		config = DefaultConnectionConfig("./database.db")
	}

	return &ConnectionManager{
		config:       config,
		stats:        &ConnectionStats{},
		pingInterval: 5 * time.Minute,
	}
}

// Initialize initializes the connection manager and establishes database connection
func (cm *ConnectionManager) Initialize(ctx context.Context) error {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	if cm.initialized {
		return nil
	}

	// Open database connection
	db, err := sql.Open("sqlite", cm.config.DatabasePath)
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}

	// Configure connection pool
	db.SetMaxOpenConns(cm.config.MaxOpenConns)
	db.SetMaxIdleConns(cm.config.MaxIdleConns)
	db.SetConnMaxLifetime(cm.config.ConnMaxLifetime)
	db.SetConnMaxIdleTime(cm.config.ConnMaxIdleTime)

	// Apply SQLite optimizations
	if err := cm.applySQLiteOptimizations(ctx, db); err != nil {
		db.Close()
		return fmt.Errorf("failed to apply SQLite optimizations: %w", err)
	}

	// Test connection
	if err := db.PingContext(ctx); err != nil {
		db.Close()
		return fmt.Errorf("failed to ping database: %w", err)
	}

	cm.db = db
	cm.initialized = true
	cm.lastPingTime = time.Now()

	// Start background monitoring
	go cm.backgroundMonitor()

	return nil
}

// applySQLiteOptimizations applies performance optimizations to SQLite
func (cm *ConnectionManager) applySQLiteOptimizations(ctx context.Context, db *sql.DB) error {
	optimizations := []struct {
		name   string
		pragma string
	}{
		{"WAL Mode", fmt.Sprintf("PRAGMA journal_mode=%s", cm.getJournalMode())},
		{"Foreign Keys", fmt.Sprintf("PRAGMA foreign_keys=%s", cm.getBooleanPragma(cm.config.EnableForeignKeys))},
		{"Synchronous", fmt.Sprintf("PRAGMA synchronous=%s", cm.config.SynchronousMode)},
		{"Cache Size", fmt.Sprintf("PRAGMA cache_size=%d", cm.config.CacheSize)},
		{"Temp Store", fmt.Sprintf("PRAGMA temp_store=%s", cm.config.TempStore)},
		{"MMAP Size", fmt.Sprintf("PRAGMA mmap_size=%d", cm.config.MMAPSize)},
		{"Busy Timeout", fmt.Sprintf("PRAGMA busy_timeout=%d", int(cm.config.BusyTimeout.Milliseconds()))},
		{"Optimize", "PRAGMA optimize"},
	}

	for _, opt := range optimizations {
		if _, err := db.ExecContext(ctx, opt.pragma); err != nil {
			return fmt.Errorf("failed to apply %s (%s): %w", opt.name, opt.pragma, err)
		}
	}

	return nil
}

// getJournalMode returns the appropriate journal mode
func (cm *ConnectionManager) getJournalMode() string {
	if cm.config.EnableWAL {
		return "WAL"
	}
	return "DELETE"
}

// getBooleanPragma converts boolean to SQLite pragma value
func (cm *ConnectionManager) getBooleanPragma(value bool) string {
	if value {
		return "ON"
	}
	return "OFF"
}

// GetConnection returns the database connection
func (cm *ConnectionManager) GetConnection() (*sql.DB, error) {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	if !cm.initialized {
		return nil, fmt.Errorf("connection manager not initialized")
	}

	return cm.db, nil
}

// ExecuteWithTimeout executes a query with timeout and statistics tracking
func (cm *ConnectionManager) ExecuteWithTimeout(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	start := time.Now()

	// Create timeout context
	timeoutCtx, cancel := context.WithTimeout(ctx, cm.config.QueryTimeout)
	defer cancel()

	// Execute query
	result, err := cm.db.ExecContext(timeoutCtx, query, args...)

	// Update statistics
	cm.updateQueryStats(time.Since(start), err)

	return result, err
}

// QueryWithTimeout executes a query with timeout and statistics tracking
func (cm *ConnectionManager) QueryWithTimeout(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	start := time.Now()

	// Create timeout context
	timeoutCtx, cancel := context.WithTimeout(ctx, cm.config.QueryTimeout)
	defer cancel()

	// Execute query
	rows, err := cm.db.QueryContext(timeoutCtx, query, args...)

	// Update statistics
	cm.updateQueryStats(time.Since(start), err)

	return rows, err
}

// QueryRowWithTimeout executes a single row query with timeout
func (cm *ConnectionManager) QueryRowWithTimeout(ctx context.Context, query string, args ...interface{}) *sql.Row {
	start := time.Now()

	// Create timeout context
	timeoutCtx, cancel := context.WithTimeout(ctx, cm.config.QueryTimeout)
	defer cancel()

	// Execute query
	row := cm.db.QueryRowContext(timeoutCtx, query, args...)

	// Update statistics (always count as successful for QueryRow)
	cm.updateQueryStats(time.Since(start), nil)

	return row
}

// updateQueryStats updates query performance statistics
func (cm *ConnectionManager) updateQueryStats(duration time.Duration, err error) {
	cm.stats.mu.Lock()
	defer cm.stats.mu.Unlock()

	cm.stats.TotalQueries++
	cm.stats.LastQueryTime = time.Now()

	if err != nil {
		cm.stats.FailedQueries++
	}

	// Update average query time (simple moving average)
	if cm.stats.TotalQueries == 1 {
		cm.stats.AverageQueryTime = duration
	} else {
		// Weighted average favoring recent queries
		weight := 0.1
		cm.stats.AverageQueryTime = time.Duration(
			float64(cm.stats.AverageQueryTime)*(1-weight) + float64(duration)*weight,
		)
	}
}

// GetStats returns current connection statistics
func (cm *ConnectionManager) GetStats() *ConnectionStats {
	cm.stats.mu.RLock()
	defer cm.stats.mu.RUnlock()

	// Get current connection pool stats
	stats := cm.db.Stats()

	// Create a copy with current data
	return &ConnectionStats{
		OpenConnections:   stats.OpenConnections,
		IdleConnections:   stats.Idle,
		TotalQueries:      cm.stats.TotalQueries,
		FailedQueries:     cm.stats.FailedQueries,
		AverageQueryTime:  cm.stats.AverageQueryTime,
		LastQueryTime:     cm.stats.LastQueryTime,
		ConnectionsOpened: int64(stats.OpenConnections), // Use current open connections
		ConnectionsClosed: 0, // Will be calculated differently
		MaxOpenReached:    stats.MaxOpenConnections,
	}
}

// Ping tests the database connection
func (cm *ConnectionManager) Ping(ctx context.Context) error {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	if !cm.initialized {
		return fmt.Errorf("connection manager not initialized")
	}

	err := cm.db.PingContext(ctx)
	if err == nil {
		cm.lastPingTime = time.Now()
	}

	return err
}

// backgroundMonitor runs background connection monitoring
func (cm *ConnectionManager) backgroundMonitor() {
	ticker := time.NewTicker(cm.pingInterval)
	defer ticker.Stop()

	for range ticker.C {
		if !cm.initialized {
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		if err := cm.Ping(ctx); err != nil {
			// Log ping failure (in a real implementation, use proper logging)
			fmt.Printf("Background ping failed: %v\n", err)
		}
		cancel()
	}
}

// OptimizeDatabase runs SQLite optimization commands
func (cm *ConnectionManager) OptimizeDatabase(ctx context.Context) error {
	optimizations := []string{
		"PRAGMA optimize",
		"PRAGMA wal_checkpoint(TRUNCATE)",
		"ANALYZE",
	}

	for _, opt := range optimizations {
		if _, err := cm.ExecuteWithTimeout(ctx, opt); err != nil {
			return fmt.Errorf("optimization failed (%s): %w", opt, err)
		}
	}

	return nil
}

// Checkpoint performs a WAL checkpoint operation
func (cm *ConnectionManager) Checkpoint(ctx context.Context, mode string) error {
	if !cm.config.EnableWAL {
		return fmt.Errorf("WAL mode not enabled")
	}

	checkpointSQL := fmt.Sprintf("PRAGMA wal_checkpoint(%s)", mode)
	_, err := cm.ExecuteWithTimeout(ctx, checkpointSQL)
	return err
}

// Close closes the connection manager and database
func (cm *ConnectionManager) Close() error {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	if !cm.initialized {
		return nil
	}

	if cm.db != nil {
		// Perform final checkpoint if WAL mode is enabled
		if cm.config.EnableWAL {
			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			cm.Checkpoint(ctx, "TRUNCATE")
			cancel()
		}

		err := cm.db.Close()
		cm.db = nil
		cm.initialized = false
		return err
	}

	return nil
}

// GetConnectionInfo returns detailed connection information
func (cm *ConnectionManager) GetConnectionInfo(ctx context.Context) (map[string]interface{}, error) {
	info := make(map[string]interface{})

	// Basic connection info
	info["initialized"] = cm.initialized
	info["database_path"] = cm.config.DatabasePath
	info["last_ping"] = cm.lastPingTime

	if !cm.initialized {
		return info, nil
	}

	// SQLite version and compile options
	var version string
	err := cm.db.QueryRowContext(ctx, "SELECT sqlite_version()").Scan(&version)
	if err == nil {
		info["sqlite_version"] = version
	}

	// Pragma information
	pragmas := map[string]string{
		"journal_mode":     "SELECT * FROM pragma_journal_mode",
		"synchronous":      "SELECT * FROM pragma_synchronous",
		"cache_size":       "SELECT * FROM pragma_cache_size",
		"mmap_size":        "SELECT * FROM pragma_mmap_size",
		"foreign_keys":     "SELECT * FROM pragma_foreign_keys",
		"temp_store":       "SELECT * FROM pragma_temp_store",
		"wal_autocheckpoint": "SELECT * FROM pragma_wal_autocheckpoint",
	}

	for name, query := range pragmas {
		var value interface{}
		if err := cm.db.QueryRowContext(ctx, query).Scan(&value); err == nil {
			info[name] = value
		}
	}

	// Connection pool stats
	info["stats"] = cm.GetStats()

	return info, nil
}