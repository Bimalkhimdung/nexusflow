package models

import (
	"time"

	"github.com/uptrace/bun"
)

// IssueType represents issue type
type IssueType string

const (
	IssueTypeEpic        IssueType = "epic"
	IssueTypeStory       IssueType = "story"
	IssueTypeTask        IssueType = "task"
	IssueTypeSubTask     IssueType = "sub_task"
	IssueTypeBug         IssueType = "bug"
	IssueTypeImprovement IssueType = "improvement"
)

// IssuePriority represents issue priority
type IssuePriority string

const (
	IssuePriorityLowest  IssuePriority = "lowest"
	IssuePriorityLow     IssuePriority = "low"
	IssuePriorityMedium  IssuePriority = "medium"
	IssuePriorityHigh    IssuePriority = "high"
	IssuePriorityHighest IssuePriority = "highest"
)

// Issue represents an issue
type Issue struct {
	bun.BaseModel `bun:"table:issues,alias:i"`

	ID          string        `bun:"id,pk,type:uuid,default:gen_random_uuid()"`
	ProjectID   string        `bun:"project_id,notnull,type:uuid"`
	Key         string        `bun:"key,notnull"`
	Summary     string        `bun:"summary,notnull"`
	Description string        `bun:"description"`
	Type        IssueType     `bun:"type,notnull,default:'task'"`
	Priority    IssuePriority `bun:"priority,notnull,default:'medium'"`
	StatusID    string        `bun:"status_id,type:uuid,nullzero"`
	AssigneeID  string        `bun:"assignee_id,type:uuid,nullzero"`
	ReporterID  string        `bun:"reporter_id,type:uuid,nullzero"`
	ParentID    string        `bun:"parent_id,type:uuid,nullzero"`
	SprintID    string        `bun:"sprint_id,type:uuid,nullzero"`
	StoryPoints int32         `bun:"story_points"`
	DueDate     time.Time     `bun:"due_date,nullzero"`
	CreatedAt   time.Time     `bun:"created_at,nullzero,notnull,default:current_timestamp"`
	UpdatedAt   time.Time     `bun:"updated_at,nullzero,notnull,default:current_timestamp"`
	DeletedAt   time.Time     `bun:"deleted_at,soft_delete,nullzero"`
	Version     int64         `bun:"version,notnull,default:1"`
}

// ProjectCounter tracks the next issue number for a project
type ProjectCounter struct {
	bun.BaseModel `bun:"table:project_counters,alias:pc"`

	ProjectID       string    `bun:"project_id,pk,type:uuid"`
	NextIssueNumber int64     `bun:"next_issue_number,notnull,default:1"`
	UpdatedAt       time.Time `bun:"updated_at,nullzero,notnull,default:current_timestamp"`
}

// CustomFieldType represents custom field type
type CustomFieldType string

const (
	CustomFieldTypeText        CustomFieldType = "text"
	CustomFieldTypeNumber      CustomFieldType = "number"
	CustomFieldTypeDate        CustomFieldType = "date"
	CustomFieldTypeSelect      CustomFieldType = "select"
	CustomFieldTypeMultiSelect CustomFieldType = "multi_select"
	CustomFieldTypeUser        CustomFieldType = "user"
	CustomFieldTypeCheckbox    CustomFieldType = "checkbox"
	CustomFieldTypeURL         CustomFieldType = "url"
	CustomFieldTypeEmail       CustomFieldType = "email"
	CustomFieldTypeTextarea    CustomFieldType = "textarea"
)

// CustomField represents a custom field definition
type CustomField struct {
	bun.BaseModel `bun:"table:custom_fields,alias:cf"`

	ID          string          `bun:"id,pk,type:uuid,default:gen_random_uuid()"`
	ProjectID   string          `bun:"project_id,notnull,type:uuid"`
	Name        string          `bun:"name,notnull"`
	Description string          `bun:"description"`
	Type        CustomFieldType `bun:"type,notnull"`
	Required    bool            `bun:"required,default:false"`
	DefaultValue interface{}    `bun:"default_value,type:jsonb"`
	Options     []string        `bun:"options,type:jsonb"`
	Config      map[string]string `bun:"config,type:jsonb"`
	CreatedAt   time.Time       `bun:"created_at,nullzero,notnull,default:current_timestamp"`
	UpdatedAt   time.Time       `bun:"updated_at,nullzero,notnull,default:current_timestamp"`
}

// IssueCustomValue represents a value for a custom field on an issue
type IssueCustomValue struct {
	bun.BaseModel `bun:"table:issue_custom_values,alias:icv"`

	IssueID   string      `bun:"issue_id,pk,type:uuid"`
	FieldID   string      `bun:"field_id,pk,type:uuid"`
	Value     interface{} `bun:"value,type:jsonb"`
	UpdatedAt time.Time   `bun:"updated_at,nullzero,notnull,default:current_timestamp"`
}

// IssueLinkType represents issue link type
type IssueLinkType string

const (
	IssueLinkTypeBlocks       IssueLinkType = "blocks"
	IssueLinkTypeBlockedBy    IssueLinkType = "blocked_by"
	IssueLinkTypeRelatesTo    IssueLinkType = "relates_to"
	IssueLinkTypeDuplicates   IssueLinkType = "duplicates"
	IssueLinkTypeDuplicatedBy IssueLinkType = "duplicated_by"
	IssueLinkTypeCauses       IssueLinkType = "causes"
	IssueLinkTypeCausedBy     IssueLinkType = "caused_by"
)

// IssueLink represents a link between issues
type IssueLink struct {
	bun.BaseModel `bun:"table:issue_links,alias:il"`

	ID            string        `bun:"id,pk,type:uuid,default:gen_random_uuid()"`
	SourceIssueID string        `bun:"source_issue_id,notnull,type:uuid"`
	TargetIssueID string        `bun:"target_issue_id,notnull,type:uuid"`
	Type          IssueLinkType `bun:"type,notnull"`
	CreatedAt     time.Time     `bun:"created_at,nullzero,notnull,default:current_timestamp"`
}

// IssueWatcher represents a user watching an issue
type IssueWatcher struct {
	bun.BaseModel `bun:"table:issue_watchers,alias:iw"`

	IssueID  string    `bun:"issue_id,pk,type:uuid"`
	UserID   string    `bun:"user_id,pk,type:uuid"`
	JoinedAt time.Time `bun:"joined_at,nullzero,notnull,default:current_timestamp"`
}
