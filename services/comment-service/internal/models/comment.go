package models

import "time"

type Comment struct {
	ID        string     `bun:"type:uuid,pk,default:uuid_generate_v4()"`
	IssueID   string     `bun:"type:uuid,notnull"`
	AuthorID  string     `bun:"type:uuid,notnull"`
	ParentID  *string    `bun:"type:uuid"`
	Content   string     `bun:"type:text,notnull"`
	CreatedAt time.Time  `bun:"type:timestamp,notnull,default:now()"`
	UpdatedAt time.Time  `bun:"type:timestamp,notnull,default:now()"`
	DeletedAt *time.Time `bun:"type:timestamp"`
}

type CommentReaction struct {
	ID        string    `bun:"type:uuid,pk,default:uuid_generate_v4()"`
	CommentID string    `bun:"type:uuid,notnull"`
	UserID    string    `bun:"type:uuid,notnull"`
	Emoji     string    `bun:"type:text,notnull"`
	CreatedAt time.Time `bun:"type:timestamp,notnull,default:now()"`
}

type CommentMention struct {
	ID              string    `bun:"type:uuid,pk,default:uuid_generate_v4()"`
	CommentID       string    `bun:"type:uuid,notnull"`
	MentionedUserID string    `bun:"type:uuid,notnull"`
	CreatedAt       time.Time `bun:"type:timestamp,notnull,default:now()"`
}
