package rbac

// Permission represents a specific action on a resource
type Permission string

const (
	// Organization permissions
	PermOrgRead   Permission = "org:read"
	PermOrgUpdate Permission = "org:update"
	PermOrgDelete Permission = "org:delete"
	
	// Member permissions
	PermMemberAdd    Permission = "member:add"
	PermMemberRemove Permission = "member:remove"
	PermMemberUpdate Permission = "member:update"
	PermMemberRead   Permission = "member:read"
	
	// Team permissions
	PermTeamCreate Permission = "team:create"
	PermTeamUpdate Permission = "team:update"
	PermTeamDelete Permission = "team:delete"
	PermTeamRead   Permission = "team:read"
	
	// Invite permissions
	PermInviteCreate Permission = "invite:create"
	PermInviteRevoke Permission = "invite:revoke"
	PermInviteRead   Permission = "invite:read"
)

// String returns the string representation of the permission
func (p Permission) String() string {
	return string(p)
}
