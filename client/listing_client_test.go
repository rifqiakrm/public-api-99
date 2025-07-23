package client_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"public-api/client"
	"public-api/model"
	"strings"
	"testing"
)

func TestFetchListings(t *testing.T) {
	tests := []struct {
		name         string
		userID       *int64
		responseCode int
		responseBody string
		expectErr    bool
	}{
		{
			name:         "success without user_id",
			userID:       nil,
			responseCode: http.StatusOK,
			responseBody: `{"result": true, "listings": [{"id":1,"user_id":2,"listing_type":"sell","price":1000}]}`,
			expectErr:    false,
		},
		{
			name:         "success with user_id",
			userID:       ptrInt64(5),
			responseCode: http.StatusOK,
			responseBody: `{"result": true, "listings": [{"id":2,"user_id":5,"listing_type":"buy","price":1500}]}`,
			expectErr:    false,
		},
		{
			name:         "non-200 response",
			responseCode: http.StatusBadRequest,
			responseBody: `{"result": false}`,
			expectErr:    true,
		},
		{
			name:         "invalid JSON",
			responseCode: http.StatusOK,
			responseBody: `{invalid json}`,
			expectErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.responseCode)
				io.WriteString(w, tt.responseBody)
			}))
			defer srv.Close()

			c := client.NewListingClient(srv.URL)
			_, err := c.FetchListings(1, 10, tt.userID)

			if tt.expectErr && err == nil {
				t.Errorf("expected error, got nil")
			}
			if !tt.expectErr && err != nil {
				t.Errorf("did not expect error, got %v", err)
			}
		})
	}
}

func TestCreateListing(t *testing.T) {
	tests := []struct {
		name         string
		responseCode int
		responseBody string
		expectErr    bool
	}{
		{
			name:         "success create listing",
			responseCode: http.StatusOK,
			responseBody: `{"result": true, "listing": {"id": 1, "user_id": 10, "listing_type": "sell", "price": 1000}}`,
			expectErr:    false,
		},
		{
			name:         "non-200 response",
			responseCode: http.StatusInternalServerError,
			responseBody: `{"result": false}`,
			expectErr:    true,
		},
		{
			name:         "invalid JSON response",
			responseCode: http.StatusOK,
			responseBody: `{invalid json}`,
			expectErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodPost {
					t.Errorf("Expected POST method, got %s", r.Method)
				}
				bodyBytes, _ := io.ReadAll(r.Body)
				bodyStr := string(bodyBytes)
				if !strings.Contains(bodyStr, "user_id") {
					t.Errorf("Expected user_id in form, got %s", bodyStr)
				}
				w.WriteHeader(tt.responseCode)
				io.WriteString(w, tt.responseBody)
			}))
			defer srv.Close()

			c := client.NewListingClient(srv.URL)
			listing := model.Listing{
				UserID:      10,
				ListingType: "sell",
				Price:       1000,
			}

			_, err := c.CreateListing(listing)
			if tt.expectErr && err == nil {
				t.Errorf("expected error, got nil")
			}
			if !tt.expectErr && err != nil {
				t.Errorf("did not expect error, got %v", err)
			}
		})
	}
}

func ptrInt64(v int64) *int64 {
	return &v
}
