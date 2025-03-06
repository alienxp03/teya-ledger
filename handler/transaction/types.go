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

type CreateDepositRequest struct {
	TransactionID string
	AccountNumber string
	Amount        int
	Currency      string
	Description   string
}

type CreateDepositResponse struct {
	Transaction Transaction
}

type CreateWithdrawalRequest struct {
	TransactionID string
	AccountNumber string
	Amount        int
	Currency      string
	Description   string
}

type CreateWithdrawalResponse struct {
	Transaction Transaction
}

type Transaction struct {
	TransactionID string
	Status        string
	Amount        int
	Currency      string
	Description   string
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

type GetBalanceRequest struct {
	AccountNumber string `json:"accountNumber" validate:"required"`
}

type GetBalanceResponse struct {
	Amount   int64  `json:"amount"`
	Currency string `json:"currency"`
}
