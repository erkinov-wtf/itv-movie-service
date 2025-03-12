package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"itv-movie/internal/api/services"
	"itv-movie/internal/models"
	"net/http"
	"strconv"
)

// LanguageHandler handles HTTP requests for Language
type LanguageHandler struct {
	languageService *services.LanguageService
}

// NewLanguageHandler creates a new Language handler
func NewLanguageHandler(languageService *services.LanguageService) *LanguageHandler {
	return &LanguageHandler{
		languageService: languageService,
	}
}

func (h *LanguageHandler) CreateLanguage(c *gin.Context) {
	var body struct {
		Name string `form:"name" binding:"required"`
		Code string `form:"code" binding:"required"`
	}

	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	newLang := models.Language{
		Name: body.Name,
		Code: body.Code,
	}

	createdLanguage, err := h.languageService.CreateLanguage(c, &newLang)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create language: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, createdLanguage)
}

func (h *LanguageHandler) GetAllLanguages(c *gin.Context) {
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

	languages, err := h.languageService.GetAllLanguages(c, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve languages: " + err.Error()})
		return
	}

	total, err := h.languageService.GetTotalLanguageCount(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve language count: " + err.Error()})
		return
	}

	totalPages := (total + limit - 1) / limit

	response := gin.H{
		"data":  languages,
		"page":  page,
		"pages": totalPages,
		"limit": limit,
	}

	c.JSON(http.StatusOK, response)
}

func (h *LanguageHandler) GetLanguage(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid language ID format"})
		return
	}

	language, err := h.languageService.GetLanguage(c, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Language not found"})
		return
	}

	c.JSON(http.StatusOK, language)
}

func (h *LanguageHandler) UpdateLanguage(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid language ID format"})
		return
	}

	var body struct {
		Name *string `json:"name,omitempty"`
		Code *string `json:"code,omitempty"`
	}

	if err = c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format: " + err.Error()})
		return
	}

	language, err := h.languageService.GetLanguage(c, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Language not found"})
		return
	}

	if body.Name != nil {
		language.Name = *body.Name
	}
	if body.Code != nil {
		language.Code = *body.Code
	}

	updatedLanguage, err := h.languageService.UpdateLanguage(c, language)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update language: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, updatedLanguage)
}

func (h *LanguageHandler) DeleteLanguage(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid language ID format"})
		return
	}

	language, err := h.languageService.GetLanguage(c, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Language not found"})
		return
	}

	if err := h.languageService.DeleteLanguage(c, language); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete language: " + err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, gin.H{"message": "Language deleted successfully"})
}
