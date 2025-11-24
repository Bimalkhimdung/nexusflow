package repository

import (
	"context"
	"fmt"
	"regexp"
	"time"

	"github.com/nexusflow/nexusflow/pkg/database"
	"github.com/nexusflow/nexusflow/pkg/logger"
	"github.com/nexusflow/nexusflow/services/comment-service/internal/models"
)

type CommentRepository struct {
	db  *database.DB
	log *logger.Logger
}

func NewCommentRepository(db *database.DB, log *logger.Logger) *CommentRepository {
	return &CommentRepository{db: db, log: log}
}

// Comment CRUD
func (r *CommentRepository) CreateComment(ctx context.Context, comment *models.Comment) error {
	comment.ID = ""
	comment.CreatedAt = time.Now()
	comment.UpdatedAt = time.Now()
	_, err := r.db.NewInsert().Model(comment).Exec(ctx)
	if err != nil {
		return fmt.Errorf("create comment: %w", err)
	}
	return nil
}

func (r *CommentRepository) GetComment(ctx context.Context, id string) (*models.Comment, error) {
	c := new(models.Comment)
	err := r.db.NewSelect().Model(c).Where("id = ? AND deleted_at IS NULL", id).Scan(ctx)
	if err != nil {
		return nil, fmt.Errorf("get comment: %w", err)
	}
	return c, nil
}

func (r *CommentRepository) ListComments(ctx context.Context, issueID string) ([]*models.Comment, error) {
	var comments []*models.Comment
	err := r.db.NewSelect().Model(&comments).
		Where("issue_id = ? AND deleted_at IS NULL", issueID).
		Order("created_at ASC").
		Scan(ctx)
	if err != nil {
		return nil, fmt.Errorf("list comments: %w", err)
	}
	return comments, nil
}

func (r *CommentRepository) UpdateComment(ctx context.Context, comment *models.Comment) error {
	comment.UpdatedAt = time.Now()
	_, err := r.db.NewUpdate().Model(comment).WherePK().Exec(ctx)
	if err != nil {
		return fmt.Errorf("update comment: %w", err)
	}
	return nil
}

func (r *CommentRepository) DeleteComment(ctx context.Context, id string) error {
	now := time.Now()
	_, err := r.db.NewUpdate().Model((*models.Comment)(nil)).
		Set("deleted_at = ?", now).
		Where("id = ?", id).
		Exec(ctx)
	if err != nil {
		return fmt.Errorf("delete comment: %w", err)
	}
	return nil
}

// Reactions
func (r *CommentRepository) AddReaction(ctx context.Context, commentID, userID, emoji string) error {
	reaction := &models.CommentReaction{
		CommentID: commentID,
		UserID:    userID,
		Emoji:     emoji,
		CreatedAt: time.Now(),
	}
	_, err := r.db.NewInsert().Model(reaction).Exec(ctx)
	if err != nil {
		return fmt.Errorf("add reaction: %w", err)
	}
	return nil
}

func (r *CommentRepository) RemoveReaction(ctx context.Context, commentID, userID, emoji string) error {
	_, err := r.db.NewDelete().Model((*models.CommentReaction)(nil)).
		Where("comment_id = ? AND user_id = ? AND emoji = ?", commentID, userID, emoji).
		Exec(ctx)
	if err != nil {
		return fmt.Errorf("remove reaction: %w", err)
	}
	return nil
}

func (r *CommentRepository) ListReactions(ctx context.Context, commentID string) ([]*models.CommentReaction, error) {
	var reactions []*models.CommentReaction
	err := r.db.NewSelect().Model(&reactions).
		Where("comment_id = ?", commentID).
		Order("created_at ASC").
		Scan(ctx)
	if err != nil {
		return nil, fmt.Errorf("list reactions: %w", err)
	}
	return reactions, nil
}

// Mentions
func (r *CommentRepository) CreateMention(ctx context.Context, mention *models.CommentMention) error {
	mention.ID = ""
	mention.CreatedAt = time.Now()
	_, err := r.db.NewInsert().Model(mention).Exec(ctx)
	if err != nil {
		return fmt.Errorf("create mention: %w", err)
	}
	return nil
}

func (r *CommentRepository) ListMentions(ctx context.Context, commentID string) ([]*models.CommentMention, error) {
	var mentions []*models.CommentMention
	err := r.db.NewSelect().Model(&mentions).
		Where("comment_id = ?", commentID).
		Scan(ctx)
	if err != nil {
		return nil, fmt.Errorf("list mentions: %w", err)
	}
	return mentions, nil
}

// ParseMentions extracts @mentions from comment content
func ParseMentions(content string) []string {
	re := regexp.MustCompile(`@([a-zA-Z0-9_-]+)`)
	matches := re.FindAllStringSubmatch(content, -1)
	
	mentions := make([]string, 0, len(matches))
	seen := make(map[string]bool)
	
	for _, match := range matches {
		if len(match) > 1 && !seen[match[1]] {
			mentions = append(mentions, match[1])
			seen[match[1]] = true
		}
	}
	return mentions
}
