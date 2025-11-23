package handler

import (
	"context"

	"github.com/nexusflow/nexusflow/pkg/logger"
	commonpb "github.com/nexusflow/nexusflow/pkg/proto/common/v1"
	pb "github.com/nexusflow/nexusflow/pkg/proto/org/v1"
	"github.com/nexusflow/nexusflow/services/org-service/internal/models"
	"github.com/nexusflow/nexusflow/services/org-service/internal/service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// OrgHandler handles gRPC requests
type OrgHandler struct {
	pb.UnimplementedOrgServiceServer
	service *service.OrgService
	log     *logger.Logger
}

// NewOrgHandler creates a new organization handler
func NewOrgHandler(service *service.OrgService, log *logger.Logger) *OrgHandler {
	return &OrgHandler{
		service: service,
		log:     log,
	}
}

// CreateOrganization creates a new organization
func (h *OrgHandler) CreateOrganization(ctx context.Context, req *pb.CreateOrganizationRequest) (*pb.CreateOrganizationResponse, error) {
	// TODO: Extract user ID from context
	userID := "00000000-0000-0000-0000-000000000000" // Placeholder

	input := service.CreateOrgInput{
		Name:        req.Name,
		Slug:        req.Slug,
		Description: req.Description,
		UserID:      userID,
	}

	org, err := h.service.CreateOrganization(ctx, input)
	if err != nil {
		h.log.Sugar().Errorw("Failed to create organization", "error", err)
		return nil, status.Errorf(codes.Internal, "failed to create organization: %v", err)
	}

	return &pb.CreateOrganizationResponse{
		Organization: h.orgToProto(org),
	}, nil
}

// GetOrganization gets an organization
func (h *OrgHandler) GetOrganization(ctx context.Context, req *pb.GetOrganizationRequest) (*pb.GetOrganizationResponse, error) {
	org, err := h.service.GetOrganization(ctx, req.Id)
	if err != nil {
		h.log.Sugar().Errorw("Failed to get organization", "error", err)
		return nil, status.Errorf(codes.Internal, "failed to get organization: %v", err)
	}
	if org == nil {
		return nil, status.Error(codes.NotFound, "organization not found")
	}

	return &pb.GetOrganizationResponse{
		Organization: h.orgToProto(org),
	}, nil
}

// UpdateOrganization updates an organization
func (h *OrgHandler) UpdateOrganization(ctx context.Context, req *pb.UpdateOrganizationRequest) (*pb.UpdateOrganizationResponse, error) {
	input := service.UpdateOrgInput{
		ID:       req.Id,
		Settings: req.Settings,
	}
	if req.Name != nil {
		input.Name = req.Name
	}
	if req.Description != nil {
		input.Description = req.Description
	}
	if req.LogoUrl != nil {
		input.LogoURL = req.LogoUrl
	}

	org, err := h.service.UpdateOrganization(ctx, input)
	if err != nil {
		h.log.Sugar().Errorw("Failed to update organization", "error", err)
		return nil, status.Errorf(codes.Internal, "failed to update organization: %v", err)
	}

	return &pb.UpdateOrganizationResponse{
		Organization: h.orgToProto(org),
	}, nil
}

// DeleteOrganization deletes an organization
func (h *OrgHandler) DeleteOrganization(ctx context.Context, req *pb.DeleteOrganizationRequest) (*pb.DeleteOrganizationResponse, error) {
	if err := h.service.DeleteOrganization(ctx, req.Id); err != nil {
		h.log.Sugar().Errorw("Failed to delete organization", "error", err)
		return nil, status.Errorf(codes.Internal, "failed to delete organization: %v", err)
	}

	return &pb.DeleteOrganizationResponse{
		Response: &commonpb.SuccessResponse{
			Success: true,
		},
	}, nil
}

// ListOrganizations lists organizations
func (h *OrgHandler) ListOrganizations(ctx context.Context, req *pb.ListOrganizationsRequest) (*pb.ListOrganizationsResponse, error) {
	// TODO: Use user ID from request or context
	userID := req.UserId
	if userID == "" {
		userID = "00000000-0000-0000-0000-000000000000" // Placeholder
	}

	page := 1
	pageSize := 10
	if req.Pagination != nil {
		page = int(req.Pagination.Page)
		pageSize = int(req.Pagination.PageSize)
	}

	orgs, count, err := h.service.ListOrganizations(ctx, userID, page, pageSize)
	if err != nil {
		h.log.Sugar().Errorw("Failed to list organizations", "error", err)
		return nil, status.Errorf(codes.Internal, "failed to list organizations: %v", err)
	}

	var pbOrgs []*pb.Organization
	for _, org := range orgs {
		pbOrgs = append(pbOrgs, h.orgToProto(org))
	}

	return &pb.ListOrganizationsResponse{
		Organizations: pbOrgs,
		Pagination: &commonpb.PaginationResponse{
			Page:       int32(page),
			PageSize:   int32(pageSize),
			TotalItems: int64(count),
			TotalPages: int32((count + pageSize - 1) / pageSize),
		},
	}, nil
}

// AddMember adds a member
func (h *OrgHandler) AddMember(ctx context.Context, req *pb.AddMemberRequest) (*pb.AddMemberResponse, error) {
	member, err := h.service.AddMember(ctx, req.OrganizationId, req.UserId, models.OrgRole(req.Role.String())) // Enum conversion might need adjustment
	// Actually, protobuf enums are int32, but we store string in DB. Need mapping.
	// For now, let's assume simple string conversion works or fix it.
	// The proto enum names are like ORG_ROLE_OWNER, but DB expects "owner".
	// We need a helper for this.
	
	role := h.protoRoleToModel(req.Role)

	member, err = h.service.AddMember(ctx, req.OrganizationId, req.UserId, role)
	if err != nil {
		h.log.Sugar().Errorw("Failed to add member", "error", err)
		return nil, status.Errorf(codes.Internal, "failed to add member: %v", err)
	}

	return &pb.AddMemberResponse{
		Member: h.memberToProto(member),
	}, nil
}

// RemoveMember removes a member
func (h *OrgHandler) RemoveMember(ctx context.Context, req *pb.RemoveMemberRequest) (*pb.RemoveMemberResponse, error) {
	if err := h.service.RemoveMember(ctx, req.OrganizationId, req.UserId); err != nil {
		h.log.Sugar().Errorw("Failed to remove member", "error", err)
		return nil, status.Errorf(codes.Internal, "failed to remove member: %v", err)
	}

	return &pb.RemoveMemberResponse{
		Response: &commonpb.SuccessResponse{
			Success: true,
		},
	}, nil
}

// UpdateMemberRole updates a member role
func (h *OrgHandler) UpdateMemberRole(ctx context.Context, req *pb.UpdateMemberRoleRequest) (*pb.UpdateMemberRoleResponse, error) {
	role := h.protoRoleToModel(req.Role)
	
	member, err := h.service.UpdateMemberRole(ctx, req.OrganizationId, req.UserId, role)
	if err != nil {
		h.log.Sugar().Errorw("Failed to update member role", "error", err)
		return nil, status.Errorf(codes.Internal, "failed to update member role: %v", err)
	}

	return &pb.UpdateMemberRoleResponse{
		Member: h.memberToProto(member),
	}, nil
}

// ListMembers lists members
func (h *OrgHandler) ListMembers(ctx context.Context, req *pb.ListMembersRequest) (*pb.ListMembersResponse, error) {
	page := 1
	pageSize := 10
	if req.Pagination != nil {
		page = int(req.Pagination.Page)
		pageSize = int(req.Pagination.PageSize)
	}

	members, count, err := h.service.ListMembers(ctx, req.OrganizationId, page, pageSize)
	if err != nil {
		h.log.Sugar().Errorw("Failed to list members", "error", err)
		return nil, status.Errorf(codes.Internal, "failed to list members: %v", err)
	}

	var pbMembers []*pb.OrgMember
	for _, m := range members {
		pbMembers = append(pbMembers, h.memberToProto(m))
	}

	return &pb.ListMembersResponse{
		Members: pbMembers,
		Pagination: &commonpb.PaginationResponse{
			Page:       int32(page),
			PageSize:   int32(pageSize),
			TotalItems: int64(count),
			TotalPages: int32((count + pageSize - 1) / pageSize),
		},
	}, nil
}

// CreateTeam creates a team
func (h *OrgHandler) CreateTeam(ctx context.Context, req *pb.CreateTeamRequest) (*pb.CreateTeamResponse, error) {
	team, err := h.service.CreateTeam(ctx, req.OrganizationId, req.Name, req.Description)
	if err != nil {
		h.log.Sugar().Errorw("Failed to create team", "error", err)
		return nil, status.Errorf(codes.Internal, "failed to create team: %v", err)
	}

	return &pb.CreateTeamResponse{
		Team: h.teamToProto(team),
	}, nil
}

// GetTeam gets a team
func (h *OrgHandler) GetTeam(ctx context.Context, req *pb.GetTeamRequest) (*pb.GetTeamResponse, error) {
	team, err := h.service.GetTeam(ctx, req.Id)
	if err != nil {
		h.log.Sugar().Errorw("Failed to get team", "error", err)
		return nil, status.Errorf(codes.Internal, "failed to get team: %v", err)
	}
	if team == nil {
		return nil, status.Error(codes.NotFound, "team not found")
	}

	return &pb.GetTeamResponse{
		Team: h.teamToProto(team),
	}, nil
}

// UpdateTeam updates a team
func (h *OrgHandler) UpdateTeam(ctx context.Context, req *pb.UpdateTeamRequest) (*pb.UpdateTeamResponse, error) {
	team, err := h.service.UpdateTeam(ctx, req.Id, req.Name, req.Description)
	if err != nil {
		h.log.Sugar().Errorw("Failed to update team", "error", err)
		return nil, status.Errorf(codes.Internal, "failed to update team: %v", err)
	}

	return &pb.UpdateTeamResponse{
		Team: h.teamToProto(team),
	}, nil
}

// DeleteTeam deletes a team
func (h *OrgHandler) DeleteTeam(ctx context.Context, req *pb.DeleteTeamRequest) (*pb.DeleteTeamResponse, error) {
	if err := h.service.DeleteTeam(ctx, req.Id); err != nil {
		h.log.Sugar().Errorw("Failed to delete team", "error", err)
		return nil, status.Errorf(codes.Internal, "failed to delete team: %v", err)
	}

	return &pb.DeleteTeamResponse{
		Response: &commonpb.SuccessResponse{
			Success: true,
		},
	}, nil
}

// ListTeams lists teams
func (h *OrgHandler) ListTeams(ctx context.Context, req *pb.ListTeamsRequest) (*pb.ListTeamsResponse, error) {
	page := 1
	pageSize := 10
	if req.Pagination != nil {
		page = int(req.Pagination.Page)
		pageSize = int(req.Pagination.PageSize)
	}

	teams, count, err := h.service.ListTeams(ctx, req.OrganizationId, page, pageSize)
	if err != nil {
		h.log.Sugar().Errorw("Failed to list teams", "error", err)
		return nil, status.Errorf(codes.Internal, "failed to list teams: %v", err)
	}

	var pbTeams []*pb.Team
	for _, t := range teams {
		pbTeams = append(pbTeams, h.teamToProto(t))
	}

	return &pb.ListTeamsResponse{
		Teams: pbTeams,
		Pagination: &commonpb.PaginationResponse{
			Page:       int32(page),
			PageSize:   int32(pageSize),
			TotalItems: int64(count),
			TotalPages: int32((count + pageSize - 1) / pageSize),
		},
	}, nil
}

// AddTeamMember adds a team member
func (h *OrgHandler) AddTeamMember(ctx context.Context, req *pb.AddTeamMemberRequest) (*pb.AddTeamMemberResponse, error) {
	if err := h.service.AddTeamMember(ctx, req.TeamId, req.UserId); err != nil {
		h.log.Sugar().Errorw("Failed to add team member", "error", err)
		return nil, status.Errorf(codes.Internal, "failed to add team member: %v", err)
	}
	
	// Return updated team
	team, _ := h.service.GetTeam(ctx, req.TeamId)
	return &pb.AddTeamMemberResponse{
		Team: h.teamToProto(team),
	}, nil
}

// RemoveTeamMember removes a team member
func (h *OrgHandler) RemoveTeamMember(ctx context.Context, req *pb.RemoveTeamMemberRequest) (*pb.RemoveTeamMemberResponse, error) {
	if err := h.service.RemoveTeamMember(ctx, req.TeamId, req.UserId); err != nil {
		h.log.Sugar().Errorw("Failed to remove team member", "error", err)
		return nil, status.Errorf(codes.Internal, "failed to remove team member: %v", err)
	}

	// Return updated team
	team, _ := h.service.GetTeam(ctx, req.TeamId)
	return &pb.RemoveTeamMemberResponse{
		Team: h.teamToProto(team),
	}, nil
}

// CreateInvite creates an invite
func (h *OrgHandler) CreateInvite(ctx context.Context, req *pb.CreateInviteRequest) (*pb.CreateInviteResponse, error) {
	// TODO: Extract user ID from context
	invitedBy := "00000000-0000-0000-0000-000000000000"

	role := h.protoRoleToModel(req.Role)

	invite, err := h.service.CreateInvite(ctx, req.OrganizationId, req.Email, role, invitedBy)
	if err != nil {
		h.log.Sugar().Errorw("Failed to create invite", "error", err)
		return nil, status.Errorf(codes.Internal, "failed to create invite: %v", err)
	}

	return &pb.CreateInviteResponse{
		Invite: h.inviteToProto(invite),
	}, nil
}

// AcceptInvite accepts an invite
func (h *OrgHandler) AcceptInvite(ctx context.Context, req *pb.AcceptInviteRequest) (*pb.AcceptInviteResponse, error) {
	member, err := h.service.AcceptInvite(ctx, req.Token, req.UserId)
	if err != nil {
		h.log.Sugar().Errorw("Failed to accept invite", "error", err)
		return nil, status.Errorf(codes.Internal, "failed to accept invite: %v", err)
	}

	return &pb.AcceptInviteResponse{
		Member: h.memberToProto(member),
	}, nil
}

// RevokeInvite revokes an invite
func (h *OrgHandler) RevokeInvite(ctx context.Context, req *pb.RevokeInviteRequest) (*pb.RevokeInviteResponse, error) {
	if err := h.service.RevokeInvite(ctx, req.Id); err != nil {
		h.log.Sugar().Errorw("Failed to revoke invite", "error", err)
		return nil, status.Errorf(codes.Internal, "failed to revoke invite: %v", err)
	}

	return &pb.RevokeInviteResponse{
		Response: &commonpb.SuccessResponse{
			Success: true,
		},
	}, nil
}

// ListInvites lists invites
func (h *OrgHandler) ListInvites(ctx context.Context, req *pb.ListInvitesRequest) (*pb.ListInvitesResponse, error) {
	page := 1
	pageSize := 10
	if req.Pagination != nil {
		page = int(req.Pagination.Page)
		pageSize = int(req.Pagination.PageSize)
	}

	invites, count, err := h.service.ListInvites(ctx, req.OrganizationId, page, pageSize)
	if err != nil {
		h.log.Sugar().Errorw("Failed to list invites", "error", err)
		return nil, status.Errorf(codes.Internal, "failed to list invites: %v", err)
	}

	var pbInvites []*pb.Invite
	for _, i := range invites {
		pbInvites = append(pbInvites, h.inviteToProto(i))
	}

	return &pb.ListInvitesResponse{
		Invites: pbInvites,
		Pagination: &commonpb.PaginationResponse{
			Page:       int32(page),
			PageSize:   int32(pageSize),
			TotalItems: int64(count),
			TotalPages: int32((count + pageSize - 1) / pageSize),
		},
	}, nil
}

// Helpers

func (h *OrgHandler) orgToProto(org *models.Organization) *pb.Organization {
	if org == nil {
		return nil
	}
	return &pb.Organization{
		Id:          org.ID,
		Name:        org.Name,
		Slug:        org.Slug,
		Description: org.Description,
		LogoUrl:     org.LogoURL,
		// Status:      pb.OrgStatus(pb.OrgStatus_value[string(org.Status)]), // Need mapping
		// Plan:        pb.OrgPlan(pb.OrgPlan_value[string(org.Plan)]),     // Need mapping
		Settings:    org.Settings,
		CreatedAt:   timestamppb.New(org.CreatedAt),
		UpdatedAt:   timestamppb.New(org.UpdatedAt),
	}
}

func (h *OrgHandler) memberToProto(m *models.OrgMember) *pb.OrgMember {
	if m == nil {
		return nil
	}
	return &pb.OrgMember{
		Id:             m.ID,
		OrganizationId: m.OrganizationID,
		UserId:         m.UserID,
		// Role:           pb.OrgRole(pb.OrgRole_value[string(m.Role)]), // Need mapping
		JoinedAt:       timestamppb.New(m.JoinedAt),
	}
}

func (h *OrgHandler) teamToProto(t *models.Team) *pb.Team {
	if t == nil {
		return nil
	}
	return &pb.Team{
		Id:             t.ID,
		OrganizationId: t.OrganizationID,
		Name:           t.Name,
		Description:    t.Description,
		CreatedAt:      timestamppb.New(t.CreatedAt),
		UpdatedAt:      timestamppb.New(t.UpdatedAt),
	}
}

func (h *OrgHandler) inviteToProto(i *models.Invite) *pb.Invite {
	if i == nil {
		return nil
	}
	return &pb.Invite{
		Id:             i.ID,
		OrganizationId: i.OrganizationID,
		Email:          i.Email,
		// Role:           pb.OrgRole(pb.OrgRole_value[string(i.Role)]), // Need mapping
		InvitedBy:      i.InvitedBy,
		Token:          i.Token,
		// Status:         pb.InviteStatus(pb.InviteStatus_value[string(i.Status)]), // Need mapping
		CreatedAt:      timestamppb.New(i.CreatedAt),
		ExpiresAt:      timestamppb.New(i.ExpiresAt),
	}
}

func (h *OrgHandler) protoRoleToModel(role pb.OrgRole) models.OrgRole {
	switch role {
	case pb.OrgRole_ORG_ROLE_OWNER:
		return models.OrgRoleOwner
	case pb.OrgRole_ORG_ROLE_ADMIN:
		return models.OrgRoleAdmin
	case pb.OrgRole_ORG_ROLE_MEMBER:
		return models.OrgRoleMember
	case pb.OrgRole_ORG_ROLE_GUEST:
		return models.OrgRoleGuest
	default:
		return models.OrgRoleMember
	}
}
