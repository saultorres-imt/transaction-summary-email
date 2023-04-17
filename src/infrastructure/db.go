package infrastructure

import (
	"log"
	"os"

	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"

	"github.com/saultorres-imt/transaction-summary-email/src/domain"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

type DBAccountRepository struct {
	db *gorm.DB
}

type DBTransactionRepository struct {
	db *gorm.DB
}

func NewDBAccountRepository(db *gorm.DB) *DBAccountRepository {
	return &DBAccountRepository{db: db}
}

func NewDBTransactionRepository(db *gorm.DB) *DBTransactionRepository {
	return &DBTransactionRepository{db: db}
}

func (repo *DBAccountRepository) FirstOrCreate(account *domain.DBAccount) error {
	result := repo.db.FirstOrCreate(account, domain.DBAccount{Name: account.Name})
	return result.Error
}

func (repo *DBTransactionRepository) Create(transaction *domain.DBTxn) error {
	result := repo.db.Create(transaction)
	return result.Error
}

func (repo *DBTransactionRepository) FindAll() ([]domain.DBTxn, error) {
	var transactions []domain.DBTxn
	result := repo.db.Find(&transactions)
	return transactions, result.Error
}

func init() {
	var err error
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	dbname := os.Getenv("DB_NAME")

	username, err := getSecretValue("db-username")
	if err != nil {
		log.Fatalf("Failed to get DB username: %v", err)
	}

	password, err := getSecretValue("db-password")
	if err != nil {
		log.Fatalf("Failed to get DB password: %v", err)
	}

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable TimeZone=UTC", host, port, username, password, dbname)
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
	return db.AutoMigrate(&domain.DBAccount{}, &domain.DBTxn{})
}

func getSecretValue(secretName string) (string, error) {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-west-2")},
	)
	if err != nil {
		return "", err
	}

	svc := secretsmanager.New(sess)
	input := &secretsmanager.GetSecretValueInput{
		SecretId: aws.String(secretName),
	}

	result, err := svc.GetSecretValue(input)
	if err != nil {
		return "", err
	}

	return *result.SecretString, nil
}
