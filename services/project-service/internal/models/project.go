package models

import (
	"time"

	"github.com/uptrace/bun"
)

// ProjectType represents project type
type ProjectType string

const (
	ProjectTypeKanban      ProjectType = "kanban"
	ProjectTypeScrum       ProjectType = "scrum"
	ProjectTypeBugTracking ProjectType = "bug_tracking"
)

// ProjectStatus represents project status
type ProjectStatus string

const (
	ProjectStatusActive   ProjectStatus = "active"
	ProjectStatusArchived ProjectStatus = "archived"
	ProjectStatusDeleted  ProjectStatus = "deleted"
)

// Project represents a project
type Project struct {
	bun.BaseModel `bun:"table:projects,alias:p"`

	ID             string            `bun:"id,pk,type:uuid,default:gen_random_uuid()"`
	OrganizationID string            `bun:"organization_id,notnull,type:uuid"`
	Key            string            `bun:"key,notnull"`
	Name           string            `bun:"name,notnull"`
	Description    string            `bun:"description"`
	AvatarURL      string            `bun:"avatar_url"`
	Type           ProjectType       `bun:"type,notnull,default:'kanban'"`
	Status         ProjectStatus     `bun:"status,notnull,default:'active'"`
	LeadID         string            `bun:"lead_id,type:uuid,nullzero"`
	Settings       map[string]string `bun:"settings,type:jsonb"`
	CreatedAt      time.Time         `bun:"created_at,nullzero,notnull,default:current_timestamp"`
	UpdatedAt      time.Time         `bun:"updated_at,nullzero,notnull,default:current_timestamp"`
	DeletedAt      time.Time         `bun:"deleted_at,soft_delete,nullzero"`
	Version        int64             `bun:"version,notnull,default:1"`
}

// ProjectRole represents member role in project
type ProjectRole string

const (
	ProjectRoleAdmin  ProjectRole = "admin"
	ProjectRoleMember ProjectRole = "member"
	ProjectRoleViewer ProjectRole = "viewer"
)

// ProjectMember represents a project member
type ProjectMember struct {
	bun.BaseModel `bun:"table:project_members,alias:pm"`

	ID        string      `bun:"id,pk,type:uuid,default:gen_random_uuid()"`
	ProjectID string      `bun:"project_id,notnull,type:uuid"`
	UserID    string      `bun:"user_id,notnull,type:uuid"`
	Role      ProjectRole `bun:"role,notnull,default:'member'"`
	JoinedAt  time.Time   `bun:"joined_at,nullzero,notnull,default:current_timestamp"`
	UpdatedAt time.Time   `bun:"updated_at,nullzero,notnull,default:current_timestamp"`
}
