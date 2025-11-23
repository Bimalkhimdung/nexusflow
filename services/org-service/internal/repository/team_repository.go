package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/nexusflow/nexusflow/pkg/database"
	"github.com/nexusflow/nexusflow/pkg/logger"
	"github.com/nexusflow/nexusflow/services/org-service/internal/models"
)

// TeamRepository handles team data access
type TeamRepository struct {
	db  *database.DB
	log *logger.Logger
}

// NewTeamRepository creates a new team repository
func NewTeamRepository(db *database.DB, log *logger.Logger) *TeamRepository {
	return &TeamRepository{
		db:  db,
		log: log,
	}
}

// Create creates a new team
func (r *TeamRepository) Create(ctx context.Context, team *models.Team) error {
	if team.ID == "" {
		team.ID = uuid.New().String()
	}
	
	_, err := r.db.NewInsert().Model(team).Exec(ctx)
	if err != nil {
		r.log.Sugar().Errorw("Failed to create team", "error", err, "name", team.Name)
		return fmt.Errorf("create team: %w", err)
	}
	
	return nil
}

// GetByID gets a team by ID
func (r *TeamRepository) GetByID(ctx context.Context, id string) (*models.Team, error) {
	team := new(models.Team)
	err := r.db.NewSelect().Model(team).Where("id = ?", id).Scan(ctx)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		r.log.Sugar().Errorw("Failed to get team", "error", err, "id", id)
		return nil, fmt.Errorf("get team: %w", err)
	}
	return team, nil
}

// Update updates a team
func (r *TeamRepository) Update(ctx context.Context, team *models.Team) error {
	team.UpdatedAt = time.Now()
	
	_, err := r.db.NewUpdate().
		Model(team).
		WherePK().
		Exec(ctx)
		
	if err != nil {
		r.log.Sugar().Errorw("Failed to update team", "error", err, "id", team.ID)
		return fmt.Errorf("update team: %w", err)
	}
	
	return nil
}

// Delete soft deletes a team
func (r *TeamRepository) Delete(ctx context.Context, id string) error {
	_, err := r.db.NewDelete().
		Model((*models.Team)(nil)).
		Where("id = ?", id).
		Exec(ctx)
		
	if err != nil {
		r.log.Sugar().Errorw("Failed to delete team", "error", err, "id", id)
		return fmt.Errorf("delete team: %w", err)
	}
	
	return nil
}

// List lists teams in an organization
func (r *TeamRepository) List(ctx context.Context, orgID string, limit, offset int) ([]*models.Team, int, error) {
	var teams []*models.Team
	
	count, err := r.db.NewSelect().
		Model(&teams).
		Where("organization_id = ?", orgID).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		ScanAndCount(ctx)
		
	if err != nil {
		r.log.Sugar().Errorw("Failed to list teams", "error", err, "org_id", orgID)
		return nil, 0, fmt.Errorf("list teams: %w", err)
	}
	
	return teams, count, nil
}

// AddMember adds a user to a team
func (r *TeamRepository) AddMember(ctx context.Context, teamID, userID string) error {
	member := &models.TeamMember{
		TeamID:   teamID,
		UserID:   userID,
		JoinedAt: time.Now(),
	}
	
	_, err := r.db.NewInsert().Model(member).Exec(ctx)
	if err != nil {
		r.log.Sugar().Errorw("Failed to add team member", "error", err, "team_id", teamID, "user_id", userID)
		return fmt.Errorf("add team member: %w", err)
	}
	
	return nil
}

// RemoveMember removes a user from a team
func (r *TeamRepository) RemoveMember(ctx context.Context, teamID, userID string) error {
	_, err := r.db.NewDelete().
		Model((*models.TeamMember)(nil)).
		Where("team_id = ? AND user_id = ?", teamID, userID).
		Exec(ctx)
		
	if err != nil {
		r.log.Sugar().Errorw("Failed to remove team member", "error", err, "team_id", teamID, "user_id", userID)
		return fmt.Errorf("remove team member: %w", err)
	}
	
	return nil
}

// GetMembers gets members of a team
func (r *TeamRepository) GetMembers(ctx context.Context, teamID string) ([]string, error) {
	var userIDs []string
	
	err := r.db.NewSelect().
		Model((*models.TeamMember)(nil)).
		Column("user_id").
		Where("team_id = ?", teamID).
		Scan(ctx, &userIDs)
		
	if err != nil {
		r.log.Sugar().Errorw("Failed to get team members", "error", err, "team_id", teamID)
		return nil, fmt.Errorf("get team members: %w", err)
	}
	
	return userIDs, nil
}
