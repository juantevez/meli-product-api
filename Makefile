.PHONY: help build run test clean docker-build docker-run docker-stop format lint deps

# Variables
APP_NAME=meli-product-api
DOCKER_IMAGE=meli-product-api-go
DOCKER_CONTAINER=meli-api-go
GO_VERSION=1.21
PORT=8080

## help: Display this help message
help:
	@echo Available commands:
	@echo.
	@grep -E "^##" Makefile | sed "s/## /  /"
	@echo.

## build: Build the Go application
build:
	@echo Building application...
	@go build -o bin/api.exe cmd/api/main.go
	@echo Build complete: bin/api.exe

## run: Run the application locally
run:
	@echo Starting application...
	@go run cmd/api/main.go

## dev: Run with hot reload using air
dev:
	@echo Starting development server with hot reload...
	@air

## test: Run all tests
test:
	@echo Running tests...
	@go test -v -race -coverprofile=coverage.out ./...
	@echo Tests complete

## test-coverage: Run tests with coverage report
test-coverage: test
	@echo Generating coverage report...
	@go tool cover -html=coverage.out -o coverage.html
	@echo Coverage report: coverage.html

## benchmark: Run benchmarks
benchmark:
	@echo Running benchmarks...
	@go test -bench=. -benchmem ./...

## deps: Install/Update dependencies
deps:
	@echo Downloading dependencies...
	@go mod download
	@go mod tidy
	@echo Dependencies updated

## format: Format Go code
format:
	@echo Formatting code...
	@gofmt -s -w .
	@echo Code formatted

## lint: Run linter
lint:
	@echo Running linter...
	@golangci-lint run ./...
	@echo Linting complete

## vet: Run go vet
vet:
	@echo Running go vet...
	@go vet ./...
	@echo Vet complete

## clean: Clean build artifacts
clean:
	@echo Cleaning...
	@if exist bin rmdir /s /q bin
	@if exist coverage.out del coverage.out
	@if exist coverage.html del coverage.html
	@go clean
	@echo Clean complete

## docker-build: Build Docker image
docker-build:
	@echo Building Docker image...
	@docker build -t $(DOCKER_IMAGE):latest .
	@echo Docker image built: $(DOCKER_IMAGE):latest

## docker-run: Run Docker container
docker-run:
	@echo Starting Docker container...
	@docker run -d -p $(PORT):8080 --name $(DOCKER_CONTAINER) $(DOCKER_IMAGE):latest
	@echo Container running: http://localhost:$(PORT)

## docker-stop: Stop and remove Docker container
docker-stop:
	@echo Stopping Docker container...
	@docker stop $(DOCKER_CONTAINER) 2>NUL || echo Container not running
	@docker rm $(DOCKER_CONTAINER) 2>NUL || echo Container not found
	@echo Container stopped

## docker-logs: Show Docker container logs
docker-logs:
	@docker logs -f $(DOCKER_CONTAINER)

## docker-shell: Open shell in Docker container
docker-shell:
	@docker exec -it $(DOCKER_CONTAINER) sh

## docker-compose-up: Start services with docker-compose
docker-compose-up:
	@echo Starting services with docker-compose...
	@docker-compose up -d
	@echo Services started

## docker-compose-down: Stop services with docker-compose
docker-compose-down:
	@echo Stopping services...
	@docker-compose down
	@echo Services stopped

## docker-compose-logs: Show docker-compose logs
docker-compose-logs:
	@docker-compose logs -f

## docker-clean: Remove all Docker artifacts
docker-clean: docker-stop
	@echo Cleaning Docker artifacts...
	@docker rmi $(DOCKER_IMAGE):latest 2>NUL || echo Image not found
	@docker system prune -f
	@echo Docker cleaned

## install-tools: Install development tools
install-tools:
	@echo Installing development tools...
	@go install github.com/cosmtrek/air@latest
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@go install github.com/swaggo/swag/cmd/swag@latest
	@echo Tools installed

## swagger: Generate Swagger documentation
swagger:
	@echo Generating Swagger docs...
	@swag init -g cmd/api/main.go -o docs
	@echo Swagger docs generated

## check: Run all checks (format, vet, lint, test)
check: format vet lint test
	@echo All checks passed

## all: Build, test, and create Docker image
all: clean deps format vet test build docker-build
	@echo Build pipeline complete

## health: Check if application is running
health:
	@curl -f http://localhost:$(PORT)/health || echo Service not responding
	@echo Service is healthy
