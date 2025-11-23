package handler

import (
	"context"

	"github.com/nexusflow/nexusflow/pkg/logger"
	pb "github.com/nexusflow/nexusflow/pkg/proto/workflow/v1"
	"github.com/nexusflow/nexusflow/services/workflow-service/internal/models"
	"github.com/nexusflow/nexusflow/services/workflow-service/internal/service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// WorkflowHandler handles gRPC requests
type WorkflowHandler struct {
	pb.UnimplementedWorkflowServiceServer
	service *service.WorkflowService
	log     *logger.Logger
}

// NewWorkflowHandler creates a new workflow handler
func NewWorkflowHandler(service *service.WorkflowService, log *logger.Logger) *WorkflowHandler {
	return &WorkflowHandler{
		service: service,
		log:     log,
	}
}

// CreateWorkflow creates a new workflow
func (h *WorkflowHandler) CreateWorkflow(ctx context.Context, req *pb.CreateWorkflowRequest) (*pb.CreateWorkflowResponse, error) {
	workflow := &models.Workflow{
		ProjectID:   req.ProjectId,
		Name:        req.Name,
		Description: req.Description,
	}

	created, err := h.service.CreateWorkflow(ctx, workflow)
	if err != nil {
		h.log.Sugar().Errorw("Failed to create workflow", "error", err)
		return nil, status.Errorf(codes.Internal, "failed to create workflow: %v", err)
	}

	return &pb.CreateWorkflowResponse{
		Workflow: h.workflowToProto(created),
	}, nil
}

// ListWorkflows lists workflows
func (h *WorkflowHandler) ListWorkflows(ctx context.Context, req *pb.ListWorkflowsRequest) (*pb.ListWorkflowsResponse, error) {
	workflows, err := h.service.ListWorkflows(ctx, req.ProjectId)
	if err != nil {
		h.log.Sugar().Errorw("Failed to list workflows", "error", err)
		return nil, status.Errorf(codes.Internal, "failed to list workflows: %v", err)
	}

	var pbWorkflows []*pb.Workflow
	for _, w := range workflows {
		pbWorkflows = append(pbWorkflows, h.workflowToProto(w))
	}

	return &pb.ListWorkflowsResponse{
		Workflows: pbWorkflows,
	}, nil
}

// CreateStatus creates a new status
func (h *WorkflowHandler) CreateStatus(ctx context.Context, req *pb.CreateStatusRequest) (*pb.CreateStatusResponse, error) {
	statusModel := &models.WorkflowStatus{
		WorkflowID:  req.WorkflowId,
		Name:        req.Name,
		Description: req.Description,
		Category:    h.protoCategoryToModel(req.Category),
		Color:       req.Color,
	}

	created, err := h.service.CreateStatus(ctx, statusModel)
	if err != nil {
		h.log.Sugar().Errorw("Failed to create status", "error", err)
		return nil, status.Errorf(codes.Internal, "failed to create status: %v", err)
	}

	return &pb.CreateStatusResponse{
		Status: h.statusToProto(created),
	}, nil
}

// CreateTransition creates a new transition
func (h *WorkflowHandler) CreateTransition(ctx context.Context, req *pb.CreateTransitionRequest) (*pb.CreateTransitionResponse, error) {
	transition := &models.WorkflowTransition{
		WorkflowID:   req.WorkflowId,
		Name:         req.Name,
		FromStatusID: req.FromStatusId,
		ToStatusID:   req.ToStatusId,
	}

	created, err := h.service.CreateTransition(ctx, transition)
	if err != nil {
		h.log.Sugar().Errorw("Failed to create transition", "error", err)
		return nil, status.Errorf(codes.Internal, "failed to create transition: %v", err)
	}

	return &pb.CreateTransitionResponse{
		Transition: h.transitionToProto(created),
	}, nil
}

// ExecuteTransition executes a transition
func (h *WorkflowHandler) ExecuteTransition(ctx context.Context, req *pb.ExecuteTransitionRequest) (*pb.ExecuteTransitionResponse, error) {
	err := h.service.ExecuteTransition(ctx, req.IssueId, req.TransitionId, req.UserId)
	if err != nil {
		h.log.Sugar().Errorw("Failed to execute transition", "error", err)
		return nil, status.Errorf(codes.Internal, "failed to execute transition: %v", err)
	}

	return &pb.ExecuteTransitionResponse{}, nil
}

// GetAvailableTransitions gets available transitions
func (h *WorkflowHandler) GetAvailableTransitions(ctx context.Context, req *pb.GetAvailableTransitionsRequest) (*pb.GetAvailableTransitionsResponse, error) {
	transitions, err := h.service.GetAvailableTransitions(ctx, req.IssueId)
	if err != nil {
		h.log.Sugar().Errorw("Failed to get available transitions", "error", err)
		return nil, status.Errorf(codes.Internal, "failed to get available transitions: %v", err)
	}

	var pbTransitions []*pb.Transition
	for _, t := range transitions {
		pbTransitions = append(pbTransitions, h.transitionToProto(t))
	}

	return &pb.GetAvailableTransitionsResponse{
		Transitions: pbTransitions,
	}, nil
}

// Helpers

func (h *WorkflowHandler) workflowToProto(w *models.Workflow) *pb.Workflow {
	if w == nil {
		return nil
	}
	return &pb.Workflow{
		Id:          w.ID,
		ProjectId:   w.ProjectID,
		Name:        w.Name,
		Description: w.Description,
		IsDefault:   w.IsDefault,
	}
}

func (h *WorkflowHandler) statusToProto(s *models.WorkflowStatus) *pb.Status {
	if s == nil {
		return nil
	}
	return &pb.Status{
		Id:          s.ID,
		Name:        s.Name,
		Description: s.Description,
		// Category:    ... // Need mapping
		Color:    s.Color,
		Position: s.Position,
	}
}

func (h *WorkflowHandler) transitionToProto(t *models.WorkflowTransition) *pb.Transition {
	if t == nil {
		return nil
	}
	return &pb.Transition{
		Id:           t.ID,
		Name:         t.Name,
		FromStatusId: t.FromStatusID,
		ToStatusId:   t.ToStatusID,
	}
}

func (h *WorkflowHandler) protoCategoryToModel(c pb.StatusCategory) models.StatusCategory {
	switch c {
	case pb.StatusCategory_STATUS_CATEGORY_TODO:
		return models.StatusCategoryTodo
	case pb.StatusCategory_STATUS_CATEGORY_IN_PROGRESS:
		return models.StatusCategoryInProgress
	case pb.StatusCategory_STATUS_CATEGORY_DONE:
		return models.StatusCategoryDone
	default:
		return models.StatusCategoryTodo
	}
}
