# NexusFlow Development Setup

This guide will help you set up your local development environment for NexusFlow.

## Prerequisites

Before you begin, ensure you have the following installed:

- **Go 1.24+** - [Download](https://go.dev/dl/)
- **Docker & Docker Compose** - [Download](https://www.docker.com/products/docker-desktop)
- **protoc (Protocol Buffers Compiler)** - [Installation Guide](https://grpc.io/docs/protoc-installation/)
- **kubectl** (optional, for Kubernetes deployment) - [Download](https://kubernetes.io/docs/tasks/tools/)
- **Helm 3+** (optional, for Kubernetes deployment) - [Download](https://helm.sh/docs/intro/install/)
- **Git** - [Download](https://git-scm.com/downloads)

## Initial Setup

### 1. Clone the Repository

```bash
git clone https://github.com/yourusername/nexusflow.git
cd nexusflow
```

### 2. Install Go Tools

Install the required Go tools for protobuf code generation and linting:

```bash
make install-tools
```

This will install:

- `protoc-gen-go` - Protocol Buffers Go plugin
- `protoc-gen-go-grpc` - gRPC Go plugin
- `golangci-lint` - Go linter

### 3. Generate Protobuf Code

Generate Go code from protobuf definitions:

```bash
make generate-proto
```

This will create Go files in `pkg/proto/` from all `.proto` files.

### 4. Sync Go Workspace

Sync the Go workspace and download dependencies:

```bash
make sync
```

## Running Locally

### Start Infrastructure Services

Start all required infrastructure services (PostgreSQL, Kafka, Elasticsearch, MinIO, etc.):

```bash
make infra-up
```

Wait for all services to be healthy. You can check the status with:

```bash
docker-compose ps
```

### Access Infrastructure Services

Once running, you can access:

- **Kafka Console**: <http://localhost:8080>
- **Traefik Dashboard**: <http://localhost:8081>
- **MinIO Console**: <http://localhost:9001> (credentials: nexusflow/nexusflow123)
- **Elasticsearch**: <http://localhost:9200>
- **Prometheus**: <http://localhost:9091>
- **Grafana**: <http://localhost:3001> (credentials: admin/admin)
- **Jaeger UI**: <http://localhost:16686>
- **Ory Hydra Admin**: <http://localhost:4445>
- **Ory Hydra Public**: <http://localhost:4444>

### Build All Services

Build all microservices:

```bash
make build-all
```

Binaries will be created in the `bin/` directory.

### Run All Services

Run all services locally (requires infrastructure to be running):

```bash
make run-all
```

This will start all services in the background. Press `Ctrl+C` to stop.

### Run Individual Services

To run a specific service:

```bash
cd services/user-service
go run ./cmd/server
```

## Development Workflow

### Creating a New Service

1. Copy the template service:

   ```bash
   cp -r services/template-service services/your-service
   ```

2. Update `services/your-service/go.mod` with the correct module name

3. Add the service to `go.work`:

   ```bash
   echo "  ./services/your-service" >> go.work
   ```

4. Implement your service logic in:
   - `internal/handler/` - gRPC handlers
   - `internal/service/` - Business logic
   - `internal/repository/` - Data access

5. Add the service to `docker-compose.yml` for local development

### Working with Protobuf

1. Edit `.proto` files in `proto/`
2. Run `make generate-proto` to regenerate Go code
3. Implement the generated service interfaces in your handlers

### Running Tests

Run all tests:

```bash
make test-all
```

Run tests for a specific service:

```bash
cd services/user-service
go test ./...
```

### Linting

Run linters on all code:

```bash
make lint
```

### Database Migrations

Database migrations will be handled by each service on startup. For manual migrations:

```bash
# Connect to PostgreSQL
docker exec -it nexusflow-postgres psql -U nexusflow -d nexusflow

# Run SQL commands
\dt  # List tables
```

## Troubleshooting

### Port Conflicts

If you encounter port conflicts, check what's running on the ports:

```bash
lsof -i :5432  # PostgreSQL
lsof -i :9092  # Kafka
lsof -i :9200  # Elasticsearch
```

Stop conflicting services or modify ports in `docker-compose.yml`.

### Docker Issues

If Docker services fail to start:

```bash
# Stop all services
make infra-down

# Remove volumes (WARNING: This will delete all data)
docker-compose down -v

# Start fresh
make infra-up
```

### Go Module Issues

If you encounter Go module errors:

```bash
# Clean Go cache
go clean -modcache

# Sync workspace
make sync
```

### Protobuf Generation Errors

If protobuf generation fails:

```bash
# Check protoc installation
protoc --version

# Reinstall Go plugins
make install-tools

# Try generating again
make generate-proto
```

## Next Steps

- Read the [Architecture Overview](../architecture/overview.md)
- Review [API Conventions](conventions.md)
- Check out the [Microservices Documentation](../architecture/microservices.md)
- Join our [Discord Community](https://discord.gg/nexusflow)

## Getting Help

- 📖 [Documentation](../../README.md)
- 💬 [GitHub Discussions](https://github.com/yourusername/nexusflow/discussions)
- 🐛 [Issue Tracker](https://github.com/yourusername/nexusflow/issues)
