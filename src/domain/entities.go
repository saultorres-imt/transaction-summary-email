package domain

import (
	"time"
)

type Txn struct {
	Id     int
	Date   time.Time
	Amount float64
}

type EmailData struct {
	Name                string `json:"name"`
	From                string `json:"from"`
	To                  string `json:"to"`
	TotalBalance        float64
	AverageDebitAmount  float64
	AverageCreditAmount float64
	TxnsPerMonth        map[string][]Txn
	SortedMonths        []string
}
