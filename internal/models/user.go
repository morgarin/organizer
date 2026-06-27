package models

import "time"

type User struct {
	ID         int       `json:"id"`
	Name       string    `json:"name"`
	TelegramID *int64    `json:"telegram_id,omitempty"`
	CreatedAt  time.Time `json:"created_at"`
}

type UserAuthorization struct {
	Name     string `json:"name"`
	Password string `json:"-"` // нельзя возвращать в JSON
}
