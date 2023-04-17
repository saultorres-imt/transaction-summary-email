package domain

type AccountRepository interface {
	FirstOrCreate(account *DBAccount) error
}

type TransactionRepository interface {
	Create(transaction *DBTxn) error
	FindAll() ([]DBTxn, error)
}

type FileRepository interface {
	GetTransactions(bucket, key string) ([]Txn, error)
}

type EmailRepository interface {
	SendEmail(from, to, subject, body string) error
}
