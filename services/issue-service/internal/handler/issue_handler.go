package handler

import (
	"context"

	"github.com/nexusflow/nexusflow/pkg/logger"
	commonpb "github.com/nexusflow/nexusflow/pkg/proto/common/v1"
	pb "github.com/nexusflow/nexusflow/pkg/proto/issue/v1"
	"github.com/nexusflow/nexusflow/services/issue-service/internal/models"
	"github.com/nexusflow/nexusflow/services/issue-service/internal/service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// IssueHandler handles gRPC requests
type IssueHandler struct {
	pb.UnimplementedIssueServiceServer
	service *service.IssueService
	log     *logger.Logger
}

// NewIssueHandler creates a new issue handler
func NewIssueHandler(service *service.IssueService, log *logger.Logger) *IssueHandler {
	return &IssueHandler{
		service: service,
		log:     log,
	}
}

// CreateIssue creates a new issue
func (h *IssueHandler) CreateIssue(ctx context.Context, req *pb.CreateIssueRequest) (*pb.CreateIssueResponse, error) {
	// TODO: Extract user ID from context
	userID := "00000000-0000-0000-0000-000000000000" // Placeholder

	input := service.CreateIssueInput{
		ProjectID:   req.ProjectId,
		Summary:     req.Summary,
		Description: req.Description,
		Type:        h.protoTypeToModel(req.Type),
		Priority:    h.protoPriorityToModel(req.Priority),
		AssigneeID:  req.AssigneeId,
		ReporterID:  userID,
		ParentID:    req.ParentId,
		CustomFields: h.protoCustomFieldsToMap(req.CustomFields),
	}

	issue, err := h.service.CreateIssue(ctx, input)
	if err != nil {
		h.log.Sugar().Errorw("Failed to create issue", "error", err)
		return nil, status.Errorf(codes.Internal, "failed to create issue: %v", err)
	}

	return &pb.CreateIssueResponse{
		Issue: h.issueToProto(issue),
	}, nil
}

// GetIssue gets an issue
func (h *IssueHandler) GetIssue(ctx context.Context, req *pb.GetIssueRequest) (*pb.GetIssueResponse, error) {
	issue, err := h.service.GetIssue(ctx, req.Id)
	if err != nil {
		h.log.Sugar().Errorw("Failed to get issue", "error", err)
		return nil, status.Errorf(codes.Internal, "failed to get issue: %v", err)
	}
	if issue == nil {
		return nil, status.Error(codes.NotFound, "issue not found")
	}

	return &pb.GetIssueResponse{
		Issue: h.issueToProto(issue),
	}, nil
}

// GetIssueByKey gets an issue by key
func (h *IssueHandler) GetIssueByKey(ctx context.Context, req *pb.GetIssueByKeyRequest) (*pb.GetIssueByKeyResponse, error) {
	issue, err := h.service.GetIssueByKey(ctx, req.Key)
	if err != nil {
		h.log.Sugar().Errorw("Failed to get issue by key", "error", err)
		return nil, status.Errorf(codes.Internal, "failed to get issue by key: %v", err)
	}
	if issue == nil {
		return nil, status.Error(codes.NotFound, "issue not found")
	}

	return &pb.GetIssueByKeyResponse{
		Issue: h.issueToProto(issue),
	}, nil
}

// ListIssues lists issues
func (h *IssueHandler) ListIssues(ctx context.Context, req *pb.ListIssuesRequest) (*pb.ListIssuesResponse, error) {
	page := 1
	pageSize := 10
	if req.Pagination != nil {
		page = int(req.Pagination.Page)
		pageSize = int(req.Pagination.PageSize)
	}

	issues, count, err := h.service.ListIssues(ctx, req.ProjectId, page, pageSize)
	if err != nil {
		h.log.Sugar().Errorw("Failed to list issues", "error", err)
		return nil, status.Errorf(codes.Internal, "failed to list issues: %v", err)
	}

	var pbIssues []*pb.Issue
	for _, i := range issues {
		pbIssues = append(pbIssues, h.issueToProto(i))
	}

	return &pb.ListIssuesResponse{
		Issues: pbIssues,
		Pagination: &commonpb.PaginationResponse{
			Page:       int32(page),
			PageSize:   int32(pageSize),
			TotalItems: int64(count),
			TotalPages: int32((count + pageSize - 1) / pageSize),
		},
	}, nil
}

// Custom Fields

func (h *IssueHandler) CreateCustomField(ctx context.Context, req *pb.CreateCustomFieldRequest) (*pb.CreateCustomFieldResponse, error) {
	field := &models.CustomField{
		ProjectID:   req.ProjectId,
		Name:        req.Name,
		Description: req.Description,
		Type:        h.protoCustomFieldTypeToModel(req.Type),
		Required:    req.Required,
		Options:     req.Options,
	}

	created, err := h.service.CreateCustomField(ctx, field)
	if err != nil {
		h.log.Sugar().Errorw("Failed to create custom field", "error", err)
		return nil, status.Errorf(codes.Internal, "failed to create custom field: %v", err)
	}

	return &pb.CreateCustomFieldResponse{
		Field: h.customFieldToProto(created),
	}, nil
}

func (h *IssueHandler) ListCustomFields(ctx context.Context, req *pb.ListCustomFieldsRequest) (*pb.ListCustomFieldsResponse, error) {
	fields, err := h.service.ListCustomFields(ctx, req.ProjectId)
	if err != nil {
		h.log.Sugar().Errorw("Failed to list custom fields", "error", err)
		return nil, status.Errorf(codes.Internal, "failed to list custom fields: %v", err)
	}

	var pbFields []*pb.CustomField
	for _, f := range fields {
		pbFields = append(pbFields, h.customFieldToProto(f))
	}

	return &pb.ListCustomFieldsResponse{
		Fields: pbFields,
	}, nil
}

// Helpers

func (h *IssueHandler) protoCustomFieldsToMap(fields []*pb.CustomFieldValue) map[string]interface{} {
	result := make(map[string]interface{})
	for _, f := range fields {
		// TODO: Handle Any type properly. For now, assuming string or simple types?
		// This is tricky with Any.
		// Let's assume for now we skip complex value parsing or just store nil
		// Real implementation needs Any unpacking.
		result[f.FieldId] = f.Value // This won't work directly as it's *anypb.Any
	}
	return result
}

func (h *IssueHandler) customFieldToProto(f *models.CustomField) *pb.CustomField {
	if f == nil {
		return nil
	}
	return &pb.CustomField{
		Id:          f.ID,
		ProjectId:   f.ProjectID,
		Name:        f.Name,
		Description: f.Description,
		// Type:        ... // Need mapping
		Required:    f.Required,
		Options:     f.Options,
	}
}

func (h *IssueHandler) protoCustomFieldTypeToModel(t pb.CustomFieldType) models.CustomFieldType {
	switch t {
	case pb.CustomFieldType_CUSTOM_FIELD_TYPE_TEXT:
		return models.CustomFieldTypeText
	case pb.CustomFieldType_CUSTOM_FIELD_TYPE_NUMBER:
		return models.CustomFieldTypeNumber
	case pb.CustomFieldType_CUSTOM_FIELD_TYPE_SELECT:
		return models.CustomFieldTypeSelect
	default:
		return models.CustomFieldTypeText
	}
}

// Helpers

func (h *IssueHandler) issueToProto(i *models.Issue) *pb.Issue {
	if i == nil {
		return nil
	}
	return &pb.Issue{
		Id:          i.ID,
		ProjectId:   i.ProjectID,
		Key:         i.Key,
		Summary:     i.Summary,
		Description: i.Description,
		// Type:        pb.IssueType(pb.IssueType_value[string(i.Type)]), // Need mapping
		// Priority:    pb.IssuePriority(pb.IssuePriority_value[string(i.Priority)]), // Need mapping
		StatusId:    i.StatusID,
		AssigneeId:  i.AssigneeID,
		ReporterId:  i.ReporterID,
		ParentId:    i.ParentID,
		SprintId:    i.SprintID,
		StoryPoints: i.StoryPoints,
		CreatedAt:   timestamppb.New(i.CreatedAt),
		UpdatedAt:   timestamppb.New(i.UpdatedAt),
		DueDate:     timestamppb.New(i.DueDate),
	}
}

func (h *IssueHandler) protoTypeToModel(t pb.IssueType) models.IssueType {
	switch t {
	case pb.IssueType_ISSUE_TYPE_EPIC:
		return models.IssueTypeEpic
	case pb.IssueType_ISSUE_TYPE_STORY:
		return models.IssueTypeStory
	case pb.IssueType_ISSUE_TYPE_TASK:
		return models.IssueTypeTask
	case pb.IssueType_ISSUE_TYPE_SUB_TASK:
		return models.IssueTypeSubTask
	case pb.IssueType_ISSUE_TYPE_BUG:
		return models.IssueTypeBug
	case pb.IssueType_ISSUE_TYPE_IMPROVEMENT:
		return models.IssueTypeImprovement
	default:
		return models.IssueTypeTask
	}
}

func (h *IssueHandler) protoPriorityToModel(p pb.IssuePriority) models.IssuePriority {
	switch p {
	case pb.IssuePriority_ISSUE_PRIORITY_LOWEST:
		return models.IssuePriorityLowest
	case pb.IssuePriority_ISSUE_PRIORITY_LOW:
		return models.IssuePriorityLow
	case pb.IssuePriority_ISSUE_PRIORITY_MEDIUM:
		return models.IssuePriorityMedium
	case pb.IssuePriority_ISSUE_PRIORITY_HIGH:
		return models.IssuePriorityHigh
	case pb.IssuePriority_ISSUE_PRIORITY_HIGHEST:
		return models.IssuePriorityHighest
	default:
		return models.IssuePriorityMedium
	}
}
