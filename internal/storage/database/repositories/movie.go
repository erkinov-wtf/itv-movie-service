package repositories

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"itv-movie/internal/models"
	"itv-movie/internal/storage/database"
	"strings"
)

// MovieRepository handles database operations for movies
type MovieRepository struct {
	db *gorm.DB
}

// NewMovieRepository creates a new movie repository
func NewMovieRepository(postgres *database.PostgresDB) *MovieRepository {
	return &MovieRepository{
		db: postgres.DB,
	}
}

func (r *MovieRepository) Create(ctx context.Context, movie *models.Movie) (*models.Movie, error) {
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(movie).Error; err != nil {
			return err
		}

		if len(movie.Genres) > 0 {
			// Build values string for PostgreSQL batch insert
			var genreValueStrings []string
			var genreValueArgs []interface{}

			for i, genre := range movie.Genres {
				genreValueStrings = append(genreValueStrings, fmt.Sprintf("($%d, $%d)", i*2+1, i*2+2))
				genreValueArgs = append(genreValueArgs, movie.ID, genre.ID)
			}

			genreQuery := fmt.Sprintf(
				"INSERT INTO movie_genres (movie_id, genre_id) VALUES %s ON CONFLICT DO NOTHING",
				strings.Join(genreValueStrings, ","),
			)

			if err := tx.Exec(genreQuery, genreValueArgs...).Error; err != nil {
				return err
			}
		}

		if len(movie.Countries) > 0 {
			var countryValueStrings []string
			var countryValueArgs []interface{}

			for i, country := range movie.Countries {
				countryValueStrings = append(countryValueStrings, fmt.Sprintf("($%d, $%d)", i*2+1, i*2+2))
				countryValueArgs = append(countryValueArgs, movie.ID, country.ID)
			}

			countryQuery := fmt.Sprintf(
				"INSERT INTO movie_countries (movie_id, country_id) VALUES %s ON CONFLICT DO NOTHING",
				strings.Join(countryValueStrings, ","),
			)

			if err := tx.Exec(countryQuery, countryValueArgs...).Error; err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return r.GetByID(ctx, movie.ID)
}

func (r *MovieRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Movie, error) {
	var movie models.Movie

	if err := r.db.WithContext(ctx).
		Preload("Language").
		Preload("Countries").
		Preload("Genres").
		Where("id = ?", id).
		First(&movie).Error; err != nil {
		return nil, err
	}

	return &movie, nil
}

func (r *MovieRepository) GetAll(ctx context.Context, page, limit int) ([]*models.Movie, error) {
	var movies []*models.Movie
	offset := (page - 1) * limit

	if err := r.db.WithContext(ctx).
		Preload("Language").
		Preload("Countries").
		Preload("Genres").
		Offset(offset).
		Limit(limit).
		Find(&movies).Error; err != nil {
		return nil, err
	}

	return movies, nil
}

func (r *MovieRepository) Update(ctx context.Context, movie *models.Movie) (*models.Movie, error) {
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Update the movie's basic fields
		if err := tx.Model(movie).Updates(map[string]interface{}{
			"title":        movie.Title,
			"director":     movie.Director,
			"year":         movie.Year,
			"plot":         movie.Plot,
			"runtime":      movie.Runtime,
			"rating":       movie.Rating,
			"poster_url":   movie.PosterURL,
			"trailer_url":  movie.TrailerURL,
			"release_date": movie.ReleaseDate,
			"language":     movie.LanguageID, // Using "language" for the column name
		}).Error; err != nil {
			return err
		}

		// Clear existing genre relationships
		if err := tx.Exec("DELETE FROM movie_genres WHERE movie_id = ?", movie.ID).Error; err != nil {
			return err
		}

		// Add new genre relationships if provided
		if len(movie.Genres) > 0 {
			var genreValueStrings []string
			var genreValueArgs []interface{}

			for i, genre := range movie.Genres {
				genreValueStrings = append(genreValueStrings, fmt.Sprintf("($%d, $%d)", i*2+1, i*2+2))
				genreValueArgs = append(genreValueArgs, movie.ID, genre.ID)
			}

			genreQuery := fmt.Sprintf(
				"INSERT INTO movie_genres (movie_id, genre_id) VALUES %s ON CONFLICT DO NOTHING",
				strings.Join(genreValueStrings, ","),
			)

			if err := tx.Exec(genreQuery, genreValueArgs...).Error; err != nil {
				return err
			}
		}

		// Clear existing country relationships
		if err := tx.Exec("DELETE FROM movie_countries WHERE movie_id = ?", movie.ID).Error; err != nil {
			return err
		}

		// Add new country relationships if provided
		if len(movie.Countries) > 0 {
			var countryValueStrings []string
			var countryValueArgs []interface{}

			for i, country := range movie.Countries {
				countryValueStrings = append(countryValueStrings, fmt.Sprintf("($%d, $%d)", i*2+1, i*2+2))
				countryValueArgs = append(countryValueArgs, movie.ID, country.ID)
			}

			countryQuery := fmt.Sprintf(
				"INSERT INTO movie_countries (movie_id, country_id) VALUES %s ON CONFLICT DO NOTHING",
				strings.Join(countryValueStrings, ","),
			)

			if err := tx.Exec(countryQuery, countryValueArgs...).Error; err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return r.GetByID(ctx, movie.ID)
}

func (r *MovieRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&models.Movie{}, id).Error
}

func (r *MovieRepository) Count(ctx context.Context) (int, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&models.Movie{}).Count(&count).Error; err != nil {
		return 0, err
	}
	return int(count), nil
}

func (r *MovieRepository) Search(ctx context.Context, query string, page, limit int) ([]*models.Movie, error) {
	var movies []*models.Movie
	offset := (page - 1) * limit

	db := r.db.WithContext(ctx).
		Preload("Language").
		Preload("Countries").
		Preload("Genres")

	if query != "" {
		db = db.Where("title ILIKE ?", "%"+query+"%").
			Or("plot ILIKE ?", "%"+query+"%").
			Or("director ILIKE ?", "%"+query+"%")
	}

	if err := db.Offset(offset).Limit(limit).Find(&movies).Error; err != nil {
		return nil, err
	}

	return movies, nil
}

func (r *MovieRepository) SearchCount(ctx context.Context, query string) (int, error) {
	var count int64
	db := r.db.WithContext(ctx).Model(&models.Movie{})

	if query != "" {
		db = db.Where("title ILIKE ?", "%"+query+"%").
			Or("director ILIKE ?", "%"+query+"%")
	}

	if err := db.Count(&count).Error; err != nil {
		return 0, err
	}

	return int(count), nil
}
