package services

import (
	"context"
	"github.com/google/uuid"
	"itv-movie/internal/models"
	"itv-movie/internal/storage/database"
)

// GenreService handles business logic for genre
type GenreService struct {
	genreRepo *database.GenreRepository
}

// NewGenreService creates a new genre service
func NewGenreService(
	genreRepo *database.GenreRepository,
) *GenreService {
	return &GenreService{
		genreRepo: genreRepo,
	}
}

func (s *GenreService) CreateGenre(ctx context.Context, newLang *models.Genre) (*models.Genre, error) {
	//TODO actual implementation
	return nil, nil
}

func (s *GenreService) GetAllCountries(ctx context.Context, page, limit int) ([]*models.Genre, error) {
	//TODO actual implementation
	return nil, nil
}

func (s *GenreService) GetCountry(ctx context.Context, id uuid.UUID) (*models.Genre, error) {
	//TODO actual implementation
	return nil, nil
}

func (s *GenreService) UpdateGenre(ctx context.Context, lang *models.Genre) (*models.Genre, error) {
	//TODO actual implementation
	return nil, nil
}

func (s *GenreService) DeleteGenre(ctx context.Context, newLang *models.Genre) error {
	//TODO actual implementation
	return nil
}
