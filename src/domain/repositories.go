package domain

type FileRepository interface {
	GetTransactions(bucket, key string) ([]Txn, error)
}

type EmailRepository interface {
	SendEmail(from, to, subject, body string) error
}
