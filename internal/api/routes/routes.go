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
		path.RegisterLanguageRoutes(api, languageHandler, authService)
		path.RegisterGenreRoutes(api, genreHandler, authService)
		path.RegisterCountryRoutes(api, countriesHandler, authService)
		path.RegisterMovieRoutes(api, moviesHandler, authService)
		path.RegisterAuthRoutes(api, authHandler, authService)
	}
}
