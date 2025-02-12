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
//	need to create history
func (s *FundService) CreateTransfer(ctx context.Context, transferData data.Transfer) ([]data.Fund, error) {
	// get the fund
	fund, err := s.DataStore.GetCapTableByID(ctx, transferData.FundID)
	if err != nil {
		return []data.Fund{}, fmt.Errorf("error getting fund")
	}

	// check for invalid transfer data
	err = validateTransferData(transferData)
	if err != nil {
		return []data.Fund{}, err
	}

	switch {
	// transfer between two owners
	case transferData.FromOwnerID != 0 && transferData.ToOwnerID != 0:
		err = s.transferFromOwnerToOwner(ctx, &fund, transferData)
	// transfer from fund to owner
	case transferData.FromOwnerID == 0 && transferData.ToOwnerID != 0:
		err = s.transferFromFundToOwner(ctx, &fund, transferData)
	// transfer from owner to fund
	case transferData.FromOwnerID != 0 && transferData.ToOwnerID == 0:
		err = s.transferFromOwnerToFund(&fund, transferData)
	}

	if err != nil {
		return []data.Fund{}, err
	}

	// create history record don't forget sorting
	// save the new fund data
	return s.DataStore.UpdateFund(ctx, fund)
}

func validateTransferData(transferData data.Transfer) error {
	if transferData.FromOwnerID == transferData.ToOwnerID {
		return fmt.Errorf("transfer data invalid: from owner must not match to owner")
	}
	if transferData.FromOwnerID == 0 && transferData.ToOwnerID == 0 {
		return fmt.Errorf("transfer data invalid: no owners specified in transfer")
	}
	if transferData.Shares <= 0 {
		return fmt.Errorf("transfer data invalid: shares must be positive")
	}
	return nil
}

func (s *FundService) transferFromOwnerToOwner(ctx context.Context, fund *data.Fund, transferData data.Transfer) (err error) {
	fromOwner, ok := fund.Owners[transferData.FromOwnerID]
	if !ok {
		return fmt.Errorf("owner doesn't exist")
	}

	if fromOwner.TotalShares < transferData.Shares {
		return fmt.Errorf("not enough shares to transfer")
	}

	toOwner, ok := fund.Owners[transferData.ToOwnerID]
	if !ok {
		// see if user exists with the same id
		if !ok {
			toOwner, err = s.createOwnerFromUserData(ctx, transferData)
			if err != nil {
				return err
			}
		}
	}

	fromOwner.TotalShares -= transferData.Shares
	toOwner.TotalShares += transferData.Shares

	fund.Owners[transferData.FromOwnerID] = fromOwner
	if fromOwner.TotalShares == 0 {
		delete(fund.Owners, fromOwner.ID)
	}

	fund.Owners[transferData.ToOwnerID] = toOwner
	return
}

func (s *FundService) transferFromFundToOwner(ctx context.Context, fund *data.Fund, transferData data.Transfer) (err error) {
	// if there is no from owner, then we are assuming an initial transfer from the fund to an owner
	// check if there is unowned shares
	unownedShares := fund.TotalShares - fund.OwnedShares
	if unownedShares <= 0 {
		return fmt.Errorf("no shares to transfer")
	}

	if unownedShares < transferData.Shares {
		return fmt.Errorf("not enough unowned shares to transfer")
	}

	toOwner, ok := fund.Owners[transferData.ToOwnerID]
	if !ok {
		// get owner from users
		toOwner, err = s.createOwnerFromUserData(ctx, transferData)
		if err != nil {
			return err
		}
	}

	toOwner.TotalShares += transferData.Shares
	fund.Owners[transferData.ToOwnerID] = toOwner

	fund.OwnedShares += transferData.Shares
	return
}

func (s *FundService) transferFromOwnerToFund(fund *data.Fund, transferData data.Transfer) (err error) {
	fromOwner, ok := fund.Owners[transferData.FromOwnerID]
	if !ok {
		return fmt.Errorf("owner not found")
	}

	if fromOwner.TotalShares < transferData.Shares {
		return fmt.Errorf("not enough shares to transfer")
	}

	fund.OwnedShares -= transferData.Shares
	fromOwner.TotalShares -= transferData.Shares
	fund.Owners[transferData.FromOwnerID] = fromOwner

	return nil
}

func (s *FundService) createOwnerFromUserData(ctx context.Context, transferData data.Transfer) (owner data.Owner, err error) {
	// get owner from users
	user, err := s.DataStore.GetUser(ctx, transferData.ToOwnerID)
	if err != nil {
		return data.Owner{}, fmt.Errorf("error getting user")
	}
	owner.ID = user.ID
	owner.Name = user.Name
	return
}

// TODO if a new user gets some shares in a transfer add the fund to the "ownedFunds"
// TODO if a owner loses all shares, remove the fund from owner's "ownedFunds" array
