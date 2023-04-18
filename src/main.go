package main

import (
	_ "embed"
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

//go:embed emailTemplate.html
var emailTemplate string

func handler(request events.APIGatewayProxyRequest, processEmailUsecase *usecase.ProcessEmailUsecase) (events.APIGatewayProxyResponse, error) {
	switch request.HTTPMethod {
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
	session := session.Must(session.NewSession(&aws.Config{
		Region: aws.String("us-east-1")},
	))
	fileRepo := infrastructure.NewS3FileRepository(session)
	emailRepo := infrastructure.NewSESEmailRepository(session)
	processEmailUsecase := usecase.NewProcessEmailUsecase(fileRepo, emailRepo)

	lambda.Start(func(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		return handler(request, processEmailUsecase)
	})
}
