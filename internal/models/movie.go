package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type Movie struct {
	ID          uuid.UUID      `gorm:"type:uuid;primaryKey"`
	Title       string         `gorm:"column:title;type:text;not null;index"`
	Director    string         `gorm:"column:director;type:text;index"`
	Year        int            `gorm:"column:year;type:integer;index"`
	Plot        string         `gorm:"column:plot;type:text"`
	Runtime     int            `gorm:"column:runtime;type:integer;comment:'Duration in minutes'"`
	Rating      float64        `gorm:"column:rating;type:decimal(3,1);default:0.0"`
	PosterURL   string         `gorm:"column:poster_url;type:text"`
	TrailerURL  string         `gorm:"column:trailer_url;type:text"`
	ReleaseDate *time.Time     `gorm:"column:release_date;type:date"`
	LanguageID  uuid.UUID      `gorm:"column:language;type:uuid;not null"`
	CreatedAt   time.Time      `gorm:"column:created_at"`
	UpdatedAt   time.Time      `gorm:"column:updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"column:deleted_at"`

	// relations
	Language  Language  `gorm:"foreignKey:LanguageID" json:"language"`
	Countries []Country `gorm:"many2many:movie_countries;" json:"countries"`
	Genres    []Genre   `gorm:"many2many:movie_genres;" json:"genres"`
}

func (m *Movie) BeforeCreate(*gorm.DB) (err error) {
	if m.ID == uuid.Nil {
		m.ID = uuid.New()
	}
	return nil
}
