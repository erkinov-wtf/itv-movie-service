package routes

import "itv-movie/internal/api/handlers"

func RegisterRoutes(router *Router, movieHandler *handlers.MovieHandler, authHandler *handlers.AuthHandler) {
	api := router.Engine().Group("/api/v1")

	// Movie routes
	movies := api.Group("/movies")
	{
		movies.GET("", movieHandler.ListMovies)
		movies.POST("", movieHandler.CreateMovie)
		movies.GET("/:id", movieHandler.GetMovie)
		movies.PUT("/:id", movieHandler.UpdateMovie)
		movies.DELETE("/:id", movieHandler.DeleteMovie)
	}

	// Auth routes
	auth := api.Group("/auth")
	{
		auth.POST("/register", authHandler.Register)
		auth.POST("/login", authHandler.Login)
		auth.POST("/logout", authHandler.Logout)
	}
}
