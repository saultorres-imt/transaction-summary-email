package infrastructure

import (
	"encoding/csv"
	"fmt"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/saultorres-imt/transaction-summary-email/src/domain"
)

type S3FileRepository struct {
	session *session.Session
}

func NewS3FileRepository(sess *session.Session) *S3FileRepository {
	return &S3FileRepository{session: sess}
}

func (repo *S3FileRepository) GetTransactions(bucket, key string) ([]domain.Txn, error) {
	txns := []domain.Txn{}

	svc := s3.New(repo.session)

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

		txn := domain.Txn{Id: id, Date: date, Amount: amount}
		txns = append(txns, txn)
	}

	return txns, nil
}
