package api

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/alienxp03/teya-ledger/handler/transaction"
)

// APIImpl provides the REST endpoints for the application.
type APIImpl struct {
	transactioner transaction.Transactioner

	once sync.Once
	mux  *http.ServeMux
}

func New(transactioner transaction.Transactioner) *APIImpl {
	return &APIImpl{
		transactioner: transactioner,
	}
}

func (a *APIImpl) setupRoutes() {
	mux := http.NewServeMux()

	// Deposits
	// Withdrawals
	// Current balance

	// Transaction history
	mux.HandleFunc("GET /api/v1/transactions", a.getTransactions)

	a.mux = mux
}

func (a *APIImpl) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.once.Do(a.setupRoutes)
	fmt.Printf("Request received: %s %s\n", r.Method, r.URL.Path)
	a.mux.ServeHTTP(w, r)
}

func (a *APIImpl) getTransactions(w http.ResponseWriter, r *http.Request) {
	params, err := getTransactionsParams(r)
	if err != nil {
		a.respondError(w, http.StatusBadRequest, err, fmt.Sprintf("Invalid body request %+v", err))
		return
	}

	result, err := a.transactioner.GetTransactions(w, params)
	if err != nil {
		a.respondError(w, http.StatusBadRequest, err, fmt.Sprintf("Failed to get transactions %+v", err))
		return
	}

	resp := GetTransactionsResponse{}
	transactions := []Transaction{}
	for _, transaction := range result.Transactions {
		transactions = append(transactions, Transaction{
			IdempotencyKey: transaction.IdempotencyKey,
			Status:         transaction.Status,
			Amount:         transaction.Amount,
			Currency:       transaction.Currency,
			Description:    transaction.Description,
			CreatedAt:      transaction.CreatedAt.Format(time.RFC3339),
			UpdatedAt:      transaction.UpdatedAt.Format(time.RFC3339),
		})
	}

	resp.Transactions = transactions

	a.respond(w, http.StatusOK, resp)
}

func getTransactionsParams(r *http.Request) (transaction.GetTransactionsRequest, error) {
	var req transaction.GetTransactionsRequest
	if err := parseRequest(r, &req); err != nil {
		return transaction.GetTransactionsRequest{}, err
	}
	return req, nil
}
