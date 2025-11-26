package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/nexusflow/nexusflow/pkg/logger"
	"github.com/nexusflow/nexusflow/services/user-service/internal/models"
	"github.com/nexusflow/nexusflow/services/user-service/internal/repository"
)

var (
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrUserNotFound       = errors.New("user not found")
	ErrEmailAlreadyExists = errors.New("email already exists")
	ErrEmailNotVerified   = errors.New("email not verified")
	ErrInvalidToken       = errors.New("invalid or expired token")
	ErrWeakPassword       = errors.New("password must be at least 8 characters")
)

// AuthService handles authentication operations
type AuthService struct {
	userRepo  *repository.UserRepository
	jwtSecret string
	jwtExpiry time.Duration
	log       *logger.Logger
}

// JWTClaims represents the JWT token claims
type JWTClaims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

// NewAuthService creates a new authentication service
func NewAuthService(userRepo *repository.UserRepository, jwtSecret string, jwtExpiry time.Duration, log *logger.Logger) *AuthService {
	return &AuthService{
		userRepo:  userRepo,
		jwtSecret: jwtSecret,
		jwtExpiry: jwtExpiry,
		log:       log,
	}
}

// Register creates a new user with email and password
func (s *AuthService) Register(ctx context.Context, email, password, displayName string) (*models.User, error) {
	s.log.Sugar().Infow("Registering new user", "email", email)

	// Validate password strength
	if len(password) < 8 {
		return nil, ErrWeakPassword
	}

	// Check if user already exists
	existingUser, err := s.userRepo.GetByEmail(ctx, email)
	if err == nil && existingUser != nil {
		return nil, ErrEmailAlreadyExists
	}

	// Create new user
	user := &models.User{
		Email:          email,
		DisplayName:    displayName,
		EmailVerified:  false,
		Status:         models.UserStatusActive,
		OAuthProvider:  "email",
		OrganizationID: uuid.New().String(),
	}

	// Set password
	if err := user.SetPassword(password); err != nil {
		s.log.Sugar().Errorw("Failed to hash password", "error", err)
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Generate verification token
	if _, err := user.GenerateVerificationToken(); err != nil {
		s.log.Sugar().Errorw("Failed to generate verification token", "error", err)
		return nil, fmt.Errorf("failed to generate verification token: %w", err)
	}

	// Save user to database
	if err := s.userRepo.Create(ctx, user); err != nil {
		s.log.Sugar().Errorw("Failed to create user", "error", err)
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	s.log.Sugar().Infow("User registered successfully", "user_id", user.ID, "email", email)
	return user, nil
}

// Login authenticates a user and returns a JWT token
func (s *AuthService) Login(ctx context.Context, email, password string) (string, *models.User, error) {
	s.log.Sugar().Infow("User login attempt", "email", email)

	// Get user by email
	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		s.log.Sugar().Warnw("User not found", "email", email)
		return "", nil, ErrInvalidCredentials
	}

	// Check if user is active
	if !user.IsActive() {
		s.log.Sugar().Warnw("Inactive user login attempt", "user_id", user.ID)
		return "", nil, errors.New("user account is inactive")
	}

	// Verify password
	if !user.CheckPassword(password) {
		s.log.Sugar().Warnw("Invalid password", "user_id", user.ID)
		return "", nil, ErrInvalidCredentials
	}

	// Optional: Check if email is verified
	// Uncomment if you want to enforce email verification
	// if !user.EmailVerified {
	// 	return "", nil, ErrEmailNotVerified
	// }

	// Update last login time
	user.UpdateLastLogin()
	if err := s.userRepo.Update(ctx, user); err != nil {
		s.log.Sugar().Errorw("Failed to update last login", "error", err)
		// Don't fail login if we can't update last login time
	}

	// Generate JWT token
	token, err := s.GenerateJWT(user.ID, user.Email)
	if err != nil {
		s.log.Sugar().Errorw("Failed to generate JWT", "error", err)
		return "", nil, fmt.Errorf("failed to generate token: %w", err)
	}

	s.log.Sugar().Infow("User logged in successfully", "user_id", user.ID)
	return token, user, nil
}

// VerifyEmail verifies a user's email with the verification token
func (s *AuthService) VerifyEmail(ctx context.Context, token string) error {
	s.log.Sugar().Infow("Verifying email", "token", token[:10]+"...")

	// Find user by verification token
	user, err := s.userRepo.GetByVerificationToken(ctx, token)
	if err != nil {
		s.log.Sugar().Warnw("Invalid verification token", "token", token[:10]+"...")
		return ErrInvalidToken
	}

	// Mark email as verified
	user.MarkEmailAsVerified()

	// Update user in database
	if err := s.userRepo.Update(ctx, user); err != nil {
		s.log.Sugar().Errorw("Failed to update user", "error", err)
		return fmt.Errorf("failed to verify email: %w", err)
	}

	s.log.Sugar().Infow("Email verified successfully", "user_id", user.ID)
	return nil
}

// RequestPasswordReset generates a reset token and returns it
func (s *AuthService) RequestPasswordReset(ctx context.Context, email string) (string, error) {
	s.log.Sugar().Infow("Password reset requested", "email", email)

	// Get user by email
	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		// Don't reveal if user exists or not for security
		s.log.Sugar().Warnw("Password reset for non-existent user", "email", email)
		return "", nil // Return success even if user doesn't exist
	}

	// Generate reset token
	token, err := user.GenerateResetToken()
	if err != nil {
		s.log.Sugar().Errorw("Failed to generate reset token", "error", err)
		return "", fmt.Errorf("failed to generate reset token: %w", err)
	}

	// Update user in database
	if err := s.userRepo.Update(ctx, user); err != nil {
		s.log.Sugar().Errorw("Failed to update user with reset token", "error", err)
		return "", fmt.Errorf("failed to save reset token: %w", err)
	}

	s.log.Sugar().Infow("Password reset token generated", "user_id", user.ID)
	return token, nil
}

