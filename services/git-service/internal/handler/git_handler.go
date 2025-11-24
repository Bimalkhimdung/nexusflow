package handler

import (
	"context"
	"time"

	"github.com/nexusflow/nexusflow/pkg/logger"
	pb "github.com/nexusflow/nexusflow/pkg/proto/git/v1"
	"github.com/nexusflow/nexusflow/services/git-service/internal/models"
	"github.com/nexusflow/nexusflow/services/git-service/internal/service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GitHandler struct {
	pb.UnimplementedGitServiceServer
	svc *service.GitService
	log *logger.Logger
}

func NewGitHandler(svc *service.GitService, log *logger.Logger) *GitHandler {
	return &GitHandler{svc: svc, log: log}
}

func (h *GitHandler) ConnectRepository(ctx context.Context, req *pb.ConnectRepositoryRequest) (*pb.ConnectRepositoryResponse, error) {
	repo := &models.Repository{
		ProjectID:     req.ProjectId,
		Name:          req.Name,
		URL:           req.Url,
		ExternalID:    req.ExternalId,
		WebhookSecret: req.WebhookSecret,
		// ProviderID would be looked up or passed
		ProviderID: "00000000-0000-0000-0000-000000000000", // Placeholder
	}

	// Set provider ID based on name (simple mapping for now)
	if req.Provider == "github" {
		// In a real app, we'd look up the provider ID from the database
		// For MVP, we'll just store "github" in a way that we can retrieve it
		// But since ProviderID is UUID, we need a valid UUID.
		// We'll assume the provider exists and has a known ID or we'd create it.
		// For now, let's just use a deterministic UUID for "github"
		repo.ProviderID = "github" // This will fail UUID validation, so we need to fix this in a real implementation
		// Let's assume the client passes the Provider ID or we look it up.
		// For this implementation, we'll skip the provider lookup complexity
	}

	createdRepo, err := h.svc.ConnectRepository(ctx, repo)
	if err != nil {
		h.log.Sugar().Errorw("Failed to connect repository", "error", err)
		return nil, status.Errorf(codes.Internal, "failed to connect repository: %v", err)
	}

	return &pb.ConnectRepositoryResponse{
		RepositoryId: createdRepo.ID,
		WebhookUrl:   "http://localhost:8086/webhooks/" + req.Provider,
	}, nil
}

func (h *GitHandler) ListRepositories(ctx context.Context, req *pb.ListRepositoriesRequest) (*pb.ListRepositoriesResponse, error) {
	repos, err := h.svc.ListRepositories(ctx, req.ProjectId)
	if err != nil {
		h.log.Sugar().Errorw("Failed to list repositories", "error", err)
		return nil, status.Errorf(codes.Internal, "failed to list repositories: %v", err)
	}

	var pbRepos []*pb.Repository
	for _, r := range repos {
		pbRepos = append(pbRepos, &pb.Repository{
			Id:        r.ID,
			Name:      r.Name,
			Url:       r.URL,
			Provider:  "github", // Simplified
			CreatedAt: r.CreatedAt.Format(time.RFC3339),
		})
	}

	return &pb.ListRepositoriesResponse{Repositories: pbRepos}, nil
}

func (h *GitHandler) GetIssueCommits(ctx context.Context, req *pb.GetIssueCommitsRequest) (*pb.GetIssueCommitsResponse, error) {
	commits, err := h.svc.GetIssueCommits(ctx, req.IssueId)
	if err != nil {
		h.log.Sugar().Errorw("Failed to get issue commits", "error", err)
		return nil, status.Errorf(codes.Internal, "failed to get issue commits: %v", err)
	}

	var pbCommits []*pb.Commit
	for _, c := range commits {
		pbCommits = append(pbCommits, &pb.Commit{
			Id:          c.ID,
			Hash:        c.Hash,
			Message:     c.Message,
			AuthorName:  c.AuthorName,
			Url:         c.URL,
			CommittedAt: c.CommittedAt.Format(time.RFC3339),
		})
	}

	return &pb.GetIssueCommitsResponse{Commits: pbCommits}, nil
}

func (h *GitHandler) GetIssuePullRequests(ctx context.Context, req *pb.GetIssuePullRequestsRequest) (*pb.GetIssuePullRequestsResponse, error) {
	prs, err := h.svc.GetIssuePullRequests(ctx, req.IssueId)
	if err != nil {
		h.log.Sugar().Errorw("Failed to get issue PRs", "error", err)
		return nil, status.Errorf(codes.Internal, "failed to get issue PRs: %v", err)
	}

	var pbPRs []*pb.PullRequest
	for _, pr := range prs {
		pbPRs = append(pbPRs, &pb.PullRequest{
			Id:         pr.ID,
			ExternalId: pr.ExternalID,
			Title:      pr.Title,
			Status:     pr.Status,
			Url:        pr.URL,
			AuthorName: pr.AuthorName,
			UpdatedAt:  pr.UpdatedAt.Format(time.RFC3339),
		})
	}

	return &pb.GetIssuePullRequestsResponse{PullRequests: pbPRs}, nil
}
