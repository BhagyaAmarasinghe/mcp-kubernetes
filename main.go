package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/BhagyaAmarasinghe/mcp-kubernetes/internal/kubernetes"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func main() {
	// Parse command line flags
	port := flag.Int("port", 3000, "Port to run the server on")
	flag.Parse()

	// Initialize Kubernetes executor
	kubeExecutor, err := kubernetes.NewExecutor()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize Kubernetes executor: %v\n", err)
		os.Exit(1)
	}

	// Create MCP server
	s := server.NewMCPServer(
		"MCP Kubernetes",
		"1.0.0",
	)

	// Add execute command tool
	executeTool := mcp.NewTool("execute",
		mcp.WithDescription("Execute a kubectl command"),
		mcp.WithString("command",
			mcp.Required(),
			mcp.Description("The kubectl command to execute (without 'kubectl' prefix)"),
		),
	)

	// Add execute tool handler
	s.AddTool(executeTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		command, ok := request.Params.Arguments["command"].(string)
		if !ok {
			return nil, errors.New("command must be a string")
		}

		output, err := kubeExecutor.Execute(ctx, command)
		if err != nil {
			return mcp.NewToolResultText(fmt.Sprintf("Error: %v\nOutput: %s", err, output)), nil
		}

		return mcp.NewToolResultText(output), nil
	})

	// Add get-contexts tool
	getContextsTool := mcp.NewTool("get-contexts",
		mcp.WithDescription("Get available Kubernetes contexts"),
	)

	// Add get-contexts tool handler
	s.AddTool(getContextsTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		contexts, err := kubeExecutor.GetContexts(ctx)
		if err != nil {
			return mcp.NewToolResultText(fmt.Sprintf("Error getting contexts: %v", err)), nil
		}

		result := "Available contexts:\n"
		for _, ctx := range contexts {
			result += fmt.Sprintf("- %s\n", ctx)
		}
		return mcp.NewToolResultText(result), nil
	})

	// Add current-context tool
	currentContextTool := mcp.NewTool("current-context",
		mcp.WithDescription("Get the current Kubernetes context"),
	)

	// Add current-context tool handler
	s.AddTool(currentContextTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		context, err := kubeExecutor.GetCurrentContext(ctx)
		if err != nil {
			return mcp.NewToolResultText(fmt.Sprintf("Error getting current context: %v", err)), nil
		}

		return mcp.NewToolResultText(fmt.Sprintf("Current context: %s", context)), nil
	})

	// Add set-context tool
	setContextTool := mcp.NewTool("set-context",
		mcp.WithDescription("Set the current Kubernetes context"),
		mcp.WithString("context",
			mcp.Required(),
			mcp.Description("Name of the context to set as current"),
		),
	)

	// Add set-context tool handler
	s.AddTool(setContextTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		contextName, ok := request.Params.Arguments["context"].(string)
		if !ok {
			return nil, errors.New("context must be a string")
		}

		err := kubeExecutor.SetContext(ctx, contextName)
		if err != nil {
			return mcp.NewToolResultText(fmt.Sprintf("Error setting context: %v", err)), nil
		}

		return mcp.NewToolResultText(fmt.Sprintf("Successfully switched to context: %s", contextName)), nil
	})

	// Setup signal handling for graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Start the server based on the port flag
	if *port > 0 {
		fmt.Printf("Starting MCP Kubernetes server on HTTP port %d...\n", *port)
		
		// Start in a goroutine so we can handle shutdown
		go func() {
			if err := s.ListenAndServe(fmt.Sprintf(":%d", *port)); err != nil {
				log.Fatalf("Server error: %v", err)
			}
		}()
	} else {
		fmt.Println("Starting MCP Kubernetes server with stdio transport...")
		
		// Run in a goroutine so we can handle shutdown
		go func() {
			if err := server.ServeStdio(s); err != nil {
				log.Fatalf("Server error: %v", err)
			}
		}()
	}

	// Wait for shutdown signal
	<-quit
	fmt.Println("Shutting down server...")
	fmt.Println("Server stopped")
}
