package path

import (
	"github.com/gin-gonic/gin"
	"itv-movie/internal/api/handlers"
)

func RegisterGenreRoutes(r *gin.RouterGroup, handler *handlers.GenreHandler) {
	genres := r.Group("/genres")
	{
		genres.POST("", handler.CreateGenre)
		genres.GET("", handler.GetAllGenres)
		genres.GET("/:id", handler.GetGenre)
		genres.PUT("/:id", handler.UpdateGenre)
		genres.DELETE("/:id", handler.DeleteGenre)
	}
}
