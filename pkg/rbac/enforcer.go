package rbac

import (
	"context"
	"fmt"

	"github.com/nexusflow/nexusflow/pkg/auth"
)

// Enforcer handles permission checks
type Enforcer struct{}

// NewEnforcer creates a new enforcer
func NewEnforcer() *Enforcer {
	return &Enforcer{}
}

// Enforce checks if the user in the context has the required permission
func (e *Enforcer) Enforce(ctx context.Context, perm Permission) error {
	userCtx, ok := auth.FromContext(ctx)
	if !ok {
		return fmt.Errorf("unauthenticated: user context missing")
	}

	role := Role(userCtx.Role)
	if !role.IsValid() {
		return fmt.Errorf("invalid role: %s", role)
	}

	if !HasPermission(role, perm) {
		return fmt.Errorf("permission denied: role %s does not have permission %s", role, perm)
	}

	return nil
}

// EnforceRole checks if the user has a specific role or higher
// This is a simplified check, usually we prefer permission-based checks
func (e *Enforcer) EnforceRole(ctx context.Context, requiredRole Role) error {
	userCtx, ok := auth.FromContext(ctx)
	if !ok {
		return fmt.Errorf("unauthenticated: user context missing")
	}

	role := Role(userCtx.Role)
	
	// Simple hierarchy check (Owner > Admin > Member > Guest)
	if role == RoleOwner {
		return nil
	}
	if role == RoleAdmin && requiredRole != RoleOwner {
		return nil
	}
	if role == RoleMember && (requiredRole == RoleMember || requiredRole == RoleGuest) {
		return nil
	}
	if role == RoleGuest && requiredRole == RoleGuest {
		return nil
	}

	return fmt.Errorf("permission denied: role %s is not sufficient", role)
}
