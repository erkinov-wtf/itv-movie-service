package handlers

import (
	"github.com/gin-gonic/gin"
	"itv-movie/internal/api/services"
)

// GenreHandler handles HTTP requests for Genre
type GenreHandler struct {
	genreService *services.GenreService
}

// NewGenreHandler creates a new Genre handler
func NewGenreHandler(genreService *services.GenreService) *GenreHandler {
	return &GenreHandler{
		genreService: genreService,
	}
}

func (h *GenreHandler) CreateGenre(c *gin.Context) {
	//TODO implement real handler logic
}

func (h *GenreHandler) GetAllCountries(c *gin.Context) {
	//TODO implement real handler logic
}

func (h *GenreHandler) UpdateGenre(c *gin.Context) {
	//TODO implement real handler logic
}

func (h *GenreHandler) DeleteGenre(c *gin.Context) {
	//TODO implement real handler logic
}
