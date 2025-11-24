package handler

import (
	"context"
	"time"

	"github.com/nexusflow/nexusflow/pkg/logger"
	pb "github.com/nexusflow/nexusflow/pkg/proto/sprint/v1"
	"github.com/nexusflow/nexusflow/services/sprint-service/internal/models"
	"github.com/nexusflow/nexusflow/services/sprint-service/internal/service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type SprintHandler struct {
	pb.UnimplementedSprintServiceServer
	svc *service.SprintService
	log *logger.Logger
}

func NewSprintHandler(svc *service.SprintService, log *logger.Logger) *SprintHandler {
	return &SprintHandler{svc: svc, log: log}
}

// Helper conversions
func sprintToProto(s *models.Sprint) *pb.Sprint {
	if s == nil {
		return nil
	}
	return &pb.Sprint{
		Id:        s.ID,
		ProjectId: s.ProjectID,
		Name:      s.Name,
		Goal:      s.Goal,
		StartDate: s.StartDate.Format(time.RFC3339),
		EndDate:   s.EndDate.Format(time.RFC3339),
		Status:    statusToProto(s.Status),
		CreatedAt: s.CreatedAt.Format(time.RFC3339),
		UpdatedAt: s.UpdatedAt.Format(time.RFC3339),
	}
}

func statusToProto(s models.SprintStatus) pb.SprintStatus {
	switch s {
	case models.SprintStatusPlanned:
		return pb.SprintStatus_SPRINT_STATUS_PLANNED
	case models.SprintStatusActive:
		return pb.SprintStatus_SPRINT_STATUS_ACTIVE
	case models.SprintStatusCompleted:
		return pb.SprintStatus_SPRINT_STATUS_COMPLETED
	default:
		return pb.SprintStatus_SPRINT_STATUS_PLANNED
	}
}

func statusFromProto(s pb.SprintStatus) models.SprintStatus {
	switch s {
	case pb.SprintStatus_SPRINT_STATUS_PLANNED:
		return models.SprintStatusPlanned
	case pb.SprintStatus_SPRINT_STATUS_ACTIVE:
		return models.SprintStatusActive
	case pb.SprintStatus_SPRINT_STATUS_COMPLETED:
		return models.SprintStatusCompleted
	default:
		return models.SprintStatusPlanned
	}
}

// RPC Methods
func (h *SprintHandler) CreateSprint(ctx context.Context, req *pb.CreateSprintRequest) (*pb.CreateSprintResponse, error) {
	sprint, err := h.svc.CreateSprint(ctx, req)
	if err != nil {
		h.log.Sugar().Errorw("Failed to create sprint", "error", err)
		return nil, status.Errorf(codes.Internal, "failed to create sprint: %v", err)
	}
	return &pb.CreateSprintResponse{Sprint: sprintToProto(sprint)}, nil
}

func (h *SprintHandler) GetSprint(ctx context.Context, req *pb.GetSprintRequest) (*pb.GetSprintResponse, error) {
	sprint, err := h.svc.GetSprint(ctx, req.Id)
	if err != nil {
		h.log.Sugar().Errorw("Failed to get sprint", "error", err)
		return nil, status.Errorf(codes.Internal, "failed to get sprint: %v", err)
	}
	if sprint == nil {
		return nil, status.Error(codes.NotFound, "sprint not found")
	}
	return &pb.GetSprintResponse{Sprint: sprintToProto(sprint)}, nil
}

func (h *SprintHandler) ListSprints(ctx context.Context, req *pb.ListSprintsRequest) (*pb.ListSprintsResponse, error) {
	sprints, err := h.svc.ListSprints(ctx, req.ProjectId, statusFromProto(req.Status))
	if err != nil {
		h.log.Sugar().Errorw("Failed to list sprints", "error", err)
		return nil, status.Errorf(codes.Internal, "failed to list sprints: %v", err)
	}

	var pbSprints []*pb.Sprint
	for _, s := range sprints {
		pbSprints = append(pbSprints, sprintToProto(s))
	}
	return &pb.ListSprintsResponse{Sprints: pbSprints}, nil
}

