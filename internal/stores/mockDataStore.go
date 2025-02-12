package stores

import (
	"context"

	"github.com/AugmentFund/internal/data"
)

type MockDataStore struct {
	users                  []data.User
	funds                  []data.Fund
	transferHistoryRecords []data.TransferHistory
}

func NewMockDataStore() *MockDataStore {
	return &MockDataStore{
		users:                  []data.User{},
		funds:                  []data.Fund{},
		transferHistoryRecords: []data.TransferHistory{},
	}
}

func (s *MockDataStore) GetUsers(ctx context.Context) ([]data.User, error) {
	return s.users, nil
}

func (s *MockDataStore) GetUser(ctx context.Context, id int) (data.User, error) {
	for _, user := range s.users {
		if user.ID == id {
			return user, nil
		}
	}
	return data.User{}, nil
}

func (s *MockDataStore) CreateUser(ctx context.Context, user data.User) ([]data.User, error) {
	s.users = append(s.users, user)
	return s.users, nil
}

func (s *MockDataStore) UpdateUser(ctx context.Context, user data.User) (data.User, error) {
	for i, dbUser := range s.users {
		if dbUser.ID == user.ID {
			s.users[i] = user
			return user, nil
		}
	}
	return data.User{}, nil
}

func (s *MockDataStore) GetCapTables(ctx context.Context) ([]data.Fund, error) {
	return []data.Fund{}, nil
}

func (s *MockDataStore) GetCapTableByID(ctx context.Context, id int) (data.Fund, error) {
	for _, fund := range s.funds {
		if fund.ID == id {
			return fund, nil
		}
	}
	return data.Fund{}, nil
}

func (s *MockDataStore) CreateFund(ctx context.Context, fund data.Fund) ([]data.Fund, error) {
	s.funds = append(s.funds, fund)
	return s.funds, nil
}

func (s *MockDataStore) UpdateFund(ctx context.Context, fund data.Fund) ([]data.Fund, error) {
	for i, dbFund := range s.funds {
		if dbFund.ID == fund.ID {
			s.funds[i] = fund
			return s.funds, nil
		}
	}
	return []data.Fund{}, nil
}

func (s *MockDataStore) GetTransferHistoryForFund(ctx context.Context, fundID int) ([]data.TransferHistory, error) {
	transferHistoryRecords := []data.TransferHistory{}
	for _, historyRecord := range s.transferHistoryRecords {
		if historyRecord.FundID == fundID {
			transferHistoryRecords = append(transferHistoryRecords, historyRecord)
		}
	}
	return transferHistoryRecords, nil
}

func (s *MockDataStore) CreateTransferHistoryRecord(ctx context.Context, transferHistoryRecord data.TransferHistory) error {
	s.transferHistoryRecords = append(s.transferHistoryRecords, transferHistoryRecord)
	return nil
}
