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
		auth.POST("/register-admin", authHandler.RegisterAdmin)
		auth.POST("/login", authHandler.Login)
		auth.POST("/refresh", authHandler.RefreshToken)

		// Protected routes
		auth.Use(middlewares.AuthMiddleware(authService))
		auth.POST("/logout", authHandler.Logout)

		// Admin routes
		admin := auth.Group("/admin")
		admin.Use(middlewares.AdminOnly())
		admin.POST("/register-director", authHandler.RegisterDirector)
		admin.GET("/users", authHandler.GetAllUsers)
		admin.PUT("/status", authHandler.UpdateStatus)
		admin.DELETE("/users/:id", authHandler.DeleteUser)
	}
}
