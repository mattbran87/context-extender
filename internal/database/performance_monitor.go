package database

import (
	"fmt"
	"sync"
	"time"
)

// PerformanceMonitor tracks database performance metrics
type PerformanceMonitor struct {
	mu              sync.RWMutex
	metrics         map[string]*OperationMetrics
	globalMetrics   *GlobalMetrics
	startTime       time.Time
	sampleInterval  time.Duration
	historySize     int
}

// OperationMetrics tracks metrics for specific database operations
type OperationMetrics struct {
	OperationType    string          `json:"operation_type"`
	TotalCalls       int64           `json:"total_calls"`
	SuccessfulCalls  int64           `json:"successful_calls"`
	FailedCalls      int64           `json:"failed_calls"`
	TotalDuration    time.Duration   `json:"total_duration"`
	MinDuration      time.Duration   `json:"min_duration"`
	MaxDuration      time.Duration   `json:"max_duration"`
	AverageDuration  time.Duration   `json:"average_duration"`
	RecentDurations  []time.Duration `json:"recent_durations"`
	LastCall         time.Time       `json:"last_call"`
	ErrorRate        float64         `json:"error_rate"`
}

// GlobalMetrics tracks overall database performance
type GlobalMetrics struct {
	TotalOperations     int64         `json:"total_operations"`
	OperationsPerSecond float64       `json:"operations_per_second"`
	AverageResponseTime time.Duration `json:"average_response_time"`
	ErrorRate           float64       `json:"error_rate"`
	UpTime              time.Duration `json:"uptime"`
	PeakOpsPerSecond    float64       `json:"peak_ops_per_second"`
	LastUpdateTime      time.Time     `json:"last_update_time"`
}

// PerformanceReport contains a comprehensive performance analysis
type PerformanceReport struct {
	GeneratedAt      time.Time                    `json:"generated_at"`
	MonitoringPeriod time.Duration               `json:"monitoring_period"`
	GlobalMetrics    *GlobalMetrics              `json:"global_metrics"`
	OperationMetrics map[string]*OperationMetrics `json:"operation_metrics"`
	TopOperations    []*OperationMetrics         `json:"top_operations"`
	SlowOperations   []*OperationMetrics         `json:"slow_operations"`
	Recommendations  []string                    `json:"recommendations"`
}

// NewPerformanceMonitor creates a new performance monitor
func NewPerformanceMonitor() *PerformanceMonitor {
	return &PerformanceMonitor{
		metrics:        make(map[string]*OperationMetrics),
		globalMetrics:  &GlobalMetrics{},
		startTime:      time.Now(),
		sampleInterval: time.Second,
		historySize:    100, // Keep last 100 operation durations
	}
}

// RecordOperation records the performance of a database operation
func (pm *PerformanceMonitor) RecordOperation(operationType string, duration time.Duration, success bool) {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	// Get or create operation metrics
	metrics, exists := pm.metrics[operationType]
	if !exists {
		metrics = &OperationMetrics{
			OperationType:   operationType,
			MinDuration:     duration,
			MaxDuration:     duration,
			RecentDurations: make([]time.Duration, 0, pm.historySize),
		}
		pm.metrics[operationType] = metrics
	}

	// Update operation metrics
	metrics.TotalCalls++
	metrics.TotalDuration += duration
	metrics.LastCall = time.Now()

	if success {
		metrics.SuccessfulCalls++
	} else {
		metrics.FailedCalls++
	}

	// Update min/max durations
	if duration < metrics.MinDuration {
		metrics.MinDuration = duration
	}
	if duration > metrics.MaxDuration {
		metrics.MaxDuration = duration
	}

	// Calculate average duration
	metrics.AverageDuration = metrics.TotalDuration / time.Duration(metrics.TotalCalls)

	// Calculate error rate
	metrics.ErrorRate = float64(metrics.FailedCalls) / float64(metrics.TotalCalls)

	// Add to recent durations (sliding window)
	metrics.RecentDurations = append(metrics.RecentDurations, duration)
	if len(metrics.RecentDurations) > pm.historySize {
		metrics.RecentDurations = metrics.RecentDurations[1:]
	}

	// Update global metrics
	pm.updateGlobalMetrics()
}

