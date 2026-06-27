package repository

import (
	"database/sql"
	"organizer/internal/models"
	"organizer/pkg/logger"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

// Создает пользователя в дб
func (r *UserRepository) Create(user *models.User, password string) error {
	query := `INSERT INTO users (password, name) VALUES ($1, $2) RETURNING id, created_at`
	return r.db.QueryRow(query, password, user.Name).Scan(&user.ID, &user.CreatedAt)
}

// Возвращает пароль пользователя
func (r *UserRepository) UserPassword(name string) (*models.UserAuthorization, error) {
	user := &models.UserAuthorization{}
	query := `SELECT password, name FROM users WHERE name = $1`
	err := r.db.QueryRow(query, name).Scan(
		&user.Password, &user.Name,
	)
	if err == sql.ErrNoRows {
		logger.Warn("Авторизация: Пользователь по указанному имени не найден")
		return nil, nil
	}
	return user, err
}

// Поиск по имени пользователя
func (r *UserRepository) FindByName(name string) (*models.User, error) {
	user := &models.User{}
	query := `SELECT id, name, telegram_id, created_at FROM users WHERE name = $1`
	err := r.db.QueryRow(query, name).Scan(
		&user.ID, &user.Name, &user.TelegramID, &user.CreatedAt,
	)
	if err == sql.ErrNoRows {
		logger.Warn("Пользователь по указанному имени не найден")
		return nil, nil
	}
	return user, err
}

// Поиск пользователя по id
func (r *UserRepository) FindByID(id int) (*models.User, error) {
	user := &models.User{}
	query := `SELECT id, name, telegram_id, created_at FROM users WHERE id = $1`
	err := r.db.QueryRow(query, id).Scan(
		&user.ID, &user.Name, &user.TelegramID, &user.CreatedAt,
	)
	if err == sql.ErrNoRows {
		logger.Warn("Пользователь по указанному id не найден")
		return nil, nil
	}
	return user, err
}
