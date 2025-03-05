package api

type GetTransactionsRequest struct {
	AccountNumber string `schema:"accountNumber,required"`
	Limit         int    `schema:"limit"`
	Page          int    `schema:"page"`
}

type GetTransactionsResponse struct {
	Transactions []Transaction `json:"transactions"`
}

type Transaction struct {
	IdempotencyKey string `json:"idempotency_key"`
	Status         string `json:"status"`
	Amount         int    `json:"amount"`
	Currency       string `json:"currency"`
	Description    string `json:"description"`
	CreatedAt      string `json:"created_at"`
	UpdatedAt      string `json:"updated_at"`
}
