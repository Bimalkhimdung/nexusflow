package service

import (
    "context"
    "fmt"
    "time"

    "github.com/nexusflow/nexusflow/pkg/kafka"
    "github.com/nexusflow/nexusflow/pkg/logger"
    pb "github.com/nexusflow/nexusflow/pkg/proto/board/v1"
    "github.com/nexusflow/nexusflow/services/board-service/internal/models"
    "github.com/nexusflow/nexusflow/services/board-service/internal/repository"
)

// BoardService handles business logic for boards and cards
type BoardService struct {
    repo     *repository.BoardRepository
    producer *kafka.Producer
    log      *logger.Logger
}

func NewBoardService(repo *repository.BoardRepository, producer *kafka.Producer, log *logger.Logger) *BoardService {
    return &BoardService{repo: repo, producer: producer, log: log}
}

// CreateBoard creates a new board
func (s *BoardService) CreateBoard(ctx context.Context, input *pb.CreateBoardRequest) (*models.Board, error) {
    b := &models.Board{
        ProjectID:   input.ProjectId,
        Name:        input.Name,
        Description: input.Description,
    }
    if err := s.repo.CreateBoard(ctx, b); err != nil {
        return nil, fmt.Errorf("create board: %w", err)
    }
    // publish event
    s.publishEvent("board.created", b.ProjectID, map[string]interface{}{"board_id": b.ID, "name": b.Name})
    return b, nil
}

func (s *BoardService) GetBoard(ctx context.Context, id string) (*models.Board, error) {
    return s.repo.GetBoard(ctx, id)
}

func (s *BoardService) ListBoards(ctx context.Context, projectID string) ([]*models.Board, error) {
    return s.repo.ListBoards(ctx, projectID)
}

func (s *BoardService) UpdateBoard(ctx context.Context, input *pb.UpdateBoardRequest) (*models.Board, error) {
    b, err := s.repo.GetBoard(ctx, input.Id)
    if err != nil {
        return nil, err
    }
    if input.Name != "" {
        b.Name = input.Name
    }
    if input.Description != "" {
        b.Description = input.Description
    }
    if err := s.repo.UpdateBoard(ctx, b); err != nil {
        return nil, fmt.Errorf("update board: %w", err)
    }
    s.publishEvent("board.updated", b.ProjectID, map[string]interface{}{"board_id": b.ID})
    return b, nil
}

func (s *BoardService) DeleteBoard(ctx context.Context, id string) error {
    // fetch to get projectID for event
    b, err := s.repo.GetBoard(ctx, id)
    if err != nil {
        return err
    }
    if err := s.repo.DeleteBoard(ctx, id); err != nil {
        return fmt.Errorf("delete board: %w", err)
    }
    s.publishEvent("board.deleted", b.ProjectID, map[string]interface{}{"board_id": id})
    return nil
}

// Card operations
func (s *BoardService) AddCard(ctx context.Context, input *pb.AddCardRequest) (*models.Card, error) {
    c := &models.Card{
        BoardID:  input.BoardId,
        IssueID:  input.IssueId,
        Position: int(input.Position),
    }
    if err := s.repo.AddCard(ctx, c); err != nil {
        return nil, fmt.Errorf("add card: %w", err)
    }
    s.publishEvent("card.added", "", map[string]interface{}{"card_id": c.ID, "board_id": c.BoardID, "issue_id": c.IssueID})
    return c, nil
}

func (s *BoardService) MoveCard(ctx context.Context, input *pb.MoveCardRequest) (*models.Card, error) {
    if err := s.repo.MoveCard(ctx, input.CardId, int(input.NewPosition)); err != nil {
        return nil, fmt.Errorf("move card: %w", err)
    }
    c, err := s.repo.GetCard(ctx, input.CardId)
    if err != nil {
        return nil, err
    }
    s.publishEvent("card.moved", "", map[string]interface{}{"card_id": c.ID, "new_position": c.Position})
    return c, nil
}

func (s *BoardService) DeleteCard(ctx context.Context, cardID string) error {
    if err := s.repo.DeleteCard(ctx, cardID); err != nil {
        return fmt.Errorf("delete card: %w", err)
    }
    s.publishEvent("card.deleted", "", map[string]interface{}{"card_id": cardID})
    return nil
}

func (s *BoardService) ListCards(ctx context.Context, boardID string) ([]*models.Card, error) {
    return s.repo.ListCardsByBoard(ctx, boardID)
}

func (s *BoardService) publishEvent(eventType, projectID string, payload map[string]interface{}) {
    if s.producer == nil {
        return
    }
    event := kafka.Event{Type: eventType, Timestamp: time.Now(), Payload: payload}
    if projectID != "" {
        payload["project_id"] = projectID
    }
    _ = s.producer.PublishEvent("board-events", event)
}
