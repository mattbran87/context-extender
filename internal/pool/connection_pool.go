package pool

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

// ConnectionPool provides high-performance connection pooling with monitoring
type ConnectionPool struct {
	// Pool configuration
	config *PoolConfig

	// Connection management
	connections chan *PooledConnection
	factory     ConnectionFactory
	mu          sync.RWMutex

	// Pool state
	created   int32
	active    int32
	idle      int32
	closed    int32

	// Statistics
	stats *PoolStatistics

	// Lifecycle
	ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup
}

// PooledConnection wraps a connection with metadata
type PooledConnection struct {
	conn       interface{}
	created    time.Time
	lastUsed   time.Time
	usageCount int64
	pool       *ConnectionPool
}

// PoolConfig contains pool configuration
type PoolConfig struct {
	// Size limits
	MinConnections int           `json:"min_connections"`
	MaxConnections int           `json:"max_connections"`
	MaxIdleTime    time.Duration `json:"max_idle_time"`
	MaxLifetime    time.Duration `json:"max_lifetime"`

	// Timeouts
	ConnectionTimeout time.Duration `json:"connection_timeout"`
	AcquireTimeout    time.Duration `json:"acquire_timeout"`

	// Health checks
	HealthCheckInterval time.Duration `json:"health_check_interval"`
	EnableHealthCheck   bool          `json:"enable_health_check"`

	// Monitoring
	EnableMetrics bool `json:"enable_metrics"`
}

// PoolStatistics tracks pool performance
type PoolStatistics struct {
	// Connection counts
	TotalCreated    uint64 `json:"total_created"`
	TotalDestroyed  uint64 `json:"total_destroyed"`
	CurrentActive   int32  `json:"current_active"`
	CurrentIdle     int32  `json:"current_idle"`

	// Usage statistics
	AcquireCount    uint64        `json:"acquire_count"`
	ReleaseCount    uint64        `json:"release_count"`
	TimeoutCount    uint64        `json:"timeout_count"`
	ErrorCount      uint64        `json:"error_count"`
	AcquireTime     time.Duration `json:"acquire_time"`
	AvgAcquireTime  time.Duration `json:"avg_acquire_time"`

	// Health statistics
	HealthCheckCount uint64 `json:"health_check_count"`
	HealthFailCount  uint64 `json:"health_fail_count"`

	// Timing
	StartTime    time.Time `json:"start_time"`
	LastActivity time.Time `json:"last_activity"`

	mu sync.RWMutex
}

// ConnectionFactory creates new connections
type ConnectionFactory interface {
	CreateConnection(ctx context.Context) (interface{}, error)
	ValidateConnection(conn interface{}) error
	CloseConnection(conn interface{}) error
}

// HealthChecker checks connection health
type HealthChecker interface {
	IsHealthy(conn interface{}) bool
}

// DefaultPoolConfig returns a default pool configuration
func DefaultPoolConfig() *PoolConfig {
	return &PoolConfig{
		MinConnections:      2,
		MaxConnections:      10,
		MaxIdleTime:         5 * time.Minute,
		MaxLifetime:         30 * time.Minute,
		ConnectionTimeout:   30 * time.Second,
		AcquireTimeout:      10 * time.Second,
		HealthCheckInterval: 1 * time.Minute,
		EnableHealthCheck:   true,
		EnableMetrics:       true,
	}
}

// NewConnectionPool creates a new connection pool
func NewConnectionPool(factory ConnectionFactory, config *PoolConfig) (*ConnectionPool, error) {
	if config == nil {
		config = DefaultPoolConfig()
	}

	if factory == nil {
		return nil, fmt.Errorf("connection factory cannot be nil")
	}

	ctx, cancel := context.WithCancel(context.Background())

	pool := &ConnectionPool{
		config:      config,
		connections: make(chan *PooledConnection, config.MaxConnections),
		factory:     factory,
		stats:       &PoolStatistics{StartTime: time.Now()},
		ctx:         ctx,
		cancel:      cancel,
	}

	// Pre-create minimum connections
	if err := pool.initializePool(); err != nil {
		cancel()
		return nil, fmt.Errorf("failed to initialize pool: %w", err)
	}

	// Start background workers
	if config.EnableHealthCheck {
		pool.wg.Add(1)
		go pool.healthCheckWorker()
	}

	if config.EnableMetrics {
		pool.wg.Add(1)
		go pool.metricsWorker()
	}

	pool.wg.Add(1)
	go pool.cleanupWorker()

	return pool, nil
}

