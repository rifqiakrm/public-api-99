package model

// User represents a user in the system
type User struct {
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	CreatedAt int64  `json:"created_at"`
	UpdatedAt int64  `json:"updated_at"`
}

// Listing represents a property listing
type Listing struct {
	ID          int64   `json:"id"`
	UserID      int64   `json:"user_id"`
	ListingType string  `json:"listing_type"`
	Price       float64 `json:"price"`
	CreatedAt   int64   `json:"created_at"`
	UpdatedAt   int64   `json:"updated_at"`
	User        *User   `json:"user,omitempty"` // Only populated in /public-api
}
