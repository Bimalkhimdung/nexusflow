package main

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"

	"github.com/nexusflow/nexusflow/pkg/config"
	"github.com/nexusflow/nexusflow/pkg/database"
	"github.com/nexusflow/nexusflow/pkg/kafka"
	"github.com/nexusflow/nexusflow/pkg/logger"
)

const serviceName = "template-service"

func main() {
	// Initialize logger
	log, err := logger.NewDefault(serviceName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}
	defer func() { _ = log.Sync() }()

	log.Info("Starting service")

	// Load configuration
	cfg, err := config.New(serviceName)
	if err != nil {
		log.Fatal("Failed to load configuration")
	}

	// Initialize database
	dbCfg := cfg.GetDatabase()
	db, err := database.New(database.Config{
		Host:            dbCfg.Host,
		Port:            dbCfg.Port,
		User:            dbCfg.User,
		Password:        dbCfg.Password,
		Database:        dbCfg.Database,
		SSLMode:         dbCfg.SSLMode,
		MaxOpenConns:    dbCfg.MaxOpenConns,
		MaxIdleConns:    dbCfg.MaxIdleConns,
		ConnMaxLifetime: time.Duration(dbCfg.ConnMaxLifetime) * time.Second,
	})
	if err != nil {
		log.Fatal("Failed to connect to database")
	}
	defer db.Close()

	log.Info("Database connection established")

	// Initialize Kafka producer
	kafkaCfg := cfg.GetKafka()
	producer, err := kafka.NewProducer(kafka.ProducerConfig{
		Brokers: kafkaCfg.Brokers,
	})
	if err != nil {
		log.Fatal("Failed to create Kafka producer")
	}
	defer producer.Close()

	log.Info("Kafka producer initialized")

	// Create gRPC server
	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			// Add interceptors here (auth, logging, tracing, etc.)
		),
	)

	// Register health check
	healthServer := health.NewServer()
	grpc_health_v1.RegisterHealthServer(grpcServer, healthServer)
	healthServer.SetServingStatus("", grpc_health_v1.HealthCheckResponse_SERVING)

	// Register reflection service (for development)
	reflection.Register(grpcServer)

	// TODO: Register your service implementation here
	// Example:
	// handler := handler.NewHandler(service, log)
	// pb.RegisterYourServiceServer(grpcServer, handler)

	// Start gRPC server
	serverCfg := cfg.GetServer()
	addr := fmt.Sprintf("%s:%d", serverCfg.Host, serverCfg.GRPCPort)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatal("Failed to listen")
	}

	// Start server in goroutine
	go func() {
		log.Info("gRPC server listening")
		if err := grpcServer.Serve(listener); err != nil {
			log.Fatal("Failed to serve")
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("Shutting down server...")

	// Graceful shutdown
	grpcServer.GracefulStop()

	log.Info("Server stopped")
}
