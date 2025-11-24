package service

import (
	"context"
	"fmt"
	"regexp"
	"time"

	"github.com/nexusflow/nexusflow/pkg/kafka"
	"github.com/nexusflow/nexusflow/pkg/logger"
	"github.com/nexusflow/nexusflow/services/git-service/internal/models"
	"github.com/nexusflow/nexusflow/services/git-service/internal/repository"
)

type GitService struct {
	repo     *repository.GitRepository
	producer *kafka.Producer
	log      *logger.Logger
	keyRegex *regexp.Regexp
}

func NewGitService(repo *repository.GitRepository, producer *kafka.Producer, log *logger.Logger) *GitService {
	// Regex to match issue keys like PROJ-123
	keyRegex := regexp.MustCompile(`([A-Z]+-\d+)`)
	return &GitService{repo: repo, producer: producer, log: log, keyRegex: keyRegex}
}

// ConnectRepository links an external repository to a project
func (s *GitService) ConnectRepository(ctx context.Context, repo *models.Repository) (*models.Repository, error) {
	// TODO: Validate provider and external ID
	if err := s.repo.CreateRepository(ctx, repo); err != nil {
		return nil, err
	}
	return repo, nil
}

// ListRepositories lists repositories for a project
func (s *GitService) ListRepositories(ctx context.Context, projectID string) ([]*models.Repository, error) {
	return s.repo.ListRepositories(ctx, projectID)
}

// ProcessCommit processes a commit from a webhook
func (s *GitService) ProcessCommit(ctx context.Context, repo *models.Repository, commit *models.Commit) error {
	// Save commit
	commit.RepositoryID = repo.ID
	if err := s.repo.CreateCommit(ctx, commit); err != nil {
		return fmt.Errorf("failed to create commit: %w", err)
	}

	// Find issue keys in message
	keys := s.findIssueKeys(commit.Message)
	for _, key := range keys {
		// TODO: Resolve issue ID from key (needs Issue Service integration or lookup table)
		// For now, we'll assume we can resolve it or store the key directly if we change the schema
		// Since we don't have direct access to Issue Service database, we'd typically call it via gRPC
		// But to keep it simple for now, we'll skip the actual linking if we don't have the ID
		// In a real implementation, we would query the Issue Service to get the ID from the Key
		
		s.log.Sugar().Infow("Found issue key in commit", "key", key, "commit", commit.Hash)
		
		// Publish event for other services to handle linking/transitioning
		s.publishEvent("git.commit_pushed", repo.ProjectID, map[string]interface{}{
			"commit_id":   commit.ID,
			"issue_key":   key,
			"message":     commit.Message,
			"repository":  repo.Name,
			"author":      commit.AuthorName,
			"url":         commit.URL,
		})
	}

	return nil
}

// GetIssueCommits gets commits linked to an issue
func (s *GitService) GetIssueCommits(ctx context.Context, issueID string) ([]*models.Commit, error) {
	return s.repo.GetIssueCommits(ctx, issueID)
}

// GetIssuePullRequests gets PRs linked to an issue
func (s *GitService) GetIssuePullRequests(ctx context.Context, issueID string) ([]*models.PullRequest, error) {
	return s.repo.GetIssuePullRequests(ctx, issueID)
}

// Helper functions

func (s *GitService) findIssueKeys(message string) []string {
	matches := s.keyRegex.FindAllString(message, -1)
	// Deduplicate
	keys := make(map[string]bool)
	var result []string
	for _, m := range matches {
		if !keys[m] {
			keys[m] = true
			result = append(result, m)
		}
	}
	return result
}

func (s *GitService) publishEvent(eventType, projectID string, payload map[string]interface{}) {
	if s.producer == nil {
		return
	}
	event := kafka.Event{Type: eventType, Timestamp: time.Now(), Payload: payload}
	if projectID != "" {
		payload["project_id"] = projectID
	}
	_ = s.producer.PublishEvent("git-events", event)
}
