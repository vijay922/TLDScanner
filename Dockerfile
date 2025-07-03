# Multi-stage build for TLD Scanner
FROM golang:1.21-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git ca-certificates tzdata

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY tldscanner.go ./

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags="-w -s" -o tldscanner .

# Final stage
FROM alpine:latest

# Install runtime dependencies
RUN apk --no-cache add ca-certificates tzdata

# Create non-root user
RUN addgroup -g 1001 scanner && \
    adduser -D -u 1001 -G scanner scanner

# Set working directory
WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/tldscanner .

# Copy wordlist
COPY wordlist.txt .

# Set ownership
RUN chown -R scanner:scanner /app

# Switch to non-root user
USER scanner

# Set entrypoint
ENTRYPOINT ["./tldscanner"]

# Default command
CMD ["-h"]

# Metadata
LABEL maintainer="TLD Scanner Team"
LABEL version="2.0.0"
LABEL description="High-performance domain enumeration tool"
LABEL org.opencontainers.image.title="TLD Scanner"
LABEL org.opencontainers.image.description="Domain enumeration tool for security research"
LABEL org.opencontainers.image.version="2.0.0"
LABEL org.opencontainers.image.vendor="TLD Scanner Team"
LABEL org.opencontainers.image.licenses="MIT"
