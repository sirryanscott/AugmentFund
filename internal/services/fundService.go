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

func (s *FundService) GetFunds(ctx context.Context) ([]data.Fund, error) {
	return s.DataStore.GetFunds(ctx)
}
