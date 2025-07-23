package service

import (
	"context"
	"errors"
	"public-api/mocks"
	"public-api/model"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestUserService_CreateUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mocks.NewMockUserClient(ctrl)
	svc := NewUserService(mockClient)

	tests := []struct {
		name         string
		inputName    string
		mockBehavior func()
		expectedUser *model.User
		expectError  bool
	}{
		{
			name:      "success",
			inputName: "Rifqi",
			mockBehavior: func() {
				mockClient.EXPECT().
					CreateUser("Rifqi").
					Return(&model.User{ID: 1, Name: "Rifqi"}, nil)
			},
			expectedUser: &model.User{ID: 1, Name: "Rifqi"},
			expectError:  false,
		},
		{
			name:      "empty name",
			inputName: "",
			mockBehavior: func() {
				// No mock expected for empty input
			},
			expectedUser: nil,
			expectError:  true,
		},
		{
			name:      "client error",
			inputName: "ErrorCase",
			mockBehavior: func() {
				mockClient.EXPECT().
					CreateUser("ErrorCase").
					Return(nil, errors.New("failed to create user"))
			},
			expectedUser: nil,
			expectError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockBehavior()

			user, err := svc.CreateUser(context.Background(), tt.inputName)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, user)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedUser, user)
			}
		})
	}
}

func TestUserService_GetUserByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mocks.NewMockUserClient(ctrl)
	svc := NewUserService(mockClient)

	tests := []struct {
		name         string
		inputID      int64
		mockBehavior func()
		expectedUser *model.User
		expectError  bool
	}{
		{
			name:    "success",
			inputID: 123,
			mockBehavior: func() {
				mockClient.EXPECT().
					FetchUserByID(int64(123)).
					Return(&model.User{ID: 123, Name: "TestUser"}, nil)
			},
			expectedUser: &model.User{ID: 123, Name: "TestUser"},
			expectError:  false,
		},
		{
			name:    "user not found",
			inputID: 999,
			mockBehavior: func() {
				mockClient.EXPECT().
					FetchUserByID(int64(999)).
					Return(nil, errors.New("not found"))
			},
			expectedUser: nil,
			expectError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockBehavior()

			user, err := svc.GetUserByID(context.Background(), tt.inputID)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, user)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedUser, user)
			}
		})
	}
}
