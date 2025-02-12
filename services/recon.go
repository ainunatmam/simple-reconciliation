package services

import (
	"math"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"fiber.misetaku.net/entity"
	"fiber.misetaku.net/helpers"
	"fiber.misetaku.net/presentation"
	"golang.org/x/exp/slices"
)

type ManualRecon struct {
    transactionMatched []entity.TransactionItem
    transactions []entity.TransactionItem
	countProcess uint
	discrepancies float64
	statements []entity.StatementItem
	StartDate time.Time
	EndDate time.Time
	statetmentUnmatched map[string][]entity.StatementItem
	transactionUnmatched []entity.TransactionItem
	TransactionPath string
	StatementPath string
}

func (m *ManualRecon) Call() presentation.ResponseBase {
	system_transaction := helpers.ReadFile(m.TransactionPath)
	statement_files, err := os.ReadDir(m.StatementPath)
	if err != nil {
		return presentation.ResponseBase{}.Failed(500, err.Error())
	}

	for _, statement_file := range statement_files {
		// Periksa apakah itu file, bukan folder
		if !statement_file.IsDir() {
			filePath := filepath.Join("./files/statements", statement_file.Name())

			content := helpers.ReadFile(filePath)
			m.parseStatementItem(content, statement_file.Name())
		}
	}

    m.parseTransactionItem(system_transaction)
    m.recon()
	m.groupingStatement()
    return presentation.ResponseBase{
        Status: 200,
        Message: "Success",
        Data: presentation.ReconciliationResponse{
			TransactionProcessed: m.countProcess,
			TransactionMatched: len(m.transactionMatched),
            Unmatched: presentation.ReconciliationUnmached{
                Count: len(m.transactionUnmatched),
				Transaction: m.transactionUnmatched,
				BankStatement: m.statetmentUnmatched,
            },
			Discrepancies: m.discrepancies,
        },
    }
}

func (m *ManualRecon) parseTransactionItem(items [][]string) {
	for _, item := range items[1:] { // mulai dari baris kedua jika baris pertama adalah header
		amount, err := strconv.ParseFloat(item[1], 32) 
		if err != nil {
			continue
		}
		
		itemTime,err := time.Parse("2006-01-02 15:04:05", item[3])
		if err != nil{
			continue
		}

		if !itemTime.After(m.StartDate.AddDate(0, 0, -1)) && !itemTime.Before(m.EndDate.AddDate(0, 0, 1)) {
			continue
		}

		transaction := entity.TransactionItem{
			TrxId: item[0],
			Amount: float32(amount),
            Type: item[2],
            TransactionTime: itemTime,
		}

		m.transactions = append(m.transactions, transaction)
	}
}

func (m *ManualRecon) parseStatementItem(items [][]string, filename string) {
	for _, item := range items[1:] { // mulai dari baris kedua jika baris pertama adalah header
		amount, err := strconv.ParseFloat(item[1], 32) 
		if err != nil {
			continue
		}

		itemTime,err := time.Parse("2006-01-02", item[2])
		if err != nil{
			continue
		}

		if !itemTime.After(m.StartDate.AddDate(0, 0, -1)) && !itemTime.Before(m.EndDate.AddDate(0, 0, 1)) {
			continue
		}

		statement := entity.StatementItem{
			Source: filename,
			UniqueIdentifier: item[0],
			Amount: float32(amount),
            Time: item[2],
			Reason: "Missing in system transaction",
		}
		m.statements = append(m.statements, statement)
	}
}

func (m *ManualRecon) recon() {
    for _, transaction:= range m.transactions {
		m.countProcess++
        check := slices.IndexFunc(m.statements, func(s entity.StatementItem) bool {
			transactionDate := transaction.TransactionTime.Format("2006-01-02")
			return s.UniqueIdentifier	 == transaction.TrxId && transactionDate == s.Time
		})
        if check == -1 {
			transaction.Reason = "Missing in statement bank"
            m.transactionUnmatched = append(m.transactionUnmatched, transaction)
            continue
        }

		if transaction.Amount != m.statements[check].Amount {
			m.statements[check].Reason = "Different in amount"
			transaction.Reason = "Different in amount"
			m.discrepancies = m.discrepancies + math.Abs(float64(transaction.Amount) - math.Abs(float64(m.statements[check].Amount)))
			m.transactionUnmatched = append(m.transactionUnmatched, transaction)
			continue
		}
		transaction.Reason = "Matched"
		m.transactionMatched = append(m.transactionMatched, transaction)
		m.statements = append(m.statements[:check], m.statements[check+1:]...)
    }
}

func (m *ManualRecon) groupingStatement() {
	m.statetmentUnmatched = make(map[string][]entity.StatementItem)
	for _, item := range m.statements {
		m.statetmentUnmatched[item.Source] = append(m.statetmentUnmatched[item.Source], item)
	}
}