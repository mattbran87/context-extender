package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"context-extender/internal/audit"
	"context-extender/internal/auth"
	"context-extender/internal/metrics"
	"context-extender/internal/middleware"
	"context-extender/internal/performance"
	"context-extender/internal/session"
)

// Server represents the HTTP API server
type Server struct {
	httpServer       *http.Server
	config          *ServerConfig
	security        *middleware.SecurityMiddleware
	sessionManager  *session.AdvancedSessionManager
	metricsCollector *metrics.Collector
	performanceMonitor *performance.Monitor
	auditLogger     *audit.AuditLogger
	rbac           *auth.RBAC
	jwtManager     *auth.JWTManager
}

// ServerConfig contains server configuration
type ServerConfig struct {
	Address         string        `json:"address"`
	Port            int           `json:"port"`
	ReadTimeout     time.Duration `json:"read_timeout"`
	WriteTimeout    time.Duration `json:"write_timeout"`
	IdleTimeout     time.Duration `json:"idle_timeout"`
	MaxHeaderBytes  int           `json:"max_header_bytes"`
	EnableTLS       bool          `json:"enable_tls"`
	TLSCertFile     string        `json:"tls_cert_file"`
	TLSKeyFile      string        `json:"tls_key_file"`
	EnableCORS      bool          `json:"enable_cors"`
	CORSOrigins     []string      `json:"cors_origins"`
	EnableMetrics   bool          `json:"enable_metrics"`
	MetricsPath     string        `json:"metrics_path"`
	EnablePprof     bool          `json:"enable_pprof"`
	PprofPath       string        `json:"pprof_path"`
	APIPrefix       string        `json:"api_prefix"`
	EnableDocs      bool          `json:"enable_docs"`
	DocsPath        string        `json:"docs_path"`
}

// APIResponse represents a standard API response
type APIResponse struct {
	Success   bool        `json:"success"`
	Data      interface{} `json:"data,omitempty"`
	Error     *APIError   `json:"error,omitempty"`
	Meta      *Meta       `json:"meta,omitempty"`
	Timestamp time.Time   `json:"timestamp"`
	RequestID string      `json:"request_id,omitempty"`
}

// APIError represents an API error
type APIError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
	Type    string `json:"type"`
}

// Meta contains response metadata
type Meta struct {
	Version    string `json:"version"`
	Pagination *Pagination `json:"pagination,omitempty"`
	RateLimit  *RateLimit  `json:"rate_limit,omitempty"`
}

// Pagination contains pagination information
type Pagination struct {
	Page       int `json:"page"`
	PerPage    int `json:"per_page"`
	Total      int `json:"total"`
	TotalPages int `json:"total_pages"`
}

// RateLimit contains rate limiting information
type RateLimit struct {
	Limit     int   `json:"limit"`
	Remaining int   `json:"remaining"`
	ResetAt   int64 `json:"reset_at"`
}

// DefaultServerConfig returns default server configuration
func DefaultServerConfig() *ServerConfig {
	return &ServerConfig{
		Address:        "0.0.0.0",
		Port:           8080,
		ReadTimeout:    15 * time.Second,
		WriteTimeout:   15 * time.Second,
		IdleTimeout:    60 * time.Second,
		MaxHeaderBytes: 1 << 20, // 1MB
		EnableTLS:      false,
		EnableCORS:     true,
		CORSOrigins:    []string{"*"},
		EnableMetrics:  true,
		MetricsPath:    "/metrics",
		EnablePprof:    false,
		PprofPath:      "/debug/pprof",
		APIPrefix:      "/api/v1",
		EnableDocs:     true,
		DocsPath:       "/docs",
	}
}

// NewServer creates a new API server
func NewServer(config *ServerConfig) *Server {
	if config == nil {
		config = DefaultServerConfig()
	}

	server := &Server{
		config: config,
	}

	// Initialize components
	server.initializeComponents()

	// Setup HTTP server
	server.setupHTTPServer()

	return server
}

// Start starts the HTTP server
func (s *Server) Start() error {
	addr := fmt.Sprintf("%s:%d", s.config.Address, s.config.Port)

	fmt.Printf("🚀 Starting API server on %s\n", addr)
	fmt.Printf("📖 API Documentation: http://%s%s\n", addr, s.config.DocsPath)
	fmt.Printf("📊 Metrics endpoint: http://%s%s\n", addr, s.config.MetricsPath)

	if s.config.EnableTLS {
		return s.httpServer.ListenAndServeTLS(s.config.TLSCertFile, s.config.TLSKeyFile)
	}

	return s.httpServer.ListenAndServe()
}

