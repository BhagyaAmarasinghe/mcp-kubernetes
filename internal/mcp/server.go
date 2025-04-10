package mcp

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

// Server represents an MCP server
type Server struct {
	name            string
	description     string
	version         string
	requestHandlers map[string]RequestHandler
	handlerMu       sync.RWMutex
}

// RequestHandler is a function that handles MCP requests
type RequestHandler func(ctx context.Context, params json.RawMessage) (interface{}, error)

// NewServer creates a new MCP server
func NewServer(name, description, version string) *Server {
	return &Server{
		name:            name,
		description:     description,
		version:         version,
		requestHandlers: make(map[string]RequestHandler),
	}
}

// RegisterRequestHandler registers a handler for a specific request type
func (s *Server) RegisterRequestHandler(requestType string, handler RequestHandler) {
	s.handlerMu.Lock()
	defer s.handlerMu.Unlock()
	s.requestHandlers[requestType] = handler
}

// HasHandler checks if a handler is registered for a specific request type
func (s *Server) HasHandler(requestType string) bool {
	s.handlerMu.RLock()
	defer s.handlerMu.RUnlock()
	_, exists := s.requestHandlers[requestType]
	return exists
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

// ServerInfo represents an MCP server info message
type ServerInfo struct {
	Type string `json:"type"`
	Info struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		Version     string `json:"version"`
		Protocol    string `json:"protocol"`
	} `json:"info"`
}

// HandleConnection processes incoming WebSocket messages for an MCP connection
func (s *Server) HandleConnection(conn *websocket.Conn) error {
	// Send server info
	serverInfo := ServerInfo{
		Type: "server_info",
	}
	serverInfo.Info.Name = s.name
	serverInfo.Info.Description = s.description
	serverInfo.Info.Version = s.version
	serverInfo.Info.Protocol = "mcp"

	infoMsg, err := json.Marshal(serverInfo)
	if err != nil {
		return fmt.Errorf("failed to marshal server info: %w", err)
	}

	if err := conn.WriteMessage(websocket.TextMessage, infoMsg); err != nil {
		return fmt.Errorf("failed to send server info: %w", err)
	}

	// Handle incoming messages
	for {
		_, message, err := conn.ReadMessage()
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
			sendErrorResponse(conn, "", "Invalid request format")
			continue
		}

		// Generate an ID if not provided
		if request.ID == "" {
			request.ID = uuid.New().String()
		}

		// Process request in a goroutine
		go func(req Request) {
			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()

			result, err := s.handleRequest(ctx, req)
			s.sendResponse(conn, req.ID, result, err)
		}(request)
	}
}

// handleRequest processes an MCP request
func (s *Server) handleRequest(ctx context.Context, req Request) (interface{}, error) {
	s.handlerMu.RLock()
	handler, exists := s.requestHandlers[req.Type]
	s.handlerMu.RUnlock()

	if !exists {
		return nil, errors.New("unknown request type")
	}

	// Call handler
	return handler(ctx, req.Parameters)
}

// sendResponse sends an MCP response back to the client
func (s *Server) sendResponse(conn *websocket.Conn, id string, result interface{}, err error) {
	response := Response{
		ID:      id,
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

	if e := conn.WriteMessage(websocket.TextMessage, responseBytes); e != nil {
		log.Printf("Error sending response: %v", e)
	}
}

// sendErrorResponse sends an error response back to the client
func sendErrorResponse(conn *websocket.Conn, id string, errorMsg string) {
	response := Response{
		ID:      id,
		Success: false,
		Error:   errorMsg,
	}

	responseBytes, err := json.Marshal(response)
	if err != nil {
		log.Printf("Error marshaling error response: %v", err)
		return
	}

	if err := conn.WriteMessage(websocket.TextMessage, responseBytes); err != nil {
		log.Printf("Error sending error response: %v", err)
	}
}
