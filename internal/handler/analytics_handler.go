package handler

import (
	"encoding/json"
	"net/http"
	"bank-api/internal/service"
	"strconv"
)

type AnalyticsHandler struct {
	analyticsService *service.AnalyticsService
}

func NewAnalyticsHandler(analyticsService *service.AnalyticsService) *AnalyticsHandler {
	return &AnalyticsHandler{
		analyticsService: analyticsService,
	}
}

func (h *AnalyticsHandler) GetAnalytics(w http.ResponseWriter, r *http.Request) {
	// Получаем ID пользователя из контекста
	userIDStr := r.Context().Value("userID").(string)
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	// Получаем количество дней для прогноза из query параметров
	forecastDaysStr := r.URL.Query().Get("forecast_days")
	forecastDays := 30 // значение по умолчанию
	if forecastDaysStr != "" {
		forecastDays, err = strconv.Atoi(forecastDaysStr)
		if err != nil || forecastDays <= 0 {
			http.Error(w, "Invalid forecast days", http.StatusBadRequest)
			return
		}
	}

	analytics, err := h.analyticsService.GetAnalytics(r.Context(), userID, forecastDays)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(analytics)
} 