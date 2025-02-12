package services

import (
	"context"
	"testing"

	"github.com/AugmentFund/internal/data"
	"github.com/AugmentFund/internal/stores"
	"github.com/stretchr/testify/assert"
)

func TestUserService_UserCrudOperations(t *testing.T) {
	ctx := context.Background()
	mockDB := stores.NewMockDataStore()
	mockDB.CreateUser(ctx, data.User{ID: 1, Name: "John Doe", OwnedFunds: []data.OwnedFund{}})

	tests := []struct {
		name    string
		userID  int
		mockDB  *stores.MockDataStore
		want    data.User
		wantErr bool
	}{
		{
			name:   "Get User",
			userID: 1,
			mockDB: mockDB,
			want: data.User{
				ID:         1,
				Name:       "John Doe",
				OwnedFunds: []data.OwnedFund{},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userService := NewUserService(mockDB)
			userResponse, err := userService.GetUser(context.Background(), tt.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("UserService.GetUser() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr {
				assert.Equal(t, tt.want, userResponse)
			}
		})
	}
}
