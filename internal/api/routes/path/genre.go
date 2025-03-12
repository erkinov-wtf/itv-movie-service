package path

import (
	"github.com/gin-gonic/gin"
	"itv-movie/internal/api/handlers"
	"itv-movie/internal/api/middlewares"
	"itv-movie/internal/api/services"
)

func RegisterGenreRoutes(r *gin.RouterGroup, handler *handlers.GenreHandler, authService *services.AuthService) {
	genres := r.Group("/genres")
	{
		genres.GET("", handler.GetAllGenres)
		genres.GET("/:id", handler.GetGenre)

		restricted := genres.Group("")
		restricted.Use(middlewares.AuthMiddleware(authService))
		restricted.Use(middlewares.AdminOrDirectorOnly())
		{
			restricted.POST("", handler.CreateGenre)
			restricted.PUT("/:id", handler.UpdateGenre)
			restricted.DELETE("/:id", handler.DeleteGenre)
		}
	}
}
