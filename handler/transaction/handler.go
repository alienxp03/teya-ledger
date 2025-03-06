package transaction

import (
	"github.com/alienxp03/teya-ledger/storage"
	"github.com/alienxp03/teya-ledger/types"
)

// Transactioner defines the interface for transaction-related operations
type Transactioner interface {
	GetTransactions(userID string, req GetTransactionsRequest) (*GetTransactionsResponse, error)
	CreateDeposit(userID string, req CreateDepositRequest) (*CreateDepositResponse, error)
	CreateWithdrawal(userID string, req CreateWithdrawalRequest) (*CreateWithdrawalResponse, error)
	GetBalance(userID string, req GetBalanceRequest) (*GetBalanceResponse, error)
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
		AccountNumber: req.AccountNumber,
		UserID:        userID,
		Status:        "pending",
		Amount:        req.Amount,
		Currency:      req.Currency,
		Description:   req.Description,
	})
	if err != nil {
		return nil, err
	}

	// Update balance
	if err := t.storage.UpdateBalance(userID, req.AccountNumber, int64(req.Amount)); err != nil {
		return nil, types.NewBadRequest(types.BadRequest, err.Error())
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

	// Check balance
	balance, err := t.storage.GetBalance(userID, req.AccountNumber)
	if err != nil {
		return nil, types.NewBadRequest(types.BadRequest, err.Error())
	}

	if balance.Amount < int64(req.Amount) {
		return nil, types.NewBadRequest(types.ErrorCodeInvalidAmount, "insufficient balance")
	}

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

	// Update balance
	if err := t.storage.UpdateBalance(userID, req.AccountNumber, int64(req.Amount)); err != nil {
		return nil, types.NewBadRequest(types.BadRequest, err.Error())
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

// GetBalance retrieves the current balance for an account
func (h *TransactionHandler) GetBalance(userID string, req GetBalanceRequest) (*GetBalanceResponse, error) {
	// Validate that the account belongs to the user
	if _, err := h.storage.GetAccount(userID, req.AccountNumber); err != nil {
		return nil, types.NewNotFound(err.Error())
	}

	// Get balance directly from storage
	balance, err := h.storage.GetBalance(userID, req.AccountNumber)
	if err != nil {
		return nil, types.NewBadRequest(types.BadRequest, err.Error())
	}

	return &GetBalanceResponse{
		Amount:   balance.Amount,
		Currency: balance.Currency,
	}, nil
}
