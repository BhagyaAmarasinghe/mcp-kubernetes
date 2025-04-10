package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/BhagyaAmarasinghe/mcp-kubernetes/internal/server"
)

func main() {
	// Parse command line flags
	port := flag.Int("port", 3000, "Port to run the server on")
	allowedCommands := flag.String("allowed-commands", "*", "Comma-separated list of allowed kubectl commands, or * for all commands")
	flag.Parse()

	// Parse allowed commands
	var allowedCommandsList []string
	if *allowedCommands != "*" {
		allowedCommandsList = strings.Split(*allowedCommands, ",")
		// Trim spaces
		for i, cmd := range allowedCommandsList {
			allowedCommandsList[i] = strings.TrimSpace(cmd)
		}
		log.Printf("Allowed kubectl commands: %v", allowedCommandsList)
	} else {
		log.Printf("All kubectl commands are allowed")
	}

	// Initialize and start the MCP server
	mcpServer, err := server.NewServer(*port, allowedCommandsList)
	if err != nil {
		log.Fatalf("Failed to create server: %v", err)
	}

	// Start the server in a goroutine
	go func() {
		if err := mcpServer.Start(); err != nil {
			log.Fatalf("Server error: %v", err)
		}
	}()

	fmt.Printf("MCP Kubernetes server running on port %d\n", *port)
	fmt.Println("Press Ctrl+C to exit")

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	fmt.Println("Shutting down server...")
	if err := mcpServer.Stop(); err != nil {
		log.Fatalf("Error during server shutdown: %v", err)
	}
	fmt.Println("Server stopped")
}
