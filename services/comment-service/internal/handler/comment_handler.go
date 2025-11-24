package handler

import (
	"context"
	"time"

	"github.com/nexusflow/nexusflow/pkg/logger"
	pb "github.com/nexusflow/nexusflow/pkg/proto/comment/v1"
	"github.com/nexusflow/nexusflow/services/comment-service/internal/models"
	"github.com/nexusflow/nexusflow/services/comment-service/internal/service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type CommentHandler struct {
	pb.UnimplementedCommentServiceServer
	svc *service.CommentService
	log *logger.Logger
}

func NewCommentHandler(svc *service.CommentService, log *logger.Logger) *CommentHandler {
	return &CommentHandler{svc: svc, log: log}
}

// Helper conversions
func commentToProto(c *models.Comment) *pb.Comment {
	if c == nil {
		return nil
	}
	parentID := ""
	if c.ParentID != nil {
		parentID = *c.ParentID
	}
	return &pb.Comment{
		Id:        c.ID,
		IssueId:   c.IssueID,
		AuthorId:  c.AuthorID,
		ParentId:  parentID,
		Content:   c.Content,
		CreatedAt: c.CreatedAt.Format(time.RFC3339),
		UpdatedAt: c.UpdatedAt.Format(time.RFC3339),
		IsDeleted: c.DeletedAt != nil,
	}
}

func reactionToProto(r *models.CommentReaction) *pb.Reaction {
	if r == nil {
		return nil
	}
	return &pb.Reaction{
		Id:        r.ID,
		CommentId: r.CommentID,
		UserId:    r.UserID,
		Emoji:     r.Emoji,
		CreatedAt: r.CreatedAt.Format(time.RFC3339),
	}
}

// RPC Methods
func (h *CommentHandler) CreateComment(ctx context.Context, req *pb.CreateCommentRequest) (*pb.CreateCommentResponse, error) {
	comment, err := h.svc.CreateComment(ctx, req)
	if err != nil {
		h.log.Sugar().Errorw("Failed to create comment", "error", err)
		return nil, status.Errorf(codes.Internal, "failed to create comment: %v", err)
	}
	return &pb.CreateCommentResponse{Comment: commentToProto(comment)}, nil
}

func (h *CommentHandler) GetComment(ctx context.Context, req *pb.GetCommentRequest) (*pb.GetCommentResponse, error) {
	comment, err := h.svc.GetComment(ctx, req.Id)
	if err != nil {
		h.log.Sugar().Errorw("Failed to get comment", "error", err)
		return nil, status.Errorf(codes.Internal, "failed to get comment: %v", err)
	}
	if comment == nil {
		return nil, status.Error(codes.NotFound, "comment not found")
	}
	return &pb.GetCommentResponse{Comment: commentToProto(comment)}, nil
}

func (h *CommentHandler) ListComments(ctx context.Context, req *pb.ListCommentsRequest) (*pb.ListCommentsResponse, error) {
	comments, err := h.svc.ListComments(ctx, req.IssueId)
	if err != nil {
		h.log.Sugar().Errorw("Failed to list comments", "error", err)
		return nil, status.Errorf(codes.Internal, "failed to list comments: %v", err)
	}

	var pbComments []*pb.Comment
	for _, c := range comments {
		pbComments = append(pbComments, commentToProto(c))
	}
	return &pb.ListCommentsResponse{Comments: pbComments}, nil
}

func (h *CommentHandler) UpdateComment(ctx context.Context, req *pb.UpdateCommentRequest) (*pb.UpdateCommentResponse, error) {
	comment, err := h.svc.UpdateComment(ctx, req)
	if err != nil {
		h.log.Sugar().Errorw("Failed to update comment", "error", err)
		return nil, status.Errorf(codes.Internal, "failed to update comment: %v", err)
	}
	return &pb.UpdateCommentResponse{Comment: commentToProto(comment)}, nil
}

func (h *CommentHandler) DeleteComment(ctx context.Context, req *pb.DeleteCommentRequest) (*pb.DeleteCommentResponse, error) {
	if err := h.svc.DeleteComment(ctx, req.Id); err != nil {
		h.log.Sugar().Errorw("Failed to delete comment", "error", err)
		return nil, status.Errorf(codes.Internal, "failed to delete comment: %v", err)
	}
	return &pb.DeleteCommentResponse{}, nil
}

func (h *CommentHandler) AddReaction(ctx context.Context, req *pb.AddReactionRequest) (*pb.AddReactionResponse, error) {
	if err := h.svc.AddReaction(ctx, req.CommentId, req.UserId, req.Emoji); err != nil {
		h.log.Sugar().Errorw("Failed to add reaction", "error", err)
		return nil, status.Errorf(codes.Internal, "failed to add reaction: %v", err)
	}
	return &pb.AddReactionResponse{}, nil
}

func (h *CommentHandler) RemoveReaction(ctx context.Context, req *pb.RemoveReactionRequest) (*pb.RemoveReactionResponse, error) {
	if err := h.svc.RemoveReaction(ctx, req.CommentId, req.UserId, req.Emoji); err != nil {
		h.log.Sugar().Errorw("Failed to remove reaction", "error", err)
		return nil, status.Errorf(codes.Internal, "failed to remove reaction: %v", err)
	}
	return &pb.RemoveReactionResponse{}, nil
}

func (h *CommentHandler) ListReactions(ctx context.Context, req *pb.ListReactionsRequest) (*pb.ListReactionsResponse, error) {
	reactions, err := h.svc.ListReactions(ctx, req.CommentId)
	if err != nil {
		h.log.Sugar().Errorw("Failed to list reactions", "error", err)
		return nil, status.Errorf(codes.Internal, "failed to list reactions: %v", err)
	}

	var pbReactions []*pb.Reaction
	for _, r := range reactions {
		pbReactions = append(pbReactions, reactionToProto(r))
	}
	return &pb.ListReactionsResponse{Reactions: pbReactions}, nil
}
