package server_test

import (
	"context"
	"net/http"
	"testing"
	"time"

	"organizer/internal/server"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestNew проверяет, что сервер создаётся с правильными настройками
func TestNew(t *testing.T) {
	cfg := server.DefaultConfig()
	cfg.Port = "9090"
	cfg.ReadTimeout = 15 * time.Second
	cfg.WriteTimeout = 20 * time.Second
	cfg.IdleTimeout = 60 * time.Second

	handler := http.NewServeMux()
	srv := server.New(handler, cfg)

	// Проверяем через геттеры
	assert.Equal(t, ":9090", srv.Addr())
	assert.Equal(t, 15*time.Second, srv.ReadTimeout())
	assert.Equal(t, 20*time.Second, srv.WriteTimeout())
	assert.Equal(t, 60*time.Second, srv.IdleTimeout())
}

// TestServerRun_ContextCancellation проверяет, что сервер корректно завершается при отмене контекста
func TestServerRun_ContextCancellation(t *testing.T) {
	// Создаём простой хендлер для проверки работоспособности
	mux := http.NewServeMux()
	mux.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("pong"))
	})

	// Используем фиксированный порт (в тестах он обычно свободен)
	const testPort = "8082"
	cfg := server.DefaultConfig()
	cfg.Port = testPort
	srv := server.New(mux, cfg)

	// Контекст с отменой
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Запускаем сервер в горутине
	errCh := make(chan error, 1)
	go func() {
		errCh <- srv.Run(ctx)
	}()

	// Даём серверу время запуститься
	time.Sleep(100 * time.Millisecond)

	// Проверяем, что сервер отвечает на запросы
	resp, err := http.Get("http://localhost:" + testPort + "/ping")
	require.NoError(t, err, "сервер должен отвечать")
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	_ = resp.Body.Close()

	// Отменяем контекст – запускаем graceful shutdown
	cancel()

	// Ждём завершения сервера с таймаутом
	select {
	case err := <-errCh:
		assert.NoError(t, err, "сервер должен завершиться без ошибок")
	case <-time.After(5 * time.Second):
		t.Fatal("сервер не остановился в течение 5 секунд")
	}

	// Проверяем, что сервер больше не отвечает
	_, err = http.Get("http://localhost:" + testPort + "/ping")
	assert.Error(t, err, "сервер должен быть недоступен после остановки")
}

// TestServerRun_ShutdownTimeout потом дописать
