package services

import (
	"context"
	"github.com/google/uuid"
	"itv-movie/internal/models"
	"itv-movie/internal/storage/database"
)

// LanguageService handles business logic for lang
type LanguageService struct {
	languageRepo *database.LanguageRepository
}

// NewLanguageService creates a new lang service
func NewLanguageService(
	languageRepo *database.LanguageRepository,
) *LanguageService {
	return &LanguageService{
		languageRepo: languageRepo,
	}
}

func (s *LanguageService) CreateLanguage(ctx context.Context, newLang *models.Language) (*models.Language, error) {
	//TODO actual implementation
	return nil, nil
}

func (s *LanguageService) GetAllLanguages(ctx context.Context, page, limit int) ([]*models.Language, error) {
	//TODO actual implementation
	return nil, nil
}

func (s *LanguageService) GetLanguage(ctx context.Context, id uuid.UUID) (*models.Language, error) {
	//TODO actual implementation
	return nil, nil
}

func (s *LanguageService) UpdateLanguage(ctx context.Context, lang *models.Language) (*models.Language, error) {
	//TODO actual implementation
	return nil, nil
}

func (s *LanguageService) DeleteLanguage(ctx context.Context, newLang *models.Language) error {
	//TODO actual implementation
	return nil
}
