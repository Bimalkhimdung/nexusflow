package service

import (
	"context"
	"fmt"
	"time"

	"github.com/nexusflow/nexusflow/pkg/kafka"
	"github.com/nexusflow/nexusflow/pkg/logger"
	pb "github.com/nexusflow/nexusflow/pkg/proto/comment/v1"
	"github.com/nexusflow/nexusflow/services/comment-service/internal/models"
	"github.com/nexusflow/nexusflow/services/comment-service/internal/repository"
)

type CommentService struct {
	repo     *repository.CommentRepository
	producer *kafka.Producer
	log      *logger.Logger
}

func NewCommentService(repo *repository.CommentRepository, producer *kafka.Producer, log *logger.Logger) *CommentService {
	return &CommentService{repo: repo, producer: producer, log: log}
}

// CreateComment creates a new comment and processes mentions
func (s *CommentService) CreateComment(ctx context.Context, req *pb.CreateCommentRequest) (*models.Comment, error) {
	comment := &models.Comment{
		IssueID:  req.IssueId,
		AuthorID: req.AuthorId,
		Content:  req.Content,
	}

	if req.ParentId != "" {
		comment.ParentID = &req.ParentId
	}

	if err := s.repo.CreateComment(ctx, comment); err != nil {
		return nil, fmt.Errorf("create comment: %w", err)
	}

	// Parse and create mentions
	mentions := repository.ParseMentions(comment.Content)
	for _, username := range mentions {
		// In a real implementation, you would resolve username to user ID
		// For now, we'll just store the username as the user ID
		mention := &models.CommentMention{
			CommentID:       comment.ID,
			MentionedUserID: username,
		}
		if err := s.repo.CreateMention(ctx, mention); err != nil {
			s.log.Sugar().Warnw("Failed to create mention", "error", err, "username", username)
		} else {
			s.publishEvent("comment.mention_created", comment.IssueID, map[string]interface{}{
				"comment_id":        comment.ID,
				"mentioned_user_id": username,
			})
		}
	}

	s.publishEvent("comment.created", comment.IssueID, map[string]interface{}{
		"comment_id": comment.ID,
		"author_id":  comment.AuthorID,
		"issue_id":   comment.IssueID,
	})

	return comment, nil
}

func (s *CommentService) GetComment(ctx context.Context, id string) (*models.Comment, error) {
	return s.repo.GetComment(ctx, id)
}

func (s *CommentService) ListComments(ctx context.Context, issueID string) ([]*models.Comment, error) {
	return s.repo.ListComments(ctx, issueID)
}

func (s *CommentService) UpdateComment(ctx context.Context, req *pb.UpdateCommentRequest) (*models.Comment, error) {
	comment, err := s.repo.GetComment(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	comment.Content = req.Content
	if err := s.repo.UpdateComment(ctx, comment); err != nil {
		return nil, fmt.Errorf("update comment: %w", err)
	}

	s.publishEvent("comment.updated", comment.IssueID, map[string]interface{}{
		"comment_id": comment.ID,
	})

	return comment, nil
}

func (s *CommentService) DeleteComment(ctx context.Context, id string) error {
	comment, err := s.repo.GetComment(ctx, id)
	if err != nil {
		return err
	}

	if err := s.repo.DeleteComment(ctx, id); err != nil {
		return fmt.Errorf("delete comment: %w", err)
	}

	s.publishEvent("comment.deleted", comment.IssueID, map[string]interface{}{
		"comment_id": id,
	})

	return nil
}

// Reactions
func (s *CommentService) AddReaction(ctx context.Context, commentID, userID, emoji string) error {
	if err := s.repo.AddReaction(ctx, commentID, userID, emoji); err != nil {
		return fmt.Errorf("add reaction: %w", err)
	}

	comment, _ := s.repo.GetComment(ctx, commentID)
	issueID := ""
	if comment != nil {
		issueID = comment.IssueID
	}

	s.publishEvent("comment.reaction_added", issueID, map[string]interface{}{
		"comment_id": commentID,
		"user_id":    userID,
		"emoji":      emoji,
	})

	return nil
}

func (s *CommentService) RemoveReaction(ctx context.Context, commentID, userID, emoji string) error {
	if err := s.repo.RemoveReaction(ctx, commentID, userID, emoji); err != nil {
		return fmt.Errorf("remove reaction: %w", err)
	}

	comment, _ := s.repo.GetComment(ctx, commentID)
	issueID := ""
	if comment != nil {
		issueID = comment.IssueID
	}

	s.publishEvent("comment.reaction_removed", issueID, map[string]interface{}{
		"comment_id": commentID,
		"user_id":    userID,
		"emoji":      emoji,
	})

	return nil
}

func (s *CommentService) ListReactions(ctx context.Context, commentID string) ([]*models.CommentReaction, error) {
	return s.repo.ListReactions(ctx, commentID)
}

func (s *CommentService) publishEvent(eventType, issueID string, payload map[string]interface{}) {
	if s.producer == nil {
		return
	}
	event := kafka.Event{Type: eventType, Timestamp: time.Now(), Payload: payload}
	if issueID != "" {
		payload["issue_id"] = issueID
	}
	_ = s.producer.PublishEvent("comment-events", event)
}
