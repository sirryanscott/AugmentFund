package data

type Fund struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	TotalShares int     `json:"totalShares"`
	OwnedShares int     `json:"ownedShares"`
	Owners      []Owner `json:"owners"`
}
