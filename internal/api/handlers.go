package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"context-extender/internal/auth"
	"context-extender/internal/performance"
)

// Authentication Handlers

// LoginRequest represents a login request
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// LoginResponse represents a login response
type LoginResponse struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	TokenType    string    `json:"token_type"`
	ExpiresIn    int       `json:"expires_in"`
	ExpiresAt    time.Time `json:"expires_at"`
	User         UserInfo  `json:"user"`
}

// UserInfo represents user information
type UserInfo struct {
	ID       string   `json:"id"`
	Username string   `json:"username"`
	Email    string   `json:"email"`
	Roles    []string `json:"roles"`
}

// RefreshRequest represents a token refresh request
type RefreshRequest struct {
	RefreshToken string `json:"refresh_token"`
}

func (s *Server) handleLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		s.writeErrorResponse(w, http.StatusMethodNotAllowed, "method_not_allowed",
			"Method not allowed", "Only POST is supported")
		return
	}

	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.writeErrorResponse(w, http.StatusBadRequest, "invalid_request",
			"Invalid JSON", err.Error())
		return
	}

	// Validate credentials (simplified for demo)
	userID, roles, valid := s.validateCredentials(req.Username, req.Password)
	if !valid {
		s.auditLogger.LogAuthentication("", req.Username, false, "invalid_credentials", map[string]interface{}{
			"ip": getClientIP(r),
			"user_agent": r.UserAgent(),
		})

		s.writeErrorResponse(w, http.StatusUnauthorized, "invalid_credentials",
			"Invalid username or password", "")
		return
	}

	// Generate JWT tokens
	token, err := s.jwtManager.GenerateToken(userID, req.Username, req.Username+"@example.com", roles)
	if err != nil {
		s.writeErrorResponse(w, http.StatusInternalServerError, "token_generation_failed",
			"Failed to generate token", err.Error())
		return
	}

	// Log successful authentication
	s.auditLogger.LogAuthentication(userID, req.Username, true, "successful", map[string]interface{}{
		"ip": getClientIP(r),
		"user_agent": r.UserAgent(),
	})

	// Create response
	response := LoginResponse{
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
		TokenType:    token.TokenType,
		ExpiresIn:    token.ExpiresIn,
		ExpiresAt:    token.ExpiresAt,
		User: UserInfo{
			ID:       userID,
			Username: req.Username,
			Email:    req.Username + "@example.com",
			Roles:    roles,
		},
	}

	s.writeJSONResponse(w, http.StatusOK, response)
}

func (s *Server) handleRefresh(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		s.writeErrorResponse(w, http.StatusMethodNotAllowed, "method_not_allowed",
			"Method not allowed", "Only POST is supported")
		return
	}

	var req RefreshRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.writeErrorResponse(w, http.StatusBadRequest, "invalid_request",
			"Invalid JSON", err.Error())
		return
	}

	// Refresh token
	newToken, err := s.jwtManager.RefreshToken(req.RefreshToken)
	if err != nil {
		s.writeErrorResponse(w, http.StatusUnauthorized, "invalid_refresh_token",
			"Invalid refresh token", err.Error())
		return
	}

	response := LoginResponse{
		AccessToken:  newToken.AccessToken,
		RefreshToken: newToken.RefreshToken,
		TokenType:    newToken.TokenType,
		ExpiresIn:    newToken.ExpiresIn,
		ExpiresAt:    newToken.ExpiresAt,
	}

	s.writeJSONResponse(w, http.StatusOK, response)
}

func (s *Server) handleLogout(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		s.writeErrorResponse(w, http.StatusMethodNotAllowed, "method_not_allowed",
			"Method not allowed", "Only POST is supported")
		return
	}

	// Extract token from Authorization header
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		s.writeErrorResponse(w, http.StatusBadRequest, "missing_token",
			"Authorization header required", "")
		return
	}

	token, err := s.jwtManager.ExtractTokenFromHeader(authHeader)
	if err != nil {
		s.writeErrorResponse(w, http.StatusBadRequest, "invalid_token_format",
			"Invalid token format", err.Error())
		return
	}

	// Revoke token (add to blacklist)
	if err := s.jwtManager.RevokeToken(token); err != nil {
		s.writeErrorResponse(w, http.StatusInternalServerError, "revocation_failed",
			"Failed to revoke token", err.Error())
		return
	}

	s.writeJSONResponse(w, http.StatusOK, map[string]string{
		"message": "Successfully logged out",
	})
}

