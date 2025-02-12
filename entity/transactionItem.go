package entity

import "time"

type TransactionItem struct {
	TrxId           string    `json:"trx_id"`
	Amount          float32   `json:"amount"`
	Type            string    `json:"type"`
	Reason string `json:"reason"`
	TransactionTime time.Time `json:"transaction_time"`

}