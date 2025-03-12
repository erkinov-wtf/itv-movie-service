package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type Genre struct {
	ID          uuid.UUID      `gorm:"type:uuid;primaryKey"`
	Name        string         `gorm:"column:name;type:text;not null;uniqueIndex"`
	Description string         `gorm:"column:description;type:text"`
	CreatedAt   time.Time      `gorm:"column:created_at"`
	UpdatedAt   time.Time      `gorm:"column:updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"column:deleted_at"`

	// relations
	Movies []Movie `gorm:"many2many:movie_genres;" json:"movies,omitempty"`
}

func (g *Genre) BeforeCreate(tx *gorm.DB) (err error) {
	if g.ID == uuid.Nil {
		g.ID = uuid.New()
	}
	return nil
}
