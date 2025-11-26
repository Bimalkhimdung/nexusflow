package handler

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/nexusflow/nexusflow/pkg/logger"
	pb "github.com/nexusflow/nexusflow/pkg/proto/user/v1"
	"github.com/nexusflow/nexusflow/services/user-service/internal/models"
	"github.com/nexusflow/nexusflow/services/user-service/internal/service"
)

// AuthHandler implements the AuthService gRPC service
type AuthHandler struct {
	pb.UnimplementedAuthServiceServer
	authService *service.AuthService
	log         *logger.Logger
}

// NewAuthHandler creates a new auth handler
func NewAuthHandler(authService *service.AuthService, log *logger.Logger) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		log:         log,
	}
}

// Register handles user registration
func (h *AuthHandler) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	h.log.Sugar().Infow("Register request", "email", req.Email)

	// Validate input
	if req.Email == "" {
		return nil, status.Error(codes.InvalidArgument, "email is required")
	}
	if req.Password == "" {
		return nil, status.Error(codes.InvalidArgument, "password is required")
	}
	if req.DisplayName == "" {
		return nil, status.Error(codes.InvalidArgument, "display name is required")
	}

	// Register user
	user, err := h.authService.Register(ctx, req.Email, req.Password, req.DisplayName)
	if err != nil {
		if err == service.ErrEmailAlreadyExists {
			return nil, status.Error(codes.AlreadyExists, "email already exists")
		}
		if err == service.ErrWeakPassword {
			return nil, status.Error(codes.InvalidArgument, "password must be at least 8 characters")
		}
		h.log.Sugar().Errorw("Failed to register user", "error", err)
		return nil, status.Error(codes.Internal, "failed to register user")
	}

	// Generate JWT token
	token, err := h.authService.GenerateJWT(user.ID, user.Email)
	if err != nil {
		h.log.Sugar().Errorw("Failed to generate token", "error", err)
		return nil, status.Error(codes.Internal, "failed to generate token")
	}

	return &pb.RegisterResponse{
		User:  h.modelToProto(user),
		Token: token,
	}, nil
}

// Login handles user authentication
func (h *AuthHandler) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	h.log.Sugar().Infow("Login request", "email", req.Email)

	// Validate input
	if req.Email == "" {
		return nil, status.Error(codes.InvalidArgument, "email is required")
	}
	if req.Password == "" {
		return nil, status.Error(codes.InvalidArgument, "password is required")
	}

	// Authenticate user
	token, user, err := h.authService.Login(ctx, req.Email, req.Password)
	if err != nil {
		if err == service.ErrInvalidCredentials {
			return nil, status.Error(codes.Unauthenticated, "invalid email or password")
		}
		if err == service.ErrEmailNotVerified {
			return nil, status.Error(codes.FailedPrecondition, "email not verified")
		}
		h.log.Sugar().Errorw("Failed to login", "error", err)
		return nil, status.Error(codes.Internal, "failed to login")
	}

	// Generate refresh token (valid for 7 days)
	refreshToken, err := h.authService.GenerateJWT(user.ID, user.Email)
	if err != nil {
		h.log.Sugar().Errorw("Failed to generate refresh token", "error", err)
		// Don't fail login if refresh token generation fails
		refreshToken = ""
	}

	return &pb.LoginResponse{
		User:         h.modelToProto(user),
		Token:        token,
		RefreshToken: refreshToken,
	}, nil
}

// VerifyEmail handles email verification
func (h *AuthHandler) VerifyEmail(ctx context.Context, req *pb.VerifyEmailRequest) (*pb.VerifyEmailResponse, error) {
	tokenLog := "empty"
	if len(req.Token) > 10 {
		tokenLog = req.Token[:10] + "..."
	} else {
		tokenLog = req.Token
	}
	h.log.Sugar().Infow("VerifyEmail request", "token", tokenLog)

	// Validate input
	if req.Token == "" {
		return nil, status.Error(codes.InvalidArgument, "token is required")
	}

	// Verify email
	if err := h.authService.VerifyEmail(ctx, req.Token); err != nil {
		if err == service.ErrInvalidToken {
			return nil, status.Error(codes.NotFound, "invalid or expired token")
		}
		h.log.Sugar().Errorw("Failed to verify email", "error", err)
		return nil, status.Error(codes.Internal, "failed to verify email")
	}

	return &pb.VerifyEmailResponse{
		Success: true,
		Message: "Email verified successfully",
	}, nil
}

