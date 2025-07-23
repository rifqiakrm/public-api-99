package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"public-api/model"
)

//go:generate mockgen -destination=../mocks/mock_user_client.go -package=mocks public-api/client UserClient
type UserClient interface {
	FetchUserByID(id int64) (*model.User, error)
	CreateUser(name string) (*model.User, error)
	FetchUsersByIDs(ids []int64) (map[int64]*model.User, error)
}

// userClientImpl handles HTTP calls to the User Service
type userClientImpl struct {
	baseURL string
	client  *http.Client
}

// NewUserClient creates a new UserClient
func NewUserClient(baseURL string) UserClient {
	return &userClientImpl{
		baseURL: baseURL,
		client:  &http.Client{},
	}
}

// FetchUserByID gets a user by ID from the User Service
func (uc *userClientImpl) FetchUserByID(id int64) (*model.User, error) {
	url := fmt.Sprintf("%s/users/%d", uc.baseURL, id)
	resp, err := uc.client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to call user service: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("user service returned status: %d", resp.StatusCode)
	}

	var result struct {
		Result bool        `json:"result"`
		User   *model.User `json:"user"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode user response: %w", err)
	}

	return result.User, nil
}

// CreateUser creates a user via POST using application/json
func (uc *userClientImpl) CreateUser(name string) (*model.User, error) {
	payload := map[string]string{
		"name": name,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal user payload: %w", err)
	}

	url := uc.baseURL + "/users"
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	resp, err := uc.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error creating user: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("user service returned status: %d", resp.StatusCode)
	}

	var result struct {
		Result bool        `json:"result"`
		User   *model.User `json:"user"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode user creation: %w", err)
	}

	return result.User, nil
}

// FetchUsersByIDs fetch users via POST by IDs using application/json
func (c *userClientImpl) FetchUsersByIDs(ids []int64) (map[int64]*model.User, error) {
	payload := map[string]interface{}{"user_ids": ids}

	body, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", c.baseURL+"/users/batch", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch users, status: %d", resp.StatusCode)
	}

	var res struct {
		Users []model.User `json:"users"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, err
	}

	userMap := make(map[int64]*model.User)
	for _, u := range res.Users {
		u := u
		userMap[u.ID] = &u
	}

	return userMap, nil
}
