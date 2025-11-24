package service

import (
	"context"
	"fmt"

	"github.com/nexusflow/nexusflow/pkg/logger"
	pb "github.com/nexusflow/nexusflow/pkg/proto/search/v1"
	"github.com/nexusflow/nexusflow/services/search-service/internal/elasticsearch"
	"github.com/nexusflow/nexusflow/services/search-service/internal/models"
)

type SearchService struct {
	es  *elasticsearch.Client
	log *logger.Logger
}

func NewSearchService(es *elasticsearch.Client, log *logger.Logger) *SearchService {
	return &SearchService{es: es, log: log}
}

// Search performs a multi-index search
func (s *SearchService) Search(ctx context.Context, req *pb.SearchRequest) (*models.SearchResponse, error) {
	indices := s.getIndices(req.EntityTypes)
	if len(indices) == 0 {
		indices = []string{"issues", "projects", "users"}
	}

	query := s.buildQuery(req)
	
	index := indices[0]
	if len(indices) > 1 {
		index = fmt.Sprintf("%s", indices[0]) // Search first index for simplicity
	}

	result, err := s.es.Search(ctx, index, query)
	if err != nil {
		return nil, fmt.Errorf("search error: %w", err)
	}

	return s.parseSearchResponse(result), nil
}

// SearchIssues searches only issues
func (s *SearchService) SearchIssues(ctx context.Context, req *pb.SearchIssuesRequest) (*models.SearchResponse, error) {
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"multi_match": map[string]interface{}{
				"query":  req.Query,
				"fields": []string{"title^2", "description", "key"},
			},
		},
		"from": req.Offset,
		"size": req.Limit,
	}

	if req.SortBy != "" {
		query["sort"] = []map[string]interface{}{
			{req.SortBy: map[string]string{"order": req.SortOrder}},
		}
	}

	result, err := s.es.Search(ctx, "issues", query)
	if err != nil {
		return nil, fmt.Errorf("search issues error: %w", err)
	}

	return s.parseSearchResponse(result), nil
}

// SearchProjects searches only projects
func (s *SearchService) SearchProjects(ctx context.Context, req *pb.SearchProjectsRequest) (*models.SearchResponse, error) {
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"multi_match": map[string]interface{}{
				"query":  req.Query,
				"fields": []string{"name^2", "description", "key"},
			},
		},
		"from": req.Offset,
		"size": req.Limit,
	}

	result, err := s.es.Search(ctx, "projects", query)
	if err != nil {
		return nil, fmt.Errorf("search projects error: %w", err)
	}

	return s.parseSearchResponse(result), nil
}

// Suggest provides autocomplete suggestions
func (s *SearchService) Suggest(ctx context.Context, query string, limit int) ([]*models.SearchResult, error) {
	esQuery := map[string]interface{}{
		"query": map[string]interface{}{
			"multi_match": map[string]interface{}{
				"query":  query,
				"fields": []string{"title", "name", "key"},
				"type":   "phrase_prefix",
			},
		},
		"size": limit,
	}

	result, err := s.es.Search(ctx, "issues,projects", esQuery)
	if err != nil {
		return nil, fmt.Errorf("suggest error: %w", err)
	}

	response := s.parseSearchResponse(result)
	return response.Results, nil
}

// Helper functions
func (s *SearchService) getIndices(entityTypes []string) []string {
	if len(entityTypes) == 0 {
		return nil
	}
	return entityTypes
}

func (s *SearchService) buildQuery(req *pb.SearchRequest) map[string]interface{}{
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"multi_match": map[string]interface{}{
				"query":  req.Query,
				"fields": []string{"title^2", "name^2", "description", "key"},
			},
		},
		"from": req.Offset,
		"size": req.Limit,
	}

	if req.SortBy != "" {
		query["sort"] = []map[string]interface{}{
			{req.SortBy: map[string]string{"order": req.SortOrder}},
		}
	}

	return query
}

func (s *SearchService) parseSearchResponse(result map[string]interface{}) *models.SearchResponse {
	hits, ok := result["hits"].(map[string]interface{})
	if !ok {
		return &models.SearchResponse{Results: []*models.SearchResult{}, Total: 0}
	}

	total := 0
	if t, ok := hits["total"].(map[string]interface{}); ok {
		if v, ok := t["value"].(float64); ok {
			total = int(v)
		}
	}

	documents, ok := hits["hits"].([]interface{})
	if !ok {
		return &models.SearchResponse{Results: []*models.SearchResult{}, Total: total}
	}

	results := make([]*models.SearchResult, 0, len(documents))
	for _, doc := range documents {
		docMap, ok := doc.(map[string]interface{})
		if !ok {
			continue
		}

		source, ok := docMap["_source"].(map[string]interface{})
		if !ok {
			continue
		}

		result := &models.SearchResult{
			ID:          getString(docMap, "_id"),
			Type:        getString(docMap, "_index"),
			Title:       getString(source, "title"),
			Description: getString(source, "description"),
			Metadata:    make(map[string]string),
		}

		if score, ok := docMap["_score"].(float64); ok {
			result.Score = score
		}

		// Add name if it's a project
		if name := getString(source, "name"); name != "" {
			result.Title = name
		}

		results = append(results, result)
	}

	return &models.SearchResponse{
		Results: results,
		Total:   total,
		Facets:  make(map[string]int),
	}
}

func getString(m map[string]interface{}, key string) string {
	if v, ok := m[key].(string); ok {
		return v
	}
	return ""
}
