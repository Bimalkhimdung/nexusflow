package main

import (
	"database/sql"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/gorilla/websocket"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"

	"github.com/nexusflow/nexusflow/pkg/config"
	"github.com/nexusflow/nexusflow/pkg/database"
	"github.com/nexusflow/nexusflow/pkg/logger"
	pb "github.com/nexusflow/nexusflow/pkg/proto/notification/v1"
	"github.com/nexusflow/nexusflow/services/notification-service/internal/handler"
	"github.com/nexusflow/nexusflow/services/notification-service/internal/repository"
	"github.com/nexusflow/nexusflow/services/notification-service/internal/service"
	ws "github.com/nexusflow/nexusflow/services/notification-service/internal/websocket"
)

const serviceName = "notification-service"

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins for development
	},
}

func main() {
	// Initialize logger
	log, err := logger.NewDefault(serviceName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}
	defer func() { _ = log.Sync() }()

	log.Sugar().Infow("Starting notification-service")

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

	// Initialize WebSocket hub
	hub := ws.NewHub(log)
	go hub.Run()

	// Initialize layers
	repo := repository.NewNotificationRepository(db, log)
	svc := service.NewNotificationService(repo, hub, log)
	h := handler.NewNotificationHandler(svc, log)

	// Create gRPC server
	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			// Add interceptors here
		),
	)

	// Register services
	pb.RegisterNotificationServiceServer(grpcServer, h)

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
		grpcPort = 50059
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

	// Start HTTP server for WebSocket
	httpPort := serverCfg.Port
	if httpPort == 0 {
		httpPort = 8085
	}
	httpAddr := fmt.Sprintf("%s:%d", serverCfg.Host, httpPort)

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		userID := r.URL.Query().Get("user_id")
		if userID == "" {
			http.Error(w, "user_id required", http.StatusBadRequest)
			return
		}

		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Sugar().Errorw("Failed to upgrade connection", "error", err)
			return
		}

		client := &ws.Client{
			UserID: userID,
			Conn:   conn,
		}

		hub.Register(client)

		// Keep connection alive
		defer hub.Unregister(client)

		for {
			_, _, err := conn.ReadMessage()
			if err != nil {
				break
			}
		}
	})

	// Start HTTP server in goroutine
	go func() {
		log.Sugar().Infow("HTTP server listening for WebSocket", "addr", httpAddr)
		if err := http.ListenAndServe(httpAddr, nil); err != nil {
			log.Sugar().Fatal("Failed to serve HTTP")
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
		MigrationsTable: "schema_migrations_notification",
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
