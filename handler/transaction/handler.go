package transaction

import (
	"net/http"

	"github.com/alienxp03/teya-ledger/storage"
)

type Transactioner interface {
	GetTransactions(w http.ResponseWriter, req GetTransactionsRequest) (GetTransactionsResponse, error)
}

type TransactionHandler struct {
	storage storage.Storage
}

func New(storage storage.Storage) *TransactionHandler {
	return &TransactionHandler{
		storage: storage,
	}
}

func (t TransactionHandler) GetTransactions(w http.ResponseWriter, req GetTransactionsRequest) (GetTransactionsResponse, error) {
	transactionsData, err := t.storage.GetTransactions(req.AccountNumber, req.Limit, req.Page)
	if err != nil {
		return GetTransactionsResponse{Transactions: []Transaction{}}, err
	}

	transactions := []Transaction{}
	for _, transaction := range transactionsData {
		transactions = append(transactions, Transaction{
			IdempotencyKey: transaction.IdempotencyKey,
			Status:         transaction.Status,
			Amount:         transaction.Amount,
			Currency:       transaction.Currency,
			Description:    transaction.Description,
			CreatedAt:      transaction.CreatedAt,
			UpdatedAt:      transaction.UpdatedAt,
		})
	}

	return GetTransactionsResponse{Transactions: transactions}, nil
}
