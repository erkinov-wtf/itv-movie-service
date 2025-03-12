package path

import (
	"github.com/gin-gonic/gin"
	"itv-movie/internal/api/handlers"
	"itv-movie/internal/api/middlewares"
	"itv-movie/internal/api/services"
)

func RegisterMovieRoutes(r *gin.RouterGroup, handler *handlers.MovieHandler, authService *services.AuthService) {
	movies := r.Group("/movies")
	{
		movies.GET("", handler.GetAllMovies)
		movies.GET("/:id", handler.GetMovie)

		restricted := movies.Group("")
		restricted.Use(middlewares.AuthMiddleware(authService))
		restricted.Use(middlewares.AdminOrDirectorOnly())
		{
			restricted.POST("", handler.CreateMovie)
			restricted.PUT("/:id", handler.UpdateMovie)
			restricted.DELETE("/:id", handler.DeleteMovie)
		}
	}
}