// ResetPassword resets a user's password using the reset token
func (s *AuthService) ResetPassword(ctx context.Context, token, newPassword string) error {
	s.log.Sugar().Infow("Resetting password", "token", token[:10]+"...")

	// Validate password strength
	if len(newPassword) < 8 {
		return ErrWeakPassword
	}

	// Find user by reset token
	user, err := s.userRepo.GetByResetToken(ctx, token)
	if err != nil {
		s.log.Sugar().Warnw("Invalid reset token", "token", token[:10]+"...")
		return ErrInvalidToken
	}

	// Check if token is still valid
	if !user.IsResetTokenValid() {
		s.log.Sugar().Warnw("Expired reset token", "user_id", user.ID)
		return ErrInvalidToken
	}

	// Set new password
	if err := user.SetPassword(newPassword); err != nil {
		s.log.Sugar().Errorw("Failed to hash new password", "error", err)
		return fmt.Errorf("failed to hash password: %w", err)
	}

	// Clear reset token
	user.ClearResetToken()

	// Update user in database
	if err := s.userRepo.Update(ctx, user); err != nil {
		s.log.Sugar().Errorw("Failed to update user password", "error", err)
		return fmt.Errorf("failed to reset password: %w", err)
	}

	s.log.Sugar().Infow("Password reset successfully", "user_id", user.ID)
	return nil
}

// GenerateJWT creates a JWT token for a user
func (s *AuthService) GenerateJWT(userID, email string) (string, error) {
	now := time.Now()
	expiresAt := now.Add(s.jwtExpiry)

	claims := JWTClaims{
		UserID: userID,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Issuer:    "nexusflow",
			Subject:   userID,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return signedToken, nil
}

// ValidateJWT validates a JWT token and returns the user ID
func (s *AuthService) ValidateJWT(tokenString string) (string, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Verify signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.jwtSecret), nil
	})

	if err != nil {
		return "", fmt.Errorf("failed to parse token: %w", err)
	}

	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims.UserID, nil
	}

	return "", ErrInvalidToken
}

// RefreshToken generates a new JWT token for a user
func (s *AuthService) RefreshToken(ctx context.Context, oldToken string) (string, error) {
	// Validate old token
	userID, err := s.ValidateJWT(oldToken)
	if err != nil {
		return "", err
	}

	// Get user to ensure they still exist and are active
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return "", ErrUserNotFound
	}

	if !user.IsActive() {
		return "", errors.New("user account is inactive")
	}

	// Generate new token
	newToken, err := s.GenerateJWT(user.ID, user.Email)
	if err != nil {
		return "", err
	}

	return newToken, nil
}
