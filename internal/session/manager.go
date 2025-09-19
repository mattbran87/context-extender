package session

// NewSessionManager creates a new database-based session manager
// This is the new interface that replaces the old file-based session manager
func NewSessionManager() (*DatabaseSessionManager, error) {
	return NewDatabaseSessionManager()
}