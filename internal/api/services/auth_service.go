package services

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"itv-movie/internal/models"
	"itv-movie/internal/storage/database/repositories"
	"time"
)

var (
	ErrInvalidCredentials = errors.New("invalid username or password")
	ErrUserInactive       = errors.New("user account is inactive")
	ErrUserExists         = errors.New("user with this email or username already exists")
	ErrSessionInvalid     = errors.New("session is invalid or expired")
)

// DefaultSessionDuration session duration: 7 days
const DefaultSessionDuration = 7 * 24 * time.Hour

type AuthService struct {
	userRepo    *repositories.UserRepository
	sessionRepo *repositories.SessionRepository
}

func NewAuthService(userRepo *repositories.UserRepository, sessionRepo *repositories.SessionRepository) *AuthService {
	return &AuthService{
		userRepo:    userRepo,
		sessionRepo: sessionRepo,
	}
}

func (s *AuthService) RegisterUser(ctx context.Context, userData *models.User) (*models.User, error) {
	existingUser, err := s.userRepo.GetByEmail(ctx, userData.Email)
	if err == nil && existingUser != nil {
		return nil, ErrUserExists
	}

	existingUser, err = s.userRepo.GetByUsername(ctx, userData.Username)
	if err == nil && existingUser != nil {
		return nil, ErrUserExists
	}

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

	session, err := user.GenerateSession(DefaultSessionDuration, userAgent, ipAddress)
	if err != nil {
		return nil, nil, err
	}

	session, err = s.sessionRepo.Create(ctx, session)
	if err != nil {
		return nil, nil, err
	}

	return user, session, nil
}

func (s *AuthService) Logout(ctx context.Context, token string) error {
	return s.sessionRepo.RevokeByToken(ctx, token)
}

func (s *AuthService) ValidateSession(ctx context.Context, token string) (*models.Session, error) {
	session, err := s.sessionRepo.GetByToken(ctx, token)
	if err != nil {
		return nil, ErrSessionInvalid
	}

	if !session.IsSessionValid() {
		return nil, ErrSessionInvalid
	}

	return session, nil
}

func (s *AuthService) UpdateUserStatus(ctx context.Context, userID uuid.UUID, active bool) error {
	return s.userRepo.UpdateStatus(ctx, userID, active)
}

func (s *AuthService) DeleteUser(ctx context.Context, userID uuid.UUID) error {
	if err := s.sessionRepo.RevokeAllForUser(ctx, userID); err != nil {
		return err
	}

	return s.userRepo.Delete(ctx, userID)
}
