package rbac

// Policy defines the mapping between roles and permissions
var Policy = map[Role][]Permission{
	RoleOwner: {
		PermOrgRead, PermOrgUpdate, PermOrgDelete,
		PermMemberAdd, PermMemberRemove, PermMemberUpdate, PermMemberRead,
		PermTeamCreate, PermTeamUpdate, PermTeamDelete, PermTeamRead,
		PermInviteCreate, PermInviteRevoke, PermInviteRead,
	},
	RoleAdmin: {
		PermOrgRead, PermOrgUpdate,
		PermMemberAdd, PermMemberRemove, PermMemberUpdate, PermMemberRead,
		PermTeamCreate, PermTeamUpdate, PermTeamDelete, PermTeamRead,
		PermInviteCreate, PermInviteRevoke, PermInviteRead,
	},
	RoleMember: {
		PermOrgRead,
		PermMemberRead,
		PermTeamRead,
		PermInviteRead,
	},
	RoleGuest: {
		PermOrgRead,
	},
}

// GetPermissions returns the permissions for a role
func GetPermissions(role Role) []Permission {
	return Policy[role]
}

// HasPermission checks if a role has a specific permission
func HasPermission(role Role, perm Permission) bool {
	perms := GetPermissions(role)
	for _, p := range perms {
		if p == perm {
			return true
		}
	}
	return false
}
