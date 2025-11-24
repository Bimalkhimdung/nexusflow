package models

import (
	"encoding/json"
	"time"
)

// NotificationType constants
const (
	NotificationTypeIssueAssigned   = "issue.assigned"
	NotificationTypeIssueUpdated    = "issue.updated"
	NotificationTypeCommentCreated  = "comment.created"
	NotificationTypeCommentMention  = "comment.mention"
	NotificationTypeSprintStarted   = "sprint.started"
	NotificationTypeSprintCompleted = "sprint.completed"
)

type Notification struct {
	ID        string          `bun:"type:uuid,pk,default:uuid_generate_v4()"`
	UserID    string          `bun:"type:uuid,notnull"`
	Type      string          `bun:"type:text,notnull"`
	Title     string          `bun:"type:text,notnull"`
	Message   string          `bun:"type:text,notnull"`
	Link      string          `bun:"type:text"`
	Metadata  json.RawMessage `bun:"type:jsonb"`
	Read      bool            `bun:"type:boolean,notnull,default:false"`
	CreatedAt time.Time       `bun:"type:timestamp,notnull,default:now()"`
}

type NotificationPreference struct {
	ID               string    `bun:"type:uuid,pk,default:uuid_generate_v4()"`
	UserID           string    `bun:"type:uuid,notnull"`
	NotificationType string    `bun:"type:text,notnull"`
	InAppEnabled     bool      `bun:"type:boolean,notnull,default:true"`
	EmailEnabled     bool      `bun:"type:boolean,notnull,default:true"`
	CreatedAt        time.Time `bun:"type:timestamp,notnull,default:now()"`
	UpdatedAt        time.Time `bun:"type:timestamp,notnull,default:now()"`
}
