package api

import (
	"encoding/json"
	"net/http"
)

// OpenAPISpec represents the OpenAPI 3.0 specification
type OpenAPISpec struct {
	OpenAPI    string                 `json:"openapi"`
	Info       OpenAPIInfo            `json:"info"`
	Servers    []OpenAPIServer        `json:"servers"`
	Paths      map[string]OpenAPIPath `json:"paths"`
	Components OpenAPIComponents      `json:"components"`
}

// OpenAPIInfo contains API information
type OpenAPIInfo struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Version     string `json:"version"`
	Contact     OpenAPIContact `json:"contact"`
	License     OpenAPILicense `json:"license"`
}

// OpenAPIContact contains contact information
type OpenAPIContact struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	URL   string `json:"url"`
}

// OpenAPILicense contains license information
type OpenAPILicense struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

// OpenAPIServer represents a server
type OpenAPIServer struct {
	URL         string `json:"url"`
	Description string `json:"description"`
}

// OpenAPIPath represents path operations
type OpenAPIPath struct {
	Get    *OpenAPIOperation `json:"get,omitempty"`
	Post   *OpenAPIOperation `json:"post,omitempty"`
	Put    *OpenAPIOperation `json:"put,omitempty"`
	Delete *OpenAPIOperation `json:"delete,omitempty"`
}

// OpenAPIOperation represents an API operation
type OpenAPIOperation struct {
	Tags        []string                        `json:"tags,omitempty"`
	Summary     string                          `json:"summary"`
	Description string                          `json:"description"`
	OperationID string                          `json:"operationId"`
	Parameters  []OpenAPIParameter              `json:"parameters,omitempty"`
	RequestBody *OpenAPIRequestBody             `json:"requestBody,omitempty"`
	Responses   map[string]OpenAPIResponse      `json:"responses"`
	Security    []map[string][]string           `json:"security,omitempty"`
}

// OpenAPIParameter represents a parameter
type OpenAPIParameter struct {
	Name        string         `json:"name"`
	In          string         `json:"in"`
	Required    bool           `json:"required"`
	Description string         `json:"description"`
	Schema      OpenAPISchema  `json:"schema"`
}

// OpenAPIRequestBody represents request body
type OpenAPIRequestBody struct {
	Description string                     `json:"description"`
	Required    bool                       `json:"required"`
	Content     map[string]OpenAPIContent  `json:"content"`
}

// OpenAPIResponse represents a response
type OpenAPIResponse struct {
	Description string                     `json:"description"`
	Content     map[string]OpenAPIContent  `json:"content,omitempty"`
	Headers     map[string]OpenAPIHeader   `json:"headers,omitempty"`
}

// OpenAPIContent represents content
type OpenAPIContent struct {
	Schema OpenAPISchema `json:"schema"`
}

// OpenAPIHeader represents a header
type OpenAPIHeader struct {
	Description string        `json:"description"`
	Schema      OpenAPISchema `json:"schema"`
}

// OpenAPISchema represents a schema
type OpenAPISchema struct {
	Type        string                    `json:"type,omitempty"`
	Format      string                    `json:"format,omitempty"`
	Properties  map[string]OpenAPISchema  `json:"properties,omitempty"`
	Items       *OpenAPISchema            `json:"items,omitempty"`
	Required    []string                  `json:"required,omitempty"`
	Example     interface{}               `json:"example,omitempty"`
	Ref         string                    `json:"$ref,omitempty"`
}

// OpenAPIComponents contains reusable components
type OpenAPIComponents struct {
	Schemas         map[string]OpenAPISchema       `json:"schemas"`
	SecuritySchemes map[string]OpenAPISecurityScheme `json:"securitySchemes"`
}

// OpenAPISecurityScheme represents a security scheme
type OpenAPISecurityScheme struct {
	Type         string `json:"type"`
	Scheme       string `json:"scheme,omitempty"`
	BearerFormat string `json:"bearerFormat,omitempty"`
	Description  string `json:"description,omitempty"`
}

// GetOpenAPISpec returns the complete OpenAPI specification
func GetOpenAPISpec() *OpenAPISpec {
	return &OpenAPISpec{
		OpenAPI: "3.0.3",
		Info: OpenAPIInfo{
			Title:       "Context Extender API",
			Description: "REST API for Context Extender CLI tool with session management, performance monitoring, and security features",
			Version:     "1.0.0",
			Contact: OpenAPIContact{
				Name:  "Context Extender Team",
				Email: "support@context-extender.dev",
				URL:   "https://github.com/context-extender/context-extender",
			},
			License: OpenAPILicense{
				Name: "MIT",
				URL:  "https://opensource.org/licenses/MIT",
			},
		},
		Servers: []OpenAPIServer{
			{
				URL:         "http://localhost:8080/api/v1",
				Description: "Development server",
			},
			{
				URL:         "https://api.context-extender.dev/v1",
				Description: "Production server",
			},
		},
		Paths: getPaths(),
		Components: OpenAPIComponents{
			Schemas:         getSchemas(),
			SecuritySchemes: getSecuritySchemes(),
		},
	}
}

