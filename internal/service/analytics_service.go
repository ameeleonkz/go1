package service

import (
	"context"
	"bank-api/internal/models"
	"bank-api/internal/repository"
	"time"
)

type AnalyticsService struct {
	accountRepo repository.AccountRepository
	creditRepo  repository.CreditRepository
}

func NewAnalyticsService(accountRepo repository.AccountRepository, creditRepo repository.CreditRepository) *AnalyticsService {
	return &AnalyticsService{
		accountRepo: accountRepo,
		creditRepo:  creditRepo,
	}
}

func (s *AnalyticsService) GetAnalytics(ctx context.Context, userID int64, forecastDays int) (*models.Analytics, error) {
	// Получаем все счета пользователя
	accounts, err := s.accountRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Получаем все кредиты пользователя
	credits, err := s.creditRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Получаем статистику за текущий месяц
	monthlyStats, err := s.getMonthlyStats(ctx, accounts)
	if err != nil {
		return nil, err
	}

	// Получаем кредитную нагрузку
	creditLoad, err := s.getCreditLoad(ctx, credits, monthlyStats.TotalIncome)
	if err != nil {
		return nil, err
	}

	// Получаем прогноз баланса
	balanceForecast, err := s.getBalanceForecast(ctx, accounts, credits, forecastDays)
	if err != nil {
		return nil, err
	}

	return &models.Analytics{
		MonthlyStats:    *monthlyStats,
		CreditLoad:      *creditLoad,
		BalanceForecast: *balanceForecast,
	}, nil
}

func (s *AnalyticsService) getMonthlyStats(ctx context.Context, accounts []*models.Account) (*models.MonthlyStats, error) {
	now := time.Now()
	startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	endOfMonth := startOfMonth.AddDate(0, 1, 0).Add(-time.Second)

	var totalIncome, totalExpenses float64
	var transactionsCount int
	categoryStats := make(map[string]*models.CategoryStats)

	for _, account := range accounts {
		transactions, err := s.accountRepo.GetTransactions(ctx, account.ID)
		if err != nil {
			return nil, err
		}

		for _, t := range transactions {
			if t.CreatedAt.Before(startOfMonth) || t.CreatedAt.After(endOfMonth) {
				continue
			}

			transactionsCount++
			if t.FromAccountID == account.ID {
				totalExpenses += t.Amount
				updateCategoryStats(categoryStats, "expense", t.Amount)
			} else {
				totalIncome += t.Amount
				updateCategoryStats(categoryStats, "income", t.Amount)
			}
		}
	}

	// Преобразуем map категорий в слайс
	topCategories := make([]models.CategoryStats, 0, len(categoryStats))
	for _, stats := range categoryStats {
		topCategories = append(topCategories, *stats)
	}

	return &models.MonthlyStats{
		Month:         startOfMonth,
		TotalIncome:   totalIncome,
		TotalExpenses: totalExpenses,
		NetIncome:     totalIncome - totalExpenses,
		Transactions:  transactionsCount,
		TopCategories: topCategories,
	}, nil
}

func (s *AnalyticsService) getCreditLoad(ctx context.Context, credits []*models.Credit, monthlyIncome float64) (*models.CreditLoad, error) {
	var totalCredits, totalDebt, monthlyPayments float64
	activeCredits := 0

	for _, credit := range credits {
		if credit.Status == "active" {
			activeCredits++
			totalCredits += credit.Amount
			totalDebt += credit.Amount // В реальном приложении нужно учитывать уже выплаченную часть

			// Получаем график платежей
			schedule, err := s.creditRepo.GetPaymentSchedule(ctx, credit.ID)
			if err != nil {
				return nil, err
			}

			// Считаем платежи за текущий месяц
			now := time.Now()
			startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
			endOfMonth := startOfMonth.AddDate(0, 1, 0).Add(-time.Second)

			for _, payment := range schedule {
				if payment.DueDate.After(startOfMonth) && payment.DueDate.Before(endOfMonth) {
					monthlyPayments += payment.Amount
				}
			}
		}
	}

	creditUtilization := 0.0
	if totalCredits > 0 {
		creditUtilization = (totalDebt / totalCredits) * 100
	}

	paymentToIncome := 0.0
	if monthlyIncome > 0 {
		paymentToIncome = (monthlyPayments / monthlyIncome) * 100
	}

	return &models.CreditLoad{
		TotalCredits:      totalCredits,
		ActiveCredits:     activeCredits,
		MonthlyPayments:   monthlyPayments,
		TotalDebt:         totalDebt,
		CreditUtilization: creditUtilization,
		PaymentToIncome:   paymentToIncome,
	}, nil
}

func (s *AnalyticsService) getBalanceForecast(ctx context.Context, accounts []*models.Account, credits []*models.Credit, days int) (*models.BalanceForecast, error) {
	now := time.Now()
	startDate := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	endDate := startDate.AddDate(0, 0, days)

	// Считаем начальный баланс
	var initialBalance float64
	for _, account := range accounts {
		initialBalance += account.Balance
	}

	// Создаем прогноз по дням
	forecastDays := make([]models.ForecastDay, days)
	currentBalance := initialBalance

	for i := 0; i < days; i++ {
		currentDate := startDate.AddDate(0, 0, i)
		var plannedIncome, plannedExpenses float64

		// Считаем запланированные платежи по кредитам
		for _, credit := range credits {
			if credit.Status != "active" {
				continue
			}

			schedule, err := s.creditRepo.GetPaymentSchedule(ctx, credit.ID)
			if err != nil {
				return nil, err
			}

			for _, payment := range schedule {
				if payment.DueDate.Equal(currentDate) {
					plannedExpenses += payment.Amount
				}
			}
		}

		// В реальном приложении здесь можно добавить учет других запланированных доходов/расходов

		currentBalance += plannedIncome - plannedExpenses

		forecastDays[i] = models.ForecastDay{
			Date:            currentDate,
			ExpectedBalance: currentBalance,
			PlannedIncome:   plannedIncome,
			PlannedExpenses: plannedExpenses,
		}
	}

	return &models.BalanceForecast{
		StartDate:      startDate,
		EndDate:        endDate,
		InitialBalance: initialBalance,
		ForecastDays:   forecastDays,
	}, nil
}

func updateCategoryStats(stats map[string]*models.CategoryStats, category string, amount float64) {
	if _, exists := stats[category]; !exists {
		stats[category] = &models.CategoryStats{
			Category: category,
			Amount:   0,
			Count:    0,
		}
	}
	stats[category].Amount += amount
	stats[category].Count++
} 