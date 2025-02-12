package data

type User struct {
	ID         int         `json:"id"`
	Name       string      `json:"name"`
	OwnedFunds []OwnedFund `json:"ownedFunds"`
}

type Owner struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	TotalShares int    `json:"totalShares,omitempty"`
	Date        string `json:"dateAcquired,omitempty"`
}

func (u *User) HasFund(fundID int) bool {
	for _, fund := range u.OwnedFunds {
		if fundID == fund.ID {
			return true
		}
	}
	return false
}
