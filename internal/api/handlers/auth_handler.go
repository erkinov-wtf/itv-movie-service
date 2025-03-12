package handlers

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/google/uuid"
	"itv-movie/internal/api/services"
	"itv-movie/internal/models"
	"itv-movie/internal/pkg/jwt"
	"itv-movie/internal/pkg/utils/constants"
	"math"
	"net/http"
	"strconv"
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
		Password  string `json:"password" binding:"required,min=4"`
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
		Password:  registerRequest.Password,
		Role:      constants.UserRole,
		Active:    true,
	}

	createdUser, err := h.authService.RegisterUser(c, newUser)
	if err != nil {
		if errors.Is(err, services.ErrUserExists) {
			c.JSON(http.StatusConflict, gin.H{"error": "User with this email or username already exists"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register user: " + err.Error()})
		}
		return
	}

	createdUser.Password = ""

	c.JSON(http.StatusCreated, createdUser)
}

func (h *AuthHandler) Login(c *gin.Context) {
	var loginRequest struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.BindJSON(&loginRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format: " + err.Error()})
		return
	}

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

	user.Password = ""

	c.JSON(http.StatusOK, gin.H{
		"user":          user,
		"access_token":  session.AccessToken,
		"refresh_token": session.RefreshToken,
		"expires_at":    session.ExpiresAt,
	})
}

func (h *AuthHandler) Logout(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Authorization header is required"})
		return
	}

	accessToken, err := jwt.ExtractBearerToken(authHeader)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err = h.authService.Logout(c, accessToken); err != nil {
		if errors.Is(err, services.ErrSessionInvalid) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid session"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to logout: " + err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Successfully logged out"})
}

func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var refreshRequest struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}

	if err := c.BindJSON(&refreshRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format: " + err.Error()})
		return
	}

	userAgent := c.GetHeader("User-Agent")
	ipAddress := c.ClientIP()

	newSession, err := h.authService.RefreshTokens(c, refreshRequest.RefreshToken, userAgent, ipAddress)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrSessionInvalid):
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid session"})
		case errors.Is(err, services.ErrRefreshTokenExpired):
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Refresh token has expired, please login again"})
		case errors.Is(err, services.ErrInvalidToken):
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid token format"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to refresh token: " + err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"access_token":  newSession.AccessToken,
		"refresh_token": newSession.RefreshToken,
		"expires_at":    newSession.ExpiresAt,
	})
}

func (h *AuthHandler) UpdateStatus(c *gin.Context) {
	var updateRequest struct {
		UserID uuid.UUID `json:"userId"`
		Active bool      `json:"active"`
	}

	if err := c.ShouldBindWith(&updateRequest, binding.JSON); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format: " + err.Error()})
		return
	}

	// Validate manually
	if updateRequest.UserID == uuid.Nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User ID is required"})
		return
	}

	if err := h.authService.UpdateUserStatus(c, updateRequest.UserID, updateRequest.Active); err != nil {
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

	if err = h.authService.DeleteUser(c, userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}

func (h *AuthHandler) RegisterAdmin(c *gin.Context) {
	exists, err := h.authService.AdminExists(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve data: " + err.Error()})
		return
	}

	if exists {
		c.JSON(http.StatusForbidden, gin.H{"error": "Cant create new Admin, already exists"})
		return
	}

	var registerRequest struct {
		FirstName string `json:"firstName" binding:"required"`
		LastName  string `json:"lastName" binding:"required"`
		Username  string `json:"username" binding:"required"`
		Email     string `json:"email" binding:"required,email"`
		Password  string `json:"password" binding:"required,min=4"`
	}

	if err = c.BindJSON(&registerRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format: " + err.Error()})
		return
	}

	newUser := &models.User{
		FirstName: registerRequest.FirstName,
		LastName:  registerRequest.LastName,
		Username:  registerRequest.Username,
		Email:     registerRequest.Email,
		Password:  registerRequest.Password,
		Role:      constants.AdminRole,
		Active:    true,
	}

	createdUser, err := h.authService.RegisterUser(c, newUser)
	if err != nil {
		if errors.Is(err, services.ErrUserExists) {
			c.JSON(http.StatusConflict, gin.H{"error": "User with this email or username already exists"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register user: " + err.Error()})
		}
		return
	}

	createdUser.Password = ""

	c.JSON(http.StatusCreated, createdUser)
}

func (h *AuthHandler) RegisterDirector(c *gin.Context) {
	var registerRequest struct {
		FirstName string `json:"firstName" binding:"required"`
		LastName  string `json:"lastName" binding:"required"`
		Username  string `json:"username" binding:"required"`
		Email     string `json:"email" binding:"required,email"`
		Password  string `json:"password" binding:"required,min=4"`
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
		Password:  registerRequest.Password,
		Role:      constants.DirectorRole,
		Active:    true,
	}

	createdUser, err := h.authService.RegisterUser(c, newUser)
	if err != nil {
		if errors.Is(err, services.ErrUserExists) {
			c.JSON(http.StatusConflict, gin.H{"error": "User with this email or username already exists"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register user: " + err.Error()})
		}
		return
	}

	createdUser.Password = ""

	c.JSON(http.StatusCreated, createdUser)
}

func (h *AuthHandler) GetAllUsers(c *gin.Context) {
	pageStr := c.Query("page")
	limitStr := c.Query("limit")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		limit = 10
	}

	users, total, err := h.authService.GetAllUsers(c, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch users: " + err.Error()})
		return
	}

	for _, user := range users {
		user.Password = "" // sanitization
	}

	lastPage := int(math.Ceil(float64(total) / float64(limit)))
	hasNextPage := page < lastPage
	hasPrevPage := page > 1

	c.JSON(http.StatusOK, gin.H{
		"data": users,
		"meta": gin.H{
			"total":        total,
			"per_page":     limit,
			"current_page": page,
			"last_page":    lastPage,
			"has_next":     hasNextPage,
			"has_prev":     hasPrevPage,
		},
	})
}
