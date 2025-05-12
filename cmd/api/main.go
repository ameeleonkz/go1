package main

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"
	"bank-api/internal/handler"
	"bank-api/internal/middleware"
	"bank-api/internal/repository"
	"bank-api/internal/service"
	"bank-api/pkg/centralbank"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	// Инициализация генератора случайных чисел
	rand.Seed(time.Now().UnixNano())

	// Загрузка .env файла
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	// Получение конфигурации
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET is required in .env file")
	}

	// Подключение к базе данных
	dbConnStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_SSL_MODE"),
	)

	db, err := sql.Open("postgres", dbConnStr)
	if err != nil {
		log.Fatal("Error connecting to database:", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatal("Error pinging database:", err)
	}

	// Инициализация клиентов
	centralBankClient := centralbank.NewClient()

	// Инициализация репозиториев
	userRepo := repository.NewUserRepository(db)
	accountRepo := repository.NewAccountRepository(db)
	cardRepo := repository.NewCardRepository(db)
	creditRepo := repository.NewCreditRepository(db)

	// Инициализация сервисов
	userService := service.NewUserService(userRepo)
	accountService := service.NewAccountService(accountRepo)
	cardService := service.NewCardService(cardRepo)
	creditService := service.NewCreditService(creditRepo, centralBankClient, accountService)
	analyticsService := service.NewAnalyticsService(accountRepo, creditRepo)

	// Инициализация middleware
	authMiddleware := middleware.NewAuthMiddleware(jwtSecret)

	// Инициализация обработчиков
	userHandler := handler.NewUserHandler(userService, authMiddleware)
	accountHandler := handler.NewAccountHandler(accountService)
	cardHandler := handler.NewCardHandler(cardService)
	creditHandler := handler.NewCreditHandler(creditService)
	analyticsHandler := handler.NewAnalyticsHandler(analyticsService)

	// Настройка маршрутизации
	r := mux.NewRouter()

	// Публичные маршруты
	r.HandleFunc("/register", userHandler.Register).Methods("POST")
	r.HandleFunc("/login", userHandler.Login).Methods("POST")

	// Защищенные маршруты
	authRouter := r.PathPrefix("/").Subrouter()
	authRouter.Use(authMiddleware.Auth)

	// Маршруты пользователей
	authRouter.HandleFunc("/profile", userHandler.GetProfile).Methods("GET")

	// Маршруты счетов
	authRouter.HandleFunc("/accounts", accountHandler.Create).Methods("POST")
	authRouter.HandleFunc("/accounts", accountHandler.GetByUserID).Methods("GET")
	authRouter.HandleFunc("/accounts/{id}", accountHandler.GetByID).Methods("GET")
	authRouter.HandleFunc("/transfer", accountHandler.Transfer).Methods("POST")
	authRouter.HandleFunc("/transactions", accountHandler.GetTransactions).Methods("GET")

	// Маршруты карт
	authRouter.HandleFunc("/cards", cardHandler.Create).Methods("POST")
	authRouter.HandleFunc("/cards/{id}", cardHandler.GetByID).Methods("GET")
	authRouter.HandleFunc("/cards", cardHandler.GetByAccountID).Methods("GET")
	authRouter.HandleFunc("/cards/{id}/deactivate", cardHandler.Deactivate).Methods("POST")

	// Маршруты кредитов
	authRouter.HandleFunc("/credits", creditHandler.Create).Methods("POST")
	authRouter.HandleFunc("/credits/{id}", creditHandler.GetByID).Methods("GET")
	authRouter.HandleFunc("/credits", creditHandler.GetByUserID).Methods("GET")
	authRouter.HandleFunc("/credits/{id}/schedule", creditHandler.GetPaymentSchedule).Methods("GET")
	authRouter.HandleFunc("/credits/{id}/payments/{payment_id}", creditHandler.ProcessPayment).Methods("POST")

	// Маршруты аналитики
	authRouter.HandleFunc("/analytics", analyticsHandler.GetAnalytics).Methods("GET")

	// Запуск сервера
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
} 