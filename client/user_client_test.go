package client_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"public-api/client"
	"public-api/model"
	"strings"
	"testing"
)

func TestFetchUserByID(t *testing.T) {
	tests := []struct {
		name       string
		userID     int64
		mockStatus int
		mockBody   string
		wantErr    bool
		wantUser   *model.User
	}{
		{
			name:       "success",
			userID:     1,
			mockStatus: http.StatusOK,
			mockBody:   `{"result": true, "user": {"id": 1, "name": "John"}}`,
			wantUser:   &model.User{ID: 1, Name: "John"},
		},
		{
			name:       "not found",
			userID:     999,
			mockStatus: http.StatusNotFound,
			mockBody:   `not found`,
			wantErr:    true,
		},
		{
			name:       "invalid json",
			userID:     2,
			mockStatus: http.StatusOK,
			mockBody:   `{invalid}`,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.mockStatus)
				fmt.Fprintln(w, tt.mockBody)
			}))
			defer server.Close()

			uc := client.NewUserClient(server.URL)
			user, err := uc.FetchUserByID(tt.userID)

			if (err != nil) != tt.wantErr {
				t.Errorf("FetchUserByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantUser != nil && (user == nil || user.ID != tt.wantUser.ID || user.Name != tt.wantUser.Name) {
				t.Errorf("FetchUserByID() got = %v, want = %v", user, tt.wantUser)
			}
		})
	}
}

func TestCreateUser(t *testing.T) {
	tests := []struct {
		name       string
		inputName  string
		mockStatus int
		mockBody   string
		wantErr    bool
		wantUser   *model.User
	}{
		{
			name:       "success",
			inputName:  "Alice",
			mockStatus: http.StatusCreated,
			mockBody:   `{"result": true, "user": {"id": 2, "name": "Alice"}}`,
			wantUser:   &model.User{ID: 2, Name: "Alice"},
		},
		{
			name:       "bad request",
			inputName:  "Bob",
			mockStatus: http.StatusBadRequest,
			mockBody:   `invalid request`,
			wantErr:    true,
		},
		{
			name:       "malformed json",
			inputName:  "Charlie",
			mockStatus: http.StatusCreated,
			mockBody:   `{invalid}`,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.mockStatus)
				fmt.Fprintln(w, tt.mockBody)
			}))
			defer server.Close()

			uc := client.NewUserClient(server.URL)
			user, err := uc.CreateUser(tt.inputName)

			if (err != nil) != tt.wantErr {
				t.Errorf("CreateUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantUser != nil && (user == nil || user.ID != tt.wantUser.ID || user.Name != tt.wantUser.Name) {
				t.Errorf("CreateUser() got = %v, want = %v", user, tt.wantUser)
			}
		})
	}
}

func TestFetchUsersByIDs(t *testing.T) {
	tests := []struct {
		name       string
		inputIDs   []int64
		mockStatus int
		mockBody   string
		wantErr    bool
		wantCount  int
	}{
		{
			name:       "success",
			inputIDs:   []int64{1, 2},
			mockStatus: http.StatusOK,
			mockBody:   `{"users": [{"id":1,"name":"John"},{"id":2,"name":"Doe"}]}`,
			wantCount:  2,
		},
		{
			name:       "invalid json",
			inputIDs:   []int64{1},
			mockStatus: http.StatusOK,
			mockBody:   `{invalid}`,
			wantErr:    true,
		},
		{
			name:       "error status",
			inputIDs:   []int64{1},
			mockStatus: http.StatusInternalServerError,
			mockBody:   `internal error`,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodPost || !strings.Contains(r.URL.Path, "/users/batch") {
					t.Errorf("unexpected request: %v %v", r.Method, r.URL.Path)
				}
				w.WriteHeader(tt.mockStatus)
				fmt.Fprintln(w, tt.mockBody)
			}))
			defer server.Close()

			uc := client.NewUserClient(server.URL)
			users, err := uc.FetchUsersByIDs(tt.inputIDs)

			if (err != nil) != tt.wantErr {
				t.Errorf("FetchUsersByIDs() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if users != nil && len(users) != tt.wantCount {
				t.Errorf("expected %d users, got %d", tt.wantCount, len(users))
			}
		})
	}
}