// Session Handlers

func (s *Server) handleSessions(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		s.handleGetSessions(w, r)
	case http.MethodPost:
		s.handleCreateSession(w, r)
	default:
		s.writeErrorResponse(w, http.StatusMethodNotAllowed, "method_not_allowed",
			"Method not allowed", "Only GET and POST are supported")
	}
}

func (s *Server) handleGetSessions(w http.ResponseWriter, r *http.Request) {
	// Get query parameters for pagination
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}

	perPage, _ := strconv.Atoi(r.URL.Query().Get("per_page"))
	if perPage < 1 || perPage > 100 {
		perPage = 20
	}

	// Get sessions (simplified - would typically query database)
	sessions := []map[string]interface{}{
		{
			"id":          "session_1",
			"user_id":     "user_123",
			"status":      "active",
			"created_at":  time.Now().Add(-2 * time.Hour),
			"updated_at":  time.Now().Add(-5 * time.Minute),
			"event_count": 45,
		},
		{
			"id":          "session_2",
			"user_id":     "user_456",
			"status":      "completed",
			"created_at":  time.Now().Add(-6 * time.Hour),
			"updated_at":  time.Now().Add(-1 * time.Hour),
			"event_count": 123,
		},
	}

	// Apply pagination (simplified)
	start := (page - 1) * perPage
	end := start + perPage
	if start >= len(sessions) {
		sessions = []map[string]interface{}{}
	} else if end > len(sessions) {
		sessions = sessions[start:]
	} else {
		sessions = sessions[start:end]
	}

	response := map[string]interface{}{
		"sessions": sessions,
		"pagination": Pagination{
			Page:       page,
			PerPage:    perPage,
			Total:      2, // Would be actual count
			TotalPages: 1,
		},
	}

	s.writeJSONResponse(w, http.StatusOK, response)
}

func (s *Server) handleCreateSession(w http.ResponseWriter, r *http.Request) {
	var req map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.writeErrorResponse(w, http.StatusBadRequest, "invalid_request",
			"Invalid JSON", err.Error())
		return
	}

	// Create session (simplified)
	sessionID := fmt.Sprintf("session_%d", time.Now().UnixNano())

	session := map[string]interface{}{
		"id":         sessionID,
		"user_id":    req["user_id"],
		"status":     "active",
		"created_at": time.Now(),
		"updated_at": time.Now(),
		"metadata":   req["metadata"],
	}

	s.writeJSONResponse(w, http.StatusCreated, session)
}

func (s *Server) handleSessionDetail(w http.ResponseWriter, r *http.Request) {
	// Extract session ID from URL
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 2 {
		s.writeErrorResponse(w, http.StatusBadRequest, "invalid_session_id",
			"Session ID required", "")
		return
	}

	sessionID := parts[len(parts)-1]

	switch r.Method {
	case http.MethodGet:
		s.handleGetSession(w, r, sessionID)
	case http.MethodPut:
		s.handleUpdateSession(w, r, sessionID)
	case http.MethodDelete:
		s.handleDeleteSession(w, r, sessionID)
	default:
		s.writeErrorResponse(w, http.StatusMethodNotAllowed, "method_not_allowed",
			"Method not allowed", "Only GET, PUT, and DELETE are supported")
	}
}

func (s *Server) handleGetSession(w http.ResponseWriter, r *http.Request, sessionID string) {
	// Get session details (simplified)
	session := map[string]interface{}{
		"id":          sessionID,
		"user_id":     "user_123",
		"status":      "active",
		"created_at":  time.Now().Add(-2 * time.Hour),
		"updated_at":  time.Now().Add(-5 * time.Minute),
		"event_count": 45,
		"metadata": map[string]interface{}{
			"project": "context-extender",
			"version": "1.0.0",
		},
		"events": []map[string]interface{}{
			{
				"id":        "event_1",
				"type":      "user_prompt",
				"timestamp": time.Now().Add(-10 * time.Minute),
				"data":      "Hello, how are you?",
			},
			{
				"id":        "event_2",
				"type":      "ai_response",
				"timestamp": time.Now().Add(-9 * time.Minute),
				"data":      "I'm doing well, thank you!",
			},
		},
	}

	s.writeJSONResponse(w, http.StatusOK, session)
}

