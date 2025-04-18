package kubernetes

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// Executor handles the execution of Kubernetes commands
type Executor struct {
	kubectlPath string
	kubeconfig  string
}

// NewExecutor creates a new Kubernetes executor
func NewExecutor() (*Executor, error) {
	// Find kubectl in PATH
	kubectlPath, err := exec.LookPath("kubectl")
	if err != nil {
		return nil, fmt.Errorf("kubectl not found in PATH: %w", err)
	}

	// Get kubeconfig from environment or use default
	kubeconfig := os.Getenv("KUBECONFIG")
	if kubeconfig == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return nil, fmt.Errorf("unable to determine user home directory: %w", err)
		}
		kubeconfig = fmt.Sprintf("%s/.kube/config", homeDir)
	}

	// Check if kubeconfig exists
	if _, err := os.Stat(kubeconfig); os.IsNotExist(err) {
		return nil, fmt.Errorf("kubeconfig not found at %s: %w", kubeconfig, err)
	}

	return &Executor{
		kubectlPath: kubectlPath,
		kubeconfig:  kubeconfig,
	}, nil
}

// Execute runs a kubectl command and returns the output
func (e *Executor) Execute(ctx context.Context, command string) (string, error) {
	// Parse and sanitize the command
	args := parseCommand(command)
	if len(args) == 0 {
		return "", fmt.Errorf("invalid kubectl command")
	}

	// Create command with context
	cmd := exec.CommandContext(ctx, e.kubectlPath, args...)
	cmd.Env = append(os.Environ(), fmt.Sprintf("KUBECONFIG=%s", e.kubeconfig))

	// Capture stdout and stderr
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	// Execute the command
	err := cmd.Run()
	if err != nil {
		// Return both stdout and stderr if there's an error
		if stderr.Len() > 0 {
			return stderr.String(), fmt.Errorf("command execution failed: %w", err)
		}
		return stdout.String(), err
	}

	return stdout.String(), nil
}

// parseCommand splits a command string into arguments
// This is a simple implementation and doesn't handle all shell parsing edge cases
func parseCommand(command string) []string {
	// Remove "kubectl" from the beginning if present
	command = strings.TrimSpace(command)
	if strings.HasPrefix(command, "kubectl") {
		command = strings.TrimSpace(strings.TrimPrefix(command, "kubectl"))
	}

	// Split by spaces, respecting quotes
	var args []string
	inQuote := false
	var currentArg strings.Builder
	
	for _, char := range command {
		if char == '"' || char == '\'' {
			inQuote = !inQuote
			continue
		}
		
		if char == ' ' && !inQuote {
			if currentArg.Len() > 0 {
				args = append(args, currentArg.String())
				currentArg.Reset()
			}
			continue
		}
		
		currentArg.WriteRune(char)
	}
	
	if currentArg.Len() > 0 {
		args = append(args, currentArg.String())
	}
	
	return args
}

// GetContexts returns a list of available Kubernetes contexts
func (e *Executor) GetContexts(ctx context.Context) ([]string, error) {
	output, err := e.Execute(ctx, "config get-contexts -o name")
	if err != nil {
		return nil, err
	}
	
	contexts := strings.Split(strings.TrimSpace(output), "\n")
	return contexts, nil
}

// GetCurrentContext returns the current Kubernetes context
func (e *Executor) GetCurrentContext(ctx context.Context) (string, error) {
	return e.Execute(ctx, "config current-context")
}

// SetContext changes the current Kubernetes context
func (e *Executor) SetContext(ctx context.Context, contextName string) error {
	_, err := e.Execute(ctx, fmt.Sprintf("config use-context %s", contextName))
	return err
}
