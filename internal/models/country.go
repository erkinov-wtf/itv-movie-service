package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type Country struct {
	ID        uuid.UUID      `gorm:"type:uuid;primaryKey"`
	Name      string         `gorm:"column:name;type:text;not null;uniqueIndex"`
	Code      string         `gorm:"column:code;type:varchar(2);not null;uniqueIndex;comment:'ISO 3166-1 alpha-2 code'"`
	Continent string         `gorm:"column:continent;type:text"`
	CreatedAt time.Time      `gorm:"column:created_at"`
	UpdatedAt time.Time      `gorm:"column:updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at"`

	// relations
	Movies []Movie `gorm:"many2many:movie_countries;" json:"movies,omitempty"`
}

func (c *Country) BeforeCreate(tx *gorm.DB) (err error) {
	if c.ID == uuid.Nil {
		c.ID = uuid.New()
	}
	return nil
}
