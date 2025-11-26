package models

import (
	"crypto/rand"
	"encoding/hex"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
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
	
	// Basic Info
	OrganizationID string                 `bun:"organization_id,notnull"`
	Email          string                 `bun:"email,notnull,unique"`
	DisplayName string                 `bun:"display_name,notnull"`
	AvatarURL   string                 `bun:"avatar_url"`
	Timezone    string                 `bun:"timezone,notnull,default:'UTC'"`
	Locale      string                 `bun:"locale,notnull,default:'en-US'"`
	Status      UserStatus             `bun:"status,notnull,default:'active'"`
	Preferences map[string]interface{} `bun:"preferences,type:jsonb,default:'{}'"`
	LastLoginAt *time.Time             `bun:"last_login_at"`
	DeletedAt   *time.Time             `bun:"deleted_at"`
	
	// Authentication Fields
	PasswordHash      string     `bun:"password_hash"`
	EmailVerified     bool       `bun:"email_verified,default:false"`
	VerificationToken string     `bun:"verification_token"`
	ResetToken        string     `bun:"reset_token"`
	ResetTokenExpiry  *time.Time `bun:"reset_token_expiry"`
	
	// OAuth Fields
	OAuthProvider    string `bun:"oauth_provider"` // 'google', 'github', 'email'
	OAuthID          string `bun:"oauth_id"`
	OAuthAccessToken string `bun:"oauth_access_token"`
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
	if u.OAuthProvider == "" {
		u.OAuthProvider = "email"
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

// SetPassword hashes and sets the user's password
func (u *User) SetPassword(password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.PasswordHash = string(hashedPassword)
	return nil
}

// CheckPassword verifies if the provided password matches the stored hash
func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password))
	return err == nil
}

// GenerateVerificationToken creates a new email verification token
func (u *User) GenerateVerificationToken() (string, error) {
	token, err := generateRandomToken(32)
	if err != nil {
		return "", err
	}
	u.VerificationToken = token
	return token, nil
}

// GenerateResetToken creates a new password reset token with expiry
func (u *User) GenerateResetToken() (string, error) {
	token, err := generateRandomToken(32)
	if err != nil {
		return "", err
	}
	u.ResetToken = token
	expiry := time.Now().Add(24 * time.Hour) // Token valid for 24 hours
	u.ResetTokenExpiry = &expiry
	return token, nil
}

// IsResetTokenValid checks if the reset token is still valid
func (u *User) IsResetTokenValid() bool {
	if u.ResetToken == "" || u.ResetTokenExpiry == nil {
		return false
	}
	return time.Now().Before(*u.ResetTokenExpiry)
}

// ClearResetToken clears the password reset token
func (u *User) ClearResetToken() {
	u.ResetToken = ""
	u.ResetTokenExpiry = nil
}

// MarkEmailAsVerified marks the user's email as verified
func (u *User) MarkEmailAsVerified() {
	u.EmailVerified = true
	u.VerificationToken = ""
}

// UpdateLastLogin updates the last login timestamp
func (u *User) UpdateLastLogin() {
	now := time.Now()
	u.LastLoginAt = &now
}

// generateRandomToken generates a random hex token of specified byte length
func generateRandomToken(byteLength int) (string, error) {
	bytes := make([]byte, byteLength)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
