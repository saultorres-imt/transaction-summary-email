package main

import (
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
	"text/template"
)

const (
	subject = "Summary of account transactions"
)

var htmlTemplate string

var months = map[string]int{
	"January": 0, "February": 1, "March": 2, "April": 3, "May": 4, "June": 5,
	"July": 6, "August": 7, "September": 8, "October": 9, "November": 10, "December": 11,
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

func processEmail() error {
	var emailData EmailData

	transactions, err := getFileTransactions()
	if err != nil {
		return fmt.Errorf("Fail to get txns %v", err)
	}

	// Scan over the transaction to summarize information
	totalBalance, txnsPerMonth, averageDebitAmount, averageCreditAmount := scanTransactions(transactions)
	fmt.Printf("totalBalance -> %v\ntxnsPerMonth -> %v\naverageDebitAmount -> %v\naverageCreditAmount -> %v\n", totalBalance, txnsPerMonth, averageDebitAmount, averageCreditAmount)

	emailData.TotalBalance = totalBalance
	emailData.AverageDebitAmount = averageDebitAmount
	emailData.AverageCreditAmount = averageCreditAmount
	emailData.TxnsPerMonth = txnsPerMonth
	emailData.SortedMonths = sortMonthsByDate(txnsPerMonth)

	tmpl := template.Must(template.New("emailTemplate").Parse(htmlTemplate))

	var emailBody strings.Builder

	err = tmpl.Execute(&emailBody, emailData)
	if err != nil {
		return fmt.Errorf("Error with html template %v", err)
	}

	if err := sendEmail("saultorres.imt@outlook.com", "saulftg.22@gmail.com", subject, emailBody.String()); err != nil {
		return fmt.Errorf("Fail to send email %v", err)
	}
	return nil
}

func sortMonthsByDate(txnsPerMonth map[string][]Txn) []string {
	sortedMonths := make([]string, 0, len(txnsPerMonth))
	for month := range txnsPerMonth {
		sortedMonths = append(sortedMonths, month)
	}

	// Sort months by date
	sort.Slice(sortedMonths, func(i, j int) bool {
		return months[sortedMonths[i]] < months[sortedMonths[j]]
	})
	return sortedMonths
}

func loadTemplate(name string) (string, error) {
	content, err := os.ReadFile(name)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

func main() {
	var err error
	htmlTemplate, err = loadTemplate("../emailTemplate.html")
	if err != nil {
		log.Fatalf("fail to process html template: %v", err)
	}
	err = processEmail()
	if err != nil {
		log.Fatalf("fail to process email: %v", err)
	}
	fmt.Print("Email sent\n")
}
