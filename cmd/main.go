package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/bhagya/mcp-kubernetes/internal/server"
)

func main() {
	// Parse command line flags
	port := flag.Int("port", 3000, "Port to run the server on")
	flag.Parse()

	// Initialize and start the MCP server
	mcpServer, err := server.NewServer(*port)
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
