package repository

import (
	"context"
	"bank-api/internal/models"
	"database/sql"
	"time"
)

type PostgresCreditRepository struct {
	db *sql.DB
}

func NewCreditRepository(db *sql.DB) CreditRepository {
	return &PostgresCreditRepository{db: db}
}

func (r *PostgresCreditRepository) Create(ctx context.Context, credit *models.Credit) error {
	query := `
		INSERT INTO credits (user_id, account_id, amount, interest_rate, term_months, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id`

	return r.db.QueryRowContext(ctx, query,
		credit.UserID,
		credit.AccountID,
		credit.Amount,
		credit.InterestRate,
		credit.TermMonths,
		credit.Status,
		time.Now(),
		time.Now(),
	).Scan(&credit.ID)
}

func (r *PostgresCreditRepository) GetByID(ctx context.Context, id int64) (*models.Credit, error) {
	credit := &models.Credit{}
	query := `
		SELECT id, user_id, account_id, amount, interest_rate, term_months, status, created_at, updated_at
		FROM credits
		WHERE id = $1`

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&credit.ID,
		&credit.UserID,
		&credit.AccountID,
		&credit.Amount,
		&credit.InterestRate,
		&credit.TermMonths,
		&credit.Status,
		&credit.CreatedAt,
		&credit.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return credit, nil
}

func (r *PostgresCreditRepository) GetByUserID(ctx context.Context, userID int64) ([]*models.Credit, error) {
	query := `
		SELECT id, user_id, account_id, amount, interest_rate, term_months, status, created_at, updated_at
		FROM credits
		WHERE user_id = $1`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var credits []*models.Credit
	for rows.Next() {
		credit := &models.Credit{}
		err := rows.Scan(
			&credit.ID,
			&credit.UserID,
			&credit.AccountID,
			&credit.Amount,
			&credit.InterestRate,
			&credit.TermMonths,
			&credit.Status,
			&credit.CreatedAt,
			&credit.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		credits = append(credits, credit)
	}

	return credits, nil
}

func (r *PostgresCreditRepository) Update(ctx context.Context, credit *models.Credit) error {
	query := `
		UPDATE credits
		SET status = $1, updated_at = $2
		WHERE id = $3`

	_, err := r.db.ExecContext(ctx, query, credit.Status, time.Now(), credit.ID)
	return err
}

func (r *PostgresCreditRepository) Delete(ctx context.Context, id int64) error {
	query := `
		DELETE FROM credits
		WHERE id = $1`

	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *PostgresCreditRepository) CreatePaymentSchedule(ctx context.Context, schedule *models.PaymentSchedule) error {
	query := `
		INSERT INTO payment_schedules (credit_id, payment_number, amount, principal, interest, due_date, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id`

	return r.db.QueryRowContext(ctx, query,
		schedule.CreditID,
		schedule.PaymentNumber,
		schedule.Amount,
		schedule.Principal,
		schedule.Interest,
		schedule.DueDate,
		schedule.Status,
		time.Now(),
		time.Now(),
	).Scan(&schedule.ID)
}

func (r *PostgresCreditRepository) GetPaymentSchedule(ctx context.Context, creditID int64) ([]*models.PaymentSchedule, error) {
	query := `
		SELECT id, credit_id, payment_number, amount, principal, interest, due_date, status, created_at, updated_at
		FROM payment_schedules
		WHERE credit_id = $1
		ORDER BY payment_number`

	rows, err := r.db.QueryContext(ctx, query, creditID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var schedules []*models.PaymentSchedule
	for rows.Next() {
		schedule := &models.PaymentSchedule{}
		err := rows.Scan(
			&schedule.ID,
			&schedule.CreditID,
			&schedule.PaymentNumber,
			&schedule.Amount,
			&schedule.Principal,
			&schedule.Interest,
			&schedule.DueDate,
			&schedule.Status,
			&schedule.CreatedAt,
			&schedule.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		schedules = append(schedules, schedule)
	}

	return schedules, nil
}

func (r *PostgresCreditRepository) UpdatePaymentStatus(ctx context.Context, paymentID int64, status string) error {
	query := `
		UPDATE payment_schedules
		SET status = $1, updated_at = $2
		WHERE id = $3`

	_, err := r.db.ExecContext(ctx, query, status, time.Now(), paymentID)
	return err
} 