package api

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/alienxp03/teya-ledger/handler/transaction"
)

func New(transactioner transaction.Transactioner) *APIImpl {
	return &APIImpl{
		transactioner: transactioner,
	}
}

func (a *APIImpl) setupRoutes() {
	mux := http.NewServeMux()

	// Deposits
	mux.Handle("POST /api/v1/deposits", AuthMiddleware(http.HandlerFunc(a.createDeposit)))

	// Withdrawals
	mux.Handle("POST /api/v1/withdrawals", AuthMiddleware(http.HandlerFunc(a.createWithdrawal)))

	// Current balance

	// Transaction history
	mux.Handle("GET /api/v1/transactions", AuthMiddleware(http.HandlerFunc(a.getTransactions)))

	a.mux = mux
}

func (a *APIImpl) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.once.Do(a.setupRoutes)
	fmt.Printf("Request received: %s %s\n", r.Method, r.URL.Path)
	a.mux.ServeHTTP(w, r)
}

func (a *APIImpl) createDeposit(w http.ResponseWriter, r *http.Request) {
	userID, _ := r.Context().Value(HeaderUserID).(string)

	params, err := createDepositParams(r)
	if err != nil {
		a.respondError(w, http.StatusBadRequest, err, fmt.Sprintf("Invalid body request %+v", err))
		return
	}

	result, err := a.transactioner.CreateDeposit(userID, *params)
	if err != nil {
		a.respondError(w, http.StatusBadRequest, err, fmt.Sprintf("Failed to create deposit: %+v", err))
		return
	}

	resp := CreateDepositResponse{
		Transaction: Transaction{
			TransactionID: result.Transaction.TransactionID,
			Status:        result.Transaction.Status,
			Amount:        result.Transaction.Amount,
			Currency:      result.Transaction.Currency,
			Description:   result.Transaction.Description,
			CreatedAt:     result.Transaction.CreatedAt.Format(time.RFC3339),
			UpdatedAt:     result.Transaction.UpdatedAt.Format(time.RFC3339),
		},
	}

	a.respond(w, http.StatusOK, resp)
}

func (a *APIImpl) createWithdrawal(w http.ResponseWriter, r *http.Request) {
	userID, _ := r.Context().Value(HeaderUserID).(string)

	params, err := createWithdrawalParams(r)
	if err != nil {
		a.respondError(w, http.StatusBadRequest, err, fmt.Sprintf("Invalid body request %+v", err))
		return
	}

	result, err := a.transactioner.CreateWithdrawal(userID, *params)
	if err != nil {
		a.respondError(w, http.StatusBadRequest, err, fmt.Sprintf("Failed to create withdrawal: %+v", err))
		return
	}

	resp := CreateWithdrawalResponse{
		Transaction: Transaction{
			TransactionID: result.Transaction.TransactionID,
			Status:        result.Transaction.Status,
			Amount:        result.Transaction.Amount,
			Currency:      result.Transaction.Currency,
			Description:   result.Transaction.Description,
			CreatedAt:     result.Transaction.CreatedAt.Format(time.RFC3339),
			UpdatedAt:     result.Transaction.UpdatedAt.Format(time.RFC3339),
		},
	}

	a.respond(w, http.StatusOK, resp)
}

func (a *APIImpl) getTransactions(w http.ResponseWriter, r *http.Request) {
	userID, _ := r.Context().Value(HeaderUserID).(string)
	params := getTransactionsParams(r)

	result, err := a.transactioner.GetTransactions(userID, *params)
	if err != nil {
		a.respondError(w, http.StatusBadRequest, err, fmt.Sprintf("Failed to get transactions %+v", err))
		return
	}

	resp := GetTransactionsResponse{}
	transactions := []Transaction{}
	for _, transaction := range result.Transactions {
		transactions = append(transactions, Transaction{
			TransactionID: transaction.TransactionID,
			Status:        transaction.Status,
			Amount:        transaction.Amount,
			Currency:      transaction.Currency,
			Description:   transaction.Description,
			CreatedAt:     transaction.CreatedAt.Format(time.RFC3339),
			UpdatedAt:     transaction.UpdatedAt.Format(time.RFC3339),
		})
	}

	resp.Transactions = transactions

	a.respond(w, http.StatusOK, resp)
}

func getTransactionsParams(r *http.Request) *transaction.GetTransactionsRequest {
	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil {
		limit = 0
	}

	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil {
		page = 0
	}

	req := &transaction.GetTransactionsRequest{
		AccountNumber: r.URL.Query().Get("accountNumber"),
		Limit:         limit,
		Page:          page,
	}

	return req
}

func createDepositParams(r *http.Request) (*transaction.CreateDepositRequest, error) {
	var req CreateDepositRequest
	if err := parseBody(r, &req); err != nil {
		return nil, err
	}

	result := &transaction.CreateDepositRequest{
		TransactionID: req.TransactionID,
		AccountNumber: req.AccountNumber,
		Amount:        req.Amount,
		Currency:      req.Currency,
		Description:   req.Description,
	}
	return result, nil
}

func createWithdrawalParams(r *http.Request) (*transaction.CreateWithdrawalRequest, error) {
	var req CreateWithdrawalRequest
	if err := parseBody(r, &req); err != nil {
		return nil, err
	}

	result := &transaction.CreateWithdrawalRequest{
		TransactionID: req.TransactionID,
		AccountNumber: req.AccountNumber,
		Amount:        req.Amount,
		Currency:      req.Currency,
		Description:   req.Description,
	}
	return result, nil
}
