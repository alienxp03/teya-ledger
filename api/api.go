package api

import (
	"net/http"

	"github.com/gorilla/schema"
)

var decoder = schema.NewDecoder()

type Transactioner interface {
	GetTransactions(w http.ResponseWriter, req GetTransactionsRequest) (GetTransactionsResponse, error)
}
