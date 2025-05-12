package repository

import (
	"context"
	"database/sql"
	"time"
	"bank-api/internal/models"
)

type PostgresCardRepository struct {
	db *sql.DB
}

func NewCardRepository(db *sql.DB) CardRepository {
	return &PostgresCardRepository{db: db}
}

func (r *PostgresCardRepository) Create(ctx context.Context, card *models.Card) error {
	query := `
		INSERT INTO cards (account_id, number, encrypted_cvv, expiry_date, cardholder_name, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id`

	return r.db.QueryRowContext(ctx, query,
		card.AccountID,
		card.Number,
		card.EncryptedCVV,
		card.ExpiryDate,
		card.CardholderName,
		card.IsActive,
		time.Now(),
		time.Now(),
	).Scan(&card.ID)
}

func (r *PostgresCardRepository) GetByID(ctx context.Context, id int64) (*models.Card, error) {
	card := &models.Card{}
	query := `
		SELECT id, account_id, number, encrypted_cvv, expiry_date, cardholder_name, is_active, created_at, updated_at
		FROM cards
		WHERE id = $1`

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&card.ID,
		&card.AccountID,
		&card.Number,
		&card.EncryptedCVV,
		&card.ExpiryDate,
		&card.CardholderName,
		&card.IsActive,
		&card.CreatedAt,
		&card.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return card, nil
}

func (r *PostgresCardRepository) GetByAccountID(ctx context.Context, accountID int64) ([]*models.Card, error) {
	query := `
		SELECT id, account_id, number, encrypted_cvv, expiry_date, cardholder_name, is_active, created_at, updated_at
		FROM cards
		WHERE account_id = $1`

	rows, err := r.db.QueryContext(ctx, query, accountID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cards []*models.Card
	for rows.Next() {
		card := &models.Card{}
		err := rows.Scan(
			&card.ID,
			&card.AccountID,
			&card.Number,
			&card.EncryptedCVV,
			&card.ExpiryDate,
			&card.CardholderName,
			&card.IsActive,
			&card.CreatedAt,
			&card.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		cards = append(cards, card)
	}

	return cards, nil
}

func (r *PostgresCardRepository) Update(ctx context.Context, card *models.Card) error {
	query := `
		UPDATE cards
		SET is_active = $1, updated_at = $2
		WHERE id = $3`

	_, err := r.db.ExecContext(ctx, query, card.IsActive, time.Now(), card.ID)
	return err
}

func (r *PostgresCardRepository) Delete(ctx context.Context, id int64) error {
	query := `
		DELETE FROM cards
		WHERE id = $1`

	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *PostgresCardRepository) Deactivate(ctx context.Context, id int64) error {
	query := `
		UPDATE cards
		SET is_active = false, updated_at = $1
		WHERE id = $2`

	_, err := r.db.ExecContext(ctx, query, time.Now(), id)
	return err
} 