package data

type User struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	OwnedFunds []OwnedFund
}

type Owner struct {
	ID          string
	Name        string
	TotalShares int
}
