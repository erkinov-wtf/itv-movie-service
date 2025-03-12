package services

import (
	"context"
	"github.com/google/uuid"
	"itv-movie/internal/models"
	"itv-movie/internal/storage/database/repositories"
)

// CountryService handles business logic for country
type CountryService struct {
	countryRepo *repositories.CountryRepository
}

// NewCountryService creates a new country service
func NewCountryService(
	countryRepo *repositories.CountryRepository,
) *CountryService {
	return &CountryService{
		countryRepo: countryRepo,
	}
}

func (s *CountryService) CreateCountry(ctx context.Context, country *models.Country) (*models.Country, error) {
	createdCountry, err := s.countryRepo.Create(ctx, country)
	if err != nil {
		return nil, err
	}

	return createdCountry, nil
}

func (s *CountryService) GetAllCountries(ctx context.Context, page, limit int) ([]*models.Country, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	countrys, err := s.countryRepo.GetAll(ctx, page, limit)
	if err != nil {
		return nil, err
	}

	return countrys, nil
}

func (s *CountryService) GetCountry(ctx context.Context, id uuid.UUID) (*models.Country, error) {
	country, err := s.countryRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return country, nil
}

func (s *CountryService) UpdateCountry(ctx context.Context, country *models.Country) (*models.Country, error) {
	_, err := s.countryRepo.GetByID(ctx, country.ID) // check if exists
	if err != nil {
		return nil, err
	}

	updatedCountry, err := s.countryRepo.Update(ctx, country)
	if err != nil {
		return nil, err
	}

	return updatedCountry, nil
}

func (s *CountryService) DeleteCountry(ctx context.Context, country *models.Country) error {
	return s.countryRepo.Delete(ctx, country.ID)
}

func (s *CountryService) GetTotalCountryCount(ctx context.Context) (int, error) {
	return s.countryRepo.Count(ctx)
}
