package handler

import (
	"encoding/json"
	"net/http"
	"bank-api/internal/models"
	"bank-api/internal/service"
	"strconv"
	"github.com/gorilla/mux"
)

type CreditHandler struct {
	creditService *service.CreditService
}

func NewCreditHandler(creditService *service.CreditService) *CreditHandler {
	return &CreditHandler{
		creditService: creditService,
	}
}

func (h *CreditHandler) Create(w http.ResponseWriter, r *http.Request) {
	// Получаем ID пользователя из контекста
	userIDStr := r.Context().Value("userID").(string)
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	var input models.CreditCreate
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Устанавливаем ID пользователя из контекста
	input.UserID = userID

	credit, err := h.creditService.Create(r.Context(), input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(credit)
}

func (h *CreditHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid credit ID", http.StatusBadRequest)
		return
	}

	credit, err := h.creditService.GetByID(r.Context(), id)
	if err != nil {
		http.Error(w, "Credit not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(credit)
}

func (h *CreditHandler) GetByUserID(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(string)
	id, err := strconv.ParseInt(userID, 10, 64)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	credits, err := h.creditService.GetByUserID(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(credits)
}

func (h *CreditHandler) GetPaymentSchedule(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	creditIDStr := vars["id"]
	creditID, err := strconv.ParseInt(creditIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid credit ID", http.StatusBadRequest)
		return
	}

	schedule, err := h.creditService.GetPaymentSchedule(r.Context(), creditID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(schedule)
}

func (h *CreditHandler) ProcessPayment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	creditIDStr := vars["id"]
	paymentIDStr := vars["payment_id"]

	creditID, err := strconv.ParseInt(creditIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid credit ID", http.StatusBadRequest)
		return
	}

	paymentID, err := strconv.ParseInt(paymentIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid payment ID", http.StatusBadRequest)
		return
	}

	if err := h.creditService.ProcessPayment(r.Context(), creditID, paymentID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
} 