# Build stage
FROM golang:1.23 AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application with CGO enabled
ENV CGO_ENABLED=1
RUN go build -o mini-ledger ./cmd/server

# Runtime stage
FROM debian:bookworm-slim

WORKDIR /app

# Install runtime dependencies
RUN apt-get update && apt-get install -y ca-certificates curl && rm -rf /var/lib/apt/lists/*

# Copy the binary from builder stage
COPY --from=builder /app/mini-ledger .

# Expose port
EXPOSE 8080

# Run the application
CMD ["./mini-ledger"]