package repository

import (
	"context"
	"bank-api/internal/models"
)

type UserRepository interface {
	Create(ctx context.Context, user *models.User) error
	GetByID(ctx context.Context, id int64) (*models.User, error)
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	GetByUsername(ctx context.Context, username string) (*models.User, error)
	Update(ctx context.Context, user *models.User) error
	Delete(ctx context.Context, id int64) error
}

type AccountRepository interface {
	Create(ctx context.Context, account *models.Account) error
	GetByID(ctx context.Context, id int64) (*models.Account, error)
	GetByUserID(ctx context.Context, userID int64) ([]*models.Account, error)
	Update(ctx context.Context, account *models.Account) error
	Delete(ctx context.Context, id int64) error
	CreateTransaction(ctx context.Context, transaction *models.Transaction) error
	GetTransactions(ctx context.Context, accountID int64) ([]*models.Transaction, error)
	UpdateBalance(ctx context.Context, accountID int64, amount float64) error
}

type CardRepository interface {
	Create(ctx context.Context, card *models.Card) error
	GetByID(ctx context.Context, id int64) (*models.Card, error)
	GetByAccountID(ctx context.Context, accountID int64) ([]*models.Card, error)
	Update(ctx context.Context, card *models.Card) error
	Delete(ctx context.Context, id int64) error
	Deactivate(ctx context.Context, id int64) error
}

type CreditRepository interface {
	Create(ctx context.Context, credit *models.Credit) error
	GetByID(ctx context.Context, id int64) (*models.Credit, error)
	GetByUserID(ctx context.Context, userID int64) ([]*models.Credit, error)
	Update(ctx context.Context, credit *models.Credit) error
	Delete(ctx context.Context, id int64) error
	CreatePaymentSchedule(ctx context.Context, schedule *models.PaymentSchedule) error
	GetPaymentSchedule(ctx context.Context, creditID int64) ([]*models.PaymentSchedule, error)
	UpdatePaymentStatus(ctx context.Context, paymentID int64, status string) error
} 