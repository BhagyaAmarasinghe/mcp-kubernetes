FROM golang:1.20-alpine AS builder
WORKDIR /app

# Install required dependencies
RUN apk add --no-cache git

# Copy go.mod and go.sum first for caching dependencies
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN go build -ldflags="-s -w" -o mcp-kubernetes ./cmd

FROM alpine:latest
WORKDIR /app

# Install kubectl
RUN apk add --no-cache curl && \
    curl -LO "https://dl.k8s.io/release/$(curl -L -s https://dl.k8s.io/release/stable.txt)/bin/linux/amd64/kubectl" && \
    chmod +x kubectl && \
    mv kubectl /usr/local/bin/

# Copy the built binary from the builder stage
COPY --from=builder /app/mcp-kubernetes ./

# Run with a default configuration
ENTRYPOINT ["./mcp-kubernetes"]
CMD ["--allowed-contexts=*"]
