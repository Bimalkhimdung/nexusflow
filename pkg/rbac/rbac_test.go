package rbac

import (
	"testing"
)

func TestHasPermission(t *testing.T) {
	tests := []struct {
		name     string
		role     Role
		perm     Permission
		expected bool
	}{
		{"Owner has OrgRead", RoleOwner, PermOrgRead, true},
		{"Owner has InviteCreate", RoleOwner, PermInviteCreate, true},
		{"Admin has OrgRead", RoleAdmin, PermOrgRead, true},
		{"Admin has InviteCreate", RoleAdmin, PermInviteCreate, true},
		{"Member has OrgRead", RoleMember, PermOrgRead, true},
		{"Member has no InviteCreate", RoleMember, PermInviteCreate, false},
		{"Guest has OrgRead", RoleGuest, PermOrgRead, true},
		{"Guest has no MemberRead", RoleGuest, PermMemberRead, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := HasPermission(tt.role, tt.perm); got != tt.expected {
				t.Errorf("HasPermission() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestRole_IsValid(t *testing.T) {
	tests := []struct {
		name string
		role Role
		want bool
	}{
		{"Owner", RoleOwner, true},
		{"Admin", RoleAdmin, true},
		{"Member", RoleMember, true},
		{"Guest", RoleGuest, true},
		{"Invalid", Role("invalid"), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.role.IsValid(); got != tt.want {
				t.Errorf("Role.IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}
