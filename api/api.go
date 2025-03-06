package api

import (
	"net/http"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

const (
	HeaderUserID = "userID"
)

type Transactioner interface {
	GetTransactions(w http.ResponseWriter, req GetTransactionsRequest) (GetTransactionsResponse, error)
	CreateDeposit(w http.ResponseWriter, req CreateDepositRequest) (CreateDepositResponse, error)
}
