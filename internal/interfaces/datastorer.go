package interfaces

import (
	"context"

	"github.com/AugmentFund/internal/data"
)

type DataStorer interface {
	GetUsers(ctx context.Context) ([]data.User, error)
	GetUser(ctx context.Context, id string) (*data.User, error)
	CreateUser(ctx context.Context, user data.User) ([]data.User, error)

	GetFunds(ctx context.Context) ([]data.Fund, error)
	GetFund(ctx context.Context, id string) (*data.Fund, error)
	CreateFund(ctx context.Context, fund data.Fund) ([]data.Fund, error)
}
