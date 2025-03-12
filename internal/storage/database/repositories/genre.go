package repositories

import (
	"context"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"itv-movie/internal/models"
	"itv-movie/internal/storage/database"
)

// GenreRepository handles database operations for genres
type GenreRepository struct {
	db *gorm.DB
}

// NewGenreRepository creates a new genre repository
func NewGenreRepository(postgres *database.PostgresDB) *GenreRepository {
	return &GenreRepository{db: postgres.DB}
}

func (r *GenreRepository) Create(ctx context.Context, lang *models.Genre) (*models.Genre, error) {
	if err := r.db.WithContext(ctx).Create(lang).Error; err != nil {
		return nil, err
	}
	return lang, nil
}

func (r *GenreRepository) GetAll(ctx context.Context, page, limit int) ([]*models.Genre, error) {
	var genres []*models.Genre
	offset := (page - 1) * limit

	if err := r.db.WithContext(ctx).Offset(offset).Limit(limit).Find(&genres).Error; err != nil {
		return nil, err
	}

	return genres, nil
}

func (r *GenreRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Genre, error) {
	var genre models.Genre

	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&genre).Error; err != nil {
		return nil, err
	}

	return &genre, nil
}

func (r *GenreRepository) Update(ctx context.Context, lang *models.Genre) (*models.Genre, error) {
	if err := r.db.WithContext(ctx).Save(lang).Error; err != nil {
		return nil, err
	}

	return lang, nil
}

func (r *GenreRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&models.Genre{}, id).Error
}

func (r *GenreRepository) Count(ctx context.Context) (int, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&models.Genre{}).Count(&count).Error; err != nil {
		return 0, err
	}
	return int(count), nil
}
