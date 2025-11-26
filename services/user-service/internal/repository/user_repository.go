package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/nexusflow/nexusflow/pkg/database"
	"github.com/nexusflow/nexusflow/pkg/logger"
	"github.com/nexusflow/nexusflow/services/user-service/internal/models"
	"github.com/uptrace/bun"
)

// UserRepository handles data access for users
type UserRepository struct {
	db  *database.MultiTenantDB
	log *logger.Logger
}

// NewUserRepository creates a new user repository
func NewUserRepository(db *database.DB, log *logger.Logger) *UserRepository {
	return &UserRepository{
		db:  database.NewMultiTenant(db),
		log: log,
	}
}

// Create creates a new user
func (r *UserRepository) Create(ctx context.Context, user *models.User) error {
	if err := user.BeforeInsert(); err != nil {
		return fmt.Errorf("before insert: %w", err)
	}

	_, err := r.db.NewInsert().
		Model(user).
		Exec(ctx)

	if err != nil {
		r.log.Sugar().Errorw("Failed to create user", "error", err, "email", user.Email)
		return fmt.Errorf("create user: %w", err)
	}

	r.log.Sugar().Infow("User created", "user_id", user.ID, "email", user.Email)
	return nil
}

// GetByID retrieves a user by ID
func (r *UserRepository) GetByID(ctx context.Context, id string) (*models.User, error) {
	var user models.User
	err := r.db.NewSelect().
		Model(&user).
		Where("id = ?", id).
		Where("deleted_at IS NULL").
		Scan(ctx)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("user not found: %s", id)
	}
	if err != nil {
		r.log.Sugar().Errorw("Failed to get user by ID", "error", err, "user_id", id)
		return nil, fmt.Errorf("get user by id: %w", err)
	}

	return &user, nil
}

// GetByEmail retrieves a user by email
func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	err := r.db.NewSelect().
		Model(&user).
		Where("email = ?", email).
		Where("deleted_at IS NULL").
		Scan(ctx)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("user not found: %s", email)
	}
	if err != nil {
		r.log.Sugar().Errorw("Failed to get user by email", "error", err, "email", email)
		return nil, fmt.Errorf("get user by email: %w", err)
	}

	return &user, nil
}

