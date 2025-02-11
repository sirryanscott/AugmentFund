package data

type Fund struct {
	ID          string           `json:"id"`
	Name        string           `json:"name"`
	TotalShares int              `json:"totalShares"`
	OwnedShares int              `json:"ownedShares"`
	Owners      map[string]Owner `json:"owners"`
}

type Transfer struct {
	FundID      string `json:"fundId"`
	FromOwnerID string `json:"fromOwnerId"`
	ToOwnerID   string `json:"toOwnerId"`
	Shares      int    `json:"shares"`
}
