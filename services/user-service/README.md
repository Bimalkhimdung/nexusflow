# User Service

User management microservice for NexusFlow.

## Overview

The user-service handles all user-related operations including:

- User CRUD operations
- User search and listing
- User preferences management
- User profile management
- Multi-tenant user isolation

## Features

- ✅ Full CRUD operations via gRPC
- ✅ Multi-tenant support (organization-based isolation)
- ✅ Soft deletes
- ✅ User preferences (JSONB storage)
- ✅ Event publishing to Kafka
- ✅ Automatic database migrations
- ✅ Optimistic locking
- ✅ Pagination and search

## Database Schema

### Users Table

```sql
CREATE TABLE users (
    id UUID PRIMARY KEY,
    organization_id UUID NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    display_name VARCHAR(255) NOT NULL,
    avatar_url VARCHAR(500),
    timezone VARCHAR(50) DEFAULT 'UTC',
    locale VARCHAR(10) DEFAULT 'en-US',
    status VARCHAR(20) NOT NULL DEFAULT 'active',
    preferences JSONB DEFAULT '{}',
    created_by UUID,
    updated_by UUID,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL,
    last_login_at TIMESTAMP WITH TIME ZONE,
    version BIGINT NOT NULL DEFAULT 1,
    deleted_at TIMESTAMP WITH TIME ZONE
);
```

## API Endpoints

All endpoints are available via gRPC. See `proto/user/v1/user.proto` for full definitions.

### User Management

- `GetUser(id)` - Get user by ID
- `GetUserByEmail(email)` - Get user by email
- `CreateUser(input)` - Create new user
- `UpdateUser(id, input)` - Update user
- `DeleteUser(id)` - Soft delete user
- `ListUsers(orgID, pagination)` - List users
- `SearchUsers(query, orgIDs, pagination)` - Search users

### User Profile

- `GetUserProfile(id)` - Get public user profile
- `UpdateUserPreferences(id, preferences)` - Update user preferences

## Events Published

The service publishes the following events to Kafka topic `nexusflow.users`:

### user.created

```json
{
  "type": "user.created",
  "organization_id": "uuid",
  "user_id": "uuid",
  "timestamp": "2025-11-23T20:00:00Z",
  "payload": {
    "user_id": "uuid",
    "email": "user@example.com",
    "display_name": "User Name"
  }
}
```

### user.updated

```json
{
  "type": "user.updated",
  "organization_id": "uuid",
  "user_id": "uuid",
  "timestamp": "2025-11-23T20:00:00Z",
  "payload": {
    "user_id": "uuid",
    "email": "user@example.com",
    "display_name": "Updated Name"
  }
}
```

### user.deleted

```json
{
  "type": "user.deleted",
  "organization_id": "uuid",
  "user_id": "uuid",
  "timestamp": "2025-11-23T20:00:00Z",
  "payload": {
    "user_id": "uuid",
    "email": "user@example.com"
  }
}
```

## Running Locally

### Prerequisites

- PostgreSQL running (via `make infra-up`)
- Kafka running (via `make infra-up`)

### Build

```bash
cd services/user-service
go build -o ../../bin/user-service ./cmd/server
```

### Run

```bash
./bin/user-service
```

The service will:

1. Connect to PostgreSQL
2. Run database migrations automatically
3. Connect to Kafka
4. Start gRPC server on port 9090

### Test with grpcurl

```bash
# List services
grpcurl -plaintext localhost:9090 list

# Create user
grpcurl -plaintext -d '{
  "email": "test@example.com",
  "display_name": "Test User",
  "timezone": "UTC",
  "locale": "en-US"
}' localhost:9090 nexusflow.user.v1.UserService.CreateUser

# Get user
grpcurl -plaintext -d '{"id": "user-uuid"}' \
  localhost:9090 nexusflow.user.v1.UserService.GetUser

# List users
grpcurl -plaintext -d '{
  "pagination": {"page": 1, "page_size": 20}
}' localhost:9090 nexusflow.user.v1.UserService.ListUsers
```

## Testing

### Run Tests

```bash
go test ./... -v -race -coverprofile=coverage.txt
```

### Run with Coverage

```bash
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out
```

## Configuration

Configuration is loaded from:

1. Environment variables (prefixed with `USER_SERVICE_`)
2. `config.yaml` file
3. Default values

See `.env.example` for all available configuration options.

## Architecture

```
cmd/server/main.go          # Entry point
├── internal/
│   ├── models/             # Domain models
│   │   └── user.go
│   ├── repository/         # Data access layer
│   │   └── user_repository.go
│   ├── service/            # Business logic
│   │   └── user_service.go
│   └── handler/            # gRPC handlers
│       └── user_handler.go
└── migrations/             # Database migrations
    └── 001_create_users_table.sql
```

## Dependencies

- `pkg/logger` - Structured logging
- `pkg/config` - Configuration management
- `pkg/database` - Database connection and ORM
- `pkg/kafka` - Event publishing
- `pkg/proto` - Protobuf definitions
- `golang-migrate` - Database migrations

## Monitoring

The service exposes:

- gRPC health check endpoint
- Metrics (TODO: Prometheus)
- Distributed tracing (TODO: Jaeger)

## Future Enhancements

- [ ] Add authentication middleware
- [ ] Add authorization checks
- [ ] Add rate limiting
- [ ] Add caching layer
- [ ] Add metrics collection
- [ ] Add distributed tracing
- [ ] Add user avatar upload
- [ ] Add email verification
- [ ] Add password management (if not using Ory)
