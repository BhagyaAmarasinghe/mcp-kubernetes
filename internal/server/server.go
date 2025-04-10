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
	port            int
	k8sExec         *kubernetes.Executor
	server          *http.Server
	upgrader        websocket.Upgrader
	mcpServer       *mcp.Server
	allowedCommands []string    // List of allowed commands, empty means all commands are allowed
	recentCommands  []CmdRecord // List of recently executed commands
}

// CmdRecord represents a record of an executed command
type CmdRecord struct {
	Command   string    `json:"command"`
	Timestamp time.Time `json:"timestamp"`
	Success   bool      `json:"success"`
	Error     string    `json:"error,omitempty"`
}

// NewServer creates a new MCP server for Kubernetes
func NewServer(port int, allowedCommands []string) (*Server, error) {
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
		port:            port,
		k8sExec:         k8sExec,
		upgrader:        upgrader,
		allowedCommands: allowedCommands,
		recentCommands:  make([]CmdRecord, 0, 100), // Pre-allocate space for 100 commands
	}

	// Create and configure MCP server
	mcpServer := mcp.NewServer("mcp-kubernetes", "MCP server for executing Kubernetes commands", "1.0.0")
	
	// Register handlers
	mcpServer.RegisterRequestHandler("execute", s.handleExecute)
	mcpServer.RegisterRequestHandler("get-contexts", s.handleGetContexts)
	mcpServer.RegisterRequestHandler("current-context", s.handleCurrentContext)
	mcpServer.RegisterRequestHandler("set-context", s.handleSetContext)
	mcpServer.RegisterRequestHandler("list-recent-commands", s.handleListRecentCommands)
	mcpServer.RegisterRequestHandler("list-allowed-commands", s.handleListAllowedCommands)
	
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

// isCommandAllowed checks if a command is allowed to be executed
func (s *Server) isCommandAllowed(command string) bool {
	// If allowedCommands is empty, all commands are allowed
	if len(s.allowedCommands) == 0 {
		return true
	}

	// Extract the main command (first word)
	parts := strings.Fields(command)
	if len(parts) == 0 {
		return false
	}
	mainCommand := parts[0]

	// Check if the main command is in the allowed list
	for _, allowed := range s.allowedCommands {
		if allowed == mainCommand {
			return true
		}
	}

	return false
}

// recordCommand adds a command to the recent commands list
func (s *Server) recordCommand(command string, success bool, errMsg string) {
	// Create a new record
	record := CmdRecord{
		Command:   command,
		Timestamp: time.Now(),
		Success:   success,
		Error:     errMsg,
	}

	// Add to the beginning of the list
	s.recentCommands = append([]CmdRecord{record}, s.recentCommands...)

	// Keep only the last 100 commands
	if len(s.recentCommands) > 100 {
		s.recentCommands = s.recentCommands[:100]
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
	command := strings.TrimSpace(p.Command)
	if command == "" {
		return nil, fmt.Errorf("command cannot be empty")
	}

	// Check if command is allowed
	if !s.isCommandAllowed(command) {
		errMsg := fmt.Sprintf("command '%s' is not allowed", command)
		s.recordCommand(command, false, errMsg)
		return nil, fmt.Errorf(errMsg)
	}

	// Execute the command
	startTime := time.Now()
	output, err := s.k8sExec.Execute(ctx, command)
	executionTime := time.Since(startTime)

	// Record the command
	if err != nil {
		s.recordCommand(command, false, err.Error())
		return map[string]interface{}{
			"success":        false,
			"error":          err.Error(),
			"output":         output,
			"execution_time": executionTime.String(),
		}, nil
	}

	s.recordCommand(command, true, "")
	return map[string]interface{}{
		"success":        true,
		"output":         output,
		"execution_time": executionTime.String(),
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
	context := strings.TrimSpace(p.Context)
	if context == "" {
		return nil, fmt.Errorf("context cannot be empty")
	}

	// Set the context
	err := s.k8sExec.SetContext(ctx, context)
	if err != nil {
		return nil, fmt.Errorf("failed to set context: %w", err)
	}

	return map[string]interface{}{
		"success": true,
	}, nil
}

// handleListRecentCommands lists recently executed commands
func (s *Server) handleListRecentCommands(ctx context.Context, params json.RawMessage) (interface{}, error) {
	// Parse parameters
	var p struct {
		Limit int `json:"limit"`
	}
	if err := json.Unmarshal(params, &p); err != nil {
		// If parameters can't be parsed, use default limit
		p.Limit = 10
	}

	// Default limit is 10 if not specified or invalid
	if p.Limit <= 0 {
		p.Limit = 10
	}

	// Limit to the number of available commands
	if p.Limit > len(s.recentCommands) {
		p.Limit = len(s.recentCommands)
	}

	return map[string]interface{}{
		"commands": s.recentCommands[:p.Limit],
	}, nil
}

// handleListAllowedCommands lists allowed commands
func (s *Server) handleListAllowedCommands(ctx context.Context, params json.RawMessage) (interface{}, error) {
	if len(s.allowedCommands) == 0 {
		return map[string]interface{}{
			"allowed_commands": "*", // All commands are allowed
		}, nil
	}

	return map[string]interface{}{
		"allowed_commands": s.allowedCommands,
	}, nil
}
