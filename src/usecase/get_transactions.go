package usecase

import (
	"encoding/json"

	"github.com/saultorres-imt/transaction-summary-email/src/domain"
	"gorm.io/gorm"
)

var db *gorm.DB

type GetTransactionsUsecase struct {
	transactionRepo domain.TransactionRepository
}

func NewGetTransactionsUsecase(transactionRepo domain.TransactionRepository) *GetTransactionsUsecase {
	return &GetTransactionsUsecase{
		transactionRepo: transactionRepo,
	}
}

func (uc *GetTransactionsUsecase) Execute() ([]byte, error) {
	txns := db.Find(uc.transactionRepo.FindAll())

	body, err := json.Marshal(txns)
	if err != nil {
		return nil, err
	}

	return body, nil
}
