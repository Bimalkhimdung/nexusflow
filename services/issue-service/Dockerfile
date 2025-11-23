# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git

# Copy shared packages
COPY pkg ./pkg
COPY proto ./proto

# Copy service code
COPY services/issue-service ./services/issue-service

# Build service
WORKDIR /app/services/issue-service
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/bin/issue-service ./cmd/server

# Final stage
FROM alpine:3.19

WORKDIR /app

# Install runtime dependencies
RUN apk add --no-cache ca-certificates tzdata

# Copy binary and config
COPY --from=builder /app/bin/issue-service .
COPY --from=builder /app/services/issue-service/config.yaml .
COPY --from=builder /app/services/issue-service/migrations ./migrations

# Expose ports
EXPOSE 8083 50054

# Run service
CMD ["./issue-service"]
