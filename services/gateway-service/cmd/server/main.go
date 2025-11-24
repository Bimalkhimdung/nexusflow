package main

import (
	"context"
	"log"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	projectv1 "github.com/nexusflow/nexusflow/pkg/proto/project/v1"
	issuev1 "github.com/nexusflow/nexusflow/pkg/proto/issue/v1"
)

func main() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Register gRPC server endpoint
	// Note: Make sure the gRPC server is running on this address
	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	// Register Project Service
	// Assuming project-service runs on localhost:50053 (default in main.go)
	err := projectv1.RegisterProjectServiceHandlerFromEndpoint(ctx, mux, "localhost:50053", opts)
	if err != nil {
		log.Fatalf("Failed to register project service handler: %v", err)
	}

	// Register Issue Service
	// Assuming issue-service runs on localhost:50054
	err = issuev1.RegisterIssueServiceHandlerFromEndpoint(ctx, mux, "localhost:50054", opts)
	if err != nil {
		log.Fatalf("Failed to register issue service handler: %v", err)
	}

	// Register other services here as we add annotations...

	// CORS Middleware
	handler := allowCORS(mux)

	log.Println("Gateway running on :8000")
	if err := http.ListenAndServe(":8000", handler); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}

func allowCORS(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		h.ServeHTTP(w, r)
	})
}
