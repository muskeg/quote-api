# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Copy Go module files
COPY go.mod go.sum ./
RUN go mod download

# Copy files
COPY *.go ./
COPY quotes-default.json ./
COPY entry.sh ./
COPY config-default.yaml ./

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o quote-api .

# Run stage
FROM alpine:latest

WORKDIR /app

# Install CA certificates for HTTPS
RUN apk --no-cache add ca-certificates

# Create non-root user
RUN adduser -D -g '' appuser
USER appuser

# Copy binary from builder
COPY --chown=appuser:appuser --from=builder /app/quote-api .
COPY --chown=appuser:appuser --from=builder /app/entry.sh .
COPY --chown=appuser:appuser --from=builder /app/config-default.yaml .
COPY --chown=appuser:appuser --from=builder /app/quotes-default.json .

# Set environment variables
ENV PORT=8080
ENV GIN_MODE=release

# Make entrypoint script executable
RUN chmod +x /app/entry.sh

# Declare persistent data directory
VOLUME ["/app/data"]

# Set entrypoint
ENTRYPOINT ["/app/entry.sh"]
