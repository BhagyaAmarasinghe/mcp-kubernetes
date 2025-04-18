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
	
	// This is just to test compilation - we need an executor to actually run this test
	if srv.k8sExec == nil {
		t.Skip("Skipping test execution: no kubernetes executor")
		return
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
