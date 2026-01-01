.PHONY: build run test clean tidy lint swagger swagger-install

# Build variables
BINARY_NAME=api
BUILD_DIR=bin
MAIN_PATH=./cmd/api

# Go commands
GOCMD=go
GOBUILD=$(GOCMD) build
GORUN=$(GOCMD) run
GOTEST=$(GOCMD) test
GOMOD=$(GOCMD) mod
GOVET=$(GOCMD) vet

# Build the application
build:
	@echo "Building..."
	@mkdir -p $(BUILD_DIR)
	$(GOBUILD) -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PATH)

# Run the application
run:
	$(GORUN) $(MAIN_PATH)

# Run with hot reload (requires air: go install github.com/cosmtrek/air@latest)
dev:
	@air

# Run tests
test:
	$(GOTEST) -v -race ./...

# Run tests with coverage
test-coverage:
	$(GOTEST) -v -race -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html

# Clean build artifacts
clean:
	@echo "Cleaning..."
	@rm -rf $(BUILD_DIR)
	@rm -f coverage.out coverage.html

# Tidy dependencies
tidy:
	$(GOMOD) tidy

# Download dependencies
deps:
	$(GOMOD) download

# Run linter (requires golangci-lint)
lint:
	@golangci-lint run ./...

# Vet the code
vet:
	$(GOVET) ./...

# Build for production (Linux)
build-linux:
	@echo "Building for Linux..."
	@mkdir -p $(BUILD_DIR)
	GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 $(MAIN_PATH)

# Docker variables
DOCKER_IMAGE=boilerplate-go
DOCKER_TAG=latest
VERSION=$(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT=$(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_TIME=$(shell date -u +"%Y-%m-%dT%H:%M:%SZ")

# Docker build (production alpine)
docker-build:
	docker build \
		--build-arg VERSION=$(VERSION) \
		--build-arg COMMIT=$(COMMIT) \
		--build-arg BUILD_TIME=$(BUILD_TIME) \
		--target production \
		-t $(DOCKER_IMAGE):$(DOCKER_TAG) .

# Docker build (distroless - smaller image)
docker-build-distroless:
	docker build \
		--build-arg VERSION=$(VERSION) \
		--build-arg COMMIT=$(COMMIT) \
		--build-arg BUILD_TIME=$(BUILD_TIME) \
		--target distroless \
		-t $(DOCKER_IMAGE):$(DOCKER_TAG)-distroless .

# Docker run
docker-run:
	docker run --rm -p 8080:8080 \
		-e APP_ENV=production \
		-e PORT=8080 \
		$(DOCKER_IMAGE):$(DOCKER_TAG)

# Docker run with env file
docker-run-env:
	docker run --rm -p 8080:8080 \
		--env-file .env \
		$(DOCKER_IMAGE):$(DOCKER_TAG)

# Docker test (runs tests in container)
docker-test:
	docker build --target test -t $(DOCKER_IMAGE):test .

# Show Docker image size
docker-size:
	@docker images $(DOCKER_IMAGE) --format "table {{.Repository}}\t{{.Tag}}\t{{.Size}}"

# Install swag CLI
swagger-install:
	go install github.com/swaggo/swag/cmd/swag@latest

# Generate Swagger documentation
swagger:
	@echo "Generating Swagger docs..."
	go run github.com/swaggo/swag/cmd/swag@latest init -g cmd/api/main.go -o docs --parseDependency --parseInternal
	@echo "Swagger docs generated in docs/"

# Format Swagger annotations
swagger-fmt:
	go run github.com/swaggo/swag/cmd/swag@latest fmt

