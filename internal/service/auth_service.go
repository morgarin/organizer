package service

import (
	"errors"
	"organizer/internal/models"

	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	repo models.UserRepositoryInterface
}

func NewAuthService(repo models.UserRepositoryInterface) *AuthService {
	return &AuthService{repo: repo}
}

// Register проверяет на существование пользователя и, если все ок создает его
func (s *AuthService) Register(password, name string) error {
	// Проверяем, существует ли уже пользователь
	existing, _ := s.repo.FindByName(name)
	if existing != nil {
		return errors.New("Пользователь с этим именем уже существует")
	}

	// Хэшируем пароль
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user := &models.User{
		Name: name,
	}
	err = s.repo.Create(user, string(hash))
	return err
}

// Логин возвращает имя юзера после авторизации
func (s *AuthService) Login(name, password string) (string, error) {
	user, err := s.repo.UserPassword(name)
	if err != nil {
		return "", err
	}
	if user == nil {
		return "", errors.New("invalid credentials")
	}

	// Сравниваем хэш пароля
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", errors.New("invalid credentials")
	}
	return user.Name, nil
}
