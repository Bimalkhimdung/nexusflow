# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY go.work go.work.sum ./
COPY services/board-service/go.mod services/board-service/go.sum ./services/board-service/
COPY pkg/ ./pkg/

# Download dependencies
RUN go mod download

# Copy source code
COPY services/board-service/ ./services/board-service/

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/board-service ./services/board-service/cmd/server

# Runtime stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /app

COPY --from=builder /bin/board-service .
COPY services/board-service/migrations ./migrations
COPY services/board-service/config.yaml .

EXPOSE 50056 8082

CMD ["./board-service"]
