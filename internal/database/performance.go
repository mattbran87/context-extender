package database

import (
	"sync"
	"time"
)

// BatchWriter provides buffered write operations for improved performance
type BatchWriter struct {
	events     []interface{}
	mu         sync.Mutex
	flushTimer *time.Timer
	flushSize  int
	maxDelay   time.Duration
}

// NewBatchWriter creates a new batch writer with specified parameters
func NewBatchWriter(flushSize int, maxDelay time.Duration) *BatchWriter {
	return &BatchWriter{
		events:    make([]interface{}, 0, flushSize),
		flushSize: flushSize,
		maxDelay:  maxDelay,
	}
}

// Add adds an item to the batch
func (bw *BatchWriter) Add(item interface{}) error {
	bw.mu.Lock()
	defer bw.mu.Unlock()

	bw.events = append(bw.events, item)

	// Start flush timer if this is the first item
	if len(bw.events) == 1 {
		bw.resetTimer()
	}

	// Flush if batch is full
	if len(bw.events) >= bw.flushSize {
		return bw.flush()
	}

	return nil
}

// flush writes all pending events to the database
func (bw *BatchWriter) flush() error {
	if len(bw.events) == 0 {
		return nil
	}

	// Process batch based on type
	for _, item := range bw.events {
		switch v := item.(type) {
		case *Event:
			if err := CreateEvent(v); err != nil {
				return err
			}
		case *Conversation:
			if err := CreateConversation(v); err != nil {
				return err
			}
		}
	}

	// Clear the batch
	bw.events = bw.events[:0]

	// Cancel timer if active
	if bw.flushTimer != nil {
		bw.flushTimer.Stop()
		bw.flushTimer = nil
	}

	return nil
}

// resetTimer resets the flush timer
func (bw *BatchWriter) resetTimer() {
	if bw.flushTimer != nil {
		bw.flushTimer.Stop()
	}

	bw.flushTimer = time.AfterFunc(bw.maxDelay, func() {
		bw.mu.Lock()
		defer bw.mu.Unlock()
		bw.flush()
	})
}

// Flush forces a flush of all pending events
func (bw *BatchWriter) Flush() error {
	bw.mu.Lock()
	defer bw.mu.Unlock()
	return bw.flush()
}

// QueryCache provides simple caching for frequently accessed data
type QueryCache struct {
	sessions map[string]*Session
	events   map[string][]*Event
	mu       sync.RWMutex
	ttl      time.Duration
	lastClean time.Time
}

// NewQueryCache creates a new query cache
func NewQueryCache(ttl time.Duration) *QueryCache {
	return &QueryCache{
		sessions:  make(map[string]*Session),
		events:    make(map[string][]*Event),
		ttl:       ttl,
		lastClean: time.Now(),
	}
}

// GetSession retrieves a session from cache or database
func (qc *QueryCache) GetSession(sessionID string) (*Session, error) {
	qc.mu.RLock()
	if session, ok := qc.sessions[sessionID]; ok {
		qc.mu.RUnlock()
		return session, nil
	}
	qc.mu.RUnlock()

	// Load from database
	session, err := GetSession(sessionID)
	if err != nil {
		return nil, err
	}

	if session != nil {
		qc.mu.Lock()
		qc.sessions[sessionID] = session
		qc.mu.Unlock()
	}

	// Clean old entries periodically
	qc.cleanIfNeeded()

	return session, nil
}

// InvalidateSession removes a session from cache
func (qc *QueryCache) InvalidateSession(sessionID string) {
	qc.mu.Lock()
	delete(qc.sessions, sessionID)
	delete(qc.events, sessionID)
	qc.mu.Unlock()
}

// cleanIfNeeded removes old cache entries
func (qc *QueryCache) cleanIfNeeded() {
	if time.Since(qc.lastClean) < qc.ttl {
		return
	}

	qc.mu.Lock()
	defer qc.mu.Unlock()

	// For simplicity, clear entire cache
	// In production, track access times per entry
	qc.sessions = make(map[string]*Session)
	qc.events = make(map[string][]*Event)
	qc.lastClean = time.Now()
}

// PerformanceMetrics tracks database performance
type PerformanceMetrics struct {
	mu                sync.RWMutex
	hookExecutionTime map[string]time.Duration
	queryCount        int64
	writeCount        int64
	cacheHits         int64
	cacheMisses       int64
}

var metrics = &PerformanceMetrics{
	hookExecutionTime: make(map[string]time.Duration),
}

// RecordHookExecution records hook execution time
func RecordHookExecution(hookName string, duration time.Duration) {
	metrics.mu.Lock()
	defer metrics.mu.Unlock()
	metrics.hookExecutionTime[hookName] = duration
}

// IncrementQueryCount increments the query counter
func IncrementQueryCount() {
	metrics.mu.Lock()
	defer metrics.mu.Unlock()
	metrics.queryCount++
}

// IncrementWriteCount increments the write counter
func IncrementWriteCount() {
	metrics.mu.Lock()
	defer metrics.mu.Unlock()
	metrics.writeCount++
}

// GetMetrics returns current performance metrics
func GetMetrics() map[string]interface{} {
	metrics.mu.RLock()
	defer metrics.mu.RUnlock()

	result := make(map[string]interface{})
	result["query_count"] = metrics.queryCount
	result["write_count"] = metrics.writeCount
	result["cache_hits"] = metrics.cacheHits
	result["cache_misses"] = metrics.cacheMisses

	// Calculate average hook execution times
	hookTimes := make(map[string]float64)
	for hook, duration := range metrics.hookExecutionTime {
		hookTimes[hook] = duration.Seconds()
	}
	result["hook_execution_times"] = hookTimes

	return result
}