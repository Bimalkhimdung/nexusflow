package handler

import (
	"context"

	"github.com/nexusflow/nexusflow/pkg/logger"
	commonpb "github.com/nexusflow/nexusflow/pkg/proto/common/v1"
	pb "github.com/nexusflow/nexusflow/pkg/proto/project/v1"
	"github.com/nexusflow/nexusflow/services/project-service/internal/models"
	"github.com/nexusflow/nexusflow/services/project-service/internal/service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// ProjectHandler handles gRPC requests
type ProjectHandler struct {
	pb.UnimplementedProjectServiceServer
	service *service.ProjectService
	log     *logger.Logger
}

// NewProjectHandler creates a new project handler
func NewProjectHandler(service *service.ProjectService, log *logger.Logger) *ProjectHandler {
	return &ProjectHandler{
		service: service,
		log:     log,
	}
}

// CreateProject creates a new project
func (h *ProjectHandler) CreateProject(ctx context.Context, req *pb.CreateProjectRequest) (*pb.CreateProjectResponse, error) {
	// TODO: Extract user ID from context
	userID := "00000000-0000-0000-0000-000000000000" // Placeholder

	inputType := h.protoTypeToModel(req.Type)

	input := service.CreateProjectInput{
		OrganizationID: req.OrganizationId,
		Key:            req.Key,
		Name:           req.Name,
		Description:    req.Description,
		Type:           inputType,
		LeadID:         req.LeadId,
		UserID:         userID,
	}

	project, err := h.service.CreateProject(ctx, input)
	if err != nil {
		h.log.Sugar().Errorw("Failed to create project", "error", err)
		return nil, status.Errorf(codes.Internal, "failed to create project: %v", err)
	}

	return &pb.CreateProjectResponse{
		Project: h.projectToProto(project),
	}, nil
}

// GetProject gets a project
func (h *ProjectHandler) GetProject(ctx context.Context, req *pb.GetProjectRequest) (*pb.GetProjectResponse, error) {
	project, err := h.service.GetProject(ctx, req.Id)
	if err != nil {
		h.log.Sugar().Errorw("Failed to get project", "error", err)
		return nil, status.Errorf(codes.Internal, "failed to get project: %v", err)
	}
	if project == nil {
		return nil, status.Error(codes.NotFound, "project not found")
	}

	return &pb.GetProjectResponse{
		Project: h.projectToProto(project),
	}, nil
}

// GetProjectByKey gets a project by key
func (h *ProjectHandler) GetProjectByKey(ctx context.Context, req *pb.GetProjectByKeyRequest) (*pb.GetProjectByKeyResponse, error) {
	project, err := h.service.GetProjectByKey(ctx, req.OrganizationId, req.Key)
	if err != nil {
		h.log.Sugar().Errorw("Failed to get project by key", "error", err)
		return nil, status.Errorf(codes.Internal, "failed to get project by key: %v", err)
	}
	if project == nil {
		return nil, status.Error(codes.NotFound, "project not found")
	}

	return &pb.GetProjectByKeyResponse{
		Project: h.projectToProto(project),
	}, nil
}

// UpdateProject updates a project
func (h *ProjectHandler) UpdateProject(ctx context.Context, req *pb.UpdateProjectRequest) (*pb.UpdateProjectResponse, error) {
	input := service.UpdateProjectInput{
		ID:       req.Id,
		Settings: req.Settings,
	}
	if req.Name != nil {
		input.Name = req.Name
	}
	if req.Description != nil {
		input.Description = req.Description
	}
	if req.AvatarUrl != nil {
		input.AvatarURL = req.AvatarUrl
	}
	if req.LeadId != nil {
		input.LeadID = req.LeadId
	}

	project, err := h.service.UpdateProject(ctx, input)
	if err != nil {
		h.log.Sugar().Errorw("Failed to update project", "error", err)
		return nil, status.Errorf(codes.Internal, "failed to update project: %v", err)
	}

	return &pb.UpdateProjectResponse{
		Project: h.projectToProto(project),
	}, nil
}

// DeleteProject deletes a project
func (h *ProjectHandler) DeleteProject(ctx context.Context, req *pb.DeleteProjectRequest) (*pb.DeleteProjectResponse, error) {
	if err := h.service.DeleteProject(ctx, req.Id); err != nil {
		h.log.Sugar().Errorw("Failed to delete project", "error", err)
		return nil, status.Errorf(codes.Internal, "failed to delete project: %v", err)
	}

	return &pb.DeleteProjectResponse{
		Response: &commonpb.SuccessResponse{
			Success: true,
		},
	}, nil
}

// ListProjects lists projects
func (h *ProjectHandler) ListProjects(ctx context.Context, req *pb.ListProjectsRequest) (*pb.ListProjectsResponse, error) {
	// TODO: Extract user ID from context for proper filtering
	userID := req.UserId // Use from request for now
	if userID == "" {
		userID = "00000000-0000-0000-0000-000000000000" // Placeholder
	}

	page := 1
	pageSize := 10
	if req.Pagination != nil {
		page = int(req.Pagination.Page)
		pageSize = int(req.Pagination.PageSize)
	}

	projects, count, err := h.service.ListProjects(ctx, req.OrganizationId, userID, page, pageSize)
	if err != nil {
		h.log.Sugar().Errorw("Failed to list projects", "error", err)
		return nil, status.Errorf(codes.Internal, "failed to list projects: %v", err)
	}

	var pbProjects []*pb.Project
	for _, p := range projects {
		pbProjects = append(pbProjects, h.projectToProto(p))
	}

	return &pb.ListProjectsResponse{
		Projects: pbProjects,
		Pagination: &commonpb.PaginationResponse{
			Page:       int32(page),
			PageSize:   int32(pageSize),
			TotalItems: int64(count),
			TotalPages: int32((count + pageSize - 1) / pageSize),
		},
	}, nil
}

// ArchiveProject archives a project
func (h *ProjectHandler) ArchiveProject(ctx context.Context, req *pb.ArchiveProjectRequest) (*pb.ArchiveProjectResponse, error) {
	project, err := h.service.ArchiveProject(ctx, req.Id)
	if err != nil {
		h.log.Sugar().Errorw("Failed to archive project", "error", err)
		return nil, status.Errorf(codes.Internal, "failed to archive project: %v", err)
	}

	return &pb.ArchiveProjectResponse{
		Project: h.projectToProto(project),
	}, nil
}

// AddProjectMember adds a member to a project
func (h *ProjectHandler) AddProjectMember(ctx context.Context, req *pb.AddProjectMemberRequest) (*pb.AddProjectMemberResponse, error) {
	role := h.protoRoleToModel(req.Role)

	member, err := h.service.AddMember(ctx, req.ProjectId, req.UserId, role)
	if err != nil {
		h.log.Sugar().Errorw("Failed to add member", "error", err)
		return nil, status.Errorf(codes.Internal, "failed to add member: %v", err)
	}

	return &pb.AddProjectMemberResponse{
		Member: h.memberToProto(member),
	}, nil
}

// RemoveProjectMember removes a member from a project
func (h *ProjectHandler) RemoveProjectMember(ctx context.Context, req *pb.RemoveProjectMemberRequest) (*pb.RemoveProjectMemberResponse, error) {
	if err := h.service.RemoveMember(ctx, req.ProjectId, req.UserId); err != nil {
		h.log.Sugar().Errorw("Failed to remove member", "error", err)
		return nil, status.Errorf(codes.Internal, "failed to remove member: %v", err)
	}

	return &pb.RemoveProjectMemberResponse{
		Response: &commonpb.SuccessResponse{
			Success: true,
		},
	}, nil
}

// UpdateProjectMemberRole updates a member's role
func (h *ProjectHandler) UpdateProjectMemberRole(ctx context.Context, req *pb.UpdateProjectMemberRoleRequest) (*pb.UpdateProjectMemberRoleResponse, error) {
	role := h.protoRoleToModel(req.Role)

	member, err := h.service.UpdateMemberRole(ctx, req.ProjectId, req.UserId, role)
	if err != nil {
		h.log.Sugar().Errorw("Failed to update member role", "error", err)
		return nil, status.Errorf(codes.Internal, "failed to update member role: %v", err)
	}

	return &pb.UpdateProjectMemberRoleResponse{
		Member: h.memberToProto(member),
	}, nil
}

// ListProjectMembers lists members of a project
func (h *ProjectHandler) ListProjectMembers(ctx context.Context, req *pb.ListProjectMembersRequest) (*pb.ListProjectMembersResponse, error) {
	page := 1
	pageSize := 10
	if req.Pagination != nil {
		page = int(req.Pagination.Page)
		pageSize = int(req.Pagination.PageSize)
	}

	members, count, err := h.service.ListMembers(ctx, req.ProjectId, page, pageSize)
	if err != nil {
		h.log.Sugar().Errorw("Failed to list members", "error", err)
		return nil, status.Errorf(codes.Internal, "failed to list members: %v", err)
	}

	var pbMembers []*pb.ProjectMember
	for _, m := range members {
		pbMembers = append(pbMembers, h.memberToProto(m))
	}

	return &pb.ListProjectMembersResponse{
		Members: pbMembers,
		Pagination: &commonpb.PaginationResponse{
			Page:       int32(page),
			PageSize:   int32(pageSize),
			TotalItems: int64(count),
			TotalPages: int32((count + pageSize - 1) / pageSize),
		},
	}, nil
}

// Helpers

func (h *ProjectHandler) projectToProto(p *models.Project) *pb.Project {
	if p == nil {
		return nil
	}
	return &pb.Project{
		Id:             p.ID,
		OrganizationId: p.OrganizationID,
		Key:            p.Key,
		Name:           p.Name,
		Description:    p.Description,
		AvatarUrl:      p.AvatarURL,
		// Type:           pb.ProjectType(pb.ProjectType_value[string(p.Type)]), // Need mapping
		// Status:         pb.ProjectStatus(pb.ProjectStatus_value[string(p.Status)]), // Need mapping
		LeadId:         p.LeadID,
		Settings:       p.Settings,
		CreatedAt:      timestamppb.New(p.CreatedAt),
		UpdatedAt:      timestamppb.New(p.UpdatedAt),
	}
}

func (h *ProjectHandler) memberToProto(m *models.ProjectMember) *pb.ProjectMember {
	if m == nil {
		return nil
	}
	return &pb.ProjectMember{
		Id:        m.ID,
		ProjectId: m.ProjectID,
		UserId:    m.UserID,
		// Role:      pb.ProjectRole(pb.ProjectRole_value[string(m.Role)]), // Need mapping
		JoinedAt:  timestamppb.New(m.JoinedAt),
	}
}

func (h *ProjectHandler) protoTypeToModel(t pb.ProjectType) models.ProjectType {
	switch t {
	case pb.ProjectType_PROJECT_TYPE_KANBAN:
		return models.ProjectTypeKanban
	case pb.ProjectType_PROJECT_TYPE_SCRUM:
		return models.ProjectTypeScrum
	case pb.ProjectType_PROJECT_TYPE_BUG_TRACKING:
		return models.ProjectTypeBugTracking
	default:
		return models.ProjectTypeKanban
	}
}

func (h *ProjectHandler) protoRoleToModel(r pb.ProjectRole) models.ProjectRole {
	switch r {
	case pb.ProjectRole_PROJECT_ROLE_ADMIN:
		return models.ProjectRoleAdmin
	case pb.ProjectRole_PROJECT_ROLE_MEMBER:
		return models.ProjectRoleMember
	case pb.ProjectRole_PROJECT_ROLE_VIEWER:
		return models.ProjectRoleViewer
	default:
		return models.ProjectRoleMember
	}
}
