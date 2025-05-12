package models

import (
	"time"
)

type Account struct {
	ID        int64     `json:"id" db:"id"`
	UserID    int64     `json:"user_id" db:"user_id"`
	Number    string    `json:"number" db:"number"`
	Balance   float64   `json:"balance" db:"balance"`
	Currency  string    `json:"currency" db:"currency"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type AccountCreate struct {
	UserID   int64  `json:"user_id"`
	Currency string `json:"currency" validate:"required,len=3"`
}

type Transaction struct {
	ID            int64     `json:"id" db:"id"`
	FromAccountID int64     `json:"from_account_id" db:"from_account_id"`
	ToAccountID   int64     `json:"to_account_id" db:"to_account_id"`
	Amount        float64   `json:"amount" db:"amount"`
	Type          string    `json:"type" db:"type"`
	Status        string    `json:"status" db:"status"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time `json:"updated_at" db:"updated_at"`
} 