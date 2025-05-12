package models

import "time"

// Analytics представляет собой общую структуру аналитики
type Analytics struct {
	MonthlyStats    MonthlyStats    `json:"monthly_stats"`
	CreditLoad      CreditLoad      `json:"credit_load"`
	BalanceForecast BalanceForecast `json:"balance_forecast"`
}

// MonthlyStats представляет статистику доходов/расходов за месяц
type MonthlyStats struct {
	Month           time.Time `json:"month"`
	TotalIncome     float64   `json:"total_income"`
	TotalExpenses   float64   `json:"total_expenses"`
	NetIncome       float64   `json:"net_income"`
	Transactions    int       `json:"transactions_count"`
	TopCategories   []CategoryStats `json:"top_categories"`
}

// CategoryStats представляет статистику по категории транзакций
type CategoryStats struct {
	Category string  `json:"category"`
	Amount   float64 `json:"amount"`
	Count    int     `json:"count"`
}

// CreditLoad представляет аналитику кредитной нагрузки
type CreditLoad struct {
	TotalCredits        float64 `json:"total_credits"`
	ActiveCredits       int     `json:"active_credits"`
	MonthlyPayments     float64 `json:"monthly_payments"`
	TotalDebt           float64 `json:"total_debt"`
	CreditUtilization   float64 `json:"credit_utilization"` // Процент использования кредитного лимита
	PaymentToIncome     float64 `json:"payment_to_income"`  // Отношение платежей к доходу
}

// BalanceForecast представляет прогноз баланса
type BalanceForecast struct {
	StartDate     time.Time         `json:"start_date"`
	EndDate       time.Time         `json:"end_date"`
	InitialBalance float64          `json:"initial_balance"`
	ForecastDays  []ForecastDay     `json:"forecast_days"`
}

// ForecastDay представляет прогноз на конкретный день
type ForecastDay struct {
	Date            time.Time `json:"date"`
	ExpectedBalance float64   `json:"expected_balance"`
	PlannedIncome   float64   `json:"planned_income"`
	PlannedExpenses float64   `json:"planned_expenses"`
} 