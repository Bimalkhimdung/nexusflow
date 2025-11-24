package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/nexusflow/nexusflow/pkg/database"
	"github.com/nexusflow/nexusflow/pkg/logger"
	"github.com/nexusflow/nexusflow/services/git-service/internal/models"
)

type GitRepository struct {
	db  *database.DB
	log *logger.Logger
}

func NewGitRepository(db *database.DB, log *logger.Logger) *GitRepository {
	return &GitRepository{db: db, log: log}
}

// Repository operations

func (r *GitRepository) CreateRepository(ctx context.Context, repo *models.Repository) error {
	repo.ID = ""
	repo.CreatedAt = time.Now()
	_, err := r.db.NewInsert().Model(repo).Exec(ctx)
	if err != nil {
		return fmt.Errorf("create repository: %w", err)
	}
	return nil
}

func (r *GitRepository) GetRepositoryByExternalID(ctx context.Context, providerID, externalID string) (*models.Repository, error) {
	repo := new(models.Repository)
	err := r.db.NewSelect().Model(repo).
		Where("provider_id = ? AND external_id = ?", providerID, externalID).
		Scan(ctx)
	if err != nil {
		return nil, fmt.Errorf("get repository: %w", err)
	}
	return repo, nil
}

func (r *GitRepository) ListRepositories(ctx context.Context, projectID string) ([]*models.Repository, error) {
	var repos []*models.Repository
	err := r.db.NewSelect().Model(&repos).
		Where("project_id = ?", projectID).
		Order("created_at DESC").
		Scan(ctx)
	if err != nil {
		return nil, fmt.Errorf("list repositories: %w", err)
	}
	return repos, nil
}

// Commit operations

func (r *GitRepository) CreateCommit(ctx context.Context, commit *models.Commit) error {
	commit.ID = ""
	commit.CreatedAt = time.Now()
	_, err := r.db.NewInsert().Model(commit).Exec(ctx)
	if err != nil {
		return fmt.Errorf("create commit: %w", err)
	}
	return nil
}

func (r *GitRepository) LinkCommitToIssue(ctx context.Context, issueID, commitID string) error {
	link := &models.IssueCommit{
		IssueID:  issueID,
		CommitID: commitID,
	}
	_, err := r.db.NewInsert().Model(link).On("CONFLICT DO NOTHING").Exec(ctx)
	if err != nil {
		return fmt.Errorf("link commit to issue: %w", err)
	}
	return nil
}

func (r *GitRepository) GetIssueCommits(ctx context.Context, issueID string) ([]*models.Commit, error) {
	var commits []*models.Commit
	err := r.db.NewSelect().Model(&commits).
		Join("JOIN issue_commits ON issue_commits.commit_id = commit.id").
		Where("issue_commits.issue_id = ?", issueID).
		Order("committed_at DESC").
		Scan(ctx)
	if err != nil {
		return nil, fmt.Errorf("get issue commits: %w", err)
	}
	return commits, nil
}

// Pull Request operations

func (r *GitRepository) CreatePullRequest(ctx context.Context, pr *models.PullRequest) error {
	pr.ID = ""
	pr.CreatedAt = time.Now()
	pr.UpdatedAt = time.Now()
	_, err := r.db.NewInsert().Model(pr).Exec(ctx)
	if err != nil {
		return fmt.Errorf("create pr: %w", err)
	}
	return nil
}

func (r *GitRepository) UpdatePullRequest(ctx context.Context, pr *models.PullRequest) error {
	pr.UpdatedAt = time.Now()
	_, err := r.db.NewUpdate().Model(pr).WherePK().Exec(ctx)
	if err != nil {
		return fmt.Errorf("update pr: %w", err)
	}
	return nil
}

func (r *GitRepository) LinkPullRequestToIssue(ctx context.Context, issueID, prID string) error {
	link := &models.IssuePullRequest{
		IssueID:       issueID,
		PullRequestID: prID,
	}
	_, err := r.db.NewInsert().Model(link).On("CONFLICT DO NOTHING").Exec(ctx)
	if err != nil {
		return fmt.Errorf("link pr to issue: %w", err)
	}
	return nil
}

func (r *GitRepository) GetIssuePullRequests(ctx context.Context, issueID string) ([]*models.PullRequest, error) {
	var prs []*models.PullRequest
	err := r.db.NewSelect().Model(&prs).
		Join("JOIN issue_pull_requests ON issue_pull_requests.pull_request_id = pull_request.id").
		Where("issue_pull_requests.issue_id = ?", issueID).
		Order("updated_at DESC").
		Scan(ctx)
	if err != nil {
		return nil, fmt.Errorf("get issue prs: %w", err)
	}
	return prs, nil
}
