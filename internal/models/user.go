package models

import (
	"fmt"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"itv-movie/internal/config"
	"itv-movie/internal/pkg/jwt"
	"time"
)

type User struct {
	ID          uuid.UUID      `gorm:"type:uuid;primaryKey"`
	FirstName   string         `gorm:"column:first_name;type:text;not null"`
	LastName    string         `gorm:"column:last_name;type:text;not null"`
	Username    string         `gorm:"column:username;type:text;not null;uniqueIndex"`
	Email       string         `gorm:"column:email;type:text;not null;uniqueIndex"`
	Password    string         `gorm:"column:password;type:text;not null"`
	Role        string         `gorm:"column:role;type:text;default:'user'"`
	Active      bool           `gorm:"column:active;default:true"`
	LastLoginAt *time.Time     `gorm:"column:last_login_at"`
	CreatedAt   time.Time      `gorm:"column:created_at"`
	UpdatedAt   time.Time      `gorm:"column:updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"column:deleted_at"`

	Sessions []Session `gorm:"foreignKey:UserID"`
}

func (u *User) BeforeCreate(*gorm.DB) (err error) {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)

	return nil
}

func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}

func (s *Session) IsSessionValid() bool {
	return !s.IsRevoked && time.Now().Before(s.ExpiresAt)
}

func (u *User) GenerateTokens(jwtConf *config.Jwt, userAgent, ipAddress string) (*Session, error) {
	accessTokenDuration := time.Duration(jwtConf.AccessTokenTTL) * time.Second
	refreshTokenDuration := time.Duration(jwtConf.RefreshTokenTTL) * time.Second

	accessToken, accessTokenExpiry, err := jwt.GenerateToken(u.ID.String(), u.Username, u.Email, u.Role, jwtConf, jwt.AccessToken, accessTokenDuration)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	refreshToken, refreshTokenExpiry, err := jwt.GenerateToken(u.ID.String(), u.Username, u.Email, u.Role, jwtConf, jwt.RefreshToken, refreshTokenDuration)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	session := Session{
		UserID:        u.ID,
		AccessToken:   accessToken,
		RefreshToken:  refreshToken,
		ExpiresAt:     accessTokenExpiry,
		RefreshExpiry: refreshTokenExpiry,
		UserAgent:     userAgent,
		IPAddress:     ipAddress,
		IsRevoked:     false,
	}

	return &session, nil
}
