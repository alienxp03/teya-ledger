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
			UserID:      "USER_ID_1",
			Email:       "user1@example.com",
			AccessToken: "USER_TOKEN_1",
			Name:        "User 1",
		},
		{
			UserID:      "USER_ID_2",
			Email:       "user2@example.com",
			AccessToken: "USER_TOKEN_2",
			Name:        "User 2",
		},
	}

	for _, _ = range users {
		// m.storage.CreateUser(user)
	}

	accounts := []storage.Account{
		{
			ID:        1,
			Number:    "ACCOUNT_NUMBER_1",
			UserID:    "USER_ID_1",
			Balance:   1000,
			CreatedAt: time.Date(2022, 2, 1, 0, 0, 0, 0, time.UTC).UTC(),
			UpdatedAt: time.Date(2025, 2, 1, 0, 0, 2, 0, time.UTC).UTC(),
		},
		{
			ID:        2,
			Number:    "ACCOUNT_NUMBER_2",
			UserID:    "USER_ID_2",
			Balance:   2000,
			CreatedAt: time.Date(2022, 2, 1, 0, 0, 0, 0, time.UTC).UTC(),
			UpdatedAt: time.Date(2025, 2, 1, 0, 0, 2, 0, time.UTC).UTC(),
		},
	}

	for _, account := range accounts {
		m.storage.CreateAccount(account)
	}

	transactions := []storage.Transaction{
		{
			ID:            1,
			TransactionID: "123456",
			Status:        "success",
			Amount:        100,
			Currency:      "MYR",
			Description:   "Payment for order 123456",
			UserID:        "USER_ID_1",
			AccountNumber: "ACCOUNT_NUMBER_1",
			CreatedAt:     time.Date(2022, 2, 1, 0, 0, 0, 0, time.UTC),
			UpdatedAt:     time.Date(2025, 2, 1, 0, 0, 2, 0, time.UTC),
		},
	}

	for _, transaction := range transactions {
		m.storage.CreateTransaction(&transaction)
	}

	return nil
}
