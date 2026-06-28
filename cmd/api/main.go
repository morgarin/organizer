package main

import (
	"context"
	"fmt"
	"net/http"
	"organizer/internal/handlers"
	"organizer/internal/repository"
	"organizer/internal/server"
	"organizer/internal/service"
	"organizer/pkg/db"
	"organizer/pkg/logger"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		fmt.Println("No .env file found:", err)
		return
	}
	logger.Init(os.Getenv("LOG_LVL"))
	logger.Info("Логгер запущен")

	// Подключаемся к БД
	dbConfig := db.LoadConfigFromEnv()
	database, err := db.NewConnection(dbConfig)
	if err != nil {
		logger.Error("Fatal: Failed to connect to DB", "error", err)
		return
	}
	defer database.Close()

	logger.Info("Database connected")

	// Инициализируем слои
	userRepo := repository.NewUserRepository(database)
	authService := service.NewAuthService(userRepo)
	authHandler := handlers.NewAuthHandler(authService)

	// Роутер
	mux := http.NewServeMux()
	mux.HandleFunc("POST /api/register", authHandler.Register)
	mux.HandleFunc("POST /api/login", authHandler.Login)

	// Настройка сервера
	cfg := server.DefaultConfig()
	cfg.Port = os.Getenv("PORT")
	if cfg.Port == "" {
		cfg.Port = "8080"
	}

	srv := server.New(mux, cfg)

	ctx := context.Background()
	if err := srv.Run(ctx); err != nil {
		logger.Error("Server error", "error", err)
	}

	// Запуск (блокируется до сигнала)
	if err := srv.Run(ctx); err != nil {
		logger.Error("Server error", "fatal", err)
	}
}
