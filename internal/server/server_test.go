package server

import (
	"context"
	"encoding/json"
	"testing"
)

func TestCreateMCPServer(t *testing.T) {
	// Create a server instance with a random port
	srv, err := NewServer(0)
	if err != nil {
		t.Fatalf("Failed to create server: %v", err)
	}
	
	// Ensure the MCP server was created
	if srv.mcpServer == nil {
		t.Fatalf("MCP server was not created")
	}
	
	// Test that handlers were registered
	handlers := []string{"execute", "get-contexts", "current-context", "set-context"}
	for _, h := range handlers {
		if !srv.mcpServer.HasHandler(h) {
			t.Errorf("Handler '%s' was not registered", h)
		}
	}
}

// TestHandleExecute mocks handling an execute request
func TestHandleExecute(t *testing.T) {
	// Create a server instance
	srv, err := NewServer(0)
	if err != nil {
		t.Skip("Skipping test: could not create server:", err)
		return
	}
	
	// Create parameters
	params, err := json.Marshal(map[string]interface{}{
		"command": "version --client",
	})
	if err != nil {
		t.Fatalf("Failed to marshal parameters: %v", err)
	}
	
	// Call the handler
	result, err := srv.handleExecute(context.Background(), params)
	if err != nil {
		t.Fatalf("Handler returned error: %v", err)
	}
	
	// Check the result
	resultMap, ok := result.(map[string]interface{})
	if !ok {
		t.Fatalf("Result is not a map: %v", result)
	}
	
	success, ok := resultMap["success"].(bool)
	if !ok || !success {
		t.Errorf("Expected success=true, got %v", resultMap["success"])
	}
}
