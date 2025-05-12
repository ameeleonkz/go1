package repository

import (
	"context"
	"bank-api/internal/models"
	"database/sql"
	"time"
)

type PostgresAccountRepository struct {
	db *sql.DB
}

func NewAccountRepository(db *sql.DB) AccountRepository {
	return &PostgresAccountRepository{db: db}
}

func (r *PostgresAccountRepository) Create(ctx context.Context, account *models.Account) error {
	query := `
		INSERT INTO accounts (user_id, number, balance, currency, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id`

	return r.db.QueryRowContext(ctx, query,
		account.UserID,
		account.Number,
		account.Balance,
		account.Currency,
		time.Now(),
		time.Now(),
	).Scan(&account.ID)
}

func (r *PostgresAccountRepository) GetByID(ctx context.Context, id int64) (*models.Account, error) {
	account := &models.Account{}
	query := `
		SELECT id, user_id, number, balance, currency, created_at, updated_at
		FROM accounts
		WHERE id = $1`

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&account.ID,
		&account.UserID,
		&account.Number,
		&account.Balance,
		&account.Currency,
		&account.CreatedAt,
		&account.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return account, nil
}

func (r *PostgresAccountRepository) GetByUserID(ctx context.Context, userID int64) ([]*models.Account, error) {
	query := `
		SELECT id, user_id, number, balance, currency, created_at, updated_at
		FROM accounts
		WHERE user_id = $1`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var accounts []*models.Account
	for rows.Next() {
		account := &models.Account{}
		err := rows.Scan(
			&account.ID,
			&account.UserID,
			&account.Number,
			&account.Balance,
			&account.Currency,
			&account.CreatedAt,
			&account.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		accounts = append(accounts, account)
	}

	return accounts, nil
}

func (r *PostgresAccountRepository) Update(ctx context.Context, account *models.Account) error {
	query := `
		UPDATE accounts
		SET balance = $1, updated_at = $2
		WHERE id = $3`

	_, err := r.db.ExecContext(ctx, query, account.Balance, time.Now(), account.ID)
	return err
}

func (r *PostgresAccountRepository) Delete(ctx context.Context, id int64) error {
	query := `
		DELETE FROM accounts
		WHERE id = $1`

	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *PostgresAccountRepository) CreateTransaction(ctx context.Context, transaction *models.Transaction) error {
	query := `
		INSERT INTO transactions (from_account_id, to_account_id, amount, type, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id`

	return r.db.QueryRowContext(ctx, query,
		transaction.FromAccountID,
		transaction.ToAccountID,
		transaction.Amount,
		transaction.Type,
		transaction.Status,
		time.Now(),
		time.Now(),
	).Scan(&transaction.ID)
}

func (r *PostgresAccountRepository) GetTransactions(ctx context.Context, accountID int64) ([]*models.Transaction, error) {
	query := `
		SELECT id, from_account_id, to_account_id, amount, type, status, created_at, updated_at
		FROM transactions
		WHERE from_account_id = $1 OR to_account_id = $1
		ORDER BY created_at DESC`

	rows, err := r.db.QueryContext(ctx, query, accountID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []*models.Transaction
	for rows.Next() {
		transaction := &models.Transaction{}
		err := rows.Scan(
			&transaction.ID,
			&transaction.FromAccountID,
			&transaction.ToAccountID,
			&transaction.Amount,
			&transaction.Type,
			&transaction.Status,
			&transaction.CreatedAt,
			&transaction.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		transactions = append(transactions, transaction)
	}

	return transactions, nil
}

func (r *PostgresAccountRepository) UpdateBalance(ctx context.Context, accountID int64, amount float64) error {
	query := `
		UPDATE accounts
		SET balance = balance + $1, updated_at = $2
		WHERE id = $3`

	_, err := r.db.ExecContext(ctx, query, amount, time.Now(), accountID)
	return err
} 