// Acquire gets a connection from the pool
func (p *ConnectionPool) Acquire(ctx context.Context) (*PooledConnection, error) {
	start := time.Now()
	defer func() {
		if p.config.EnableMetrics {
			p.stats.mu.Lock()
			p.stats.AcquireCount++
			duration := time.Since(start)
			p.stats.AcquireTime += duration
			p.stats.AvgAcquireTime = time.Duration(int64(p.stats.AcquireTime) / int64(p.stats.AcquireCount))
			p.stats.LastActivity = time.Now()
			p.stats.mu.Unlock()
		}
	}()

	// Create context with timeout
	acquireCtx, cancel := context.WithTimeout(ctx, p.config.AcquireTimeout)
	defer cancel()

	select {
	case conn := <-p.connections:
		// Got a connection from pool
		conn.lastUsed = time.Now()
		atomic.AddInt64(&conn.usageCount, 1)
		atomic.AddInt32(&p.idle, -1)
		atomic.AddInt32(&p.active, 1)
		return conn, nil

	case <-acquireCtx.Done():
		// Try to create a new connection if under limit
		if atomic.LoadInt32(&p.created) < int32(p.config.MaxConnections) {
			return p.createConnection(ctx)
		}

		if p.config.EnableMetrics {
			atomic.AddUint64(&p.stats.TimeoutCount, 1)
		}
		return nil, fmt.Errorf("acquire timeout: no connections available")

	default:
		// Try to create a new connection if under limit
		if atomic.LoadInt32(&p.created) < int32(p.config.MaxConnections) {
			return p.createConnection(ctx)
		}

		// Wait for a connection to become available
		select {
		case conn := <-p.connections:
			conn.lastUsed = time.Now()
			atomic.AddInt64(&conn.usageCount, 1)
			atomic.AddInt32(&p.idle, -1)
			atomic.AddInt32(&p.active, 1)
			return conn, nil

		case <-acquireCtx.Done():
			if p.config.EnableMetrics {
				atomic.AddUint64(&p.stats.TimeoutCount, 1)
			}
			return nil, fmt.Errorf("acquire timeout: no connections available")
		}
	}
}

// Release returns a connection to the pool
func (p *ConnectionPool) Release(conn *PooledConnection) error {
	if conn == nil {
		return fmt.Errorf("cannot release nil connection")
	}

	// Check if connection is still valid
	if p.shouldDiscardConnection(conn) {
		return p.discardConnection(conn)
	}

	// Return to pool
	select {
	case p.connections <- conn:
		atomic.AddInt32(&p.active, -1)
		atomic.AddInt32(&p.idle, 1)
		if p.config.EnableMetrics {
			atomic.AddUint64(&p.stats.ReleaseCount, 1)
		}
		return nil
	default:
		// Pool is full, discard the connection
		return p.discardConnection(conn)
	}
}

// GetConnection returns the underlying connection
func (pc *PooledConnection) GetConnection() interface{} {
	return pc.conn
}

// GetStats returns current pool statistics
func (p *ConnectionPool) GetStats() *PoolStatistics {
	p.stats.mu.RLock()
	defer p.stats.mu.RUnlock()

	stats := *p.stats
	stats.CurrentActive = atomic.LoadInt32(&p.active)
	stats.CurrentIdle = atomic.LoadInt32(&p.idle)

	return &stats
}

// Close shuts down the pool and closes all connections
func (p *ConnectionPool) Close() error {
	p.cancel()
	p.wg.Wait()

	// Close all connections in the pool
	close(p.connections)
	for conn := range p.connections {
		p.factory.CloseConnection(conn.conn)
	}

	return nil
}

// Internal methods

func (p *ConnectionPool) initializePool() error {
	for i := 0; i < p.config.MinConnections; i++ {
		conn, err := p.createConnection(context.Background())
		if err != nil {
			return fmt.Errorf("failed to create initial connection %d: %w", i+1, err)
		}

		// Put in pool
		p.connections <- conn
		atomic.AddInt32(&p.active, -1)
		atomic.AddInt32(&p.idle, 1)
	}

	return nil
}

