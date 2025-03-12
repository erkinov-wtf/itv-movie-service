package handlers

import (
	"fmt"
	"itv-movie/internal/pkg/utils/constants"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"itv-movie/internal/api/services"
	"itv-movie/internal/models"
)

// MovieHandler handles HTTP requests for Movies
type MovieHandler struct {
	movieService *services.MovieService
}

// NewMovieHandler creates a new Movie handler
func NewMovieHandler(movieService *services.MovieService) *MovieHandler {
	return &MovieHandler{
		movieService: movieService,
	}
}

func (h *MovieHandler) CreateMovie(c *gin.Context) {
	var body struct {
		Title       string   `json:"title" binding:"required"`
		Director    string   `json:"director" binding:"required"`
		Year        int      `json:"year" binding:"required"`
		Plot        string   `json:"plot" binding:"required"`
		Runtime     int      `json:"runtime" binding:"required"`
		Rating      *float32 `json:"rating" binding:"omitempty"`
		PosterUrl   string   `json:"posterUrl" binding:"required"`
		TrailerUrl  string   `json:"trailerUrl" binding:"required"`
		ReleaseDate string   `json:"releaseDate" binding:"required"`
		Language    string   `json:"language" binding:"required"`
		Genres      []string `json:"genres" binding:"omitempty"`
		Countries   []string `json:"countries" binding:"omitempty"`
	}

	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format: " + err.Error()})
		return
	}

	releaseDate, err := time.Parse(constants.DateFormat, body.ReleaseDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid release date format. Use YYYY-MM-DD"})
		return
	}

	language, err := h.movieService.GetLangByCode(c, body.Language)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Language with code '%s' not found", body.Language)})
		return
	}

	newMovie := &models.Movie{
		Title:       body.Title,
		Director:    body.Director,
		Year:        body.Year,
		Plot:        body.Plot,
		Runtime:     body.Runtime,
		PosterURL:   body.PosterUrl,
		TrailerURL:  body.TrailerUrl,
		ReleaseDate: &releaseDate,
		LanguageID:  language.ID,
	}

	if body.Rating != nil {
		newMovie.Rating = *body.Rating
	}

	if len(body.Genres) > 0 {
		genreList := make([]models.Genre, 0, len(body.Genres))
		for _, genreName := range body.Genres {
			genre, err := h.movieService.GetGenreByName(c, genreName)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Genre '%s' not found", genreName)})
				return
			}
			genreList = append(genreList, *genre)
		}
		newMovie.Genres = genreList
	}

	if len(body.Countries) > 0 {
		countryList := make([]models.Country, 0, len(body.Countries))
		for _, countryCode := range body.Countries {
			country, err := h.movieService.GetCountryByCode(c, countryCode)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Country with code '%s' not found", countryCode)})
				return
			}
			countryList = append(countryList, *country)
		}
		newMovie.Countries = countryList
	}

	createdMovie, err := h.movieService.CreateMovie(c, newMovie)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create movie: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, createdMovie)
}

func (h *MovieHandler) GetAllMovies(c *gin.Context) {
	pageStr := c.Query("page")
	limitStr := c.Query("limit")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		limit = 10
	}

	query := c.Query("search")

	var movies []*models.Movie
	var total int

	if query != "" {
		movies, total, err = h.movieService.SearchMovies(c, query, page, limit)
	} else {
		movies, err = h.movieService.GetAllMovies(c, page, limit)
		if err == nil {
			total, err = h.movieService.GetTotalMovieCount(c)
		}
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve movies: " + err.Error()})
		return
	}

	totalPages := (total + limit - 1) / limit

	response := gin.H{
		"data":  movies,
		"page":  page,
		"pages": totalPages,
		"limit": limit,
	}

	c.JSON(http.StatusOK, response)
}

func (h *MovieHandler) GetMovie(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid movie ID format"})
		return
	}

	movie, err := h.movieService.GetMovie(c, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Movie not found"})
		return
	}

	c.JSON(http.StatusOK, movie)
}

// UpdateMovie updates an existing movie
func (h *MovieHandler) UpdateMovie(c *gin.Context) {
	// Parse ID from URL parameter
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid movie ID format"})
		return
	}

	// Get existing movie
	movie, err := h.movieService.GetMovie(c, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Movie not found"})
		return
	}

	// Define structure for partial updates
	var update struct {
		Title       *string  `json:"title,omitempty"`
		Director    *string  `json:"director,omitempty"`
		Year        *int     `json:"year,omitempty"`
		Plot        *string  `json:"plot,omitempty"`
		Runtime     *int     `json:"runtime,omitempty"`
		Rating      *float32 `json:"rating,omitempty"`
		PosterUrl   *string  `json:"posterUrl,omitempty"`
		TrailerUrl  *string  `json:"trailerUrl,omitempty"`
		ReleaseDate *string  `json:"releaseDate,omitempty"`
		Language    *string  `json:"language,omitempty"`
		Genres      []string `json:"genres,omitempty"`
		Countries   []string `json:"countries,omitempty"`
	}

	if err := c.ShouldBindJSON(&update); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format: " + err.Error()})
		return
	}

	if update.Title != nil {
		movie.Title = *update.Title
	}

	if update.Director != nil {
		movie.Director = *update.Director
	}

	if update.Year != nil {
		movie.Year = *update.Year
	}

	if update.Plot != nil {
		movie.Plot = *update.Plot
	}

	if update.Runtime != nil {
		movie.Runtime = *update.Runtime
	}

	if update.Rating != nil {
		movie.Rating = *update.Rating
	}

	if update.PosterUrl != nil {
		movie.PosterURL = *update.PosterUrl
	}

	if update.TrailerUrl != nil {
		movie.TrailerURL = *update.TrailerUrl
	}

	if update.ReleaseDate != nil {
		// Parse the date string into time.Time
		releaseDate, err := time.Parse(constants.DateFormat, *update.ReleaseDate)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid release date format. Use YYYY-MM-DD"})
			return
		}
		movie.ReleaseDate = &releaseDate
	}

	// Update language if provided
	if update.Language != nil {
		language, err := h.movieService.GetLangByCode(c, *update.Language)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Language with code '%s' not found", *update.Language)})
			return
		}
		movie.LanguageID = language.ID
	}

	// Update genres if provided
	// This uses a field presence check rather than nil check since it's a slice
	if update.Genres != nil {
		genreList := make([]models.Genre, 0, len(update.Genres))
		for _, genreName := range update.Genres {
			genre, err := h.movieService.GetGenreByName(c, genreName)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Genre '%s' not found", genreName)})
				return
			}
			genreList = append(genreList, *genre)
		}
		movie.Genres = genreList
	}

	// Update countries if provided
	if update.Countries != nil {
		countryList := make([]models.Country, 0, len(update.Countries))
		for _, countryCode := range update.Countries {
			country, err := h.movieService.GetCountryByCode(c, countryCode)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Country with code '%s' not found", countryCode)})
				return
			}
			countryList = append(countryList, *country)
		}
		movie.Countries = countryList
	}

	updatedMovie, err := h.movieService.UpdateMovie(c, movie)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update movie: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, updatedMovie)
}

func (h *MovieHandler) DeleteMovie(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid movie ID format"})
		return
	}

	if err = h.movieService.DeleteMovie(c, id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete movie: " + err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, gin.H{"message": "Movie deleted successfully"})
}
