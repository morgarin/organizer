package models

// Интерфейс для инкапсуляции и возможости юзать тесты с моками
type UserRepositoryInterface interface {
	Create(user *User, password string) error
	UserPassword(name string) (*UserAuthorization, error)
	FindByName(name string) (*User, error)
	FindByID(id int) (*User, error)
}

// Интерфейс для хендлеров
type AuthServiceInterface interface {
	Register(password, name string) error
	Login(name, password string) (string, error)
}
