package services

import (
	"context"
	"testing"
	"time"

	"github.com/AugmentFund/internal/data"
	"github.com/AugmentFund/internal/stores"
	"github.com/stretchr/testify/assert"
)

func TestFundService_CreateTransfer(t *testing.T) {
	ctx := context.Background()
	date := time.Now().Format("2006-01-02 15:04:05")

	mockDB := stores.NewMockDataStore()
	mockDB.CreateUser(ctx, data.User{ID: 1, Name: "John Doe", OwnedFunds: []data.OwnedFund{}})
	mockDB.CreateFund(ctx, data.Fund{
		ID:          1,
		Name:        "test fund",
		TotalShares: 1000,
		OwnedShares: 0,
		Owners:      make(map[int]data.Owner),
	})

	tests := []struct {
		name         string
		transferData data.Transfer
		mockDBSetup  func()
		want         []data.Fund
		wantErr      bool
	}{
		{
			name: "create transfer from fund to owner",
			transferData: data.Transfer{
				FundID:      1,
				FromOwnerID: 0,
				ToOwnerID:   1,
				Shares:      500,
			},
			mockDBSetup: func() {},
			want: []data.Fund{
				{
					ID:          1,
					Name:        "test fund",
					TotalShares: 1000,
					OwnedShares: 500,
					Owners: map[int]data.Owner{
						1: {
							ID:          1,
							Name:        "John Doe",
							TotalShares: 500,
							Date:        date,
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "create transfer some shares from owner to fund",
			transferData: data.Transfer{
				FundID:      1,
				FromOwnerID: 1,
				ToOwnerID:   0,
				Shares:      100,
			},
			mockDBSetup: func() {
				mockDB.UpdateFund(context.Background(), data.Fund{
					ID:          1,
					Name:        "test fund",
					TotalShares: 1000,
					OwnedShares: 500,
					Owners: map[int]data.Owner{
						1: {
							ID:          1,
							Name:        "John Doe",
							TotalShares: 500,
							Date:        date,
						},
					},
				})
			},
			want: []data.Fund{
				{
					ID:          1,
					Name:        "test fund",
					TotalShares: 1000,
					OwnedShares: 400,
					Owners: map[int]data.Owner{
						1: {
							ID:          1,
							Name:        "John Doe",
							TotalShares: 400,
							Date:        date,
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "create transfer all shares from owner to fund",
			transferData: data.Transfer{
				FundID:      1,
				FromOwnerID: 1,
				ToOwnerID:   0,
				Shares:      500,
			},
			mockDBSetup: func() {
				mockDB.UpdateFund(context.Background(), data.Fund{
					ID:          1,
					Name:        "test fund",
					TotalShares: 1000,
					OwnedShares: 500,
					Owners: map[int]data.Owner{
						1: {
							ID:          1,
							Name:        "John Doe",
							TotalShares: 500,
							Date:        date,
						},
					},
				})
			},
			want: []data.Fund{
				{
					ID:          1,
					Name:        "test fund",
					TotalShares: 1000,
					OwnedShares: 0,
					Owners:      map[int]data.Owner{},
				},
			},
			wantErr: false,
		},
		{
			name: "create transfer all shares from fund to owner",
			transferData: data.Transfer{
				FundID:      1,
				FromOwnerID: 0,
				ToOwnerID:   1,
				Shares:      1000,
			},
			mockDBSetup: func() {
				mockDB.UpdateFund(context.Background(), data.Fund{
					ID:          1,
					Name:        "test fund",
					TotalShares: 1000,
					OwnedShares: 0,
					Owners:      map[int]data.Owner{},
				})
			},
			want: []data.Fund{
				{
					ID:          1,
					Name:        "test fund",
					TotalShares: 1000,
					OwnedShares: 1000,
					Owners: map[int]data.Owner{
						1: {
							ID:          1,
							Name:        "John Doe",
							TotalShares: 1000,
							Date:        date,
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "create transfer some shares from owner to new owner",
			transferData: data.Transfer{
				FundID:      1,
				FromOwnerID: 1,
				ToOwnerID:   2,
				Shares:      400,
			},
			mockDBSetup: func() {
				mockDB.CreateUser(ctx, data.User{ID: 2, Name: "Bill Murray", OwnedFunds: []data.OwnedFund{}})
				mockDB.UpdateFund(context.Background(), data.Fund{
					ID:          1,
					Name:        "test fund",
					TotalShares: 1000,
					OwnedShares: 1000,
					Owners: map[int]data.Owner{
						1: {
							ID:          1,
							Name:        "John Doe",
							TotalShares: 1000,
							Date:        date,
						},
					},
				})
			},
			want: []data.Fund{
				{
					ID:          1,
					Name:        "test fund",
					TotalShares: 1000,
					OwnedShares: 1000,
					Owners: map[int]data.Owner{
						1: {
							ID:          1,
							Name:        "John Doe",
							TotalShares: 600,
							Date:        date,
						},
						2: {
							ID:          2,
							Name:        "Bill Murray",
							TotalShares: 400,
							Date:        date,
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "create transfer some shares from owner to existing owner",
			transferData: data.Transfer{
				FundID:      1,
				FromOwnerID: 1,
				ToOwnerID:   2,
				Shares:      200,
			},
			mockDBSetup: func() {
				mockDB.CreateUser(ctx, data.User{ID: 2, Name: "Bill Murray", OwnedFunds: []data.OwnedFund{}})
				mockDB.UpdateFund(context.Background(), data.Fund{
					ID:          1,
					Name:        "test fund",
					TotalShares: 1000,
					OwnedShares: 1000,
					Owners: map[int]data.Owner{
						1: {
							ID:          1,
							Name:        "John Doe",
							TotalShares: 400,
							Date:        date,
						},
						2: {
							ID:          2,
							Name:        "Bill Murray",
							TotalShares: 600,
							Date:        date,
						},
					},
				})
			},
			want: []data.Fund{
				{
					ID:          1,
					Name:        "test fund",
					TotalShares: 1000,
					OwnedShares: 1000,
					Owners: map[int]data.Owner{
						1: {
							ID:          1,
							Name:        "John Doe",
							TotalShares: 200,
							Date:        date,
						},
						2: {
							ID:          2,
							Name:        "Bill Murray",
							TotalShares: 800,
							Date:        date,
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "create transfer all shares from owner to owner",
			transferData: data.Transfer{
				FundID:      1,
				FromOwnerID: 1,
				ToOwnerID:   2,
				Shares:      1000,
			},
			mockDBSetup: func() {
				mockDB.CreateUser(ctx, data.User{ID: 2, Name: "Bill Murray", OwnedFunds: []data.OwnedFund{}})
				mockDB.UpdateFund(context.Background(), data.Fund{
					ID:          1,
					Name:        "test fund",
					TotalShares: 1000,
					OwnedShares: 1000,
					Owners: map[int]data.Owner{
						1: {
							ID:          1,
							Name:        "John Doe",
							TotalShares: 1000,
							Date:        date,
						},
					},
				})
			},
			want: []data.Fund{
				{
					ID:          1,
					Name:        "test fund",
					TotalShares: 1000,
					OwnedShares: 1000,
					Owners: map[int]data.Owner{
						2: {
							ID:          2,
							Name:        "Bill Murray",
							TotalShares: 1000,
							Date:        date,
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "create transfer shares from owner to owner - invalid amount: 0",
			transferData: data.Transfer{
				FundID:      1,
				FromOwnerID: 1,
				ToOwnerID:   2,
				Shares:      0,
			},
			mockDBSetup: func() {
				mockDB.CreateUser(ctx, data.User{ID: 2, Name: "Bill Murray", OwnedFunds: []data.OwnedFund{}})
				mockDB.UpdateFund(context.Background(), data.Fund{
					ID:          1,
					Name:        "test fund",
					TotalShares: 1000,
					OwnedShares: 1000,
					Owners: map[int]data.Owner{
						1: {
							ID:          1,
							Name:        "John Doe",
							TotalShares: 1000,
							Date:        date,
						},
					},
				})
			},
			want:    []data.Fund{},
			wantErr: true,
		},
		{
			name: "create transfer shares from owner to owner - invalid amount: negative",
			transferData: data.Transfer{
				FundID:      1,
				FromOwnerID: 1,
				ToOwnerID:   2,
				Shares:      -100,
			},
			mockDBSetup: func() {
				mockDB.CreateUser(ctx, data.User{ID: 2, Name: "Bill Murray", OwnedFunds: []data.OwnedFund{}})
				mockDB.UpdateFund(context.Background(), data.Fund{
					ID:          1,
					Name:        "test fund",
					TotalShares: 1000,
					OwnedShares: 1000,
					Owners: map[int]data.Owner{
						1: {
							ID:          1,
							Name:        "John Doe",
							TotalShares: 1000,
							Date:        date,
						},
					},
				})
			},
			want:    []data.Fund{},
			wantErr: true,
		},
		{
			name: "create transfer shares from owner to owner - invalid data: same owner transfer",
			transferData: data.Transfer{
				FundID:      1,
				FromOwnerID: 1,
				ToOwnerID:   1,
				Shares:      100,
			},
			mockDBSetup: func() {
				mockDB.CreateUser(ctx, data.User{ID: 2, Name: "Bill Murray", OwnedFunds: []data.OwnedFund{}})
				mockDB.UpdateFund(context.Background(), data.Fund{
					ID:          1,
					Name:        "test fund",
					TotalShares: 1000,
					OwnedShares: 1000,
					Owners: map[int]data.Owner{
						1: {
							ID:          1,
							Name:        "John Doe",
							TotalShares: 1000,
							Date:        date,
						},
					},
				})
			},
			want:    []data.Fund{},
			wantErr: true,
		},
		{
			name: "create transfer shares from owner to owner - invalid data: user doesn't exist",
			transferData: data.Transfer{
				FundID:      1,
				FromOwnerID: 3,
				ToOwnerID:   2,
				Shares:      100,
			},
			mockDBSetup: func() {
				mockDB.CreateUser(ctx, data.User{ID: 2, Name: "Bill Murray", OwnedFunds: []data.OwnedFund{}})
				mockDB.UpdateFund(context.Background(), data.Fund{
					ID:          1,
					Name:        "test fund",
					TotalShares: 1000,
					OwnedShares: 1000,
					Owners: map[int]data.Owner{
						1: {
							ID:          1,
							Name:        "John Doe",
							TotalShares: 1000,
							Date:        date,
						},
					},
				})
			},
			want:    []data.Fund{},
			wantErr: true,
		},
		{
			name: "create transfer shares from owner to owner - invalid amount: more than available",
			transferData: data.Transfer{
				FundID:      1,
				FromOwnerID: 1,
				ToOwnerID:   2,
				Shares:      100000,
			},
			mockDBSetup: func() {
				mockDB.CreateUser(ctx, data.User{ID: 2, Name: "Bill Murray", OwnedFunds: []data.OwnedFund{}})
				mockDB.UpdateFund(context.Background(), data.Fund{
					ID:          1,
					Name:        "test fund",
					TotalShares: 1000,
					OwnedShares: 1000,
					Owners: map[int]data.Owner{
						1: {
							ID:          1,
							Name:        "John Doe",
							TotalShares: 1000,
							Date:        date,
						},
					},
				})
			},
			want:    []data.Fund{},
			wantErr: true,
		},
		{
			name: "create transfer shares from fund to owner - invalid amount: more than available",
			transferData: data.Transfer{
				FundID:      1,
				FromOwnerID: 0,
				ToOwnerID:   2,
				Shares:      100000,
			},
			mockDBSetup: func() {
				mockDB.CreateUser(ctx, data.User{ID: 2, Name: "Bill Murray", OwnedFunds: []data.OwnedFund{}})
				mockDB.UpdateFund(context.Background(), data.Fund{
					ID:          1,
					Name:        "test fund",
					TotalShares: 1000,
					OwnedShares: 1000,
					Owners: map[int]data.Owner{
						1: {
							ID:          1,
							Name:        "John Doe",
							TotalShares: 1000,
							Date:        date,
						},
					},
				})
			},
			want:    []data.Fund{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockDBSetup()
			fundService := NewFundService(mockDB)
			capTable, err := fundService.CreateTransfer(context.Background(), tt.transferData)
			if (err != nil) != tt.wantErr {
				t.Errorf("FundService.CreateTransfer() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr {
				assert.Equal(t, tt.want, capTable)
			}
		})
	}
}
