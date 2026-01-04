package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gal1996/vibe_coding_with_architecture/usecase/interactor"
)

// UserHandler handles HTTP requests for users
type UserHandler struct {
	userUseCase *interactor.UserUseCase
}

// NewUserHandler creates a new user handler
func NewUserHandler(userUseCase *interactor.UserUseCase) *UserHandler {
	return &UserHandler{
		userUseCase: userUseCase,
	}
}

// RegisterRequest represents the request body for user registration
type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Password string `json:"password" binding:"required,min=6"`
	IsAdmin  bool   `json:"is_admin"`
}

// Register handles POST /register
func (h *UserHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	input := interactor.RegisterInput{
		Username: req.Username,
		Password: req.Password,
		IsAdmin:  req.IsAdmin,
	}

	user, err := h.userUseCase.Register(c.Request.Context(), input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Don't return the password hash
	c.JSON(http.StatusCreated, gin.H{
		"id":       user.ID,
		"username": user.Username,
		"is_admin": user.IsAdmin,
	})
}

// LoginRequest represents the request body for user login
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// Login handles POST /login
func (h *UserHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	input := interactor.LoginInput{
		Username: req.Username,
		Password: req.Password,
	}

	output, err := h.userUseCase.Login(c.Request.Context(), input)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": output.Token,
		"user": gin.H{
			"id":       output.User.ID,
			"username": output.User.Username,
			"is_admin": output.User.IsAdmin,
		},
	})
}

// GetProfile handles GET /profile
func (h *UserHandler) GetProfile(c *gin.Context) {
	userID := c.Param("id")

	user, err := h.userUseCase.GetUser(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":       user.ID,
		"username": user.Username,
		"is_admin": user.IsAdmin,
	})
}