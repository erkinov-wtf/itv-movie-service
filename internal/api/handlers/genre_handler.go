package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"itv-movie/internal/api/services"
	"itv-movie/internal/models"
	"net/http"
	"strconv"
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
	var body struct {
		Name        string `form:"name" binding:"required"`
		Description string `form:"code" binding:"required"`
	}

	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	newLang := models.Genre{
		Name:        body.Name,
		Description: body.Description,
	}

	createdGenre, err := h.genreService.CreateGenre(c, &newLang)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create genre: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, createdGenre)
}

func (h *GenreHandler) GetAllGenres(c *gin.Context) {
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

	genres, err := h.genreService.GetAllGenres(c, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve genres: " + err.Error()})
		return
	}

	total, err := h.genreService.GetTotalGenreCount(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve genre count: " + err.Error()})
		return
	}

	totalPages := (total + limit - 1) / limit

	response := gin.H{
		"data":  genres,
		"page":  page,
		"pages": totalPages,
		"limit": limit,
	}

	c.JSON(http.StatusOK, response)
}

func (h *GenreHandler) GetGenre(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid genre ID format"})
		return
	}

	genre, err := h.genreService.GetGenre(c, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Genre not found"})
		return
	}

	c.JSON(http.StatusOK, genre)
}

func (h *GenreHandler) UpdateGenre(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid genre ID format"})
		return
	}

	var body struct {
		Name        *string `json:"name,omitempty"`
		Description *string `json:"code,omitempty"`
	}

	if err = c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format: " + err.Error()})
		return
	}

	genre, err := h.genreService.GetGenre(c, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Genre not found"})
		return
	}

	if body.Name != nil {
		genre.Name = *body.Name
	}
	if body.Description != nil {
		genre.Description = *body.Description
	}

	updatedGenre, err := h.genreService.UpdateGenre(c, genre)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update genre: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, updatedGenre)
}

func (h *GenreHandler) DeleteGenre(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid genre ID format"})
		return
	}

	genre, err := h.genreService.GetGenre(c, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Genre not found"})
		return
	}

	if err := h.genreService.DeleteGenre(c, genre); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete genre: " + err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, gin.H{"message": "Genre deleted successfully"})
}
