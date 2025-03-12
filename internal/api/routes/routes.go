package routes

import (
	"itv-movie/internal/api/handlers"
	"itv-movie/internal/api/routes/path"
	"itv-movie/internal/api/services"
)

func RegisterRoutes(router *Router,
	languageHandler *handlers.LanguageHandler,
	genreHandler *handlers.GenreHandler,
	countriesHandler *handlers.CountryHandler,
	moviesHandler *handlers.MovieHandler,
	authHandler *handlers.AuthHandler,
	authService *services.AuthService,
) {
	api := router.Engine().Group("/api/v1")
	{
		path.RegisterLanguageRoutes(api, languageHandler)
		path.RegisterGenreRoutes(api, genreHandler)
		path.RegisterCountryRoutes(api, countriesHandler)
		path.RegisterMovieRoutes(api, moviesHandler)
		path.RegisterAuthRoutes(api, authHandler, authService)
	}
}
