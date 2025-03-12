package services

import (
	"context"
	"github.com/google/uuid"
	"itv-movie/internal/models"
	"itv-movie/internal/storage/database/repositories"
)

// GenreService handles business logic for genre
type GenreService struct {
	genreRepo *repositories.GenreRepository
}

// NewGenreService creates a new genre service
func NewGenreService(
	genreRepo *repositories.GenreRepository,
) *GenreService {
	return &GenreService{
		genreRepo: genreRepo,
	}
}

func (s *GenreService) CreateGenre(ctx context.Context, genre *models.Genre) (*models.Genre, error) {
	createdGenre, err := s.genreRepo.Create(ctx, genre)
	if err != nil {
		return nil, err
	}

	return createdGenre, nil
}

func (s *GenreService) GetAllGenres(ctx context.Context, page, limit int) ([]*models.Genre, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	genres, err := s.genreRepo.GetAll(ctx, page, limit)
	if err != nil {
		return nil, err
	}

	return genres, nil
}

func (s *GenreService) GetGenre(ctx context.Context, id uuid.UUID) (*models.Genre, error) {
	genre, err := s.genreRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return genre, nil
}

func (s *GenreService) UpdateGenre(ctx context.Context, genre *models.Genre) (*models.Genre, error) {
	_, err := s.genreRepo.GetByID(ctx, genre.ID) // check if exists
	if err != nil {
		return nil, err
	}

	updatedGenre, err := s.genreRepo.Update(ctx, genre)
	if err != nil {
		return nil, err
	}

	return updatedGenre, nil
}

func (s *GenreService) DeleteGenre(ctx context.Context, genre *models.Genre) error {
	return s.genreRepo.Delete(ctx, genre.ID)
}

func (s *GenreService) GetTotalGenreCount(ctx context.Context) (int, error) {
	return s.genreRepo.Count(ctx)
}
