package handlers

import (
	"github.com/gin-gonic/gin"
	"itv-movie/internal/api/services"
)

// CountryHandler handles HTTP requests for country
type CountryHandler struct {
	countryService *services.CountryService
}

// NewCountryHandler creates a new country handler
func NewCountryHandler(countryService *services.CountryService) *CountryHandler {
	return &CountryHandler{
		countryService: countryService,
	}
}

func (h *CountryHandler) CreateCountry(c *gin.Context) {
	//TODO implement real handler logic
}

func (h *CountryHandler) GetAllCountries(c *gin.Context) {
	//TODO implement real handler logic
}

func (h *CountryHandler) UpdateCountry(c *gin.Context) {
	//TODO implement real handler logic
}

func (h *CountryHandler) DeleteCountry(c *gin.Context) {
	//TODO implement real handler logic
}
