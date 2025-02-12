package services

import (
	"context"
	"fmt"
	"sort"
	"time"

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
		err = s.transferFromOwnerToFund(ctx, &fund, transferData)
	}

	if err != nil {
		return []data.Fund{}, err
	}

	// create history record don't forget sorting
	err = s.createTransferHistoryRecord(ctx, fund, transferData)
	if err != nil {
		return []data.Fund{}, err
	}
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
		// remove fund from user
		err = s.removeFundFromUser(ctx, fromOwner.ID, fund.ID)
		if err != nil {
			return err
		}
	}

	toOwner.Date = time.Now().Format("2006-01-02 15:04:05")

	fund.Owners[transferData.ToOwnerID] = toOwner

	err = s.updateUserTransferData(ctx, fromOwner.ID, toOwner.ID, transferData)
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
	toOwner.Date = time.Now().Format("2006-01-02 15:04:05")

	fund.Owners[transferData.ToOwnerID] = toOwner

	fund.OwnedShares += transferData.Shares

	err = s.updateUserTransferData(ctx, 0, toOwner.ID, transferData)
	return
}

func (s *FundService) transferFromOwnerToFund(ctx context.Context, fund *data.Fund, transferData data.Transfer) (err error) {
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

	if fromOwner.TotalShares == 0 {
		delete(fund.Owners, fromOwner.ID)
		// remove fund from user
		err = s.removeFundFromUser(ctx, fromOwner.ID, fund.ID)
		if err != nil {
			return
		}
	} else {
		err = s.updateUserTransferData(ctx, fromOwner.ID, 0, transferData)
	}

	return
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

func (s *FundService) removeFundFromUser(ctx context.Context, userId, fundId int) error {
	user, err := s.DataStore.GetUser(ctx, userId)
	if err != nil {
		return fmt.Errorf("error getting user")
	}

	for i, fund := range user.OwnedFunds {
		if fund.ID == fundId {
			user.OwnedFunds = append(user.OwnedFunds[:i], user.OwnedFunds[i+1:]...)
			break
		}
	}

	_, err = s.DataStore.UpdateUser(ctx, user)
	return err
}

func (s *FundService) addFundToUser(ctx context.Context, user *data.User, fundId int, shares int) error {
	fund, err := s.DataStore.GetCapTableByID(ctx, fundId)
	if err != nil {
		return fmt.Errorf("error getting fund cap table")
	}

	user.OwnedFunds = append(user.OwnedFunds, data.OwnedFund{
		ID:       fund.ID,
		FundName: fund.Name,
		Shares:   shares,
		Date:     time.Now().Format("2006-01-02 15:04:05"),
	})

	return nil
}

func (s *FundService) updateUserTransferData(ctx context.Context, fromUserID, toUserID int, transferData data.Transfer) (err error) {
	var fromUser data.User
	var toUser data.User

	if fromUserID != 0 {
		fromUser, err = s.DataStore.GetUser(ctx, fromUserID)
		if err != nil {
			return fmt.Errorf("error getting from user data")
		}
	}

	if toUserID != 0 {
		toUser, err = s.DataStore.GetUser(ctx, toUserID)
		if err != nil {
			return fmt.Errorf("error getting to user data")
		}
	}

	if !toUser.HasFund(transferData.FundID) && toUserID != 0 {
		err = s.addFundToUser(ctx, &toUser, transferData.FundID, transferData.Shares)
		if err != nil {
			return fmt.Errorf("error adding fund to user data")
		}
	} else {
		updateUsersOwnedFundsData(&fromUser, &toUser, transferData)
	}

	if fromUserID != 0 {
		_, err = s.DataStore.UpdateUser(ctx, fromUser)
		if err != nil {
			return fmt.Errorf("error updating from user data")
		}
	}

	if toUserID != 0 {
		_, err = s.DataStore.UpdateUser(ctx, toUser)
		if err != nil {
			return fmt.Errorf("error updating from user data")
		}
	}

	return nil
}

func updateUsersOwnedFundsData(fromUser, toUser *data.User, transferData data.Transfer) {
	if fromUser.ID != 0 {
		for i, fund := range fromUser.OwnedFunds {
			if fund.ID == transferData.FundID {
				fromUser.OwnedFunds[i].Shares -= transferData.Shares
				fromUser.OwnedFunds[i].Date = time.Now().Format("2006-01-02 15:04:05")
			}
		}
	}

	if toUser.ID != 0 {
		for i, fund := range toUser.OwnedFunds {
			if fund.ID == transferData.FundID {
				toUser.OwnedFunds[i].Shares += transferData.Shares
				toUser.OwnedFunds[i].Date = time.Now().Format("2006-01-02 15:04:05")
			}
		}
	}
}

func (s *FundService) createTransferHistoryRecord(ctx context.Context, fund data.Fund, transferData data.Transfer) error {
	fromOwner := data.Owner{
		ID:   transferData.FromOwnerID,
		Name: fund.Name,
	}

	if fromOwner.ID != 0 {
		owner, err := s.DataStore.GetUser(ctx, fromOwner.ID)
		if err != nil {
			return fmt.Errorf("error getting from owner data")
		}
		fromOwner.Name = owner.Name
	}

	toOwner := data.Owner{
		ID:   transferData.ToOwnerID,
		Name: fund.Name,
	}

	if toOwner.ID != 0 {
		owner, err := s.DataStore.GetUser(ctx, toOwner.ID)
		if err != nil {
			return fmt.Errorf("error getting to owner data")
		}
		toOwner.Name = owner.Name
	}

	transferHistoryRecord := data.TransferHistory{
		FundID:    transferData.FundID,
		FundName:  fund.Name,
		FromOwner: fromOwner,
		ToOwner:   toOwner,
		Shares:    transferData.Shares,
		Date:      time.Now().Format("2006-01-02 15:04:05"),
	}

	return s.DataStore.CreateTransferHistoryRecord(ctx, transferHistoryRecord)
}

func (s *FundService) GetTransferHistoryForFund(ctx context.Context, fundID int) ([]data.TransferHistory, error) {
	transferHistoryRecords, err := s.DataStore.GetTransferHistoryForFund(ctx, fundID)
	if err != nil {
		return []data.TransferHistory{}, err
	}

	sort.Slice(transferHistoryRecords, func(i, j int) bool {
		return transferHistoryRecords[i].Date > transferHistoryRecords[j].Date
	})

	return transferHistoryRecords, nil
}
