#!/bin/bash

# Setup script for MELI Product API

set -e

echo "🚀 MELI Product API - Setup Script"
echo "==================================="
echo ""

# Check Go installation
if ! command -v go &> /dev/null; then
    echo "❌ Go is not installed. Please install Go 1.21+ first."
    exit 1
fi

echo "✓ Go version: $(go version)"

# Install dependencies
echo ""
echo "📦 Installing dependencies..."
go mod download
go mod tidy

# Install development tools
echo ""
echo "🔧 Installing development tools..."
go install github.com/cosmtrek/air@latest
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Build application
echo ""
echo "🔨 Building application..."
go build -o bin/api cmd/api/main.go

# Run tests
echo ""
echo "🧪 Running tests..."
go test ./...

echo ""
echo "✅ Setup complete!"
echo ""
echo "Next steps:"
echo "  • Run locally:        make run"
echo "  • Run with Docker:    make docker-compose-up"
echo "  • View all commands:  make help"
echo ""
