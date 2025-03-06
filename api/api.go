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

type Transactioner interface {
	GetTransactions(w http.ResponseWriter, req GetTransactionsRequest) (GetTransactionsResponse, error)
	CreateDeposit(w http.ResponseWriter, req CreateDepositRequest) (CreateDepositResponse, error)
	CreateWithdrawal(w http.ResponseWriter, req CreateWithdrawalRequest) (CreateWithdrawalResponse, error)
}

type APIImpl struct {
	transactioner transaction.Transactioner

	once sync.Once
	mux  *http.ServeMux
}
