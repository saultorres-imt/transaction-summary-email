package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"gorm.io/gorm"
)

type DBAccount struct {
	Id    uint `gorm:"primaryKey"`
	Name  string
	Email string
}

type DBTxn struct {
	ID        uint `gorm:"primaryKey"`
	AccountID uint
	Account   DBAccount
	Date      time.Time
	Amount    float64
}

var db *gorm.DB

func updateDB(name string) uint {
	account := &DBAccount{Name: name}
	db.FirstOrCreate(account, DBAccount{Name: name})
	return account.Id
}

func storeTransaction(txn Txn, accountId uint) {
	transaction := &DBTxn{
		AccountID: accountId,
		Date:      txn.Date,
		Amount:    txn.Amount,
	}
	db.Create(transaction)
}

func getTransactionsHandler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var transactions []DBTxn
	result := db.Find(&transactions)

	if result.Error != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       result.Error.Error(),
		}, nil
	}

	body, err := json.Marshal(transactions)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       "Failed to marshal transactions",
		}, nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       string(body),
		Headers:    map[string]string{"Content-Type": "application/json"},
	}, nil
}

/*
func init() {
	var err error
	dsn := "host=myhost user=myuser password=mypassword dbname=mydbname port=myport sslmode=disable TimeZone=UTC"
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}

	err = createTables()
	if err != nil {
		log.Fatalf("Failed to create tables: %v", err)
	}
}

func createTables() error {
	return db.AutoMigrate(&DBAccount{}, &DBTxn{})
}
*/
