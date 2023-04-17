package usecase

import (
	"fmt"
	"sort"
	"strings"
	"text/template"

	"github.com/saultorres-imt/transaction-summary-email/src/domain"
)

const (
	subject = "Summary of account transactions"
)

var htmlTemplate string

var months = map[string]int{
	"January": 0, "February": 1, "March": 2, "April": 3, "May": 4, "June": 5,
	"July": 6, "August": 7, "September": 8, "October": 9, "November": 10, "December": 11,
}

type ProcessEmailUsecase struct {
	accountRepo     domain.AccountRepository
	transactionRepo domain.TransactionRepository
	fileRepo        domain.FileRepository
	emailRepo       domain.EmailRepository
}

func NewProcessEmailUsecase(accountRepo domain.AccountRepository, transactionRepo domain.TransactionRepository, fileRepo domain.FileRepository, emailRepo domain.EmailRepository) *ProcessEmailUsecase {
	return &ProcessEmailUsecase{
		accountRepo:     accountRepo,
		transactionRepo: transactionRepo,
		fileRepo:        fileRepo,
		emailRepo:       emailRepo,
	}
}

func (uc *ProcessEmailUsecase) Execute(bucket, key, emailTemplate string) error {
	var emailData domain.EmailData

	// Upodate DB with name account
	account := &domain.DBAccount{Name: emailData.Name}
	err := uc.accountRepo.FirstOrCreate(account)
	if err != nil {
		return fmt.Errorf("Failed to create or find account: %v", err)
	}
	accountID := account.Id

	// Get transactions from a CVS file in a S3 Bucket
	transactions, err := uc.fileRepo.GetTransactions(bucket, key)
	if err != nil {
		return fmt.Errorf("Failed to read CSV file, bucket %v, key %v, err %v", bucket, key, err)
	}

	// Scan over the transaction to summarize information
	totalBalance, txnsPerMonth, averageDebitAmount, averageCreditAmount, err := scanTransactions(accountID, transactions, uc.transactionRepo)

	if err != nil {
		return fmt.Errorf("Failed to scan transactions, err %v", err)
	}

	emailData.TotalBalance = totalBalance
	emailData.AverageDebitAmount = averageDebitAmount
	emailData.AverageCreditAmount = averageCreditAmount
	emailData.TxnsPerMonth = txnsPerMonth
	emailData.SortedMonths = sortMonthsByDate(txnsPerMonth)

	tmpl := template.Must(template.New("emailTemplate").Parse(emailTemplate))

	var emailBody strings.Builder

	err = tmpl.Execute(&emailBody, emailData)
	if err != nil {
		return fmt.Errorf("Failed to process HTML template, err %v", err)
	}

	// Send email
	if err := uc.emailRepo.SendEmail(emailData.From, emailData.To, subject, emailBody.String()); err != nil {
		return fmt.Errorf("Failed to send email, err %v", err)
	}

	return nil
}

func sortMonthsByDate(txnsPerMonth map[string][]domain.Txn) []string {
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

func scanTransactions(accountID uint, txns []domain.Txn, transactionRepo domain.TransactionRepository) (float64, map[string][]domain.Txn, float64, float64, error) {
	totalBalance := 0.0
	txnsPerMonth := map[string][]domain.Txn{}
	averageDebitAmount := 0.0
	averageCreditAmount := 0.0
	debitTxns := 0
	creditTxns := 0

	for _, txn := range txns {
		totalBalance += txn.Amount
		month := txn.Date.Format("January")
		txnsPerMonth[month] = append(txnsPerMonth[month], txn)

		if txn.Amount < 0 {
			debitTxns++
			averageDebitAmount += txn.Amount
		} else {
			creditTxns++
			averageCreditAmount += txn.Amount
		}

		dbTxn := &domain.DBTxn{
			ID:        accountID,
			AccountID: uint(txn.Id),
			Date:      txn.Date,
			Amount:    txn.Amount,
		}
		err := transactionRepo.Create(dbTxn)
		if err != nil {
			return 0.0, nil, 0.0, 0.0, err
		}
	}

	return totalBalance, txnsPerMonth, averageDebitAmount / float64(debitTxns), averageCreditAmount / float64(creditTxns), nil
}
