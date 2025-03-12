package handlers

import (
	"github.com/gin-gonic/gin"
	"itv-movie/internal/api/services"
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
	//TODO implement real handler logic
}

func (h *LanguageHandler) GetAllLanguages(c *gin.Context) {
	//TODO implement real handler logic
}

func (h *LanguageHandler) UpdateLanguage(c *gin.Context) {
	//TODO implement real handler logic
}

func (h *LanguageHandler) DeleteLanguage(c *gin.Context) {
	//TODO implement real handler logic
}
