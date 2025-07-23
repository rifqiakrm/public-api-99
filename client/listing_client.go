package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"public-api/model"
	"strconv"
)

//go:generate mockgen -destination=../mocks/mock_listing_client.go -package=mocks public-api/client ListingClient
type ListingClient interface {
	FetchListings(page, size int, userID *int64) ([]model.Listing, error)
	CreateListing(l model.Listing) (*model.Listing, error)
}

// listingClientImpl talks to the Listing Service
type listingClientImpl struct {
	baseURL string
	client  *http.Client
}

// NewListingClient creates a new ListingClient
func NewListingClient(baseURL string) ListingClient {
	return &listingClientImpl{
		baseURL: baseURL,
		client:  &http.Client{},
	}
}

// FetchListings fetches listings, optionally filtered by user_id
func (lc *listingClientImpl) FetchListings(page, size int, userID *int64) ([]model.Listing, error) {
	q := url.Values{}
	q.Set("page_num", strconv.Itoa(page))
	q.Set("page_size", strconv.Itoa(size))
	if userID != nil {
		q.Set("user_id", strconv.FormatInt(*userID, 10))
	}

	url := fmt.Sprintf("%s/listings?%s", lc.baseURL, q.Encode())
	resp, err := lc.client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to call listing service: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("listing service returned status: %d", resp.StatusCode)
	}

	var result struct {
		Result   bool            `json:"result"`
		Listings []model.Listing `json:"listings"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decode listing failed: %w", err)
	}

	return result.Listings, nil
}

// CreateListing creates a listing
func (lc *listingClientImpl) CreateListing(l model.Listing) (*model.Listing, error) {
	form := url.Values{}
	form.Set("user_id", strconv.FormatInt(l.UserID, 10))
	form.Set("listing_type", l.ListingType)
	form.Set("price", strconv.FormatInt(int64(l.Price), 10))

	url := lc.baseURL + "/listings"
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBufferString(form.Encode()))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := lc.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error creating listing: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("listing service returned status: %d", resp.StatusCode)
	}

	var result struct {
		Result  bool           `json:"result"`
		Listing *model.Listing `json:"listing"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode listing creation: %w", err)
	}

	return result.Listing, nil
}
