package service

import (
	"context"
	"fmt"
	"public-api/client"
	"public-api/model"
)

//go:generate mockgen -destination=../mocks/mock_user_service.go -package=mocks public-api/service UserService
type UserService interface {
	CreateUser(ctx context.Context, name string) (*model.User, error)
	GetUserByID(ctx context.Context, id int64) (*model.User, error)
}

// UserService handles user-related operations for the public API
type userServiceImpl struct {
	client client.UserClient
}

// NewUserService constructs a new UserService
func NewUserService(client client.UserClient) UserService {
	return &userServiceImpl{client: client}
}

// CreateUser creates a user by delegating to the user-service
func (us *userServiceImpl) CreateUser(ctx context.Context, name string) (*model.User, error) {
	if name == "" {
		return nil, fmt.Errorf("name is required")
	}
	return us.client.CreateUser(name)
}

// GetUserByID fetches a user by ID
func (us *userServiceImpl) GetUserByID(ctx context.Context, id int64) (*model.User, error) {
	return us.client.FetchUserByID(id)
}
