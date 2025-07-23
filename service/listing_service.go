package service

import (
	"context"
	"fmt"
	"log"
	"public-api/client"
	"public-api/model"
)

//go:generate mockgen -destination=../mocks/mock_listing_service.go -package=mocks public-api/service ListingService
type ListingService interface {
	CreateListing(ctx context.Context, l model.Listing) (*model.Listing, error)
	GetListings(ctx context.Context, page, size int, userID *int64) ([]model.Listing, error)
}

// listingServiceImpl handles listing-related logic for public API
type listingServiceImpl struct {
	listingClient client.ListingClient
	userClient    client.UserClient
}

// NewListingService constructs a new ListingService
func NewListingService(lc client.ListingClient, uc client.UserClient) ListingService {
	return &listingServiceImpl{
		listingClient: lc,
		userClient:    uc,
	}
}

// CreateListing creates a new listing via listing-service
func (ls *listingServiceImpl) CreateListing(ctx context.Context, l model.Listing) (*model.Listing, error) {
	if l.UserID == 0 || l.Price <= 0 || l.ListingType == "" {
		return nil, fmt.Errorf("user_id, price, and listing_type are required")
	}
	return ls.listingClient.CreateListing(l)
}

// GetListings fetches listings and attaches user info to each one
func (ls *listingServiceImpl) GetListings(ctx context.Context, page, size int, userID *int64) ([]model.Listing, error) {
	listings, err := ls.listingClient.FetchListings(page, size, userID)
	if err != nil {
		return nil, err
	}

	if len(listings) == 0 {
		return listings, nil
	}

	// 1. Collect unique user IDs
	userIDSet := make(map[int64]struct{})
	for _, l := range listings {
		userIDSet[l.UserID] = struct{}{}
	}

	var userIDs []int64
	for id := range userIDSet {
		userIDs = append(userIDs, id)
	}

	// 2. Fetch all users in batch
	usersMap, err := ls.userClient.FetchUsersByIDs(userIDs)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch users: %w", err)
	}

	// 3. Attach user info to listings
	for i, l := range listings {
		if user, ok := usersMap[l.UserID]; ok {
			listings[i].User = user
		} else {
			log.Println("user not found", l.UserID)
		}
	}

	return listings, nil
}
