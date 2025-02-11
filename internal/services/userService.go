package services

import (
	"context"

	"github.com/AugmentFund/internal/data"
	"github.com/AugmentFund/internal/interfaces"
)

type UserService struct {
	DataStore interfaces.DataStorer
}

func NewUserService(dataStore interfaces.DataStorer) *UserService {
	return &UserService{DataStore: dataStore}
}

func (s *UserService) GetUsers(ctx context.Context) ([]data.User, error) {
	return s.DataStore.GetUsers(ctx)
}

func (s *UserService) GetUser(ctx context.Context, id int) (data.User, error) {
	return s.DataStore.GetUser(ctx, id)
}

func (s *UserService) CreateUser(ctx context.Context, user data.User) ([]data.User, error) {
	return s.DataStore.CreateUser(ctx, user)
}
