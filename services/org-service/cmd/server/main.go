package main

import (
	"database/sql"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"

	"github.com/nexusflow/nexusflow/pkg/config"
	"github.com/nexusflow/nexusflow/pkg/database"
	"github.com/nexusflow/nexusflow/pkg/kafka"
	"github.com/nexusflow/nexusflow/pkg/logger"
	"github.com/nexusflow/nexusflow/pkg/middleware"
	pb "github.com/nexusflow/nexusflow/pkg/proto/org/v1"
	"github.com/nexusflow/nexusflow/services/org-service/internal/handler"
	"github.com/nexusflow/nexusflow/services/org-service/internal/repository"
	"github.com/nexusflow/nexusflow/services/org-service/internal/service"
)

const serviceName = "org-service"

func main() {
	// Initialize logger
	log, err := logger.NewDefault(serviceName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}
	defer func() { _ = log.Sync() }()

	log.Sugar().Infow("Starting org-service")

	// Load configuration
	cfg, err := config.New(serviceName)
	if err != nil {
		log.Sugar().Fatal("Failed to load configuration")
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
		log.Sugar().Fatal("Failed to connect to database")
	}
	defer db.Close()

	log.Sugar().Infow("Database connection established")

	// Run database migrations
	if err := runMigrations(db.GetSQLDB(), log); err != nil {
		log.Sugar().Fatal("Failed to run migrations")
	}

	// Initialize Kafka producer
	kafkaCfg := cfg.GetKafka()
	var producer *kafka.Producer
	producer, err = kafka.NewProducer(kafka.ProducerConfig{
		Brokers: kafkaCfg.Brokers,
	})
	if err != nil {
		log.Sugar().Warnw("Failed to create Kafka producer, continuing without events", "error", err)
		producer = nil
	} else {
		defer producer.Close()
		log.Sugar().Infow("Kafka producer initialized")
	}

	// Initialize layers
	orgRepo := repository.NewOrgRepository(db, log)
	teamRepo := repository.NewTeamRepository(db, log)
	inviteRepo := repository.NewInviteRepository(db, log)
	
	orgService := service.NewOrgService(orgRepo, teamRepo, inviteRepo, producer, log)
	orgHandler := handler.NewOrgHandler(orgService, log)

	// Create gRPC server with auth interceptor
	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			middleware.AuthInterceptor(),
		),
	)

	// Register services
	pb.RegisterOrgServiceServer(grpcServer, orgHandler)

	// Register health check
	healthServer := health.NewServer()
	grpc_health_v1.RegisterHealthServer(grpcServer, healthServer)
	healthServer.SetServingStatus("", grpc_health_v1.HealthCheckResponse_SERVING)

	// Register reflection service (for development)
	reflection.Register(grpcServer)

	// Start gRPC server
	serverCfg := cfg.GetServer()
	// Use a different port than user-service (50051)
	// Default to 50052 if not specified
	port := serverCfg.GRPCPort
	if port == 0 || port == 9090 { // Default in config package might be 9090
		port = 50052
	}
	
	addr := fmt.Sprintf("%s:%d", serverCfg.Host, port)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Sugar().Fatalw("Failed to listen", "error", err, "addr", addr)
	}

	// Start server in goroutine
	go func() {
		log.Sugar().Infow("gRPC server listening", "addr", addr)
		if err := grpcServer.Serve(listener); err != nil {
			log.Sugar().Fatal("Failed to serve")
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

// runMigrations runs database migrations
func runMigrations(db *sql.DB, log *logger.Logger) error {
	driver, err := postgres.WithInstance(db, &postgres.Config{
		MigrationsTable: "schema_migrations_org",
	})
	if err != nil {
		return fmt.Errorf("failed to create migration driver: %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"postgres",
		driver,
	)
	if err != nil {
		return fmt.Errorf("failed to create migration instance: %w", err)
	}

	log.Sugar().Infow("Running database migrations...")

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	version, dirty, err := m.Version()
	if err != nil && err != migrate.ErrNilVersion {
		return fmt.Errorf("failed to get migration version: %w", err)
	}

	log.Sugar().Infow("Database migrations complete", "version", version, "dirty", dirty)
	return nil
}
