package models

import "time"

type GitProvider struct {
	ID        string    `bun:"type:uuid,pk,default:uuid_generate_v4()"`
	Name      string    `bun:"type:text,notnull"` // github, gitlab
	BaseURL   string    `bun:"type:text,notnull"`
	CreatedAt time.Time `bun:"type:timestamp,notnull,default:now()"`
}

type Repository struct {
	ID            string    `bun:"type:uuid,pk,default:uuid_generate_v4()"`
	ProviderID    string    `bun:"type:uuid,notnull"`
	ExternalID    string    `bun:"type:text,notnull"`
	Name          string    `bun:"type:text,notnull"`
	URL           string    `bun:"type:text,notnull"`
	ProjectID     string    `bun:"type:uuid,notnull"`
	WebhookSecret string    `bun:"type:text"`
	CreatedAt     time.Time `bun:"type:timestamp,notnull,default:now()"`
}

type Commit struct {
	ID           string    `bun:"type:uuid,pk,default:uuid_generate_v4()"`
	RepositoryID string    `bun:"type:uuid,notnull"`
	Hash         string    `bun:"type:text,notnull"`
	Message      string    `bun:"type:text,notnull"`
	AuthorName   string    `bun:"type:text,notnull"`
	AuthorEmail  string    `bun:"type:text,notnull"`
	URL          string    `bun:"type:text,notnull"`
	CommittedAt  time.Time `bun:"type:timestamp,notnull"`
	CreatedAt    time.Time `bun:"type:timestamp,notnull,default:now()"`
}

type PullRequest struct {
	ID           string    `bun:"type:uuid,pk,default:uuid_generate_v4()"`
	RepositoryID string    `bun:"type:uuid,notnull"`
	ExternalID   string    `bun:"type:text,notnull"`
	Title        string    `bun:"type:text,notnull"`
	Status       string    `bun:"type:text,notnull"` // open, closed, merged
	URL          string    `bun:"type:text,notnull"`
	AuthorName   string    `bun:"type:text,notnull"`
	CreatedAt    time.Time `bun:"type:timestamp,notnull,default:now()"`
	UpdatedAt    time.Time `bun:"type:timestamp,notnull,default:now()"`
}

type IssueCommit struct {
	IssueID  string `bun:"type:uuid,pk"`
	CommitID string `bun:"type:uuid,pk"`
}

type IssuePullRequest struct {
	IssueID       string `bun:"type:uuid,pk"`
	PullRequestID string `bun:"type:uuid,pk"`
}
