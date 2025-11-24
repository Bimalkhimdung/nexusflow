package repository

import (
    "context"
    "fmt"
    "time"

    "github.com/nexusflow/nexusflow/pkg/database"
    "github.com/nexusflow/nexusflow/pkg/logger"
    "github.com/nexusflow/nexusflow/services/board-service/internal/models"
)

// BoardRepository handles board and card persistence
type BoardRepository struct {
    db  *database.DB
    log *logger.Logger
}

func NewBoardRepository(db *database.DB, log *logger.Logger) *BoardRepository {
    return &BoardRepository{db: db, log: log}
}

// Board CRUD
func (r *BoardRepository) CreateBoard(ctx context.Context, board *models.Board) error {
    board.ID = ""
    board.CreatedAt = time.Now()
    board.UpdatedAt = time.Now()
    _, err := r.db.NewInsert().Model(board).Exec(ctx)
    if err != nil {
        return fmt.Errorf("create board: %w", err)
    }
    return nil
}

func (r *BoardRepository) GetBoard(ctx context.Context, id string) (*models.Board, error) {
    b := new(models.Board)
    err := r.db.NewSelect().Model(b).Where("id = ?", id).Scan(ctx)
    if err != nil {
        return nil, fmt.Errorf("get board: %w", err)
    }
    return b, nil
}

func (r *BoardRepository) ListBoards(ctx context.Context, projectID string) ([]*models.Board, error) {
    var boards []*models.Board
    err := r.db.NewSelect().Model(&boards).Where("project_id = ?", projectID).Order("created_at DESC").Scan(ctx)
    if err != nil {
        return nil, fmt.Errorf("list boards: %w", err)
    }
    return boards, nil
}

func (r *BoardRepository) UpdateBoard(ctx context.Context, board *models.Board) error {
    board.UpdatedAt = time.Now()
    _, err := r.db.NewUpdate().Model(board).WherePK().Exec(ctx)
    if err != nil {
        return fmt.Errorf("update board: %w", err)
    }
    return nil
}

func (r *BoardRepository) DeleteBoard(ctx context.Context, id string) error {
    _, err := r.db.NewDelete().Model((*models.Board)(nil)).Where("id = ?", id).Exec(ctx)
    if err != nil {
        return fmt.Errorf("delete board: %w", err)
    }
    return nil
}

// Card CRUD
func (r *BoardRepository) AddCard(ctx context.Context, card *models.Card) error {
    card.ID = ""
    card.CreatedAt = time.Now()
    card.UpdatedAt = time.Now()
    _, err := r.db.NewInsert().Model(card).Exec(ctx)
    if err != nil {
        return fmt.Errorf("add card: %w", err)
    }
    return nil
}

func (r *BoardRepository) GetCard(ctx context.Context, id string) (*models.Card, error) {
    c := new(models.Card)
    err := r.db.NewSelect().Model(c).Where("id = ?", id).Scan(ctx)
    if err != nil {
        return nil, fmt.Errorf("get card: %w", err)
    }
    return c, nil
}

func (r *BoardRepository) ListCardsByBoard(ctx context.Context, boardID string) ([]*models.Card, error) {
    var cards []*models.Card
    err := r.db.NewSelect().Model(&cards).Where("board_id = ?", boardID).Order("position ASC").Scan(ctx)
    if err != nil {
        return nil, fmt.Errorf("list cards: %w", err)
    }
    return cards, nil
}

func (r *BoardRepository) MoveCard(ctx context.Context, cardID string, newPosition int) error {
    _, err := r.db.NewUpdate().Model((*models.Card)(nil)).Set("position = ?", newPosition).Where("id = ?", cardID).Exec(ctx)
    if err != nil {
        return fmt.Errorf("move card: %w", err)
    }
    return nil
}

func (r *BoardRepository) DeleteCard(ctx context.Context, id string) error {
    _, err := r.db.NewDelete().Model((*models.Card)(nil)).Where("id = ?", id).Exec(ctx)
    if err != nil {
        return fmt.Errorf("delete card: %w", err)
    }
    return nil
}