// getPaths returns all API paths
func getPaths() map[string]OpenAPIPath {
	return map[string]OpenAPIPath{
		"/health": {
			Get: &OpenAPIOperation{
				Tags:        []string{"Health"},
				Summary:     "Health check",
				Description: "Check API server health status",
				OperationID: "getHealth",
				Responses: map[string]OpenAPIResponse{
					"200": {
						Description: "Server is healthy",
						Content: map[string]OpenAPIContent{
							"application/json": {
								Schema: OpenAPISchema{
									Ref: "#/components/schemas/HealthResponse",
								},
							},
						},
					},
				},
			},
		},
		"/auth/login": {
			Post: &OpenAPIOperation{
				Tags:        []string{"Authentication"},
				Summary:     "User login",
				Description: "Authenticate user and return JWT tokens",
				OperationID: "login",
				RequestBody: &OpenAPIRequestBody{
					Description: "Login credentials",
					Required:    true,
					Content: map[string]OpenAPIContent{
						"application/json": {
							Schema: OpenAPISchema{
								Ref: "#/components/schemas/LoginRequest",
							},
						},
					},
				},
				Responses: map[string]OpenAPIResponse{
					"200": {
						Description: "Login successful",
						Content: map[string]OpenAPIContent{
							"application/json": {
								Schema: OpenAPISchema{
									Ref: "#/components/schemas/LoginResponse",
								},
							},
						},
					},
					"401": {
						Description: "Invalid credentials",
						Content: map[string]OpenAPIContent{
							"application/json": {
								Schema: OpenAPISchema{
									Ref: "#/components/schemas/ErrorResponse",
								},
							},
						},
					},
				},
			},
		},
		"/auth/refresh": {
			Post: &OpenAPIOperation{
				Tags:        []string{"Authentication"},
				Summary:     "Refresh token",
				Description: "Refresh access token using refresh token",
				OperationID: "refreshToken",
				RequestBody: &OpenAPIRequestBody{
					Description: "Refresh token",
					Required:    true,
					Content: map[string]OpenAPIContent{
						"application/json": {
							Schema: OpenAPISchema{
								Ref: "#/components/schemas/RefreshRequest",
							},
						},
					},
				},
				Responses: map[string]OpenAPIResponse{
					"200": {
						Description: "Token refreshed successfully",
						Content: map[string]OpenAPIContent{
							"application/json": {
								Schema: OpenAPISchema{
									Ref: "#/components/schemas/TokenResponse",
								},
							},
						},
					},
					"401": {
						Description: "Invalid refresh token",
						Content: map[string]OpenAPIContent{
							"application/json": {
								Schema: OpenAPISchema{
									Ref: "#/components/schemas/ErrorResponse",
								},
							},
						},
					},
				},
			},
		},
		"/sessions": {
			Get: &OpenAPIOperation{
				Tags:        []string{"Sessions"},
				Summary:     "List sessions",
				Description: "Get list of user sessions with optional filtering",
				OperationID: "listSessions",
				Security: []map[string][]string{
					{"bearerAuth": {}},
				},
				Parameters: []OpenAPIParameter{
					{
						Name:        "page",
						In:          "query",
						Required:    false,
						Description: "Page number for pagination",
						Schema: OpenAPISchema{
							Type:    "integer",
							Example: 1,
						},
					},
					{
						Name:        "limit",
						In:          "query",
						Required:    false,
						Description: "Number of items per page",
						Schema: OpenAPISchema{
							Type:    "integer",
							Example: 20,
						},
					},
					{
						Name:        "status",
						In:          "query",
						Required:    false,
						Description: "Filter by session status",
						Schema: OpenAPISchema{
							Type:    "string",
							Example: "active",
						},
					},
				},
				Responses: map[string]OpenAPIResponse{
					"200": {
						Description: "List of sessions",
						Content: map[string]OpenAPIContent{
							"application/json": {
								Schema: OpenAPISchema{
									Ref: "#/components/schemas/SessionListResponse",
								},
							},
						},
					},
					"401": {
						Description: "Unauthorized",
						Content: map[string]OpenAPIContent{
							"application/json": {
								Schema: OpenAPISchema{
									Ref: "#/components/schemas/ErrorResponse",
								},
							},
						},
					},
				},
			},
		},
		"/metrics": {
			Get: &OpenAPIOperation{
				Tags:        []string{"Metrics"},
				Summary:     "Get metrics",
				Description: "Retrieve performance and usage metrics",
				OperationID: "getMetrics",
				Security: []map[string][]string{
					{"bearerAuth": {}},
				},
				Parameters: []OpenAPIParameter{
					{
						Name:        "start_time",
						In:          "query",
						Required:    false,
						Description: "Start time for metrics (ISO 8601)",
						Schema: OpenAPISchema{
							Type:   "string",
							Format: "date-time",
						},
					},
					{
						Name:        "end_time",
						In:          "query",
						Required:    false,
						Description: "End time for metrics (ISO 8601)",
						Schema: OpenAPISchema{
							Type:   "string",
							Format: "date-time",
						},
					},
				},
				Responses: map[string]OpenAPIResponse{
					"200": {
						Description: "Metrics data",
						Content: map[string]OpenAPIContent{
							"application/json": {
								Schema: OpenAPISchema{
									Ref: "#/components/schemas/MetricsResponse",
								},
							},
						},
					},
				},
			},
		},
	}
}

