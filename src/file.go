package main

import (
	"encoding/csv"
	"fmt"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

type Txn struct {
	Id     int
	Date   time.Time
	Amount float64
}

func getFileTransactions(bucket, key string) ([]Txn, error) {
	txns := []Txn{}

	session := session.Must(session.NewSession(&aws.Config{
		Region: aws.String("us-west-2")},
	))

	svc := s3.New(session)

	obj, err := svc.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})

	if err != nil {
		return nil, err
	}

	defer obj.Body.Close()

	reader := csv.NewReader(obj.Body)

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

func scanTransactions(accountID uint, txns []Txn) (float64, map[string][]Txn, float64, float64) {
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

		// Sotre the transaction in database
		storeTransaction(txn, accountID)
	}

	return totalBalance, txnsPerMonth, averageDebitAmount / float64(debitTxns), averageCreditAmount / float64(creditTxns)
}
