package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type Session struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey"`
	UserID    uuid.UUID `gorm:"type:uuid;not null;index"`
	Token     string    `gorm:"column:token;type:text;not null;uniqueIndex"`
	ExpiresAt time.Time `gorm:"column:expires_at;not null"`
	UserAgent string    `gorm:"column:user_agent;type:text"`
	IPAddress string    `gorm:"column:ip_address;type:text"`
	IsRevoked bool      `gorm:"column:is_revoked;default:false"`
	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
}

func (s *Session) BeforeCreate(*gorm.DB) (err error) {
	if s.ID == uuid.Nil {
		s.ID = uuid.New()
	}
	return nil
}
