package handler

import (
	"net/http"
	"public-api/model"
	"public-api/service"

	"github.com/gin-gonic/gin"
)

// UserHandler handles HTTP requests related to users
type UserHandler struct {
	service service.UserService
}

// NewUserHandler constructs a new UserHandler
func NewUserHandler(s service.UserService) *UserHandler {
	return &UserHandler{service: s}
}

// CreateUser handles POST /public-api/users
func (h *UserHandler) CreateUser(c *gin.Context) {
	var req model.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	user, err := h.service.CreateUser(c.Request.Context(), req.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"user": user})
}
