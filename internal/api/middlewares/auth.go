package middlewares

import (
	"github.com/gin-gonic/gin"
	"itv-movie/internal/api/services"
	"net/http"
)

func AuthMiddleware(authService *services.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := c.Cookie("session_token")
		if err != nil {
			// If not in cookie, try from Authorization header
			token = c.GetHeader("Authorization")
			if token == "" {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
				c.Abort()
				return
			}
		}

		// TODO: Implement actual token validation with JWT
		// For now, we just pass through as a placeholder

		session, err := authService.ValidateSession(c, token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired session"})
			c.Abort()
			return
		}

		c.Set("userID", session.UserID)
		c.Set("sessionToken", token)

		c.Next()
	}
}

// AdminOnly middleware ensures the user is an admin
func AdminOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: Implement actual admin check
		// This is just a placeholder

		// For now, we'll just pass through
		c.Next()
	}
}
