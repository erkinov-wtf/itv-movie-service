// cmd/movie-service/main.go
package main

import (
	"context"
	"go.uber.org/fx"
	"itv-movie/internal/api/handlers"
	"itv-movie/internal/api/routes"
	"itv-movie/internal/api/services"
	"itv-movie/internal/config"
	"itv-movie/internal/pkg/utils/logger"
	"itv-movie/internal/storage/database"
	"itv-movie/internal/storage/database/repositories"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

func provideLoggerEnv(cfg *config.Config) string {
	env := config.LocalEnv
	if cfg.Env != "" {
		env = cfg.Env
	}
	return env
}

func main() {
	app := fx.New(
		// Providing all dependencies
		fx.Provide(
			config.MustLoad,
			provideLoggerEnv,
			logger.SetupLogger,
			database.MustLoadDB,

			// Repositories
			repositories.NewLanguageRepository,
			repositories.NewGenreRepository,
			repositories.NewCountryRepository,
			repositories.NewMovieRepository,

			// Services
			services.NewLanguageService,
			services.NewGenreService,
			services.NewCountryService,
			services.NewMovieService,

			// Handlers setup
			handlers.NewLanguageHandler,
			handlers.NewGenreHandler,
			handlers.NewCountryHandler,
			handlers.NewMovieHandler,
			//TODO add others

			// Router
			routes.NewRouter,
		),

		// Register hooks
		fx.Invoke(routes.RegisterRoutes),

		// Lifecycle hooks
		fx.Invoke(registerHooks),
		fx.Invoke(startHTTPServer),
	)

	startCtx, cancel := context.WithTimeout(context.Background(), config.DefaultTimeout)
	defer cancel()

	if err := app.Start(startCtx); err != nil {
		log.Fatal(err)
	}

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	// Graceful shutdown
	stopCtx, cancel := context.WithTimeout(context.Background(), config.DefaultTimeout)
	defer cancel()

	if err := app.Stop(stopCtx); err != nil {
		log.Fatal(err)
	}
}

func startHTTPServer(lc fx.Lifecycle, router *routes.Router, log *slog.Logger) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			log.Info("Starting HTTP server")
			// Start the server in a goroutine so it doesn't block
			go func() {
				if err := router.Run(); err != nil {
					log.Error("Failed to start HTTP server", "error", err)
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			log.Info("Stopping HTTP server")
			// The HTTP server will stop when the application stops
			return nil
		},
	})
}

func registerHooks(lc fx.Lifecycle, db *database.PostgresDB, cfg *config.Config, log *slog.Logger) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			log.Info("Starting movie service")
			return nil
		},
		OnStop: func(ctx context.Context) error {
			log.Info("Shutting down movie service")
			return db.Close()
		},
	})
}
