package interfaces

import (
	"context"

	"github.com/AugmentFund/internal/data"
)

type DataStorer interface {
	GetUsers(ctx context.Context) ([]data.User, error)
	GetUser(ctx context.Context, id int) (data.User, error)
	CreateUser(ctx context.Context, user data.User) ([]data.User, error)

	GetCapTables(ctx context.Context) ([]data.Fund, error)
	GetCapTableByID(ctx context.Context, id int) (data.Fund, error)
	CreateFund(ctx context.Context, fund data.Fund) ([]data.Fund, error)
	UpdateFund(ctx context.Context, fund data.Fund) ([]data.Fund, error)
}
