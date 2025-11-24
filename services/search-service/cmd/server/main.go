package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"

	"github.com/nexusflow/nexusflow/pkg/config"
	"github.com/nexusflow/nexusflow/pkg/logger"
	pb "github.com/nexusflow/nexusflow/pkg/proto/search/v1"
	"github.com/nexusflow/nexusflow/services/search-service/internal/elasticsearch"
	"github.com/nexusflow/nexusflow/services/search-service/internal/handler"
	"github.com/nexusflow/nexusflow/services/search-service/internal/service"
)

const serviceName = "search-service"

func main() {
	// Initialize logger
	log, err := logger.NewDefault(serviceName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}
	defer func() { _ = log.Sync() }()

	log.Sugar().Infow("Starting search-service")

	// Load configuration
	cfg, err := config.New(serviceName)
	if err != nil {
		log.Sugar().Fatal("Failed to load configuration")
	}

	// Initialize Elasticsearch client
	esAddresses := []string{"http://localhost:9200"}
	esClient, err := elasticsearch.NewClient(esAddresses, log)
	if err != nil {
		log.Sugar().Fatalw("Failed to connect to Elasticsearch", "error", err)
	}

	// Initialize indices
	if err := initializeIndices(esClient, log); err != nil {
		log.Sugar().Warnw("Failed to initialize indices", "error", err)
	}

	// Initialize layers
	svc := service.NewSearchService(esClient, log)
	h := handler.NewSearchHandler(svc, log)

	// Create gRPC server
	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			// Add interceptors here
		),
	)

	// Register services
	pb.RegisterSearchServiceServer(grpcServer, h)

	// Register health check
	healthServer := health.NewServer()
	grpc_health_v1.RegisterHealthServer(grpcServer, healthServer)
	healthServer.SetServingStatus("", grpc_health_v1.HealthCheckResponse_SERVING)

	// Register reflection service (for development)
	reflection.Register(grpcServer)

	// Start gRPC server
	serverCfg := cfg.GetServer()
	grpcPort := serverCfg.GRPCPort
	if grpcPort == 0 || grpcPort == 9090 {
		grpcPort = 50060
	}

	grpcAddr := fmt.Sprintf("%s:%d", serverCfg.Host, grpcPort)
	listener, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		log.Sugar().Fatalw("Failed to listen", "error", err, "addr", grpcAddr)
	}

	// Start gRPC server in goroutine
	go func() {
		log.Sugar().Infow("gRPC server listening", "addr", grpcAddr)
		if err := grpcServer.Serve(listener); err != nil {
			log.Sugar().Fatal("Failed to serve gRPC")
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Sugar().Infow("Shutting down server...")

	// Graceful shutdown
	grpcServer.GracefulStop()

	log.Sugar().Infow("Server stopped")
}

func initializeIndices(es *elasticsearch.Client, log *logger.Logger) error {
	ctx := context.Background()

	// Issues index mapping
	issuesMapping := map[string]interface{}{
		"mappings": map[string]interface{}{
			"properties": map[string]interface{}{
				"key":         map[string]string{"type": "keyword"},
				"title":       map[string]string{"type": "text"},
				"description": map[string]string{"type": "text"},
				"status":      map[string]string{"type": "keyword"},
				"priority":    map[string]string{"type": "keyword"},
				"type":        map[string]string{"type": "keyword"},
				"project_id":  map[string]string{"type": "keyword"},
				"assignee_id": map[string]string{"type": "keyword"},
				"reporter_id": map[string]string{"type": "keyword"},
				"labels":      map[string]string{"type": "keyword"},
				"created_at":  map[string]string{"type": "date"},
				"updated_at":  map[string]string{"type": "date"},
			},
		},
	}

	if err := es.CreateIndex(ctx, "issues", issuesMapping); err != nil {
		return fmt.Errorf("failed to create issues index: %w", err)
	}

	// Projects index mapping
	projectsMapping := map[string]interface{}{
		"mappings": map[string]interface{}{
			"properties": map[string]interface{}{
				"key":         map[string]string{"type": "keyword"},
				"name":        map[string]string{"type": "text"},
				"description": map[string]string{"type": "text"},
				"org_id":      map[string]string{"type": "keyword"},
				"lead_id":     map[string]string{"type": "keyword"},
				"created_at":  map[string]string{"type": "date"},
				"updated_at":  map[string]string{"type": "date"},
			},
		},
	}

	if err := es.CreateIndex(ctx, "projects", projectsMapping); err != nil {
		return fmt.Errorf("failed to create projects index: %w", err)
	}

	// Users index mapping
	usersMapping := map[string]interface{}{
		"mappings": map[string]interface{}{
			"properties": map[string]interface{}{
				"email":      map[string]string{"type": "keyword"},
				"name":       map[string]string{"type": "text"},
				"created_at": map[string]string{"type": "date"},
			},
		},
	}

	if err := es.CreateIndex(ctx, "users", usersMapping); err != nil {
		return fmt.Errorf("failed to create users index: %w", err)
	}

	log.Sugar().Infow("Initialized Elasticsearch indices")
	return nil
}
