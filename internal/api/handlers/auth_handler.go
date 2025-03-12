package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"itv-movie/internal/api/services"
	"itv-movie/internal/models"
	"net/http"
	"time"
)

type AuthHandler struct {
	authService *services.AuthService
}

func NewAuthHandler(authService *services.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// Register handles user registration
func (h *AuthHandler) Register(c *gin.Context) {
	var registerRequest struct {
		FirstName string `json:"firstName" binding:"required"`
		LastName  string `json:"lastName" binding:"required"`
		Username  string `json:"username" binding:"required"`
		Email     string `json:"email" binding:"required,email"`
		Password  string `json:"password" binding:"required,min=8"`
	}

	if err := c.BindJSON(&registerRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format: " + err.Error()})
		return
	}

	newUser := &models.User{
		FirstName: registerRequest.FirstName,
		LastName:  registerRequest.LastName,
		Username:  registerRequest.Username,
		Email:     registerRequest.Email,
		Password:  registerRequest.Password, // Will be hashed in BeforeCreate hook
		Role:      "user",                   // Default role
		Active:    true,
	}

	createdUser, err := h.authService.RegisterUser(c, newUser)
	if err != nil {
		if err == services.ErrUserExists {
			c.JSON(http.StatusConflict, gin.H{"error": "User with this email or username already exists"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register user: " + err.Error()})
		}
		return
	}

	// Don't return the password hash
	createdUser.Password = ""

	c.JSON(http.StatusCreated, createdUser)
}

// Login handles user authentication
func (h *AuthHandler) Login(c *gin.Context) {
	var loginRequest struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.BindJSON(&loginRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format: " + err.Error()})
		return
	}

	// Get IP and user agent for the session
	userAgent := c.GetHeader("User-Agent")
	ipAddress := c.ClientIP()

	user, session, err := h.authService.Login(c, loginRequest.Username, loginRequest.Password, userAgent, ipAddress)
	if err != nil {
		switch err {
		case services.ErrInvalidCredentials:
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		case services.ErrUserInactive:
			c.JSON(http.StatusForbidden, gin.H{"error": "User account is inactive"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Login failed: " + err.Error()})
		}
		return
	}

	// Don't return the password hash
	user.Password = ""

	// Set session token in cookie
	c.SetCookie(
		"session_token",
		session.Token,
		int(time.Until(session.ExpiresAt).Seconds()),
		"/",
		"",    // Domain
		false, // Not secure for local dev (change to true in production)
		true,  // HTTP only
	)

	c.JSON(http.StatusOK, gin.H{
		"user":      user,
		"token":     session.Token,
		"expiresAt": session.ExpiresAt,
	})
}

// Logout handles user logout
func (h *AuthHandler) Logout(c *gin.Context) {
	// Get token from cookie or header
	token, err := c.Cookie("session_token")
	if err != nil {
		// If not in cookie, try from Authorization header
		token = c.GetHeader("Authorization")
		if token == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "No session token provided"})
			return
		}
	}

	if err := h.authService.Logout(c, token); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to logout: " + err.Error()})
		return
	}

	// Clear the cookie
	c.SetCookie("session_token", "", -1, "/", "", false, true)

	c.JSON(http.StatusOK, gin.H{"message": "Successfully logged out"})
}

func (h *AuthHandler) UpdateStatus(c *gin.Context) {
	var updateRequest struct {
		UserID string `json:"userId" binding:"required"`
		Active bool   `json:"active" binding:"required"`
	}

	if err := c.BindJSON(&updateRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format: " + err.Error()})
		return
	}

	userID, err := uuid.Parse(updateRequest.UserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID format"})
		return
	}

	// TODO: Verify that the requester is an admin
	// This would normally be done by middleware or an auth check here

	if err := h.authService.UpdateUserStatus(c, userID, updateRequest.Active); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user status: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User status updated successfully"})
}

func (h *AuthHandler) DeleteUser(c *gin.Context) {
	idStr := c.Param("id")
	userID, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID format"})
		return
	}

	// TODO: Verify that the requester is an admin or the user themself
	// This would normally be done by middleware or an auth check here

	if err := h.authService.DeleteUser(c, userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}
