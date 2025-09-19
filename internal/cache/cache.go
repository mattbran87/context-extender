package cache

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"
)

// Cache provides high-performance caching with TTL and statistics
type Cache struct {
	items map[string]*Item
	mu    sync.RWMutex

	// Configuration
	config *Config

	// Statistics
	stats *Statistics

	// Cleanup
	stopCleanup chan bool
}

// Item represents a cached item
type Item struct {
	Key        string      `json:"key"`
	Value      interface{} `json:"value"`
	ExpiresAt  time.Time   `json:"expires_at"`
	CreatedAt  time.Time   `json:"created_at"`
	AccessedAt time.Time   `json:"accessed_at"`
	AccessCount uint64     `json:"access_count"`
	Size       int64       `json:"size"`
}

// Config contains cache configuration
type Config struct {
	DefaultTTL      time.Duration `json:"default_ttl"`
	MaxSize         int64         `json:"max_size"`
	MaxItems        int           `json:"max_items"`
	CleanupInterval time.Duration `json:"cleanup_interval"`
	EnableStats     bool          `json:"enable_stats"`
	EvictionPolicy  EvictionPolicy `json:"eviction_policy"`
}

// EvictionPolicy defines how items are evicted when cache is full
type EvictionPolicy string

const (
	EvictionLRU   EvictionPolicy = "lru"   // Least Recently Used
	EvictionLFU   EvictionPolicy = "lfu"   // Least Frequently Used
	EvictionFIFO  EvictionPolicy = "fifo"  // First In, First Out
	EvictionRandom EvictionPolicy = "random" // Random eviction
)

// Statistics tracks cache performance metrics
type Statistics struct {
	Hits           uint64    `json:"hits"`
	Misses         uint64    `json:"misses"`
	Sets           uint64    `json:"sets"`
	Deletes        uint64    `json:"deletes"`
	Evictions      uint64    `json:"evictions"`
	CurrentItems   int       `json:"current_items"`
	CurrentSize    int64     `json:"current_size"`
	HitRate        float64   `json:"hit_rate"`
	LastCleanup    time.Time `json:"last_cleanup"`
	StartTime      time.Time `json:"start_time"`

	mu sync.RWMutex
}

// DefaultConfig returns a default cache configuration
func DefaultConfig() *Config {
	return &Config{
		DefaultTTL:      5 * time.Minute,
		MaxSize:         100 * 1024 * 1024, // 100MB
		MaxItems:        10000,
		CleanupInterval: 1 * time.Minute,
		EnableStats:     true,
		EvictionPolicy:  EvictionLRU,
	}
}

// New creates a new cache instance
func New(config *Config) *Cache {
	if config == nil {
		config = DefaultConfig()
	}

	c := &Cache{
		items:       make(map[string]*Item),
		config:      config,
		stats:       &Statistics{StartTime: time.Now()},
		stopCleanup: make(chan bool),
	}

	// Start cleanup goroutine
	go c.cleanupWorker()

	return c
}

// Set stores a value in the cache with the default TTL
func (c *Cache) Set(key string, value interface{}) error {
	return c.SetWithTTL(key, value, c.config.DefaultTTL)
}

// SetWithTTL stores a value in the cache with a specific TTL
func (c *Cache) SetWithTTL(key string, value interface{}, ttl time.Duration) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Calculate item size
	size := c.calculateSize(value)

	// Check if we need to make room
	if err := c.makeRoom(size); err != nil {
		return fmt.Errorf("failed to make room in cache: %w", err)
	}

	now := time.Now()
	item := &Item{
		Key:        key,
		Value:      value,
		ExpiresAt:  now.Add(ttl),
		CreatedAt:  now,
		AccessedAt: now,
		AccessCount: 0,
		Size:       size,
	}

	// Remove existing item if present
	if existing, exists := c.items[key]; exists {
		c.stats.CurrentSize -= existing.Size
		c.stats.CurrentItems--
	}

	c.items[key] = item
	c.stats.CurrentSize += size
	c.stats.CurrentItems++

	if c.config.EnableStats {
		c.stats.mu.Lock()
		c.stats.Sets++
		c.stats.mu.Unlock()
	}

	return nil
}