func (s *Server) handleUpdateSession(w http.ResponseWriter, r *http.Request, sessionID string) {
	var req map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.writeErrorResponse(w, http.StatusBadRequest, "invalid_request",
			"Invalid JSON", err.Error())
		return
	}

	// Update session (simplified)
	session := map[string]interface{}{
		"id":         sessionID,
		"status":     req["status"],
		"updated_at": time.Now(),
		"metadata":   req["metadata"],
	}

	s.writeJSONResponse(w, http.StatusOK, session)
}

func (s *Server) handleDeleteSession(w http.ResponseWriter, r *http.Request, sessionID string) {
	// Delete session (simplified)
	s.writeJSONResponse(w, http.StatusOK, map[string]string{
		"message": fmt.Sprintf("Session %s deleted successfully", sessionID),
	})
}

// User Handlers

func (s *Server) handleUsers(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		s.writeErrorResponse(w, http.StatusMethodNotAllowed, "method_not_allowed",
			"Method not allowed", "Only GET is supported")
		return
	}

	users := s.rbac.ListUsers()
	userList := make([]map[string]interface{}, len(users))

	for i, user := range users {
		roles, _ := s.rbac.GetUserRoles(user.ID)
		roleNames := make([]string, len(roles))
		for j, role := range roles {
			roleNames[j] = role.Name
		}

		userList[i] = map[string]interface{}{
			"id":         user.ID,
			"username":   user.Username,
			"email":      user.Email,
			"roles":      roleNames,
			"active":     user.Active,
			"created_at": user.CreatedAt,
			"last_login": user.LastLogin,
		}
	}

	s.writeJSONResponse(w, http.StatusOK, map[string]interface{}{
		"users": userList,
	})
}

func (s *Server) handleUserDetail(w http.ResponseWriter, r *http.Request) {
	// Extract user ID from URL
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 2 {
		s.writeErrorResponse(w, http.StatusBadRequest, "invalid_user_id",
			"User ID required", "")
		return
	}

	userID := parts[len(parts)-1]

	user, err := s.rbac.GetUser(userID)
	if err != nil {
		s.writeErrorResponse(w, http.StatusNotFound, "user_not_found",
			"User not found", err.Error())
		return
	}

	roles, _ := s.rbac.GetUserRoles(user.ID)
	permissions, _ := s.rbac.GetUserPermissions(user.ID)

	roleNames := make([]string, len(roles))
	for i, role := range roles {
		roleNames[i] = role.Name
	}

	permNames := make([]string, len(permissions))
	for i, perm := range permissions {
		permNames[i] = perm.Name
	}

	response := map[string]interface{}{
		"id":          user.ID,
		"username":    user.Username,
		"email":       user.Email,
		"roles":       roleNames,
		"permissions": permNames,
		"active":      user.Active,
		"created_at":  user.CreatedAt,
		"updated_at":  user.UpdatedAt,
		"last_login":  user.LastLogin,
		"attributes":  user.Attributes,
	}

	s.writeJSONResponse(w, http.StatusOK, response)
}

// Metrics Handlers

func (s *Server) handleMetrics(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		s.writeErrorResponse(w, http.StatusMethodNotAllowed, "method_not_allowed",
			"Method not allowed", "Only GET is supported")
		return
	}

	metrics, err := s.metricsCollector.GetCurrentMetrics()
	if err != nil {
		s.writeErrorResponse(w, http.StatusInternalServerError, "metrics_error",
			"Failed to get metrics", err.Error())
		return
	}

	s.writeJSONResponse(w, http.StatusOK, metrics)
}

