package db

import (
	"database/sql"
	"fmt"
	"organizer/pkg/logger"
	"os"
	"time"

	_ "github.com/lib/pq" // драйвер PostgreSQL
)

// Config содержит параметры подключения к БД
type Config struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
}

// NewConnection открывает подключение к Postgree
func NewConnection(cfg Config) (*sql.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName,
	)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		logger.Error("Ошибка подключения к БД")
		return nil, err
	}

	// Настройка пула соединений
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	if err := db.Ping(); err != nil {
		logger.Error("Ошибка соединения с БД")
		return nil, err
	}

	return db, nil
}

// LoadConfigFromEnv читает настройки из переменных окружения
func LoadConfigFromEnv() Config {
	return Config{
		Host:     os.Getenv("HOST"),
		Port:     os.Getenv("DB_PORT"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		DBName:   os.Getenv("DB_NAME"),
	}
}
