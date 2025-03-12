package repositories

import (
	"context"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"itv-movie/internal/models"
	"itv-movie/internal/storage/database"
)

// CountryRepository handles database operations for countries
type CountryRepository struct {
	db *gorm.DB
}

// NewCountryRepository creates a new country repository
func NewCountryRepository(postgres *database.PostgresDB) *CountryRepository {
	return &CountryRepository{db: postgres.DB}
}

func (r *CountryRepository) Create(ctx context.Context, lang *models.Country) (*models.Country, error) {
	if err := r.db.WithContext(ctx).Create(lang).Error; err != nil {
		return nil, err
	}
	return lang, nil
}

func (r *CountryRepository) GetAll(ctx context.Context, page, limit int) ([]*models.Country, error) {
	var countries []*models.Country
	offset := (page - 1) * limit

	if err := r.db.WithContext(ctx).Offset(offset).Limit(limit).Find(&countries).Error; err != nil {
		return nil, err
	}

	return countries, nil
}

func (r *CountryRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Country, error) {
	var country models.Country

	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&country).Error; err != nil {
		return nil, err
	}

	return &country, nil
}

func (r *CountryRepository) Update(ctx context.Context, lang *models.Country) (*models.Country, error) {
	if err := r.db.WithContext(ctx).Save(lang).Error; err != nil {
		return nil, err
	}

	return lang, nil
}

func (r *CountryRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&models.Country{}, id).Error
}

func (r *CountryRepository) Count(ctx context.Context) (int, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&models.Country{}).Count(&count).Error; err != nil {
		return 0, err
	}
	return int(count), nil
}
