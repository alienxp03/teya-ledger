package db

import "time"

type User struct {
	ID          string
	Email       string
	AccessToken string
	Name        string
}

type Account struct {
	ID        string
	Number    string
	UserID    string
	Currency  string
	Balance   int
	CreatedAt time.Time
	UpdatedAt time.Time
}
