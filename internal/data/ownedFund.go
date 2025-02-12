package data

type OwnedFund struct {
	ID       int    `json:"id"`
	FundName string `json:"fundName"`
	Shares   int    `json:"shares"`
	Date     string `json:"dateAcquired"`
}
