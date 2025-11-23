# NexusFlow Development Conventions

This document outlines coding conventions, best practices, and standards for NexusFlow development.

## API Design Conventions

### REST API

- Use RESTful resource naming (plural nouns)
- HTTP methods: GET (read), POST (create), PUT (update), DELETE (delete), PATCH (partial update)
- Use proper HTTP status codes
- Version APIs in URL: `/api/v1/`
- Use query parameters for filtering, sorting, pagination

**Example:**

```
GET    /api/v1/projects              # List projects
GET    /api/v1/projects/{id}         # Get project
POST   /api/v1/projects              # Create project
PUT    /api/v1/projects/{id}         # Update project
DELETE /api/v1/projects/{id}         # Delete project
GET    /api/v1/projects/{id}/issues  # List project issues
```

### gRPC API

- Use protobuf for service definitions
- Follow Google API design guide
- Use standard request/response patterns
- Include pagination in list operations
- Use proper error codes

### GraphQL API

- Use schema-first approach
- Implement DataLoader for N+1 prevention
- Use relay-style pagination
- Provide comprehensive error messages

## Error Handling

### Error Codes

Use consistent error codes across services:

```go
const (
    ErrCodeInvalidInput    = "INVALID_INPUT"
    ErrCodeNotFound        = "NOT_FOUND"
    ErrCodeUnauthorized    = "UNAUTHORIZED"
    ErrCodeForbidden       = "FORBIDDEN"
    ErrCodeConflict        = "CONFLICT"
    ErrCodeInternal        = "INTERNAL_ERROR"
)
```

### Error Responses

```json
{
  "error": {
    "code": "INVALID_INPUT",
    "message": "Invalid email format",
    "details": {
      "field": "email",
      "value": "invalid-email"
    }
  }
}
```

## Database Conventions

### Table Naming

- Use `snake_case` for table and column names
- Use plural nouns for table names
- Include `organization_id` for multi-tenant tables

### Indexes

- Add indexes for foreign keys
- Add indexes for frequently queried columns
- Use composite indexes for multi-column queries
- Name indexes: `idx_{table}_{column(s)}`

### Migrations

- Use sequential numbering: `001_initial_schema.sql`
- Include both up and down migrations
- Never modify existing migrations
- Test migrations on production-like data

## Logging Conventions

### Log Levels

- **DEBUG**: Detailed information for debugging
- **INFO**: General informational messages
- **WARN**: Warning messages for potential issues
- **ERROR**: Error messages for failures
- **FATAL**: Critical errors that require immediate attention

### Structured Logging

Always use structured logging with fields:

```go
log.Info("User created",
    "user_id", userID,
    "email", email,
    "organization_id", orgID,
)
```

### Log Context

Include trace context in all logs:

```go
log = log.WithContext(ctx)
log.Info("Processing request")
```

## Testing Conventions

### Unit Tests

- Test file naming: `{file}_test.go`
- Test function naming: `Test{FunctionName}`
- Use table-driven tests for multiple cases
- Mock external dependencies
- Aim for >80% coverage

**Example:**

```go
func TestCreateUser(t *testing.T) {
    tests := []struct {
        name    string
        input   CreateUserInput
        want    *User
        wantErr bool
    }{
        {
            name: "valid user",
            input: CreateUserInput{Email: "test@example.com"},
            want: &User{Email: "test@example.com"},
            wantErr: false,
        },
        // More test cases...
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := CreateUser(tt.input)
            if (err != nil) != tt.wantErr {
                t.Errorf("CreateUser() error = %v, wantErr %v", err, tt.wantErr)
                return
            }
            if !reflect.DeepEqual(got, tt.want) {
                t.Errorf("CreateUser() = %v, want %v", got, tt.want)
            }
        })
    }
}
```

### Integration Tests

- Use separate test database
- Clean up test data after each test
- Use realistic test data
- Test happy path and error cases

## Event Conventions

### Event Naming

Use dot notation: `{resource}.{action}`

Examples:

- `issue.created`
- `issue.updated`
- `issue.deleted`
- `issue.assigned`
- `project.created`

### Event Payload

Include all relevant information:

```json
{
  "id": "event-uuid",
  "type": "issue.created",
  "organization_id": "org-uuid",
  "project_id": "project-uuid",
  "user_id": "user-uuid",
  "timestamp": "2025-11-23T20:00:00Z",
  "payload": {
    "issue_id": "issue-uuid",
    "issue_key": "PROJ-123",
    "summary": "New issue"
  }
}
```

## Security Conventions

### Authentication

- Always validate JWT tokens
- Check token expiration
- Verify token signature
- Extract user context from token

### Authorization

- Check permissions before operations
- Use RBAC for access control
- Validate organization membership
- Log authorization failures

### Input Validation

- Validate all user input
- Sanitize HTML/SQL input
- Use parameterized queries
- Limit input size

### Secrets Management

- Never commit secrets to git
- Use environment variables
- Rotate secrets regularly
- Use Kubernetes secrets in production

## Performance Conventions

### Database Queries

- Use connection pooling
- Avoid N+1 queries
- Use indexes appropriately
- Limit result sets
- Use pagination

### Caching

- Cache frequently accessed data
- Set appropriate TTLs
- Invalidate cache on updates
- Use Redis for distributed cache

### API Performance

- Implement rate limiting
- Use compression (gzip)
- Minimize payload size
- Use HTTP/2 where possible

## Documentation Conventions

### Code Comments

```go
// CreateUser creates a new user in the system.
// It validates the input, checks for duplicates, and publishes
// a user.created event to Kafka.
//
// Parameters:
//   - ctx: Request context with trace information
//   - input: User creation input with email and display name
//
// Returns:
//   - *User: Created user object
//   - error: Error if creation fails
func CreateUser(ctx context.Context, input CreateUserInput) (*User, error) {
    // Implementation...
}
```

### API Documentation

- Document all endpoints
- Include request/response examples
- List possible error codes
- Provide authentication requirements

## Monitoring Conventions

### Metrics

Expose the following metrics:

- Request count
- Request duration
- Error rate
- Database query duration
- Kafka publish/consume metrics

### Health Checks

Implement health check endpoints:

- `/health/live` - Liveness probe
- `/health/ready` - Readiness probe

### Alerts

Set up alerts for:

- High error rate (>5%)
- High latency (p95 > 500ms)
- Database connection failures
- Kafka consumer lag

## Deployment Conventions

### Container Images

- Use multi-stage builds
- Run as non-root user
- Include health checks
- Tag with git commit SHA
- Keep images small (<100MB)

### Kubernetes

- Set resource limits
- Configure liveness/readiness probes
- Use ConfigMaps for configuration
- Use Secrets for sensitive data
- Implement graceful shutdown

### Environment Variables

Use consistent naming:

- `{SERVICE}_HOST` - Service host
- `{SERVICE}_PORT` - Service port
- `{SERVICE}_DATABASE_URL` - Database connection string
- `LOG_LEVEL` - Logging level
- `ENVIRONMENT` - Environment (dev/staging/prod)
