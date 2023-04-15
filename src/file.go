package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"time"
)

type Txn struct {
	Id     int
	Date   time.Time
	Amount float64
}

func getFileTransactions() ([]Txn, error) {
	txns := []Txn{}

	// Open the file
	file, err := os.Open("../mock-data/txns.csv")
	if err != nil {
		return nil, fmt.Errorf("Could not open file, %v", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)

	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("Could not read a file, %v", err)
	}

	for i, record := range records {
		if i == 0 {
			continue
		}

		id, err := strconv.Atoi(record[0])
		if err != nil {
			return nil, fmt.Errorf("Error converting id, %v", err)
		}

		date, err := time.Parse("1/2", record[1])
		if err != nil {
			return nil, fmt.Errorf("Error converting date, %v", err)
		}

		amount, err := strconv.ParseFloat(record[2], 64)
		if err != nil {
			return nil, fmt.Errorf("Error converting amount, %v", err)
		}

		txn := Txn{id, date, amount}
		txns = append(txns, txn)
	}

	return txns, nil
}

func scanTransactions(txns []Txn) (float64, map[string][]Txn, float64, float64) {
	totalBalance := 0.0
	txnsPerMonth := map[string][]Txn{}
	averageDebitAmount := 0.0
	averageCreditAmount := 0.0
	debitTxns := 0
	creditTxns := 0

	for _, txn := range txns {
		totalBalance += txn.Amount
		month := txn.Date.Format("January")
		txnsPerMonth[month] = append(txnsPerMonth[month], txn)

		if txn.Amount < 0 {
			debitTxns++
			averageDebitAmount += txn.Amount
		} else {
			creditTxns++
			averageCreditAmount += txn.Amount
		}
	}

	return totalBalance, txnsPerMonth, averageDebitAmount / float64(debitTxns), averageCreditAmount / float64(creditTxns)
}
