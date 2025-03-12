package database

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"itv-movie/internal/config"
	"log"
	"log/slog"
)

func MustLoadDB(cfg *config.Config, logger *slog.Logger) *gorm.DB {
	var err error
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%v search_path=%s TimeZone=%s sslmode=disable",
		cfg.Internal.Database.Host,
		cfg.Internal.Database.User,
		cfg.Internal.Database.Password,
		cfg.Internal.Database.Name,
		cfg.Internal.Database.Port,
		cfg.Internal.Database.Schema,
		cfg.Internal.Database.Timezone,
	)

	var db *gorm.DB

	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err.Error())
	}

	logger.Info("database connected successfully")

	return db
}
