package main

import (
	_ "embed"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/saultorres-imt/transaction-summary-email/src/infrastructure"
	"github.com/saultorres-imt/transaction-summary-email/src/usecase"
)

//go:embed emailTemplate.html
var emailTemplate string

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
		err := processEmailUsecase.Execute(bucket, key, emailTemplate, request)
		if err != nil {
			return events.APIGatewayProxyResponse{
				StatusCode: http.StatusInternalServerError,
				Body:       err.Error(),
			}, nil
		} else {
			return events.APIGatewayProxyResponse{
				StatusCode:      http.StatusOK,
				Body:            string("Email sent"),
				IsBase64Encoded: false,
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
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	ddb := dynamodb.New(sess)

	session := session.Must(session.NewSession(&aws.Config{
		Region: aws.String("us-east-2")},
	))

	transactionRepo := infrastructure.NewDBTransactionRepository(ddb)
	fileRepo := infrastructure.NewS3FileRepository(session)
	emailRepo := infrastructure.NewSESEmailRepository(session)

	processEmailUsecase := usecase.NewProcessEmailUsecase(transactionRepo, fileRepo, emailRepo)
	getTransactionsUsecase := usecase.NewGetTransactionsUsecase(transactionRepo)

	lambda.Start(func(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		return handler(request, processEmailUsecase, getTransactionsUsecase)
	})
}