func (s *Server) handleMetricsExport(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		s.writeErrorResponse(w, http.StatusMethodNotAllowed, "method_not_allowed",
			"Method not allowed", "Only GET is supported")
		return
	}

	format := r.URL.Query().Get("format")
	if format == "" {
		format = "json"
	}

	var exportFormat performance.ExportFormat
	switch format {
	case "json":
		exportFormat = performance.FormatJSON
	case "prometheus":
		exportFormat = performance.FormatPrometheus
	case "csv":
		exportFormat = performance.FormatCSV
	default:
		s.writeErrorResponse(w, http.StatusBadRequest, "invalid_format",
			"Invalid export format", "Supported formats: json, prometheus, csv")
		return
	}

	data, err := s.metricsCollector.ExportMetrics(exportFormat)
	if err != nil {
		s.writeErrorResponse(w, http.StatusInternalServerError, "export_error",
			"Failed to export metrics", err.Error())
		return
	}

	// Set appropriate content type
	switch format {
	case "json":
		w.Header().Set("Content-Type", "application/json")
	case "prometheus":
		w.Header().Set("Content-Type", "text/plain")
	case "csv":
		w.Header().Set("Content-Type", "text/csv")
	}

	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

// Performance Handlers

func (s *Server) handlePerformance(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		s.writeErrorResponse(w, http.StatusMethodNotAllowed, "method_not_allowed",
			"Method not allowed", "Only GET is supported")
		return
	}

	summary := s.performanceMonitor.GetSummary()
	systemMetrics := s.performanceMonitor.GetSystemMetrics()

	response := map[string]interface{}{
		"summary":        summary,
		"system_metrics": systemMetrics,
	}

	s.writeJSONResponse(w, http.StatusOK, response)
}

func (s *Server) handlePerformanceReport(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		s.writeErrorResponse(w, http.StatusMethodNotAllowed, "method_not_allowed",
			"Method not allowed", "Only GET is supported")
		return
	}

	report := s.performanceMonitor.ExportReport()
	s.writeJSONResponse(w, http.StatusOK, report)
}

// Admin Handlers

func (s *Server) handleAdminUsers(w http.ResponseWriter, r *http.Request) {
	// Similar to handleUsers but with additional admin functionality
	s.handleUsers(w, r)
}

func (s *Server) handleAdminAudit(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		s.writeErrorResponse(w, http.StatusMethodNotAllowed, "method_not_allowed",
			"Method not allowed", "Only GET is supported")
		return
	}

	// Get query parameters for filtering
	startTimeStr := r.URL.Query().Get("start_time")
	endTimeStr := r.URL.Query().Get("end_time")
	eventType := r.URL.Query().Get("event_type")
	actorID := r.URL.Query().Get("actor_id")

	filter := map[string]interface{}{
		"event_type": eventType,
		"actor_id":   actorID,
	}

	if startTimeStr != "" {
		filter["start_time"] = startTimeStr
	}
	if endTimeStr != "" {
		filter["end_time"] = endTimeStr
	}

	// Mock audit entries (would query actual audit storage)
	auditEntries := []map[string]interface{}{
		{
			"id":         "audit_1",
			"timestamp":  time.Now().Add(-1 * time.Hour),
			"event_type": "authentication",
			"actor":      map[string]string{"id": "user_123", "username": "john"},
			"action":     "login",
			"result":     map[string]interface{}{"status": "success"},
			"risk":       "low",
		},
		{
			"id":         "audit_2",
			"timestamp":  time.Now().Add(-30 * time.Minute),
			"event_type": "authorization",
			"actor":      map[string]string{"id": "user_456", "username": "jane"},
			"action":     "access_admin",
			"result":     map[string]interface{}{"status": "denied"},
			"risk":       "medium",
		},
	}

	response := map[string]interface{}{
		"audit_entries": auditEntries,
		"filter":        filter,
		"total":         len(auditEntries),
	}

	s.writeJSONResponse(w, http.StatusOK, response)
}

// Health and System Handlers

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	health := map[string]interface{}{
		"status":    "healthy",
		"timestamp": time.Now(),
		"uptime":    time.Since(time.Now().Add(-2 * time.Hour)), // Mock uptime
		"version":   "1.0.0",
		"components": map[string]string{
			"database":    "healthy",
			"auth":        "healthy",
			"performance": "healthy",
			"audit":       "healthy",
		},
	}

	s.writeJSONResponse(w, http.StatusOK, health)
}

