package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/nexusflow/nexusflow/pkg/database"
)

// UserStatus represents the status of a user
type UserStatus string

const (
	UserStatusActive    UserStatus = "active"
	UserStatusInactive  UserStatus = "inactive"
	UserStatusSuspended UserStatus = "suspended"
)

// User represents a user in the system
type User struct {
	database.BaseModel `bun:",embed"`
	
	Email       string                 `bun:"email,notnull,unique"`
	DisplayName string                 `bun:"display_name,notnull"`
	AvatarURL   string                 `bun:"avatar_url"`
	Timezone    string                 `bun:"timezone,notnull,default:'UTC'"`
	Locale      string                 `bun:"locale,notnull,default:'en-US'"`
	Status      UserStatus             `bun:"status,notnull,default:'active'"`
	Preferences map[string]interface{} `bun:"preferences,type:jsonb,default:'{}'"`
	LastLoginAt *time.Time             `bun:"last_login_at"`
	DeletedAt   *time.Time             `bun:"deleted_at"`
}

// TableName returns the table name for the User model
func (User) TableName() string {
	return "users"
}

// BeforeInsert is called before inserting a user
func (u *User) BeforeInsert() error {
	if u.ID == "" {
		u.ID = uuid.New().String()
	}
	if u.Timezone == "" {
		u.Timezone = "UTC"
	}
	if u.Locale == "" {
		u.Locale = "en-US"
	}
	if u.Status == "" {
		u.Status = UserStatusActive
	}
	if u.Preferences == nil {
		u.Preferences = make(map[string]interface{})
	}
	return nil
}

// IsActive returns true if the user is active
func (u *User) IsActive() bool {
	return u.Status == UserStatusActive && u.DeletedAt == nil
}

// IsDeleted returns true if the user is soft deleted
func (u *User) IsDeleted() bool {
	return u.DeletedAt != nil
}

// SoftDelete marks the user as deleted
func (u *User) SoftDelete() {
	now := time.Now()
	u.DeletedAt = &now
}