// Update updates a user
func (r *UserRepository) Update(ctx context.Context, user *models.User) error {
	result, err := r.db.NewUpdate().
		Model(user).
		Where("id = ?", user.ID).
		Where("version = ?", user.Version).
		Where("deleted_at IS NULL").
		Exec(ctx)

	if err != nil {
		r.log.Sugar().Errorw("Failed to update user", "error", err, "user_id", user.ID)
		return fmt.Errorf("update user: %w", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("user not found or version mismatch: %s", user.ID)
	}

	r.log.Sugar().Infow("User updated", "user_id", user.ID)
	return nil
}

// Delete soft deletes a user
func (r *UserRepository) Delete(ctx context.Context, id string) error {
	now := time.Now()
	result, err := r.db.NewUpdate().
		Model((*models.User)(nil)).
		Set("deleted_at = ?", now).
		Where("id = ?", id).
		Where("deleted_at IS NULL").
		Exec(ctx)

	if err != nil {
		r.log.Sugar().Errorw("Failed to delete user", "error", err, "user_id", id)
		return fmt.Errorf("delete user: %w", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("user not found: %s", id)
	}

	r.log.Sugar().Infow("User deleted", "user_id", id)
	return nil
}

// List retrieves users with pagination
func (r *UserRepository) List(ctx context.Context, orgID string, limit, offset int) ([]*models.User, int, error) {
	var users []*models.User

	query := r.db.NewSelect().
		Model(&users).
		Where("organization_id = ?", orgID).
		Where("deleted_at IS NULL").
		Order("created_at DESC")

	// Get total count
	count, err := query.Count(ctx)
	if err != nil {
		r.log.Sugar().Errorw("Failed to count users", "error", err, "org_id", orgID)
		return nil, 0, fmt.Errorf("count users: %w", err)
	}

	// Get paginated results
	err = query.
		Limit(limit).
		Offset(offset).
		Scan(ctx)

	if err != nil {
		r.log.Sugar().Errorw("Failed to list users", "error", err, "org_id", orgID)
		return nil, 0, fmt.Errorf("list users: %w", err)
	}

	return users, count, nil
}

// Search searches for users by query
func (r *UserRepository) Search(ctx context.Context, query string, orgIDs []string, limit, offset int) ([]*models.User, int, error) {
	var users []*models.User

	q := r.db.NewSelect().
		Model(&users).
		Where("deleted_at IS NULL")

	if len(orgIDs) > 0 {
		q = q.Where("organization_id IN (?)", bun.In(orgIDs))
	}

	if query != "" {
		searchPattern := "%" + query + "%"
		q = q.Where("(email ILIKE ? OR display_name ILIKE ?)", searchPattern, searchPattern)
	}

	q = q.Order("created_at DESC")

	// Get total count
	count, err := q.Count(ctx)
	if err != nil {
		r.log.Sugar().Errorw("Failed to count search results", "error", err, "query", query)
		return nil, 0, fmt.Errorf("count search results: %w", err)
	}

	// Get paginated results
	err = q.
		Limit(limit).
		Offset(offset).
		Scan(ctx)

	if err != nil {
		r.log.Sugar().Errorw("Failed to search users", "error", err, "query", query)
		return nil, 0, fmt.Errorf("search users: %w", err)
	}

	return users, count, nil
}

// UpdatePreferences updates user preferences
func (r *UserRepository) UpdatePreferences(ctx context.Context, id string, preferences map[string]interface{}) error {
	result, err := r.db.NewUpdate().
		Model((*models.User)(nil)).
		Set("preferences = ?", preferences).
		Where("id = ?", id).
		Where("deleted_at IS NULL").
		Exec(ctx)

	if err != nil {
		r.log.Sugar().Errorw("Failed to update preferences", "error", err, "user_id", id)
		return fmt.Errorf("update preferences: %w", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("user not found: %s", id)
	}

	r.log.Sugar().Infow("User preferences updated", "user_id", id)
	return nil
}

// UpdateLastLogin updates the last login timestamp
func (r *UserRepository) UpdateLastLogin(ctx context.Context, id string) error {
	now := time.Now()
	result, err := r.db.NewUpdate().
		Model((*models.User)(nil)).
		Set("last_login_at = ?", now).
		Where("id = ?", id).
		Where("deleted_at IS NULL").
		Exec(ctx)

	if err != nil {
		r.log.Sugar().Errorw("Failed to update last login", "error", err, "user_id", id)
		return fmt.Errorf("update last login: %w", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("user not found: %s", id)
	}

	return nil
}

// GetByVerificationToken retrieves a user by verification token
func (r *UserRepository) GetByVerificationToken(ctx context.Context, token string) (*models.User, error) {
	var user models.User
	err := r.db.NewSelect().
		Model(&user).
		Where("verification_token = ?", token).
		Where("deleted_at IS NULL").
		Scan(ctx)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("user not found with verification token")
	}
	if err != nil {
		r.log.Sugar().Errorw("Failed to get user by verification token", "error", err)
		return nil, fmt.Errorf("get user by verification token: %w", err)
	}

	return &user, nil
}

// GetByResetToken retrieves a user by password reset token
func (r *UserRepository) GetByResetToken(ctx context.Context, token string) (*models.User, error) {
	var user models.User
	err := r.db.NewSelect().
		Model(&user).
		Where("reset_token = ?", token).
		Where("deleted_at IS NULL").
		Scan(ctx)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("user not found with reset token")
	}
	if err != nil {
		r.log.Sugar().Errorw("Failed to get user by reset token", "error", err)
		return nil, fmt.Errorf("get user by reset token: %w", err)
	}

	return &user, nil
}
