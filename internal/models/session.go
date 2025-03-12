package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type Session struct {
	ID            uuid.UUID `gorm:"type:uuid;primaryKey"`
	UserID        uuid.UUID `gorm:"type:uuid;not null;index"`
	AccessToken   string    `gorm:"column:access_token;type:text;not null;uniqueIndex"`
	RefreshToken  string    `gorm:"column:refresh_token;type:text;not null;uniqueIndex"`
	ExpiresAt     time.Time `gorm:"column:expires_at;not null"`
	RefreshExpiry time.Time `gorm:"column:refresh_expiry;not null"`
	UserAgent     string    `gorm:"column:user_agent;type:text"`
	IPAddress     string    `gorm:"column:ip_address;type:text"`
	IsRevoked     bool      `gorm:"column:is_revoked;default:false"`
	CreatedAt     time.Time `gorm:"column:created_at"`
	UpdatedAt     time.Time `gorm:"column:updated_at"`

	User User `gorm:"foreignKey:UserID"`
}

func (s *Session) BeforeCreate(*gorm.DB) (err error) {
	if s.ID == uuid.Nil {
		s.ID = uuid.New()
	}
	return nil
}

// IsAccessTokenValid checks if the access token is still valid
func (s *Session) IsAccessTokenValid() bool {
	return !s.IsRevoked && time.Now().Before(s.ExpiresAt)
}

// IsRefreshTokenValid checks if the refresh token is still valid
func (s *Session) IsRefreshTokenValid() bool {
	return !s.IsRevoked && time.Now().Before(s.RefreshExpiry)
}
