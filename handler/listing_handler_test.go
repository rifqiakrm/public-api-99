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

func TestListingHandler_CreateListing(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		requestBody    interface{}
		mockService    func(s *mocks.MockListingService)
		expectedCode   int
		expectedResult string
	}{
		{
			name: "success",
			requestBody: model.CreateListingRequest{
				UserID:      1,
				ListingType: "rent",
				Price:       100,
			},
			mockService: func(s *mocks.MockListingService) {
				s.EXPECT().
					CreateListing(gomock.Any(), gomock.Any()).
					Return(&model.Listing{
						ID:          1,
						UserID:      1,
						ListingType: "rent",
						Price:       100,
					}, nil)
			},
			expectedCode:   http.StatusCreated,
			expectedResult: `"listing"`,
		},
		{
			name:           "invalid json",
			requestBody:    `{invalid-json}`,                     // malformed JSON string
			mockService:    func(s *mocks.MockListingService) {}, // no call expected
			expectedCode:   http.StatusBadRequest,
			expectedResult: `"error"`,
		},
		{
			name: "internal error",
			requestBody: model.CreateListingRequest{
				UserID:      1,
				ListingType: "sale",
				Price:       200,
			},
			mockService: func(s *mocks.MockListingService) {
				s.EXPECT().
					CreateListing(gomock.Any(), gomock.Any()).
					Return(&model.Listing{}, errors.New("db error"))
			},
			expectedCode:   http.StatusInternalServerError,
			expectedResult: `"error"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockSvc := mocks.NewMockListingService(ctrl)
			tt.mockService(mockSvc)

			router := gin.Default()
			h := handler.NewListingHandler(mockSvc)
			router.POST("/public-api/listings", h.CreateListing)

			var body []byte
			switch v := tt.requestBody.(type) {
			case string:
				body = []byte(v)
			default:
				body, _ = json.Marshal(v)
			}

			req := httptest.NewRequest(http.MethodPost, "/public-api/listings", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			resp := httptest.NewRecorder()

			router.ServeHTTP(resp, req)

			assert.Equal(t, tt.expectedCode, resp.Code)
			assert.Contains(t, resp.Body.String(), tt.expectedResult)
		})
	}
}

func TestListingHandler_GetListings(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		query          string
		mockService    func(s *mocks.MockListingService)
		expectedCode   int
		expectedResult string
	}{
		{
			name:  "success without userID",
			query: "/public-api/listings?page_num=1&page_size=2",
			mockService: func(s *mocks.MockListingService) {
				s.EXPECT().
					GetListings(gomock.Any(), 1, 2, nil).
					Return([]model.Listing{{ID: 1}, {ID: 2}}, nil)
			},
			expectedCode:   http.StatusOK,
			expectedResult: `"listings"`,
		},
		{
			name:  "success with userID",
			query: "/public-api/listings?page_num=1&page_size=2&user_id=5",
			mockService: func(s *mocks.MockListingService) {
				uid := int64(5)
				s.EXPECT().
					GetListings(gomock.Any(), 1, 2, &uid).
					Return([]model.Listing{{ID: 5}}, nil)
			},
			expectedCode:   http.StatusOK,
			expectedResult: `"listings"`,
		},
		{
			name:  "internal error",
			query: "/public-api/listings?page_num=1&page_size=2",
			mockService: func(s *mocks.MockListingService) {
				s.EXPECT().
					GetListings(gomock.Any(), 1, 2, nil).
					Return(nil, errors.New("internal error"))
			},
			expectedCode:   http.StatusInternalServerError,
			expectedResult: `"error"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockSvc := mocks.NewMockListingService(ctrl)
			tt.mockService(mockSvc)

			router := gin.Default()
			h := handler.NewListingHandler(mockSvc)
			router.GET("/public-api/listings", h.GetListings)

			req := httptest.NewRequest(http.MethodGet, tt.query, nil)
			resp := httptest.NewRecorder()

			router.ServeHTTP(resp, req)

			assert.Equal(t, tt.expectedCode, resp.Code)
			assert.Contains(t, resp.Body.String(), tt.expectedResult)
		})
	}
}
