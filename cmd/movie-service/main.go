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
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	app := fx.New(
		// Providing all dependencies
		fx.Provide(
			config.MustLoad,
			logger.SetupLogger,
			database.MustLoadDB,

			// Repositories
			database.NewMovieRepository,
			database.NewGenreRepository,
			database.NewCountryRepository,
			database.NewLanguageRepository,
			database.NewUserRepository,

			// Services
			services.NewMovieService,
			services.NewAuthService,

			// Handlers setup
			handlers.NewMovieHandler,
			handlers.NewAuthHandler,
			//TODO add others

			// Router
			routes.NewRouter,
		),

		// Register hooks
		fx.Invoke(routes.RegisterRoutes),

		// Lifecycle hooks
		fx.Invoke(registerHooks),
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
