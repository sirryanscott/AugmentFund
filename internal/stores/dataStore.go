package stores

import (
	"context"

	"github.com/AugmentFund/internal/data"
)

// TODO load from db
type DataStore struct {
	Users []data.User
	Funds []data.Fund
}

func NewDataStore() *DataStore {
	return &DataStore{
		Users: []data.User{
			{
				ID:   "1",
				Name: "John Doe",
			},
			{
				ID:   "2",
				Name: "Jane Doe",
			},
		},
		Funds: []data.Fund{
			{
				ID:          "1",
				Name:        "Fund 1",
				TotalShares: 1000,
			},
		},
	}
}

func (s *DataStore) GetUsers(ctx context.Context) ([]data.User, error) {
	if s.Users == nil {
		s.Users = []data.User{}
	}
	return s.Users, nil
}

func (s *DataStore) GetUser(ctx context.Context, id string) (*data.User, error) {
	for _, user := range s.Users {
		if user.ID == id {
			return &user, nil
		}
	}
	return &data.User{}, nil
}

func (s *DataStore) CreateUser(ctx context.Context, user data.User) ([]data.User, error) {
	for _, u := range s.Users {
		if u.ID == user.ID {
			return s.Users, nil
		}
	}
	s.Users = append(s.Users, user)
	return s.Users, nil
}

func (s *DataStore) GetFunds(ctx context.Context) ([]data.Fund, error) {
	return []data.Fund{}, nil
}

func (s *DataStore) GetFund(ctx context.Context, id string) (*data.Fund, error) {
	return &data.Fund{}, nil
}

func (s *DataStore) CreateFund(ctx context.Context, fund data.Fund) ([]data.Fund, error) {
	for _, f := range s.Funds {
		if f.ID == fund.ID {
			return s.Funds, nil
		}
	}
	s.Funds = append(s.Funds, fund)
	return s.Funds, nil
}
