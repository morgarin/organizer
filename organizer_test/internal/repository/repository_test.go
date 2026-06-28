package repository_test

import (
	"database/sql"
	"organizer/internal/models"
	"organizer/internal/repository"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestUserRepository_Create(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock: %v", err)
	}
	defer db.Close()

	repo := repository.NewUserRepository(db)

	user := &models.User{
		Name: "testuser",
	}
	password := "hashedpassword"

	// Ожидаемый запрос
	mock.ExpectQuery(`INSERT INTO users \(password, name\) VALUES \(\$1, \$2\) RETURNING id, created_at`).
		WithArgs(password, user.Name).
		WillReturnRows(sqlmock.NewRows([]string{"id", "created_at"}).AddRow(1, time.Now()))

	err = repo.Create(user, password)
	assert.NoError(t, err)
	assert.Equal(t, 1, user.ID)
	assert.NotZero(t, user.CreatedAt)

	// Проверка, что все ожидания выполнены
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestUserRepository_UserPassword(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock: %v", err)
	}
	defer db.Close()

	repo := repository.NewUserRepository(db)

	name := "testuser"
	expectedPassword := "hashedpass"

	rows := sqlmock.NewRows([]string{"password", "name"}).AddRow(expectedPassword, name)
	mock.ExpectQuery(`SELECT password, name FROM users WHERE name = \$1`).
		WithArgs(name).
		WillReturnRows(rows)

	auth, err := repo.UserPassword(name)
	assert.NoError(t, err)
	assert.NotNil(t, auth)
	assert.Equal(t, expectedPassword, auth.Password)
	assert.Equal(t, name, auth.Name)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	// Тест на отсутствие пользователя
	mock.ExpectQuery(`SELECT password, name FROM users WHERE name = \$1`).
		WithArgs("unknown").
		WillReturnError(sql.ErrNoRows)

	auth, err = repo.UserPassword("unknown")
	assert.NoError(t, err) // функция возвращает nil, nil
	assert.Nil(t, auth)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestUserRepository_FindByName(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock: %v", err)
	}
	defer db.Close()

	repo := repository.NewUserRepository(db)

	name := "testuser"
	now := time.Now()
	telegramID := int64(12345)

	rows := sqlmock.NewRows([]string{"id", "name", "telegram_id", "created_at"}).
		AddRow(1, name, telegramID, now)

	mock.ExpectQuery(`SELECT id, name, telegram_id, created_at FROM users WHERE name = \$1`).
		WithArgs(name).
		WillReturnRows(rows)

	user, err := repo.FindByName(name)
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, 1, user.ID)
	assert.Equal(t, name, user.Name)
	assert.Equal(t, telegramID, *user.TelegramID) // если TelegramID не nil
	assert.Equal(t, now, user.CreatedAt)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	// Тест на отсутствие пользователя
	mock.ExpectQuery(`SELECT id, name, telegram_id, created_at FROM users WHERE name = \$1`).
		WithArgs("unknown").
		WillReturnError(sql.ErrNoRows)

	user, err = repo.FindByName("unknown")
	assert.NoError(t, err)
	assert.Nil(t, user)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestUserRepository_FindByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock: %v", err)
	}
	defer db.Close()

	repo := repository.NewUserRepository(db)

	id := 1
	name := "testuser"
	now := time.Now()
	var telegramID *int64 = nil

	rows := sqlmock.NewRows([]string{"id", "name", "telegram_id", "created_at"}).
		AddRow(id, name, telegramID, now)

	mock.ExpectQuery(`SELECT id, name, telegram_id, created_at FROM users WHERE id = \$1`).
		WithArgs(id).
		WillReturnRows(rows)

	user, err := repo.FindByID(id)
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, id, user.ID)
	assert.Equal(t, name, user.Name)
	assert.Nil(t, user.TelegramID)
	assert.Equal(t, now, user.CreatedAt)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	// Тест на отсутствие пользователя
	mock.ExpectQuery(`SELECT id, name, telegram_id, created_at FROM users WHERE id = \$1`).
		WithArgs(999).
		WillReturnError(sql.ErrNoRows)

	user, err = repo.FindByID(999)
	assert.NoError(t, err)
	assert.Nil(t, user)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
