package model

// CreateUserRequest represents the payload to create a user
type CreateUserRequest struct {
	Name string `json:"name" binding:"required"`
}

// CreateListingRequest represents the payload to create a listing
type CreateListingRequest struct {
	UserID      int64   `json:"user_id" binding:"required"`
	ListingType string  `json:"listing_type" binding:"required"`
	Price       float64 `json:"price" binding:"required"`
}
