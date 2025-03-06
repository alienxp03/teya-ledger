package transaction

import (
	"github.com/alienxp03/teya-ledger/storage"
	"github.com/alienxp03/teya-ledger/types"
)

type Transactioner interface {
	GetTransactions(userID string, req GetTransactionsRequest) (*GetTransactionsResponse, error)
	CreateDeposit(userID string, req CreateDepositRequest) (*CreateDepositResponse, error)
	CreateWithdrawal(userID string, req CreateWithdrawalRequest) (*CreateWithdrawalResponse, error)
}

type TransactionHandler struct {
	storage storage.Storage
}

func New(storage storage.Storage) *TransactionHandler {
	return &TransactionHandler{
		storage: storage,
	}
}

func (t TransactionHandler) CreateDeposit(userID string, req CreateDepositRequest) (*CreateDepositResponse, error) {
	if _, err := t.storage.GetAccount(userID, req.AccountNumber); err != nil {
		return nil, types.NewNotFound(err.Error())
	}

	transaction, err := t.storage.CreateDeposit(&storage.Transaction{
		TransactionID: req.TransactionID,
		Status:        "pending",
		Amount:        req.Amount,
		Currency:      req.Currency,
		Description:   req.Description,
	})
	if err != nil {
		return nil, err
	}

	return &CreateDepositResponse{Transaction: Transaction{
		TransactionID: transaction.TransactionID,
		Status:        transaction.Status,
		Amount:        transaction.Amount,
		Currency:      transaction.Currency,
		Description:   transaction.Description,
		CreatedAt:     transaction.CreatedAt,
		UpdatedAt:     transaction.UpdatedAt,
	}}, nil
}

func (t TransactionHandler) GetTransactions(userID string, req GetTransactionsRequest) (*GetTransactionsResponse, error) {
	transactionsData, err := t.storage.GetTransactions(userID, req.AccountNumber, req.Limit, req.Page)
	if err != nil {
		return nil, err
	}

	transactions := []Transaction{}
	for _, transaction := range transactionsData {
		transactions = append(transactions, Transaction{
			TransactionID: transaction.TransactionID,
			Status:        transaction.Status,
			Amount:        transaction.Amount,
			Currency:      transaction.Currency,
			Description:   transaction.Description,
			CreatedAt:     transaction.CreatedAt,
			UpdatedAt:     transaction.UpdatedAt,
		})
	}

	return &GetTransactionsResponse{Transactions: transactions}, nil
}

func (t TransactionHandler) CreateWithdrawal(userID string, req CreateWithdrawalRequest) (*CreateWithdrawalResponse, error) {
	if _, err := t.storage.GetAccount(userID, req.AccountNumber); err != nil {
		return nil, types.NewNotFound(err.Error())
	}

	// TODO: Add check balance

	transaction, err := t.storage.CreateWithdrawal(&storage.Transaction{
		TransactionID: req.TransactionID,
		Status:        "pending",
		Amount:        req.Amount,
		Currency:      req.Currency,
		Description:   req.Description,
		UserID:        userID,
		AccountNumber: req.AccountNumber,
	})
	if err != nil {
		return nil, err
	}

	return &CreateWithdrawalResponse{Transaction: Transaction{
		TransactionID: transaction.TransactionID,
		Status:        transaction.Status,
		Amount:        transaction.Amount,
		Currency:      transaction.Currency,
		Description:   transaction.Description,
		CreatedAt:     transaction.CreatedAt,
		UpdatedAt:     transaction.UpdatedAt,
	}}, nil
}
