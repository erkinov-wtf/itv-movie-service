package repositories

import (
	"context"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"itv-movie/internal/models"
	"itv-movie/internal/storage/database"
)

// LanguageRepository handles database operations for languages
type LanguageRepository struct {
	db *gorm.DB
}

// NewLanguageRepository creates a new language repository
func NewLanguageRepository(postgres *database.PostgresDB) *LanguageRepository {
	return &LanguageRepository{db: postgres.DB}
}

func (r *LanguageRepository) Create(ctx context.Context, lang *models.Language) (*models.Language, error) {
	if err := r.db.WithContext(ctx).Create(lang).Error; err != nil {
		return nil, err
	}
	return lang, nil
}

func (r *LanguageRepository) GetAll(ctx context.Context, page, limit int) ([]*models.Language, error) {
	var languages []*models.Language
	offset := (page - 1) * limit

	if err := r.db.WithContext(ctx).Offset(offset).Limit(limit).Find(&languages).Error; err != nil {
		return nil, err
	}

	return languages, nil
}

func (r *LanguageRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Language, error) {
	var language models.Language

	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&language).Error; err != nil {
		return nil, err
	}

	return &language, nil
}

func (r *LanguageRepository) Update(ctx context.Context, lang *models.Language) (*models.Language, error) {
	if err := r.db.WithContext(ctx).Save(lang).Error; err != nil {
		return nil, err
	}

	return lang, nil
}

func (r *LanguageRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&models.Language{}, id).Error
}

func (r *LanguageRepository) Count(ctx context.Context) (int, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&models.Language{}).Count(&count).Error; err != nil {
		return 0, err
	}
	return int(count), nil
}

func (r *LanguageRepository) GetByCode(ctx context.Context, code string) (*models.Language, error) {
	var language models.Language

	if err := r.db.WithContext(ctx).Where("code = ?", code).First(&language).Error; err != nil {
		return nil, err
	}

	return &language, nil
}
