package service

import (
	"context"
	"bank-api/internal/models"
	"bank-api/internal/repository"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"math/rand"
	"time"
)

type CardService struct {
	repo repository.CardRepository
}

func NewCardService(repo repository.CardRepository) *CardService {
	return &CardService{repo: repo}
}

func (s *CardService) Create(ctx context.Context, input models.CardCreate) (*models.Card, error) {
	// Генерация номера карты по алгоритму Луна
	cardNumber := generateCardNumber()
	
	// Генерация CVV
	cvv := fmt.Sprintf("%03d", rand.Intn(1000))
	
	// Хеширование CVV
	hashedCVV, err := bcrypt.GenerateFromPassword([]byte(cvv), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// Генерация срока действия (текущая дата + 4 года)
	expiryDate := time.Now().AddDate(4, 0, 0).Format("01/06")

	card := &models.Card{
		AccountID:      input.AccountID,
		Number:         cardNumber,
		EncryptedCVV:   string(hashedCVV),
		ExpiryDate:     expiryDate,
		CardholderName: input.CardholderName,
		IsActive:       true,
	}

	if err := s.repo.Create(ctx, card); err != nil {
		return nil, err
	}

	return card, nil
}

func (s *CardService) GetByID(ctx context.Context, id int64) (*models.Card, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *CardService) GetByAccountID(ctx context.Context, accountID int64) ([]*models.Card, error) {
	return s.repo.GetByAccountID(ctx, accountID)
}

func (s *CardService) Deactivate(ctx context.Context, id int64) error {
	return s.repo.Deactivate(ctx, id)
}

// Вспомогательные функции

func generateCardNumber() string {
	// Генерация 16-значного номера карты
	number := make([]byte, 16)
	for i := 0; i < 16; i++ {
		number[i] = byte(rand.Intn(10)) + '0'
	}

	// Применение алгоритма Луна
	sum := 0
	for i := len(number) - 1; i >= 0; i-- {
		digit := int(number[i] - '0')
		if (len(number)-i)%2 == 0 {
			digit *= 2
			if digit > 9 {
				digit -= 9
			}
		}
		sum += digit
	}

	// Корректировка последней цифры
	checkDigit := (10 - (sum % 10)) % 10
	number[15] = byte(checkDigit) + '0'

	return string(number)
}

func computeHMAC(data string, secret []byte) string {
	h := hmac.New(sha256.New, secret)
	h.Write([]byte(data))
	return hex.EncodeToString(h.Sum(nil))
} 