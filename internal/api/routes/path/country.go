package path

import (
	"github.com/gin-gonic/gin"
	"itv-movie/internal/api/handlers"
	"itv-movie/internal/api/middlewares"
	"itv-movie/internal/api/services"
)

func RegisterCountryRoutes(r *gin.RouterGroup, handler *handlers.CountryHandler, authService *services.AuthService) {
	countries := r.Group("/countries")
	{
		countries.GET("", handler.GetAllCountries)
		countries.GET("/:id", handler.GetCountry)

		restricted := countries.Group("")
		restricted.Use(middlewares.AuthMiddleware(authService))
		restricted.Use(middlewares.AdminOrDirectorOnly())
		{
			restricted.POST("", handler.CreateCountry)
			restricted.PUT("/:id", handler.UpdateCountry)
			restricted.DELETE("/:id", handler.DeleteCountry)
		}
	}
}
