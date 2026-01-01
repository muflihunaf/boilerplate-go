# =============================================================================
# Boilerplate Go - Makefile
# =============================================================================

.PHONY: help build run dev test test-coverage clean tidy deps lint vet \
        build-linux docker-build docker-build-distroless docker-run \
        docker-run-env docker-test docker-size swagger swagger-fmt

# Default target
.DEFAULT_GOAL := help

# =============================================================================
# Variables
# =============================================================================

# Application
BINARY_NAME := api
BUILD_DIR := bin
MAIN_PATH := ./cmd/api

# Go
GOCMD := go
GOBUILD := $(GOCMD) build
GORUN := $(GOCMD) run
GOTEST := $(GOCMD) test
GOMOD := $(GOCMD) mod
GOVET := $(GOCMD) vet

# Docker
DOCKER_IMAGE := boilerplate-go
DOCKER_TAG := latest

# Version info (from git)
VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_TIME := $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")

# Build flags
LDFLAGS := -ldflags="-s -w -X main.Version=$(VERSION) -X main.Commit=$(COMMIT) -X main.BuildTime=$(BUILD_TIME)"

# =============================================================================
# Help
# =============================================================================

help: ## Show this help message
	@echo "Boilerplate Go - Available Commands"
	@echo "===================================="
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

# =============================================================================
# Development
# =============================================================================

run: ## Run the application
	$(GORUN) $(MAIN_PATH)

dev: ## Run with hot reload (requires air)
	@air

build: ## Build the application
	@echo "Building..."
	@mkdir -p $(BUILD_DIR)
	$(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PATH)
	@echo "Binary: $(BUILD_DIR)/$(BINARY_NAME)"

build-linux: ## Build for Linux (cross-compile)
	@echo "Building for Linux..."
	@mkdir -p $(BUILD_DIR)
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 $(MAIN_PATH)
	@echo "Binary: $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64"

clean: ## Clean build artifacts
	@echo "Cleaning..."
	@rm -rf $(BUILD_DIR)
	@rm -f coverage.out coverage.html
	@echo "Done."

# =============================================================================
# Dependencies
# =============================================================================

deps: ## Download dependencies
	$(GOMOD) download

tidy: ## Tidy and verify dependencies
	$(GOMOD) tidy
	$(GOMOD) verify

# =============================================================================
# Testing
# =============================================================================

test: ## Run tests
	$(GOTEST) -v -race ./...

test-coverage: ## Run tests with coverage report
	$(GOTEST) -v -race -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report: coverage.html"

# =============================================================================
# Code Quality
# =============================================================================

lint: ## Run linter (requires golangci-lint)
	@golangci-lint run ./...

vet: ## Run go vet
	$(GOVET) ./...

fmt: ## Format code
	@gofmt -s -w .

# =============================================================================
# Swagger
# =============================================================================

swagger: ## Generate Swagger documentation
	@echo "Generating Swagger docs..."
	@go run github.com/swaggo/swag/cmd/swag@latest init -g cmd/api/main.go -o docs --parseDependency --parseInternal
	@echo "Swagger docs generated in docs/"

swagger-fmt: ## Format Swagger annotations
	@go run github.com/swaggo/swag/cmd/swag@latest fmt

# =============================================================================
# Docker
# =============================================================================

docker-build: ## Build Docker image (alpine)
	docker build \
		--build-arg VERSION=$(VERSION) \
		--build-arg COMMIT=$(COMMIT) \
		--build-arg BUILD_TIME=$(BUILD_TIME) \
		--target production \
		-t $(DOCKER_IMAGE):$(DOCKER_TAG) .
	@echo "Image: $(DOCKER_IMAGE):$(DOCKER_TAG)"

docker-build-distroless: ## Build Docker image (distroless, smaller)
	docker build \
		--build-arg VERSION=$(VERSION) \
		--build-arg COMMIT=$(COMMIT) \
		--build-arg BUILD_TIME=$(BUILD_TIME) \
		--target distroless \
		-t $(DOCKER_IMAGE):$(DOCKER_TAG)-distroless .
	@echo "Image: $(DOCKER_IMAGE):$(DOCKER_TAG)-distroless"

docker-run: ## Run Docker container
	docker run --rm -p 8080:8080 \
		-e APP_ENV=production \
		-e PORT=8080 \
		$(DOCKER_IMAGE):$(DOCKER_TAG)

docker-run-env: ## Run Docker container with .env file
	docker run --rm -p 8080:8080 \
		--env-file .env \
		$(DOCKER_IMAGE):$(DOCKER_TAG)

docker-test: ## Run tests in Docker container
	docker build --target test -t $(DOCKER_IMAGE):test .

docker-size: ## Show Docker image sizes
	@docker images $(DOCKER_IMAGE) --format "table {{.Repository}}\t{{.Tag}}\t{{.Size}}"

# =============================================================================
# All-in-one
# =============================================================================

all: tidy lint test build ## Run tidy, lint, test, and build
