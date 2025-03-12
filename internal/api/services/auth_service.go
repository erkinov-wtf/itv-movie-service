package services

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"itv-movie/internal/config"
	"itv-movie/internal/models"
	jwtpkg "itv-movie/internal/pkg/jwt"
	"itv-movie/internal/storage/database/repositories"
)

var (
	ErrInvalidCredentials  = errors.New("invalid username or password")
	ErrUserInactive        = errors.New("user account is inactive")
	ErrUserExists          = errors.New("user with this email or username already exists")
	ErrSessionInvalid      = errors.New("session is invalid or expired")
	ErrInvalidToken        = errors.New("invalid token format")
	ErrRefreshTokenExpired = errors.New("refresh token has expired")
)

type AuthService struct {
	userRepo    *repositories.UserRepository
	sessionRepo *repositories.SessionRepository
	config      *config.Config
}

func NewAuthService(userRepo *repositories.UserRepository, sessionRepo *repositories.SessionRepository, config *config.Config) *AuthService {
	return &AuthService{
		userRepo:    userRepo,
		sessionRepo: sessionRepo,
		config:      config,
	}
}

// RegisterUser creates a new user account
func (s *AuthService) RegisterUser(ctx context.Context, userData *models.User) (*models.User, error) {
	// Check if user already exists with the same email
	existingUser, err := s.userRepo.GetByEmail(ctx, userData.Email)
	if err == nil && existingUser != nil {
		return nil, ErrUserExists
	}

	// Check if user already exists with the same username
	existingUser, err = s.userRepo.GetByUsername(ctx, userData.Username)
	if err == nil && existingUser != nil {
		return nil, ErrUserExists
	}

	// Create the new user
	user, err := s.userRepo.Create(ctx, userData)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// Login authenticates a user and creates a new session with JWT tokens
func (s *AuthService) Login(ctx context.Context, username, password, userAgent, ipAddress string) (*models.User, *models.Session, error) {
	// Find the user by username
	user, err := s.userRepo.GetByUsername(ctx, username)
	if err != nil {
		return nil, nil, ErrInvalidCredentials
	}

	// Check if user is active
	if !user.Active {
		return nil, nil, ErrUserInactive
	}

	// Verify password
	if !user.CheckPassword(password) {
		return nil, nil, ErrInvalidCredentials
	}

	// Update last login timestamp
	if err := s.userRepo.UpdateLastLogin(ctx, user.ID); err != nil {
		return nil, nil, err
	}

	// Generate JWT tokens and create a session
	session, err := user.GenerateTokens(&s.config.Internal.Jwt, userAgent, ipAddress)
	if err != nil {
		return nil, nil, err
	}

	// Save session to database
	session, err = s.sessionRepo.Create(ctx, session)
	if err != nil {
		return nil, nil, err
	}

	return user, session, nil
}

// Logout revokes the current session
func (s *AuthService) Logout(ctx context.Context, accessToken string) error {
	// Get the session by access token
	session, err := s.sessionRepo.GetByAccessToken(ctx, accessToken)
	if err != nil {
		return ErrSessionInvalid
	}

	// Revoke the session
	return s.sessionRepo.RevokeByID(ctx, session.ID)
}

// ValidateAccessToken validates a JWT access token and returns the claims
func (s *AuthService) ValidateAccessToken(ctx context.Context, accessToken string) (*jwtpkg.CustomClaims, error) {
	// Validate the JWT token
	claims, err := jwtpkg.ValidateToken(accessToken, s.config.JWT.SecretKey)
	if err != nil {
		return nil, err
	}

	// Ensure it's an access token
	if claims.TokenType != jwtpkg.AccessToken {
		return nil, ErrInvalidToken
	}

	// Check if the token is in a valid session (not revoked)
	session, err := s.sessionRepo.GetByAccessToken(ctx, accessToken)
	if err != nil || session.IsRevoked || !session.IsAccessTokenValid() {
		return nil, ErrSessionInvalid
	}

	return claims, nil
}

// RefreshTokens generates new access and refresh tokens using a valid refresh token
func (s *AuthService) RefreshTokens(ctx context.Context, refreshToken, userAgent, ipAddress string) (*models.Session, error) {
	// Validate the refresh token
	claims, err := jwtpkg.ValidateToken(refreshToken, s.config.JWT.SecretKey)
	if err != nil {
		if errors.Is(err, jwtpkg.ErrExpiredToken) {
			return nil, ErrRefreshTokenExpired
		}
		return nil, err
	}

	// Ensure it's a refresh token
	if claims.TokenType != jwtpkg.RefreshToken {
		return nil, ErrInvalidToken
	}

	// Get the session by refresh token
	oldSession, err := s.sessionRepo.GetByRefreshToken(ctx, refreshToken)
	if err != nil {
		return nil, ErrSessionInvalid
	}

	// Check if refresh token is valid
	if !oldSession.IsRefreshTokenValid() {
		return nil, ErrRefreshTokenExpired
	}

	// Get the user
	userID, err := uuid.Parse(claims.UserID)
	if err != nil {
		return nil, ErrInvalidToken
	}

	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Revoke the old session
	err = s.sessionRepo.RevokeByID(ctx, oldSession.ID)
	if err != nil {
		return nil, err
	}

	// Generate new tokens
	newSession, err := user.GenerateTokens(
		s.config.JWT.AccessTokenTTL,
		s.config.JWT.RefreshTokenTTL,
		s.config.JWT.SecretKey,
		userAgent,
		ipAddress,
	)
	if err != nil {
		return nil, err
	}

	// Save the new session
	newSession, err = s.sessionRepo.Create(ctx, newSession)
	if err != nil {
		return nil, err
	}

	return newSession, nil
}

// UpdateUserStatus changes a user's active status
func (s *AuthService) UpdateUserStatus(ctx context.Context, userID uuid.UUID, active bool) error {
	return s.userRepo.UpdateStatus(ctx, userID, active)
}

// DeleteUser soft-deletes a user account and revokes all sessions
func (s *AuthService) DeleteUser(ctx context.Context, userID uuid.UUID) error {
	// First revoke all sessions for this user
	if err := s.sessionRepo.RevokeAllForUser(ctx, userID); err != nil {
		return err
	}

	// Then delete the user (soft delete)
	return s.userRepo.Delete(ctx, userID)
}
