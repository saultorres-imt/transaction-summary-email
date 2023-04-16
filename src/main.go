package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"sort"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
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

func processEmail(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var emailData EmailData

	err := json.Unmarshal([]byte(request.Body), &emailData)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       fmt.Sprintf("Invalid request boody %v, err %v", request.Body, err),
		}, nil
	}

	bucket := os.Getenv("S3_BUCKET")
	key := os.Getenv("S3_KEY")

	// Upodate DB with name account
	accountID := updateDB(emailData.Name)

	// Get transactions from a CVS file in a S3 Bucket
	transactions, err := getFileTransactions(bucket, key)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       fmt.Sprintf("Fail to read CVS file, bucket %v, key %v, err %v", bucket, key, err),
		}, nil
	}

	// Scan over the transaction to summarize information
	totalBalance, txnsPerMonth, averageDebitAmount, averageCreditAmount := scanTransactions(accountID, transactions)

	emailData.TotalBalance = totalBalance
	emailData.AverageDebitAmount = averageDebitAmount
	emailData.AverageCreditAmount = averageCreditAmount
	emailData.TxnsPerMonth = txnsPerMonth
	emailData.SortedMonths = sortMonthsByDate(txnsPerMonth)

	tmpl := template.Must(template.New("emailTemplate").Parse(htmlTemplate))

	var emailBody strings.Builder

	err = tmpl.Execute(&emailBody, emailData)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       fmt.Sprintf("Fail to process html template, err %v", err),
		}, nil
	}

	// Send email
	if err := sendEmail(emailData.From, emailData.To, subject, emailBody.String()); err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       fmt.Sprintf("Fail to process html template, err %v", err),
		}, nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       "Email sent",
	}, nil
}

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	switch request.HTTPMethod {
	case "GET":
		return getTransactionsHandler(request)
	case "POST":
		return processEmail(request)
	default:
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusMethodNotAllowed,
			Body:       "Method Not Allowed",
		}, nil
	}
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

func init() {
	var err error
	htmlTemplate, err = loadTemplate("emailTemplate.html")
	if err != nil {
		log.Fatalf("fail to process html template: %v", err)
	}
}

func main() {
	lambda.Start(handler)
}
