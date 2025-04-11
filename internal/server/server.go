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
	"github.com/gorilla/websocket"
	mcp "github.com/modelcontextprotocol/sdk/go"
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
	mcpServer, err := s.createMCPServer()
	if err != nil {
		return nil, fmt.Errorf("failed to create MCP server: %w", err)
	}
	s.mcpServer = mcpServer

	return s, nil
}

// createMCPServer sets up the MCP server with handlers
func (s *Server) createMCPServer() (*mcp.Server, error) {
	// Configure MCP server
	mcpServer := mcp.NewServer("mcp-kubernetes", &mcp.ServerOptions{
		LogLevel: mcp.LogLevelInfo,
	})

	// Register handlers
	mcpServer.RegisterTool("execute", "Execute a kubectl command", map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"command": {
				"type":        "string",
				"description": "The kubectl command to execute (without the 'kubectl' prefix)",
			},
		},
		"required": []string{"command"},
	}, s.handleExecute)

	mcpServer.RegisterTool("get-contexts", "Get available Kubernetes contexts", map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{},
	}, s.handleGetContexts)

	mcpServer.RegisterTool("current-context", "Get the current Kubernetes context", map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{},
	}, s.handleCurrentContext)

	mcpServer.RegisterTool("set-context", "Set the current Kubernetes context", map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"context": {
				"type":        "string",
				"description": "The context to set as current",
			},
		},
		"required": []string{"context"},
	}, s.handleSetContext)

	return mcpServer, nil
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

	// Create MCP connection
	mcpConn := mcp.NewConnection(conn)
	
	// Handle the MCP connection
	if err := s.mcpServer.HandleConnection(mcpConn); err != nil {
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
			"content": []map[string]interface{}{
				{
					"type": "text",
					"text": fmt.Sprintf("Error: %s\nOutput: %s", err.Error(), output),
				},
			},
			"isError": true,
		}, nil
	}

	return map[string]interface{}{
		"content": []map[string]interface{}{
			{
				"type": "text",
				"text": output,
			},
		},
		"isError": false,
	}, nil
}

// handleGetContexts gets available Kubernetes contexts
func (s *Server) handleGetContexts(ctx context.Context, params json.RawMessage) (interface{}, error) {
	contexts, err := s.k8sExec.GetContexts(ctx)
	if err != nil {
		return map[string]interface{}{
			"content": []map[string]interface{}{
				{
					"type": "text",
					"text": fmt.Sprintf("Failed to get contexts: %s", err.Error()),
				},
			},
			"isError": true,
		}, nil
	}

	return map[string]interface{}{
		"content": []map[string]interface{}{
			{
				"type": "text",
				"text": fmt.Sprintf("Available contexts:\n%s", strings.Join(contexts, "\n")),
			},
		},
		"isError": false,
		"metadata": map[string]interface{}{
			"contexts": contexts,
		},
	}, nil
}

// handleCurrentContext gets the current Kubernetes context
func (s *Server) handleCurrentContext(ctx context.Context, params json.RawMessage) (interface{}, error) {
	currentContext, err := s.k8sExec.GetCurrentContext(ctx)
	if err != nil {
		return map[string]interface{}{
			"content": []map[string]interface{}{
				{
					"type": "text",
					"text": fmt.Sprintf("Failed to get current context: %s", err.Error()),
				},
			},
			"isError": true,
		}, nil
	}

	return map[string]interface{}{
		"content": []map[string]interface{}{
			{
				"type": "text",
				"text": fmt.Sprintf("Current context: %s", strings.TrimSpace(currentContext)),
			},
		},
		"isError": false,
		"metadata": map[string]interface{}{
			"context": strings.TrimSpace(currentContext),
		},
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
		"content": []map[string]interface{}{
			{
				"type": "text",
				"text": fmt.Sprintf("Successfully switched to context '%s'", p.Context),
			},
		},
		"isError": false,
	}, nil
}
