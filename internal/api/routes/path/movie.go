package path

import (
	"github.com/gin-gonic/gin"
	"itv-movie/internal/api/handlers"
)

func RegisterMovieRoutes(r *gin.RouterGroup, handler *handlers.MovieHandler) {
	movies := r.Group("/movies")
	{
		movies.POST("", handler.CreateMovie)
		movies.GET("", handler.GetAllMovies)
		movies.GET("/:id", handler.GetMovie)
		movies.PUT("/:id", handler.UpdateMovie)
		movies.DELETE("/:id", handler.DeleteMovie)
	}
}
