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

func (u *User) GenerateSession(expiresIn time.Duration, userAgent, ipAddress string) (*Session, error) {
	session := Session{
		UserID:    u.ID,
		Token:     uuid.New().String(), // TODO replace with actual jwt
		ExpiresAt: time.Now().Add(expiresIn),
		UserAgent: userAgent,
		IPAddress: ipAddress,
		IsRevoked: false,
	}

	return &session, nil
}

func (s *Session) IsSessionValid() bool {
	return !s.IsRevoked && time.Now().Before(s.ExpiresAt)
}

func (u *User) GenerateTokens(jwtConf *config.Jwt, userAgent, ipAddress string) (*Session, error) {
	// Generate an access token
	accessToken, accessTokenExpiry, err := jwt.GenerateToken(u, jwtConf, jwt.AccessToken, time.Duration(jwtConf.AccessTokenTTL))
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	// Generate a refresh token
	refreshToken, refreshTokenExpiry, err := jwt.GenerateToken(u, jwtConf, jwt.RefreshToken, time.Duration(jwtConf.RefreshTokenTTL))
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	// Create and return a new session
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
