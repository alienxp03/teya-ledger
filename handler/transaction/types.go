package transaction

import "time"

type GetTransactionsRequest struct {
	AccountNumber string
	Limit         int
	Page          int
}

type GetTransactionsResponse struct {
	Transactions []Transaction
}

type Transaction struct {
	IdempotencyKey string
	Status         string
	Amount         int
	Currency       string
	Description    string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}
