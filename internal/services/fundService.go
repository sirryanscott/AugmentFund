package services

import (
	"context"

	"github.com/AugmentFund/internal/data"
	"github.com/AugmentFund/internal/interfaces"
)

type FundService struct {
	DataStore interfaces.DataStorer
}

func NewFundService(dataStore interfaces.DataStorer) *FundService {
	return &FundService{DataStore: dataStore}
}

func (s *FundService) GetCapTables(ctx context.Context) ([]data.Fund, error) {
	return s.DataStore.GetCapTables(ctx)
}

func (s *FundService) GetCapTableByID(ctx context.Context, id string) (*data.Fund, error) {
	return s.DataStore.GetCapTableByID(ctx, id)
}

func (s *FundService) CreateFund(ctx context.Context, fund data.Fund) ([]data.Fund, error) {
	return s.DataStore.CreateFund(ctx, fund)
}
