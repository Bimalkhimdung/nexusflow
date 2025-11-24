package models

import "time"

type Board struct {
    ID          string    `bun:"type:uuid,default:uuid_generate_v4()"`
    ProjectID   string    `bun:"type:uuid,notnull"`
    Name        string    `bun:"type:text,notnull"`
    Description string    `bun:"type:text"`
    CreatedAt   time.Time `bun:"type:timestamp,default:now()"`
    UpdatedAt   time.Time `bun:"type:timestamp,default:now()"`
}

type Card struct {
    ID        string    `bun:"type:uuid,default:uuid_generate_v4()"`
    BoardID   string    `bun:"type:uuid,notnull"`
    IssueID   string    `bun:"type:uuid,notnull"`
    Position  int       `bun:"type:int,notnull"`
    CreatedAt time.Time `bun:"type:timestamp,default:now()"`
    UpdatedAt time.Time `bun:"type:timestamp,default:now()"`
}
