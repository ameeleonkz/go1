package service

import (
	"context"
	"bank-api/internal/models"
	"bank-api/internal/repository"
	"bank-api/pkg/centralbank"
	"errors"
	"math"
	"time"
)

type CreditService struct {
	repo           repository.CreditRepository
	centralBank    *centralbank.Client
	accountService *AccountService
}

func NewCreditService(repo repository.CreditRepository, centralBank *centralbank.Client, accountService *AccountService) *CreditService {
	return &CreditService{
		repo:           repo,
		centralBank:    centralBank,
		accountService: accountService,
	}
}

func (s *CreditService) Create(ctx context.Context, input models.CreditCreate) (*models.Credit, error) {
	// Получение ключевой ставки ЦБ
	baseRate, err := s.centralBank.GetKeyRate()
	if err != nil {
		return nil, err
	}

	// Добавление маржи банка (5%)
	interestRate := baseRate + 5.0

	credit := &models.Credit{
		UserID:       input.UserID,
		AccountID:    input.AccountID,
		Amount:       input.Amount,
		InterestRate: interestRate,
		TermMonths:   input.TermMonths,
		Status:       "active",
	}

	if err := s.repo.Create(ctx, credit); err != nil {
		return nil, err
	}

	// Увеличиваем баланс счета на сумму кредита
	if err := s.accountService.UpdateBalance(ctx, input.AccountID, input.Amount); err != nil {
		// В случае ошибки отменяем создание кредита
		_ = s.repo.Delete(ctx, credit.ID)
		return nil, err
	}

	// Создание графика платежей
	if err := s.createPaymentSchedule(ctx, credit); err != nil {
		// В случае ошибки отменяем создание кредита и возвращаем средства
		_ = s.accountService.UpdateBalance(ctx, input.AccountID, -input.Amount)
		_ = s.repo.Delete(ctx, credit.ID)
		return nil, err
	}

	return credit, nil
}

func (s *CreditService) GetByID(ctx context.Context, id int64) (*models.Credit, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *CreditService) GetByUserID(ctx context.Context, userID int64) ([]*models.Credit, error) {
	return s.repo.GetByUserID(ctx, userID)
}

func (s *CreditService) GetPaymentSchedule(ctx context.Context, creditID int64) ([]*models.PaymentSchedule, error) {
	return s.repo.GetPaymentSchedule(ctx, creditID)
}

func (s *CreditService) ProcessPayment(ctx context.Context, creditID int64, paymentID int64) error {
	schedule, err := s.repo.GetPaymentSchedule(ctx, creditID)
	if err != nil {
		return err
	}

	var payment *models.PaymentSchedule
	for _, p := range schedule {
		if p.ID == paymentID {
			payment = p
			break
		}
	}

	if payment == nil {
		return errors.New("payment not found")
	}

	if payment.Status != "pending" {
		return errors.New("payment already processed")
	}

	// Проверка достаточности средств
	account, err := s.accountService.GetByID(ctx, payment.CreditID)
	if err != nil {
		return err
	}

	if account.Balance < payment.Amount {
		// Начисление штрафа за просрочку
		payment.Amount *= 1.1 // +10% к сумме
		payment.Status = "overdue"
	} else {
		payment.Status = "completed"
	}

	return s.repo.UpdatePaymentStatus(ctx, paymentID, payment.Status)
}

// Вспомогательные функции

func (s *CreditService) createPaymentSchedule(ctx context.Context, credit *models.Credit) error {
	// Расчет аннуитетного платежа
	monthlyRate := credit.InterestRate / 12 / 100
	annuityPayment := credit.Amount * (monthlyRate * math.Pow(1+monthlyRate, float64(credit.TermMonths))) / (math.Pow(1+monthlyRate, float64(credit.TermMonths)) - 1)

	remainingPrincipal := credit.Amount
	dueDate := time.Now().AddDate(0, 1, 0)

	for i := 1; i <= credit.TermMonths; i++ {
		interest := remainingPrincipal * monthlyRate
		principal := annuityPayment - interest
		remainingPrincipal -= principal

		payment := &models.PaymentSchedule{
			CreditID:      credit.ID,
			PaymentNumber: i,
			Amount:        annuityPayment,
			Principal:     principal,
			Interest:      interest,
			DueDate:       dueDate,
			Status:        "pending",
		}

		if err := s.repo.CreatePaymentSchedule(ctx, payment); err != nil {
			return err
		}

		dueDate = dueDate.AddDate(0, 1, 0)
	}

	return nil
} 