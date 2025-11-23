package models

import (
	"time"

	"github.com/uptrace/bun"
)

// StatusCategory represents status category
type StatusCategory string

const (
	StatusCategoryTodo       StatusCategory = "todo"
	StatusCategoryInProgress StatusCategory = "in_progress"
	StatusCategoryDone       StatusCategory = "done"
)

// Workflow represents a workflow
type Workflow struct {
	bun.BaseModel `bun:"table:workflows,alias:w"`

	ID          string    `bun:"id,pk,type:uuid,default:gen_random_uuid()"`
	ProjectID   string    `bun:"project_id,notnull,type:uuid"`
	Name        string    `bun:"name,notnull"`
	Description string    `bun:"description"`
	IsDefault   bool      `bun:"is_default,default:false"`
	CreatedAt   time.Time `bun:"created_at,nullzero,notnull,default:current_timestamp"`
	UpdatedAt   time.Time `bun:"updated_at,nullzero,notnull,default:current_timestamp"`
}

// WorkflowStatus represents a status in a workflow
type WorkflowStatus struct {
	bun.BaseModel `bun:"table:workflow_statuses,alias:ws"`

	ID          string         `bun:"id,pk,type:uuid,default:gen_random_uuid()"`
	WorkflowID  string         `bun:"workflow_id,notnull,type:uuid"`
	Name        string         `bun:"name,notnull"`
	Description string         `bun:"description"`
	Category    StatusCategory `bun:"category,notnull"`
	Color       string         `bun:"color"`
	Position    int32          `bun:"position,default:0"`
	CreatedAt   time.Time      `bun:"created_at,nullzero,notnull,default:current_timestamp"`
	UpdatedAt   time.Time      `bun:"updated_at,nullzero,notnull,default:current_timestamp"`
}

// WorkflowTransition represents a transition between statuses
type WorkflowTransition struct {
	bun.BaseModel `bun:"table:workflow_transitions,alias:wt"`

	ID            string      `bun:"id,pk,type:uuid,default:gen_random_uuid()"`
	WorkflowID    string      `bun:"workflow_id,notnull,type:uuid"`
	Name          string      `bun:"name,notnull"`
	FromStatusID  string      `bun:"from_status_id,notnull,type:uuid"`
	ToStatusID    string      `bun:"to_status_id,notnull,type:uuid"`
	Rules         interface{} `bun:"rules,type:jsonb"`
	Validators    interface{} `bun:"validators,type:jsonb"`
	PostFunctions interface{} `bun:"post_functions,type:jsonb"`
	CreatedAt     time.Time   `bun:"created_at,nullzero,notnull,default:current_timestamp"`
	UpdatedAt     time.Time   `bun:"updated_at,nullzero,notnull,default:current_timestamp"`
}
