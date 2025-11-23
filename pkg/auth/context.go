package auth

import (
	"context"
	"fmt"

	"google.golang.org/grpc/metadata"
)

type contextKey string

const (
	userContextKey contextKey = "user_context"
)

// UserContext holds user identity information
type UserContext struct {
	UserID         string
	OrganizationID string
	Role           string
	Email          string
}

// FromContext extracts UserContext from context
func FromContext(ctx context.Context) (*UserContext, bool) {
	u, ok := ctx.Value(userContextKey).(*UserContext)
	return u, ok
}

// NewContext injects UserContext into context
func NewContext(ctx context.Context, u *UserContext) context.Context {
	return context.WithValue(ctx, userContextKey, u)
}

// GetUserID returns the user ID from context
func GetUserID(ctx context.Context) (string, error) {
	u, ok := FromContext(ctx)
	if !ok {
		return "", fmt.Errorf("user context not found")
	}
	return u.UserID, nil
}

// GetOrganizationID returns the organization ID from context
func GetOrganizationID(ctx context.Context) (string, error) {
	u, ok := FromContext(ctx)
	if !ok {
		return "", fmt.Errorf("user context not found")
	}
	return u.OrganizationID, nil
}

// ExtractMetadata extracts metadata from gRPC context
func ExtractMetadata(ctx context.Context) metadata.MD {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil
	}
	return md
}
