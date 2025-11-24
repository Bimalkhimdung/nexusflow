module nexusflow/services/gateway-service

go 1.22

require (
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.19.1
	google.golang.org/grpc v1.62.1
	github.com/nexusflow/nexusflow/pkg/proto v0.0.0
)

replace github.com/nexusflow/nexusflow/pkg/proto => ../../pkg/proto
