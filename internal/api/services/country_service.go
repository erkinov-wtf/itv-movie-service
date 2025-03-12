package services

import (
	"context"
	"itv-movie/internal/models"
	"itv-movie/internal/storage/database"
)

// CountryService handles business logic for country
type CountryService struct {
	countryRepo *database.CountryRepository
}

// NewCountryService creates a new lang country
func NewCountryService(
	countryRepo *database.CountryRepository,
) *CountryService {
	return &CountryService{
		countryRepo: countryRepo,
	}
}

func (s *CountryService) CreateCountry(ctx context.Context, newLang *models.Country) (*models.Country, error) {
	//TODO actual implementation
	return nil, nil
}

func (s *CountryService) GetAllCountries(ctx context.Context, page, limit int) ([]*models.Country, error) {
	//TODO actual implementation
	return nil, nil
}

func (s *CountryService) UpdateCountry(ctx context.Context, lang *models.Country) (*models.Country, error) {
	//TODO actual implementation
	return nil, nil
}

func (s *CountryService) DeleteCountry(ctx context.Context, newLang *models.Country) error {
	//TODO actual implementation
	return nil
}
