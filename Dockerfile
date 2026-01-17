# ============================================
# STAGE 1: Build
# ============================================
FROM golang:1.24-alpine AS builder

LABEL maintainer="Juan - Backend Developer"
LABEL description="MELI Product API - Build Stage"

RUN apk add --no-cache git make

WORKDIR /app

# Copy go modules
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build application (output to /app/api directamente)
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-w -s" \
    -o api \
    ./cmd/api/main.go

# ============================================
# STAGE 2: Runtime
# ============================================
FROM alpine:latest

LABEL maintainer="Juan - Backend Developer"
LABEL description="MELI Product API - Production Image"
LABEL version="1.0.0"

# Install dependencies
RUN apk --no-cache add ca-certificates tzdata curl

# Create non-root user
RUN addgroup -S app && adduser -S app -G app

WORKDIR /app

# Copy binary from builder (ahora está en /app/api)
COPY --from=builder --chmod=755 /app/api .

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