func (h *SprintHandler) UpdateSprint(ctx context.Context, req *pb.UpdateSprintRequest) (*pb.UpdateSprintResponse, error) {
	sprint, err := h.svc.UpdateSprint(ctx, req)
	if err != nil {
		h.log.Sugar().Errorw("Failed to update sprint", "error", err)
		return nil, status.Errorf(codes.Internal, "failed to update sprint: %v", err)
	}
	return &pb.UpdateSprintResponse{Sprint: sprintToProto(sprint)}, nil
}

func (h *SprintHandler) DeleteSprint(ctx context.Context, req *pb.DeleteSprintRequest) (*pb.DeleteSprintResponse, error) {
	if err := h.svc.DeleteSprint(ctx, req.Id); err != nil {
		h.log.Sugar().Errorw("Failed to delete sprint", "error", err)
		return nil, status.Errorf(codes.Internal, "failed to delete sprint: %v", err)
	}
	return &pb.DeleteSprintResponse{}, nil
}

func (h *SprintHandler) AddIssueToSprint(ctx context.Context, req *pb.AddIssueToSprintRequest) (*pb.AddIssueToSprintResponse, error) {
	if err := h.svc.AddIssueToSprint(ctx, req.SprintId, req.IssueId); err != nil {
		h.log.Sugar().Errorw("Failed to add issue to sprint", "error", err)
		return nil, status.Errorf(codes.Internal, "failed to add issue to sprint: %v", err)
	}
	return &pb.AddIssueToSprintResponse{}, nil
}

func (h *SprintHandler) RemoveIssueFromSprint(ctx context.Context, req *pb.RemoveIssueFromSprintRequest) (*pb.RemoveIssueFromSprintResponse, error) {
	if err := h.svc.RemoveIssueFromSprint(ctx, req.SprintId, req.IssueId); err != nil {
		h.log.Sugar().Errorw("Failed to remove issue from sprint", "error", err)
		return nil, status.Errorf(codes.Internal, "failed to remove issue from sprint: %v", err)
	}
	return &pb.RemoveIssueFromSprintResponse{}, nil
}

func (h *SprintHandler) StartSprint(ctx context.Context, req *pb.StartSprintRequest) (*pb.StartSprintResponse, error) {
	sprint, err := h.svc.StartSprint(ctx, req.SprintId)
	if err != nil {
		h.log.Sugar().Errorw("Failed to start sprint", "error", err)
		return nil, status.Errorf(codes.Internal, "failed to start sprint: %v", err)
	}
	return &pb.StartSprintResponse{Sprint: sprintToProto(sprint)}, nil
}

func (h *SprintHandler) CompleteSprint(ctx context.Context, req *pb.CompleteSprintRequest) (*pb.CompleteSprintResponse, error) {
	sprint, err := h.svc.CompleteSprint(ctx, req.SprintId)
	if err != nil {
		h.log.Sugar().Errorw("Failed to complete sprint", "error", err)
		return nil, status.Errorf(codes.Internal, "failed to complete sprint: %v", err)
	}
	return &pb.CompleteSprintResponse{Sprint: sprintToProto(sprint)}, nil
}

func (h *SprintHandler) GetSprintIssues(ctx context.Context, req *pb.GetSprintIssuesRequest) (*pb.GetSprintIssuesResponse, error) {
	issueIDs, err := h.svc.GetSprintIssues(ctx, req.SprintId)
	if err != nil {
		h.log.Sugar().Errorw("Failed to get sprint issues", "error", err)
		return nil, status.Errorf(codes.Internal, "failed to get sprint issues: %v", err)
	}
	return &pb.GetSprintIssuesResponse{IssueIds: issueIDs}, nil
}
