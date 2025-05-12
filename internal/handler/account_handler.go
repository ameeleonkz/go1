package handler

import (
	"encoding/json"
	"net/http"
	"bank-api/internal/models"
	"bank-api/internal/service"
	"strconv"
	"github.com/gorilla/mux"
)

type AccountHandler struct {
	accountService *service.AccountService
}

func NewAccountHandler(accountService *service.AccountService) *AccountHandler {
	return &AccountHandler{
		accountService: accountService,
	}
}

func (h *AccountHandler) Create(w http.ResponseWriter, r *http.Request) {
	// Получаем ID пользователя из контекста
	userIDStr := r.Context().Value("userID").(string)
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	var input models.AccountCreate
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Устанавливаем ID пользователя из контекста
	input.UserID = userID

	account, err := h.accountService.Create(r.Context(), input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(account)
}

func (h *AccountHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid account ID", http.StatusBadRequest)
		return
	}

	account, err := h.accountService.GetByID(r.Context(), id)
	if err != nil {
		http.Error(w, "Account not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(account)
}

func (h *AccountHandler) GetByUserID(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(string)
	id, err := strconv.ParseInt(userID, 10, 64)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	accounts, err := h.accountService.GetByUserID(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(accounts)
}

func (h *AccountHandler) Transfer(w http.ResponseWriter, r *http.Request) {
	var input struct {
		FromAccountID int64   `json:"from_account_id"`
		ToAccountID   int64   `json:"to_account_id"`
		Amount        float64 `json:"amount"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.accountService.Transfer(r.Context(), input.FromAccountID, input.ToAccountID, input.Amount); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *AccountHandler) GetTransactions(w http.ResponseWriter, r *http.Request) {
	// Получаем ID пользователя из контекста
	userIDStr := r.Context().Value("userID").(string)
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	// Получаем все счета пользователя
	accounts, err := h.accountService.GetByUserID(r.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Собираем транзакции со всех счетов
	transactionMap := make(map[int64]*models.Transaction)
	for _, account := range accounts {
		transactions, err := h.accountService.GetTransactions(r.Context(), account.ID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		for _, transaction := range transactions {
			transactionMap[transaction.ID] = transaction
		}
	}

	// Преобразуем map обратно в слайс
	allTransactions := make([]*models.Transaction, 0, len(transactionMap))
	for _, transaction := range transactionMap {
		allTransactions = append(allTransactions, transaction)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(allTransactions)
} 