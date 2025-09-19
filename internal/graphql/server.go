package graphql

import (
	"encoding/json"
	"fmt"
	"net/http"

	"context-extender/internal/database"
	"github.com/graphql-go/graphql"
)

// GraphQLServer handles GraphQL HTTP requests
type GraphQLServer struct {
	schema *graphql.Schema
	port   int
}

// NewGraphQLServer creates a new GraphQL server
func NewGraphQLServer(port int) (*GraphQLServer, error) {
	// Initialize database
	config := database.DefaultConfig()
	if err := database.Initialize(config); err != nil {
		return nil, fmt.Errorf("failed to initialize database: %w", err)
	}

	// Initialize schema
	if err := InitializeSchema(); err != nil {
		return nil, fmt.Errorf("failed to initialize schema: %w", err)
	}

	// Setup resolvers
	SetupResolvers()

	return &GraphQLServer{
		schema: GraphQLSchema,
		port:   port,
	}, nil
}

// Start starts the GraphQL server
func (s *GraphQLServer) Start() error {
	http.HandleFunc("/graphql", s.handleGraphQL)
	http.HandleFunc("/", s.handlePlayground)

	addr := fmt.Sprintf(":%d", s.port)
	fmt.Printf("🚀 GraphQL server starting on http://localhost%s\n", addr)
	fmt.Printf("📊 GraphQL endpoint: http://localhost%s/graphql\n", addr)
	fmt.Printf("🎮 GraphQL playground: http://localhost%s/\n", addr)

	return http.ListenAndServe(addr, nil)
}

// handleGraphQL handles GraphQL queries
func (s *GraphQLServer) handleGraphQL(w http.ResponseWriter, r *http.Request) {
	// Set CORS headers
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == "OPTIONS" {
		return
	}

	if r.Method != "POST" && r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var query string
	var variables map[string]interface{}

	if r.Method == "POST" {
		var requestBody struct {
			Query     string                 `json:"query"`
			Variables map[string]interface{} `json:"variables"`
		}

		if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		query = requestBody.Query
		variables = requestBody.Variables
	} else {
		// GET request
		query = r.URL.Query().Get("query")
		variablesParam := r.URL.Query().Get("variables")
		if variablesParam != "" {
			if err := json.Unmarshal([]byte(variablesParam), &variables); err != nil {
				http.Error(w, "Invalid variables", http.StatusBadRequest)
				return
			}
		}
	}

	// Execute GraphQL query
	result := graphql.Do(graphql.Params{
		Schema:         *s.schema,
		RequestString:  query,
		VariableValues: variables,
	})

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

// handlePlayground serves a simple GraphQL playground
func (s *GraphQLServer) handlePlayground(w http.ResponseWriter, r *http.Request) {
	playground := `
<!DOCTYPE html>
<html>
<head>
    <title>Context Extender GraphQL</title>
    <style>
        body {
            margin: 0;
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
        }
        .container {
            display: flex;
            height: 100vh;
        }
        .query-panel {
            flex: 1;
            padding: 20px;
            border-right: 1px solid #eee;
        }
        .result-panel {
            flex: 1;
            padding: 20px;
            background: #f8f9fa;
        }
        textarea {
            width: 100%;
            height: 300px;
            font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', monospace;
            border: 1px solid #ddd;
            padding: 10px;
            border-radius: 4px;
        }
        button {
            background: #007bff;
            color: white;
            border: none;
            padding: 10px 20px;
            border-radius: 4px;
            cursor: pointer;
            margin-top: 10px;
        }
        button:hover {
            background: #0056b3;
        }
        .result {
            background: white;
            border: 1px solid #ddd;
            border-radius: 4px;
            padding: 10px;
            height: 300px;
            overflow: auto;
            white-space: pre-wrap;
            font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', monospace;
        }
        .examples {
            margin-top: 20px;
        }
        .example {
            background: #e9ecef;
            padding: 10px;
            margin: 5px 0;
            border-radius: 4px;
            cursor: pointer;
            font-size: 12px;
        }
        .example:hover {
            background: #dee2e6;
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="query-panel">
            <h2>🔍 Context Extender GraphQL</h2>
            <textarea id="query" placeholder="Enter your GraphQL query here...">
query {
  stats {
    totalSessions
    totalConversations
    totalEvents
    oldestSession
    newestSession
  }
}</textarea>
            <button onclick="executeQuery()">Execute Query</button>

            <div class="examples">
                <h3>📚 Example Queries</h3>
                <div class="example" onclick="loadExample(this)">
query { stats { totalSessions totalConversations } }
                </div>
                <div class="example" onclick="loadExample(this)">
query { sessions(limit: 5) { id status createdAt } }
                </div>
                <div class="example" onclick="loadExample(this)">
query { conversations(limit: 3) { content messageType timestamp } }
                </div>
                <div class="example" onclick="loadExample(this)">
query { search(query: "hello") { totalCount conversations { content } } }
                </div>
            </div>
        </div>

        <div class="result-panel">
            <h2>📊 Result</h2>
            <div id="result" class="result">Execute a query to see results...</div>
        </div>
    </div>

    <script>
        function executeQuery() {
            const query = document.getElementById('query').value;
            const resultDiv = document.getElementById('result');

            resultDiv.textContent = 'Executing...';

            fetch('/graphql', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({ query: query })
            })
            .then(response => response.json())
            .then(data => {
                resultDiv.textContent = JSON.stringify(data, null, 2);
            })
            .catch(error => {
                resultDiv.textContent = 'Error: ' + error.message;
            });
        }

        function loadExample(element) {
            document.getElementById('query').value = element.textContent.trim();
        }

        // Execute initial query
        window.onload = () => executeQuery();
    </script>
</body>
</html>`

	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(playground))
}

// ExecuteQuery executes a GraphQL query and returns the result
func ExecuteQuery(query string, variables map[string]interface{}) *graphql.Result {
	return graphql.Do(graphql.Params{
		Schema:         *GraphQLSchema,
		RequestString:  query,
		VariableValues: variables,
	})
}