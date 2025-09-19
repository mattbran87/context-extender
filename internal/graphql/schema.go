package graphql

import (
	"github.com/graphql-go/graphql"
)

// GraphQLSchema defines the complete GraphQL schema
var GraphQLSchema *graphql.Schema

// SessionType represents a session in GraphQL
var SessionType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Session",
	Fields: graphql.Fields{
		"id": &graphql.Field{
			Type: graphql.String,
		},
		"createdAt": &graphql.Field{
			Type: graphql.String,
		},
		"updatedAt": &graphql.Field{
			Type: graphql.String,
		},
		"status": &graphql.Field{
			Type: graphql.String,
		},
		"metadata": &graphql.Field{
			Type: graphql.String,
		},
		"events": &graphql.Field{
			Type: graphql.NewList(EventType),
			Args: graphql.FieldConfigArgument{
				"limit": &graphql.ArgumentConfig{
					Type: graphql.Int,
				},
				"offset": &graphql.ArgumentConfig{
					Type: graphql.Int,
				},
			},
		},
		"conversations": &graphql.Field{
			Type: graphql.NewList(ConversationType),
			Args: graphql.FieldConfigArgument{
				"limit": &graphql.ArgumentConfig{
					Type: graphql.Int,
				},
				"offset": &graphql.ArgumentConfig{
					Type: graphql.Int,
				},
			},
		},
	},
})

// EventType represents an event in GraphQL
var EventType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Event",
	Fields: graphql.Fields{
		"id": &graphql.Field{
			Type: graphql.Int,
		},
		"sessionId": &graphql.Field{
			Type: graphql.String,
		},
		"eventType": &graphql.Field{
			Type: graphql.String,
		},
		"eventData": &graphql.Field{
			Type: graphql.String,
		},
		"timestamp": &graphql.Field{
			Type: graphql.String,
		},
		"sequenceNumber": &graphql.Field{
			Type: graphql.Int,
		},
	},
})

// ConversationType represents a conversation message in GraphQL
var ConversationType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Conversation",
	Fields: graphql.Fields{
		"id": &graphql.Field{
			Type: graphql.Int,
		},
		"sessionId": &graphql.Field{
			Type: graphql.String,
		},
		"messageType": &graphql.Field{
			Type: graphql.String,
		},
		"content": &graphql.Field{
			Type: graphql.String,
		},
		"timestamp": &graphql.Field{
			Type: graphql.String,
		},
		"tokenCount": &graphql.Field{
			Type: graphql.Int,
		},
		"modelInfo": &graphql.Field{
			Type: graphql.String,
		},
	},
})

// SearchResultType represents search results
var SearchResultType = graphql.NewObject(graphql.ObjectConfig{
	Name: "SearchResult",
	Fields: graphql.Fields{
		"sessions": &graphql.Field{
			Type: graphql.NewList(SessionType),
		},
		"conversations": &graphql.Field{
			Type: graphql.NewList(ConversationType),
		},
		"totalCount": &graphql.Field{
			Type: graphql.Int,
		},
	},
})

// StatsType represents database statistics
var StatsType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Stats",
	Fields: graphql.Fields{
		"totalSessions": &graphql.Field{
			Type: graphql.Int,
		},
		"totalEvents": &graphql.Field{
			Type: graphql.Int,
		},
		"totalConversations": &graphql.Field{
			Type: graphql.Int,
		},
		"totalImports": &graphql.Field{
			Type: graphql.Int,
		},
		"oldestSession": &graphql.Field{
			Type: graphql.String,
		},
		"newestSession": &graphql.Field{
			Type: graphql.String,
		},
	},
})

// QueryType defines the root query
var QueryType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Query",
	Fields: graphql.Fields{
		"session": &graphql.Field{
			Type:        SessionType,
			Description: "Get a session by ID",
			Args: graphql.FieldConfigArgument{
				"id": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
			},
		},
		"sessions": &graphql.Field{
			Type:        graphql.NewList(SessionType),
			Description: "Get all sessions with optional filtering",
			Args: graphql.FieldConfigArgument{
				"limit": &graphql.ArgumentConfig{
					Type:         graphql.Int,
					DefaultValue: 50,
				},
				"offset": &graphql.ArgumentConfig{
					Type:         graphql.Int,
					DefaultValue: 0,
				},
				"status": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
				"sortBy": &graphql.ArgumentConfig{
					Type:         graphql.String,
					DefaultValue: "created_at",
				},
				"sortOrder": &graphql.ArgumentConfig{
					Type:         graphql.String,
					DefaultValue: "DESC",
				},
			},
		},
		"events": &graphql.Field{
			Type:        graphql.NewList(EventType),
			Description: "Get events with optional filtering",
			Args: graphql.FieldConfigArgument{
				"sessionId": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
				"eventType": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
				"limit": &graphql.ArgumentConfig{
					Type:         graphql.Int,
					DefaultValue: 100,
				},
				"offset": &graphql.ArgumentConfig{
					Type:         graphql.Int,
					DefaultValue: 0,
				},
			},
		},
		"conversations": &graphql.Field{
			Type:        graphql.NewList(ConversationType),
			Description: "Get conversations with optional filtering",
			Args: graphql.FieldConfigArgument{
				"sessionId": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
				"messageType": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
				"limit": &graphql.ArgumentConfig{
					Type:         graphql.Int,
					DefaultValue: 100,
				},
				"offset": &graphql.ArgumentConfig{
					Type:         graphql.Int,
					DefaultValue: 0,
				},
			},
		},
		"search": &graphql.Field{
			Type:        SearchResultType,
			Description: "Search across conversations and sessions",
			Args: graphql.FieldConfigArgument{
				"query": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
				"limit": &graphql.ArgumentConfig{
					Type:         graphql.Int,
					DefaultValue: 50,
				},
				"searchSessions": &graphql.ArgumentConfig{
					Type:         graphql.Boolean,
					DefaultValue: true,
				},
				"searchConversations": &graphql.ArgumentConfig{
					Type:         graphql.Boolean,
					DefaultValue: true,
				},
			},
		},
		"stats": &graphql.Field{
			Type:        StatsType,
			Description: "Get database statistics",
		},
	},
})

// InitializeSchema creates and initializes the GraphQL schema
func InitializeSchema() error {
	var err error

	// Create the schema with query type and resolvers
	schema, err := graphql.NewSchema(graphql.SchemaConfig{
		Query: QueryType,
	})
	GraphQLSchema = &schema

	if err != nil {
		return err
	}

	return nil
}