// updateGlobalMetrics updates the global performance metrics
func (pm *PerformanceMonitor) updateGlobalMetrics() {
	totalOps := int64(0)
	totalDuration := time.Duration(0)
	totalErrors := int64(0)

	for _, metrics := range pm.metrics {
		totalOps += metrics.TotalCalls
		totalDuration += metrics.TotalDuration
		totalErrors += metrics.FailedCalls
	}

	pm.globalMetrics.TotalOperations = totalOps
	pm.globalMetrics.UpTime = time.Since(pm.startTime)
	pm.globalMetrics.LastUpdateTime = time.Now()

	if totalOps > 0 {
		pm.globalMetrics.AverageResponseTime = totalDuration / time.Duration(totalOps)
		pm.globalMetrics.ErrorRate = float64(totalErrors) / float64(totalOps)

		// Calculate operations per second
		uptimeSeconds := pm.globalMetrics.UpTime.Seconds()
		if uptimeSeconds > 0 {
			currentOpsPerSecond := float64(totalOps) / uptimeSeconds
			pm.globalMetrics.OperationsPerSecond = currentOpsPerSecond

			if currentOpsPerSecond > pm.globalMetrics.PeakOpsPerSecond {
				pm.globalMetrics.PeakOpsPerSecond = currentOpsPerSecond
			}
		}
	}
}

// GetMetrics returns a copy of current metrics for a specific operation
func (pm *PerformanceMonitor) GetMetrics(operationType string) *OperationMetrics {
	pm.mu.RLock()
	defer pm.mu.RUnlock()

	metrics, exists := pm.metrics[operationType]
	if !exists {
		return nil
	}

	// Return a copy to avoid race conditions
	copy := *metrics
	copy.RecentDurations = make([]time.Duration, len(metrics.RecentDurations))
	for i, d := range metrics.RecentDurations {
		copy.RecentDurations[i] = d
	}

	return &copy
}

// GetAllMetrics returns a copy of all current metrics
func (pm *PerformanceMonitor) GetAllMetrics() map[string]*OperationMetrics {
	pm.mu.RLock()
	defer pm.mu.RUnlock()

	result := make(map[string]*OperationMetrics)
	for opType := range pm.metrics {
		result[opType] = pm.GetMetrics(opType)
	}

	return result
}

// GetGlobalMetrics returns a copy of global metrics
func (pm *PerformanceMonitor) GetGlobalMetrics() *GlobalMetrics {
	pm.mu.RLock()
	defer pm.mu.RUnlock()

	copy := *pm.globalMetrics
	return &copy
}

// GenerateReport creates a comprehensive performance report
func (pm *PerformanceMonitor) GenerateReport() *PerformanceReport {
	pm.mu.RLock()
	defer pm.mu.RUnlock()

	report := &PerformanceReport{
		GeneratedAt:      time.Now(),
		MonitoringPeriod: time.Since(pm.startTime),
		GlobalMetrics:    pm.GetGlobalMetrics(),
		OperationMetrics: pm.GetAllMetrics(),
	}

	// Find top operations by volume
	topOps := make([]*OperationMetrics, 0, len(pm.metrics))
	for _, metrics := range pm.metrics {
		topOps = append(topOps, pm.GetMetrics(metrics.OperationType))
	}

	// Sort by total calls (simple sort for now)
	for i := 0; i < len(topOps)-1; i++ {
		for j := i + 1; j < len(topOps); j++ {
			if topOps[i].TotalCalls < topOps[j].TotalCalls {
				topOps[i], topOps[j] = topOps[j], topOps[i]
			}
		}
	}

	// Take top 5
	if len(topOps) > 5 {
		report.TopOperations = topOps[:5]
	} else {
		report.TopOperations = topOps
	}

	// Find slow operations (operations with high average duration)
	slowOps := make([]*OperationMetrics, 0)
	threshold := 100 * time.Millisecond // Consider operations >100ms as slow

	for _, metrics := range pm.metrics {
		if metrics.AverageDuration > threshold {
			slowOps = append(slowOps, pm.GetMetrics(metrics.OperationType))
		}
	}

	report.SlowOperations = slowOps

	// Generate recommendations
	report.Recommendations = pm.generateRecommendations(report)

	return report
}

