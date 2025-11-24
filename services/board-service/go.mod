module github.com/nexusflow/nexusflow/services/board-service

go 1.24

require (
    github.com/uptrace/bun v1.1.17
    github.com/gorilla/websocket v1.5.0
    github.com/nexusflow/nexusflow/pkg/kafka v0.0.0
    github.com/nexusflow/nexusflow/pkg/logger v0.0.0
    github.com/nexusflow/nexusflow/pkg/proto v0.0.0
    github.com/nexusflow/nexusflow/pkg/database v0.0.0
    github.com/nexusflow/nexusflow/pkg/config v0.0.0
    google.golang.org/grpc v1.77.0
)

replace (
    github.com/nexusflow/nexusflow/pkg/config => ../../pkg/config
    github.com/nexusflow/nexusflow/pkg/database => ../../pkg/database
    github.com/nexusflow/nexusflow/pkg/kafka => ../../pkg/kafka
    github.com/nexusflow/nexusflow/pkg/logger => ../../pkg/logger
    github.com/nexusflow/nexusflow/pkg/proto => ../../pkg/proto
)
