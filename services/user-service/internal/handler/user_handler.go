package handler

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/nexusflow/nexusflow/pkg/logger"
	commonpb "github.com/nexusflow/nexusflow/pkg/proto/common/v1"
	pb "github.com/nexusflow/nexusflow/pkg/proto/user/v1"
	"github.com/nexusflow/nexusflow/services/user-service/internal/models"
	"github.com/nexusflow/nexusflow/services/user-service/internal/service"
)

// UserHandler implements the UserService gRPC service
type UserHandler struct {
	pb.UnimplementedUserServiceServer
	service *service.UserService
	log     *logger.Logger
}

// NewUserHandler creates a new user handler
func NewUserHandler(service *service.UserService, log *logger.Logger) *UserHandler {
	return &UserHandler{
		service: service,
		log:     log,
	}
}

// GetUser retrieves a user by ID
func (h *UserHandler) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	user, err := h.service.GetUser(ctx, req.Id)
	if err != nil {
		h.log.Sugar().Errorw("Failed to get user", "error", err, "user_id", req.Id)
		return nil, status.Error(codes.NotFound, "user not found")
	}

	return &pb.GetUserResponse{
		User: h.modelToProto(user),
	}, nil
}

// GetUserByEmail retrieves a user by email
func (h *UserHandler) GetUserByEmail(ctx context.Context, req *pb.GetUserByEmailRequest) (*pb.GetUserByEmailResponse, error) {
	user, err := h.service.GetUserByEmail(ctx, req.Email)
	if err != nil {
		h.log.Sugar().Errorw("Failed to get user by email", "error", err, "email", req.Email)
		return nil, status.Error(codes.NotFound, "user not found")
	}

	return &pb.GetUserByEmailResponse{
		User: h.modelToProto(user),
	}, nil
}

// CreateUser creates a new user
func (h *UserHandler) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	// TODO: Extract organization_id from auth context
	// For now, use a default UUID
	defaultOrgID := "00000000-0000-0000-0000-000000000000"
	systemUserID := "00000000-0000-0000-0000-000000000001"
	
	input := service.CreateUserInput{
		OrganizationID: defaultOrgID,
		Email:          req.Email,
		DisplayName:    req.DisplayName,
		AvatarURL:      req.AvatarUrl,
		Timezone:       req.Timezone,
		Locale:         req.Locale,
		CreatedBy:      systemUserID,
	}

	user, err := h.service.CreateUser(ctx, input)
	if err != nil {
		h.log.Sugar().Errorw("Failed to create user", "error", err, "email", req.Email)
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	return &pb.CreateUserResponse{
		User: h.modelToProto(user),
	}, nil
}

// UpdateUser updates a user
func (h *UserHandler) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UpdateUserResponse, error) {
	systemUserID := "00000000-0000-0000-0000-000000000001"
	
	input := service.UpdateUserInput{
		UpdatedBy: systemUserID,
	}

	if req.DisplayName != nil {
		input.DisplayName = req.DisplayName
	}
	if req.AvatarUrl != nil {
		input.AvatarURL = req.AvatarUrl
	}
	if req.Timezone != nil {
		input.Timezone = req.Timezone
	}
	if req.Locale != nil {
		input.Locale = req.Locale
	}
	if req.Status != nil {
		status := models.UserStatus(*req.Status)
		input.Status = &status
	}

	user, err := h.service.UpdateUser(ctx, req.Id, input)
	if err != nil {
		h.log.Sugar().Errorw("Failed to update user", "error", err, "user_id", req.Id)
		return nil, status.Error(codes.Internal, "failed to update user")
	}

	return &pb.UpdateUserResponse{
		User: h.modelToProto(user),
	}, nil
}

// DeleteUser deletes a user
func (h *UserHandler) DeleteUser(ctx context.Context, req *pb.DeleteUserRequest) (*pb.DeleteUserResponse, error) {
	if err := h.service.DeleteUser(ctx, req.Id); err != nil {
		h.log.Sugar().Errorw("Failed to delete user", "error", err, "user_id", req.Id)
		return nil, status.Error(codes.Internal, "failed to delete user")
	}

	return &pb.DeleteUserResponse{
		Response: &commonpb.SuccessResponse{
			Success: true,
			Message: "User deleted successfully",
		},
	}, nil
}

