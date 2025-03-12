package services

import (
	"context"
	"github.com/google/uuid"
	"itv-movie/internal/models"
	"itv-movie/internal/storage/database/repositories"
)

// LanguageService handles business logic for lang
type LanguageService struct {
	languageRepo *repositories.LanguageRepository
}

// NewLanguageService creates a new lang service
func NewLanguageService(
	languageRepo *repositories.LanguageRepository,
) *LanguageService {
	return &LanguageService{
		languageRepo: languageRepo,
	}
}

func (s *LanguageService) CreateLanguage(ctx context.Context, newLang *models.Language) (*models.Language, error) {
	createdLang, err := s.languageRepo.Create(ctx, newLang)
	if err != nil {
		return nil, err
	}

	return createdLang, nil
}

func (s *LanguageService) GetAllLanguages(ctx context.Context, page, limit int) ([]*models.Language, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	languages, err := s.languageRepo.GetAll(ctx, page, limit)
	if err != nil {
		return nil, err
	}

	return languages, nil
}

func (s *LanguageService) GetLanguage(ctx context.Context, id uuid.UUID) (*models.Language, error) {
	language, err := s.languageRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return language, nil
}

func (s *LanguageService) UpdateLanguage(ctx context.Context, lang *models.Language) (*models.Language, error) {
	_, err := s.languageRepo.GetByID(ctx, lang.ID) // check if exists
	if err != nil {
		return nil, err
	}

	updatedLang, err := s.languageRepo.Update(ctx, lang)
	if err != nil {
		return nil, err
	}

	return updatedLang, nil
}

func (s *LanguageService) DeleteLanguage(ctx context.Context, lang *models.Language) error {
	return s.languageRepo.Delete(ctx, lang.ID)
}

func (s *LanguageService) GetTotalLanguageCount(ctx context.Context) (int, error) {
	return s.languageRepo.Count(ctx)
}