// getSchemas returns component schemas
func getSchemas() map[string]OpenAPISchema {
	return map[string]OpenAPISchema{
		"HealthResponse": {
			Type: "object",
			Properties: map[string]OpenAPISchema{
				"status": {
					Type:    "string",
					Example: "healthy",
				},
				"version": {
					Type:    "string",
					Example: "1.0.0",
				},
				"uptime": {
					Type:    "string",
					Example: "2h30m15s",
				},
				"components": {
					Type: "object",
					Properties: map[string]OpenAPISchema{
						"database": {
							Type:    "string",
							Example: "healthy",
						},
						"cache": {
							Type:    "string",
							Example: "healthy",
						},
					},
				},
			},
			Required: []string{"status", "version"},
		},
		"LoginRequest": {
			Type: "object",
			Properties: map[string]OpenAPISchema{
				"username": {
					Type:    "string",
					Example: "admin",
				},
				"password": {
					Type:    "string",
					Example: "password123",
				},
			},
			Required: []string{"username", "password"},
		},
		"LoginResponse": {
			Type: "object",
			Properties: map[string]OpenAPISchema{
				"access_token": {
					Type:    "string",
					Example: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
				},
				"refresh_token": {
					Type:    "string",
					Example: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
				},
				"expires_in": {
					Type:    "integer",
					Example: 900,
				},
				"token_type": {
					Type:    "string",
					Example: "Bearer",
				},
				"user": {
					Ref: "#/components/schemas/User",
				},
			},
			Required: []string{"access_token", "expires_in", "token_type"},
		},
		"RefreshRequest": {
			Type: "object",
			Properties: map[string]OpenAPISchema{
				"refresh_token": {
					Type:    "string",
					Example: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
				},
			},
			Required: []string{"refresh_token"},
		},
		"TokenResponse": {
			Type: "object",
			Properties: map[string]OpenAPISchema{
				"access_token": {
					Type:    "string",
					Example: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
				},
				"expires_in": {
					Type:    "integer",
					Example: 900,
				},
				"token_type": {
					Type:    "string",
					Example: "Bearer",
				},
			},
			Required: []string{"access_token", "expires_in", "token_type"},
		},
		"User": {
			Type: "object",
			Properties: map[string]OpenAPISchema{
				"id": {
					Type:    "string",
					Example: "user123",
				},
				"username": {
					Type:    "string",
					Example: "admin",
				},
				"email": {
					Type:    "string",
					Example: "admin@example.com",
				},
				"roles": {
					Type: "array",
					Items: &OpenAPISchema{
						Type: "string",
					},
					Example: []string{"admin"},
				},
			},
			Required: []string{"id", "username"},
		},
		"Session": {
			Type: "object",
			Properties: map[string]OpenAPISchema{
				"id": {
					Type:    "string",
					Example: "session123",
				},
				"user_id": {
					Type:    "string",
					Example: "user123",
				},
				"status": {
					Type:    "string",
					Example: "active",
				},
				"created_at": {
					Type:   "string",
					Format: "date-time",
				},
				"updated_at": {
					Type:   "string",
					Format: "date-time",
				},
				"metadata": {
					Type: "object",
				},
			},
			Required: []string{"id", "user_id", "status"},
		},
		"SessionListResponse": {
			Type: "object",
			Properties: map[string]OpenAPISchema{
				"sessions": {
					Type: "array",
					Items: &OpenAPISchema{
						Ref: "#/components/schemas/Session",
					},
				},
				"total": {
					Type:    "integer",
					Example: 100,
				},
				"page": {
					Type:    "integer",
					Example: 1,
				},
				"per_page": {
					Type:    "integer",
					Example: 20,
				},
			},
			Required: []string{"sessions", "total"},
		},
		"MetricsResponse": {
			Type: "object",
			Properties: map[string]OpenAPISchema{
				"performance": {
					Type: "object",
					Properties: map[string]OpenAPISchema{
						"total_operations": {
							Type:    "integer",
							Example: 1000,
						},
						"success_rate": {
							Type:    "number",
							Example: 99.5,
						},
						"average_duration": {
							Type:    "number",
							Example: 150.5,
						},
					},
				},
				"cache": {
					Type: "object",
					Properties: map[string]OpenAPISchema{
						"hit_rate": {
							Type:    "number",
							Example: 85.2,
						},
						"total_hits": {
							Type:    "integer",
							Example: 500,
						},
						"total_misses": {
							Type:    "integer",
							Example: 50,
						},
					},
				},
			},
		},
		"ErrorResponse": {
			Type: "object",
			Properties: map[string]OpenAPISchema{
				"success": {
					Type:    "boolean",
					Example: false,
				},
				"error": {
					Type: "object",
					Properties: map[string]OpenAPISchema{
						"code": {
							Type:    "integer",
							Example: 400,
						},
						"type": {
							Type:    "string",
							Example: "validation_error",
						},
						"message": {
							Type:    "string",
							Example: "Invalid request parameters",
						},
						"details": {
							Type:    "string",
							Example: "Username is required",
						},
					},
					Required: []string{"code", "type", "message"},
				},
				"timestamp": {
					Type:   "string",
					Format: "date-time",
				},
				"request_id": {
					Type:    "string",
					Example: "req_12345",
				},
			},
			Required: []string{"success", "error", "timestamp"},
		},
	}
}

