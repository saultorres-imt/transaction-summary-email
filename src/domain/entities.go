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

type DBAccount struct {
	Id    uint `gorm:"primaryKey"`
	Name  string
	Email string
}

type DBTxn struct {
	ID        uint `gorm:"primaryKey"`
	AccountID uint
	Date      time.Time
	Amount    float64
}
