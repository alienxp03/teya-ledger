package db

import (
	"time"

	"github.com/alienxp03/teya-ledger/storage"
)

type DB interface {
	Initialize() error
	GetStorage() storage.Storage
	SeedData() error
}

type MemoryDB struct {
	storage storage.Storage
}

func NewMemoryStorage() *MemoryDB {
	return &MemoryDB{
		storage: storage.NewMemoryStorage(),
	}
}

func (m *MemoryDB) Initialize() error {
	return nil
}

func (m *MemoryDB) GetStorage() storage.Storage {
	return m.storage
}

func (m *MemoryDB) SeedData() error {
	users := []User{
		{
			ID:          "USER_ID_1",
			Email:       "user1@example.com",
			AccessToken: "USER_TOKEN_1",
			Name:        "User 1",
		},
	}

	for _, _ = range users {
		// m.storage.CreateUser(user)
	}

	accounts := []Account{
		{
			ID:        "ACCOUNT_ID_1",
			Number:    "ACCOUNT_NUMBER_1",
			UserID:    "USER_ID_1",
			Currency:  "USD",
			Balance:   1000,
			CreatedAt: time.Date(2022, 2, 1, 0, 0, 0, 0, time.UTC),
			UpdatedAt: time.Date(2025, 2, 1, 0, 0, 2, 0, time.UTC),
		},
	}

	for _, _ = range accounts {
	}

	transactions := []storage.Transaction{
		{
			ID:             1,
			IdempotencyKey: "123456",
			Status:         "success",
			Amount:         100,
			Currency:       "USD",
			Description:    "Payment for order 123456",
			AccountNumber:  "ACCOUNT_NUMBER_1",
			CreatedAt:      time.Date(2022, 2, 1, 0, 0, 0, 0, time.UTC),
			UpdatedAt:      time.Date(2025, 2, 1, 0, 0, 2, 0, time.UTC),
		},
	}

	for _, transaction := range transactions {
		// m.storage.CreateTransaction(transaction)
		m.storage.CreateTransaction(&transaction)
	}

	return nil
}
