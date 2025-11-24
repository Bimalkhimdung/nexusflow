package models

import "time"

type SprintStatus string

const (
	SprintStatusPlanned   SprintStatus = "planned"
	SprintStatusActive    SprintStatus = "active"
	SprintStatusCompleted SprintStatus = "completed"
)

type Sprint struct {
	ID        string       `bun:"type:uuid,pk,default:uuid_generate_v4()"`
	ProjectID string       `bun:"type:uuid,notnull"`
	Name      string       `bun:"type:text,notnull"`
	Goal      string       `bun:"type:text"`
	StartDate time.Time    `bun:"type:timestamp"`
	EndDate   time.Time    `bun:"type:timestamp"`
	Status    SprintStatus `bun:"type:text,notnull,default:'planned'"`
	CreatedAt time.Time    `bun:"type:timestamp,notnull,default:now()"`
	UpdatedAt time.Time    `bun:"type:timestamp,notnull,default:now()"`
}

type SprintIssue struct {
	ID       string    `bun:"type:uuid,pk,default:uuid_generate_v4()"`
	SprintID string    `bun:"type:uuid,notnull"`
	IssueID  string    `bun:"type:uuid,notnull"`
	AddedAt  time.Time `bun:"type:timestamp,notnull,default:now()"`
}
