package models

import (
	"time"
)

type Credit struct {
	ID            int64     `json:"id" db:"id"`
	UserID        int64     `json:"user_id" db:"user_id"`
	AccountID     int64     `json:"account_id" db:"account_id"`
	Amount        float64   `json:"amount" db:"amount"`
	InterestRate  float64   `json:"interest_rate" db:"interest_rate"`
	TermMonths    int       `json:"term_months" db:"term_months"`
	Status        string    `json:"status" db:"status"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time `json:"updated_at" db:"updated_at"`
}

type CreditCreate struct {
	UserID       int64   `json:"user_id" validate:"required"`
	AccountID    int64   `json:"account_id" validate:"required"`
	Amount       float64 `json:"amount" validate:"required,gt=0"`
	TermMonths   int     `json:"term_months" validate:"required,gt=0"`
}

type PaymentSchedule struct {
	ID            int64     `json:"id" db:"id"`
	CreditID      int64     `json:"credit_id" db:"credit_id"`
	PaymentNumber int       `json:"payment_number" db:"payment_number"`
	Amount        float64   `json:"amount" db:"amount"`
	Principal     float64   `json:"principal" db:"principal"`
	Interest      float64   `json:"interest" db:"interest"`
	DueDate       time.Time `json:"due_date" db:"due_date"`
	Status        string    `json:"status" db:"status"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time `json:"updated_at" db:"updated_at"`
} 