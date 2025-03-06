package api

import (
	"net/http"
	"sync"

	"github.com/alienxp03/teya-ledger/handler/transaction"
	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

const (
	HeaderUserID = "userID"
)

// Transactioner defines the interface for transaction-related operations
type Transactioner interface {
	GetTransactions(userID string, req transaction.GetTransactionsRequest) (*transaction.GetTransactionsResponse, error)
	CreateDeposit(userID string, req transaction.CreateDepositRequest) (*transaction.CreateDepositResponse, error)
	CreateWithdrawal(userID string, req transaction.CreateWithdrawalRequest) (*transaction.CreateWithdrawalResponse, error)
	GetBalance(userID string, req transaction.GetBalanceRequest) (*transaction.GetBalanceResponse, error)
}

type APIImpl struct {
	transactioner transaction.Transactioner

	once sync.Once
	mux  *http.ServeMux
}
