package models

import (
	"time"

	"github.com/uptrace/bun"
)

// OrgStatus represents organization status
type OrgStatus string

const (
	OrgStatusActive    OrgStatus = "active"
	OrgStatusSuspended OrgStatus = "suspended"
	OrgStatusDeleted   OrgStatus = "deleted"
)

// OrgPlan represents organization plan
type OrgPlan string

const (
	OrgPlanFree       OrgPlan = "free"
	OrgPlanTeam       OrgPlan = "team"
	OrgPlanEnterprise OrgPlan = "enterprise"
)

// Organization represents an organization
type Organization struct {
	bun.BaseModel `bun:"table:organizations,alias:o"`

	ID          string            `bun:"id,pk,type:uuid,default:gen_random_uuid()"`
	Name        string            `bun:"name,notnull"`
	Slug        string            `bun:"slug,notnull,unique"`
	Description string            `bun:"description"`
	LogoURL     string            `bun:"logo_url"`
	Status      OrgStatus         `bun:"status,notnull,default:'active'"`
	Plan        OrgPlan           `bun:"plan,notnull,default:'free'"`
	Settings    map[string]string `bun:"settings,type:jsonb"`
	CreatedAt   time.Time         `bun:"created_at,nullzero,notnull,default:current_timestamp"`
	UpdatedAt   time.Time         `bun:"updated_at,nullzero,notnull,default:current_timestamp"`
	DeletedAt   time.Time         `bun:"deleted_at,soft_delete,nullzero"`
	Version     int64             `bun:"version,notnull,default:1"`
}

// OrgRole represents member role
type OrgRole string

const (
	OrgRoleOwner  OrgRole = "owner"
	OrgRoleAdmin  OrgRole = "admin"
	OrgRoleMember OrgRole = "member"
	OrgRoleGuest  OrgRole = "guest"
)

// OrgMember represents an organization member
type OrgMember struct {
	bun.BaseModel `bun:"table:org_members,alias:om"`

	ID             string    `bun:"id,pk,type:uuid,default:gen_random_uuid()"`
	OrganizationID string    `bun:"organization_id,notnull,type:uuid"`
	UserID         string    `bun:"user_id,notnull,type:uuid"`
	Role           OrgRole   `bun:"role,notnull,default:'member'"`
	JoinedAt       time.Time `bun:"joined_at,nullzero,notnull,default:current_timestamp"`
	UpdatedAt      time.Time `bun:"updated_at,nullzero,notnull,default:current_timestamp"`
}

// Team represents a team within an organization
type Team struct {
	bun.BaseModel `bun:"table:teams,alias:t"`

	ID             string    `bun:"id,pk,type:uuid,default:gen_random_uuid()"`
	OrganizationID string    `bun:"organization_id,notnull,type:uuid"`
	Name           string    `bun:"name,notnull"`
	Description    string    `bun:"description"`
	CreatedAt      time.Time `bun:"created_at,nullzero,notnull,default:current_timestamp"`
	UpdatedAt      time.Time `bun:"updated_at,nullzero,notnull,default:current_timestamp"`
	DeletedAt      time.Time `bun:"deleted_at,soft_delete,nullzero"`
}

// TeamMember represents a team member
type TeamMember struct {
	bun.BaseModel `bun:"table:team_members,alias:tm"`

	TeamID   string    `bun:"team_id,pk,type:uuid"`
	UserID   string    `bun:"user_id,pk,type:uuid"`
	JoinedAt time.Time `bun:"joined_at,nullzero,notnull,default:current_timestamp"`
}

// InviteStatus represents invite status
type InviteStatus string

const (
	InviteStatusPending  InviteStatus = "pending"
	InviteStatusAccepted InviteStatus = "accepted"
	InviteStatusExpired  InviteStatus = "expired"
	InviteStatusRevoked  InviteStatus = "revoked"
)

// Invite represents an organization invitation
type Invite struct {
	bun.BaseModel `bun:"table:invites,alias:i"`

	ID             string       `bun:"id,pk,type:uuid,default:gen_random_uuid()"`
	OrganizationID string       `bun:"organization_id,notnull,type:uuid"`
	Email          string       `bun:"email,notnull"`
	Role           OrgRole      `bun:"role,notnull,default:'member'"`
	InvitedBy      string       `bun:"invited_by,notnull,type:uuid"`
	Token          string       `bun:"token,notnull,unique"`
	Status         InviteStatus `bun:"status,notnull,default:'pending'"`
	CreatedAt      time.Time    `bun:"created_at,nullzero,notnull,default:current_timestamp"`
	ExpiresAt      time.Time    `bun:"expires_at,notnull"`
	AcceptedAt     time.Time    `bun:"accepted_at,nullzero"`
}
