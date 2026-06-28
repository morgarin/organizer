package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"organizer/pkg/logger"
)

type Config struct {
	Port         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
}

func DefaultConfig() Config {
	return Config{
		Port:         "8080",
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  30 * time.Second,
	}
}

type Server struct {
	httpServer *http.Server
	config     Config
}

// Геттеры
func (s *Server) Addr() string                { return s.httpServer.Addr }
func (s *Server) ReadTimeout() time.Duration  { return s.httpServer.ReadTimeout }
func (s *Server) WriteTimeout() time.Duration { return s.httpServer.WriteTimeout }
func (s *Server) IdleTimeout() time.Duration  { return s.httpServer.IdleTimeout }

func New(handler http.Handler, config Config) *Server {
	return &Server{
		httpServer: &http.Server{
			Addr:         ":" + config.Port,
			Handler:      handler,
			ReadTimeout:  config.ReadTimeout,
			WriteTimeout: config.WriteTimeout,
			IdleTimeout:  config.IdleTimeout,
		},
		config: config,
	}
}

// Run запускает сервер и ожидает сигнал завершения или отмены контекста
func (s *Server) Run(ctx context.Context) error {
	serverErrors := make(chan error, 1)

	go func() {
		logger.Info("Starting HTTP server", "port", s.config.Port)
		if err := s.httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			serverErrors <- err
		}
	}()

	// Канал для сигналов ОС
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	// Ожидаем завершения
	var shutdownErr error
	select {
	case err := <-serverErrors:
		return fmt.Errorf("server error: %w", err)

	case <-ctx.Done():
		logger.Info("Shutting down due to context cancellation")
		shutdownErr = s.gracefulShutdown(30 * time.Second)

	case sig := <-shutdown:
		logger.Info("Shutdown signal received", "signal", sig)
		shutdownErr = s.gracefulShutdown(30 * time.Second)
	}

	return shutdownErr
}

// gracefulShutdown выполняет мягкое завершение с таймаутом
func (s *Server) gracefulShutdown(timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	logger.Info("Shutting down gracefully...")
	if err := s.httpServer.Shutdown(ctx); err != nil {
		logger.Error("Graceful shutdown failed, forcing close", "error", err)
		if closeErr := s.httpServer.Close(); closeErr != nil {
			return fmt.Errorf("force close error: %w", closeErr)
		}
		return err
	}
	logger.Info("Server stopped gracefully")
	return nil
}

// Shutdown позволяет остановить сервер вручную (например, из тестов)
func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
