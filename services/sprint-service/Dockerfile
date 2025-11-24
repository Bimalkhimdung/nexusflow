# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY go.work go.work.sum ./
COPY services/sprint-service/go.mod services/sprint-service/go.sum ./services/sprint-service/
COPY pkg/ ./pkg/

# Download dependencies
RUN go mod download

# Copy source code
COPY services/sprint-service/ ./services/sprint-service/

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/sprint-service ./services/sprint-service/cmd/server

# Runtime stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /app

COPY --from=builder /bin/sprint-service .
COPY services/sprint-service/migrations ./migrations
COPY services/sprint-service/config.yaml .

EXPOSE 50057

CMD ["./sprint-service"]
