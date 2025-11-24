package models

type SearchQuery struct {
	Query        string
	EntityTypes  []string          // issues, projects, users
	Filters      map[string]string // status:open, assignee:user-id
	SortBy       string            // relevance, created_at, updated_at
	SortOrder    string            // asc, desc
	Limit        int
	Offset       int
}

type SearchResult struct {
	ID          string
	Type        string  // issue, project, user
	Title       string
	Description string
	Metadata    map[string]string
	Score       float64
}

type SearchResponse struct {
	Results []*SearchResult
	Total   int
	Facets  map[string]int
}

// Issue document for Elasticsearch
type IssueDocument struct {
	ID          string   `json:"id"`
	Key         string   `json:"key"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Status      string   `json:"status"`
	Priority    string   `json:"priority"`
	Type        string   `json:"type"`
	ProjectID   string   `json:"project_id"`
	AssigneeID  string   `json:"assignee_id"`
	ReporterID  string   `json:"reporter_id"`
	Labels      []string `json:"labels"`
	CreatedAt   string   `json:"created_at"`
	UpdatedAt   string   `json:"updated_at"`
}

// Project document for Elasticsearch
type ProjectDocument struct {
	ID          string `json:"id"`
	Key         string `json:"key"`
	Name        string `json:"name"`
	Description string `json:"description"`
	OrgID       string `json:"org_id"`
	LeadID      string `json:"lead_id"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

// User document for Elasticsearch
type UserDocument struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	Name      string `json:"name"`
	CreatedAt string `json:"created_at"`
}
