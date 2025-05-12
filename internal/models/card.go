package models

import (
	"time"
)

type Card struct {
	ID            int64     `json:"id" db:"id"`
	AccountID     int64     `json:"account_id" db:"account_id"`
	Number        string    `json:"number" db:"number"`
	EncryptedCVV  string    `json:"-" db:"encrypted_cvv"`
	ExpiryDate    string    `json:"expiry_date" db:"expiry_date"`
	CardholderName string   `json:"cardholder_name" db:"cardholder_name"`
	IsActive      bool      `json:"is_active" db:"is_active"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time `json:"updated_at" db:"updated_at"`
}

type CardCreate struct {
	AccountID      int64  `json:"account_id" validate:"required"`
	CardholderName string `json:"cardholder_name" validate:"required"`
}

type CardResponse struct {
	ID            int64     `json:"id"`
	AccountID     int64     `json:"account_id"`
	Number        string    `json:"number"`
	ExpiryDate    string    `json:"expiry_date"`
	CardholderName string   `json:"cardholder_name"`
	IsActive      bool      `json:"is_active"`
	CreatedAt     time.Time `json:"created_at"`
} 