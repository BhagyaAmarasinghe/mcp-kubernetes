FROM golang:1.23-alpine AS builder

WORKDIR /app

# Copy go.mod and go.sum files first to leverage Docker caching
COPY go.mod go.sum* ./
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o mcp-kubernetes

# Use a small image for the final container
FROM alpine:latest

# Install kubectl
RUN apk add --no-cache curl && \
    curl -LO "https://dl.k8s.io/release/$(curl -L -s https://dl.k8s.io/release/stable.txt)/bin/linux/amd64/kubectl" && \
    chmod +x kubectl && \
    mv kubectl /usr/local/bin/

# Copy the binary from the builder stage
COPY --from=builder /app/mcp-kubernetes /usr/local/bin/

# Create a directory for Kubernetes config
RUN mkdir -p /root/.kube

# The kubeconfig should be mounted at runtime
VOLUME ["/root/.kube"]

# Run the server
ENTRYPOINT ["mcp-kubernetes"]
CMD ["-port", "8080"]
