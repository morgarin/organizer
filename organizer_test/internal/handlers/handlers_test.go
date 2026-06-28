package handlers_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"organizer/internal/handlers"
	"testing"

	"github.com/stretchr/testify/assert"
)

// mockAuthService – мок для интерфейса service.AuthServiceInterface
type mockAuthService struct {
	registerFunc func(password, name string) error
	loginFunc    func(name, password string) (string, error)
}

func (m *mockAuthService) Register(password, name string) error {
	if m.registerFunc != nil {
		return m.registerFunc(password, name)
	}
	return nil
}

func (m *mockAuthService) Login(name, password string) (string, error) {
	if m.loginFunc != nil {
		return m.loginFunc(name, password)
	}
	return name, nil
}

func TestAuthHandler_Register(t *testing.T) {
	tests := []struct {
		name           string
		payload        interface{} // чтобы можно было передать невалидный JSON
		mockRegister   func(password, name string) error
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name:    "успешная регистрация",
			payload: map[string]string{"name": "testuser", "password": "pass123"},
			mockRegister: func(password, name string) error {
				return nil
			},
			expectedStatus: http.StatusCreated,
			expectedBody:   map[string]interface{}{"success": true},
		},
		{
			name:    "ошибка регистрации (имя занято)",
			payload: map[string]string{"name": "existing", "password": "pass123"},
			mockRegister: func(password, name string) error {
				return errors.New("Пользователь с этим именем уже существует")
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   map[string]interface{}{"error": "Пользователь с этим именем уже существует"},
		},
		{
			name:           "невалидный JSON",
			payload:        `{invalid}`,
			mockRegister:   nil, // не будет вызван
			expectedStatus: http.StatusBadRequest,
			expectedBody:   nil, // тело содержит текст ошибки
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc := &mockAuthService{
				registerFunc: tt.mockRegister,
			}
			handler := handlers.NewAuthHandler(mockSvc)

			var reqBody []byte
			switch v := tt.payload.(type) {
			case string:
				reqBody = []byte(v) // невалидный JSON
			case map[string]string:
				reqBody, _ = json.Marshal(v)
			default:
				reqBody = []byte{}
			}

			req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewReader(reqBody))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			handler.Register(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectedBody != nil {
				var resp map[string]interface{}
				err := json.NewDecoder(w.Body).Decode(&resp)
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedBody, resp)
			} else {
				// Для невалидного JSON проверяем текст ошибки
				assert.Contains(t, w.Body.String(), "Invalid request")
			}
		})
	}
}

// Тесты для Login аналогично
func TestAuthHandler_Login(t *testing.T) {
	tests := []struct {
		name           string
		payload        interface{} // чтобы можно было передать невалидный JSON
		mockLogin      func(password, name string) (string, error)
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name:    "успешный вход",
			payload: map[string]string{"name": "testuser", "password": "pass123"},
			mockLogin: func(name, password string) (string, error) {
				return name, nil
			},
			expectedStatus: http.StatusOK,
			expectedBody:   map[string]interface{}{"success": true},
		},
		{
			name:    "неверные учетные данные",
			payload: map[string]string{"name": "existing", "password": "wrong"},
			mockLogin: func(name, password string) (string, error) {
				return "", errors.New("invalid credentials")
			},
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   map[string]interface{}{"error": "invalid credentials"},
		},
		{
			name:           "невалидный JSON",
			payload:        `{invalid}`,
			mockLogin:      nil,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc := &mockAuthService{
				loginFunc: tt.mockLogin,
			}
			handler := handlers.NewAuthHandler(mockSvc)

			var reqBody []byte
			switch v := tt.payload.(type) {
			case string:
				reqBody = []byte(v) // невалидный JSON
			case map[string]string:
				reqBody, _ = json.Marshal(v)
			default:
				reqBody = []byte{}
			}

			req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(reqBody))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			handler.Login(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectedBody != nil {
				var resp map[string]interface{}
				err := json.NewDecoder(w.Body).Decode(&resp)
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedBody, resp)
			} else {
				// Для невалидного JSON проверяем текст ошибки
				assert.Contains(t, w.Body.String(), "Invalid request")
			}
		})
	}
}
