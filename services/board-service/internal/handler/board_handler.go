package handler

import (
    "context"
    "time"
    
    "github.com/nexusflow/nexusflow/pkg/logger"
    pb "github.com/nexusflow/nexusflow/pkg/proto/board/v1"
    "github.com/nexusflow/nexusflow/services/board-service/internal/models"
    "github.com/nexusflow/nexusflow/services/board-service/internal/service"
    "google.golang.org/grpc/codes"
    "google.golang.org/grpc/status"
)

// BoardHandler implements the BoardService gRPC server.
type BoardHandler struct {
    pb.UnimplementedBoardServiceServer
    svc *service.BoardService
    log *logger.Logger
}

func NewBoardHandler(svc *service.BoardService, log *logger.Logger) *BoardHandler {
    return &BoardHandler{svc: svc, log: log}
}

// Helper conversions
func boardToProto(b *models.Board) *pb.Board {
    if b == nil {
        return nil
    }
    return &pb.Board{
        Id:          b.ID,
        ProjectId:   b.ProjectID,
        Name:        b.Name,
        Description: b.Description,
        CreatedAt:   b.CreatedAt.Format(time.RFC3339),
        UpdatedAt:   b.UpdatedAt.Format(time.RFC3339),
    }
}

func cardToProto(c *models.Card) *pb.Card {
    if c == nil {
        return nil
    }
    return &pb.Card{
        Id:        c.ID,
        BoardId:   c.BoardID,
        IssueId:   c.IssueID,
        Position:  int32(c.Position),
        CreatedAt: c.CreatedAt.Format(time.RFC3339),
        UpdatedAt: c.UpdatedAt.Format(time.RFC3339),
    }
}

// CreateBoard
func (h *BoardHandler) CreateBoard(ctx context.Context, req *pb.CreateBoardRequest) (*pb.CreateBoardResponse, error) {
    board, err := h.svc.CreateBoard(ctx, req)
    if err != nil {
        h.log.Sugar().Errorw("Failed to create board", "error", err)
        return nil, status.Errorf(codes.Internal, "failed to create board: %v", err)
    }
    return &pb.CreateBoardResponse{Board: boardToProto(board)}, nil
}

func (h *BoardHandler) GetBoard(ctx context.Context, req *pb.GetBoardRequest) (*pb.GetBoardResponse, error) {
    board, err := h.svc.GetBoard(ctx, req.Id)
    if err != nil {
        h.log.Sugar().Errorw("Failed to get board", "error", err)
        return nil, status.Errorf(codes.Internal, "failed to get board: %v", err)
    }
    if board == nil {
        return nil, status.Error(codes.NotFound, "board not found")
    }
    return &pb.GetBoardResponse{Board: boardToProto(board)}, nil
}

func (h *BoardHandler) ListBoards(ctx context.Context, req *pb.ListBoardsRequest) (*pb.ListBoardsResponse, error) {
    boards, err := h.svc.ListBoards(ctx, req.ProjectId)
    if err != nil {
        h.log.Sugar().Errorw("Failed to list boards", "error", err)
        return nil, status.Errorf(codes.Internal, "failed to list boards: %v", err)
    }
    var pbBoards []*pb.Board
    for _, b := range boards {
        pbBoards = append(pbBoards, boardToProto(b))
    }
    return &pb.ListBoardsResponse{Boards: pbBoards}, nil
}

func (h *BoardHandler) UpdateBoard(ctx context.Context, req *pb.UpdateBoardRequest) (*pb.UpdateBoardResponse, error) {
    board, err := h.svc.UpdateBoard(ctx, req)
    if err != nil {
        h.log.Sugar().Errorw("Failed to update board", "error", err)
        return nil, status.Errorf(codes.Internal, "failed to update board: %v", err)
    }
    return &pb.UpdateBoardResponse{Board: boardToProto(board)}, nil
}

func (h *BoardHandler) DeleteBoard(ctx context.Context, req *pb.DeleteBoardRequest) (*pb.DeleteBoardResponse, error) {
    if err := h.svc.DeleteBoard(ctx, req.Id); err != nil {
        h.log.Sugar().Errorw("Failed to delete board", "error", err)
        return nil, status.Errorf(codes.Internal, "failed to delete board: %v", err)
    }
    return &pb.DeleteBoardResponse{}, nil
}

// Card operations
func (h *BoardHandler) AddCard(ctx context.Context, req *pb.AddCardRequest) (*pb.AddCardResponse, error) {
    card, err := h.svc.AddCard(ctx, req)
    if err != nil {
        h.log.Sugar().Errorw("Failed to add card", "error", err)
        return nil, status.Errorf(codes.Internal, "failed to add card: %v", err)
    }
    return &pb.AddCardResponse{Card: cardToProto(card)}, nil
}

func (h *BoardHandler) MoveCard(ctx context.Context, req *pb.MoveCardRequest) (*pb.MoveCardResponse, error) {
    card, err := h.svc.MoveCard(ctx, req)
    if err != nil {
        h.log.Sugar().Errorw("Failed to move card", "error", err)
        return nil, status.Errorf(codes.Internal, "failed to move card: %v", err)
    }
    return &pb.MoveCardResponse{Card: cardToProto(card)}, nil
}

func (h *BoardHandler) DeleteCard(ctx context.Context, req *pb.DeleteCardRequest) (*pb.DeleteCardResponse, error) {
    if err := h.svc.DeleteCard(ctx, req.CardId); err != nil {
        h.log.Sugar().Errorw("Failed to delete card", "error", err)
        return nil, status.Errorf(codes.Internal, "failed to delete card: %v", err)
    }
    return &pb.DeleteCardResponse{}, nil
}

func (h *BoardHandler) ListCards(ctx context.Context, req *pb.ListCardsRequest) (*pb.ListCardsResponse, error) {
    cards, err := h.svc.ListCards(ctx, req.BoardId)
    if err != nil {
        h.log.Sugar().Errorw("Failed to list cards", "error", err)
        return nil, status.Errorf(codes.Internal, "failed to list cards: %v", err)
    }
    var pbCards []*pb.Card
    for _, c := range cards {
        pbCards = append(pbCards, cardToProto(c))
    }
    return &pb.ListCardsResponse{Cards: pbCards}, nil
}