// Stop gracefully stops the HTTP server
func (s *Server) Stop(ctx context.Context) error {
	fmt.Println("🛑 Shutting down API server...")

	// Stop metrics collection
	if s.metricsCollector != nil {
		s.metricsCollector.Stop()
	}

	// Close performance monitor
	if s.performanceMonitor != nil {
		s.performanceMonitor.Close()
	}

	return s.httpServer.Shutdown(ctx)
}

// Private methods

func (s *Server) initializeComponents() {
	// Initialize JWT manager
	s.jwtManager = auth.NewJWTManager(auth.DefaultJWTConfig())

	// Initialize RBAC
	s.rbac = auth.NewRBAC()

	// Initialize audit logger
	s.auditLogger = audit.NewAuditLogger(
		audit.DefaultAuditConfig(),
		audit.NewInMemoryStorage(),
	)

	// Initialize performance monitor
	s.performanceMonitor = performance.NewMonitor(performance.DefaultConfig())

	// Initialize session manager with encryption
	s.sessionManager = session.NewAdvancedSessionManager(session.DefaultAdvancedConfig())

	// Initialize metrics collector
	s.metricsCollector = metrics.NewCollector(metrics.DefaultCollectorConfig())
	s.metricsCollector.RegisterMonitor(s.performanceMonitor)

	// Initialize security middleware
	securityConfig := middleware.DefaultSecurityConfig()
	s.security = middleware.NewSecurityMiddleware(securityConfig)

	fmt.Println("✅ Initialized all API server components")
}

func (s *Server) setupHTTPServer() {
	mux := http.NewServeMux()

	// Setup routes
	s.setupRoutes(mux)

	// Setup middleware chain
	handler := s.setupMiddleware(mux)

	s.httpServer = &http.Server{
		Addr:           fmt.Sprintf("%s:%d", s.config.Address, s.config.Port),
		Handler:        handler,
		ReadTimeout:    s.config.ReadTimeout,
		WriteTimeout:   s.config.WriteTimeout,
		IdleTimeout:    s.config.IdleTimeout,
		MaxHeaderBytes: s.config.MaxHeaderBytes,
	}
}

func (s *Server) setupRoutes(mux *http.ServeMux) {
	// Health check endpoint
	mux.HandleFunc("/health", s.handleHealth)

	// API routes
	apiPrefix := s.config.APIPrefix

	// Authentication routes
	mux.HandleFunc(apiPrefix+"/auth/login", s.handleLogin)
	mux.HandleFunc(apiPrefix+"/auth/refresh", s.handleRefresh)
	mux.HandleFunc(apiPrefix+"/auth/logout", s.handleLogout)

	// Session routes (protected)
	mux.HandleFunc(apiPrefix+"/sessions", s.protectedHandler(s.handleSessions, "session", "read"))
	mux.HandleFunc(apiPrefix+"/sessions/", s.protectedHandler(s.handleSessionDetail, "session", "read"))

	// User routes (protected)
	mux.HandleFunc(apiPrefix+"/users", s.protectedHandler(s.handleUsers, "users", "read"))
	mux.HandleFunc(apiPrefix+"/users/", s.protectedHandler(s.handleUserDetail, "users", "read"))

	// Metrics routes (protected)
	mux.HandleFunc(apiPrefix+"/metrics", s.protectedHandler(s.handleMetrics, "metrics", "read"))
	mux.HandleFunc(apiPrefix+"/metrics/export", s.protectedHandler(s.handleMetricsExport, "metrics", "read"))

	// Performance routes (protected)
	mux.HandleFunc(apiPrefix+"/performance", s.protectedHandler(s.handlePerformance, "performance", "read"))
	mux.HandleFunc(apiPrefix+"/performance/report", s.protectedHandler(s.handlePerformanceReport, "performance", "read"))

	// Admin routes (admin only)
	mux.HandleFunc(apiPrefix+"/admin/users", s.protectedHandler(s.handleAdminUsers, "admin", "manage"))
	mux.HandleFunc(apiPrefix+"/admin/audit", s.protectedHandler(s.handleAdminAudit, "admin", "read"))

	// Metrics endpoint for Prometheus
	if s.config.EnableMetrics {
		mux.HandleFunc(s.config.MetricsPath, s.handlePrometheusMetrics)
	}

	// Documentation endpoint
	if s.config.EnableDocs {
		mux.HandleFunc(s.config.DocsPath, s.handleDocs)
		mux.HandleFunc(s.config.DocsPath+"/", s.handleDocs)
	}

	// Pprof endpoints (if enabled)
	if s.config.EnablePprof {
		mux.HandleFunc(s.config.PprofPath+"/", s.handlePprof)
	}
}