// getSecuritySchemes returns security schemes
func getSecuritySchemes() map[string]OpenAPISecurityScheme {
	return map[string]OpenAPISecurityScheme{
		"bearerAuth": {
			Type:         "http",
			Scheme:       "bearer",
			BearerFormat: "JWT",
			Description:  "JWT Bearer token authentication",
		},
	}
}

// handleDocs serves the OpenAPI documentation
func (s *Server) handleDocs(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case s.config.DocsPath:
		// Serve Swagger UI HTML
		s.serveSwaggerUI(w, r)
	case s.config.DocsPath + "/openapi.json":
		// Serve OpenAPI spec JSON
		s.serveOpenAPISpec(w, r)
	default:
		// Serve Swagger UI assets (simplified)
		s.serveSwaggerUI(w, r)
	}
}

// serveOpenAPISpec serves the OpenAPI specification as JSON
func (s *Server) serveOpenAPISpec(w http.ResponseWriter, r *http.Request) {
	spec := GetOpenAPISpec()

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	if err := json.NewEncoder(w).Encode(spec); err != nil {
		s.writeErrorResponse(w, http.StatusInternalServerError, "encoding_error",
			"Failed to encode OpenAPI specification", err.Error())
		return
	}
}

// serveSwaggerUI serves a simple Swagger UI
func (s *Server) serveSwaggerUI(w http.ResponseWriter, r *http.Request) {
	html := `<!DOCTYPE html>
<html>
<head>
    <title>Context Extender API Documentation</title>
    <link rel="stylesheet" type="text/css" href="https://unpkg.com/swagger-ui-dist@3.52.5/swagger-ui.css" />
    <style>
        html {
            box-sizing: border-box;
            overflow: -moz-scrollbars-vertical;
            overflow-y: scroll;
        }
        *, *:before, *:after {
            box-sizing: inherit;
        }
        body {
            margin:0;
            background: #fafafa;
        }
    </style>
</head>
<body>
    <div id="swagger-ui"></div>
    <script src="https://unpkg.com/swagger-ui-dist@3.52.5/swagger-ui-bundle.js"></script>
    <script src="https://unpkg.com/swagger-ui-dist@3.52.5/swagger-ui-standalone-preset.js"></script>
    <script>
        window.onload = function() {
            const ui = SwaggerUIBundle({
                url: '` + s.config.DocsPath + `/openapi.json',
                dom_id: '#swagger-ui',
                deepLinking: true,
                presets: [
                    SwaggerUIBundle.presets.apis,
                    SwaggerUIStandalonePreset
                ],
                plugins: [
                    SwaggerUIBundle.plugins.DownloadUrl
                ],
                layout: "StandaloneLayout"
            });
        };
    </script>
</body>
</html>`

	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(html))
}