// Get retrieves a value from the cache
func (c *Cache) Get(key string) (interface{}, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	item, exists := c.items[key]
	if !exists {
		if c.config.EnableStats {
			c.stats.mu.Lock()
			c.stats.Misses++
			c.updateHitRate()
			c.stats.mu.Unlock()
		}
		return nil, false
	}

	// Check if expired
	if time.Now().After(item.ExpiresAt) {
		delete(c.items, key)
		c.stats.CurrentSize -= item.Size
		c.stats.CurrentItems--
		if c.config.EnableStats {
			c.stats.mu.Lock()
			c.stats.Misses++
			c.updateHitRate()
			c.stats.mu.Unlock()
		}
		return nil, false
	}

	// Update access information
	item.AccessedAt = time.Now()
	item.AccessCount++

	if c.config.EnableStats {
		c.stats.mu.Lock()
		c.stats.Hits++
		c.updateHitRate()
		c.stats.mu.Unlock()
	}

	return item.Value, true
}

// GetWithInfo retrieves a value and its metadata from the cache
func (c *Cache) GetWithInfo(key string) (*Item, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	item, exists := c.items[key]
	if !exists || time.Now().After(item.ExpiresAt) {
		return nil, false
	}

	// Return a copy to avoid external modification
	return &Item{
		Key:        item.Key,
		Value:      item.Value,
		ExpiresAt:  item.ExpiresAt,
		CreatedAt:  item.CreatedAt,
		AccessedAt: item.AccessedAt,
		AccessCount: item.AccessCount,
		Size:       item.Size,
	}, true
}

// Delete removes an item from the cache
func (c *Cache) Delete(key string) bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	item, exists := c.items[key]
	if !exists {
		return false
	}

	delete(c.items, key)
	c.stats.CurrentSize -= item.Size
	c.stats.CurrentItems--

	if c.config.EnableStats {
		c.stats.mu.Lock()
		c.stats.Deletes++
		c.stats.mu.Unlock()
	}

	return true
}

// Clear removes all items from the cache
func (c *Cache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.items = make(map[string]*Item)
	c.stats.CurrentSize = 0
	c.stats.CurrentItems = 0
}

// Keys returns all cache keys
func (c *Cache) Keys() []string {
	c.mu.RLock()
	defer c.mu.RUnlock()

	keys := make([]string, 0, len(c.items))
	for key := range c.items {
		keys = append(keys, key)
	}
	return keys
}

// Size returns the current cache size in bytes
func (c *Cache) Size() int64 {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.stats.CurrentSize
}

// ItemCount returns the current number of items
func (c *Cache) ItemCount() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.stats.CurrentItems
}

// GetStats returns current cache statistics
func (c *Cache) GetStats() *Statistics {
	c.stats.mu.RLock()
	defer c.stats.mu.RUnlock()

	// Return a copy
	return &Statistics{
		Hits:         c.stats.Hits,
		Misses:       c.stats.Misses,
		Sets:         c.stats.Sets,
		Deletes:      c.stats.Deletes,
		Evictions:    c.stats.Evictions,
		CurrentItems: c.stats.CurrentItems,
		CurrentSize:  c.stats.CurrentSize,
		HitRate:      c.stats.HitRate,
		LastCleanup:  c.stats.LastCleanup,
		StartTime:    c.stats.StartTime,
	}
}

// Close shuts down the cache and cleanup workers
func (c *Cache) Close() {
	close(c.stopCleanup)
}

// Internal methods

func (c *Cache) makeRoom(newItemSize int64) error {
	// Check size limit
	if c.stats.CurrentSize+newItemSize > c.config.MaxSize {
		if err := c.evictBySize(newItemSize); err != nil {
			return err
		}
	}

	// Check item count limit
	if c.stats.CurrentItems >= c.config.MaxItems {
		if err := c.evictByCount(1); err != nil {
			return err
		}
	}

	return nil
}

func (c *Cache) evictBySize(requiredSpace int64) error {
	spaceToFree := c.stats.CurrentSize + requiredSpace - c.config.MaxSize

	switch c.config.EvictionPolicy {
	case EvictionLRU:
		return c.evictLRU(spaceToFree)
	case EvictionLFU:
		return c.evictLFU(spaceToFree)
	case EvictionFIFO:
		return c.evictFIFO(spaceToFree)
	case EvictionRandom:
		return c.evictRandom(spaceToFree)
	default:
		return c.evictLRU(spaceToFree)
	}
}

