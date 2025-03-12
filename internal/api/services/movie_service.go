package services

import (
	"context"
	"github.com/google/uuid"
	"itv-movie/internal/models"
	"itv-movie/internal/storage/database/repositories"
)

// MovieService handles business logic for movies
type MovieService struct {
	movieRepo    *repositories.MovieRepository
	languageRepo *repositories.LanguageRepository
	countryRepo  *repositories.CountryRepository
	genreRepo    *repositories.GenreRepository
}

// NewMovieService creates a new movie service
func NewMovieService(
	movieRepo *repositories.MovieRepository,
	languageRepo *repositories.LanguageRepository,
	countryRepo *repositories.CountryRepository,
	genreRepo *repositories.GenreRepository,
) *MovieService {
	return &MovieService{
		movieRepo:    movieRepo,
		languageRepo: languageRepo,
		countryRepo:  countryRepo,
		genreRepo:    genreRepo,
	}
}

func (s *MovieService) CreateMovie(ctx context.Context, movie *models.Movie) (*models.Movie, error) {
	createdMovie, err := s.movieRepo.Create(ctx, movie)
	if err != nil {
		return nil, err
	}

	return createdMovie, nil
}

func (s *MovieService) GetAllMovies(ctx context.Context, page, limit int) ([]*models.Movie, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10 // Default limit
	}

	movies, err := s.movieRepo.GetAll(ctx, page, limit)
	if err != nil {
		return nil, err
	}

	return movies, nil
}

func (s *MovieService) GetMovie(ctx context.Context, id uuid.UUID) (*models.Movie, error) {
	movie, err := s.movieRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return movie, nil
}

func (s *MovieService) UpdateMovie(ctx context.Context, movie *models.Movie) (*models.Movie, error) {
	_, err := s.movieRepo.GetByID(ctx, movie.ID)
	if err != nil {
		return nil, err
	}

	updatedMovie, err := s.movieRepo.Update(ctx, movie)
	if err != nil {
		return nil, err
	}

	return updatedMovie, nil
}

func (s *MovieService) DeleteMovie(ctx context.Context, id uuid.UUID) error {
	return s.movieRepo.Delete(ctx, id)
}

func (s *MovieService) GetTotalMovieCount(ctx context.Context) (int, error) {
	return s.movieRepo.Count(ctx)
}

func (s *MovieService) SearchMovies(ctx context.Context, query string, page, limit int) ([]*models.Movie, int, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10 // Default limit
	}

	movies, err := s.movieRepo.Search(ctx, query, page, limit)
	if err != nil {
		return nil, 0, err
	}

	total, err := s.movieRepo.SearchCount(ctx, query)
	if err != nil {
		return nil, 0, err
	}

	return movies, total, nil
}

func (s *MovieService) GetLangByCode(ctx context.Context, code string) (*models.Language, error) {
	return s.languageRepo.GetByCode(ctx, code)
}

func (s *MovieService) GetGenreByName(ctx context.Context, name string) (*models.Genre, error) {
	return s.genreRepo.GetByName(ctx, name)
}

func (s *MovieService) GetCountryByCode(ctx context.Context, code string) (*models.Country, error) {
	return s.countryRepo.GetByCode(ctx, code)
}
