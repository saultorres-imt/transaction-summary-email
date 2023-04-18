package infrastructure

import (
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"

	"github.com/saultorres-imt/transaction-summary-email/src/domain"
)

var tableName string

type DBTransactionRepository struct {
	ddb *dynamodb.DynamoDB
}

func NewDBTransactionRepository(ddb *dynamodb.DynamoDB) *DBTransactionRepository {
	return &DBTransactionRepository{ddb: ddb}
}

func (repo *DBTransactionRepository) Create(transaction *domain.DBTxn) error {
	av, err := dynamodbattribute.MarshalMap(transaction)

	if err != nil {
		return err
	}

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(tableName),
	}

	_, err = repo.ddb.PutItem(input)
	if err != nil {
		return err
	}
	return nil
}

func (repo *DBTransactionRepository) FindAll() ([]domain.DBTxn, error) {
	proj := expression.NamesList(expression.Name("Id"), expression.Name("AccountName"), expression.Name("Date"), expression.Name("Amount"))
	expr, err := expression.NewBuilder().WithProjection(proj).Build()
	if err != nil {
		return nil, err
	}

	input := &dynamodb.ScanInput{
		TableName:                aws.String(tableName),
		ExpressionAttributeNames: expr.Names(),
		ProjectionExpression:     expr.Projection(),
	}
	// Execute the scan operation
	result, err := repo.ddb.Scan(input)
	if err != nil {
		return nil, err
	}

	// Convert the DynamoDB items to domain.DBTxn objects
	transactions := make([]domain.DBTxn, 0)
	for _, item := range result.Items {
		txn := domain.DBTxn{}
		err = dynamodbattribute.UnmarshalMap(item, &txn)
		if err != nil {
			return nil, err
		}
		transactions = append(transactions, txn)
	}
	return transactions, nil
}

func init() {
	tableName = os.Getenv("DYNAMO_DB_TABLE")
}
