.PHONY: help generate-proto build-all test-all lint clean run-all docker-build docker-push

# Default target
help:
	@echo "NexusFlow - Available Make Targets"
	@echo "===================================="
	@echo "  generate-proto    Generate Go code from protobuf definitions"
	@echo "  build-all         Build all microservices"
	@echo "  test-all          Run tests for all services"
	@echo "  lint              Run golangci-lint on all services"
	@echo "  clean             Clean build artifacts"
	@echo "  run-all           Run all services locally"
	@echo "  docker-build      Build Docker images for all services"
	@echo "  docker-push       Push Docker images to registry"
	@echo "  infra-up          Start infrastructure services (docker-compose)"
	@echo "  infra-down        Stop infrastructure services"

# Generate protobuf code
generate-proto:
	@echo "Generating protobuf code..."
	@chmod +x scripts/generate-proto.sh
	@./scripts/generate-proto.sh

# Build all services
build-all:
	@echo "Building all services..."
	@for service in services/*-service; do \
		if [ -d "$$service" ]; then \
			echo "Building $$service..."; \
			cd $$service && go build -o ../../bin/$$(basename $$service) ./cmd/server && cd ../..; \
		fi \
	done

# Run tests for all services
test-all:
	@echo "Running tests..."
	@for dir in pkg/* services/*; do \
		if [ -d "$$dir" ] && [ -f "$$dir/go.mod" ]; then \
			echo "Testing $$dir..."; \
			cd $$dir && go test -v -race -coverprofile=coverage.txt -covermode=atomic ./... && cd ../..; \
		fi \
	done

# Run linter
lint:
	@echo "Running golangci-lint..."
	@if command -v golangci-lint > /dev/null; then \
		for dir in pkg/* services/*; do \
			if [ -d "$$dir" ] && [ -f "$$dir/go.mod" ] && [ "$$dir" != "pkg/proto" ]; then \
				echo "Linting $$dir..."; \
				cd $$dir && golangci-lint run ./... && cd ../..; \
			fi \
		done \
	else \
		echo "golangci-lint not installed. Install with: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"; \
	fi

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	@rm -rf bin/
	@rm -rf dist/
	@rm -f coverage.txt
	@find . -name "*.pb.go" -type f -delete
	@find . -name "*_grpc.pb.go" -type f -delete

# Run all services locally
run-all:
	@echo "Starting all services..."
	@echo "Make sure infrastructure is running: make infra-up"
	@for service in services/*-service; do \
		if [ -d "$$service" ]; then \
			echo "Starting $$service..."; \
			cd $$service && go run ./cmd/server & cd ../..; \
		fi \
	done
	@echo "All services started. Press Ctrl+C to stop."

# Build Docker images
docker-build:
	@echo "Building Docker images..."
	@for service in services/*-service; do \
		if [ -d "$$service" ]; then \
			echo "Building Docker image for $$service..."; \
			docker build -t nexusflow/$$(basename $$service):latest -f $$service/Dockerfile $$service; \
		fi \
	done

# Push Docker images
docker-push:
	@echo "Pushing Docker images..."
	@for service in services/*-service; do \
		if [ -d "$$service" ]; then \
			echo "Pushing Docker image for $$service..."; \
			docker push nexusflow/$$(basename $$service):latest; \
		fi \
	done

# Start infrastructure services
infra-up:
	@echo "Starting infrastructure services..."
	@docker-compose up -d
	@echo "Waiting for services to be healthy..."
	@sleep 5
	@docker-compose ps

# Stop infrastructure services
infra-down:
	@echo "Stopping infrastructure services..."
	@docker-compose down

# Install development tools
install-tools:
	@echo "Installing development tools..."
	@go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	@go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@echo "Tools installed successfully"

# Sync Go workspace
sync:
	@echo "Syncing Go workspace..."
	@go work sync
	@go mod tidy -C pkg/logger
	@go mod tidy -C pkg/config
	@go mod tidy -C pkg/database
	@go mod tidy -C pkg/kafka
