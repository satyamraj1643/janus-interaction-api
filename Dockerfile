# Build stage
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache gcc musl-dev

# Copy dependency files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN go build -o janus-api main.go

# Run stage
FROM alpine:latest

# Install runtime dependencies (e.g., for healthcheck)
RUN apk add --no-cache curl

WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/janus-api .

# Expose the application port
EXPOSE 8080

# Environment variables (can be overridden at runtime)
ENV SERVER_PORT=8080

# Healthcheck
HEALTHCHECK --interval=30s --timeout=5s --start-period=5s --retries=3 \
    CMD curl -f http://localhost:${SERVER_PORT}/health || exit 1

# Run the application
CMD ["./janus-api"]
