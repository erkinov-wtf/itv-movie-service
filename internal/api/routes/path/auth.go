package path

import (
	"github.com/gin-gonic/gin"
	"itv-movie/internal/api/handlers"
	"itv-movie/internal/api/middlewares"
	"itv-movie/internal/api/services"
)

func RegisterAuthRoutes(router *gin.RouterGroup, authHandler *handlers.AuthHandler, authService *services.AuthService) {
	auth := router.Group("/auth")
	{
		// Public routes
		auth.POST("/register", authHandler.Register)
		auth.POST("/login", authHandler.Login)

		// Protected routes
		auth.Use(middlewares.AuthMiddleware(authService))
		auth.POST("/logout", authHandler.Logout)

		// Admin routes
		admin := auth.Group("/admin")
		admin.Use(middlewares.AdminOnly())
		admin.PUT("/status", authHandler.UpdateStatus)
		admin.DELETE("/users/:id", authHandler.DeleteUser)
	}
}