// ListUsers lists users with pagination
func (h *UserHandler) ListUsers(ctx context.Context, req *pb.ListUsersRequest) (*pb.ListUsersResponse, error) {
	page := int(req.Pagination.Page)
	pageSize := int(req.Pagination.PageSize)

	// TODO: Extract from auth context
	orgID := "00000000-0000-0000-0000-000000000000"
	if len(req.OrganizationIds) > 0 {
		orgID = req.OrganizationIds[0]
	}

	users, total, err := h.service.ListUsers(ctx, orgID, page, pageSize)
	if err != nil {
		h.log.Sugar().Errorw("Failed to list users", "error", err)
		return nil, status.Error(codes.Internal, "failed to list users")
	}

	totalPages := int32((total + pageSize - 1) / pageSize)

	return &pb.ListUsersResponse{
		Users: h.modelsToProtos(users),
		Pagination: &commonpb.PaginationResponse{
			Page:       int32(page),
			PageSize:   int32(pageSize),
			TotalItems: int64(total),
			TotalPages: totalPages,
			HasNext:    page < int(totalPages),
			HasPrevious: page > 1,
		},
	}, nil
}

// SearchUsers searches for users
func (h *UserHandler) SearchUsers(ctx context.Context, req *pb.SearchUsersRequest) (*pb.SearchUsersResponse, error) {
	page := int(req.Pagination.Page)
	pageSize := int(req.Pagination.PageSize)

	users, total, err := h.service.SearchUsers(ctx, req.Query, req.OrganizationIds, page, pageSize)
	if err != nil {
		h.log.Sugar().Errorw("Failed to search users", "error", err, "query", req.Query)
		return nil, status.Error(codes.Internal, "failed to search users")
	}

	totalPages := int32((total + pageSize - 1) / pageSize)

	return &pb.SearchUsersResponse{
		Users: h.modelsToProtos(users),
		Pagination: &commonpb.PaginationResponse{
			Page:       int32(page),
			PageSize:   int32(pageSize),
			TotalItems: int64(total),
			TotalPages: totalPages,
			HasNext:    page < int(totalPages),
			HasPrevious: page > 1,
		},
	}, nil
}

// GetUserProfile retrieves a user profile
func (h *UserHandler) GetUserProfile(ctx context.Context, req *pb.GetUserProfileRequest) (*pb.GetUserProfileResponse, error) {
	user, err := h.service.GetUser(ctx, req.Id)
	if err != nil {
		h.log.Sugar().Errorw("Failed to get user profile", "error", err, "user_id", req.Id)
		return nil, status.Error(codes.NotFound, "user not found")
	}

	return &pb.GetUserProfileResponse{
		Profile: &pb.UserProfile{
			Id:          user.ID,
			DisplayName: user.DisplayName,
			AvatarUrl:   user.AvatarURL,
		},
	}, nil
}

// UpdateUserPreferences updates user preferences
func (h *UserHandler) UpdateUserPreferences(ctx context.Context, req *pb.UpdateUserPreferencesRequest) (*pb.UpdateUserPreferencesResponse, error) {
	preferences := make(map[string]interface{})
	for k, v := range req.Preferences {
		preferences[k] = v
	}

	if err := h.service.UpdateUserPreferences(ctx, req.Id, preferences); err != nil {
		h.log.Sugar().Errorw("Failed to update preferences", "error", err, "user_id", req.Id)
		return nil, status.Error(codes.Internal, "failed to update preferences")
	}

	user, _ := h.service.GetUser(ctx, req.Id)
	return &pb.UpdateUserPreferencesResponse{
		User: h.modelToProto(user),
	}, nil
}

// modelToProto converts a user model to protobuf
func (h *UserHandler) modelToProto(user *models.User) *pb.User {
	if user == nil {
		return nil
	}

	pbUser := &pb.User{
		Id:              user.ID,
		Email:           user.Email,
		DisplayName:     user.DisplayName,
		AvatarUrl:       user.AvatarURL,
		Timezone:        user.Timezone,
		Locale:          user.Locale,
		Status:          string(user.Status),
		OrganizationIds: []string{user.OrganizationID},
		EmailVerified:   user.EmailVerified,
		OauthProvider:   user.OAuthProvider,
		CreatedAt:       timestamppb.New(user.CreatedAt),
		UpdatedAt:       timestamppb.New(user.UpdatedAt),
	}

	if user.LastLoginAt != nil {
		pbUser.LastLoginAt = timestamppb.New(*user.LastLoginAt)
	}
	
	if user.DeletedAt != nil {
		pbUser.DeletedAt = timestamppb.New(*user.DeletedAt)
	}

	// Convert preferences
	if user.Preferences != nil {
		pbUser.Preferences = make(map[string]string)
		for k, v := range user.Preferences {
			if str, ok := v.(string); ok {
				pbUser.Preferences[k] = str
			}
		}
	}

	return pbUser
}

// modelsToProtos converts multiple user models to protobuf
func (h *UserHandler) modelsToProtos(users []*models.User) []*pb.User {
	pbUsers := make([]*pb.User, len(users))
	for i, user := range users {
		pbUsers[i] = h.modelToProto(user)
	}
	return pbUsers
}
