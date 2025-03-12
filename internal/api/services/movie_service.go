package services

import (
	"context"
	"github.com/google/uuid"
	"itv-movie/internal/models"
	"itv-movie/internal/storage/database"
)

// MovieService handles business logic for movies
type MovieService struct {
	movieRepo    *database.MovieRepository
	genreRepo    *database.GenreRepository
	languageRepo *database.LanguageRepository
	countryRepo  *database.CountryRepository
}

// NewMovieService creates a new movie service
func NewMovieService(
	movieRepo *database.MovieRepository,
	genreRepo *database.GenreRepository,
	languageRepo *database.LanguageRepository,
	countryRepo *database.CountryRepository,
) *MovieService {
	return &MovieService{
		movieRepo:    movieRepo,
		genreRepo:    genreRepo,
		languageRepo: languageRepo,
		countryRepo:  countryRepo,
	}
}

func (s *MovieService) GetMovieByID(ctx context.Context, id uuid.UUID) (*models.Movie, error) {
	//TODO actual implementation
	return nil, nil
}

func (s *MovieService) CreateMovie(ctx context.Context, movie *models.Movie) (*models.Movie, error) {
	//TODO actual implementation
	return nil, nil
}

func (s *MovieService) UpdateMovie(ctx context.Context, movie *models.Movie) (*models.Movie, error) {
	//TODO actual implementation
	return nil, nil
}

func (s *MovieService) DeleteMovie(ctx context.Context, id uuid.UUID) error {
	//TODO actual implementation
	return nil
}

func (s *MovieService) ListMovies(ctx context.Context, page, pageSize int) ([]*models.Movie, int64, error) {
	//TODO actual implementation
	return nil, 0, nil
}