func (p *ConnectionPool) createConnection(ctx context.Context) (*PooledConnection, error) {
	createCtx, cancel := context.WithTimeout(ctx, p.config.ConnectionTimeout)
	defer cancel()

	conn, err := p.factory.CreateConnection(createCtx)
	if err != nil {
		if p.config.EnableMetrics {
			atomic.AddUint64(&p.stats.ErrorCount, 1)
		}
		return nil, fmt.Errorf("failed to create connection: %w", err)
	}

	pooledConn := &PooledConnection{
		conn:       conn,
		created:    time.Now(),
		lastUsed:   time.Now(),
		usageCount: 0,
		pool:       p,
	}

	atomic.AddInt32(&p.created, 1)
	atomic.AddInt32(&p.active, 1)

	if p.config.EnableMetrics {
		atomic.AddUint64(&p.stats.TotalCreated, 1)
	}

	return pooledConn, nil
}

func (p *ConnectionPool) shouldDiscardConnection(conn *PooledConnection) bool {
	now := time.Now()

	// Check max lifetime
	if p.config.MaxLifetime > 0 && now.Sub(conn.created) > p.config.MaxLifetime {
		return true
	}

	// Check max idle time
	if p.config.MaxIdleTime > 0 && now.Sub(conn.lastUsed) > p.config.MaxIdleTime {
		return true
	}

	// Validate connection
	if err := p.factory.ValidateConnection(conn.conn); err != nil {
		return true
	}

	return false
}

func (p *ConnectionPool) discardConnection(conn *PooledConnection) error {
	atomic.AddInt32(&p.active, -1)
	atomic.AddInt32(&p.created, -1)

	if p.config.EnableMetrics {
		atomic.AddUint64(&p.stats.TotalDestroyed, 1)
	}

	return p.factory.CloseConnection(conn.conn)
}

func (p *ConnectionPool) healthCheckWorker() {
	defer p.wg.Done()

	ticker := time.NewTicker(p.config.HealthCheckInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			p.performHealthCheck()
		case <-p.ctx.Done():
			return
		}
	}
}

func (p *ConnectionPool) performHealthCheck() {
	// Check all idle connections
	idleConnections := make([]*PooledConnection, 0)

	// Drain idle connections for health check
	for {
		select {
		case conn := <-p.connections:
			idleConnections = append(idleConnections, conn)
		default:
			goto checkConnections
		}
	}

checkConnections:
	healthyConnections := make([]*PooledConnection, 0)

	for _, conn := range idleConnections {
		if p.config.EnableMetrics {
			atomic.AddUint64(&p.stats.HealthCheckCount, 1)
		}

		// Check if connection is healthy
		if checker, ok := p.factory.(HealthChecker); ok {
			if !checker.IsHealthy(conn.conn) {
				// Connection is unhealthy, discard it
				p.discardConnection(conn)
				if p.config.EnableMetrics {
					atomic.AddUint64(&p.stats.HealthFailCount, 1)
				}
				continue
			}
		} else {
			// Use validation as health check
			if err := p.factory.ValidateConnection(conn.conn); err != nil {
				p.discardConnection(conn)
				if p.config.EnableMetrics {
					atomic.AddUint64(&p.stats.HealthFailCount, 1)
				}
				continue
			}
		}

		healthyConnections = append(healthyConnections, conn)
	}

	// Return healthy connections to pool
	for _, conn := range healthyConnections {
		select {
		case p.connections <- conn:
		default:
			// Pool is full, discard excess connections
			p.discardConnection(conn)
		}
	}

	// Ensure minimum connections
	p.ensureMinimumConnections()
}

func (p *ConnectionPool) ensureMinimumConnections() {
	current := atomic.LoadInt32(&p.created)
	if int(current) < p.config.MinConnections {
		needed := p.config.MinConnections - int(current)
		for i := 0; i < needed; i++ {
			conn, err := p.createConnection(context.Background())
			if err != nil {
				// Log error but continue
				continue
			}

			select {
			case p.connections <- conn:
				atomic.AddInt32(&p.active, -1)
				atomic.AddInt32(&p.idle, 1)
			default:
				// Pool is full
				p.discardConnection(conn)
				break
			}
		}
	}
}

