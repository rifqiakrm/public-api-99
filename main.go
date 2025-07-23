package main

import (
	"log"
	"public-api/config"

	"public-api/client"
	"public-api/handler"
	"public-api/router"
	"public-api/service"
)

func main() {
	// Load from ENV or use fallback
	cfg := config.Load()

	// Init clients
	listingClient := client.NewListingClient(cfg.ListingServiceURL)
	userClient := client.NewUserClient(cfg.UserServiceURL)

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