// RequestPasswordReset handles password reset requests
func (h *AuthHandler) RequestPasswordReset(ctx context.Context, req *pb.RequestPasswordResetRequest) (*pb.RequestPasswordResetResponse, error) {
	h.log.Sugar().Infow("RequestPasswordReset request", "email", req.Email)

	// Validate input
	if req.Email == "" {
		return nil, status.Error(codes.InvalidArgument, "email is required")
	}

	// Request password reset
	token, err := h.authService.RequestPasswordReset(ctx, req.Email)
	if err != nil {
		h.log.Sugar().Errorw("Failed to request password reset", "error", err)
		return nil, status.Error(codes.Internal, "failed to request password reset")
	}

	// Note: In production, don't return the token in the response
	// Instead, send it via email. For now, we return it for testing.
	return &pb.RequestPasswordResetResponse{
		Success: true,
		Message: "Password reset email sent",
		Token:   token, // Remove this in production
	}, nil
}

// ResetPassword handles password reset with token
func (h *AuthHandler) ResetPassword(ctx context.Context, req *pb.ResetPasswordRequest) (*pb.ResetPasswordResponse, error) {
	tokenLog := "empty"
	if len(req.Token) > 10 {
		tokenLog = req.Token[:10] + "..."
	} else {
		tokenLog = req.Token
	}
	h.log.Sugar().Infow("ResetPassword request", "token", tokenLog)

	// Validate input
	if req.Token == "" {
		return nil, status.Error(codes.InvalidArgument, "token is required")
	}
	if req.NewPassword == "" {
		return nil, status.Error(codes.InvalidArgument, "new password is required")
	}

	// Reset password
	if err := h.authService.ResetPassword(ctx, req.Token, req.NewPassword); err != nil {
		if err == service.ErrInvalidToken {
			return nil, status.Error(codes.NotFound, "invalid or expired token")
		}
		if err == service.ErrWeakPassword {
			return nil, status.Error(codes.InvalidArgument, "password must be at least 8 characters")
		}
		h.log.Sugar().Errorw("Failed to reset password", "error", err)
		return nil, status.Error(codes.Internal, "failed to reset password")
	}

	return &pb.ResetPasswordResponse{
		Success: true,
		Message: "Password reset successfully",
	}, nil
}

// RefreshToken handles token refresh
func (h *AuthHandler) RefreshToken(ctx context.Context, req *pb.RefreshTokenRequest) (*pb.RefreshTokenResponse, error) {
	h.log.Sugar().Infow("RefreshToken request")

	// Validate input
	if req.RefreshToken == "" {
		return nil, status.Error(codes.InvalidArgument, "refresh token is required")
	}

	// Refresh token
	newToken, err := h.authService.RefreshToken(ctx, req.RefreshToken)
	if err != nil {
		if err == service.ErrInvalidToken {
			return nil, status.Error(codes.Unauthenticated, "invalid or expired refresh token")
		}
		if err == service.ErrUserNotFound {
			return nil, status.Error(codes.NotFound, "user not found")
		}
		h.log.Sugar().Errorw("Failed to refresh token", "error", err)
		return nil, status.Error(codes.Internal, "failed to refresh token")
	}

	return &pb.RefreshTokenResponse{
		Token: newToken,
	}, nil
}

// Logout handles user logout (placeholder for future session management)
func (h *AuthHandler) Logout(ctx context.Context, req *pb.LogoutRequest) (*pb.LogoutResponse, error) {
	h.log.Sugar().Infow("Logout request")

	// TODO: Implement session invalidation/token blacklisting
	// For now, logout is handled client-side by removing the token

	return &pb.LogoutResponse{
		Success: true,
		Message: "Logged out successfully",
	}, nil
}

// modelToProto converts a user model to protobuf
func (h *AuthHandler) modelToProto(user *models.User) *pb.User {
	if user == nil {
		return nil
	}

	var lastLoginAt *timestamppb.Timestamp
	if user.LastLoginAt != nil {
		lastLoginAt = timestamppb.New(*user.LastLoginAt)
	}

	var deletedAt *timestamppb.Timestamp
	if user.DeletedAt != nil {
		deletedAt = timestamppb.New(*user.DeletedAt)
	}

	return &pb.User{
		Id:            user.ID,
		Email:         user.Email,
		DisplayName:   user.DisplayName,
		AvatarUrl:     user.AvatarURL,
		Timezone:      user.Timezone,
		Locale:        user.Locale,
		Status:        string(user.Status),
		EmailVerified: user.EmailVerified,
		OauthProvider: user.OAuthProvider,
		LastLoginAt:   lastLoginAt,
		CreatedAt:     timestamppb.New(user.CreatedAt),
		UpdatedAt:     timestamppb.New(user.UpdatedAt),
		DeletedAt:     deletedAt,
	}
}
