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

func (s *AuthService) Login(ctx context.Context, username, password, userAgent, ipAddress string) (*models.User, *models.Session, error) {
	user, err := s.userRepo.GetByUsername(ctx, username)
	if err != nil {
		return nil, nil, ErrInvalidCredentials
	}

	if !user.Active {
		return nil, nil, ErrUserInactive
	}

	if !user.CheckPassword(password) {
		return nil, nil, ErrInvalidCredentials
	}

	if err = s.userRepo.UpdateLastLogin(ctx, user.ID); err != nil {
		return nil, nil, err
	}

	session, err := user.GenerateTokens(&s.config.Internal.Jwt, userAgent, ipAddress)
	if err != nil {
		return nil, nil, err
	}

	session, err = s.sessionRepo.Create(ctx, session)
	if err != nil {
		return nil, nil, err
	}

	return user, session, nil
}

func (s *AuthService) Logout(ctx context.Context, accessToken string) error {
	session, err := s.sessionRepo.GetByAccessToken(ctx, accessToken)
	if err != nil {
		return ErrSessionInvalid
	}

	return s.sessionRepo.RevokeByID(ctx, session.ID)
}

func (s *AuthService) ValidateAccessToken(ctx context.Context, accessToken string) (*jwtpkg.CustomClaims, error) {
	claims, err := jwtpkg.ValidateToken(accessToken, s.config.Internal.Jwt.Secret)
	if err != nil {
		return nil, err
	}

	if claims.TokenType != jwtpkg.AccessToken {
		return nil, ErrInvalidToken
	}

	session, err := s.sessionRepo.GetByAccessToken(ctx, accessToken)
	if err != nil || session.IsRevoked || !session.IsAccessTokenValid() {
		return nil, ErrSessionInvalid
	}

	return claims, nil
}

func (s *AuthService) RefreshTokens(ctx context.Context, refreshToken, userAgent, ipAddress string) (*models.Session, error) {
	claims, err := jwtpkg.ValidateToken(refreshToken, s.config.Internal.Jwt.Secret)
	if err != nil {
		if errors.Is(err, jwtpkg.ErrExpiredToken) {
			return nil, ErrRefreshTokenExpired
		}
		return nil, err
	}

	if claims.TokenType != jwtpkg.RefreshToken {
		return nil, ErrInvalidToken
	}

	oldSession, err := s.sessionRepo.GetByRefreshToken(ctx, refreshToken)
	if err != nil {
		return nil, ErrSessionInvalid
	}

	if !oldSession.IsRefreshTokenValid() {
		return nil, ErrRefreshTokenExpired
	}

	userID, err := uuid.Parse(claims.UserID)
	if err != nil {
		return nil, ErrInvalidToken
	}

	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	err = s.sessionRepo.RevokeByID(ctx, oldSession.ID)
	if err != nil {
		return nil, err
	}

	// Generate new tokens
	newSession, err := user.GenerateTokens(&s.config.Internal.Jwt, userAgent, ipAddress)
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

func (s *AuthService) UpdateUserStatus(ctx context.Context, userID uuid.UUID, active bool) error {
	return s.userRepo.UpdateStatus(ctx, userID, active)
}

func (s *AuthService) DeleteUser(ctx context.Context, userID uuid.UUID) error {
	// First revoke all sessions for this user
	if err := s.sessionRepo.RevokeAllForUser(ctx, userID); err != nil {
		return err
	}

	// Then delete the user (soft delete)
	return s.userRepo.Delete(ctx, userID)
}

func (s *AuthService) AdminExists(ctx context.Context) (bool, error) {
	count, err := s.userRepo.FindAdmin(ctx)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (s *AuthService) GetAllUsers(ctx context.Context, page, limit int) ([]*models.User, int64, error) {
	if page < 1 {
		page = 1
	}

	if limit < 1 || limit > 1000 {
		limit = 10
	}

	return s.userRepo.GetAllUsers(ctx, page, limit)
}
