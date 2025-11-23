package auth

import (
	"context"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// AuthInterceptor is a gRPC interceptor for authentication
type AuthInterceptor struct {
	// In the future, we can add a token validator here
	// validator TokenValidator
}

// NewAuthInterceptor creates a new auth interceptor
func NewAuthInterceptor() *AuthInterceptor {
	return &AuthInterceptor{}
}

// Unary returns a unary server interceptor
func (i *AuthInterceptor) Unary() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		// Skip auth for health check and reflection
		if isPublicMethod(info.FullMethod) {
			return handler(ctx, req)
		}

		newCtx, err := i.authenticate(ctx)
		if err != nil {
			return nil, err
		}

		return handler(newCtx, req)
	}
}

// authenticate validates the token and extracts user info
func (i *AuthInterceptor) authenticate(ctx context.Context) (context.Context, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "metadata is not provided")
	}

	// TODO: Implement real token validation (JWT/Ory)
	// For now, we'll look for x-user-id header for development
	
	values := md["x-user-id"]
	if len(values) == 0 {
		// For development, allow unauthenticated if configured, or fail
		// return nil, status.Error(codes.Unauthenticated, "authorization token is not provided")
		
		// Temporary: Create a guest/system context if no header
		// This allows us to test without passing headers for now, but we should enforce it later
		return ctx, nil
	}

	userID := values[0]
	
	// Extract other fields if present
	orgID := ""
	if v := md["x-org-id"]; len(v) > 0 {
		orgID = v[0]
	}
	
	role := ""
	if v := md["x-role"]; len(v) > 0 {
		role = v[0]
	}

	userCtx := &UserContext{
		UserID:         userID,
		OrganizationID: orgID,
		Role:           role,
	}

	return NewContext(ctx, userCtx), nil
}

// isPublicMethod checks if the method is public
func isPublicMethod(method string) bool {
	publicMethods := []string{
		"/grpc.health.v1.Health/Check",
		"/grpc.health.v1.Health/Watch",
		"/grpc.reflection.v1alpha.ServerReflection/ServerReflectionInfo",
		"/grpc.reflection.v1.ServerReflection/ServerReflectionInfo",
		// Add login/register methods here
	}

	for _, m := range publicMethods {
		if method == m {
			return true
		}
	}
	
	return strings.HasPrefix(method, "/grpc.reflection")
}
