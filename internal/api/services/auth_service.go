package services

import (
	"context"
	"itv-movie/internal/models"
	"itv-movie/internal/storage/database"
)

// AuthService handles business logic for auth
type AuthService struct {
	userRepo *database.UserRepository
}

// NewAuthService creates a new auth service
func NewAuthService(
	userRepo *database.UserRepository,
) *AuthService {
	return &AuthService{
		userRepo: userRepo,
	}
}

func (s *AuthService) Login(ctx context.Context, username, password string) (*models.User, error) {
	//TODO actual implementation
	return nil, nil
}

func (s *AuthService) Register(ctx context.Context, newUser *models.User) (*models.User, error) {
	//TODO actual implementation
	return nil, nil
}

func (s *AuthService) Logout(ctx context.Context, user *models.User) error {
	//TODO actual implementation
	return nil
}

func (s *AuthService) UpdateStatus(ctx context.Context, user *models.User) (*models.User, error) {
	//TODO actual implementation
	return nil, nil
}

func (s *AuthService) DeleteUser(ctx context.Context, user *models.User) error {
	//TODO actual implementation
	return nil
}
