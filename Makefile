.PHONY: help build run test clean docker-build docker-run docker-stop format lint deps

# Variables
APP_NAME=meli-product-api
DOCKER_IMAGE=meli-product-api-go
DOCKER_CONTAINER=meli-api-go
GO_VERSION=1.21
PORT=8080

# Colors for output
GREEN=\033[0;32m
YELLOW=\033[1;33m
RED=\033[0;31m
NC=\033[0m # No Color

## help: Display this help message
help:
	@echo "$(GREEN)Available commands:$(NC)"
	@echo ""
	@grep -E '^##' Makefile | sed 's/## /  /'
	@echo ""

## build: Build the Go application
build:
	@echo "$(YELLOW)Building application...$(NC)"
	@go build -o bin/api cmd/api/main.go
	@echo "$(GREEN)✓ Build complete: bin/api$(NC)"

## run: Run the application locally
run:
	@echo "$(YELLOW)Starting application...$(NC)"
	@go run cmd/api/main.go

## dev: Run with hot reload using air (install: go install github.com/cosmtrek/air@latest)
dev:
	@echo "$(YELLOW)Starting development server with hot reload...$(NC)"
	@air

## test: Run all tests
test:
	@echo "$(YELLOW)Running tests...$(NC)"
	@go test -v -race -coverprofile=coverage.out ./...
	@echo "$(GREEN)✓ Tests complete$(NC)"

## test-coverage: Run tests with coverage report
test-coverage: test
	@echo "$(YELLOW)Generating coverage report...$(NC)"
	@go tool cover -html=coverage.out -o coverage.html
	@echo "$(GREEN)✓ Coverage report: coverage.html$(NC)"

## benchmark: Run benchmarks
benchmark:
	@echo "$(YELLOW)Running benchmarks...$(NC)"
	@go test -bench=. -benchmem ./...

## deps: Install/Update dependencies
deps:
	@echo "$(YELLOW)Downloading dependencies...$(NC)"
	@go mod download
	@go mod tidy
	@echo "$(GREEN)✓ Dependencies updated$(NC)"

## format: Format Go code
format:
	@echo "$(YELLOW)Formatting code...$(NC)"
	@gofmt -s -w .
	@echo "$(GREEN)✓ Code formatted$(NC)"

## lint: Run linter (install: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest)
lint:
	@echo "$(YELLOW)Running linter...$(NC)"
	@golangci-lint run ./...
	@echo "$(GREEN)✓ Linting complete$(NC)"

## vet: Run go vet
vet:
	@echo "$(YELLOW)Running go vet...$(NC)"
	@go vet ./...
	@echo "$(GREEN)✓ Vet complete$(NC)"

## clean: Clean build artifacts
clean:
	@echo "$(YELLOW)Cleaning...$(NC)"
	@rm -rf bin/
	@rm -f coverage.out coverage.html
	@go clean
	@echo "$(GREEN)✓ Clean complete$(NC)"

## docker-build: Build Docker image
docker-build:
	@echo "$(YELLOW)Building Docker image...$(NC)"
	@docker build -t $(DOCKER_IMAGE):latest .
	@echo "$(GREEN)✓ Docker image built: $(DOCKER_IMAGE):latest$(NC)"

## docker-run: Run Docker container
docker-run:
	@echo "$(YELLOW)Starting Docker container...$(NC)"
	@docker run -d \
		-p $(PORT):8080 \
		--name $(DOCKER_CONTAINER) \
		$(DOCKER_IMAGE):latest
	@echo "$(GREEN)✓ Container running: http://localhost:$(PORT)$(NC)"

## docker-stop: Stop and remove Docker container
docker-stop:
	@echo "$(YELLOW)Stopping Docker container...$(NC)"
	@docker stop $(DOCKER_CONTAINER) 2>/dev/null || true
	@docker rm $(DOCKER_CONTAINER) 2>/dev/null || true
	@echo "$(GREEN)✓ Container stopped$(NC)"

## docker-logs: Show Docker container logs
docker-logs:
	@docker logs -f $(DOCKER_CONTAINER)

## docker-shell: Open shell in Docker container
docker-shell:
	@docker exec -it $(DOCKER_CONTAINER) sh

## docker-compose-up: Start services with docker-compose
docker-compose-up:
	@echo "$(YELLOW)Starting services with docker-compose...$(NC)"
	@docker-compose up -d
	@echo "$(GREEN)✓ Services started$(NC)"

## docker-compose-down: Stop services with docker-compose
docker-compose-down:
	@echo "$(YELLOW)Stopping services...$(NC)"
	@docker-compose down
	@echo "$(GREEN)✓ Services stopped$(NC)"

## docker-compose-logs: Show docker-compose logs
docker-compose-logs:
	@docker-compose logs -f

## docker-clean: Remove all Docker artifacts
docker-clean: docker-stop
	@echo "$(YELLOW)Cleaning Docker artifacts...$(NC)"
	@docker rmi $(DOCKER_IMAGE):latest 2>/dev/null || true
	@docker system prune -f
	@echo "$(GREEN)✓ Docker cleaned$(NC)"

## install-tools: Install development tools
install-tools:
	@echo "$(YELLOW)Installing development tools...$(NC)"
	@go install github.com/cosmtrek/air@latest
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@go install github.com/swaggo/swag/cmd/swag@latest
	@echo "$(GREEN)✓ Tools installed$(NC)"

## swagger: Generate Swagger documentation
swagger:
	@echo "$(YELLOW)Generating Swagger docs...$(NC)"
	@swag init -g cmd/api/main.go -o docs
	@echo "$(GREEN)✓ Swagger docs generated$(NC)"

## check: Run all checks (format, vet, lint, test)
check: format vet lint test
	@echo "$(GREEN)✓ All checks passed$(NC)"

## all: Build, test, and create Docker image
all: clean deps format vet test build docker-build
	@echo "$(GREEN)✓ Build pipeline complete$(NC)"

## health: Check if application is running
health:
	@curl -f http://localhost:$(PORT)/health || (echo "$(RED)✗ Service not responding$(NC)" && exit 1)
	@echo "$(GREEN)✓ Service is healthy$(NC)"

## load-test: Simple load test (requires apache bench: apt-get install apache2-utils)
load-test:
	@echo "$(YELLOW)Running load test...$(NC)"
	@ab -n 1000 -c 10 http://localhost:$(PORT)/api/v1/products/MLA123456
	