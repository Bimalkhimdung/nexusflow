package middleware

import (
	"context"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// ContextKey is a custom type for context keys
type ContextKey string

const (
	// UserIDKey is the context key for user ID
	UserIDKey ContextKey = "user_id"
	// OrgIDKey is the context key for organization ID
	OrgIDKey ContextKey = "org_id"
)

// AuthInterceptor extracts user information from gRPC metadata and adds it to context
func AuthInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		// Extract metadata from context
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			// For now, allow requests without metadata (development mode)
			// In production, this should return an error
			return handler(ctx, req)
		}

		// Extract user ID from metadata
		userIDs := md.Get("x-user-id")
		if len(userIDs) > 0 {
			ctx = context.WithValue(ctx, UserIDKey, userIDs[0])
		}

		// Extract organization ID from metadata
		orgIDs := md.Get("x-org-id")
		if len(orgIDs) > 0 {
			ctx = context.WithValue(ctx, OrgIDKey, orgIDs[0])
		}

		// Extract from authorization header (JWT token)
		authHeaders := md.Get("authorization")
		if len(authHeaders) > 0 {
			token := strings.TrimPrefix(authHeaders[0], "Bearer ")
			// TODO: Validate JWT token and extract user ID
			// For now, we'll use the x-user-id header
			_ = token
		}

		return handler(ctx, req)
	}
}

// GetUserIDFromContext extracts user ID from context
func GetUserIDFromContext(ctx context.Context) (string, error) {
	userID, ok := ctx.Value(UserIDKey).(string)
	if !ok || userID == "" {
		return "", status.Error(codes.Unauthenticated, "user not authenticated")
	}
	return userID, nil
}

// GetOrgIDFromContext extracts organization ID from context
func GetOrgIDFromContext(ctx context.Context) (string, error) {
	orgID, ok := ctx.Value(OrgIDKey).(string)
	if !ok || orgID == "" {
		return "", status.Error(codes.InvalidArgument, "organization ID not provided")
	}
	return orgID, nil
}

// GetUserIDFromContextOrDefault extracts user ID from context or returns default
func GetUserIDFromContextOrDefault(ctx context.Context, defaultID string) string {
	userID, err := GetUserIDFromContext(ctx)
	if err != nil {
		return defaultID
	}
	return userID
}
