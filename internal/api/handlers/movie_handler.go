package handlers

import (
	"github.com/gin-gonic/gin"
	"itv-movie/internal/api/services"
)

// MovieHandler handles HTTP requests for movies
type MovieHandler struct {
	movieService *services.MovieService
}

// NewMovieHandler creates a new movie handler
func NewMovieHandler(movieService *services.MovieService) *MovieHandler {
	return &MovieHandler{
		movieService: movieService,
	}
}

func (h *MovieHandler) GetMovie(c *gin.Context) {
	//TODO implement real handler logic
}

func (h *MovieHandler) CreateMovie(c *gin.Context) {
	//TODO implement real handler logic
}

func (h *MovieHandler) UpdateMovie(c *gin.Context) {
	//TODO implement real handler logic
}

func (h *MovieHandler) DeleteMovie(c *gin.Context) {
	//TODO implement real handler logic
}

func (h *MovieHandler) ListMovies(c *gin.Context) {
	//TODO implement real handler logic
}
