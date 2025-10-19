# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Copy Go module files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY *.go ./
COPY quotes.json ./

RUN cat go.mod

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
COPY --from=builder /app/quote-api .
COPY --from=builder /app/quotes.json .

# Expose API port
EXPOSE 8080

# Run the application
CMD ["./quote-api"]