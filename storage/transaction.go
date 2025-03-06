package storage

import (
	"errors"
	"time"
)

func (m *MemoryStorage) CreateDeposit(transaction *Transaction) (*Transaction, error) {
	// Ideally should be handled by a unique constraint
	for _, transactionData := range m.transactions {
		if transactionData.TransactionID == transaction.TransactionID {
			return nil, errors.New("transaction already exists")
		}
	}

	now := time.Now()
	transaction.CreatedAt = now
	transaction.UpdatedAt = now

	m.transactions = append(m.transactions, transaction)

	return transaction, nil
}

func (m *MemoryStorage) GetTransactions(userID, accountNumber string, limit, page int) ([]*Transaction, error) {
	if limit <= 0 {
		limit = 10
	}

	result := []*Transaction{}

	for _, transaction := range m.transactions {
		if transaction.UserID == userID && transaction.AccountNumber == accountNumber {
			result = append(result, transaction)
		}
	}

	return result, nil
}

// CreateTransaction creates a new transaction
func (m *MemoryStorage) CreateTransaction(transaction *Transaction) error {
	// Ideally should be handled by a unique constraint
	for _, transactionData := range m.transactions {
		if transactionData.ID == transaction.ID {
			return errors.New("transaction already exists")
		}
	}

	// Set timestamps
	now := time.Now()
	transaction.CreatedAt = now
	transaction.UpdatedAt = now

	m.transactions = append(m.transactions, transaction)
	return nil
}

func (m *MemoryStorage) CreateWithdrawal(transaction *Transaction) (*Transaction, error) {
	// Ideally should be handled by a unique constraint
	for _, transactionData := range m.transactions {
		if transactionData.TransactionID == transaction.TransactionID {
			return nil, errors.New("transaction already exists")
		}
	}

	now := time.Now()
	transaction.CreatedAt = now
	transaction.UpdatedAt = now

	m.transactions = append(m.transactions, transaction)

	return transaction, nil
}
