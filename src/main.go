package main

import (
	"fmt"
	"log"
)

const (
	subject = "Summary of account transactions"
)

func processEmail() {
	transactions, err := getFileTransactions()
	if err != nil {
		fmt.Errorf("Fail to get txns")
	}

	// Scan over the transaction to summarize information
	totalBalance, txnsPerMonth, averageDebitAmount, averageCreditAmount := scanTransactions(transactions)
	fmt.Printf("totalBalance -> %v\ntxnsPerMonth -> %v\naverageDebitAmount -> %v\naverageCreditAmount -> %v\n", totalBalance, txnsPerMonth, averageDebitAmount, averageCreditAmount)

}

func main() {
	processEmail()
	if err := sendEmail("saultorres.imt@outlook.com", "saulftg.22@gmail.com", subject, "Testing SES"); err != nil {
		log.Fatalf("Fail to send email %v", err)
	}
	fmt.Print("Email sent\n")
}
