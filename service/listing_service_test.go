// service/listing_service_test.go
package service_test

import (
	"context"
	"errors"
	"public-api/mocks"
	"public-api/model"
	"public-api/service"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestCreateListing(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	listingClient := mocks.NewMockListingClient(ctrl)
	userClient := mocks.NewMockUserClient(ctrl)
	svc := service.NewListingService(listingClient, userClient)

	tests := []struct {
		name       string
		input      model.Listing
		mock       func()
		wantErr    bool
		assertFunc func(t *testing.T, res *model.Listing, err error)
	}{
		{
			name: "missing fields",
			input: model.Listing{
				UserID:      0,
				Price:       0,
				ListingType: "",
			},
			mock:    func() {},
			wantErr: true,
			assertFunc: func(t *testing.T, res *model.Listing, err error) {
				assert.Nil(t, res)
				assert.Error(t, err)
			},
		},
		{
			name: "success create listing",
			input: model.Listing{
				UserID:      1,
				Price:       1000,
				ListingType: "house",
			},
			mock: func() {
				listingClient.EXPECT().
					CreateListing(gomock.Any()).
					Return(&model.Listing{ID: 1, UserID: 1, Price: 1000, ListingType: "house"}, nil)
			},
			wantErr: false,
			assertFunc: func(t *testing.T, res *model.Listing, err error) {
				assert.NoError(t, err)
				assert.Equal(t, int64(1), res.ID)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			res, err := svc.CreateListing(context.Background(), tt.input)
			tt.assertFunc(t, res, err)
		})
	}
}

func TestGetListings(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	listingClient := mocks.NewMockListingClient(ctrl)
	userClient := mocks.NewMockUserClient(ctrl)
	svc := service.NewListingService(listingClient, userClient)

	tests := []struct {
		name       string
		page, size int
		userID     *int64
		mock       func()
		wantErr    bool
		assertFunc func(t *testing.T, res []model.Listing, err error)
	}{
		{
			name:   "listing client error",
			page:   1,
			size:   10,
			userID: nil,
			mock: func() {
				listingClient.EXPECT().
					FetchListings(1, 10, nil).
					Return(nil, errors.New("listing error"))
			},
			wantErr: true,
			assertFunc: func(t *testing.T, res []model.Listing, err error) {
				assert.Nil(t, res)
				assert.Error(t, err)
			},
		},
		{
			name:   "empty listings",
			page:   1,
			size:   10,
			userID: nil,
			mock: func() {
				listingClient.EXPECT().
					FetchListings(1, 10, nil).
					Return([]model.Listing{}, nil)
			},
			wantErr: false,
			assertFunc: func(t *testing.T, res []model.Listing, err error) {
				assert.Empty(t, res)
				assert.NoError(t, err)
			},
		},
		{
			name:   "user client error",
			page:   1,
			size:   10,
			userID: nil,
			mock: func() {
				listingClient.EXPECT().
					FetchListings(1, 10, nil).
					Return([]model.Listing{
						{ID: 1, UserID: 123, Price: 999},
					}, nil)

				userClient.EXPECT().
					FetchUsersByIDs([]int64{123}).
					Return(nil, errors.New("user error"))
			},
			wantErr: true,
			assertFunc: func(t *testing.T, res []model.Listing, err error) {
				assert.Nil(t, res)
				assert.Error(t, err)
			},
		},
		{
			name:   "success with user mapping",
			page:   1,
			size:   10,
			userID: nil,
			mock: func() {
				listingClient.EXPECT().
					FetchListings(1, 10, nil).
					Return([]model.Listing{
						{ID: 1, UserID: 123, Price: 999},
					}, nil)

				userClient.EXPECT().
					FetchUsersByIDs([]int64{123}).
					Return(map[int64]*model.User{
						123: {ID: 123, Name: "John"},
					}, nil)
			},
			wantErr: false,
			assertFunc: func(t *testing.T, res []model.Listing, err error) {
				assert.NoError(t, err)
				assert.Equal(t, int64(1), res[0].ID)
				assert.Equal(t, "John", res[0].User.Name)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			res, err := svc.GetListings(context.Background(), tt.page, tt.size, tt.userID)
			tt.assertFunc(t, res, err)
		})
	}
}
