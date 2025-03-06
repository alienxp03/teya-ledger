package storage

import "time"

type User struct {
	ID        int
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Account struct {
	ID        int
	Number    string
	CreatedAt time.Time
	UpdatedAt time.Time
	Balance   int
	UserID    string
}

type Transaction struct {
	ID            int
	TransactionID string
	Status        string
	Amount        int
	Currency      string
	UserID        string
	Description   string
	AccountNumber string
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

type Balance struct {
	UserID        string
	AccountNumber string
	Amount        int64
	Currency      string
}
