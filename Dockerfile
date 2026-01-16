# ============================================
# STAGE 1: Build
# ============================================
FROM golang:1.21-alpine AS builder

# Metadata
LABEL maintainer="Juan - Backend Developer"
LABEL description="MELI Product API - Build Stage"

# Install build dependencies
RUN apk add --no-cache git make

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build application
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-w -s" \
    -o /app/bin/api \
    ./cmd/api

# ============================================
# STAGE 2: Runtime
# ============================================
FROM alpine:latest

# Metadata
LABEL maintainer="Juan - Backend Developer"
LABEL description="MELI Product API - Production Image"
LABEL version="1.0.0"

# Install runtime dependencies
RUN apk --no-cache add ca-certificates tzdata curl

# Create app user (security best practice)
RUN addgroup -S app && adduser -S app -G app

# Set working directory
WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/bin/api .

# Copy data files
COPY --from=builder /app/data ./data

# Change ownership
RUN chown -R app:app /app

# Switch to non-root user
USER app

# Expose port
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=40s --retries=3 \
  CMD curl -f http://localhost:8080/health || exit 1

# Run application
ENTRYPOINT ["./api"]
