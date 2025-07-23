package handler

import (
	"net/http"
	"public-api/model"
	"public-api/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

// ListingHandler handles HTTP requests related to listings
type ListingHandler struct {
	service service.ListingService
}

// NewListingHandler constructs a new ListingHandler
func NewListingHandler(s service.ListingService) *ListingHandler {
	return &ListingHandler{service: s}
}

// CreateListing handles POST /public-api/listings
func (h *ListingHandler) CreateListing(c *gin.Context) {
	var req model.CreateListingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	l := model.Listing{
		UserID:      req.UserID,
		ListingType: req.ListingType,
		Price:       req.Price,
	}

	created, err := h.service.CreateListing(c.Request.Context(), l)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"listing": created})
}

// GetListings handles GET /public-api/listings
func (h *ListingHandler) GetListings(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page_num", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	var userIDPtr *int64

	if userIDStr := c.Query("user_id"); userIDStr != "" {
		if id, err := strconv.ParseInt(userIDStr, 10, 64); err == nil {
			userIDPtr = &id
		}
	}

	listings, err := h.service.GetListings(c.Request.Context(), page, size, userIDPtr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"listings": listings})
}
