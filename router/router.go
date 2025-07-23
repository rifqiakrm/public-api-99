package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"public-api/handler"
)

// SetupRouter initializes all routes and handlers
func SetupRouter(
	userHandler *handler.UserHandler,
	listingHandler *handler.ListingHandler,
) *gin.Engine {
	r := gin.Default()

	// Health check
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "pong"})
	})

	api := r.Group("/api/v1")
	{
		// User routes
		api.POST("/users", userHandler.CreateUser)

		// Listing routes
		api.POST("/listings", listingHandler.CreateListing)
		api.GET("/listings", listingHandler.GetListings)
	}

	return r
}
