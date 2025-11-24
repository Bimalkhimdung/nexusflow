package handler

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/nexusflow/nexusflow/pkg/logger"
	"github.com/nexusflow/nexusflow/services/git-service/internal/models"
	"github.com/nexusflow/nexusflow/services/git-service/internal/repository"
	"github.com/nexusflow/nexusflow/services/git-service/internal/service"
)

type WebhookHandler struct {
	svc  *service.GitService
	repo *repository.GitRepository
	log  *logger.Logger
}

func NewWebhookHandler(svc *service.GitService, repo *repository.GitRepository, log *logger.Logger) *WebhookHandler {
	return &WebhookHandler{svc: svc, repo: repo, log: log}
}

// GitHubPushPayload represents a subset of GitHub push event
type GitHubPushPayload struct {
	Ref        string `json:"ref"`
	Before     string `json:"before"`
	After      string `json:"after"`
	Repository struct {
		ID       int64  `json:"id"`
		Name     string `json:"name"`
		FullName string `json:"full_name"`
		HTMLURL  string `json:"html_url"`
	} `json:"repository"`
	Pusher struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	} `json:"pusher"`
	Commits []struct {
		ID        string `json:"id"`
		Message   string `json:"message"`
		Timestamp string `json:"timestamp"`
		URL       string `json:"url"`
		Author    struct {
			Name  string `json:"name"`
			Email string `json:"email"`
		} `json:"author"`
	} `json:"commits"`
}

func (h *WebhookHandler) HandleGitHub(w http.ResponseWriter, r *http.Request) {
	eventType := r.Header.Get("X-GitHub-Event")
	signature := r.Header.Get("X-Hub-Signature-256")

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read body", http.StatusBadRequest)
		return
	}

	// TODO: Verify signature using webhook secret from DB
	// For now, we'll skip signature verification for simplicity in this MVP
	// In production, we would look up the repo by ID and verify the secret

	if eventType == "push" {
		var payload GitHubPushPayload
		if err := json.Unmarshal(body, &payload); err != nil {
			h.log.Sugar().Errorw("Failed to unmarshal GitHub payload", "error", err)
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		// Find repository in DB
		repoID := fmt.Sprintf("%d", payload.Repository.ID)
		repo, err := h.repo.GetRepositoryByExternalID(r.Context(), "github", repoID)
		if err != nil {
			h.log.Sugar().Warnw("Repository not found", "external_id", repoID)
			// Don't error to GitHub, just ignore
			w.WriteHeader(http.StatusOK)
			return
		}

		// Process commits
		for _, c := range payload.Commits {
			committedAt, _ := time.Parse(time.RFC3339, c.Timestamp)
			commit := &models.Commit{
				Hash:        c.ID,
				Message:     c.Message,
				AuthorName:  c.Author.Name,
				AuthorEmail: c.Author.Email,
				URL:         c.URL,
				CommittedAt: committedAt,
			}

			if err := h.svc.ProcessCommit(r.Context(), repo, commit); err != nil {
				h.log.Sugar().Errorw("Failed to process commit", "error", err, "hash", c.ID)
			}
		}
	}

	w.WriteHeader(http.StatusOK)
}

func verifySignature(payload []byte, secret, signature string) bool {
	if !strings.HasPrefix(signature, "sha256=") {
		return false
	}
	sig := signature[7:] // Remove "sha256=" prefix

	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(payload)
	expectedMAC := mac.Sum(nil)
	expectedSig := hex.EncodeToString(expectedMAC)

	return hmac.Equal([]byte(sig), []byte(expectedSig))
}
