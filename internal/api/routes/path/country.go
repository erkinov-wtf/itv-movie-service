package path

import (
	"github.com/gin-gonic/gin"
	"itv-movie/internal/api/handlers"
)

func RegisterCountryRoutes(r *gin.RouterGroup, handler *handlers.CountryHandler) {
	genres := r.Group("/countries")
	{
		genres.POST("", handler.CreateCountry)
		genres.GET("", handler.GetAllCountries)
		genres.GET("/:id", handler.GetCountry)
		genres.PUT("/:id", handler.UpdateCountry)
		genres.DELETE("/:id", handler.DeleteCountry)
	}
}
