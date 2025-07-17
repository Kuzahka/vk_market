package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"

	"vk/internal/adapter/handler"
	_ "vk/internal/adapter/repository"
	"vk/internal/infrastructure/postgres"
	"vk/internal/usecase"
)

func main() {
	// Загрузка переменных окружения из .env файла
	if err := godotenv.Load(); err != nil {
		log.Println("Файл .env не найден, предполагается, что переменные окружения установлены.")
	}

	// Получаем порт из переменных окружения, по умолчанию 8080
	port := os.Getenv("PORT")

	tokenSecretKey := os.Getenv("JWT_SECRET_KEY")
	if tokenSecretKey == "" {
		log.Fatal("Переменная окружения JWT_SECRET_KEY не установлена.")
	}

	// Устанавливаем время жизни токена
	tokenExpiration := 24 * time.Hour

	// Инициализация базы данных PostgreSQL
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("Переменная окружения DATABASE_URL не установлена.")
	}
	db, err := postgres.NewPostgresDB(dbURL)
	if err != nil {
		log.Fatalf("Не удалось подключиться к базе данных: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Printf("Ошибка при закрытии соединения с базой данных: %v", err)
		}
	}()
	log.Println("Успешно подключено к базе данных PostgreSQL.")

	// Инициализация репозиториев
	userRepo := postgres.NewPGUserRepository(db)
	adRepo := postgres.NewPGAdRepository(db)

	// Инициализация Use Cases
	authUseCase := usecase.NewAuthUseCase(userRepo, tokenSecretKey, tokenExpiration) // Передаем tokenSecretKey
	adUseCase := usecase.NewAdUseCase(adRepo)

	// Инициализация HTTP-обработчиков
	authHandler := handler.NewAuthHandler(authUseCase)
	adHandler := handler.NewAdHandler(adUseCase)

	// Настройка маршрутизатора
	router := http.NewServeMux()
	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "OK")
	})

	// Маршруты для аутентификации и регистрации
	router.HandleFunc("POST /auth/register", authHandler.RegisterUser)
	router.HandleFunc("POST /auth/login", authHandler.LoginUser)

	// Маршруты для объявлений
	router.Handle("POST /ads", handler.AuthMiddleware(tokenSecretKey, authUseCase, http.HandlerFunc(adHandler.CreateAd))) // Передаем tokenSecretKey
	router.HandleFunc("GET /ads", adHandler.GetAdsFeed)                                                                   // Лента объявлений не требует авторизации, но может использовать userID из контекста

	server := &http.Server{
		Addr:         ":" + port,
		Handler:      router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	log.Printf("Сервер слушает на порту %s...", port)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Не удалось слушать на порту %s: %v\n", port, err)
	}
}
