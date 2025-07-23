package handler_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"public-api/handler"
	"public-api/mocks"
	"public-api/model"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestUserHandler_CreateUser(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockUserService(ctrl)
	handler := handler.NewUserHandler(mockService)

	tests := []struct {
		name           string
		requestBody    interface{}
		mockSetup      func()
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "success",
			requestBody: model.CreateUserRequest{
				Name: "John Doe",
			},
			mockSetup: func() {
				mockService.EXPECT().
					CreateUser(gomock.Any(), "John Doe").
					Return(&model.User{
						ID:   1,
						Name: "John Doe",
					}, nil)
			},
			expectedStatus: http.StatusCreated,
			expectedBody:   `"name":"John Doe"`,
		},
		{
			name:           "invalid JSON",
			requestBody:    `{invalid}`,
			mockSetup:      func() {}, // No call expected
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `Invalid request`,
		},
		{
			name: "internal error",
			requestBody: model.CreateUserRequest{
				Name: "Failing User",
			},
			mockSetup: func() {
				mockService.EXPECT().
					CreateUser(gomock.Any(), "Failing User").
					Return(nil, errors.New("failed to create user"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `failed to create user`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mocks
			tt.mockSetup()

			// Prepare body
			var reqBodyBytes []byte
			var req *http.Request
			if str, ok := tt.requestBody.(string); ok {
				// malformed case
				req = httptest.NewRequest(http.MethodPost, "/public-api/users", bytes.NewBufferString(str))
			} else {
				reqBodyBytes, _ = json.Marshal(tt.requestBody)
				req = httptest.NewRequest(http.MethodPost, "/public-api/users", bytes.NewBuffer(reqBodyBytes))
			}

			// Setup context
			rec := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(rec)
			ctx.Request = req

			// Invoke handler
			handler.CreateUser(ctx)

			// Assert
			assert.Equal(t, tt.expectedStatus, rec.Code)
			assert.Contains(t, rec.Body.String(), tt.expectedBody)
		})
	}
}
