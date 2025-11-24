package handler

import (
	"context"

	"github.com/nexusflow/nexusflow/pkg/logger"
	pb "github.com/nexusflow/nexusflow/pkg/proto/search/v1"
	"github.com/nexusflow/nexusflow/services/search-service/internal/service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type SearchHandler struct {
	pb.UnimplementedSearchServiceServer
	svc *service.SearchService
	log *logger.Logger
}

func NewSearchHandler(svc *service.SearchService, log *logger.Logger) *SearchHandler {
	return &SearchHandler{svc: svc, log: log}
}

// Search performs a multi-entity search
func (h *SearchHandler) Search(ctx context.Context, req *pb.SearchRequest) (*pb.SearchResponse, error) {
	if req.Limit <= 0 {
		req.Limit = 20
	}

	result, err := h.svc.Search(ctx, req)
	if err != nil {
		h.log.Sugar().Errorw("Failed to search", "error", err)
		return nil, status.Errorf(codes.Internal, "failed to search: %v", err)
	}

	var pbResults []*pb.SearchResult
	for _, r := range result.Results {
		pbResults = append(pbResults, &pb.SearchResult{
			Id:          r.ID,
			Type:        r.Type,
			Title:       r.Title,
			Description: r.Description,
			Metadata:    r.Metadata,
			Score:       float32(r.Score),
		})
	}

	facets := make(map[string]int32)
	for k, v := range result.Facets {
		facets[k] = int32(v)
	}

	return &pb.SearchResponse{
		Results: pbResults,
		Total:   int32(result.Total),
		Facets:  facets,
	}, nil
}

// SearchIssues searches only issues
func (h *SearchHandler) SearchIssues(ctx context.Context, req *pb.SearchIssuesRequest) (*pb.SearchIssuesResponse, error) {
	if req.Limit <= 0 {
		req.Limit = 20
	}

	result, err := h.svc.SearchIssues(ctx, req)
	if err != nil {
		h.log.Sugar().Errorw("Failed to search issues", "error", err)
		return nil, status.Errorf(codes.Internal, "failed to search issues: %v", err)
	}

	var pbResults []*pb.SearchResult
	for _, r := range result.Results {
		pbResults = append(pbResults, &pb.SearchResult{
			Id:          r.ID,
			Type:        r.Type,
			Title:       r.Title,
			Description: r.Description,
			Metadata:    r.Metadata,
			Score:       float32(r.Score),
		})
	}

	return &pb.SearchIssuesResponse{
		Results: pbResults,
		Total:   int32(result.Total),
	}, nil
}

// SearchProjects searches only projects
func (h *SearchHandler) SearchProjects(ctx context.Context, req *pb.SearchProjectsRequest) (*pb.SearchProjectsResponse, error) {
	if req.Limit <= 0 {
		req.Limit = 20
	}

	result, err := h.svc.SearchProjects(ctx, req)
	if err != nil {
		h.log.Sugar().Errorw("Failed to search projects", "error", err)
		return nil, status.Errorf(codes.Internal, "failed to search projects: %v", err)
	}

	var pbResults []*pb.SearchResult
	for _, r := range result.Results {
		pbResults = append(pbResults, &pb.SearchResult{
			Id:          r.ID,
			Type:        r.Type,
			Title:       r.Title,
			Description: r.Description,
			Metadata:    r.Metadata,
			Score:       float32(r.Score),
		})
	}

	return &pb.SearchProjectsResponse{
		Results: pbResults,
		Total:   int32(result.Total),
	}, nil
}

// Suggest provides autocomplete suggestions
func (h *SearchHandler) Suggest(ctx context.Context, req *pb.SuggestRequest) (*pb.SuggestResponse, error) {
	if req.Limit <= 0 {
		req.Limit = 5
	}

	results, err := h.svc.Suggest(ctx, req.Query, int(req.Limit))
	if err != nil {
		h.log.Sugar().Errorw("Failed to get suggestions", "error", err)
		return nil, status.Errorf(codes.Internal, "failed to get suggestions: %v", err)
	}

	var suggestions []*pb.Suggestion
	for _, r := range results {
		suggestions = append(suggestions, &pb.Suggestion{
			Text:  r.Title,
			Type:  r.Type,
			Score: float32(r.Score),
		})
	}

	return &pb.SuggestResponse{Suggestions: suggestions}, nil
}
