package main

import (
	"net/http"
	"os"

	"gorm.io/gorm"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/saultorres-imt/transaction-summary-email/src/infrastructure"
	"github.com/saultorres-imt/transaction-summary-email/src/usecase"
)

var db *gorm.DB

func handler(request events.APIGatewayProxyRequest, processEmailUsecase *usecase.ProcessEmailUsecase, getTransactionsUsecase *usecase.GetTransactionsUsecase) (events.APIGatewayProxyResponse, error) {
	switch request.HTTPMethod {
	case "GET":
		transactions, err := getTransactionsUsecase.Execute()
		if err != nil {
			return events.APIGatewayProxyResponse{
				StatusCode: http.StatusInternalServerError,
				Body:       err.Error(),
			}, nil
		} else {
			return events.APIGatewayProxyResponse{
				StatusCode: http.StatusOK,
				Body:       string(transactions),
				Headers:    map[string]string{"Content-Type": "application/json"},
			}, nil
		}
	case "POST":
		bucket := os.Getenv("S3_BUCKET")
		key := os.Getenv("S3_KEY")
		err := processEmailUsecase.Execute(bucket, key)
		if err != nil {
			return events.APIGatewayProxyResponse{
				StatusCode: http.StatusInternalServerError,
				Body:       err.Error(),
			}, nil
		} else {
			return events.APIGatewayProxyResponse{
				StatusCode: http.StatusOK,
				Body:       string("Email sent"),
				Headers:    map[string]string{"Content-Type": "application/json"},
			}, nil
		}
	default:
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusMethodNotAllowed,
			Body:       "Method Not Allowed",
		}, nil
	}
}

func main() {
	session := session.Must(session.NewSession(&aws.Config{
		Region: aws.String("us-west-2")},
	))

	accountRepo := infrastructure.NewDBAccountRepository(db)
	transactionRepo := infrastructure.NewDBTransactionRepository(db)
	fileRepo := infrastructure.NewS3FileRepository(session)
	emailRepo := infrastructure.NewSESEmailRepository(session)

	processEmailUsecase := usecase.NewProcessEmailUsecase(accountRepo, transactionRepo, fileRepo, emailRepo)
	getTransactionsUsecase := usecase.NewGetTransactionsUsecase(transactionRepo)

	lambda.Start(func(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		return handler(request, processEmailUsecase, getTransactionsUsecase)
	})
}
