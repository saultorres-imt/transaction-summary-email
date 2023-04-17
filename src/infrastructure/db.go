package infrastructure

import (
	"log"

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
	return db.AutoMigrate(&domain.DBAccount{}, &domain.DBTxn{})
}
