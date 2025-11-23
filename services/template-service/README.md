# Service Template

This directory serves as a template for creating new microservices in NexusFlow.

## Structure

```
template-service/
├── cmd/
│   └── server/
│       └── main.go          # Service entry point
├── internal/
│   ├── handler/             # gRPC handlers
│   │   └── handler.go
│   ├── service/             # Business logic
│   │   └── service.go
│   └── repository/          # Data access layer
│       └── repository.go
├── Dockerfile               # Multi-stage Docker build
├── go.mod                   # Go module definition
├── go.sum                   # Go module checksums
└── .env.example             # Example environment variables
```

## Creating a New Service

1. Copy this template directory:

   ```bash
   cp -r services/template-service services/your-service
   ```

2. Update `go.mod` with the correct module name

3. Implement your gRPC handlers in `internal/handler/`

4. Implement business logic in `internal/service/`

5. Implement data access in `internal/repository/`

6. Update the Dockerfile if needed

7. Add the service to `go.work` in the root directory

8. Add the service to `docker-compose.yml` for local development

## Environment Variables

See `.env.example` for required environment variables.

## Building

```bash
go build -o bin/your-service ./cmd/server
```

## Running

```bash
./bin/your-service
```
