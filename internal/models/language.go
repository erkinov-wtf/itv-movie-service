package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type Language struct {
	ID        uuid.UUID      `gorm:"type:uuid;primaryKey"`
	Name      string         `gorm:"column:name;type:text;not null;index"`
	Code      string         `gorm:"column:code;type:text;not null;uniqueIndex"`
	CreatedAt time.Time      `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time      `gorm:"column:updated_at;autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;index"`
}

func (l *Language) BeforeCreate(tx *gorm.DB) (err error) {
	if l.ID == uuid.Nil {
		l.ID = uuid.New()
	}
	return nil
}