func (s *Server) setupMiddleware(handler http.Handler) http.Handler {
	// Apply middleware in reverse order (last applied = first executed)

	// Logging middleware
	handler = s.loggingMiddleware(handler)

	// Performance monitoring
	handler = s.performanceMiddleware(handler)

	// Security headers
	handler = s.security.SecureHeadersHandler(handler.ServeHTTP)

	// CORS
	if s.config.EnableCORS {
		handler = s.security.CORSHandler(handler.ServeHTTP)
	}

	// Rate limiting
	handler = s.security.RateLimitHandler(handler.ServeHTTP)

	// Request ID middleware
	handler = s.requestIDMiddleware(handler)

	// Recovery middleware (should be last)
	handler = s.recoveryMiddleware(handler)

	return handler
}

func (s *Server) protectedHandler(handler http.HandlerFunc, resource, action string) http.HandlerFunc {
	// Apply authentication and authorization
	return s.security.Chain(
		s.security.AuthHandler,
		s.security.RBACHandler(resource, action),
	)(handler)
}

// Middleware implementations

func (s *Server) loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Create a custom response writer to capture status code
		rw := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		next.ServeHTTP(rw, r)

		duration := time.Since(start)

		// Log request
		fmt.Printf("%s %s %d %v\n", r.Method, r.URL.Path, rw.statusCode, duration)

		// Log to audit system if authentication endpoints
		if s.isAuthEndpoint(r.URL.Path) {
			s.auditLogger.LogEntry(audit.AuditEntry{
				ID:        generateRequestID(),
				Timestamp: time.Now(),
				EventType: audit.EventTypeAccess,
				Actor: audit.Actor{
					Type:      "user",
					IPAddress: getClientIP(r),
					UserAgent: r.UserAgent(),
				},
				Resource: audit.Resource{
					Type: "endpoint",
					Path: r.URL.Path,
				},
				Action: r.Method,
				Result: audit.Result{
					Status:     getStatusString(rw.statusCode),
					StatusCode: rw.statusCode,
				},
				Risk: audit.RiskLevelLow,
			})
		}
	})
}

func (s *Server) performanceMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		timer := performance.NewTimerWithMonitor(
			fmt.Sprintf("http_%s_%s", r.Method, sanitizePath(r.URL.Path)),
			s.performanceMonitor,
		)

		next.ServeHTTP(w, r)

		timer.Stop(true) // Assume success for now
	})
}

func (s *Server) requestIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := generateRequestID()
		ctx := context.WithValue(r.Context(), "request_id", requestID)
		w.Header().Set("X-Request-ID", requestID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (s *Server) recoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				fmt.Printf("🚨 Panic recovered: %v\n", err)

				s.writeErrorResponse(w, http.StatusInternalServerError, "internal_error",
					"Internal server error", "")
			}
		}()

		next.ServeHTTP(w, r)
	})
}

// Helper types and functions

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func (s *Server) writeJSONResponse(w http.ResponseWriter, statusCode int, data interface{}) {
	response := APIResponse{
		Success:   statusCode < 400,
		Data:      data,
		Timestamp: time.Now(),
		Meta: &Meta{
			Version: "1.0.0",
		},
	}

	// Add request ID if available
	if requestID := w.Header().Get("X-Request-ID"); requestID != "" {
		response.RequestID = requestID
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(response)
}

func (s *Server) writeErrorResponse(w http.ResponseWriter, statusCode int, errorType, message, details string) {
	response := APIResponse{
		Success: false,
		Error: &APIError{
			Code:    statusCode,
			Type:    errorType,
			Message: message,
			Details: details,
		},
		Timestamp: time.Now(),
		Meta: &Meta{
			Version: "1.0.0",
		},
	}

	// Add request ID if available
	if requestID := w.Header().Get("X-Request-ID"); requestID != "" {
		response.RequestID = requestID
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(response)
}

func generateRequestID() string {
	return fmt.Sprintf("req_%d_%d", time.Now().UnixNano(), randomInt())
}

func getClientIP(r *http.Request) string {
	// Implementation similar to middleware
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		return xff
	}
	if xri := r.Header.Get("X-Real-IP"); xri != "" {
		return xri
	}
	return r.RemoteAddr
}

func getStatusString(code int) string {
	if code >= 200 && code < 300 {
		return "success"
	} else if code >= 400 && code < 500 {
		return "client_error"
	} else if code >= 500 {
		return "server_error"
	}
	return "unknown"
}

func (s *Server) isAuthEndpoint(path string) bool {
	authPaths := []string{"/auth/login", "/auth/refresh", "/auth/logout"}
	for _, authPath := range authPaths {
		if path == s.config.APIPrefix+authPath {
			return true
		}
	}
	return false
}

func sanitizePath(path string) string {
	// Replace path parameters with placeholders for consistent metrics
	// This is a simple implementation - in production, use a more sophisticated approach
	if len(path) > 50 {
		return "long_path"
	}
	return path
}

func randomInt() int {
	return int(time.Now().UnixNano() % 1000000)
}