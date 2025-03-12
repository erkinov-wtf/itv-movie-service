package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"itv-movie/internal/api/services"
	"itv-movie/internal/pkg/jwt"
	"itv-movie/internal/pkg/utils/constants"
	"net/http"
)

func AuthMiddleware(authService *services.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get token from Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			c.Abort()
			return
		}

		// Extract token from header
		tokenString, err := jwt.ExtractBearerToken(authHeader)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			c.Abort()
			return
		}

		// Validate the token
		claims, err := authService.ValidateAccessToken(c, tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		// Set user info in context for downstream handlers
		userID, _ := uuid.Parse(claims.UserID)
		c.Set("userID", userID)
		c.Set("username", claims.Username)
		c.Set("email", claims.Email)
		c.Set("role", claims.Role)
		c.Set("accessToken", tokenString)

		c.Next()
	}
}

// RoleMiddleware checks if the user has a specific role
func RoleMiddleware(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, exists := c.Get("role")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User role not found in context"})
			c.Abort()
			return
		}

		// Check if user's role is in the allowed roles
		roleMatched := false
		for _, role := range roles {
			if role == userRole.(string) {
				roleMatched = true
				break
			}
		}

		if !roleMatched {
			c.JSON(http.StatusForbidden, gin.H{"error": "You don't have permission to access this resource"})
			c.Abort()
			return
		}

		c.Next()
	}
}

// AdminOnly middleware ensures the user is an admin
func AdminOnly() gin.HandlerFunc {
	return RoleMiddleware(constants.AdminRole)
}
