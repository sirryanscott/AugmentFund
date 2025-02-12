package data

type Fund struct {
	ID          int           `json:"id"`
	Name        string        `json:"name"`
	TotalShares int           `json:"totalShares"`
	OwnedShares int           `json:"ownedShares"`
	Owners      map[int]Owner `json:"owners"`
}

type Transfer struct {
	FundID      int `json:"fundId"`
	FromOwnerID int `json:"fromOwnerId"`
	ToOwnerID   int `json:"toOwnerId"`
	Shares      int `json:"shares"`
}

type TransferHistory struct {
	ID        int    `json:"id"`
	FundID    int    `json:"fundId"`
	FundName  string `json:"fundName"`
	FromOwner Owner  `json:"fromOwner"`
	ToOwner   Owner  `json:"toOwner"`
	Shares    int    `json:"sharesTransfered"`
	Date      string `json:"transferDate"`
}
