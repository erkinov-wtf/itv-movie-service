package routes

import (
	"itv-movie/internal/api/handlers"
	"itv-movie/internal/api/routes/path"
)

func RegisterRoutes(router *Router,
	languageHandler *handlers.LanguageHandler, genreHandler *handlers.GenreHandler, countriesHandler *handlers.CountryHandler, moviesHandler *handlers.MovieHandler,
) {
	api := router.Engine().Group("/api/v1")
	{
		path.RegisterLanguageRoutes(api, languageHandler)
		path.RegisterGenreRoutes(api, genreHandler)
		path.RegisterCountryRoutes(api, countriesHandler)
		path.RegisterMoviesRoutes(api, moviesHandler)
	}
}
