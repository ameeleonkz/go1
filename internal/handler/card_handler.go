package handler

import (
	"encoding/json"
	"net/http"
	"bank-api/internal/models"
	"bank-api/internal/service"
	"strconv"
	"github.com/gorilla/mux"
)

type CardHandler struct {
	cardService *service.CardService
}

func NewCardHandler(cardService *service.CardService) *CardHandler {
	return &CardHandler{
		cardService: cardService,
	}
}

func (h *CardHandler) Create(w http.ResponseWriter, r *http.Request) {
	var input models.CardCreate
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	card, err := h.cardService.Create(r.Context(), input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Преобразуем в CardResponse для безопасности
	response := models.CardResponse{
		ID:            card.ID,
		AccountID:     card.AccountID,
		Number:        card.Number,
		ExpiryDate:    card.ExpiryDate,
		CardholderName: card.CardholderName,
		IsActive:      card.IsActive,
		CreatedAt:     card.CreatedAt,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *CardHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid card ID", http.StatusBadRequest)
		return
	}

	card, err := h.cardService.GetByID(r.Context(), id)
	if err != nil {
		http.Error(w, "Card not found", http.StatusNotFound)
		return
	}

	// Преобразуем в CardResponse для безопасности
	response := models.CardResponse{
		ID:            card.ID,
		AccountID:     card.AccountID,
		Number:        card.Number,
		ExpiryDate:    card.ExpiryDate,
		CardholderName: card.CardholderName,
		IsActive:      card.IsActive,
		CreatedAt:     card.CreatedAt,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *CardHandler) GetByAccountID(w http.ResponseWriter, r *http.Request) {
	accountIDStr := r.URL.Query().Get("account_id")
	accountID, err := strconv.ParseInt(accountIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid account ID", http.StatusBadRequest)
		return
	}

	cards, err := h.cardService.GetByAccountID(r.Context(), accountID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Преобразуем в CardResponse для безопасности
	responses := make([]models.CardResponse, len(cards))
	for i, card := range cards {
		responses[i] = models.CardResponse{
			ID:            card.ID,
			AccountID:     card.AccountID,
			Number:        card.Number,
			ExpiryDate:    card.ExpiryDate,
			CardholderName: card.CardholderName,
			IsActive:      card.IsActive,
			CreatedAt:     card.CreatedAt,
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(responses)
}

func (h *CardHandler) Deactivate(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid card ID", http.StatusBadRequest)
		return
	}

	if err := h.cardService.Deactivate(r.Context(), id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
} 