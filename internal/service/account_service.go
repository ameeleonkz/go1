package service

import (
	"context"
	"bank-api/internal/models"
	"bank-api/internal/repository"
	"errors"
	"math/rand"
)

type AccountService struct {
	repo repository.AccountRepository
}

func NewAccountService(repo repository.AccountRepository) *AccountService {
	return &AccountService{repo: repo}
}

func (s *AccountService) Create(ctx context.Context, input models.AccountCreate) (*models.Account, error) {
	account := &models.Account{
		UserID:   input.UserID,
		Currency: input.Currency,
		Balance:  0,
		Number:   generateAccountNumber(),
	}

	if err := s.repo.Create(ctx, account); err != nil {
		return nil, err
	}

	return account, nil
}

func (s *AccountService) GetByID(ctx context.Context, id int64) (*models.Account, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *AccountService) GetByUserID(ctx context.Context, userID int64) ([]*models.Account, error) {
	return s.repo.GetByUserID(ctx, userID)
}

func (s *AccountService) Transfer(ctx context.Context, fromAccountID, toAccountID int64, amount float64) error {
	if amount <= 0 {
		return errors.New("invalid amount")
	}

	fromAccount, err := s.repo.GetByID(ctx, fromAccountID)
	if err != nil {
		return err
	}

	if fromAccount.Balance < amount {
		return errors.New("insufficient funds")
	}

	// Проверяем существование счета получателя
	if _, err := s.repo.GetByID(ctx, toAccountID); err != nil {
		return errors.New("recipient account not found")
	}

	// Обновляем балансы
	if err := s.repo.UpdateBalance(ctx, fromAccountID, -amount); err != nil {
		return err
	}

	if err := s.repo.UpdateBalance(ctx, toAccountID, amount); err != nil {
		// В случае ошибки возвращаем средства на первый счет
		_ = s.repo.UpdateBalance(ctx, fromAccountID, amount)
		return err
	}

	// Создаем транзакцию
	transaction := &models.Transaction{
		FromAccountID: fromAccountID,
		ToAccountID:   toAccountID,
		Amount:        amount,
		Type:          "transfer",
		Status:        "completed",
	}

	return s.repo.CreateTransaction(ctx, transaction)
}

func (s *AccountService) GetTransactions(ctx context.Context, accountID int64) ([]*models.Transaction, error) {
	return s.repo.GetTransactions(ctx, accountID)
}

func (s *AccountService) UpdateBalance(ctx context.Context, accountID int64, amount float64) error {
	return s.repo.UpdateBalance(ctx, accountID, amount)
}

// Вспомогательные функции

func generateAccountNumber() string {
	// Генерация 20-значного номера счета
	number := make([]byte, 20)
	for i := 0; i < 20; i++ {
		number[i] = byte(rand.Intn(10)) + '0'
	}
	return string(number)
} 