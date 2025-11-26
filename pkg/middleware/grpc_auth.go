package middleware

import (
	"context"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// GRPCAuthConfig holds gRPC auth interceptor configuration
type GRPCAuthConfig struct {
	Secret        string
	SkipMethods   []string // Methods that don't require authentication
}

// UnaryServerInterceptor creates a gRPC unary interceptor for JWT authentication
func UnaryServerInterceptor(config *GRPCAuthConfig) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		// Check if method should skip authentication
		if shouldSkipMethod(info.FullMethod, config.SkipMethods) {
			return handler(ctx, req)
		}

		// Extract token from metadata
		token, err := extractTokenFromMetadata(ctx)
		if err != nil {
			return nil, status.Error(codes.Unauthenticated, "missing or invalid authorization token")
		}

		// Validate token
		claims, err := validateToken(token, config.Secret)
		if err != nil {
			return nil, status.Error(codes.Unauthenticated, "invalid token")
		}

		// Add user info to context
		ctx = context.WithValue(ctx, UserIDKey, claims.UserID)
		ctx = context.WithValue(ctx, UserEmailKey, claims.Email)

		return handler(ctx, req)
	}
}

// StreamServerInterceptor creates a gRPC stream interceptor for JWT authentication
func StreamServerInterceptor(config *GRPCAuthConfig) grpc.StreamServerInterceptor {
	return func(
		srv interface{},
		ss grpc.ServerStream,
		info *grpc.StreamServerInfo,
		handler grpc.StreamHandler,
	) error {
		// Check if method should skip authentication
		if shouldSkipMethod(info.FullMethod, config.SkipMethods) {
			return handler(srv, ss)
		}

		// Extract token from metadata
		token, err := extractTokenFromMetadata(ss.Context())
		if err != nil {
			return status.Error(codes.Unauthenticated, "missing or invalid authorization token")
		}

		// Validate token
		claims, err := validateToken(token, config.Secret)
		if err != nil {
			return status.Error(codes.Unauthenticated, "invalid token")
		}

		// Create new context with user info
		ctx := context.WithValue(ss.Context(), UserIDKey, claims.UserID)
		ctx = context.WithValue(ctx, UserEmailKey, claims.Email)

		// Wrap the stream with new context
		wrappedStream := &wrappedServerStream{
			ServerStream: ss,
			ctx:          ctx,
		}

		return handler(srv, wrappedStream)
	}
}

// extractTokenFromMetadata extracts JWT token from gRPC metadata
func extractTokenFromMetadata(ctx context.Context) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", &AuthError{Message: "missing metadata"}
	}

	// Try "authorization" header
	values := md.Get("authorization")
	if len(values) == 0 {
		return "", &AuthError{Message: "missing authorization header"}
	}

	authHeader := values[0]
	
	// Expected format: "Bearer <token>"
	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		return "", &AuthError{Message: "invalid authorization header format"}
	}

	return parts[1], nil
}

// shouldSkipMethod checks if the gRPC method should skip authentication
func shouldSkipMethod(method string, skipMethods []string) bool {
	for _, skipMethod := range skipMethods {
		if strings.HasSuffix(method, skipMethod) || method == skipMethod {
			return true
		}
	}
	return false
}

// wrappedServerStream wraps grpc.ServerStream with a custom context
type wrappedServerStream struct {
	grpc.ServerStream
	ctx context.Context
}

func (w *wrappedServerStream) Context() context.Context {
	return w.ctx
}
