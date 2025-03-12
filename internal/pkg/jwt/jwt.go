package jwt

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"itv-movie/internal/config"
	"time"
)

var (
	ErrInvalidToken = errors.New("invalid token")
	ErrExpiredToken = errors.New("token has expired")
)

type TokenType string

const (
	AccessToken  TokenType = "access"
	RefreshToken TokenType = "refresh"
)

// CustomClaims contains the claims we want in our JWT tokens
type CustomClaims struct {
	UserID    string    `json:"user_id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	TokenType TokenType `json:"token_type"`
	jwt.RegisteredClaims
}

func GenerateToken(userId, username, email, role string, jwtConfig *config.Jwt, tokenType TokenType, duration time.Duration) (string, time.Time, error) {
	now := time.Now()

	if duration <= 0 {
		if tokenType == AccessToken {
			duration = 30 * time.Minute
		} else {
			duration = 7 * 24 * time.Hour
		}
		fmt.Printf("Warning: Invalid token duration provided. Using default: %v\n", duration)
	}

	expirationTime := now.Add(duration)

	if expirationTime.Equal(now) {
		expirationTime = now.Add(30 * time.Minute)
		fmt.Println("Warning: Expiration time equals issue time. Setting to now+30min.")
	}

	claims := CustomClaims{
		UserID:    userId,
		Username:  username,
		Email:     email,
		Role:      role,
		TokenType: tokenType,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Issuer:    jwtConfig.Domain,
			Subject:   jwtConfig.Realm,
			ID:        uuid.New().String(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(jwtConfig.Secret))
	if err != nil {
		return "", time.Time{}, err
	}

	return tokenString, expirationTime, nil
}

func ValidateToken(tokenString, secretKey string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(
		tokenString,
		&CustomClaims{},
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(secretKey), nil
		},
	)

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) || errors.Is(err, jwt.ErrTokenNotValidYet) {
			return nil, ErrExpiredToken
		}
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	claims, ok := token.Claims.(*CustomClaims)
	if !ok || !token.Valid {
		return nil, ErrInvalidToken
	}

	return claims, nil
}

// ExtractBearerToken extracts the token from the "Bearer <token>" format
func ExtractBearerToken(authHeader string) (string, error) {
	if len(authHeader) < 7 || authHeader[:7] != "Bearer " {
		return "", errors.New("bearer token not found in authorization header")
	}
	return authHeader[7:], nil
}