func (p *ConnectionPool) cleanupWorker() {
	defer p.wg.Done()

	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			p.cleanupExpiredConnections()
		case <-p.ctx.Done():
			return
		}
	}
}

func (p *ConnectionPool) cleanupExpiredConnections() {
	// Similar to health check but focuses on expired connections
	idleConnections := make([]*PooledConnection, 0)

	// Drain pool
	for {
		select {
		case conn := <-p.connections:
			idleConnections = append(idleConnections, conn)
		default:
			goto cleanup
		}
	}

cleanup:
	validConnections := make([]*PooledConnection, 0)

	for _, conn := range idleConnections {
		if p.shouldDiscardConnection(conn) {
			p.discardConnection(conn)
		} else {
			validConnections = append(validConnections, conn)
		}
	}

	// Return valid connections
	for _, conn := range validConnections {
		select {
		case p.connections <- conn:
		default:
			p.discardConnection(conn)
		}
	}
}

func (p *ConnectionPool) metricsWorker() {
	defer p.wg.Done()

	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			p.updateMetrics()
		case <-p.ctx.Done():
			return
		}
	}
}

func (p *ConnectionPool) updateMetrics() {
	p.stats.mu.Lock()
	defer p.stats.mu.Unlock()

	p.stats.CurrentActive = atomic.LoadInt32(&p.active)
	p.stats.CurrentIdle = atomic.LoadInt32(&p.idle)
}

// Utility methods for monitoring and debugging

// GetUtilization returns pool utilization percentage
func (p *ConnectionPool) GetUtilization() float64 {
	active := atomic.LoadInt32(&p.active)
	return float64(active) / float64(p.config.MaxConnections) * 100
}

// GetEfficiency returns pool efficiency metrics
func (p *ConnectionPool) GetEfficiency() *EfficiencyMetrics {
	stats := p.GetStats()

	var hitRate float64
	if stats.AcquireCount > 0 {
		hitRate = float64(stats.AcquireCount-stats.TimeoutCount) / float64(stats.AcquireCount) * 100
	}

	var healthRate float64
	if stats.HealthCheckCount > 0 {
		healthRate = float64(stats.HealthCheckCount-stats.HealthFailCount) / float64(stats.HealthCheckCount) * 100
	}

	return &EfficiencyMetrics{
		HitRate:            hitRate,
		HealthRate:         healthRate,
		Utilization:        p.GetUtilization(),
		AvgAcquireTime:     stats.AvgAcquireTime,
		ConnectionTurnover: float64(stats.TotalDestroyed) / float64(stats.TotalCreated) * 100,
	}
}

// EfficiencyMetrics contains pool efficiency information
type EfficiencyMetrics struct {
	HitRate            float64       `json:"hit_rate"`
	HealthRate         float64       `json:"health_rate"`
	Utilization        float64       `json:"utilization"`
	AvgAcquireTime     time.Duration `json:"avg_acquire_time"`
	ConnectionTurnover float64       `json:"connection_turnover"`
}

// Resize changes the pool size dynamically
func (p *ConnectionPool) Resize(newMin, newMax int) error {
	if newMin < 0 || newMax < newMin {
		return fmt.Errorf("invalid pool size: min=%d, max=%d", newMin, newMax)
	}

	p.mu.Lock()
	defer p.mu.Unlock()

	oldMin := p.config.MinConnections
	oldMax := p.config.MaxConnections

	p.config.MinConnections = newMin
	p.config.MaxConnections = newMax

	// If we're reducing max connections, close excess
	if newMax < oldMax {
		current := atomic.LoadInt32(&p.created)
		if int(current) > newMax {
			excess := int(current) - newMax
			for i := 0; i < excess; i++ {
				select {
				case conn := <-p.connections:
					p.discardConnection(conn)
				default:
					break
				}
			}
		}
	}

	// If we're increasing min connections, create more
	if newMin > oldMin {
		p.ensureMinimumConnections()
	}

	return nil
}