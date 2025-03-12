package database

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"itv-movie/internal/config"
	"log/slog"
)

type PostgresDB struct {
	DB *gorm.DB
}

func MustLoadDB(cfg *config.Config, logger *slog.Logger) (*PostgresDB, error) {
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
		return nil, fmt.Errorf("failed to connect to database: %v", err.Error())
	}

	logger.Info("database connected successfully")

	return &PostgresDB{DB: db}, nil
}

func (p *PostgresDB) Close() error {
	sqlDB, err := p.DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}
