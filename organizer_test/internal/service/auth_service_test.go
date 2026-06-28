package service_test

import (
	"errors"
	"organizer/internal/models"
	"organizer/internal/service"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

// mockUserRepository – реализация интерфейса UserRepositoryInterface для тестов
type mockUserRepository struct {
	// можно хранить данные в map
	users            map[string]*models.User
	passwords        map[string]string
	findByIDFunc     func(id int) (*models.User, error)
	createFunc       func(user *models.User, password string) error
	userPasswordFunc func(name string) (*models.UserAuthorization, error)
}

func (m *mockUserRepository) Create(user *models.User, password string) error {
	if m.createFunc != nil {
		return m.createFunc(user, password)
	}
	if _, exists := m.users[user.Name]; exists {
		return errors.New("user already exists")
	}
	m.users[user.Name] = user
	m.passwords[user.Name] = password
	user.ID = len(m.users)
	return nil
}

func (m *mockUserRepository) UserPassword(name string) (*models.UserAuthorization, error) {
	if m.userPasswordFunc != nil {
		return m.userPasswordFunc(name)
	}

	pass, ok := m.passwords[name]
	if !ok {
		return nil, nil
	}

	return &models.UserAuthorization{
		Name:     name,
		Password: pass,
	}, nil
}

func (m *mockUserRepository) FindByName(name string) (*models.User, error) {
	user, ok := m.users[name]
	if !ok {
		return nil, nil
	}
	return user, nil
}

func (m *mockUserRepository) FindByID(id int) (*models.User, error) {
	if m.findByIDFunc != nil {
		return m.findByIDFunc(id)
	}
	for _, user := range m.users {
		if user.ID == id {
			return user, nil
		}
	}
	return nil, nil
}

// Тесты

func TestAuthService_Register(t *testing.T) {
	mockRepo := &mockUserRepository{
		users:     make(map[string]*models.User),
		passwords: make(map[string]string),
	}

	service := service.NewAuthService(mockRepo)

	// Успешная регистрация
	err := service.Register("password123", "testuser")
	assert.NoError(t, err)

	// Проверяем, что пользователь создался
	user, _ := mockRepo.FindByName("testuser")
	assert.NotNil(t, user)
	assert.Equal(t, "testuser", user.Name)

	// Проверяем, что пароль захэширован
	pass, _ := mockRepo.UserPassword("testuser")
	err = bcrypt.CompareHashAndPassword([]byte(pass.Password), []byte("password123"))
	assert.NoError(t, err)

	// Повторная регистрация с тем же именем
	err = service.Register("anotherpass", "testuser")
	assert.Error(t, err)
	assert.Equal(t, "Пользователь с этим именем уже существует", err.Error())

	// Ошибка хэширования
	mockRepo.createFunc = func(user *models.User, password string) error {
		return errors.New("db error")
	}
	err = service.Register("pass", "newuser")
	assert.Error(t, err)
	assert.Equal(t, "db error", err.Error())
}

func TestAuthService_Login(t *testing.T) {
	mockRepo := &mockUserRepository{
		users:     make(map[string]*models.User),
		passwords: make(map[string]string),
	}
	service := service.NewAuthService(mockRepo)

	// Сначала зарегистрируем пользователя вручную через мок
	hashed, _ := bcrypt.GenerateFromPassword([]byte("correct"), bcrypt.DefaultCost)
	mockRepo.users["testuser"] = &models.User{Name: "testuser", ID: 1}
	mockRepo.passwords["testuser"] = string(hashed)

	// Успешный логин
	name, err := service.Login("testuser", "correct")
	assert.NoError(t, err)
	assert.Equal(t, "testuser", name)

	// Неверный пароль
	name, err = service.Login("testuser", "wrong")
	assert.Error(t, err)
	assert.Equal(t, "invalid credentials", err.Error())
	assert.Empty(t, name)

	// Несуществующий пользователь
	name, err = service.Login("unknown", "pass")
	assert.Error(t, err)
	assert.Equal(t, "invalid credentials", err.Error())
	assert.Empty(t, name)

	// 4. Ошибка репозитория (например, ошибка БД)
	mockRepo.userPasswordFunc = func(name string) (*models.UserAuthorization, error) {
		return nil, errors.New("database error")
	}
}