func (c *Cache) evictByCount(itemsToEvict int) error {
	switch c.config.EvictionPolicy {
	case EvictionLRU:
		return c.evictLRUCount(itemsToEvict)
	case EvictionLFU:
		return c.evictLFUCount(itemsToEvict)
	case EvictionFIFO:
		return c.evictFIFOCount(itemsToEvict)
	case EvictionRandom:
		return c.evictRandomCount(itemsToEvict)
	default:
		return c.evictLRUCount(itemsToEvict)
	}
}

func (c *Cache) evictLRU(spaceToFree int64) error {
	// Find items to evict based on least recently used
	items := make([]*Item, 0, len(c.items))
	for _, item := range c.items {
		items = append(items, item)
	}

	// Sort by access time (oldest first)
	for i := 0; i < len(items)-1; i++ {
		for j := i+1; j < len(items); j++ {
			if items[i].AccessedAt.After(items[j].AccessedAt) {
				items[i], items[j] = items[j], items[i]
			}
		}
	}

	freedSpace := int64(0)
	evicted := 0
	for _, item := range items {
		if freedSpace >= spaceToFree {
			break
		}

		delete(c.items, item.Key)
		freedSpace += item.Size
		c.stats.CurrentSize -= item.Size
		c.stats.CurrentItems--
		evicted++
	}

	if c.config.EnableStats {
		c.stats.mu.Lock()
		c.stats.Evictions += uint64(evicted)
		c.stats.mu.Unlock()
	}

	return nil
}

func (c *Cache) evictLFU(spaceToFree int64) error {
	// Find items to evict based on least frequently used
	items := make([]*Item, 0, len(c.items))
	for _, item := range c.items {
		items = append(items, item)
	}

	// Sort by access count (lowest first)
	for i := 0; i < len(items)-1; i++ {
		for j := i+1; j < len(items); j++ {
			if items[i].AccessCount > items[j].AccessCount {
				items[i], items[j] = items[j], items[i]
			}
		}
	}

	freedSpace := int64(0)
	evicted := 0
	for _, item := range items {
		if freedSpace >= spaceToFree {
			break
		}

		delete(c.items, item.Key)
		freedSpace += item.Size
		c.stats.CurrentSize -= item.Size
		c.stats.CurrentItems--
		evicted++
	}

	if c.config.EnableStats {
		c.stats.mu.Lock()
		c.stats.Evictions += uint64(evicted)
		c.stats.mu.Unlock()
	}

	return nil
}

func (c *Cache) evictFIFO(spaceToFree int64) error {
	// Find items to evict based on creation time (oldest first)
	items := make([]*Item, 0, len(c.items))
	for _, item := range c.items {
		items = append(items, item)
	}

	// Sort by creation time (oldest first)
	for i := 0; i < len(items)-1; i++ {
		for j := i+1; j < len(items); j++ {
			if items[i].CreatedAt.After(items[j].CreatedAt) {
				items[i], items[j] = items[j], items[i]
			}
		}
	}

	freedSpace := int64(0)
	evicted := 0
	for _, item := range items {
		if freedSpace >= spaceToFree {
			break
		}

		delete(c.items, item.Key)
		freedSpace += item.Size
		c.stats.CurrentSize -= item.Size
		c.stats.CurrentItems--
		evicted++
	}

	if c.config.EnableStats {
		c.stats.mu.Lock()
		c.stats.Evictions += uint64(evicted)
		c.stats.mu.Unlock()
	}

	return nil
}

func (c *Cache) evictRandom(spaceToFree int64) error {
	freedSpace := int64(0)
	evicted := 0

	for key, item := range c.items {
		if freedSpace >= spaceToFree {
			break
		}

		delete(c.items, key)
		freedSpace += item.Size
		c.stats.CurrentSize -= item.Size
		c.stats.CurrentItems--
		evicted++
	}

	if c.config.EnableStats {
		c.stats.mu.Lock()
		c.stats.Evictions += uint64(evicted)
		c.stats.mu.Unlock()
	}

	return nil
}

func (c *Cache) evictLRUCount(count int) error {
	// Similar to evictLRU but evicts a specific number of items
	items := make([]*Item, 0, len(c.items))
	for _, item := range c.items {
		items = append(items, item)
	}

	// Sort by access time (oldest first)
	for i := 0; i < len(items)-1; i++ {
		for j := i+1; j < len(items); j++ {
			if items[i].AccessedAt.After(items[j].AccessedAt) {
				items[i], items[j] = items[j], items[i]
			}
		}
	}

	evicted := 0
	for _, item := range items {
		if evicted >= count {
			break
		}

		delete(c.items, item.Key)
		c.stats.CurrentSize -= item.Size
		c.stats.CurrentItems--
		evicted++
	}

	if c.config.EnableStats {
		c.stats.mu.Lock()
		c.stats.Evictions += uint64(evicted)
		c.stats.mu.Unlock()
	}

	return nil
}

