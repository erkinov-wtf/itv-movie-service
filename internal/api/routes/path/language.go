package path

import (
	"github.com/gin-gonic/gin"
	"itv-movie/internal/api/handlers"
	"itv-movie/internal/api/middlewares"
	"itv-movie/internal/api/services"
)

func RegisterLanguageRoutes(r *gin.RouterGroup, handler *handlers.LanguageHandler, authService *services.AuthService) {
	languages := r.Group("/languages")
	{
		languages.GET("", handler.GetAllLanguages)
		languages.GET("/:id", handler.GetLanguage)

		restricted := languages.Group("")
		restricted.Use(middlewares.AuthMiddleware(authService))
		restricted.Use(middlewares.AdminOrDirectorOnly())
		{
			restricted.POST("", handler.CreateLanguage)
			restricted.PUT("/:id", handler.UpdateLanguage)
			restricted.DELETE("/:id", handler.DeleteLanguage)
		}
	}
}
