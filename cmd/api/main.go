package main

import (
	"fmt"
	"net/http"
	"organizer/internal/handlers"
	"organizer/internal/repository"
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

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	logger.Info("Server starting", "port", port)
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		logger.Error("Server failed", "error", err)
	}
}
