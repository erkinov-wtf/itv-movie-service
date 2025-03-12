package database

import "gorm.io/gorm"

type MovieRepository struct {
	db *gorm.DB
}

func NewMovieRepository(postgres *PostgresDB) *MovieRepository {
	return &MovieRepository{db: postgres.DB}
}

type CountryRepository struct {
	db *gorm.DB
}

func NewCountryRepository(postgres *PostgresDB) *CountryRepository {
	return &CountryRepository{db: postgres.DB}
}

type GenreRepository struct {
	db *gorm.DB
}

func NewGenreRepository(postgres *PostgresDB) *GenreRepository {
	return &GenreRepository{db: postgres.DB}
}

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(postgres *PostgresDB) *UserRepository {
	return &UserRepository{db: postgres.DB}
}
