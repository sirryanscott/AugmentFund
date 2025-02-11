package stores

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strconv"

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
	users, err := loadDataFromFile[[]data.User](usersFilePath, []data.User{})
	if err != nil {
		return []data.User{}, err
	}

	return users, nil
}

func (s *DataStore) GetUser(ctx context.Context, id string) (*data.User, error) {
	users, err := loadDataFromFile[[]data.User](usersFilePath, []data.User{})
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
	users, err := loadDataFromFile[[]data.User](usersFilePath, []data.User{})
	if err != nil {
		return []data.User{}, err
	}

	id := strconv.Itoa(len(users) + 1)
	user.ID = id

	users = append(users, user)

	err = s.saveDataToFile(usersFilePath, users)
	if err != nil {
		return []data.User{}, err
	}

	return users, nil
}

func (s *DataStore) GetCapTables(ctx context.Context) ([]data.Fund, error) {
	funds, err := loadDataFromFile[[]data.Fund](fundsFilePath, []data.Fund{})
	if err != nil {
		return []data.Fund{}, err
	}

	return funds, nil
}

func (s *DataStore) GetCapTableByID(ctx context.Context, id string) (*data.Fund, error) {
	funds, err := loadDataFromFile[[]data.Fund](fundsFilePath, []data.Fund{})
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
	funds, err := loadDataFromFile[[]data.Fund](fundsFilePath, []data.Fund{})
	if err != nil {
		return []data.Fund{}, err
	}

	id := strconv.Itoa(len(funds) + 1)
	fund.ID = id

	funds = append(funds, fund)

	err = s.saveDataToFile(fundsFilePath, funds)
	if err != nil {
		return []data.Fund{}, err
	}

	return funds, nil
}

func loadDataFromFile[T any](filePath string, data T) (T, error) {
	currDir, err := os.Getwd()
	if err != nil {
		return data, err
	}

	path := filepath.Join(currDir, filePath)

	dataBytes, err := os.ReadFile(path)
	if err != nil {
		return data, err
	}

	err = json.Unmarshal(dataBytes, &data)
	if err != nil {
		return data, err
	}

	return data, nil
}

func (s *DataStore) saveDataToFile(filePath string, data interface{}) error {
	currDir, err := os.Getwd()
	if err != nil {
		return err
	}

	path := filepath.Join(currDir, filePath)

	dataBytes, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}

	// overwrite file
	err = os.WriteFile(path, dataBytes, 0644) // set normal write permissions
	if err != nil {
		return err
	}

	return nil
}
