package kubernetes

import (
	"context"
	"strings"
	"testing"
	"time"
)

func TestParseCommand(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []string
	}{
		{
			name:     "simple command",
			input:    "get pods",
			expected: []string{"get", "pods"},
		},
		{
			name:     "command with kubectl prefix",
			input:    "kubectl get pods",
			expected: []string{"get", "pods"},
		},
		{
			name:     "command with namespace",
			input:    "get pods -n kube-system",
			expected: []string{"get", "pods", "-n", "kube-system"},
		},
		{
			name:     "command with quoted argument",
			input:    "get pods -l \"app=nginx\"",
			expected: []string{"get", "pods", "-l", "app=nginx"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parseCommand(tt.input)
			
			if len(result) != len(tt.expected) {
				t.Errorf("Expected %d args, got %d: %v", len(tt.expected), len(result), result)
				return
			}
			
			for i, arg := range result {
				if arg != tt.expected[i] {
					t.Errorf("Arg %d: expected '%s', got '%s'", i, tt.expected[i], arg)
				}
			}
		})
	}
}

// TestExecutor tests the executor if kubectl is available
// This test is skipped if kubectl is not installed
func TestExecutor(t *testing.T) {
	// Create executor
	exec, err := NewExecutor()
	if err != nil {
		t.Skip("Skipping test: kubectl not available:", err)
	}
	
	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	// Test version command
	output, err := exec.Execute(ctx, "version --client")
	if err != nil {
		t.Fatalf("Failed to execute version command: %v", err)
	}
	
	if !strings.Contains(output, "Client Version") {
		t.Errorf("Expected output to contain 'Client Version', got: %s", output)
	}
}
