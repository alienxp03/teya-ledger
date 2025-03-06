package api

type GetTransactionsRequest struct {
	AccountNumber string `json:"accountNumber"`
	Limit         string `json:"limit"`
	Page          string `json:"page"`
}

type GetTransactionsResponse struct {
	Transactions []Transaction `json:"transactions"`
}

type CreateDepositRequest struct {
	TransactionID string `validate:"required"`
	AccountNumber string `validate:"required"`
	Amount        int64  `validate:"required,gte=0"`
	Currency      string `validate:"required,oneof=MYR"`
	Description   string `validate:"required"`
}

type CreateDepositResponse struct {
	Transaction Transaction `json:"transaction"`
}

type CreateWithdrawalRequest struct {
	TransactionID string `validate:"required"`
	AccountNumber string `validate:"required"`
	Amount        int64  `validate:"required,lte=0"`
	Currency      string `validate:"required,oneof=MYR"`
	Description   string `validate:"required"`
}

type CreateWithdrawalResponse struct {
	Transaction Transaction `json:"transaction"`
}

type GetBalanceRequest struct {
	AccountNumber string `json:"accountNumber" validate:"required"`
}

type GetBalanceResponse struct {
	Balance Balance `json:"balance"`
}

type GetTransactionRequest struct {
	TransactionID string `json:"transactionID" validate:"required"`
}

type GetTransactionResponse struct {
	Transaction Transaction `json:"transaction"`
}

type Balance struct {
	Amount   int64  `json:"amount"`
	Currency string `json:"currency"`
}

type Transaction struct {
	TransactionID string `json:"transactionID"`
	Status        string `json:"status"`
	Amount        int64  `json:"amount"`
	Currency      string `json:"currency"`
	Description   string `json:"description"`
	CreatedAt     string `json:"createdAt"`
	UpdatedAt     string `json:"updatedAt"`
}
