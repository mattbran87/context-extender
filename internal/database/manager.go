package database

import (
	"context"
	"fmt"
	"log"
	"sync"
)

// Manager handles backend selection and lifecycle
type Manager struct {
	config          *DatabaseConfig
	backend         DatabaseBackend
	mu              sync.RWMutex
	registry        map[BackendType]BackendFactory
	initialized     bool
}

// NewManager creates a new database manager
func NewManager(config *DatabaseConfig) *Manager {
	if config == nil {
		config = DefaultDatabaseConfig()
	}

	m := &Manager{
		config:   config,
		registry: make(map[BackendType]BackendFactory),
	}

	// Register available backends
	m.registerBackends()

	return m
}

// registerBackends registers all available database backends
func (m *Manager) registerBackends() {
	// Register Pure Go SQLite backend
	m.registry[BackendPureGoSQLite] = &PureGoSQLiteFactory{}

	// Register CGO SQLite backend if available (legacy support)
	if m.isCGOAvailable() {
		m.registry[BackendCGOSQLite] = &CGOSQLiteFactory{}
	}
}

// isCGOAvailable checks if CGO SQLite is available
func (m *Manager) isCGOAvailable() bool {
	// For now, assume CGO SQLite is not available in pure Go mode
	// This can be enhanced with build tags or runtime detection
	return false
}

// Initialize initializes the database manager and backend
func (m *Manager) Initialize(ctx context.Context) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.initialized {
		return nil
	}

	backendType := m.config.Backend
	if backendType == BackendAuto {
		backendType = m.autoSelectBackend()
	}

	factory, exists := m.registry[backendType]
	if !exists {
		return fmt.Errorf("backend %s not available", backendType)
	}

	backend, err := factory.CreateBackend(m.config)
	if err != nil {
		return fmt.Errorf("failed to create backend: %w", err)
	}

	if err := backend.Initialize(ctx, m.config); err != nil {
		return fmt.Errorf("failed to initialize backend: %w", err)
	}

	m.backend = backend
	m.initialized = true

	log.Printf("Database manager initialized with backend: %s", backendType)
	return nil
}

// autoSelectBackend automatically selects the best available backend
func (m *Manager) autoSelectBackend() BackendType {
	// Priority order for auto-selection
	priorities := []BackendType{
		BackendPureGoSQLite, // Preferred: No CGO dependency
		BackendCGOSQLite,    // Fallback: If CGO available
	}

	for _, backend := range priorities {
		if factory, exists := m.registry[backend]; exists && factory.IsAvailable() {
			log.Printf("Auto-selected backend: %s", backend)
			return backend
		}
	}

	// Default to Pure Go SQLite
	return BackendPureGoSQLite
}

// GetBackend returns the initialized backend
func (m *Manager) GetBackend() (DatabaseBackend, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if !m.initialized {
		return nil, fmt.Errorf("database manager not initialized")
	}

	return m.backend, nil
}

// Close closes the database manager and backend
func (m *Manager) Close() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.backend != nil {
		if err := m.backend.Close(); err != nil {
			return err
		}
		m.backend = nil
	}

	m.initialized = false
	return nil
}

// GetBackendInfo returns information about the current backend
func (m *Manager) GetBackendInfo() (*BackendInfo, error) {
	backend, err := m.GetBackend()
	if err != nil {
		return nil, err
	}

	return backend.GetBackendInfo(), nil
}

// GetAvailableBackends returns a list of available backends
func (m *Manager) GetAvailableBackends() []BackendType {
	var available []BackendType

	for backendType, factory := range m.registry {
		if factory.IsAvailable() {
			available = append(available, backendType)
		}
	}

	return available
}

// SwitchBackend switches to a different backend (requires reinitialization)
func (m *Manager) SwitchBackend(ctx context.Context, newBackendType BackendType) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Close current backend
	if m.backend != nil {
		if err := m.backend.Close(); err != nil {
			return fmt.Errorf("failed to close current backend: %w", err)
		}
	}

	// Update config
	m.config.Backend = newBackendType
	m.initialized = false

	// Reinitialize with new backend
	m.mu.Unlock() // Unlock for Initialize call
	err := m.Initialize(ctx)
	m.mu.Lock()   // Re-lock

	return err
}

// Global manager instance (for backward compatibility)
var globalManager *Manager
var globalOnce sync.Once

// InitializeGlobal initializes the global database manager
func InitializeGlobal(ctx context.Context, config *DatabaseConfig) error {
	var err error
	globalOnce.Do(func() {
		globalManager = NewManager(config)
		err = globalManager.Initialize(ctx)
	})
	return err
}

// GetGlobalManager returns the global database manager
func GetGlobalManager() (*Manager, error) {
	if globalManager == nil {
		return nil, fmt.Errorf("global database manager not initialized")
	}
	return globalManager, nil
}

// GetGlobalBackend returns the global backend instance
func GetGlobalBackend() (DatabaseBackend, error) {
	if globalManager == nil {
		return nil, fmt.Errorf("global database manager not initialized")
	}
	return globalManager.GetBackend()
}

// CloseGlobal closes the global database manager
func CloseGlobal() error {
	if globalManager != nil {
		return globalManager.Close()
	}
	return nil
}

// PureGoSQLiteFactory creates pure Go SQLite backend instances
type PureGoSQLiteFactory struct{}

func (f *PureGoSQLiteFactory) CreateBackend(config *DatabaseConfig) (DatabaseBackend, error) {
	backend := NewPureGoSQLiteBackend()
	return backend, nil
}

func (f *PureGoSQLiteFactory) IsAvailable() bool {
	// Pure Go SQLite is always available
	return true
}

func (f *PureGoSQLiteFactory) GetCapabilities() *BackendCapabilities {
	return &BackendCapabilities{
		SupportsEncryption:   true,
		SupportsFullText:     true,
		SupportsTransactions: true,
		RequiresCGO:          false,
		PlatformSupport:      []string{"windows", "linux", "darwin"},
	}
}

// Legacy CGO factory (kept for compatibility)
type CGOSQLiteFactory struct{}

func (f *CGOSQLiteFactory) CreateBackend(config *DatabaseConfig) (DatabaseBackend, error) {
	// Legacy CGO implementation
	return nil, fmt.Errorf("CGO SQLite backend not available in pure Go mode")
}

func (f *CGOSQLiteFactory) IsAvailable() bool {
	return false // Not available in pure Go mode
}

func (f *CGOSQLiteFactory) GetCapabilities() *BackendCapabilities {
	return &BackendCapabilities{
		SupportsEncryption:   true, // SQLite with CGO
		SupportsFullText:     true, // SQLite FTS
		SupportsTransactions: true, // SQLite transactions
		RequiresCGO:          true, // Requires CGO
		PlatformSupport:      []string{"windows", "linux", "darwin"},
	}
}