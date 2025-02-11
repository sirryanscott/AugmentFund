package data

type User struct {
	ID         int         `json:"id"`
	Name       string      `json:"name"`
	OwnedFunds []OwnedFund `json:"ownedFunds"`
}

type Owner struct {
	ID          int
	Name        string
	TotalShares int
}
