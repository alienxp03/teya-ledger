package storage

import "time"

type Transaction struct {
	ID             int
	IdempotencyKey string
	Status         string
	Amount         int
	Currency       string
	Description    string
	AccountNumber  string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}
