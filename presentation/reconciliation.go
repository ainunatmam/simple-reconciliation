package presentation

import "fiber.misetaku.net/entity"

type ReconciliationResponse struct {
	TransactionProcessed uint                   `json:"transaction_processed"`
	TransactionMatched   int                    `json:"transaction_matched"`
	Unmatched            ReconciliationUnmached `json:"unmatched_transactions"`
	Discrepancies        float64                `json:"discrepancies"`
}

type ReconciliationUnmached struct {
	Count       int  `json:"count"`
	Transaction []entity.TransactionItem `json:"system_transactions"`
	BankStatement map[string][]entity.StatementItem `json:"bank_statements"`
}