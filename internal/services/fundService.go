package services

import (
	"context"
	"fmt"

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

func (s *FundService) GetCapTableByID(ctx context.Context, id int) (data.Fund, error) {
	return s.DataStore.GetCapTableByID(ctx, id)
}

func (s *FundService) CreateFund(ctx context.Context, fund data.Fund) ([]data.Fund, error) {
	return s.DataStore.CreateFund(ctx, fund)
}

// TODO:
//
//	clean up this function with helper functions
//	need to create history
//	need to update the ownedFunds on the users
func (s *FundService) CreateTransfer(ctx context.Context, transferData data.Transfer) ([]data.Fund, error) {
	// get the fund
	fund, err := s.DataStore.GetCapTableByID(ctx, transferData.FundID)
	if err != nil {
		return []data.Fund{}, fmt.Errorf("error getting fund")
	}

	if transferData.FromOwnerID == 0 && transferData.ToOwnerID == 0 {
		return []data.Fund{}, fmt.Errorf("incomplete transfer data")
	}

	// update the shares from the transfer data
	fromOwner, ok := fund.Owners[transferData.FromOwnerID]
	if !ok {
		// if there is no from owner, then we are assuming an initial transfer from the fund to an owner
		// check if there is unowned shares
		unownedShares := fund.TotalShares - fund.OwnedShares
		if unownedShares <= 0 {
			return []data.Fund{}, fmt.Errorf("no shares to transfer")
		}

		if unownedShares < transferData.Shares {
			return []data.Fund{}, fmt.Errorf("not enough unowned shares to transfer")
		}

		toOwner, ok := fund.Owners[transferData.ToOwnerID]
		if !ok {
			// get owner from users
			user, err := s.DataStore.GetUser(ctx, transferData.ToOwnerID)
			if err != nil {
				return []data.Fund{}, fmt.Errorf("error getting user")
			}
			toOwner.ID = user.ID
			toOwner.Name = user.Name
		}

		toOwner.TotalShares += transferData.Shares
		fund.Owners[transferData.ToOwnerID] = toOwner

		fund.OwnedShares += transferData.Shares
		return s.DataStore.UpdateFund(ctx, fund)

	}

	// from owner needs to have enough shares to transer
	if fromOwner.TotalShares < transferData.Shares {
		return []data.Fund{}, fmt.Errorf("not enough shares to transfer")
	}

	if transferData.ToOwnerID == 0 {
		// transfer to unonwed shares
		fund.OwnedShares -= transferData.Shares
		fromOwner.TotalShares -= transferData.Shares
		fund.Owners[transferData.FromOwnerID] = fromOwner
		// TODO: save data
		return s.DataStore.UpdateFund(ctx, fund)
	}

	toOwner, ok := fund.Owners[transferData.ToOwnerID]
	if !ok {
		// see if user exists with the same id
		if !ok {
			// get owner from users
			user, err := s.DataStore.GetUser(ctx, transferData.ToOwnerID)
			if err != nil {
				return []data.Fund{}, fmt.Errorf("error getting user")
			}
			toOwner.ID = user.ID
			toOwner.Name = user.Name
		}
	}

	fromOwner.TotalShares -= transferData.Shares
	toOwner.TotalShares += transferData.Shares

	fund.Owners[transferData.FromOwnerID] = fromOwner
	fund.Owners[transferData.ToOwnerID] = toOwner

	// create history record

	// save the new fund data
	return s.DataStore.UpdateFund(ctx, fund)
}
