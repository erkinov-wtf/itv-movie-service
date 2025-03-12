package path

import (
	"github.com/gin-gonic/gin"
	"itv-movie/internal/api/handlers"
)

func RegisterLanguageRoutes(r *gin.RouterGroup, handler *handlers.LanguageHandler) {
	languages := r.Group("/languages")
	{
		languages.POST("", handler.CreateLanguage)
		languages.GET("", handler.GetAllLanguages)
		languages.GET("/:id", handler.GetLanguage)
		languages.PUT("/:id", handler.UpdateLanguage)
		languages.DELETE("/:id", handler.DeleteLanguage)
	}
}
