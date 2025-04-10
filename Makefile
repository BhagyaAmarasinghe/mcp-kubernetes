.PHONY: build test clean run

# Binary name
BINARY_NAME=mcp-kubernetes

# Build the application
build:
	go build -o bin/$(BINARY_NAME) ./cmd

# Run tests
test:
	go test -v ./...

# Clean build artifacts
clean:
	rm -rf bin/

# Run the application
run: build
	./bin/$(BINARY_NAME)

# Install the application
install:
	go install ./cmd

# Format code
fmt:
	go fmt ./...

# Run linter
lint:
	go vet ./...