func (c *Cache) evictLFUCount(count int) error {
	return c.evictLRUCount(count) // Simplified for demo
}

func (c *Cache) evictFIFOCount(count int) error {
	return c.evictLRUCount(count) // Simplified for demo
}

func (c *Cache) evictRandomCount(count int) error {
	evicted := 0
	for key, item := range c.items {
		if evicted >= count {
			break
		}

		delete(c.items, key)
		c.stats.CurrentSize -= item.Size
		c.stats.CurrentItems--
		evicted++
	}

	if c.config.EnableStats {
		c.stats.mu.Lock()
		c.stats.Evictions += uint64(evicted)
		c.stats.mu.Unlock()
	}

	return nil
}

func (c *Cache) calculateSize(value interface{}) int64 {
	// Simple size calculation - in production, use more sophisticated methods
	if value == nil {
		return 8 // pointer size
	}

	switch v := value.(type) {
	case string:
		return int64(len(v))
	case []byte:
		return int64(len(v))
	case int, int32, int64, uint, uint32, uint64, float32, float64:
		return 8
	case bool:
		return 1
	default:
		// For complex objects, try JSON serialization to estimate size
		if data, err := json.Marshal(value); err == nil {
			return int64(len(data))
		}
		return 64 // Default estimate
	}
}

func (c *Cache) updateHitRate() {
	total := c.stats.Hits + c.stats.Misses
	if total > 0 {
		c.stats.HitRate = float64(c.stats.Hits) / float64(total) * 100
	}
}

func (c *Cache) cleanupWorker() {
	ticker := time.NewTicker(c.config.CleanupInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			c.cleanup()
		case <-c.stopCleanup:
			return
		}
	}
}

func (c *Cache) cleanup() {
	c.mu.Lock()
	defer c.mu.Unlock()

	now := time.Now()
	expired := make([]string, 0)

	for key, item := range c.items {
		if now.After(item.ExpiresAt) {
			expired = append(expired, key)
		}
	}

	for _, key := range expired {
		item := c.items[key]
		delete(c.items, key)
		c.stats.CurrentSize -= item.Size
		c.stats.CurrentItems--
	}

	c.stats.mu.Lock()
	c.stats.LastCleanup = now
	c.stats.mu.Unlock()
}

// Utility functions

// GetOrSet retrieves a value or sets it if not present
func (c *Cache) GetOrSet(key string, factory func() (interface{}, error)) (interface{}, error) {
	// Try to get first
	if value, exists := c.Get(key); exists {
		return value, nil
	}

	// Create the value
	value, err := factory()
	if err != nil {
		return nil, err
	}

	// Set in cache
	if err := c.Set(key, value); err != nil {
		return value, err // Return the value even if caching failed
	}

	return value, nil
}

// GetOrSetWithTTL retrieves a value or sets it with specific TTL if not present
func (c *Cache) GetOrSetWithTTL(key string, ttl time.Duration, factory func() (interface{}, error)) (interface{}, error) {
	// Try to get first
	if value, exists := c.Get(key); exists {
		return value, nil
	}

	// Create the value
	value, err := factory()
	if err != nil {
		return nil, err
	}

	// Set in cache with TTL
	if err := c.SetWithTTL(key, value, ttl); err != nil {
		return value, err // Return the value even if caching failed
	}

	return value, nil
}

// Exists checks if a key exists in the cache
func (c *Cache) Exists(key string) bool {
	_, exists := c.Get(key)
	return exists
}

// Touch updates the access time of an item without retrieving it
func (c *Cache) Touch(key string) bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	item, exists := c.items[key]
	if !exists || time.Now().After(item.ExpiresAt) {
		return false
	}

	item.AccessedAt = time.Now()
	item.AccessCount++
	return true
}

// Extend extends the TTL of an existing item
func (c *Cache) Extend(key string, additionalTTL time.Duration) bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	item, exists := c.items[key]
	if !exists || time.Now().After(item.ExpiresAt) {
		return false
	}

	item.ExpiresAt = item.ExpiresAt.Add(additionalTTL)
	return true
}