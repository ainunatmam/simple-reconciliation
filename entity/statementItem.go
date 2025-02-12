package entity

type StatementItem struct {
	Source           string  `json:"source"`
	UniqueIdentifier string  `json:"unique_identifier"`
	Amount           float32 `json:"amount"`
	Time             string  `json:"time"`
	Reason           string  `json:"reason"`
}