// generateRecommendations analyzes metrics and provides optimization recommendations
func (pm *PerformanceMonitor) generateRecommendations(report *PerformanceReport) []string {
	recommendations := make([]string, 0)

	// Check error rate
	if report.GlobalMetrics.ErrorRate > 0.05 { // 5% error rate
		recommendations = append(recommendations,
			fmt.Sprintf("High error rate detected (%.2f%%). Review failed operations and add error handling.",
				report.GlobalMetrics.ErrorRate*100))
	}

	// Check response time
	if report.GlobalMetrics.AverageResponseTime > 50*time.Millisecond {
		recommendations = append(recommendations,
			fmt.Sprintf("Average response time is high (%v). Consider query optimization or indexing.",
				report.GlobalMetrics.AverageResponseTime))
	}

	// Check for slow operations
	if len(report.SlowOperations) > 0 {
		recommendations = append(recommendations,
			fmt.Sprintf("Found %d slow operations. Review query efficiency and consider optimization.",
				len(report.SlowOperations)))
	}

	// Check operations per second
	if report.GlobalMetrics.OperationsPerSecond < 10 {
		recommendations = append(recommendations,
			"Low throughput detected. Consider connection pooling optimization or hardware upgrade.")
	}

	// Check for operations with high variance
	for _, metrics := range report.OperationMetrics {
		if len(metrics.RecentDurations) > 10 {
			variance := pm.calculateVariance(metrics.RecentDurations)
			if variance > 100*time.Millisecond {
				recommendations = append(recommendations,
					fmt.Sprintf("High variance in %s operation times. Investigate performance inconsistencies.",
						metrics.OperationType))
			}
		}
	}

	if len(recommendations) == 0 {
		recommendations = append(recommendations, "Performance metrics are within acceptable ranges.")
	}

	return recommendations
}

// calculateVariance calculates the variance of operation durations
func (pm *PerformanceMonitor) calculateVariance(durations []time.Duration) time.Duration {
	if len(durations) == 0 {
		return 0
	}

	// Calculate mean
	var sum time.Duration
	for _, d := range durations {
		sum += d
	}
	mean := sum / time.Duration(len(durations))

	// Calculate variance
	var varianceSum int64
	for _, d := range durations {
		diff := d - mean
		varianceSum += int64(diff * diff)
	}

	variance := time.Duration(varianceSum / int64(len(durations)))
	return variance
}

// Reset clears all metrics and starts fresh monitoring
func (pm *PerformanceMonitor) Reset() {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	pm.metrics = make(map[string]*OperationMetrics)
	pm.globalMetrics = &GlobalMetrics{}
	pm.startTime = time.Now()
}

// GetTopOperationsByDuration returns operations sorted by average duration
func (pm *PerformanceMonitor) GetTopOperationsByDuration(limit int) []*OperationMetrics {
	pm.mu.RLock()
	defer pm.mu.RUnlock()

	operations := make([]*OperationMetrics, 0, len(pm.metrics))
	for _, metrics := range pm.metrics {
		operations = append(operations, pm.GetMetrics(metrics.OperationType))
	}

	// Sort by average duration (descending)
	for i := 0; i < len(operations)-1; i++ {
		for j := i + 1; j < len(operations); j++ {
			if operations[i].AverageDuration < operations[j].AverageDuration {
				operations[i], operations[j] = operations[j], operations[i]
			}
		}
	}

	if limit > 0 && len(operations) > limit {
		operations = operations[:limit]
	}

	return operations
}

// MonitoredExecute wraps a database operation with performance monitoring
func (pm *PerformanceMonitor) MonitoredExecute(operationType string, fn func() error) error {
	start := time.Now()
	err := fn()
	duration := time.Since(start)

	pm.RecordOperation(operationType, duration, err == nil)
	return err
}

// MonitoredExecuteWithResult wraps a database operation that returns a result
func (pm *PerformanceMonitor) MonitoredExecuteWithResult(operationType string, fn func() (interface{}, error)) (interface{}, error) {
	start := time.Now()
	result, err := fn()
	duration := time.Since(start)

	pm.RecordOperation(operationType, duration, err == nil)
	return result, err
}