package models

import "time"

type Attachment struct {
	ID               string    `bun:"type:uuid,pk,default:uuid_generate_v4()"`
	EntityType       string    `bun:"type:text,notnull"`
	EntityID         string    `bun:"type:uuid,notnull"`
	Filename         string    `bun:"type:text,notnull"`
	OriginalFilename string    `bun:"type:text,notnull"`
	ContentType      string    `bun:"type:text,notnull"`
	Size             int64     `bun:"type:bigint,notnull"`
	StoragePath      string    `bun:"type:text,notnull"`
	UploaderID       string    `bun:"type:uuid,notnull"`
	CreatedAt        time.Time `bun:"type:timestamp,notnull,default:now()"`
}

type UploadMetadata struct {
	EntityType  string
	EntityID    string
	Filename    string
	ContentType string
	UploaderID  string
}
