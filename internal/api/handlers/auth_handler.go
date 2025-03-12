package handlers

import (
	"github.com/gin-gonic/gin"
	"itv-movie/internal/api/services"
)

// AuthHandler handles HTTP requests for Auths
type AuthHandler struct {
	AuthService *services.AuthService
}

// NewAuthHandler creates a new Auth handler
func NewAuthHandler(AuthService *services.AuthService) *AuthHandler {
	return &AuthHandler{
		AuthService: AuthService,
	}
}

func (h *AuthHandler) Login(c *gin.Context) {
	//TODO implement real handler logic
}

func (h *AuthHandler) Register(c *gin.Context) {
	//TODO implement real handler logic
}

func (h *AuthHandler) Logout(c *gin.Context) {
	//TODO implement real handler logic
}

func (h *AuthHandler) UpdateStatus(c *gin.Context) {
	//TODO implement real handler logic
}

func (h *AuthHandler) DeleteUser(c *gin.Context) {
	//TODO implement real handler logic
}