func (s *Server) handlePrometheusMetrics(w http.ResponseWriter, r *http.Request) {
	data, err := s.performanceMonitor.Export(performance.FormatPrometheus)
	if err != nil {
		http.Error(w, "Failed to export metrics", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.Write(data)
}

func (s *Server) handleDocs(w http.ResponseWriter, r *http.Request) {
	// Simple API documentation
	docs := `
<!DOCTYPE html>
<html>
<head>
    <title>Context Extender API Documentation</title>
    <style>
        body { font-family: Arial, sans-serif; margin: 40px; }
        h1, h2 { color: #333; }
        .endpoint { margin: 20px 0; padding: 15px; border-left: 4px solid #007cba; background: #f9f9f9; }
        .method { font-weight: bold; color: #007cba; }
        .path { font-family: monospace; }
    </style>
</head>
<body>
    <h1>🚀 Context Extender API v1.0</h1>

    <h2>Authentication Endpoints</h2>
    <div class="endpoint">
        <span class="method">POST</span> <span class="path">/api/v1/auth/login</span><br>
        Authenticate user and get JWT tokens
    </div>
    <div class="endpoint">
        <span class="method">POST</span> <span class="path">/api/v1/auth/refresh</span><br>
        Refresh access token using refresh token
    </div>
    <div class="endpoint">
        <span class="method">POST</span> <span class="path">/api/v1/auth/logout</span><br>
        Logout and revoke tokens
    </div>

    <h2>Session Endpoints</h2>
    <div class="endpoint">
        <span class="method">GET</span> <span class="path">/api/v1/sessions</span><br>
        List user sessions with pagination
    </div>
    <div class="endpoint">
        <span class="method">POST</span> <span class="path">/api/v1/sessions</span><br>
        Create new session
    </div>
    <div class="endpoint">
        <span class="method">GET</span> <span class="path">/api/v1/sessions/{id}</span><br>
        Get session details
    </div>

    <h2>Metrics Endpoints</h2>
    <div class="endpoint">
        <span class="method">GET</span> <span class="path">/api/v1/metrics</span><br>
        Get current system metrics
    </div>
    <div class="endpoint">
        <span class="method">GET</span> <span class="path">/api/v1/metrics/export?format=json|prometheus|csv</span><br>
        Export metrics in specified format
    </div>

    <h2>Performance Endpoints</h2>
    <div class="endpoint">
        <span class="method">GET</span> <span class="path">/api/v1/performance</span><br>
        Get performance summary and system metrics
    </div>
    <div class="endpoint">
        <span class="method">GET</span> <span class="path">/api/v1/performance/report</span><br>
        Get comprehensive performance report
    </div>

    <h2>System Endpoints</h2>
    <div class="endpoint">
        <span class="method">GET</span> <span class="path">/health</span><br>
        Health check endpoint
    </div>
    <div class="endpoint">
        <span class="method">GET</span> <span class="path">/metrics</span><br>
        Prometheus metrics endpoint
    </div>
</body>
</html>`

	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(docs))
}

func (s *Server) handlePprof(w http.ResponseWriter, r *http.Request) {
	// Pprof endpoint (simplified)
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Pprof endpoints would be available here in production"))
}

// Helper functions

func (s *Server) validateCredentials(username, password string) (string, []string, bool) {
	// Simplified credential validation
	credentials := map[string]struct {
		userID   string
		password string
		roles    []string
	}{
		"admin": {
			userID:   "user_admin",
			password: "admin123",
			roles:    []string{"admin"},
		},
		"user": {
			userID:   "user_regular",
			password: "user123",
			roles:    []string{"user"},
		},
		"viewer": {
			userID:   "user_viewer",
			password: "viewer123",
			roles:    []string{"viewer"},
		},
	}

	if cred, exists := credentials[username]; exists && cred.password == password {
		return cred.userID, cred.roles, true
	}

	return "", nil, false
}