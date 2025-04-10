package mcp

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// ServerConfig contains configuration options for the MCP server
type ServerConfig struct {
	Name        string
	Description string
	Version     string
}

// RequestHandler is a function that handles MCP requests
type RequestHandler func(ctx context.Context, params json.RawMessage) (interface{}, error)

// Server represents an MCP server
type Server struct {
	config           *ServerConfig
	requestHandlers  map[string]RequestHandler
	requestHandlerMu sync.RWMutex
}

// NewServer creates a new MCP server with the given configuration
func NewServer(config *ServerConfig) *Server {
	return &Server{
		config:          config,
		requestHandlers: make(map[string]RequestHandler),
	}
}

// RegisterRequestHandler registers a handler for a specific request type
func (s *Server) RegisterRequestHandler(requestType string, handler RequestHandler) {
	s.requestHandlerMu.Lock()
	defer s.requestHandlerMu.Unlock()
	s.requestHandlers[requestType] = handler
}

// HasHandler checks if a handler is registered for a request type
func (s *Server) HasHandler(requestType string) bool {
	s.requestHandlerMu.RLock()
	defer s.requestHandlerMu.RUnlock()
	_, exists := s.requestHandlers[requestType]
	return exists
}

// Connection represents an MCP WebSocket connection
type Connection struct {
	conn   *websocket.Conn
	server *Server
}

// NewConnection creates a new MCP connection
func NewConnection(conn *websocket.Conn) *Connection {
	return &Connection{
		conn: conn,
	}
}

// Request represents an MCP request
type Request struct {
	ID         string          `json:"id"`
	Type       string          `json:"type"`
	Parameters json.RawMessage `json:"parameters,omitempty"`
}

// Response represents an MCP response
type Response struct {
	ID      string          `json:"id"`
	Success bool            `json:"success"`
	Result  json.RawMessage `json:"result,omitempty"`
	Error   string          `json:"error,omitempty"`
}

// HandleConnection processes incoming messages for an MCP connection
func (s *Server) HandleConnection(conn *Connection) error {
	conn.server = s

	// Send server info
	serverInfo := map[string]interface{}{
		"name":        s.config.Name,
		"description": s.config.Description,
		"version":     s.config.Version,
		"protocol":    "mcp",
	}
	
	infoMsg, err := json.Marshal(map[string]interface{}{
		"type": "server_info",
		"info": serverInfo,
	})
	if err != nil {
		return fmt.Errorf("failed to marshal server info: %w", err)
	}

	if err := conn.conn.WriteMessage(websocket.TextMessage, infoMsg); err != nil {
		return fmt.Errorf("failed to send server info: %w", err)
	}

	// Handle incoming messages
	for {
		_, message, err := conn.conn.ReadMessage()
		if err != nil {
			// Normal close
			if websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway) {
				return nil
			}
			return fmt.Errorf("error reading message: %w", err)
		}

		// Parse request
		var request Request
		if err := json.Unmarshal(message, &request); err != nil {
			log.Printf("Error parsing request: %v", err)
			continue
		}

		// Process request in a goroutine
		go func(req Request) {
			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()

			result, err := s.handleRequest(ctx, req)
			
			response := Response{
				ID:      req.ID,
				Success: err == nil,
			}

			if err != nil {
				response.Error = err.Error()
			} else if result != nil {
				resultBytes, e := json.Marshal(result)
				if e != nil {
					response.Success = false
					response.Error = fmt.Sprintf("failed to marshal result: %v", e)
				} else {
					response.Result = resultBytes
				}
			}

			responseBytes, e := json.Marshal(response)
			if e != nil {
				log.Printf("Error marshaling response: %v", e)
				return
			}

			if err := conn.conn.WriteMessage(websocket.TextMessage, responseBytes); err != nil {
				log.Printf("Error sending response: %v", err)
			}
		}(request)
	}
}

// handleRequest processes an MCP request
func (s *Server) handleRequest(ctx context.Context, req Request) (interface{}, error) {
	// Get handler
	s.requestHandlerMu.RLock()
	handler, exists := s.requestHandlers[req.Type]
	s.requestHandlerMu.RUnlock()

	if !exists {
		return nil, errors.New("unknown request type")
	}

	// Call handler
	return handler(ctx, req.Parameters)
}
