package stores

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/AugmentFund/internal/data"
)

const (
	usersFilePath = "internal/stores/Data/users.json"
	fundsFilePath = "internal/stores/Data/funds.json"
)

// TODO load from db
type DataStore struct{}

func NewDataStore() *DataStore {
	return &DataStore{}
}

func (s *DataStore) GetUsers(ctx context.Context) ([]data.User, error) {
	users, err := s.loadUsersFromFile()
	if err != nil {
		return []data.User{}, err
	}
	return users, nil
}

func (s *DataStore) GetUser(ctx context.Context, id string) (*data.User, error) {
	users, err := s.loadUsersFromFile()
	if err != nil {
		return &data.User{}, err
	}
	for _, user := range users {
		if user.ID == id {
			return &user, nil
		}
	}
	return &data.User{}, nil
}

func (s *DataStore) CreateUser(ctx context.Context, user data.User) ([]data.User, error) {
	users, err := s.loadUsersFromFile()
	if err != nil {
		return []data.User{}, err
	}
	for _, u := range users {
		if u.ID == user.ID {
			return users, nil
		}
	}
	users = append(users, user)
	return users, nil
}

func (s *DataStore) loadUsersFromFile() ([]data.User, error) {
	currDir, err := os.Getwd()
	if err != nil {
		return []data.User{}, err
	}

	path := filepath.Join(currDir, usersFilePath)

	usersData, err := os.ReadFile(path)
	if err != nil {
		return []data.User{}, err
	}

	var users []data.User
	err = json.Unmarshal(usersData, &users)
	if err != nil {
		return []data.User{}, err
	}

	return users, nil
}

func (s *DataStore) GetCapTables(ctx context.Context) ([]data.Fund, error) {
	funds, err := s.loadFundsFromFile()
	if err != nil {
		return []data.Fund{}, err
	}
	return funds, nil
}

func (s *DataStore) GetCapTableByID(ctx context.Context, id string) (*data.Fund, error) {
	funds, err := s.loadFundsFromFile()
	if err != nil {
		return &data.Fund{}, err
	}
	for _, fund := range funds {
		if fund.ID == id {
			return &fund, nil
		}
	}
	return &data.Fund{}, nil
}

func (s *DataStore) CreateFund(ctx context.Context, fund data.Fund) ([]data.Fund, error) {
	funds, err := s.loadFundsFromFile()
	if err != nil {
		return []data.Fund{}, err
	}
	for _, f := range funds {
		if f.ID == fund.ID {
			return funds, nil
		}
	}
	funds = append(funds, fund)
	return funds, nil
}

func (s *DataStore) loadFundsFromFile() ([]data.Fund, error) {
	currDir, err := os.Getwd()
	if err != nil {
		return []data.Fund{}, err
	}

	path := filepath.Join(currDir, fundsFilePath)

	fundsData, err := os.ReadFile(path)
	if err != nil {
		return []data.Fund{}, err
	}

	var funds []data.Fund
	err = json.Unmarshal(fundsData, &funds)
	if err != nil {
		return []data.Fund{}, err
	}

	return funds, nil
}
