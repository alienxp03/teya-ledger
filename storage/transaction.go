package storage

import (
	"errors"
	"time"
)

func (m *MemoryStorage) GetTransactions(accountNumber string, limit, page int) ([]*Transaction, error) {
	if limit <= 0 {
		limit = 10
	}

	result := []*Transaction{}

	for _, transaction := range m.transactions {
		if transaction.AccountNumber == accountNumber {
			result = append(result, transaction)
		}
	}

	return result, nil
}

// CreateTransaction creates a new transaction
func (m *MemoryStorage) CreateTransaction(transaction *Transaction) error {
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
