package middlewares

import (
	"github.com/gin-gonic/gin"
	"itv-movie/internal/pkg/utils/constants"
)

// AdminOnly middleware ensures the user is an admin
func AdminOnly() gin.HandlerFunc {
	return RoleMiddleware(constants.AdminRole)
}

// AdminOrDirectorOnly checks if the user is either an admin or a director
func AdminOrDirectorOnly() gin.HandlerFunc {
	return RoleMiddleware(constants.AdminRole, constants.DirectorRole)
}

// DirectorOnly ensures the user is a director
func DirectorOnly() gin.HandlerFunc {
	return RoleMiddleware(constants.DirectorRole)
}
