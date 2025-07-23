package main

import (
	"log"
	"os"

	"public-api/client"
	"public-api/handler"
	"public-api/router"
	"public-api/service"
)

func main() {
	// Load from ENV or use fallback
	listingURL := os.Getenv("LISTING_SERVICE_URL")
	if listingURL == "" {
		listingURL = "http://localhost:6000"
	}

	userURL := os.Getenv("USER_SERVICE_URL")
	if userURL == "" {
		userURL = "http://localhost:6001"
	}

	// Init clients
	listingClient := client.NewListingClient(listingURL)
	userClient := client.NewUserClient(userURL)

	// Init services
	listingService := service.NewListingService(listingClient, userClient)
	userService := service.NewUserService(userClient)

	// Init handlers
	listingHandler := handler.NewListingHandler(listingService)
	userHandler := handler.NewUserHandler(userService)

	// Setup and run router
	r := router.SetupRouter(userHandler, listingHandler)

	log.Println("🚀 Public API is running at :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
