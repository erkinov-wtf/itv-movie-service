package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"itv-movie/internal/api/services"
	"itv-movie/internal/models"
	"net/http"
	"strconv"
)

// CountryHandler handles HTTP requests for Country
type CountryHandler struct {
	countryService *services.CountryService
}

// NewCountryHandler creates a new Country handler
func NewCountryHandler(countryService *services.CountryService) *CountryHandler {
	return &CountryHandler{
		countryService: countryService,
	}
}

func (h *CountryHandler) CreateCountry(c *gin.Context) {
	var body struct {
		Name      string `form:"name" binding:"required"`
		Code      string `form:"code" binding:"required"`
		Continent string `form:"content" binding:"required"`
	}

	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	newLang := models.Country{
		Name:      body.Name,
		Code:      body.Code,
		Continent: body.Continent,
	}

	createdCountry, err := h.countryService.CreateCountry(c, &newLang)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create country: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, createdCountry)
}

func (h *CountryHandler) GetAllCountries(c *gin.Context) {
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

	countrys, err := h.countryService.GetAllCountrys(c, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve countrys: " + err.Error()})
		return
	}

	total, err := h.countryService.GetTotalCountryCount(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve country count: " + err.Error()})
		return
	}

	totalPages := (total + limit - 1) / limit

	response := gin.H{
		"data":  countrys,
		"page":  page,
		"pages": totalPages,
		"limit": limit,
	}

	c.JSON(http.StatusOK, response)
}

func (h *CountryHandler) GetCountry(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid country ID format"})
		return
	}

	country, err := h.countryService.GetCountry(c, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Country not found"})
		return
	}

	c.JSON(http.StatusOK, country)
}

func (h *CountryHandler) UpdateCountry(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid country ID format"})
		return
	}

	var body struct {
		Name      *string `json:"name,omitempty"`
		Code      *string `json:"code,omitempty"`
		Continent *string `json:"continent,omitempty"`
	}

	if err = c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format: " + err.Error()})
		return
	}

	country, err := h.countryService.GetCountry(c, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Country not found"})
		return
	}

	if body.Name != nil {
		country.Name = *body.Name
	}
	if body.Code != nil {
		country.Code = *body.Code
	}
	if body.Continent != nil {
		country.Continent = *body.Continent
	}

	updatedCountry, err := h.countryService.UpdateCountry(c, country)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update country: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, updatedCountry)
}

func (h *CountryHandler) DeleteCountry(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid country ID format"})
		return
	}

	country, err := h.countryService.GetCountry(c, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Country not found"})
		return
	}

	if err := h.countryService.DeleteCountry(c, country); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete country: " + err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, gin.H{"message": "Country deleted successfully"})
}
