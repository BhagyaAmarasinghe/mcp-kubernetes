package server

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/BhagyaAmarasinghe/mcp-kubernetes/internal/kubernetes"
	"github.com/BhagyaAmarasinghe/mcp-kubernetes/internal/mcp"
	"github.com/gorilla/websocket"
)

// Server represents the MCP server for Kubernetes
type Server struct {
	port      int
	k8sExec   *kubernetes.Executor
	server    *http.Server
	upgrader  websocket.Upgrader
	mcpServer *mcp.Server
}

// NewServer creates a new MCP server for Kubernetes
func NewServer(port int) (*Server, error) {
	// Initialize Kubernetes executor
	k8sExec, err := kubernetes.NewExecutor()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize Kubernetes executor: %w", err)
	}

	// Create websocket upgrader
	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true // Allow all origins for now
		},
	}

	// Create server instance
	s := &Server{
		port:     port,
		k8sExec:  k8sExec,
		upgrader: upgrader,
	}

	// Create and configure MCP server
	mcpServer := mcp.NewServer("mcp-kubernetes", "MCP server for executing Kubernetes commands", "1.0.0")
	
	// Register handlers
	mcpServer.RegisterRequestHandler("execute", s.handleExecute)
	mcpServer.RegisterRequestHandler("get-contexts", s.handleGetContexts)
	mcpServer.RegisterRequestHandler("current-context", s.handleCurrentContext)
	mcpServer.RegisterRequestHandler("set-context", s.handleSetContext)
	
	s.mcpServer = mcpServer

	return s, nil
}

// Start starts the MCP server
func (s *Server) Start() error {
	// Create an HTTP server
	mux := http.NewServeMux()
	
	// Register MCP WebSocket handler
	mux.HandleFunc("/ws", s.handleWebSocket)
	
	// Add a simple health check endpoint
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("MCP Kubernetes server is running"))
	})

	// Create the HTTP server
	s.server = &http.Server{
		Addr:    ":" + strconv.Itoa(s.port),
		Handler: mux,
	}

	// Start the HTTP server
	return s.server.ListenAndServe()
}

// Stop gracefully stops the server
func (s *Server) Stop() error {
	if s.server != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		return s.server.Shutdown(ctx)
	}
	return nil
}

// handleWebSocket upgrades HTTP connections to WebSocket and handles MCP messages
func (s *Server) handleWebSocket(w http.ResponseWriter, r *http.Request) {
	// Upgrade the HTTP connection to WebSocket
	conn, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Failed to upgrade connection: %v", err)
		return
	}
	defer conn.Close()

	// Handle the MCP connection
	if err := s.mcpServer.HandleConnection(conn); err != nil {
		log.Printf("Error handling MCP connection: %v", err)
	}
}

// handleExecute executes a kubectl command
func (s *Server) handleExecute(ctx context.Context, params json.RawMessage) (interface{}, error) {
	// Parse parameters
	var p struct {
		Command string `json:"command"`
	}
	if err := json.Unmarshal(params, &p); err != nil {
		return nil, fmt.Errorf("invalid parameters: %w", err)
	}

	// Validate command
	if strings.TrimSpace(p.Command) == "" {
		return nil, fmt.Errorf("command cannot be empty")
	}

	// Execute the command
	output, err := s.k8sExec.Execute(ctx, p.Command)
	if err != nil {
		return map[string]interface{}{
			"success": false,
			"error":   err.Error(),
			"output":  output,
		}, nil
	}

	return map[string]interface{}{
		"success": true,
		"output":  output,
	}, nil
}

// handleGetContexts gets available Kubernetes contexts
func (s *Server) handleGetContexts(ctx context.Context, params json.RawMessage) (interface{}, error) {
	contexts, err := s.k8sExec.GetContexts(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get contexts: %w", err)
	}

	return map[string]interface{}{
		"contexts": contexts,
	}, nil
}

// handleCurrentContext gets the current Kubernetes context
func (s *Server) handleCurrentContext(ctx context.Context, params json.RawMessage) (interface{}, error) {
	context, err := s.k8sExec.GetCurrentContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get current context: %w", err)
	}

	return map[string]interface{}{
		"context": strings.TrimSpace(context),
	}, nil
}

// handleSetContext sets the current Kubernetes context
func (s *Server) handleSetContext(ctx context.Context, params json.RawMessage) (interface{}, error) {
	// Parse parameters
	var p struct {
		Context string `json:"context"`
	}
	if err := json.Unmarshal(params, &p); err != nil {
		return nil, fmt.Errorf("invalid parameters: %w", err)
	}

	// Validate context
	if strings.TrimSpace(p.Context) == "" {
		return nil, fmt.Errorf("context cannot be empty")
	}

	// Set the context
	err := s.k8sExec.SetContext(ctx, p.Context)
	if err != nil {
		return nil, fmt.Errorf("failed to set context: %w", err)
	}

	return map[string]interface{}{
		"success": true,
	}, nil
}
