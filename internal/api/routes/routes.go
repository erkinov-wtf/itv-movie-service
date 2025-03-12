package routes

import "itv-movie/internal/api/handlers"

func RegisterRoutes(router *Router, movieHandler *handlers.MovieHandler, authHandler *handlers.AuthHandler) {
	api := router.Engine().Group("/api/v1")

	// Movie routes
	movies := api.Group("/movies")
	{
		//TODO add routes
	}
}
